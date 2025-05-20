package controllers

import (
	"context"

	examplev1 "github.com/esterry-golang-operator-demo/api/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NamespaceReconciler watches Namespace objects and applies labels based on NamespaceConfig
// Implements a two-pass reconcile: creates default config on first pass, labels on second.
type NamespaceReconciler struct {
	client.Client
}

// Reconcile is called on every Namespace event
func (r *NamespaceReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	logger := zap.New()

	// Fetch namespace
	var ns corev1.Namespace
	if err := r.Get(ctx, req.NamespacedName, &ns); err != nil {
		logger.Error(err, "Failed to get namespace")
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	// Fetch or create the default NamespaceConfig
	var nsConfig examplev1.NamespaceConfig
	nsConfigKey := client.ObjectKey{Name: "default"}
	if err := r.Get(ctx, nsConfigKey, &nsConfig); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("No NamespaceConfig found, creating default config and requeueing")

			defaultConfig := examplev1.NamespaceConfig{
				TypeMeta:   metav1.TypeMeta{APIVersion: "example.com/v1", Kind: "NamespaceConfig"},
				ObjectMeta: metav1.ObjectMeta{Name: "default"},
				Spec:       examplev1.NamespaceConfigSpec{Labels: map[string]string{"environment": "default", "owner": "admin"}}, // These are example labels, obviously this is a PoC so this should be dynamic. that's more work though.
			}
			if err := r.Create(ctx, &defaultConfig); err != nil {
				logger.Error(err, "Failed to create default NamespaceConfig")
				return reconcile.Result{}, err
			}

			// First pass: requeue to apply labels next
			return reconcile.Result{Requeue: true}, nil
		}
		// Other errors fetching config
		logger.Error(err, "Error fetching NamespaceConfig")
		return reconcile.Result{}, err
	}

	// Second pass: apply labels from NamespaceConfig to the Namespace
	if ns.Labels == nil {
		ns.Labels = make(map[string]string)
	}
	for key, val := range nsConfig.Spec.Labels {
		ns.Labels[key] = val
	}
	if err := r.Update(ctx, &ns); err != nil {
		logger.Error(err, "Failed to update namespace with labels")
		return reconcile.Result{}, err
	}

	logger.Info("Namespace labeled successfully", "name", req.NamespacedName)
	return reconcile.Result{}, nil
}
