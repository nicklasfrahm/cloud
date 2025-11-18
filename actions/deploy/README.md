# Deploy

This GitHub Action automates deployment via GitOps by updating the image tag in the configuration repository and pushing the changes.

## Inputs

| Name             | Description                                                   | Required |
| ---------------- | ------------------------------------------------------------- | -------- |
| `service`        | The name of the service to deploy                             | Yes      |
| `environment`    | The target environment for deployment (dev, stg, prd)         | Yes      |
| `tag`            | The image tag to deploy                                       | Yes      |
| `github-ssh-key` | An SSH key with write access to use for GitHub authentication | Yes      |

## Usage

```yaml
name: Release

on:
  push:
    tags:
      - "*"

jobs:
  deploy:
    name: Deploy
    runs-on: ubuntu-latest
    steps:
      - name: Deploy service
        uses: nicklasfrahm-dev/platform/.github/actions/deploy@main
        with:
          service: my-service
          environment: prd
          tag: ${{ github.ref_name }}
          github-ssh-key: ${{ secrets.GITOPS_SSH_PRIVATE_KEY }}
```
