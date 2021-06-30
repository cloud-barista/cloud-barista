function applyMonitoringConfig(){
    var agentItv = $("#agentInterval").val();
    var colItv = $("#collectorInterval").val();
    var maxhostcnt = $("#maxHostCount").val();
    console.log(agentItv + ", " + colItv + ", " + maxhostcnt);
    var message = "Set monitoring config.<br><br>Agent Interval : " + agentItv + "<br>Collector Interval : " + colItv + "<br>Max Host Count : " + maxhostcnt + "<br><br>Are you sure?"; 
    console.log(message);
    commonConfirmMsgOpen("monitoringConfigPolicyConfig", message);
}

function regMonitoringConfigPolicy() {
    // 여기에서 value 값을 담아서 PUT던진다?
    var url = "/operation/policies/monitoringconfig/policy/put"        
    console.log("Monitoring Policy Reg URL : ",url)
    var agentItv = $("#agentInterval").val();
    var colItv = $("#collectorInterval").val();
    var maxhostcnt = $("#maxHostCount").val();
    var obj = {        
        agent_interval : Number(agentItv),
        collector_interval : Number(colItv),
        max_host_count : Number(maxhostcnt)
    }
    console.log("info Monitoring Policy obj Data : ", obj);
    
    if (agentItv) {
        axios.put(url, obj, {
            headers: {
                'Content-type': 'application/json',
                // 'Authorization': apiInfo,
            }
        }).then(result => {
            console.log("result Monitoring Policy : ", result);
            var data = result.data;
            console.log("모니터링 : ", data);
            
            if (data.status == 200 || data.status == 201) {
                commonAlert("Success Setting Monitoring Policy!!")
                var resultData = data.MonitoringConfig;
                $("#agentInterval").val(resultData.agent_interval);
                $("#collectorInterval").val(resultData.collector_interval);
                $("#maxHostCount").val(resultData.max_host_count);
                
                //displayVNetInfo("REG_SUCCESS")
                
            } else {
                commonAlert("Fail Set Monitoring Policy" + data.message)
            }

        }).catch((error) => {
            console.log(error.response) 
            var errorMessage = error.response.data.error;
            var statusCode = error.response.status;
            commonErrorAlert(statusCode, errorMessage) 
        });
    } else {
        commonAlert("Input Monitoring Policy Data")
        $("#agentInterval").focus()
        return;
    }
}