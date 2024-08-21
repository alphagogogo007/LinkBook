.PHONY: docker
docker:
	@rm webook || true
	@docker rmi -f hx13/webook:v0.0.1
	@go mod tidy
	@GOOS=linux GOARCH=amd64 go build -tags=k8s -o webook .
	@docker build -t hx13/webook:v0.0.1 .