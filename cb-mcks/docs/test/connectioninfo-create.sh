#!/bin/bash
# ------------------------------------------------------------------------------
# usage
if [ "$#" -lt 1 ]; then 
	echo "./connectioninfo-create.sh [AWS/GCP/AZURE/ALIBABA/TENCENT/OPENSTACK] <option>"
	echo "./connectioninfo-create.sh GCP"
	echo "./connectioninfo-create.sh AWS add"
	exit 0
fi

source ./conf.env

# ------------------------------------------------------------------------------
# const



# ------------------------------------------------------------------------------
# variables

# 1. CSP
if [ "$#" -gt 0 ]; then v_CSP="$1"; else	v_CSP="${CSP}"; fi
if [ "${v_CSP}" == "" ]; then 
	read -e -p "Cloud ? [AWS(default) or GCP or AZURE or ALIBABA or TENCENT or OPENSTACK] : "  v_CSP
fi

if [ "${v_CSP}" == "" ]; then v_CSP="AWS"; fi
if [ "${v_CSP}" != "GCP" ] && [ "${v_CSP}" != "AWS" ] && [ "${v_CSP}" != "AZURE" ] && [ "${v_CSP}" != "ALIBABA" ] && [ "${v_CSP}" != "TENCENT" ] && [ "${v_CSP}" != "OPENSTACK" ]; then echo "[ERROR] missing <cloud>"; exit -1;fi

v_CSP_LOWER="$(echo ${v_CSP} | tr [:upper:] [:lower:])"

# 2. option
if [ "$#" -gt 1 ]; then v_OPTION="$2"; else v_OPTION=""; fi
if [ "${v_OPTION}" != "" ] && [ "${v_OPTION}" != "add" ]; then echo "[ERROR] not valid <option>"; v_OPTION="" ;fi

# driver
if [ "${v_CSP}" == "GCP" ]; then 
	v_DRIVER="${c_GCP_DRIVER}"
elif [ "${v_CSP}" == "AWS" ]; then 
	v_DRIVER="${c_AWS_DRIVER}"
elif [ "${v_CSP}" == "AZURE" ]; then 
	v_DRIVER="${c_AZURE_DRIVER}"
elif [ "${v_CSP}" == "ALIBABA" ]; then 
	v_DRIVER="${c_ALIBABA_DRIVER}"
elif [ "${v_CSP}" == "TENCENT" ]; then 
	v_DRIVER="${c_TENCENT_DRIVER}"
elif [ "${v_CSP}" == "OPENSTACK" ]; then 
	v_DRIVER="${c_OPENSTACK_DRIVER}"
fi

