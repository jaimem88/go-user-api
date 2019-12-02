FROM alpine:3.10

EXPOSE 8899
COPY version.json .
COPY bin/go-user-api /go/bin/

ENTRYPOINT ["/go/bin/go-user-api"]
