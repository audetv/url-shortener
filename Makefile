check:
	golangci-lint run -v

check-fix:
	golangci-lint run --fix

docker-build:
	docker --log-level=debug build --pull --file=./Dockerfile --tag=${REGISTRY}/url-shortener:${IMAGE_TAG} .

push:
	docker push ${REGISTRY}/url-shortener:${IMAGE_TAG}
