
echo "####################################################################"
echo "## VM Test Scripts for CB-Spider IID Working Version - 2020.04.21."
echo "##   VM: StartVM "
echo "####################################################################"

curl -sX POST http://localhost:1024/spider/vm -H 'Content-Type: application/json' -d '{ "ConnectionName": "'${CONN_CONFIG}'", "ReqInfo": { "Name": "'${CONN_CONFIG}'-vm-01", "ImageName": "'${IMAGE_NAME}'", "VPCName": "NA", "SubnetName": "NA", "SecurityGroupNames": [ "NA" ], "VMSpecName": "'${SPEC_NAME}'", "KeyPairName": "NA", "VMUserId": "NA", "VMUserPasswd": "NA" } }' |json_pp