if [ "${v_OPTION}" != "add" ]; then 

	# credential
	# GCP
	if [ "${v_CSP}" == "GCP" ]; then 

		# Project
		v_GCP_PROJECT="${GCP_PROJECT}"
		if [ "${v_GCP_PROJECT}" == "" ]; then 
			read -e -p "Project ? [예:kore3-etri-cloudbarista] : "  v_GCP_PROJECT
			if [ "${v_GCP_PROJECT}" == "" ]; then echo "[ERROR] missing gcp <project_id>"; exit -1;fi
		fi

		# private key
		v_GCP_PKEY="${GCP_PKEY}"
		if [ "${v_GCP_PKEY}" == "" ]; then 
			read -e -p "Private Key ? [예:-----BEGIN PRIVATE KEY-----\n....] : "  v_GCP_PKEY
			if [ "${v_GCP_PKEY}" == "" ]; then echo "[ERROR] missing gcp <private_key>"; exit -1;fi
		fi

		# system account
		v_GCP_SA="${GCP_SA}"
		if [ "${v_GCP_SA}" == "" ]; then 
			read -e -p "Service account (client email) ? [예:331829771895-compute@developer.gserviceaccount.com] : "  v_GCP_SA
			if [ "${v_GCP_SA}" == "" ]; then echo "[ERROR] missing gcp <client_email>"; exit -1;fi
		fi

		# region
		v_REGION="${GCP_REGION}"
		if [ "${v_REGION}" == "" ]; then 
			read -e -p "region ? [예:asia-northeast3] : "  v_REGION
			if [ "${v_REGION}" == "" ]; then echo "[ERROR] missing region"; exit -1;fi
		fi

		# zone
		v_ZONE="${GCP_ZONE}"
		if [ "${v_ZONE}" == "" ]; then 
			read -e -p "zone ? [예:asia-northeast3-a] : "  v_ZONE
			if [ "${v_ZONE}" == "" ]; then v_ZONE="${v_REGION}-a";fi
		fi
	fi

	# AWS
	if [ "${v_CSP}" == "AWS" ]; then 

		v_AWS_ACCESS_KEY="${AWS_KEY}"
		if [ "${v_AWS_ACCESS_KEY}" == "" ]; then 
			read -e -p "Access Key ? [예:AH24UUA2ZGNOP6DKKIA6] : "  v_AWS_ACCESS_KEY
			if [ "${v_AWS_ACCESS_KEY}" == "" ]; then echo "[ERROR] missing <aws_access_key_id>"; exit -1;fi
		fi

		v_AWS_SECRET="${AWS_SECRET}"
		if [ "${v_AWS_SECRET}" == "" ]; then 
			read -e -p "Access-key Secret ? [예:y76ZWz6A/vwqGanDAI926TTPCJrrMo1VbPOh8X7K] : "  v_AWS_SECRET
			if [ "${v_AWS_SECRET}" == "" ]; then echo "[ERROR] missing <aws_secret_access_key>"; exit -1;fi
		fi

		# region
		v_REGION="${AWS_REGION}"
		if [ "${v_REGION}" == "" ]; then 
			read -e -p "region ? [예:asia-northeast3] : "  v_REGION
			if [ "${v_REGION}" == "" ]; then echo "[ERROR] missing region"; exit -1;fi
		fi

		# zone
		v_ZONE="${AWS_ZONE}"
		if [ "${v_ZONE}" == "" ]; then 
			read -e -p "zone ? [예:asia-northeast3-a] : "  v_ZONE
			if [ "${v_ZONE}" == "" ]; then v_ZONE="${v_REGION}a";fi
		fi
	fi

	# AZURE
	if [ "${v_CSP}" == "AZURE" ]; then 

		# client id
		v_AZURE_CLIENT_ID="${AZURE_CLIENT_ID}"
		if [ "${v_AZURE_CLIENT_ID}" == "" ]; then 
			read -e -p "client id ? [예:123445-dfef-s9df-9292-c9d9d01030] : "  v_AZURE_CLIENT_ID
			if [ "${v_AZURE_CLIENT_ID}" == "" ]; then echo "[ERROR] missing <azure_client_id>"; exit -1;fi
		fi	

		# client secret
		v_AZURE_CLIENT_SECRET="${AZURE_CLIENT_SECRET}"
		if [ "${v_AZURE_CLIENT_SECRET}" == "" ]; then 
			read -e -p "client secret ? [예:239DLKJFSJ=DFLKJSFK-FDSLKJFS0d] : "  v_AZURE_CLIENT_SECRET
			if [ "${v_AZURE_CLIENT_SECRET}" == "" ]; then echo "[ERROR] missing <azure_client_secret>"; exit -1;fi
		fi	

		# tenant id
		v_AZURE_TENANT_ID="${AZURE_TENANT_ID}"
		if [ "${v_AZURE_TENANT_ID}" == "" ]; then 
			read -e -p "tenant id ? [예:123445-dfef-s9df-9292-c9d9d01030] : "  v_AZURE_TENANT_ID
			if [ "${v_AZURE_TENANT_ID}" == "" ]; then echo "[ERROR] missing <azure_tenant_id>"; exit -1;fi
		fi	

		# subscription id
		v_AZURE_SUBSCRIPTION_ID="${AZURE_SUBSCRIPTION_ID}"
		if [ "${v_AZURE_SUBSCRIPTION_ID}" == "" ]; then 
			read -e -p "subscription id ? [예:123445-dfef-s9df-9292-c9d9d01030] : "  v_AZURE_SUBSCRIPTION_ID
			if [ "${v_AZURE_SUBSCRIPTION_ID}" == "" ]; then echo "[ERROR] missing <azure_subscription_id>"; exit -1;fi
		fi	

		# region
		v_REGION="${AZURE_REGION}"
		if [ "${v_REGION}" == "" ]; then 
			read -e -p "region ? [예:asia-northeast3] : "  v_REGION
			if [ "${v_REGION}" == "" ]; then echo "[ERROR] missing region"; exit -1;fi
		fi

		# resource group
		v_RESOURCE_GROUP="${AZURE_RESOURCE_GROUP}"
		if [ "${v_RESOURCE_GROUP}" == "" ]; then 
			read -e -p "resource group ? [예:cb-mcksRG] : "  v_RESOURCE_GROUP
			if [ "${v_RESOURCE_GROUP}" == "" ]; then echo "[ERROR] missing resource group"; exit -1;fi
		fi
	fi


	# ALIBABA
	if [ "${v_CSP}" == "ALIBABA" ]; then 

		v_ALIBABA_ACCESS_KEY="${ALIBABA_KEY}"
		if [ "${v_ALIBABA_ACCESS_KEY}" == "" ]; then 
			read -e -p "Access Key ? [예:AH24UUA2ZGNOP6DKKIA6] : "  v_ALIBABA_ACCESS_KEY
			if [ "${v_ALIBABA_ACCESS_KEY}" == "" ]; then echo "[ERROR] missing <alibaba_access_key_id>"; exit -1;fi
		fi

		v_ALIBABA_SECRET="${ALIBABA_SECRET}"
		if [ "${v_ALIBABA_SECRET}" == "" ]; then 
			read -e -p "Access-key Secret ? [예:y76ZWz6A/vwqGanDAI926TTPCJrrMo1VbPOh8X7K] : "  v_ALIBABA_SECRET
			if [ "${v_ALIBABA_SECRET}" == "" ]; then echo "[ERROR] missing <alibaba_access_key_secret>"; exit -1;fi
		fi

		# region
		v_REGION="${ALIBABA_REGION}"
		if [ "${v_REGION}" == "" ]; then 
			read -e -p "region ? [예:asia-northeast3] : "  v_REGION
			if [ "${v_REGION}" == "" ]; then echo "[ERROR] missing region"; exit -1;fi
		fi

		# zone
		v_ZONE="${ALIBABA_ZONE}"
		if [ "${v_ZONE}" == "" ]; then 
			read -e -p "zone ? [예:asia-northeast3-a] : "  v_ZONE
			if [ "${v_ZONE}" == "" ]; then v_ZONE="${v_REGION}a";fi
		fi
	fi

	# TENCENT
	if [ "${v_CSP}" == "TENCENT" ]; then 

		v_TENCENT_ACCESS_KEY="${TENCENT_KEY}"
		if [ "${v_TENCENT_ACCESS_KEY}" == "" ]; then 
			read -e -p "Access Key ? [예:AH24UUA2ZGNOP6DKKIA6] : "  v_TENCENT_ACCESS_KEY
			if [ "${v_TENCENT_ACCESS_KEY}" == "" ]; then echo "[ERROR] missing <tencent_access_key_id>"; exit -1;fi
		fi

		v_TENCENT_SECRET="${TENCENT_SECRET}"
		if [ "${v_TENCENT_SECRET}" == "" ]; then 
			read -e -p "Access-key Secret ? [예:y76ZWz6A/vwqGanDAI926TTPCJrrMo1VbPOh8X7K] : "  v_TENCENT_SECRET
			if [ "${v_TENCENT_SECRET}" == "" ]; then echo "[ERROR] missing <tencent_access_key_secret>"; exit -1;fi
		fi

		# region
		v_REGION="${TENCENT_REGION}"
		if [ "${v_REGION}" == "" ]; then 
			read -e -p "region ? [예:ap-seoul] : "  v_REGION
			if [ "${v_REGION}" == "" ]; then echo "[ERROR] missing region"; exit -1;fi
		fi

		# zone
		v_ZONE="${TENCENT_ZONE}"
		if [ "${v_ZONE}" == "" ]; then 
			read -e -p "zone ? [예:ap-seoul-1] : "  v_ZONE
			if [ "${v_ZONE}" == "" ]; then v_ZONE="${v_REGION}a";fi
		fi
	fi

	# OPENSTACK
	if [ "${v_CSP}" == "OPENSTACK" ]; then 

		v_OPENSTACK_ENDPOINT="${OS_AUTH_URL}"
		if [ "${v_OPENSTACK_ENDPOINT}" == "" ]; then 
			read -e -p "Identity Endpoint ? [예:http://123.456.789.123:5000/v3] : "  v_OPENSTACK_ENDPOINT
			if [ "${v_OPENSTACK_ENDPOINT}" == "" ]; then echo "[ERROR] missing <openstack identity endpoint>"; exit -1;fi
		fi

		v_OPENSTACK_USERNAME="${OS_USERNAME}"
		if [ "${v_OPENSTACK_USERNAME}" == "" ]; then 
			read -e -p "Username ? [예:mcks] : "  v_OPENSTACK_USERNAME
			if [ "${v_OPENSTACK_USERNAME}" == "" ]; then echo "[ERROR] missing <openstack username>"; exit -1;fi
		fi

		v_OPENSTACK_PASSWORD="${OS_PASSWORD}"
		if [ "${v_OPENSTACK_PASSWORD}" == "" ]; then 
			read -e -p "Password ? [예:asdfqwer12] : "  v_OPENSTACK_PASSWORD
			if [ "${v_OPENSTACK_PASSWORD}" == "" ]; then echo "[ERROR] missing <openstack password>"; exit -1;fi
		fi

		v_OPENSTACK_DOMAINNAME="${OS_USER_DOMAIN_NAME}"
		if [ "${v_OPENSTACK_DOMAINNAME}" == "" ]; then 
			read -e -p "DomainName ? [예:default] : "  v_OPENSTACK_DOMAINNAME
			if [ "${v_OPENSTACK_DOMAINNAME}" == "" ]; then echo "[ERROR] missing <openstack domainname>"; exit -1;fi
		fi

		v_OPENSTACK_PROJECTID="${OS_PROJECT_ID}"
		if [ "${v_OPENSTACK_PROJECTID}" == "" ]; then 
			read -e -p "ProjectID ? [예:kdjf1k12jkdjf2kjskjf] : "  v_OPENSTACK_PROJECTID
			if [ "${v_OPENSTACK_PROJECTID}" == "" ]; then echo "[ERROR] missing <openstack projectid>"; exit -1;fi
		fi

		# region
		v_REGION="${OS_REGION}"
		if [ "${v_REGION}" == "" ]; then 
			read -e -p "region ? [예:ap-seoul] : "  v_REGION
			if [ "${v_REGION}" == "" ]; then echo "[ERROR] missing region"; exit -1;fi
		fi

		# zone
		v_ZONE="${OS_ZONE}"
		if [ "${v_ZONE}" == "" ]; then 
			read -e -p "zone ? [예:ap-seoul-1] : "  v_ZONE
			if [ "${v_ZONE}" == "" ]; then v_ZONE="${v_REGION}a";fi
		fi
	fi

