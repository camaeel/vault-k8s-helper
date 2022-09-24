FROM --platform=$BUILDPLATFORM golang:1.19-alpine as builder


# Install our build tools
RUN apk add --update ca-certificates

WORKDIR /app

ARG DEBUG
ARG TARGETOS
ARG TARGETARCH
ENV LDFLAGS "-X 'main.VERSION=${RELEASE_VERSION}' "

COPY . ./



RUN if [ DEBUG -eq 1 ]; then export DEBUGFLAGS='-gcflags=all="-N -l"'; fi && \
  CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build $DEBUGFLAGS -o bin/setup-tls github.com/camaeel/vault-k8s-helper/cmd/setupTls && \
  CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build $DEBUGFLAGS -o bin/vault-autounseal github.com/camaeel/vault-k8s-helper/cmd/vaultAutounseal

FROM golang:1.19-alpine as debug

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bin/setup-tls /setup-tls
COPY --from=builder /app/bin/vault-autounseal /vault-autounseal
RUN go install github.com/go-delve/delve/cmd/dlv@latest

ENTRYPOINT ["dlv"]

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bin/setup-tls /setup-tls
COPY --from=builder /app/bin/vault-autounseal /vault-autounseal

ENTRYPOINT ["/setup-tls"]
