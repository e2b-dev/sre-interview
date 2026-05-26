.PHONY: build

build:
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		-t ghcr.io/e2b-dev/sre-interview:latest \
		--push \
		.