fi

v_REGION_LOWER="$(echo ${v_REGION} | tr [:upper:] [:lower:])"

NM_CREDENTIAL="credential-${v_CSP_LOWER}"
NM_REGION="region-${v_CSP_LOWER}-${v_REGION_LOWER}"
NM_CONFIG="config-${v_CSP_LOWER}-${v_REGION_LOWER}"

# ------------------------------------------------------------------------------
# print info.
echo ""
echo "[INFO]"
echo "- Cloud                      is '${v_CSP}'"
echo "- Driver                     is '${v_DRIVER}'"
echo "- Region                     is '${v_REGION}'"
if [ "${v_CSP}" == "GCP" ]; then 
	echo "- Zone                       is '${v_ZONE}'"
	echo "- Project                    is '${v_GCP_PROJECT}'"
	echo "- private key                is '${v_GCP_PKEY}'"
	echo "- Service account            is '${v_GCP_SA}'"
elif [ "${v_CSP}" == "AWS" ]; then 
	echo "- Zone                       is '${v_ZONE}'"
 	echo "- aws_access_key_id          is '${v_AWS_ACCESS_KEY}'"
	echo "- aws_secret_access_key      is '${v_AWS_SECRET}'"
elif [ "${v_CSP}" == "AZURE" ]; then 
	echo "- Resource Group             is '${v_RESOURCE_GROUP}'"
	echo "- azure_client_id            is '${v_AZURE_CLIENT_ID}'"
	echo "- azure_client_secret        is '${v_AZURE_CLIENT_SECRET}'"
	echo "- azure_tenant_id            is '${v_AZURE_TENANT_ID}'"
	echo "- azure_subscription_id      is '${v_AZURE_SUBSCRIPTION_ID}'"
