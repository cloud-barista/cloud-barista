#!/bin/bash
source ../setup.env


num=0
for NAME in "${CONNECT_NAMES[@]}"
do
        #ID=`curl -H "${AUTH}" -sX GET http://$RESTSERVER:1024/publicip?connection_name=${NAME} |json_pp |grep "\"Name\" :" |awk '{print $3}' | head -n 1 |sed 's/"//g' |sed 's/,//g'`
	ID=publicipt${num}-powerkim
        curl -H "${AUTH}" -sX GET http://$RESTSERVER:1024/publicip/${ID}?connection_name=${NAME} |json_pp &

	num=`expr $num + 1`
done

