# ArgoCD Application Actions

This action will sync ArgoCD application.

## Usage

### Example workflow

This example replaces syncs ArgoCD application.

```yaml
name: My Workflow
on: [ push, pull_request ]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Sync ArgoCD Application
        uses: butterfly1924/argocd-actions@master
        with:
          address: "vault.example.com"
          token: ${{ secrets.ARGOCD_TOKEN }}
          insecure: "true"
          appName: "my-example-app"
          imageTag: "dev"
```

### Inputs

| Input | Description|
| --- | --- |
| `address` | ArgoCD server address. |
| `token` | ArgoCD Token. |
| `inseucre` | Whether argocd is insecure |
| `appName` | Application name to sync. |

## Examples

### Sync Application

You can sync ArgoCD application after building an image etc.

```yaml
name: My Workflow
on: [ push, pull_request ]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Sync ArgoCD Application
        uses: butterfly1924/argocd-actions@master
        with:
          address: "vault.example.com"
          token: ${{ secrets.ARGOCD_TOKEN }}
          insecure: "true"
          appName: "my-example-app"
          imageTag: "dev"
```

## Publishing

To publish a new version of this Action we need to update the Docker image tag in `action.yml` and also create a new
release on GitHub.

- Work out the next tag version number.
- Update the Docker image in `action.yml`.
- Create a new release on GitHub with the same tag.
