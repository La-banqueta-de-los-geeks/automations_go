FROM golang:1.15.2-buster
WORKDIR /app
RUN go mod init go_redis_example
RUN go get github.com/go-redis/redis/v8
COPY . .
#ENTRYPOINT ["reflex”, "-c", "reflex.conf"]
CMD bash