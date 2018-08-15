FROM golang:1.8 as builder
ADD . /go/src/github.com/previousnext/mysql-toolkit
WORKDIR /go/src/github.com/previousnext/mysql-toolkit
RUN go get github.com/mitchellh/gox
RUN make operator

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/previousnext/mysql-toolkit/bin/mtk-operator_linux_amd64 /usr/local/bin/mtk-operator

CMD ["mtk-operator"]