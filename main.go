package main

import (
	"context"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

type NamespaceConfigSpec struct {
	Labels map[string]string `json:"labels,omitempty"`
}

type NamespaceConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              NamespaceConfigSpec `json:"spec,omitempty"`
}

type NamespaceReconciler struct {
	client client.Client
}

func (r *NamespaceReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	logger := zap.New()
	logger.Info("Reconciling namespace", "name", req.NamespacedName)

	var ns corev1.Namespace
	if err := r.client.Get(ctx, req.NamespacedName, &ns); err != nil {
		logger.Error(err, "Failed to get namespace")
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	var nsConfig NamespaceConfig
	nsConfigKey := client.ObjectKey{Name: "default"} // Assume a default config
	if err := r.client.Get(ctx, nsConfigKey, &nsConfig); err != nil {
		logger.Info("No NamespaceConfig found, creating default")
		defaultConfig := NamespaceConfig{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "example.com/v1",
				Kind:       "NamespaceConfig",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "default",
			},
			Spec: NamespaceConfigSpec{
				Labels: map[string]string{
					"environment": "default",
					"owner":       "admin",
				},
			},
		}
		if err := r.client.Create(ctx, &defaultConfig); err != nil {
			logger.Error(err, "Failed to create default NamespaceConfig")
			return reconcile.Result{}, err
		}
		logger.Info("Default NamespaceConfig created successfully")
		return reconcile.Result{}, nil // Requeue to ensure processing with the new config
	}

	if ns.Labels == nil {
		ns.Labels = make(map[string]string)
	}
	for key, value := range nsConfig.Spec.Labels {
		ns.Labels[key] = value
	}
	if err := r.client.Update(ctx, &ns); err != nil {
		logger.Error(err, "Failed to update namespace with labels")
		return reconcile.Result{}, err
	}
	logger.Info("Namespace labeled successfully", "name", req.NamespacedName)

	return reconcile.Result{}, nil
}

func main() {
	logger := zap.New()
	mgr, err := manager.New(manager.GetConfigOrDie(), manager.Options{})
	if err != nil {
		logger.Error(err, "Unable to create manager")
		os.Exit(1)
	}

	ctrl, err := controller.New("namespace-controller", mgr, controller.Options{
		Reconciler: &NamespaceReconciler{client: mgr.GetClient()},
	})
	if err != nil {
		logger.Error(err, "Unable to create controller")
		os.Exit(1)
	}

	logger.Info("Starting Kubernetes Operator")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		logger.Error(err, "Problem running manager")
		os.Exit(1)
	}
}
