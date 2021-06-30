#!/bin/bash
# -----------------------------------------------------------------
# usage
if [ "$#" -lt 1 ]; then 
	echo "./cluster-delete.sh <namespace> <clsuter name>"
	echo "./cluster-delete.sh cb-ladybug-ns cb-cluster"
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

# 2. Cluster Name
if [ "$#" -gt 1 ]; then v_CLUSTER_NAME="$2"; else	v_CLUSTER_NAME="${CLUSTER_NAME}"; fi
if [ "${v_CLUSTER_NAME}" == "" ]; then 
	read -e -p "Cluster name  ? : "  v_CLUSTER_NAME
fi
if [ "${v_CLUSTER_NAME}" == "" ]; then echo "[ERROR] missing <cluster name>"; exit -1; fi


c_URL_LADYBUG_NS="${c_URL_LADYBUG}/ns/${v_NAMESPACE}"


# ------------------------------------------------------------------------------
# print info.
echo ""
echo "[INFO]"
echo "- Namespace                  is '${v_NAMESPACE}'"
echo "- Cluster name               is '${v_CLUSTER_NAME}'"


# ------------------------------------------------------------------------------
# Delete a cluster
delete() {

	curl -sX DELETE ${c_URL_LADYBUG_NS}/clusters/${v_CLUSTER_NAME}    -H "${c_CT}" | jq;

}

# ------------------------------------------------------------------------------
if [ "$1" != "-h" ]; then 
	echo ""
	echo "------------------------------------------------------------------------------"
	delete;
fi
