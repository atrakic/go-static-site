# go-static-site

[![golangci](https://github.com/atrakic/go-static-site/actions/workflows/ci.yml/badge.svg)](https://github.com/atrakic/go-static-site/actions/workflows/ci.yml)
[![license](https://img.shields.io/github/license/atrakic/go-static-site.svg)](https://github.com/atrakic/go-static-site/blob/main/LICENSE)
[![release](https://img.shields.io/github/release/atrakic/go-static-site/all.svg)](https://github.com/atrakic/go-static-site/releases)
[![Release Charts](https://github.com/atrakic/go-static-site/actions/workflows/chart-release.yml/badge.svg)](https://github.com/atrakic/go-static-site/actions/workflows/chart-release.yml)

> A minimal static file server with go.
> Site content is autogenerated from https://github.com/atrakic/static-generator.

## Kubernetes usage

[Helm](https://helm.sh) must be installed to use the charts.

Please refer to Helm's [documentation](https://helm.sh/docs/) to get started.

Once Helm is set up properly, add the repo as follows:

```bash
helm repo add go-static-site https://atrakic.github.io/go-static-site
```

go-static-site can now be installed with the following command:

```bash
helm install go-static-site --namespace go-static-site go-static-site/go-static-site --create-namespace
```

If you have custom options or values you want to override:

```bash
helm install go-static-site --namespace go-static-site -f my-values.yaml go-static-site/go-static-site
```
