all: clean build

clean:
	rm -rf bin || true

build: build_setup_tls

build_setup_tls: build_setup_tls_amd64 build_setup_tls_arm64

build_setup_tls_amd64:
	GOOS=linux GOARCH=amd64 go build -o bin/setup-tls-linux-amd64 github.com/camaeel/vault-k8s-helper/cmd/setupTls

build_setup_tls_arm64:
	GOOS=linux GOARCH=arm64 go build -o bin/setup-tls-linux-arm64 github.com/camaeel/vault-k8s-helper/cmd/setupTls/

docker:
	docker buildx build -t vault-k8s-helper:local .
