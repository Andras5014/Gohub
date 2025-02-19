FROM ubuntu:20.04
COPY gohub /app/gohub
WORKDIR /app

ENTRYPOINT ["/app/gohub"]