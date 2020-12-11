#!/bin/bash
source ../setup.env

#for NAME in "${CONNECT_NAMES[@]}"
#do 
#       ID=`curl -H "${AUTH}" -sX GET http://$RESTSERVER:1024/publicip?connection_name=${NAME} |json_pp |grep "\"Name\" :" |awk '{print $3}' | head -n 1 |sed 's/"//g' |sed 's/,//g'`
#       curl -H "${AUTH}" -sX DELETE http://$RESTSERVER:1024/publicip/${ID}?connection_name=${NAME} 
#done

TB_IMAGE_IDS=`curl -H "${AUTH}" -sX GET http://$TUMBLEBUG_IP:1323/ns/$NS_ID/resources/image | jq -r '.image[].id'`
#echo $TB_IMAGE_IDS | json_pp

if [ -n "$TB_IMAGE_IDS" ]
then
        #TB_IMAGE_IDS=`curl -H "${AUTH}" -sX GET http://$TUMBLEBUG_IP:1323/ns/$NS_ID/resources/image | jq -r '.image[].id'`
        for TB_IMAGE_ID in ${TB_IMAGE_IDS}
        do
                echo ....Delete ${TB_IMAGE_ID} ...
                curl -H "${AUTH}" -sX DELETE http://$TUMBLEBUG_IP:1323/ns/$NS_ID/resources/image/${TB_IMAGE_ID} | json_pp
        done
else
        echo ....no images found
fi
