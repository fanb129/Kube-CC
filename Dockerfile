# 第一阶段：构建应用程序
FROM golang:alpine AS builder

WORKDIR /go/src/github.com/fanb129/Kube-CC

COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env \
    && go mod tidy \
    && go build -o kubecc .

# 第二阶段：运行应用程序
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /go/src/github.com/fanb129/Kube-CC/kubecc .

EXPOSE 8080

CMD ["./kubecc"]