elif [ "${v_CSP}" == "ALIBABA" ]; then 
	echo "- Zone                       is '${v_ZONE}'"
 	echo "- alibaba_access_key_id      is '${v_ALIBABA_ACCESS_KEY}'"
	echo "- alibaba_access_key_secret  is '${v_ALIBABA_SECRET}'"
elif [ "${v_CSP}" == "TENCENT" ]; then 
	echo "- Zone                       is '${v_ZONE}'"
 	echo "- tencent_access_key_id      is '${v_TENCENT_ACCESS_KEY}'"
	echo "- tencent_access_key_secret  is '${v_TENCENT_SECRET}'"
elif [ "${v_CSP}" == "OPENSTACK" ]; then 
	echo "- Zone                        is '${v_ZONE}'"
 	echo "- openstack_identity_endpoint is '${v_OPENSTACK_ENDPOINT}'"
	echo "- openstack_username  			  is '${v_OPENSTACK_USERNAME}'"
	echo "- openstack_password  			  is '${v_OPENSTACK_PASSWORD}'"
	echo "- openstack_domainname			  is '${v_OPENSTACK_DOMAINNAME}'"
	echo "- openstack_projectid		 	    is '${v_OPENSTACK_PROJECTID}'"
fi
echo "- (Name of credential)       is '${NM_CREDENTIAL}'"
echo "- (Name of region)           is '${NM_REGION}'"
echo "- (Name of Connection Info.) is '${NM_CONFIG}'"


