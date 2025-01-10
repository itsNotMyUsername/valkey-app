IMAGE=emvoo/valkey-app
TAG=latest

docker-build:
	docker build -t ${IMAGE}:${TAG} .

docker-push:
	docker push ${IMAGE}:${TAG}
