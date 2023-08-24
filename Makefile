check:
	golangci-lint run -v

check-fix:
	golangci-lint run --fix

docker-build: build-url-shortener build-postgres build-backup

build-url-shortener:
	DOCKER_BUILDKIT=1 docker --log-level=debug build --pull --build-arg BUILDKIT_INLINE_CACHE=1 \
		--target builder \
		--cache-from ${REGISTRY}/url-shortener:cache-builder \
		--tag ${REGISTRY}/url-shortener:cache-builder \
		--file ./Dockerfile .

	DOCKER_BUILDKIT=1 docker --log-level=debug build --pull --build-arg BUILDKIT_INLINE_CACHE=1 \
    	--cache-from ${REGISTRY}/url-shortener:cache-builder \
    	--cache-from ${REGISTRY}/url-shortener:cache \
    	--tag ${REGISTRY}/url-shortener:cache \
    	--tag ${REGISTRY}/url-shortener:${IMAGE_TAG} \
    	--file ./Dockerfile .

build-postgres:
	DOCKER_BUILDKIT=1 docker --log-level=debug build --pull --build-arg BUILDKIT_INLINE_CACHE=1 \
        --cache-from ${REGISTRY}/url-shortener-postgres:cache \
        --tag ${REGISTRY}/url-shortener-postgres:cache \
        --tag ${REGISTRY}/url-shortener-postgres:${IMAGE_TAG} \
        --file docker/postgres/Dockerfile docker/postgres

build-backup:
	docker --log-level=debug build --pull --file=docker/common/postgres-backup/Dockerfile --tag=${REGISTRY}/url-shortener-postgres-backup:${IMAGE_TAG} docker/common

push-build-cache:
	docker push ${REGISTRY}/url-shortener:cache-builder
	docker push ${REGISTRY}/url-shortener:cache
	docker push ${REGISTRY}/url-shortener-postgres:cache

push:
	docker push ${REGISTRY}/url-shortener:${IMAGE_TAG}
	docker push ${REGISTRY}/url-shortener-postgres:${IMAGE_TAG}
	docker push ${REGISTRY}/url-shortener-postgres-backup:${IMAGE_TAG}

deploy:
	ssh -o StrictHostKeyChecking=no deploy@${HOST} -p ${PORT} 'docker network create --driver=overlay traefik-public || true'
	ssh -o StrictHostKeyChecking=no deploy@${HOST} -p ${PORT} 'rm -rf url_shortener_${BUILD_NUMBER} && mkdir url_shortener_${BUILD_NUMBER}'

	envsubst < docker-compose-production.yml > docker-compose-production-env.yml
	scp -o StrictHostKeyChecking=no -P ${PORT} docker-compose-production-env.yml deploy@${HOST}:url_shortener_${BUILD_NUMBER}/docker-compose.yml
	rm -f docker-compose-production-env.yml

	ssh -o StrictHostKeyChecking=no deploy@${HOST} -p ${PORT} 'mkdir url_shortener_${BUILD_NUMBER}/secrets'
	ssh -o StrictHostKeyChecking=no deploy@${HOST} -p ${PORT} 'cp .secrets_url_shortener/* url_shortener_${BUILD_NUMBER}/secrets'
	ssh -o StrictHostKeyChecking=no deploy@${HOST} -p ${PORT} 'cd url_shortener_${BUILD_NUMBER} && docker stack deploy --compose-file docker-compose.yml url-shortener --with-registry-auth --prune'
