#!/bin/bash

if [ "$1" = "" ]; then
	echo
	echo -e 'usage: '$0' mock|aws|azure|gcp|alibaba|tencent|ibm|openstack|cloudit|ncp|nhncloud'
	echo -e '\n\tex) '$0' aws'
	echo
	exit 0;
fi

# common setup.env path
SETUP_PATH=$CBSPIDER_ROOT/test/functional-test-cli/common
source $SETUP_PATH/setup.env $1

SUBNET_NAME=${SUBNET_NAME}-Added
# SUBNET_CIDR=192.168.0.0/24
SUBNET_CIDR=`echo $SUBNET_CIDR | sed 's/0\./1\./g'`

# for Cloudit
if [ $SUBNET_CIDR_ADD ]; then
	SUBNET_CIDR=$SUBNET_CIDR_ADD
fi

echo "============== before add Subnet: '${VPC_NAME}' : '${SUBNET_NAME}'"

$CLIPATH/spctl --config $CLIPATH/spctl.conf vpc add-subnet -i json -d \
    '{
      "ConnectionName":"'${CONN_CONFIG}'",
      "VPCName": "'${VPC_NAME}'",
      "ReqInfo": {
        "Name": "'${SUBNET_NAME}'",
        "IPv4_CIDR": "'${SUBNET_CIDR}'"
      }
    }'

echo "============== after add Subnet: '${VPC_NAME}' : '${SUBNET_NAME}'"

echo -e "\n\n"

