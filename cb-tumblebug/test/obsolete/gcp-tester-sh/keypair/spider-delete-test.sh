#!/bin/bash
source ../setup.env

#for NAME in "${CONNECT_NAMES[@]}"
#do
        NAME=${CONNECT_NAMES[0]}
        curl -H "${AUTH}" -sX DELETE http://$RESTSERVER:1024/keypair/mcb-keypair-powerkim?connection_name=${NAME} |json_pp &
#done
