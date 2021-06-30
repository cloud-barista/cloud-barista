#!/bin/bash
# ------------------------------------------------------------------------------
# usage
if [ "$1" == "-h" ]; then 
	echo "./savekey.sh <namespace> <connection info>"
	echo "./savekey.sh cb-ladybug-ns config-aws-ap-northeast-1"
	exit 0
fi


# ------------------------------------------------------------------------------
# const
c_URL_TUMBLEBUG="http://localhost:1323/tumblebug"
c_CT="Content-Type: application/json"
c_AUTH="Authorization: Basic $(echo -n default:default | base64)"

# ------------------------------------------------------------------------------
# paramter

# 1. namespace
if [ "$#" -gt 0 ]; then v_NAMESPACE="$1"; else	v_NAMESPACE="${NAMESPACE}"; fi
if [ "${v_NAMESPACE}" == "" ]; then 
	read -e -p "Namespace ? : " v_NAMESPACE
fi
if [ "${v_NAMESPACE}" == "" ]; then echo "[ERROR] missing <namespace>"; exit -1; fi

# 2. connection info
if [ "$#" -gt 1 ]; then v_CONFIG="$2"; else	v_CONFIG="${CONNECTION_CONFIG}"; fi
if [ "${v_CONFIG}" == "" ]; then 
	read -e -p "connection info  ? : "  v_CONFIG
fi
if [ "${v_CONFIG}" == "" ]; then echo "[ERROR] missing <connection info>"; exit -1; fi


# variable - name
NM_SSH_KEY="${v_CONFIG/config-/}-sshkey"

c_URL_TUMBLEBUG_NS="${c_URL_TUMBLEBUG}/ns/${v_NAMESPACE}"

# ------------------------------------------------------------------------------
# print info.
echo "[INFO]"
echo "- Namespace                  is '${v_NAMESPACE}'"
echo "- Connection Info            is '${v_CONFIG}'"
echo "- (Name of ssh key)          is '${NM_SSH_KEY}'"


# ------------------------------------------------------------------------------
# get Infrastructure
get() {
	rm -f ${NM_SSH_KEY}.pem
	curl -sX GET ${c_URL_TUMBLEBUG_NS}/resources/sshKey/${NM_SSH_KEY}   -H "${c_AUTH}" -H "${c_CT}" | jq -r ".privateKey" > ${NM_SSH_KEY}.pem
	chmod 400 ${NM_SSH_KEY}.pem
	cat ${NM_SSH_KEY}.pem
}


# ------------------------------------------------------------------------------
if [ "$1" != "-h" ]; then 
	echo ""
	echo "------------------------------------------------------------------------------"
	get;
fi
