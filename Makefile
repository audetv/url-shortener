check:
	golangci-lint run -v

check-fix:
	golangci-lint run --fix

docker-build:
	DOCKER_BUILDKIT=1 docker --log-level=debug build --pull --build-arg BUILDKIT_INLINE_CACHE=1 \
		--target builder \
		--cache-from ${REGISTRY}/url-shortener:cache-builder \
		--tag ${REGISTRY}/url-shortener:cache-builder \
		--file ./Dockerfile app

	DOCKER_BUILDKIT=1 docker --log-level=debug build --pull --build-arg BUILDKIT_INLINE_CACHE=1 \
    	--cache-from ${REGISTRY}/url-shortener:cache-builder \
    	--cache-from ${REGISTRY}/url-shortener:cache \
    	--tag ${REGISTRY}/url-shortener:cache \
    	--tag ${REGISTRY}/url-shortener:${IMAGE_TAG} \
    	--file ./Dockerfile app

	docker --log-level=debug build --pull --file=./Dockerfile --tag=${REGISTRY}/url-shortener:${IMAGE_TAG} .

push-build-cache:
	docker push ${REGISTRY}/url-shortener:cache-builder
	docker push ${REGISTRY}/url-shortener:cache

push:
	docker push ${REGISTRY}/url-shortener:${IMAGE_TAG}
