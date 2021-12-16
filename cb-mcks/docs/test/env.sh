#!/bin/bash
# -----------------------------------------------------------------
# usage
if [ "$#" -lt 1 ]; then 
	echo "./env.sh [AWS/GCP/AZURE/ALBIABA/TENCENT/OPENSTACK] <credential file>"
	echo "./env.sh AWS ~/.aws/credential"
	exit 0; 
fi


# ------------------------------------------------------------------------------
# parameter

# 1. CSP
if [ "$#" -gt 0 ]; then v_CSP="$1"; else	v_CSP="${CSP}"; fi
if [ "${v_CSP}" = "" ]; then
	read -e -p "Cloud ? [AWS(default) or GCP or AZURE or ALIBABA or TENCENT or OPENSTACK] : "  v_CSP
fi

if [ "${v_CSP}" != "GCP" ] && [ "${v_CSP}" != "AWS" ] && [ "${v_CSP}" != "AZURE" ] && [ "${v_CSP}" != "ALIBABA" ] && [ "${v_CSP}" != "TENCENT" ] && [ "${v_CSP}" != "OPENSTACK" ]; then echo "[ERROR] missing <cloud>"; exit -1;fi

# credential file
if [ "$#" -gt 1 ]; then v_FILE="$2"; else	v_FILE="${CRT_FILE}"; fi
if [ "${v_FILE}" = "" ]; then
	read -e -p "credential file path ? : "  v_FILE
fi
if [ "${v_FILE}" = "" ]; then echo "[ERROR] missing <credential file>"; exit -1;fi

# credential (gcp)
if [ "${v_CSP}" = "GCP" ]; then

	export GCP_PROJECT=$(cat ${v_FILE} | jq ".project_id" | tr -d '"')
	export GCP_PKEY=$(cat ${v_FILE} | jq ".private_key" | tr -d '"')
	export GCP_SA=$(cat ${v_FILE} | jq ".client_email" | tr -d '"')

fi

# credential (aws)
if [ "${v_CSP}" = "AWS" ]; then

	export AWS_KEY="$(head -n 2 ${v_FILE} | tail -n 1 | sed  '/^$/d; s/\r//; s/aws_access_key_id = //g')"
	export AWS_SECRET="$(head -n 3 ${v_FILE} | tail -n 1 | sed  '/^$/d; s/\r//; s/aws_secret_access_key = //g')"

fi

# credential (azure)
if [ "${v_CSP}" = "AZURE" ]; then

	export AZURE_CLIENT_ID="$(cat ${v_FILE} | jq '.clientId' | sed  '/^$/d; s/\r//; s/"//g')"
	export AZURE_CLIENT_SECRET="$(cat ${v_FILE} | jq '.clientSecret' | sed  '/^$/d; s/\r//; s/"//g')"
	export AZURE_TENANT_ID="$(cat ${v_FILE} | jq '.tenantId' | sed  '/^$/d; s/\r//; s/"//g')"
	export AZURE_SUBSCRIPTION_ID="$(cat ${v_FILE} | jq '.subscriptionId' | sed  '/^$/d; s/\r//; s/"//g')"
	 
fi


# credential (alibaba)
if [ "${v_CSP}" = "ALIBABA" ]; then

	export ALIBABA_KEY="$(head -n 2 ${v_FILE} | tail -n 1 | cut -d ',' -f 1)"
	export ALIBABA_SECRET="$(head -n 2 ${v_FILE} | tail -n 1 | cut -d ',' -f 2)"

fi

# credential (tencent)
if [ "${v_CSP}" = "TENCENT" ]; then

	export TENCENT_KEY="$(cat ${v_FILE} | jq '.secretId' | sed  '/^$/d; s/\r//; s/"//g')"
	export TENCENT_SECRET="$(cat ${v_FILE} | jq '.secretKey' | sed  '/^$/d; s/\r//; s/"//g')"

fi

# credential (openstack)
if [ "${v_CSP}" = "OPENSTACK" ]; then

	source ${v_FILE}

fi

# ------------------------------------------------------------------------------
# print info.
echo ""
echo "[Env.]"
echo "GCP"
echo -E "- PROJECT is '${GCP_PROJECT}'"
echo -E "- PKEY    is '${GCP_PKEY}'"
echo -E "- SA      is '${GCP_SA}'"
echo "AWS"
echo -E "- KEY     is '${AWS_KEY}'"
echo -E "- SECRET  is '${AWS_SECRET}'"
echo "AZURE"
echo -E "- CLIENT_ID       is '${AZURE_CLIENT_ID}'"
echo -E "- CLIENT_SECRET   is '${AZURE_CLIENT_SECRET}'"
echo -E "- TENANT_ID       is '${AZURE_TENANT_ID}'"
echo -E "- SUBSCRIPTION_ID is '${AZURE_SUBSCRIPTION_ID}'"
echo "ALIBABA"
echo -E "- KEY     is '${ALIBABA_KEY}'"
echo -E "- SECRET  is '${ALIBABA_SECRET}'"
echo "TENCENT"
echo -E "- KEY     is '${TENCENT_KEY}'"
echo -E "- SECRET  is '${TENCENT_SECRET}'"
echo "OPENSTACK"
echo -E "- IdentityEndpoint  is '${OS_AUTH_URL}'"
echo -E "- Username          is '${OS_USERNAME}'"
echo -E "- Password          is '${OS_PASSWORD}'"
echo -E "- DomainName        is '${OS_USER_DOMAIN_NAME}'"
echo -E "- ProjectID         is '${OS_PROJECT_ID}'"
