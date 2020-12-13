#!/bin/bash
# ------------------------------------------------------------------------------
# usage
if [ "$1" == "-h" ]; then 
	echo "./list.sh [GCP/AWS] [all/config/ns/vpc/fm/mcis]"
	echo "./list.sh GCP ns"
	echo "./list.sh GCP ns,spec"
	exit 0
fi


# ------------------------------------------------------------------------------
# const

c_URL_SPIDER="http://localhost:1024/spider"
c_URL_TUMBLEBUG="http://localhost:1323/tumblebug"
c_CT="Content-Type: application/json"
c_AUTH="Authorization: Basic $(echo -n default:default | base64)"


# ------------------------------------------------------------------------------
# argument

# 1. CSP
if [ "$#" -gt 0 ]; then v_CSP="$1"; else	v_CSP="${CSP}"; fi
if [ "${v_CSP}" == "" ]; then 
	read -e -p "Cloud ? [AWS(default) or GCP] : "  v_CSP
fi

if [ "${v_CSP}" == "" ]; then v_CSP="AWS"; fi
if [ "${v_CSP}" != "GCP" ] && [ "${v_CSP}" != "AWS" ]; then echo "[ERROR] missing <cloud>"; exit -1;fi

# PREFIX
if [ "${v_CSP}" == "GCP" ]; then 
	v_PREFIX="cb-gcp"
else
	v_PREFIX="cb-aws"
fi

# # PREFIX
# if [ "$#" -gt 0 ]; then v_PREFIX="$1"; else	v_PREFIX="${PREFIX}"; fi

# if [ "${v_PREFIX}" == "" ]; then 
# 	read -e -p "Name prefix ? : "  v_PREFIX
# fi
# if [ "${v_PREFIX}" == "" ]; then echo "[ERROR] missing <prefix>"; exit -1; fi


# query
if [ "$#" -gt 1 ]; then v_QUERY="$2"; fi

if [ "${v_QUERY}" == "" ]; then 
	read -e -p "Query ? [all/config/ns/vpc/mcis/ssh] : "  v_QUERY
fi
if [ "${v_QUERY}" == "" ]; then echo "[ERROR] missing <query>"; exit -1; fi
if [ "${v_QUERY}" == "all" ]; then v_QUERY="config,ns,vpc,fm,mcis,ssh"; fi


# variable - name
NM_NAMESPACE="${v_PREFIX}-namespace"
NM_CONFIG="${v_PREFIX}-config"
NM_VPC="${v_PREFIX}-vpc"
NM_FW="${v_PREFIX}-allow-external"
NM_SSH_KEY="${v_PREFIX}-ssh-key"
NM_IMAGE="${v_PREFIX}-image"
NM_MACHINE="${v_PREFIX}-spec"
NM_MCIS="${v_PREFIX}" 

c_URL_TUMBLEBUG_NS="${c_URL_TUMBLEBUG}/ns/${NM_NAMESPACE}"



# ------------------------------------------------------------------------------
# list
list() {
	if [[ "${v_QUERY}" == *"config"* ]]; then	echo "@_CONFIG_@";		curl -sX GET ${c_URL_SPIDER}/connectionconfig              -H "${c_AUTH}" | jq; fi
	if [[ "${v_QUERY}" == *"ns"* ]]; then		echo "@_NAMESPACE_@";	curl -sX GET ${c_URL_TUMBLEBUG}/ns                         -H "${c_AUTH}" | jq; fi
	if [[ "${v_QUERY}" == *"vpc"* ]]; then		echo "@_VPC_@";			curl -sX GET ${c_URL_TUMBLEBUG_NS}/resources/vNet          -H "${c_AUTH}" | jq; fi
	if [[ "${v_QUERY}" == *"fw"* ]]; then		echo "@_FW_@";			curl -sX GET ${c_URL_TUMBLEBUG_NS}/resources/securityGroup -H "${c_AUTH}" | jq; fi
	if [[ "${v_QUERY}" == *"ssh"* ]]; then		echo "@_FW_@";			curl -sX GET ${c_URL_TUMBLEBUG_NS}/resources/securityGroup -H "${c_AUTH}" | jq; fi
	if [[ "${v_QUERY}" == *"ssh"* ]]; then		echo "@_SSH_@";			curl -sX GET ${c_URL_TUMBLEBUG_NS}/resources/sshKey        -H "${c_AUTH}" -H "${c_CT}" -d '{"connectionName" : "'${NM_CONFIG}'"}' | jq; fi
	if [[ "${v_QUERY}" == *"mcis"* ]]; then		echo "@_MCIS_@";		curl -sX GET ${c_URL_TUMBLEBUG_NS}/mcis                    -H "${c_AUTH}" | jq; fi
}


if [ "$1" == "-h" ]; then 
	echo "./list.sh <prefix> [all/ns/vpc/fw/mcis]"
	echo "./list.sh cb-aws mcis"
fi
if [ "$1" != "-h" ]; then 
	list;
fi
