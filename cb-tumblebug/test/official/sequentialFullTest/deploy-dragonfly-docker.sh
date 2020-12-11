#!/bin/bash

#function deploy_nginx_to_mcis() {
    FILE=../conf.env
    if [ ! -f "$FILE" ]; then
        echo "$FILE does not exist."
        exit
    fi

	echo "[Check jq package (if not, install)]"
	if ! dpkg-query -W -f='${Status}' jq  | grep "ok installed"; then sudo apt install -y jq; fi
	

	source ../conf.env
	AUTH="Authorization: Basic $(echo -n $ApiUsername:$ApiPassword | base64)"

	echo "####################################################################"
	echo "## Command (SSH) to MCIS "
	echo "####################################################################"

	CSP=${1}
	REGION=${2:-1}
	POSTFIX=${3:-developer}
	if [ "${CSP}" == "all" ]; then
		echo "[Test for all CSP regions (AWS, Azure, GCP, Alibaba, ...)]"
		CSP="aws"
		INDEX=0
	elif [ "${CSP}" == "aws" ]; then
		echo "[Test for AWS]"
		INDEX=1
	elif [ "${CSP}" == "azure" ]; then
		echo "[Test for Azure]"
		INDEX=2
	elif [ "${CSP}" == "gcp" ]; then
		echo "[Test for GCP]"
		INDEX=3
	elif [ "${CSP}" == "alibaba" ]; then
		echo "[Test for Alibaba]"
		INDEX=4
	else
		echo "[No acceptable argument was provided (all, aws, azure, gcp, alibaba, ...). Default: Test for AWS]"
		CSP="aws"
		INDEX=1
	fi


	MCISID=${CONN_CONFIG[$INDEX,$REGION]}-${POSTFIX}

	if [ "${INDEX}" == "0" ]; then
		MCISPREFIX=avengers
		MCISID=${MCISPREFIX}-${POSTFIX}
	fi

	curl -H "${AUTH}" -sX POST http://$TumblebugServer/tumblebug/ns/$NS_ID/cmd/mcis/$MCISID -H 'Content-Type: application/json' -d \
		'{
			"command": "wget https://raw.githubusercontent.com/cloud-barista/cb-tumblebug/master/assets/scripts/setcbdf.sh -O ~/setcbdf.sh; chmod +x ~/setcbdf.sh; ~/setcbdf.sh"
		}' | json_pp #|| return 1

	MCISINFO=`curl -H "${AUTH}" -sX GET http://$TumblebugServer/tumblebug/ns/$NS_ID/mcis/${MCISID}?action=status`
	MASTERIP=$(jq -r '.status.masterIp' <<< "$MCISINFO")
	MASTERVM=$(jq -r '.status.masterVmId' <<< "$MCISINFO")
	
	echo "MASTERIP: $MASTERIP"
	echo "MASTERVM: $MASTERVM"

	echo "[Update Tumblebug Environment for Dragonfly with following command]"
	PARAM="DRAGONFLY_REST_URL http://${MASTERIP}:9090/dragonfly"
	echo $PARAM
	../2.configureTumblebug/update-config.sh $PARAM
	echo ""
	echo "[You can test Dragonfly with following command]"
	echo " ../9.monitoring/install-agent.sh ${CSP} ${REGION} ${POSTFIX}"
	echo " ../9.monitoring/get-monitoring-data.sh ${CSP} ${REGION} ${POSTFIX} cpu"
#}

#deploy_cb-df_to_mcis