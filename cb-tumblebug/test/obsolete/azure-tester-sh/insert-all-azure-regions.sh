#!/bin/bash
RESTSERVER=localhost

LOCS=(`cat azure-locations-list.txt |grep "name" |awk '{print $2}' |sed 's/",//g' |sed 's/"//g'`)

for LOC in "${LOCS[@]}"
do
	echo $LOC

	curl -H "${AUTH}" -sX POST http://$RESTSERVER:1024/region -H 'Content-Type: application/json' -d '{"RegionName":"azure-'$LOC'","ProviderName":"AZURE", "KeyValueInfoList": [{"Key":"location", "Value":"'$LOC'"}, {"Key":"ResourceGroup", "Value":"jhseo-test-'$LOC'"}]}'
	curl -H "${AUTH}" -sX POST http://$RESTSERVER:1024/connectionconfig -H 'Content-Type: application/json' -d '{"ConfigName":"azure-'$LOC'-config","ProviderName":"AZURE", "DriverName":"azure-driver01", "CredentialName":"azure-credential01", "RegionName":"azure-'$LOC'"}'

done
