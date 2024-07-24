.PHONY: docker
docker:
	-del /f webook
	@go build -o webook .
	@docker rmi -f andras/webook:v0.0.1
	@docker build -t andras/webook:v0.0.1 .
