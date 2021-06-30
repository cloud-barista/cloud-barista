#!/bin/bash
# ------------------------------------------------------------------------------
# usage
if [ "$#" -lt 1 ]; then 
	echo "./mcir-list.sh <namespace> [all/image/spec/ssh/sg/vpc]"
	echo "./mcir-list.sh cb-ladybug-ns all"
	exit 0
fi


# ------------------------------------------------------------------------------
# const
c_URL_TUMBLEBUG="http://localhost:1323/tumblebug"
c_CT="Content-Type: application/json"
c_AUTH="Authorization: Basic $(echo -n default:default | base64)"


# ------------------------------------------------------------------------------
# variables

# 1. namespace
if [ "$#" -gt 0 ]; then v_NAMESPACE="$1"; fi
if [ "${v_NAMESPACE}" == "" ]; then 
	read -e -p "Namespace ? : " v_NAMESPACE
fi
if [ "${v_NAMESPACE}" == "" ]; then echo "[ERROR] missing <namespace>"; exit -1; fi

# 2. query
if [ "$#" -gt 1 ]; then v_QUERY="$2"; fi
if [ "${v_QUERY}" == "" ]; then 
	read -e -p "Query ? [all/image/spec/ssh/sg/vpc] : "  v_QUERY
fi
if [ "${v_QUERY}" == "" ]; then echo "[ERROR] missing <query>"; exit -1; fi
if [ "${v_QUERY}" == "all" ]; then v_QUERY="image,spec,ssh,sg,vpc"; fi


# ------------------------------------------------------------------------------
# print info.
echo ""
echo "[INFO]"
echo "- Namespace                  is '${v_NAMESPACE}'"
echo "- Query                      is '${v_QUERY}'"

NM_TUMBLEBUG_NS="${c_URL_TUMBLEBUG}/ns/${v_NAMESPACE}"


# ------------------------------------------------------------------------------
# list
list() {
	if [[ "${v_QUERY}" == *"image"* ]]; then echo "IMAGE";     curl -sX GET ${NM_TUMBLEBUG_NS}/resources/image 						-H "${c_AUTH}" -H "${c_CT}" | jq; fi
	if [[ "${v_QUERY}" == *"spec"* ]]; then echo "SPEC";      curl -sX GET ${NM_TUMBLEBUG_NS}/resources/spec   -H "${c_AUTH}" -H "${c_CT}" | jq; fi
	if [[ "${v_QUERY}" == *"ssh"* ]]; then	echo "SSHKEY";     curl -sX GET ${NM_TUMBLEBUG_NS}/resources/sshKey   -H "${c_AUTH}"	-H "${c_CT}" | jq; fi
	if [[ "${v_QUERY}" == *"sg"* ]]; then	echo "SECURITYGROUP";     curl -sX GET ${NM_TUMBLEBUG_NS}/resources/securityGroup   -H "${c_AUTH}"	-H "${c_CT}" | jq; fi
	if [[ "${v_QUERY}" == *"vpc"* ]]; then	echo "VPC";     curl -sX GET ${NM_TUMBLEBUG_NS}/resources/vNet   -H "${c_AUTH}" -H "${c_AUTH}" -H "${c_CT}" | jq; fi
}

# ------------------------------------------------------------------------------
if [ "$1" != "-h" ]; then 
	echo ""
	echo "------------------------------------------------------------------------------"
	list;
fi
