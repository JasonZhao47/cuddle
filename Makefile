.PHONY: docker
docker:
	@rm cuddle || true
	@go mod tidy
	@GOOS=linux GOARCH=arm64 go build -tags=k8s -o cuddle .
	@docker rmi -f jasonzhao47/cuddle:v0.2
	@docker build -t jasonzhao47/cuddle:v0.2 .