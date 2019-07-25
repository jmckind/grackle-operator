.PHONY: build

IMAGE_REPO ?= quay.io/jmckind
IMAGE_NAME ?= grackle-operator
IMAGE_TAG  ?= latest
IMAGE_URL  := $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)
NAMESPACE  ?= grackle

build:
	operator-sdk build $(IMAGE_URL)
	docker tag $(IMAGE_URL) $(IMAGE_REPO)/$(IMAGE_NAME):latest

docker-push:
	docker push $(IMAGE_URL)
	docker push $(IMAGE_REPO)/$(IMAGE_NAME):latest

run-local:
	operator-sdk up local --namespace $(NAMESPACE)
