#!/bin/bash

#function delete_ns() {


    TestSetFile=${4:-../testSet.env}
    if [ ! -f "$TestSetFile" ]; then
        echo "$TestSetFile does not exist."
        exit
    fi
	source $TestSetFile
    source ../conf.env
    
    echo "####################################################################"
    echo "## 0. Namespace: Delete"
    echo "####################################################################"

    INDEX=${1}

    $CBTUMBLEBUG_ROOT/src/api/grpc/cbadm/cbadm namespace delete --config $CBTUMBLEBUG_ROOT/src/api/grpc/cbadm/grpc_conf.yaml -o json --ns $NSID | jq ''
    echo ""
#}

#delete_ns