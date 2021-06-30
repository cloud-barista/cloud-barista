// connection 정보가 바뀌었을 때, 변경 될 object : 원래는 각각 만들어야 하나, 가져오는게 spec만 있어서 plane, worker 같이 씀.
function changeConnectionInfo(configName, targetObjId){
    console.log("config name : ",configName)
    if( configName == ""){
        // 0번째면 selectbox들을 초기화한다.(vmInfo, sshKey, image 등)
    }
    
    getSpecInfo(configName, targetObjId);
}

// connection에 맞는 spec들 조회
function getSpecInfo(configName, targetObjId){
    var configName = configName;
    if(!configName){
        configName = $("#nodeConnectionName option:selected").val();
    }

    var url = "/setting/resources/vmspec/list"
    var html = "";
    axios.get(url,{
        // headers:{
        // 	'Authorization': apiInfo
        // }
    }).then(result=>{
        // console.log(result.data)
        var data = result.data.VmSpecList
        console.log("spec result : ",data)
        if(data){
            html +="<option value=''>Select SpecName</option>"
            data.filter(csp => csp.connectionName === configName).map(item =>(
                html += '<option value="'+item.cspSpecName+'">'+item.name+'('+item.cspSpecName+')</option>'	
            ))

        }else{
            html +=""
        }       
      
        $("#" + targetObjId).empty();
        $("#" + targetObjId).append(html);        
    })
}

// mcks , node deploy
// 우선 mcks 부터
function nodeDone_btn(){
    var mcksID = $("#mcksID").val();
    var mcksName = $("#mcksName").val();
    
    var workerCountLength = $("input[name='workerCount']").length;
    console.log("workerCountLength1 " + workerCountLength)
    var workerConnectionData = new Array();
    var workerCountData = new Array();
    var workerSpecIdData = new Array();
    for(var i=0; i<workerCountLength; i++){   
        var workerId = $("input[name='workerCount']").eq(i).attr("id");
        console.log("workerId " + workerId)
        if( workerId.indexOf("hidden_worker") > -1) continue;// 복사를 위한 영역이 있으므로
        console.log("aa " + workerId)
        workerConnectionData.push($("select[name='workerConnectionName']")[i].value);
        workerCountData.push($("input[name='workerCount']")[i].value);
        workerSpecIdData.push($("select[name='workerSpecId']")[i].value);

        if( !workerConnectionData[i]){
            commonAlert("Please Select Connection " + i)
            return;
        }

        if( !workerCountData[i]){
            commonAlert("Please Input Worker Count" + i)
            return;
        }

        if( !workerSpecIdData[i]){
            commonAlert("Please Select Worker Spec" + i)
            return;
        }
    }
    console.log(workerConnectionData)
    console.log(workerCountData)
    console.log(workerSpecIdData)

    var new_obj = {}
    
    // VM추가시에는 controlPlane 없음.
    // var controlPlanes = new Array(controlPlaneLength);
    // console.log("controlPlaneConnectionLength " + controlPlaneLength)
    // for(var i=0; i<controlPlaneLength; i++){
    //     console.log("controlPlane " + i)
    //     var new_controlPlane = {}
    //     new_controlPlane['connection'] = controlPlaneConnectionData[i];
    //     new_controlPlane['count'] = Number(controlPlaneCountData[i])
    //     new_controlPlane['spec'] = controlPlaneSpecIdData[i]
    //     controlPlanes[i] = new_controlPlane
    // }
    // new_obj['controlPlane'] = controlPlanes;

    var workers = new Array(workerCountData.length);
    for(var i=0; i<workerCountData.length; i++){
        console.log("workerCountLength " + i)
        var new_worker = {}
        new_worker['connection'] = workerConnectionData[i];
        new_worker['count'] = Number(workerCountData[i])
        new_worker['spec'] = workerSpecIdData[i]
        
        // new_worker['config'] = workerConnectionData[i];
        // new_worker['workerNodeCount'] = Number(workerCountData[i])
        // new_worker['workerNodeSpec'] = workerSpecIdData[i]

        workers[i] = new_worker
        // new_obj = new_worker;

        // new_obj['config'] = workerConnectionData[i];
        // new_obj['workerNodeCount'] = Number(workerCountData[i])
        // new_obj['workerNodeSpec'] = workerSpecIdData[i]
    }
    new_obj['worker'] = workers;
    // new_obj = workers;
   
    console.log(new_obj);

    try{
        // configurer 는 mcks 선택하고 들어옴. : TODO : MCKS create 와 node create는 버튼 액션을 달리해야
        // /operation/manages/mcksmng/:clusteruID/:clusterName/reg/proc
        var url = "/operation/manages/mcksmng/" + mcksID + "/" + mcksName + "/reg/proc";
        axios.post(url,new_obj,{
            headers :{
                },
        }).then(result=>{
            console.log("data : ",result);
            console.log("Result Status : ",result.status); 

            var statusCode = result.data.status;
            var message = result.data.message;
            console.log("Result Status : ",statusCode); 
            console.log("Result message : ",message); 

            if(result.status == 201 || result.status == 200){
                commonAlert("Node Add Success")
                var targetUrl = "/operation/manages/mcksmng/mngform"
                changePage(targetUrl);
            
            }else{
                commonErrorAlert(statusCode, message) 
            }
        }).catch((error) => {
            console.log(error);
            console.log(error.response)
            var errorMessage = error.response.data.error;
            commonErrorAlert(statusCode, errorMessage) 
        })
    }finally{
        
    }
}

