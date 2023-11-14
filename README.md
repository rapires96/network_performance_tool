# Network performance tool

Network performance application to measure round trip time (rtt) and throughput, using ping and iperf3 respectively. Then it sends the real-time statistics to a time series database (influxDB). The you may also set up grafana to produce better visualizations.

## To run the app

For running the project in your local machine you must first create a docker network, if this address range is already taken then change this command and all docker-compose files to a suitable network addressing.
```bash
sudo docker network create --gateway 172.100.1.1 --subnet 172.100.1.0/24 app_subnet
```

First run the iperf server and the tig stack
```bash
cd iperfServer && sudo docker compose up -d # iperf server
cd ../tig_stack && sudo docker compose up -d # influxdb and grafana
```

Then you start the client.
```bash
cd app && sudo docker compose up
```
If you wish to change the measurement parameters you may modify /app/config yaml files. For instance changing the running boolean to false makes the script to not run a measurement. 