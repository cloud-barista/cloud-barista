#!/bin/bash

#function register_spec() {
    FILE=../conf.env
    if [ ! -f "$FILE" ]; then
        echo "$FILE does not exist."
        exit
    fi

	source ../conf.env
	AUTH="Authorization: Basic $(echo -n $ApiUsername:$ApiPassword | base64)"

	echo "####################################################################"
	echo "## 7. spec: Update"
	echo "####################################################################"

	curl -H "${AUTH}" -sX PUT http://$TumblebugServer/tumblebug/ns/$NS_ID/resources/spec/aws-us-east-1-m5ad.2xlarge -H 'Content-Type: application/json' -d \
		'{ 
			"id": "aws-us-east-1-m5ad.2xlarge", 
			"description": "UpdateSpec() test"
		}' | json_pp #|| return 1
#}

#register_spec
