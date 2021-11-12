.PHONY: build
build:
	docker build -t gamboa/biopipe-cli:latest .

.PHONY: push
push:
	docker push gamboa/biopipe-cli:latest