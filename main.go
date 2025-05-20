package main

import (
	"os"

	examplev1 "github.com/esterry-golang-operator-demo/api/v1"
	"github.com/esterry-golang-operator-demo/controllers"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	// Register built-in Kubernetes types
	_ = clientgoscheme.AddToScheme(scheme)
	// Register our CRD types
	_ = examplev1.AddToScheme(scheme)
}

func main() {
	// Set up structured logging
	ctrl.SetLogger(zap.New())

	// Create the manager with the registered scheme
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
	})
	if err != nil {
		ctrl.Log.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// Wire up the namespace controller from the controllers package
	if err := ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Namespace{}).
		Complete(&controllers.NamespaceReconciler{Client: mgr.GetClient()}); err != nil {
		ctrl.Log.Error(err, "unable to create controller", "controller", "Namespace")
		os.Exit(1)
	}

	// Start the manager
	ctrl.Log.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		ctrl.Log.Error(err, "problem running manager")
		os.Exit(1)
	}
}
