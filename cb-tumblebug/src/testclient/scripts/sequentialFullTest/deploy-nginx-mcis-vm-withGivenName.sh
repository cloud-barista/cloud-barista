#!/bin/bash

#function deploy_nginx_to_mcis() {


	TestSetFile=${4:-../testSet.env}
    if [ ! -f "$TestSetFile" ]; then
        echo "$TestSetFile does not exist."
        exit
    fi
	source $TestSetFile
    source ../conf.env
	
	echo "####################################################################"
	echo "## Command (SSH) to MCIS VM with given ID"
	echo "####################################################################"


	MCISID=${1:-no}
	VMID=${2:-no}

	if [ "${MCISID}" != "no" ]; then

		if [ "${VMID}" != "no" ]; then

			curl -H "${AUTH}" -sX POST http://$TumblebugServer/tumblebug/ns/$NSID/cmd/mcis/$MCISID/vm/$VMID -H 'Content-Type: application/json' -d \
			'{
				"command": "wget https://raw.githubusercontent.com/cloud-barista/cb-tumblebug/main/assets/scripts/setweb.sh -O ~/setweb.sh; chmod +x ~/setweb.sh; sudo ~/setweb.sh"
			}' | jq '' #|| return 1

		fi

	fi



#}

#deploy_nginx_to_mcis