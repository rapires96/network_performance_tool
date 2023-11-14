export DOCKER_GRAFANA_INIT_PASSWORD=$DOCKER_GRAFANA_INIT_PASSWORD
grafana-cli --config "/etc/configuration/" admin reset-admin-password ${DOCKER_GRAFANA_INIT_PASSWORD}