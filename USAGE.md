# USAGE

## Prerequisites

- **Go ≥ 1.24**  
- **Docker**  
- **kubectl** configured for your target cluster  
- **AWS CLI** (for EKS)
- **Docker Compose** (for local smoke-run)
- **A pre-existing ECR repository** (for EKS)
- **Kind** (for local end-to-end testing)  

---

## Quickstart

1. **Clone the repo**  
   ```bash
   git clone https://github.com/esterry-golang-operator-demo.git
   cd esterry-golang-operator-demo


2. **Run unit tests**

   ```bash
   go mod tidy
   go test ./controllers
   ```

---

## Build & Smoke-Run via Docker Compose

1. **Build the Docker image**

   ```bash
   docker build -t local-operator:latest .
   ```

2. **Start the operator locally**

   ```bash
   docker-compose up -d
   docker-compose logs -f operator
   ```

3. **Test against your kubeconfig**

   ```bash
   kubectl create namespace smoke-test
   kubectl get namespace smoke-test --show-labels
   ```

4. **Stop the operator**

   ```bash
   docker-compose down
   ```

---

## End-to-End in Kind

1. **Create a Kind cluster**
   2. This local testing procedure relies on [Kind](https://kind.sigs.k8s.io/) to create a local Kubernetes cluster. If you don't have Kind installed follow the link to install it for your machine. You can deploy and test this in any local k8s cluster following a similar procedure

      ```bash
      kind create cluster --name operator-test
      ```

2. **Build & load your image into Kind**

   ```bash
   docker build -t local-operator:latest .
   kind load docker-image local-operator:latest --name operator-test
   ```

3. **Apply manifests**

   ```bash
   kubectl apply -f config/crd/namespaceconfig.yaml --context kind-operator-test
   kubectl apply -f config/rbac/role.yaml    --context kind-operator-test
   kubectl apply -f config/rbac/rolebinding.yaml --context kind-operator-test
   kubectl apply -f deploy/operator-deployment.yaml --context kind-operator-test
   ```

4. **Verify rollout & functionality**

   ```bash
   kubectl rollout status deployment/namespace-operator --context kind-operator-test
   kubectl create namespace kind-smoke --context kind-operator-test
   kubectl get namespace kind-smoke --show-labels --context kind-operator-test
   ```

5. **Tear down**

   ```bash
   kind delete cluster --name operator-test
   ```

---

## Push & Deploy to EKS

1. **Tag & push to ECR**

   ```bash
   export AWS_ACCOUNT_ID=<your-aws-account-id>
   export AWS_REGION=us-west-2
   export ECR_REPO=golang-operator # or whatever you named your ECR repo
   export IMAGE_URI=${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${ECR_REPO}:latest

   docker tag local-operator:latest $IMAGE_URI
   aws ecr get-login-password --region $AWS_REGION \
     | docker login --username AWS --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com
   docker push $IMAGE_URI
   ```

2. **Switch kubectl to EKS sandbox context**

   ```bash
   kubectl config use-context <your_eks_context>
   ```

3. **Apply CRD & RBAC**

   ```bash
   kubectl apply -f config/crd/namespaceconfig.yaml
   kubectl apply -f config/rbac/role.yaml
   kubectl apply -f config/rbac/rolebinding.yaml
   ```

4. **Deploy the operator**

   ```bash
   # Edit deploy/operator-deployment.yaml to use $IMAGE_URI
   kubectl apply -f deploy/operator-deployment.yaml
   kubectl rollout status deployment/namespace-operator
   ```

5. **Smoke-test**

   ```bash
   kubectl create namespace namespace-test
   kubectl get namespace namespace-test --show-labels
   ```
## Teardown (EKS Sandbox)

When you’re done, you can remove the operator and related resources:

kubectl delete deployment namespace-operator
kubectl delete serviceaccount default
kubectl delete clusterrolebinding namespace-operator
kubectl delete clusterrole namespace-operator
kubectl delete crd namespaceconfigs.example.com

### (Optional) Clean up ECR image
```bash
aws ecr batch-delete-image \
--repository-name $ECR_REPO \
--image-ids imageTag=latest \
--region $AWS_REGION
```
