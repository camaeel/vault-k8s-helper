FROM --platform=$BUILDPLATFORM golang:1.19-alpine as builder


# Install our build tools
RUN apk add --update ca-certificates

WORKDIR /app

ARG TARGETOS
ARG TARGETARCH

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN \
  CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o bin/setup-tls github.com/camaeel/vault-k8s-helper/cmd/setupTls && \
  CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o bin/vault-autounseal github.com/camaeel/vault-k8s-helper/cmd/vaultAutounseal

#####################################################
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bin/setup-tls /setup-tls
COPY --from=builder /app/bin/vault-autounseal /vault-autounseal

ENTRYPOINT ["/setup-tls"]
