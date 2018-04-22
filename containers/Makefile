#!/usr/bin/make -f

IMAGE=previousnext/mysql
TAG=latest

build:
	docker build -t ${IMAGE}:${TAG} mysql

push:
	docker push ${IMAGE}:${TAG}

.PHONY: build push
