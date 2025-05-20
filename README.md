# Golang Kubernetes Namespace Labeling Operator

## Overview
This PoC is designed to prove core instrumentation that allows an event driven operator based around GoLang, interfacing with namespaces as a relatively easy proof-of-concept for further iteration.
This demonstrates a minimal **Kubernetes Operator** written in Go that:

1. Defines a **CRD** (`NamespaceConfig`) to hold a map of labels.
2. Watches **Namespace** resources and, on each event:
    - **Pass 1**: Creates a default `NamespaceConfig/default` (with `environment=default` and `owner=admin`) if none exists, then requeues.
    - **Pass 2**: Reads that config and applies **all** its `spec.labels` onto the Namespace.
3. Includes:
    - Unit tests (fake client, two-pass reconcile)
    - Docker container build & local smoke-run via Docker Compose
    - End-to-end testing in a local Kind cluster (See USAGE.md for local testing)
    - Raw manifests for CRD, RBAC, and Deployment to run in any Kubernetes (e.g. EKS)

---

This Golang Operator is specifically oriented toward an EKS deployed k8s cluster for deployment, but can be utilized and tested in any local k8s cluster (e.g. Kind, Minikube, Docker Desktop) for development and testing purposes.
