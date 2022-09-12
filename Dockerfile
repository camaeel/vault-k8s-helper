FROM --platform=$BUILDPLATFORM golang:1.19-alpine as builder


# Install our build tools
RUN apk add --update ca-certificates

WORKDIR /app

ARG TARGETOS
ARG TARGETARCH
ENV LDFLAGS "-X 'main.VERSION=${RELEASE_VERSION}' "

COPY . ./

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o bin/setup-tls github.com/camaeel/vault-k8s-helper/cmd/setup-tls

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bin/setup-tls /setup-tls

ENTRYPOINT ["/setup-tls"]
