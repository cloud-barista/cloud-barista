#!/bin/bash
source ../setup.env

#for NAME in "${CONNECT_NAMES[@]}"
#do
#        echo ========================== $NAME
#
#	VM_ID=vm-powerkim01
#	echo ....terminate ${VM_ID} ...
#	curl -H "${AUTH}" -sX DELETE http://$RESTSERVER:1024/vm/${VM_ID}?connection_name=${NAME} &
#done

MCIS_ID=`curl -H "${AUTH}" -sX GET http://$TUMBLEBUG_IP:1323/ns/$NS_ID/mcis | jq -r '.mcis[].id'`
curl -H "${AUTH}" -sX DELETE http://$TUMBLEBUG_IP:1323/ns/$NS_ID/mcis/$MCIS_ID
