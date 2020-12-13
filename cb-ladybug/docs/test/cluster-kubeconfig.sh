#!/bin/bash
# ------------------------------------------------------------------------------
# usage
if [ "$1" == "-h" ]; then 
	echo "./cluster-kubeconfig.sh [GCP/AWS] <cluster name>"
	echo "./cluster-kubeconfig.sh AWS cb-cluster"
	exit 0
fi


# ------------------------------------------------------------------------------
# const

c_URL_LADYBUG="http://localhost:8080/ladybug"
c_CT="Content-Type: application/json"

# ------------------------------------------------------------------------------
# paramter

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

# # 1. PREFIX
# if [ "$#" -gt 0 ]; then v_PREFIX="$1"; else	v_PREFIX="${PREFIX}"; fi

# if [ "${v_PREFIX}" == "" ]; then 
# 	read -e -p "Name prefix ? : "  v_PREFIX
# fi
# if [ "${v_PREFIX}" == "" ]; then echo "[ERROR] missing <prefix>"; exit -1; fi

# 3. Cluster Name
v_CLUSTER="$2"
if [ "${v_CLUSTER}" == "" ]; then read -e -p "Cluster name  ? : "  v_CLUSTER;	fi
if [ "${v_CLUSTER}" == "" ]; then echo "[ERROR] missing <cluster name>"; exit -1; fi


# variable - name
NM_NAMESPACE="${v_PREFIX}-namespace"
c_URL_LADYBUG_NS="${c_URL_LADYBUG}/ns/${NM_NAMESPACE}"

# ------------------------------------------------------------------------------
# print info.
echo "[INFO]"
echo "- Prefix                     is '${v_PREFIX}'"
echo "- Namespace                  is '${NM_NAMESPACE}'"
echo "- (Name of cluster)          is '${v_CLUSTER}'"


# ------------------------------------------------------------------------------
# get Infrastructure
get() {

	rm -f "kubeconfig.yaml"
	# curl -sX GET ${c_URL_LADYBUG_NS}/clusters/${v_CLUSTER} -H "${c_CT}"
	curl -sX GET ${c_URL_LADYBUG_NS}/clusters/${v_CLUSTER} -H "${c_CT}" | jq -r ".clusterConfig" > kubeconfig.yaml
	echo "kubectl get nodes --kubeconfig=./kubeconfig.yaml --insecure-skip-tls-verify=true"
#	chmod 400 ${v_CLUSTER_NAME}.pem
#	cat ${v_CLUSTER_NAME}.pem
}


# ------------------------------------------------------------------------------
if [ "$1" != "-h" ]; then 
	echo ""
	echo "------------------------------------------------------------------------------"
	get;
fi
