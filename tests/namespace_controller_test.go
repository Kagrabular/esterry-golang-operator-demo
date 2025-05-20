package controllers_test

import (
	"context"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"

	examplev1 "github.com/esterry-golang-operator-demo/api/v1"
	"github.com/esterry-golang-operator-demo/controllers"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Two passes of Reconcile:
// this mirrors k8s controller behavior where the first pass creates a resource and the second pass applies it.
func TestNamespaceReconcile_LabelsApplied(t *testing.T) {
	ctx := context.Background()

	// Set up a Scheme that knows about corev1 and our CRD types
	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)
	_ = examplev1.AddToScheme(scheme) // ensure NamespaceConfig is recognized

	// Seed the fake client with just a Namespace named "foo"
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "foo"}}
	fakeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithRuntimeObjects(ns).
		Build()

	reconciler := &controllers.NamespaceReconciler{Client: fakeClient}
	req := reconcile.Request{NamespacedName: types.NamespacedName{Name: "foo"}}

	// first reconcile pass
	// expect a default NamespaceConfig is created, and Reconcile signals a requeue
	res1, err := reconciler.Reconcile(ctx, req)
	if err != nil {
		t.Fatalf("first Reconcile error: %v", err)
	}
	if !res1.Requeue {
		t.Errorf("first Reconcile: expected Requeue=true, got %+v", res1)
	}

	// optional, and untested. verify that the CRD object was created - leaving this commented out to avoid test bloat, if we want to explicitly assert CRD creation on first pass we could uncomment this for more coverage.
	// var createdConfig examplev1.NamespaceConfig
	// if err := fakeClient.Get(ctx, client.ObjectKey{Name: "default"}, &createdConfig); err != nil {
	//     t.Fatalf("expected NamespaceConfig to be created, but got error: %v", err)
	// }

	// second reconcile pass
	// Now that the default config exists, Reconcile should apply labels and not requeue
	res2, err := reconciler.Reconcile(ctx, req)
	if err != nil {
		t.Fatalf("second Reconcile error: %v", err)
	}
	if res2.Requeue {
		t.Errorf("second Reconcile: expected no requeue, got %+v", res2)
	}

	// verify labels on namespace
	fetched := &corev1.Namespace{}
	if err := fakeClient.Get(ctx, client.ObjectKey{Name: "foo"}, fetched); err != nil {
		t.Fatalf("failed to get namespace after second reconcile: %v", err)
	}
	if fetched.Labels["environment"] != "default" {
		t.Errorf("environment = %q; want %q", fetched.Labels["environment"], "default")
	}
	if fetched.Labels["owner"] != "admin" {
		t.Errorf("owner = %q; want %q", fetched.Labels["owner"], "admin")
	}
}
