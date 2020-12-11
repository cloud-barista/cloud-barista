#!/bin/bash

#function deploy_nginx_to_mcis() {
    FILE=../conf.env
    if [ ! -f "$FILE" ]; then
        echo "$FILE does not exist."
        exit
    fi

	source ../conf.env
	AUTH="Authorization: Basic $(echo -n $ApiUsername:$ApiPassword | base64)"

	echo "####################################################################"
	echo "## Command (SSH) to MCIS VM with given ID"
	echo "####################################################################"


	MCISID=${1:-no}
	VMID=${2:-no}

	if [ "${MCISID}" != "no" ]; then

		if [ "${VMID}" != "no" ]; then

			curl -H "${AUTH}" -sX POST http://$TumblebugServer/tumblebug/ns/$NS_ID/cmd/mcis/$MCISID/vm/$VMID -H 'Content-Type: application/json' -d \
			'{
				"command": "wget https://raw.githubusercontent.com/cloud-barista/cb-tumblebug/master/assets/scripts/setweb.sh -O ~/setweb.sh; chmod +x ~/setweb.sh; sudo ~/setweb.sh"
			}' | json_pp #|| return 1

		fi

	fi



#}

#deploy_nginx_to_mcis