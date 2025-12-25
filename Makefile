IMAGE_NAME=eugeniypetrov/browscap-go
VERSION?=$(shell git describe --tags --always --dirty)

build:
	DOCKER_BUILDKIT=1 docker buildx build \
		--ssh default \
		--platform linux/amd64 \
		-t $(IMAGE_NAME):$(VERSION) \
		.

	docker tag $(IMAGE_NAME):$(VERSION) $(IMAGE_NAME):latest

push:
	docker push $(IMAGE_NAME):$(VERSION)
	docker push $(IMAGE_NAME):latest
