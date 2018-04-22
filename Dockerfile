FROM golang:1.8 as builder
ADD . /go/src/github.com/previousnext/mysql-toolkit
WORKDIR /go/src/github.com/previousnext/mysql-toolkit
RUN go get github.com/mitchellh/gox
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/previousnext/mysql-toolkit/bin/mysql-toolkit_linux_amd64 /usr/local/bin/mysql-toolkit

ENV SOURCE_DOCKERFILE=/root/Dockerfile
ENV SOURCE_BUILDSPEC=/root/buildspec.yml
ENV MYSQL_CONFIG=/root/config.yml

ADD example/Dockerfile /root/Dockerfile
ADD example/buildspec.yml /root/buildspec.yml
ADD example/config.yml /root/config.yml

CMD ["mysql-toolkit"]
