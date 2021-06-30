#!/bin/bash
# -----------------------------------------------------------------
# usage
if [ "$#" -lt 1 ]; then 
	echo "./cluster-list.sh <namespace>"
	echo "./cluster-list.sh cb-ladybug-ns"
	exit 0; 
fi


# ------------------------------------------------------------------------------
# const
c_URL_LADYBUG="http://localhost:8080/ladybug"
c_CT="Content-Type: application/json"


# -----------------------------------------------------------------
# parameter

# 1. namespace
if [ "$#" -gt 0 ]; then v_NAMESPACE="$1"; else	v_NAMESPACE="${NAMESPACE}"; fi
if [ "${v_NAMESPACE}" == "" ]; then 
	read -e -p "Namespace ? : " v_NAMESPACE
fi
if [ "${v_NAMESPACE}" == "" ]; then echo "[ERROR] missing <namespace>"; exit -1; fi

c_URL_LADYBUG_NS="${c_URL_LADYBUG}/ns/${v_NAMESPACE}"


# ------------------------------------------------------------------------------
# print info.
echo ""
echo "[INFO]"
echo "- Namespace			             is '${v_NAMESPACE}'"


# ------------------------------------------------------------------------------
# list
list() {
	curl -sX GET ${c_URL_LADYBUG_NS}/clusters -H "${c_CT}" | jq;
}


# ------------------------------------------------------------------------------
if [ "$1" != "-h" ]; then 
	echo ""
	echo "------------------------------------------------------------------------------"
	list;
fi
