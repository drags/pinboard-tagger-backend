.PHONY: build push release

TAG ?= "latest"

build:
	docker build -t drags/pinboard-tagger-backend:$(TAG) .

push:
	docker push drags/pinboard-tagger-backend:$(TAG)

release: build push
