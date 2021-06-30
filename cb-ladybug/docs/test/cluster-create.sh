#!/bin/bash
# -----------------------------------------------------------------
# usage
if [ "$#" -lt 1 ]; then 
	echo "./cluster-create.sh <namespace> <clsuter name>"
	echo "./cluster-create.sh cb-ladybug-ns cluster-01"
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
# Create a cluster
create() {

	resp=$(curl -sX POST ${c_URL_LADYBUG_NS}/clusters -H "${c_CT}" -d @- <<EOF
	{
		"name": "${v_CLUSTER_NAME}",
		"config": {
			"kubernetes": {
				"networkCni": "kilo",
				"podCidr": "10.244.0.0/16",
				"serviceCidr": "10.96.0.0/12",
				"serviceDnsDomain": "cluster.local"
			}
		},
		"controlPlane": [
			{
				"connection": "config-aws-ap-northeast-1",
				"count": 1,
				"spec": "t2.medium"
			}
		],
		"worker": [
			{
				"connection": "config-gcp-asia-northeast3",
				"count": 1,
				"spec": "n1-standard-2"
			}
		]
	}
EOF
	); echo ${resp} | jq
}


# ------------------------------------------------------------------------------
if [ "$1" != "-h" ]; then 
	echo ""
	echo "------------------------------------------------------------------------------"
	create;
fi
