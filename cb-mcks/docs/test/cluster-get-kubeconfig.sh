#!/bin/bash
# ------------------------------------------------------------------------------
# usage
if [ "$#" -lt 1 ]; then 
	echo "./cluster-get-kubeconfig.sh <namespace> <cluster name>"
	echo "./cluster-get-kubeconfig.sh cb-mcks-ns cluster-01"
	exit 0
fi

source ./conf.env

# ------------------------------------------------------------------------------
# const


# ------------------------------------------------------------------------------
# paramter

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


c_URL_MCKS_NS="${c_URL_MCKS}/ns/${v_NAMESPACE}"

# ------------------------------------------------------------------------------
# print info.
echo ""
echo "[INFO]"
echo "- Namespace                  is '${v_NAMESPACE}'"
echo "- Cluster name               is '${v_CLUSTER_NAME}'"


# ------------------------------------------------------------------------------
# get Infrastructure
get() {

	if [ "$MCKS_CALL_METHOD" == "REST" ]; then
		
		rm -f "kubeconfig.yaml"
		curl -sX GET ${c_URL_MCKS_NS}/clusters/${v_CLUSTER_NAME} -H "${c_CT}" | jq -r ".clusterConfig" > kubeconfig.yaml

		echo "export KUBECONFIG=$(pwd)/kubeconfig.yaml"
		echo "kubectl get nodes"	

	elif [ "$MCKS_CALL_METHOD" == "GRPC" ]; then

		rm -f "kubeconfig.yaml"
		$APP_ROOT/src/grpc-api/cbadm/cbadm cluster get --config $APP_ROOT/src/grpc-api/cbadm/grpc_conf.yaml -o json --ns ${v_NAMESPACE} --cluster ${v_CLUSTER_NAME} | jq -r ".clusterConfig" > kubeconfig.yaml

		echo "export KUBECONFIG=$(pwd)/kubeconfig.yaml"
		echo "kubectl get nodes"	
		
	else
		echo "[ERROR] missing MCKS_CALL_METHOD"; exit -1;
	fi
	
}


# ------------------------------------------------------------------------------
if [ "$1" != "-h" ]; then 
	echo ""
	echo "------------------------------------------------------------------------------"
	get;
fi
