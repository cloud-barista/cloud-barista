#!/bin/bash
# ------------------------------------------------------------------------------
# usage
if [ "$#" -lt 1 ]; then 
	echo "./connectioninfo-list.sh [all/driver/credential/region/config]"
	echo "./connectioninfo-list.sh all"
	exit 0
fi

source ./conf.env

# ------------------------------------------------------------------------------
# const



# ------------------------------------------------------------------------------
# variables

# 1. query
if [ "$#" -gt 0 ]; then v_QUERY="$1"; fi
if [ "${v_QUERY}" == "" ]; then 
	read -e -p "Query ? [all/image/spec/ssh/sg/vpc] : "  v_QUERY
fi
if [ "${v_QUERY}" == "" ]; then echo "[ERROR] missing <query>"; exit -1; fi
if [ "${v_QUERY}" == "all" ]; then v_QUERY="driver,credential,region,config"; fi


# ------------------------------------------------------------------------------
# print info.
echo ""
echo "[INFO]"
echo "- Query                      is '${v_QUERY}'"

# ------------------------------------------------------------------------------
# show init result
list() {
	if [[ "${v_QUERY}" == *"driver"* ]]; then echo "DRIVER";     curl -sX GET ${c_URL_SPIDER}/driver							-H "${c_CT}" | jq; fi
	if [[ "${v_QUERY}" == *"credential"* ]]; then echo "CREDENTIAL"; curl -sX GET ${c_URL_SPIDER}/credential			-H "${c_CT}" | jq; fi
	if [[ "${v_QUERY}" == *"region"* ]]; then echo "REGION";     curl -sX GET ${c_URL_SPIDER}/region							-H "${c_CT}" | jq; fi
	if [[ "${v_QUERY}" == *"config"* ]]; then echo "CONFIG";     curl -sX GET ${c_URL_SPIDER}/connectionconfig		-H "${c_CT}" | jq; fi
}

# ------------------------------------------------------------------------------
if [ "$1" != "-h" ]; then 
	echo ""
	echo "------------------------------------------------------------------------------"
	list;
fi
