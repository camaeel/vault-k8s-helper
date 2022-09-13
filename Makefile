all: clean build

clean:
	rm -rf bin || true

build: build_setup_tls build_vault_autounseal

build_setup_tls: build_setup_tls_amd64 build_setup_tls_arm64

build_vault_autounseal: build_vault_autounseal_amd64 build_vault_autounseal_arm64

build_vault_autounseal_amd64:
	GOOS=linux GOARCH=amd64 go build -o bin/vault-autounseal-amd64 github.com/camaeel/vault-k8s-helper/cmd/vaultAutounseal

build_vault_autounseal_arm64:
	GOOS=linux GOARCH=arm64 go build -o bin/vault-autounseal-arm64 github.com/camaeel/vault-k8s-helper/cmd/vaultAutounseal

build_setup_tls_amd64:
	GOOS=linux GOARCH=amd64 go build -o bin/setup-tls-linux-amd64 github.com/camaeel/vault-k8s-helper/cmd/setupTls

build_setup_tls_arm64:
	GOOS=linux GOARCH=arm64 go build -o bin/setup-tls-linux-arm64 github.com/camaeel/vault-k8s-helper/cmd/setupTls/

docker:
	docker buildx build -t vault-k8s-helper:local --load .

autounseal_kind: docker docker_kind_load
	kubectl run --rm -it --image vault-k8s-helper:local test --command -- /vault-autounseal

docker_kind_load:
	kind load docker-image vault-k8s-helper:local
