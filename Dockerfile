FROM golang:alpine as build

WORKDIR /root
COPY . /root
RUN go build .

FROM alpine:latest

ENV SWITCH_IP_ADDRESS="192.168.10.10" LISTEN_ADDRESS="0.0.0.0" LISTEN_PORT="9452"

RUN addgroup -S mystrom \
    && adduser -S mystrom -G mystrom \
    && mkdir /app \
    && chown -R mystrom:mystrom /app

WORKDIR /app
USER mystrom
COPY --from=build /root/mystrom-exporter /app

CMD /app/mystrom-exporter -web.listen-address $LISTEN_ADDRESS:$LISTEN_PORT -switch.ip-address $SWITCH_IP_ADDRESS