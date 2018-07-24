FROM golang:1.8 as builder
ADD . /go/src/github.com/previousnext/mysql-toolkit
WORKDIR /go/src/github.com/previousnext/mysql-toolkit
RUN go get github.com/mitchellh/gox
RUN make cli

FROM previousnext/mysql:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/previousnext/mysql-toolkit/bin/mtk_linux_amd64 /usr/local/bin/mtk

ENV AWS_CODEBUILD_DOCKERFILE=/root/Dockerfile
ENV AWS_CODEBUILD_SPEC=/root/buildspec.yml
ENV MYSQL_CONFIG=/root/config.yml

ADD example/codebuild/Dockerfile /root/Dockerfile
ADD example/codebuild/buildspec.yml /root/buildspec.yml
ADD example/config.yml /root/config.yml

ADD containers/cli/scripts/database-sanitize.sh /usr/local/bin/database-sanitize
RUN chmod +x /usr/local/bin/*

CMD ["mtk"]