#!/bin/bash
# ------------------------------------------------------------------------------
# usage
if [ "$#" -lt 1 ]; then 
	echo "./connectioninfo-delete.sh [config/region/credential/driver] <name>"
	echo "./connectioninfo-delete.sh config config-aws-ap-northeast-1"
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
	read -e -p "Query ? [config/region/credential/driver/ns] : "  v_QUERY
fi
if [ "${v_QUERY}" == "" ]; then echo "[ERROR] missing <query>"; exit -1; fi


# 2. name
if [ "$#" -gt 1 ]; then v_NAME="$2"; fi
if [ "${v_NAME}" == "" ]; then 
	read -e -p "Name ? : " v_NAME
fi
if [ "${v_NAME}" == "" ]; then echo "[ERROR] missing <name>"; exit -1; fi


# ------------------------------------------------------------------------------
# print info.
echo ""
echo "[INFO]"
echo "- Name                       is '${v_NAME}'"
echo "- Query                      is '${v_QUERY}'"


# ------------------------------------------------------------------------------
# delete
delete() {

	# driver
	if [[ "${v_QUERY}" == *"driver"* ]]; then	echo "@_DRIVER_@";		curl -sX DELETE ${c_URL_SPIDER}/driver/${v_NAME}  -H "${c_CT}" -o /dev/null -w "DRIVER.delete():%{http_code}\n"; fi

	# credential
	if [[ "${v_QUERY}" == *"credential"* ]]; then	echo "@_CREDENTIAL_@";		curl -sX DELETE ${c_URL_SPIDER}/credential/${v_NAME} -H "${c_CT}" -o /dev/null -w "CREDENTIAL.delete():%{http_code}\n"; fi

	# region
	if [[ "${v_QUERY}" == *"region"* ]]; then	echo "@_REGION_@";		curl -sX DELETE ${c_URL_SPIDER}/region/${v_NAME} -H "${c_CT}" -o /dev/null -w "REGION.delete():%{http_code}\n"; fi

	# config
	if [[ "${v_QUERY}" == *"config"* ]]; then	echo "@_CONFIG_@";		curl -sX DELETE ${c_URL_SPIDER}/connectionconfig/${v_NAME} -H "${c_AUTH}" -H "${c_CT}" -o /dev/null -w "CONFIG.delete():%{http_code}\n"; fi

}


# ------------------------------------------------------------------------------
# show 
show() {
	if [[ "${v_QUERY}" == *"driver"* ]]; then echo "DRIVER";     curl -sX GET ${c_URL_SPIDER}/driver							-H "${c_CT}" | jq; fi
	if [[ "${v_QUERY}" == *"credential"* ]]; then echo "CREDENTIAL"; curl -sX GET ${c_URL_SPIDER}/credential			-H "${c_CT}" | jq; fi
	if [[ "${v_QUERY}" == *"region"* ]]; then	echo "REGION";     curl -sX GET ${c_URL_SPIDER}/region							-H "${c_CT}" | jq; fi
	if [[ "${v_QUERY}" == *"config"* ]]; then	echo "CONFIG";     curl -sX GET ${c_URL_SPIDER}/connectionconfig		-H "${c_CT}" | jq; fi
}

# ------------------------------------------------------------------------------
if [ "$1" != "-h" ]; then 
	echo ""
	echo "------------------------------------------------------------------------------"
	delete;	show;
fi
