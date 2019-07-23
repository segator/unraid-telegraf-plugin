FROM golang:latest as builder
COPY . /build/
WORKDIR /build
RUN go get -d ./ && go build -o unraid-influxdb-line . && chmod +x unraid-influxdb-line

FROM telegraf:alpine

RUN apk --update add --no-cache --virtual smartmontools && \
    apk --update add --no-cache --virtual ipmitool && \
    apk --update add --no-cache --virtual apcupsd

COPY --from=builder /build/unraid-influxdb-line /app/
