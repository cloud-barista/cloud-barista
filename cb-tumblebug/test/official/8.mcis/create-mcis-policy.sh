#!/bin/bash

#function create_mcis_policy() {
    FILE=../conf.env
    if [ ! -f "$FILE" ]; then
        echo "$FILE does not exist."
        exit
    fi

	source ../conf.env
	AUTH="Authorization: Basic $(echo -n $ApiUsername:$ApiPassword | base64)"

	echo "####################################################################"
	echo "## 8. Create MCIS Policy"
	echo "####################################################################"

	CSP=${1}
	REGION=${2:-1}
	POSTFIX=${3:-developer}
	MCISNAME=${4:-noname}
	if [ "${CSP}" == "all" ]; then
		echo "[Test for all CSP regions (AWS, Azure, GCP, Alibaba, ...)]"
		CSP="aws"
		INDEX=0
	elif [ "${CSP}" == "aws" ]; then
		INDEX=1
	elif [ "${CSP}" == "azure" ]; then
		INDEX=2
	elif [ "${CSP}" == "gcp" ]; then
		INDEX=3
	elif [ "${CSP}" == "alibaba" ]; then
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

	if [ "${MCISNAME}" != "noname" ]; then
		echo "[MCIS name is given]"
		MCISID=${MCISNAME}
	fi

	curl -H "${AUTH}" -sX POST http://$TumblebugServer/tumblebug/ns/$NS_ID/policy/mcis/$MCISID -H 'Content-Type: application/json' -d \
		'{
			"description": "Tumblebug Auto Control Demo",
			"policy": [
				{
					"autoCondition": {
						"metric": "cpu",
						"operator": ">=",
						"operand": "80",
						"evaluationPeriod": "10"
					},
					"autoAction": {
						"actionType": "ScaleOut",
						"placement_algo": "random",
						"vm": {
							"name": "AutoGen"
						},
						"postCommand": {
							"command": "wget https://raw.githubusercontent.com/cloud-barista/cb-tumblebug/master/assets/scripts/setweb.sh -O ~/setweb.sh; chmod +x ~/setweb.sh; sudo ~/setweb.sh; wget https://raw.githubusercontent.com/cloud-barista/cb-tumblebug/master/assets/scripts/runLoadMaker.sh -O ~/runLoadMaker.sh; chmod +x ~/runLoadMaker.sh; sudo ~/runLoadMaker.sh"
						}
					}
				},				
				{
					"autoCondition": {
						"metric": "cpu",
						"operator": "<=",
						"operand": "60",
						"evaluationPeriod": "10"
					},
					"autoAction": {
						"actionType": "ScaleIn"
					}
				}
			]
		}' | json_pp || return 1
#}

#create_mcis