// WorkNode 추가
var lastWorkerId = "";
function addWorkNode(){
    console.log("addWorkNode start");
    try{
    // 마지막 name의 index 추출
    var maxWorkerId = "";
    var nameCount = 0;
    $("input[name='workerCount']").each(function (i) {
        var currentWorkCountID = $(this).attr('id');
        console.log("::: " + currentWorkCountID)
        if( currentWorkCountID.indexOf("hidden") == -1){
            //console.log( i + "번째  : " + $("input[name='workerCount']").eq(i).attr("value") );
            console.log("currentWorkCountID=" + currentWorkCountID)
            maxWorkerId = currentWorkCountID;
            nameCount++;
        }
   });

    if( lastWorkerId == "" ){
        lastWorkerId = maxWorkerId;
    }

    var lastIndexArr = lastWorkerId.split ("_")
    var lastIndex = lastIndexArr[lastIndexArr.length-1];

    var maxIndexArr = maxWorkerId.split ("_")
    var maxIndex = maxIndexArr[maxIndexArr.length-1];
    // console.log( lastWorkerId + " <> " + maxWorkerId)
    // console.log(maxIndexArr)
    // console.log( lastIndex + " : " + maxIndex + " : " + nameCount)
    if( lastIndex <= maxIndex){
        nameCount = Number(maxIndex) + 1;
        lastWorkerId = maxWorkerId;
    }
//    var lastIndexArr = lastWorkerId.split ("_")
//    var lastIndex = lastIndexArr[lastIndexArr.length-1];
//    console.log("lastIndex=" + lastIndex)
    // var addWorkerIndex = Number(lastIndex) +1;
    var addWorkerIndex = Number(nameCount);
    
    // console.log("addWorkerIndex=" + addWorkerIndex)
    var addWorkerHtml = $('#hidden_work_area').clone();
    // console.log(addWorkerHtml.html());    
    var addW = "";
    // 최초 1번만 .html() 이 먹고 다음부터는 string으로 인식함.
    addW = addWorkerHtml.html().replace('hidden_mcks_Worker_list', 'mcks_Worker_list_' + addWorkerIndex);
    addW = addW.replace(/hidden_workerProvider/gi, 'workerProvider_' + addWorkerIndex);
    addW = addW.replace(/hidden_workerConnectionName/gi, 'workerConnectionName_' + addWorkerIndex);
    addW = addW.replace(/hidden_workerCount/gi, 'workerCount_' + addWorkerIndex);
    addW = addW.replace(/hidden_workerSpecId/gi, 'workerSpecId_' + addWorkerIndex);
    addW = addW.replace(/hidden_workerRemove/gi, 'workerRemove_' + addWorkerIndex);
    addW = addW.replace(/hidden_workerAddCount/gi, 'workerAddCount_' + addWorkerIndex);
    

    //$("#mcks_Worker_list").append(addW); 
    $("#mcksNodeArea").append(addW);    
    $("#mcks_Worker_list_" + addWorkerIndex).css("display", "block");
    //$("#aa").css("display", "block");
    // console.log($("#mcks_Worker_list").html())

    $("#workerAddCount_" + addWorkerIndex).text(addWorkerIndex);
    }catch(e){
        console.log(e);
    }
}

function removeWorkerNode(removeWorkerId){
    console.log("removeWorkerId " + removeWorkerId)
    
    var workerArr = removeWorkerId.split("_");
    console.log(workerArr)
    var workerIndex = workerArr[workerArr.length-1];
    console.log("removeWorkerNode " + workerIndex)
    $("#mcks_Worker_list_" + workerIndex).remove();
}