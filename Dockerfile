FROM alpine:latest

RUN apk add go musl-dev imagemagick ffmpeg x264

RUN mkdir /app
WORKDIR /app

ADD . /app
RUN go build -o gopic cmd/gopic/main.go

ENTRYPOINT /app/gopic
