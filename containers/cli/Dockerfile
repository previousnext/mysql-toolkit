FROM golang:1.10 as builder
ADD . /go/src/github.com/previousnext/mysql-toolkit
WORKDIR /go/src/github.com/previousnext/mysql-toolkit
RUN go get github.com/mitchellh/gox
RUN make cli

FROM previousnext/mysql:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/previousnext/mysql-toolkit/bin/mtk_linux_amd64 /usr/local/bin/mtk

ENV AWS_CODEBUILD_DOCKERFILE=/etc/mtk/codebuild/Dockerfile
ENV AWS_CODEBUILD_SPEC=/etc/mtk/codebuild/buildspec.yml
ENV MYSQL_CONFIG=/etc/mtk/config.yml

RUN mkdir -p /etc/mtk/codebuild
ADD example/codebuild/Dockerfile /etc/mtk/codebuild/Dockerfile
ADD example/codebuild/buildspec.yml /etc/mtk/codebuild/buildspec.yml
ADD example/mtk.yml /etc/mtk/config.yml

# Helper script for importing, sanitizing and exporting a database.
ADD containers/cli/scripts/database-sanitize.sh /usr/local/bin/database-sanitize
RUN chmod +x /usr/local/bin/*

CMD ["mtk"]