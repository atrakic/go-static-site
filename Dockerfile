FROM golang:1.23 AS builder
WORKDIR /app
ADD . .

# Statically compile our app for use in a distroless container
RUN CGO_ENABLED=0 go build -ldflags="-w -s -X main.SHA=$(git rev-parse HEAD)" -v -o app .
#RUN sh -c 'CGO_ENABLED=0 go build -ldflags "-w -s -X main.SHA=$(git rev-parse HEAD)" -o app .'

# A distroless container image with some basics like SSL certificates
# https://github.com/GoogleContainerTools/distroless
FROM gcr.io/distroless/static
COPY _site /_site
COPY --from=builder /app/app ./
EXPOSE 80
ENTRYPOINT ["/app"]
