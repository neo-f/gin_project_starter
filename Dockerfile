FROM golang:alpine AS build-stage
WORKDIR /go/src/gin_project_starter
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -tags=jsoniter -o server ./src

FROM alpine
ENV TZ=Asia/Shanghai
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
  && apk --no-cache --update add tzdata ca-certificates

WORKDIR /server
COPY --from=build-stage /go/src/gin_project_starter/server .
COPY --from=build-stage /go/src/gin_project_starter/configs ./configs

CMD ["./server"]
