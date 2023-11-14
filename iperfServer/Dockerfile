FROM alpine:3.10

RUN apk update && apk add wget
RUN apk add iperf3

EXPOSE 5201
EXPOSE 5001

CMD ["iperf3", "-s", "-i", "0.2", "-p", "5201"]