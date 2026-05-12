# renovate: datasource=docker depName=golang
ARG GO_VERSION=1.26
FROM golang:${GO_VERSION} AS builder
WORKDIR /build
ADD . .
RUN CGO_ENABLED=0 go build -ldflags="-w -s -X main.SHA=$(git rev-parse HEAD)" -v -o app .

FROM scratch as final
COPY _site /_site
COPY --from=builder /build/app /static-site
EXPOSE 8080
ENTRYPOINT ["/static-site"]
