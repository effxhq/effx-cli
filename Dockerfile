FROM golang:1.13.5-alpine as builder
MAINTAINER Effx Engineering

RUN apk add --no-cache ca-certificates git

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GO111MODULE=on
ENV GOPRIVATE="github.com/effxhq"

ARG GITHUB_TOKEN
ARG GITHUB_USER

RUN git config --global url."https://$GITHUB_USER:$GITHUB_TOKEN@github.com/".insteadOf "https://github.com/"

# set up working directory
RUN mkdir /build
ADD . /build/
WORKDIR /build

RUN go mod download
RUN go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o effx effx.go

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/effx /
WORKDIR /

ENTRYPOINT ["/effx"]
