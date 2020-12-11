#!/bin/bash
source ../setup.env

for NAME in "${CONNECT_NAMES[@]}"
do
	curl -H "${AUTH}" -sX POST http://$RESTSERVER:1024/vpc?connection_name=${NAME} -H 'Content-Type: application/json' -d '{"Name":"CB-VNet-powerkim"}' |json_pp &
done
