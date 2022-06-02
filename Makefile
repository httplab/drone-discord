build:
	CGO_ENABLED=0 go build

push:
	docker build -t httplab/drone-discord .
	docker push httplab/drone-discord

release:
	$(MAKE) build
	$(MAKE) push
