.PHONY: docker
docker:
	@rm cuddle || true
	@go mod tidy
	@GOOS=linux GOARCH=arm64 go build -tags=k8s -o cuddle .
	@docker rmi -f jasonzhao47/cuddle:v0.2
	@docker build -t jasonzhao47/cuddle:v0.2 .

.PHONY: mock
mock:
	@mockgen -source=internal/service/article.go -destination=internal/service/mocks/article.mock.go -package=svcmock
	@mockgen -source=internal/repository/dao/article.go -destination=internal/dao/mocks/article.mock.go -package=daomock
	@mockgen -source=internal/service/user_activity.go -destination=internal/service/mocks/user_activity.mock.go -package=svcmock
	@mockgen -source=pkg/ginx/middleware/ratelimit/limiter.go -destination=internal/pkg/mocks/limiter.mock.go -package=pkgmock
	@mockgen -source=internal/service/sms/sms_service.go -destination=internal/service/mocks/sms_service.mock.go -package=svcmock
