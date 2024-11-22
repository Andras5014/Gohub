.PHONY: docker



docker:
	-del /f webook
	@go build -o webook .
	@docker rmi -f andras/webook:v0.0.1
	@docker build -t andras/webook:v0.0.1 .

.PHONY: mock
mock:

	@mockgen -source=./internal/service/user.go -destination=./internal/service/mocks/user.go -package=svcmocks
	@mockgen -source=./internal/repository/user.go -destination=./internal/repository/mocks/user.go -package=repomocks
	@mockgen -source=./internal/repository/dao/user.go -destination=./internal/repository/dao/mocks/user.go -package=daomocks
	@mockgen -source=./internal/repository/cache/user.go -destination=./internal/repository/cache/mocks/user.go -package=cachemocks


	@mockgen -source=./internal/service/article.go -destination=./internal/service/mocks/article.go -package=svcmocks
	@mockgen -source=./internal/repository/article/article.go -destination=./internal/repository/article/mocks/article.go -package=artrepomocks
	@mockgen -source=./internal/repository/article/article_author.go -destination=./internal/repository/article/mocks/article_author.go -package=artrepomocks
	@mockgen -source=./internal/repository/article/article_reader.go -destination=./internal/repository/article/mocks/article_reader.go -package=artrepomocks


	@mockgen -source=./internal/service/code.go -destination=./internal/service/mocks/code.go -package=svcmocks
	@mockgen -source=./internal/repository/code.go -destination=./internal/repository/mocks/code.go -package=repomocks

	@mockgen  -package=redismocks -destination=./internal/repository/cache/redismocks/cmdable.go github.com/redis/go-redis/v9 Cmdable

	@go mod tidy
