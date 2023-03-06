FROM golang:1.20-buster

#Install iperf3
RUN apt-get update && apt-get install -y iperf3

WORKDIR /app
#Copy scripts Into Container
COPY ./src .

RUN go build -o run

ENV APP_CONFIG "config/default.yaml"
ENTRYPOINT ["./run"]