# Makefile

build-node:
	docker build --tag tbb-node --rm .

run-node:
	docker run -d --name tbb-node tbb-node
