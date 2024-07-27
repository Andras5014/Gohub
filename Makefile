.PHONY: docker
docker:
	-del /f webook
	@go build -o webook .
	@docker rmi -f andras/webook:v0.0.1
	@docker build -t andras/webook:v0.0.1 .

.PHONY: mock
mock:
	@mockgen -source=./internal/service/user.go -destination=./internal/service/mocks/user.go -package=svcmocks
	@mockgen -source=./internal/service/code.go -destination=./internal/service/mocks/code.go -package=svcmocks
	@go mod tidy