# ------------------------------------------------------------------------------
# Configuration Spider
create() {

if [ "${v_OPTION}" != "add" ]; then 

		# driver
		curl -sX DELETE ${c_URL_SPIDER}/driver/${v_DRIVER}  -H "${c_CT}" -o /dev/null -w "DRIVER.delete():%{http_code}\n"
		curl -sX POST   ${c_URL_SPIDER}/driver              -H "${c_CT}" -o /dev/null -w "DRIVER.regist():%{http_code}\n" -d @- <<EOF
		{
		"DriverName"        : "${v_DRIVER}",
		"ProviderName"      : "${v_CSP}",
		"DriverLibFileName" : "${v_DRIVER}.so"
		}
EOF

		# credential
		if [ "${v_CSP}" == "GCP" ]; then
			curl -sX DELETE ${c_URL_SPIDER}/credential/${NM_CREDENTIAL} -H "${c_CT}" -o /dev/null -w "CREDENTIAL.delete():%{http_code}\n"
			curl -sX POST   ${c_URL_SPIDER}/credential                  -H "${c_CT}" -o /dev/null -w "CREDENTIAL.regist():%{http_code}\n" -d @- <<EOF
			{
			"CredentialName"   : "${NM_CREDENTIAL}",
			"ProviderName"     : "${v_CSP}",
			"KeyValueInfoList" : [
				{"Key" : "ClientEmail", "Value" : "${v_GCP_SA}"},
				{"Key" : "ProjectID",   "Value" : "${v_GCP_PROJECT}"},
				{"Key" : "PrivateKey",  "Value" : "${v_GCP_PKEY}"}
			]
			}
EOF
		elif [ "${v_CSP}" == "AWS" ]; then
			curl -sX DELETE ${c_URL_SPIDER}/credential/${NM_CREDENTIAL} -H "${c_CT}" -o /dev/null -w "CREDENTIAL.delete():%{http_code}\n"
			curl -sX POST   ${c_URL_SPIDER}/credential                  -H "${c_CT}" -o /dev/null -w "CREDENTIAL.regist():%{http_code}\n" -d @- <<EOF
			{
			"CredentialName"   : "${NM_CREDENTIAL}",
			"ProviderName"     : "${v_CSP}",
			"KeyValueInfoList" : [
				{"Key" : "ClientId",       "Value" : "${v_AWS_ACCESS_KEY}"},
				{"Key" : "ClientSecret",   "Value" : "${v_AWS_SECRET}"}
			]
			}
EOF
		elif [ "${v_CSP}" == "AZURE" ]; then
			curl -sX DELETE ${c_URL_SPIDER}/credential/${NM_CREDENTIAL} -H "${c_CT}" -o /dev/null -w "CREDENTIAL.delete():%{http_code}\n"
			curl -sX POST   ${c_URL_SPIDER}/credential                  -H "${c_CT}" -o /dev/null -w "CREDENTIAL.regist():%{http_code}\n" -d @- <<EOF
			{
			"CredentialName"   : "${NM_CREDENTIAL}",
			"ProviderName"     : "${v_CSP}",
			"KeyValueInfoList" : [
				{"Key" : "ClientId",        "Value" : "${v_AZURE_CLIENT_ID}"},
				{"Key" : "ClientSecret",    "Value" : "${v_AZURE_CLIENT_SECRET}"},
				{"Key" : "TenantId",        "Value" : "${v_AZURE_TENANT_ID}"},
				{"Key" : "SubscriptionId",  "Value" : "${v_AZURE_SUBSCRIPTION_ID}"}
			]
			}
EOF
		elif [ "${v_CSP}" == "ALIBABA" ]; then
			curl -sX DELETE ${c_URL_SPIDER}/credential/${NM_CREDENTIAL} -H "${c_CT}" -o /dev/null -w "CREDENTIAL.delete():%{http_code}\n"
			curl -sX POST   ${c_URL_SPIDER}/credential                  -H "${c_CT}" -o /dev/null -w "CREDENTIAL.regist():%{http_code}\n" -d @- <<EOF
			{
			"CredentialName"   : "${NM_CREDENTIAL}",
			"ProviderName"     : "${v_CSP}",
			"KeyValueInfoList" : [
				{"Key" : "ClientId",       "Value" : "${v_ALIBABA_ACCESS_KEY}"},
				{"Key" : "ClientSecret",   "Value" : "${v_ALIBABA_SECRET}"}
			]
			}
EOF
		elif [ "${v_CSP}" == "TENCENT" ]; then
			curl -sX DELETE ${c_URL_SPIDER}/credential/${NM_CREDENTIAL} -H "${c_CT}" -o /dev/null -w "CREDENTIAL.delete():%{http_code}\n"
			curl -sX POST   ${c_URL_SPIDER}/credential                  -H "${c_CT}" -o /dev/null -w "CREDENTIAL.regist():%{http_code}\n" -d @- <<EOF
			{
			"CredentialName"   : "${NM_CREDENTIAL}",
			"ProviderName"     : "${v_CSP}",
			"KeyValueInfoList" : [
				{"Key" : "ClientId",       "Value" : "${v_TENCENT_ACCESS_KEY}"},
				{"Key" : "ClientSecret",   "Value" : "${v_TENCENT_SECRET}"}
			]
			}
EOF
		elif [ "${v_CSP}" == "OPENSTACK" ]; then
			curl -sX DELETE ${c_URL_SPIDER}/credential/${NM_CREDENTIAL} -H "${c_CT}" -o /dev/null -w "CREDENTIAL.delete():%{http_code}\n"
			curl -sX POST   ${c_URL_SPIDER}/credential                  -H "${c_CT}" -o /dev/null -w "CREDENTIAL.regist():%{http_code}\n" -d @- <<EOF
			{
			"CredentialName"   : "${NM_CREDENTIAL}",
			"ProviderName"     : "${v_CSP}",
			"KeyValueInfoList" : [
				{"Key" : "IdentityEndpoint",	"Value" : "${v_OPENSTACK_ENDPOINT}"},
				{"Key" : "Username",    			"Value" : "${v_OPENSTACK_USERNAME}"},
				{"Key" : "Password",					"Value" : "${v_OPENSTACK_PASSWORD}"},
				{"Key" : "DomainName",				"Value" : "${v_OPENSTACK_DOMAINNAME}"},
				{"Key" : "ProjectID",					"Value" : "${v_OPENSTACK_PROJECTID}"}
			]
			}
EOF
		fi

fi

	# region
	if [ "${v_CSP}" == "AZURE" ]; then
		curl -sX DELETE ${c_URL_SPIDER}/region/${NM_REGION} -H "${c_CT}" -o /dev/null -w "REGION.delete():%{http_code}\n"
		curl -sX POST   ${c_URL_SPIDER}/region              -H "${c_CT}" -o /dev/null -w "REGION.regist():%{http_code}\n" -d @- <<EOF
		{
		"RegionName"       : "${NM_REGION}",
		"ProviderName"     : "${v_CSP}", 
		"KeyValueInfoList" : [
			{"Key" : "location", "Value" : "${v_REGION}"},
			{"Key" : "ResourceGroup", "Value" : "${v_RESOURCE_GROUP}"}
		]
		}
EOF
	else
		curl -sX DELETE ${c_URL_SPIDER}/region/${NM_REGION} -H "${c_CT}" -o /dev/null -w "REGION.delete():%{http_code}\n"
		curl -sX POST   ${c_URL_SPIDER}/region              -H "${c_CT}" -o /dev/null -w "REGION.regist():%{http_code}\n" -d @- <<EOF
		{
		"RegionName"       : "${NM_REGION}",
		"ProviderName"     : "${v_CSP}", 
		"KeyValueInfoList" : [
			{"Key" : "Region", "Value" : "${v_REGION}"},
			{"Key" : "Zone",   "Value" : "${v_ZONE}"}
		]
		}
EOF
	fi

	# config
	curl -sX DELETE ${c_URL_SPIDER}/connectionconfig/${NM_CONFIG} -H "${c_CT}" -o /dev/null -w "CONFIG.delete():%{http_code}\n"
	curl -sX POST   ${c_URL_SPIDER}/connectionconfig              -H "${c_CT}" -o /dev/null -w "CONFIG.regist():%{http_code}\n" -d @- <<EOF
	{
	"ConfigName"     : "${NM_CONFIG}",
	"ProviderName"   : "${v_CSP}", 
	"DriverName"     : "${v_DRIVER}", 
	"CredentialName" : "${NM_CREDENTIAL}", 
	"RegionName"     : "${NM_REGION}"
	}
EOF

}


# ------------------------------------------------------------------------------
# show init result
show() {
	echo "DRIVER";     curl -sX GET ${c_URL_SPIDER}/driver/${v_DRIVER}					-H "${c_CT}" | jq
	echo "CREDENTIAL"; curl -sX GET ${c_URL_SPIDER}/credential/${NM_CREDENTIAL}			-H "${c_CT}" | jq
	echo "REGION";     curl -sX GET ${c_URL_SPIDER}/region/${NM_REGION}					-H "${c_CT}" | jq
	echo "CONFIG";     curl -sX GET ${c_URL_SPIDER}/connectionconfig/${NM_CONFIG}		-H "${c_CT}" | jq
}

# ------------------------------------------------------------------------------
if [ "$1" != "-h" ]; then 
	echo ""
	echo "------------------------------------------------------------------------------"
	create;	show;
fi
