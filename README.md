# Network performance tool

## To run the app 

### 1. docker compose

Simply run the command below ans see the logs displayed in the terminal

```bash
sudo docker compose up
```

### 2. docker build and docker run 
Here is how we build the docker image. 

* 1. build the image
```bash
sudo docker build . -t vcu_app
```
* 2. make sure you have started the iperf server, for exaple in your local machine

```bash
iperf3 -s -i 0.2 
```
In this particular case, you make sure you now what is the ip address of the local machine with respect to the docker container. You may do so by running the following command and using the gateway address.

```bash
sudo docker network inspect bridge
```

* 3. Finally run the docker image and the statinstics should be shown in real time.

```bash
sudo docker run --rm -v ./src/config:/app/config -e APP_CONFIG="config/default.yaml" vcu_app
```

You may edit the default.yaml file in the config folder or create a new one within the folder and change the APP_CONFIG environment variable.


