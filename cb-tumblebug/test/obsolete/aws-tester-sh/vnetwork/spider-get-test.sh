#!/bin/bash
source ../setup.env

for NAME in "${CONNECT_NAMES[@]}"
do
	ID=`curl -H "${AUTH}" -sX GET http://$RESTSERVER:1024/vpc?connection_name=${NAME} |json_pp |grep "\"Id\"" |awk '{print $3}' |sed 's/"//g' |sed 's/,//g'`
	curl -H "${AUTH}" -sX GET http://$RESTSERVER:1024/vpc/${ID}?connection_name=${NAME} |json_pp &
done
