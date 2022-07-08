
echo "####################################################################"
echo "## NLB Test Scripts for CB-Spider - 2022.06."
echo "##   NLB: CreateNLB "
echo "####################################################################"

echo ""


echo "#####---------- CreateNLB ----------####"
curl -sX POST http://localhost:1024/spider/nlb -H 'Content-Type: application/json' -d \
	'{
		"ConnectionName": "'${CONN_CONFIG}'", 
		"ReqInfo": {
			"Name": "spider-nlb-01", 
			"VPCName": "vpc-01", 
			"Type": "PUBLIC", 
			"Scope": "REGION", 
			"Listener": {
				"Protocol" : "TCP",
	       			"Port" : "80"
			},
		        "VMGroup": {
	       			"Protocol" : "TCP", 	       
	       			"Port" : "80", 	       
	       			"VMs" : ["vm-01", "vm-02"]
		 	}, 
			"HealthChecker": {
	       			"Protocol" : "TCP", 	       
	       			"Port" : "80", 	       
	       			"Interval" : "10", 	       
	       			"Timeout" : "10", 	       
	       			"Threshold" : "3"       
			}
		}
	}' | json_pp

