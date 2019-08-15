#第一步构建
FROM golang:1.12 as builder
ENV GOPROXY https://goproxy.io
ENV GO111MODULE on

WORKDIR /buider
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo ./cmd/conv-server


#第二步发布
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /go/
COPY --from=builder /buider/conv-server .
COPY --from=builder /buider/dev.yaml .
CMD ["./conv-server","--conf","dev.yaml"]
