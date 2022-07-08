#!/bin/bash

function CallTB() {
	echo "- Register sshKey in ${MCIRRegionName}"

	curl -H "${AUTH}" -sX POST http://$TumblebugServer/tumblebug/ns/$NSID/resources/sshKey?option=register -H 'Content-Type: application/json' -d \
		'{ 
			"connectionName": "'${CONN_CONFIG[$INDEX,$REGION]}'", 
			"name": "'${CONN_CONFIG[$INDEX,$REGION]}'-'${POSTFIX}'", 
			"cspSshKeyId": "jhseo-test",
			"fingerprint": "test-fingerprint",
			"username": "test-username",
			"publicKey": "test-public-key",
			"privateKey": "test-private-key"
		}' | jq ''
}

#function register_sshKey() {

	echo "####################################################################"
	echo "## 5. sshKey: Register"
	echo "####################################################################"

	source ../init.sh

	if [ "${INDEX}" == "0" ]; then
        echo "[Parallel execution for all CSP regions]"
        INDEXX=${NumCSP}
        for ((cspi = 1; cspi <= INDEXX; cspi++)); do
            INDEXY=${NumRegion[$cspi]}
            CSP=${CSPType[$cspi]}
            echo "[$cspi] $CSP details"
            for ((cspj = 1; cspj <= INDEXY; cspj++)); do
                echo "[$cspi,$cspj] ${RegionName[$cspi,$cspj]}"

				MCIRRegionName=${RegionName[$cspi,$cspj]}

				CallTB

			done

		done
		wait

	else
		echo ""
		
		MCIRRegionName=${CONN_CONFIG[$INDEX,$REGION]}

		CallTB

	fi
	
#}

#register_sshKey