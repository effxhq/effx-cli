ORG		 := effx
APP    := effx-cli
NAME   := ${ORG}/${APP}
TAG    := $$(git tag --points-at HEAD)
IMG    := ${NAME}:${TAG}
LATEST := ${NAME}:latest

docker/registry/login:
	@docker login

docker/build:
	@docker build -t ${IMG} . --build-arg GITHUB_TOKEN=${GITHUB_PERSONAL_TOKEN} --build-arg GITHUB_USER=${GITHUB_USER}
	@docker tag ${IMG} ${LATEST}

docker/image/tag:
	@docker tag ${LATEST} effx/${APP}:latest
	@docker tag ${LATEST} effx/${APP}:${TAG}

docker/registry/push:
	@docker push effx/${APP}:latest
	@docker push effx/${APP}:${TAG}

docker/build-and-push: docker/registry/login docker/build docker/image/tag docker/registry/push
