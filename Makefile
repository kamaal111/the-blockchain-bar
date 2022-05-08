# Makefile

build-node:
	docker build --file docker/Dockerfile.node --tag tbb-node --force-rm .

run-node:
	docker run -d --name tbb-node tbb-node

log-node:
	docker logs tbb-node -n 10 --timestamps