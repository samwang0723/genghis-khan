# Stage 1 (to create a "build" image, 866MB)
FROM golang:1.11 AS Builder

ARG VERSION
ENV VERSION=${VERSION}
ARG GIT_COMMIT_HASH
ENV GIT_COMMIT_HASH=${GIT_COMMIT_HASH}

COPY . /go/src/github.com/samwang0723/genghis-khan/
WORKDIR /go/src/github.com/samwang0723/genghis-khan/
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -tags=jsoniter -a -o bin/genghis-khan

# Stage 2 (to create a downsized "container executable", 11.6MB)

# If you need SSL certificates for HTTPS, replace `FROM SCRATCH` with:
#
#   FROM alpine:3.7
#   RUN apk --no-cache add ca-certificates
#
FROM alpine
WORKDIR /go/bin
ENV GOPATH=/go
RUN apk add --update ca-certificates

COPY --from=builder /go/src/github.com/samwang0723/genghis-khan/bin/genghis-khan /go/bin/genghis-khan

EXPOSE 8080
CMD ["/go/bin/genghis-khan"]