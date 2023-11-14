FROM golang:1.20-buster

#Install iperf3
RUN apt-get update && apt-get install -y iperf3

EXPOSE 5201
EXPOSE 5001

#Copy scripts Into Container
COPY . /app/src
#COPY ../config /app/config

WORKDIR /app/src
RUN go build -o run

#ENV APP_CONFIG "./config/measurement.yaml"
#ENV DB_CONFIG "config/database.yaml"
ENTRYPOINT ["./run"]