#!/bin/sh

# ------------------------------------------------------------------------------
# env.
# ${CRT_FILE} : name

# ------------------------------------------------------------------------------
# parameter

# 1. CSP
if [ "$#" -gt 0 ]; then v_CSP="$1"; else	v_CSP="${CSP}"; fi
if [ "${v_CSP}" = "" ]; then
	read -e -p "Cloud ? [AWS(default) or GCP] : "  v_CSP
fi

# credential file
if [ "$#" -gt 1 ]; then v_FILE="$2"; else	v_FILE="${CRT_FILE}"; fi
if [ "${v_FILE}" = "" ]; then
	read -e -p "credential file path ? : "  v_FILE
fi
if [ "${v_FILE}" = "" ]; then echo "[ERROR] missing <credential file>"; exit -1;fi

# credential (gcp)
if [ "${v_CSP}" = "GCP" ]; then

	export PROJECT=$(cat ${v_FILE} | jq -r ".project_id")
	export PKEY=$(cat ${v_FILE} | jq -r ".private_key" | while read line; do	if [[ "$line" != "" ]]; then	echo -n "$line\n";	fi; done )
	export SA=$(cat ${v_FILE} | jq -r ".client_email")

fi

# credential (aws)
if [ "${v_CSP}" = "AWS" ]; then

	export KEY="$(head -n 3 ${v_FILE} | tail -n 1 | sed  '/^$/d; s/\r//; s/aws_access_key_id = //g')"
	export SECRET="$(head -n 2 ${v_FILE} | tail -n 1 | sed  '/^$/d; s/\r//; s/aws_secret_access_key = //g')"

fi


# ------------------------------------------------------------------------------
# print info.
echo ""
echo "[Env.]"
echo "GCP"
echo "- PROJECT is '${PROJECT}'"
echo "- PKEY    is '${PKEY}'"
echo "- SA      is '${SA}'"
echo "AWS"
echo "- KEY     is '${KEY}'"
echo "- SECRET  is '${SECRET}'"
