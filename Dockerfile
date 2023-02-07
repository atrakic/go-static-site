FROM golang:1.19 as builder
WORKDIR /app
COPY . /app

RUN go get -d -v

# Statically compile our app for use in a distroless container
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -v -o app .

# A distroless container image with some basics like SSL certificates
# https://github.com/GoogleContainerTools/distroless
# https://iximiuz.com/en/posts/containers-distroless-images/
FROM gcr.io/distroless/static
LABEL org.opencontainers.image.source=https://github.com/atrakic/go-static-site
LABEL org.opencontainers.image.description="Static site with golang"
LABEL org.opencontainers.image.licenses=MIT

COPY _site /_site
COPY --from=builder /app/app ./
EXPOSE 8080
ENTRYPOINT ["/app"]
