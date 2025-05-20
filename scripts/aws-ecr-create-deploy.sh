#!/usr/bin/env bash
# run TAG=$(git rev-parse --short HEAD) after you run this script in the repo you're deploying from to set the tag. Then docker build -t <whateverthiscomesoutas>.dkr.ecr.us-west-2.amazonaws.com/demo-pe-service-java:$TAG .
set -euo pipefail

# defaults
AWS_PROFILE="<aws-profile-name>"
AWS_REGION="<aws-region>"
REPO_NAME="<whatever-you-want-to-name-this-ecr-repo>"

usage() {
  cat <<EOF
Usage: $0 [--name REPO_NAME] [--region AWS_REGION] [--profile AWS_PROFILE]

  --name     Name of the ECR repository to create (default: $REPO_NAME)
  --region   AWS region (default: $AWS_REGION)
  --profile  AWS CLI profile (default: $AWS_PROFILE)
  -h|--help  Show this help
EOF
  exit 1
}

# parse args
while [[ $# -gt 0 ]]; do
  case "$1" in
    --name)    REPO_NAME="$2";   shift 2;;
    --region)  AWS_REGION="$2";  shift 2;;
    --profile) AWS_PROFILE="$2"; shift 2;;
    -h|--help) usage;;
    *) echo "Unknown flag: $1" >&2; usage;;
  esac
done

echo "→ Using AWS profile: $AWS_PROFILE"
echo "→ AWS region:       $AWS_REGION"
echo "→ Repository name:  $REPO_NAME"

# check if repo exists
if aws ecr describe-repositories \
     --repository-names "$REPO_NAME" \
     --region "$AWS_REGION" \
     --profile "$AWS_PROFILE" \
     &>/dev/null; then
  echo "✔ Repository '$REPO_NAME' already exists."
else
  echo "→ Creating repository '$REPO_NAME'..."
  aws ecr create-repository \
    --repository-name "$REPO_NAME" \
    --image-scanning-configuration scanOnPush=true \
    --region "$AWS_REGION" \
    --profile "$AWS_PROFILE" \
    >/dev/null
  echo "✔ Created ECR repository."
fi

# fetch and display the URI
REPO_URI=$(aws ecr describe-repositories \
  --repository-names "$REPO_NAME" \
  --query "repositories[0].repositoryUri" \
  --output text \
  --region "$AWS_REGION" \
  --profile "$AWS_PROFILE")

echo
echo "Repository URI: $REPO_URI"
echo
echo "You can now docker build & tag with:"
echo "  docker build -t $REPO_URI:<TAG> ."
echo "  docker push $REPO_URI:<TAG>"
