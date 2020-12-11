#!/bin/bash
source ../setup.env

for NAME in "${CONNECT_NAMES[@]}"
do
	curl -H "${AUTH}" -sX POST http://$TUMBLEBUG_IP:1323/ns/$NS_ID/resources/vNet -H 'Content-Type: application/json' -d '{"connectionName":"'$NAME'", "cspNetworkName":"jhseo-test"}' | json_pp &
done
