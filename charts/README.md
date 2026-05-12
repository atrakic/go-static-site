## Kubernetes usage

[Helm](https://helm.sh) must be installed to use the charts.

Please refer to Helm's [documentation](https://helm.sh/docs/) to get started.

Once Helm is set up properly, add the repo as follows:

```bash
helm repo add static-site https://atrakic.github.io/static-site
```

go-static-site can now be installed with the following command:

```bash
helm install static-site --namespace static-site static-site/static-site --create-namespace
```

If you have custom options or values you want to override:

```bash
helm install static-site --namespace static-site -f my-values.yaml static-site/go-static-site
```
