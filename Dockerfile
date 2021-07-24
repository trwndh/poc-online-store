FROM golang:1.15.6-alpine AS builder
ENV GO111MODULE=on

RUN apk add --update gcc openssh git bash libc-dev ca-certificates make g++
ENV BUILDDIR /go/src/poc-online-store

ENV ENV production
COPY . /go/src/poc-online-store
WORKDIR /go/src/poc-online-store
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /go/src/poc-online-store/poc-online-store /go/src/poc-online-store/main.go

# Stage Runtime Applications
FROM alpine:latest

# Setting timezone
ENV TZ=Asia/Jakarta
RUN apk add -U tzdata
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# add ca-certificates
RUN apk add --no-cache ca-certificates

ENV BUILDDIR /go/src/poc-online-store

# Setting folder workdir
WORKDIR /opt/
# Copy Data App
COPY --from=builder $BUILDDIR/poc-online-store .
COPY --from=builder $BUILDDIR/config.env .

EXPOSE 9999

ENTRYPOINT ["./poc-online-store"]