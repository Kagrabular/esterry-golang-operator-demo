# USAGE.md

This script automates the creation (if needed) and retrieval of an Amazon ECR repository URI in your AWS account, and displays commands for building and pushing Docker images.

## Prerequisites

* **AWS CLI v2** installed and configured

    * Ensure you have run `aws configure --profile <profile>` or have `$AWS_PROFILE` environment variable set.
* **AWS IAM permissions**:

    * `ecr:DescribeRepositories`
    * `ecr:CreateRepository`
    * (`ecr:PutImage` and `ecr:InitiateLayerUpload` etc. for pushing)
* **Docker** installed
* **Git** (optional) if using the `git rev-parse --short HEAD` command for tagging

## Installation

1. Clone or download this repository to your local machine.
2. Make the script executable:

   ```bash
   chmod +x scripts/aws-ecr-create-deploy.sh
   ```

## Usage

```bash
./scripts/aws-ecr-create-deploy.sh [--name REPO_NAME] [--region AWS_REGION] [--profile AWS_PROFILE]
```

### Options

| Flag           | Description                                    | Default                          |
| -------------- | ---------------------------------------------- | -------------------------------- |
| `--name`       | Name of the ECR repository to create or lookup | Value of `REPO_NAME` in script   |
| `--region`     | AWS region for the repository                  | Value of `AWS_REGION` in script  |
| `--profile`    | AWS CLI profile to use                         | Value of `AWS_PROFILE` in script |
| `-h`, `--help` | Show this help message                         | N/A                              |

### Example

```bash
# Set the Docker tag to your current Git commit SHA
TAG=$(git rev-parse --short HEAD)

# Run the script (override defaults if desired)
./scripts/aws-ecr-create-deploy.sh --name my-ecr-repo --region us-west-2 --profile my-aws-profile

# Output will include:
# → Using AWS profile: my-aws-profile
# → AWS region:       us-west-2
# → Repository name:  my-ecr-repo
# Repository URI: 123456789012.dkr.ecr.us-west-2.amazonaws.com/my-ecr-repo

# Build and push your Docker image:
docker build -t 123456789012.dkr.ecr.us-west-2.amazonaws.com/my-ecr-repo:$TAG .
docker push 123456789012.dkr.ecr.us-west-2.amazonaws.com/my-ecr-repo:$TAG
```

## Considerations

* **Existing repositories**: If the specified repo already exists, the script will skip creation and simply output the URI.
* **Credentials**: Ensure your AWS profile has valid credentials and the necessary ECR permissions.
* **Region consistency**: Use the same `--region` when building, tagging, and pushing to avoid cross-region mismatches.
* **Tagging strategy**: I recommend using immutable tags (e.g., Git commit SHAs) to track image versions reliably.
* **Encryption & Scanning**: The script enables image scanning on push by default. Ensure your account has the Amazon ECR scan permissions if you rely on that feature.

---
