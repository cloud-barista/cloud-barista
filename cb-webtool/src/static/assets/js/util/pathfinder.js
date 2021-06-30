// path 와 매핑되는 controller의 이름 = key가 되어 
// 해당 key입력 시 main.go의 path를 return
// 필요한 param을 path에 적용하여 호출 url return
// leftmenu에서 script import



// map에 담긴 Key를 value로 바꿔 url을 return한다.
// url에는 main.go 에서 사용하는 path를 넣는다.
function setUrlByParam(url, urlParamMap){
    //resultVmCreateMap.set(resultVmKey, resultStatus)
    // var url = "/operation/manages/mcksmng/:clusteruID/:clusterName/del/:nodeID/:nodeName";    
    var returnUrl = url;
    for (let key of urlParamMap.keys()) { 
        console.log("urlParamMap " + key + " : " + urlParamMap.get(key) );
        
        var urlParamValue = urlParamMap.get(key)
        returnUrl = returnUrl.replace(key, urlParamValue);        
    }
    return returnUrl;
}

// conteroller의 methodName으로 main.go에 정의된 url값을 가져온다.
function getWebToolUrl(controllerKeyName){
    // ex ) monitoringGroup.GET("/operation/monitorings/mcismonitoring/mngform", controller.McisMonitoringMngForm)    
    let controllerMethodNameMap = new Map(
        [
            ["McisMonitoringMngForm", "/operation/monitorings/mcismonitoring/mngform"],
            ["VmMonitoringAgentRegForm", "/operation/monitorings/mcismonitoring/:mcisID/vm/:vmID/agent/mngform"],
            ["RemoteCommandVmOfMcis", "/operation/manages/mcismng/cmd/mcis/:mcisID/vm/:vmID"],
        ]
    );

    var webtoolUrl = controllerMethodNameMap.get(controllerKeyName);
    
    return webtoolUrl;
}

// main 화면인 경우에는 apitest로 보내고
// 그 외에는 helpArea를 보여준다.
// helpKey가 있는 경우에는 해당 key에 맞는 help 정보를 보여준다.
function showHelp(helpKey){
    var path = window.location.pathname;
    if( path == "/main"){
        location.href="/main/apitestmng"
    }else{
        //$("#helpArea").modal()        
        changePage("/operation/about/about");// About으로 이동
    }
}

//////////////// api -> local server -> target api  호출 ///////////////
// 한 화면에서 서로다른 형태로 호출이 가능하므로 caller(호출자) 를 callback에 같이 넘겨서 구분할 수 있게 함.
function getCommonNameSpaceList(caller){
    var url = "/setting/namespaces/namespace/list";
    axios.get(url,{
        headers:{
            'Content-Type' : "application/json"
        }
    }).then(result=>{
        console.log("get NameSpace Data : ",result.data);
        // var data = result.data.ns;
        var data = result.data;
        
        getNameSpaceListCallbackSuccess(caller, data);
            
    }).catch((error) => {
        console.warn(error);
        console.log(error.response)
        var errorMessage = error.response.data.error;
        var statusCode = error.response.status;
        // commonErrorAlert(statusCode, errorMessage) 
        
        getNameSpaceListCallbackFail(caller, error);
        
    });
}

function getCommonCloudConnectionList(caller, sortType){
    var url = "/setting/connections/cloudconnectionconfig/list";
    axios.get(url,{
        headers:{
                'Content-Type' : "application/json"
        }
    }).then(result=>{
        console.log("get CloudConnection Data : ",result.data);
        var data = result.data.ConnectionConfig;
        getCloudConnectionListCallbackSuccess(caller, data, sortType);
    }).catch((error) => {
        console.warn(error);
        console.log(error.response)
        var errorMessage = error.response.data.error;
        getCloudConnectionListCallbackFail(caller, error);
    });
}

function getCommonCredentialList(caller){
    var url = "/setting/connections/credential";
    axios.get(url,{
        headers:{
                'Content-Type' : "application/json"
        }
    }).then(result=>{
        console.log("get Credential Data : ",result.data);
        var data = result.data.Credential;
        getCredentialListCallbackSuccess(caller, data);
    }).catch((error) => {
        console.warn(error);
        console.log(error.response)
        var errorMessage = error.response.data.error;
        getCredentialListCallbackFail(caller, error);
    });
}


function getCommonRegionList(caller){
    var url = "/setting/connections/region"
    axios.get(url,{

    }).then(result=>{
        console.log("get Region Data : ",result.data);
        var data = result.data.Region;
        getRegionListCallbackSuccess(caller, data);
    }).catch((error) => {
        console.warn(error);
        console.log(error.response)
        // var errorMessage = error.response.data.error;
        // commonErrorAlert(statusCode, errorMessage) 
        getRegionListCallbackFail(caller, error);
    });
}


function getCommonDriverList(caller){
    var url = "/setting/connections"+"/driver";
    axios.get(url,{
        // headers:{
        //     'Authorization': "{{ .apiInfo}}",
        //     'Content-Type' : "application/json"
        // }
    }).then(result=>{
        console.log("get Driver Data : ",result.data);
        var data = result.data.Driver;
        getDriverListCallbackSuccess(caller, data);
    }).catch((error) => {
        console.warn(error);
        console.log(error.response)
        // var errorMessage = error.response.data.error;
        // var statusCode = error.response.status;
        // commonErrorAlert(statusCode, errorMessage) 
        getDriverListCallbackFail(caller, error);
    });
}

function getCommonNetworkList(caller){
    console.log("vnet : ");
    
    var url = "/setting/resources/network/list"
    var html = "";
    axios.get(url,{
        headers:{
            // 'Authorization': apiInfo
        }
    }).then(result=>{
        data = result.data.VNetList;
        console.log("vNetwork Info : ",result);
        console.log("vNetwork data : ",data);
        getNetworkListCallbackSuccess(caller, data);
    }).catch(error => {
        console.warn(error);
        console.log(error.response)
        // var errorMessage = error.response.data.error;
        // var statusCode = error.response.status;
        // commonErrorAlert(statusCode, errorMessage) 
        getNetworkListCallbackFail(caller, error);
    });
}


function getCommonSecurityGroupList(caller, sortType) {
    var url = "/setting/resources/securitygroup/list";
    axios.get(url, {
        headers: {
            // 'Authorization': "{{ .apiInfo}}",
            'Content-Type': "application/json"
        }
    }).then(result => {
        console.log("get SG Data : ", result.data);
        var data = result.data.SecurityGroupList; // exception case : if null 
        
        console.log("Data : ", data);
        if( caller == "securitygroupmng"){
			console.log("return get Data securitygroupmng")
			setSecurityGroupListAtServerImage(data, sortType)			
		}else if( caller == "mcissimpleconfigure"){
			console.log("return get Data")
			setSecurityGroupListAtSimpleConfigure(data)			
		}else if( caller == "mainsecuritygroup"){
			console.log("return get Data")
			getSecurityGroupListCallbackSuccess(caller, data)			
		}

	}).catch(error => {
		console.warn(error);
		console.log(error.response) 
        getSecurityGroupListCallbackFail(error)
	});
}

function getCommonSshKeyList(caller) {
    var url = "/setting/resources/sshkey/list"
    axios.get(url, {
        headers: {
            // 'Authorization': "{{ .apiInfo}}",
            'Content-Type': "application/json"
        }
    }).then(result => {
        console.log("get SSH Data : ", result.data);
        var data = result.data.SshKeyList; // exception case : if null 
        getSshKeyListCallbackSuccess(caller, data)
    }).catch(error => {
        console.warn(error);
        console.log(error.response)
        var errorMessage = error.response.data.error;
        var statusCode = error.response.status;
        // commonErrorAlert(statusCode, errorMessage);
        getSshKeyListCallbackFail(caller, error)
    });
}



// connection 정보가 바뀔 때 해당 connection에 등록 된 vmi(virtual machine image) 목록 조회.
// 공통으로 사용해야하므로 호출후 결과만 리턴... 그러나, ajax로 호출이라 결과 받기 전에 return되므로 해결방안 필요
function getCommonVirtualMachineImageList(caller, sortType) {
    var sortType = sortType;
    // var url = CommonURL + "/ns/" + NAMESPACE + "/resources/image";
    var url = "/setting/resources" + "/machineimage/list"
    axios.get(url, {
        headers: {
            // 'Authorization': "{{ .apiInfo}}",
            'Content-Type': "application/json"
        }
    }).then(result => {
        console.log("get Image List : ", result.data);
        
        var data = result.data.VirtualMachineImageList;
        // Data가져온 뒤 set할 method 호출
		if( caller == "virtualmachineimagemng"){
			console.log("return get Data")
			setVirtualMachineImageListAtServerImage(data, sortType)			
		}else if( caller == "mcissimpleconfigure"){
			console.log("return get Data")
			setVirtualMachineImageListAtSimpleConfigure(data, sortType)			
		}else if( caller == "mainimage"){
			console.log("return get Data")
			getImageListCallbackSuccess(caller, data)		
		}
    // }).catch(function(error){
    //     console.log("list error : ",error);        
    // });
	}).catch(error => {
		console.warn(error);
		console.log(error.response) 
	});
}


function getCommonVirtualMachineSpecList(caller, sortType) {
    console.log("CommonSpecCaller : " + caller);
    console.log("CommonSpecSortType : " + sortType);
    // var url = CommonURL + "/ns/" + NAMESPACE + "/resources/image";
    var url = "/setting/resources" + "/vmspec/list"

    axios.get(url, {
        headers: {
            // 'Authorization': "{{ .apiInfo}}",
            'Content-Type': "application/json"
        }
    }).then(result => {
        console.log("get Spec List : ", result.data);
        
        var data = result.data.VmSpecList;

        if ( caller == "virtualmachinespecmng") {
            console.log("return get Data");
            virtualMachineSpecListCallbackSuccess(caller, data, sortType);	
            // setVirtualMachineSpecListAtServerSpec(data, sortType);
        }else if( caller == "mainspec"){
			console.log("return get Data")
			getSpecListCallbackSuccess(caller, data)		
		}
    }).catch(error => {
		console.warn(error);
		console.log(error.response) 
        getSpecListCallbackFail(error)
	});
}

// /lookupSpecs
function getCommonLookupSpecList(caller, connectionName) {    
    var url = "/setting/resources/vmspec/lookupvmspec"
    axios.get(url, {
        headers: {
            // 'Authorization': "{{ .apiInfo}}",
            'Content-Type': "application/json"
        },
        params: {
            connectionName: connectionName
        }
    }).then(result => {
        console.log("get Image List : ", result.data);
        
        var data = result.data.CspVmSpecList;
        
		// Data가져온 뒤 set할 method 호출
		if( caller == "vmspecmng"){
			console.log("return get Data")			
			lookupSpecListCallbackSuccess(caller, data)		
		}
    // }).catch(function(error){
    //     console.log("list error : ",error);        
    // });
	}).catch(error => {
		console.warn(error);
		console.log(error.response) 
        lookupSpecListCallbackSuccess(error)
	});
}

// 현재 선택 된 
function putFetchSpecs(connectionName){
    var url = "/setting/resources/vmspec/fetchvmspec"
    axios.post(url, {
        headers: {
            'Content-Type': "application/json"
        },
        params: {
            connectionName: connectionName
        }      
    }).then(result => {
        console.log(result);
        if(result.data.status == 200 || result.data.status == 201){
            commonAlert("Spec Fetched");                
            // Region 갱신 
            getRegionList();   
        }else{
            commonAlert("Fail to Spec Fetched");
        }
	}).catch(error => {
		console.warn(error);
		console.log(error.response) 
        var errorMessage = error.response.data.error;
        var statusCode = error.response.status;
        commonErrorAlert(statusCode, errorMessage);
	});
}

function getCommonFilterSpecsByRange(caller, searchObj){
    var url = "/setting/resources/vmspec/filterspecsbyrange";

    // 똑같은데... 얘는 param을 못받음
    // axios.post(url, {    
    //     headers: { 
    //                 'Content-type': 'application/json',
    //             },
    //     searchObj       
    axios.post(url,searchObj,{
        headers: { 
            'Content-type': 'application/json',
            // 'Authorization': apiInfo, 
        }

    }).then(result => {
        console.log(result);
        // if(result.data.status == 200 || result.data.status == 201){
        //     var data = result.data.VmSpec
        //     // commonAlert("Spec Searched");                            
        // }else{
        //     // commonAlert("Fail to Spec Searched");
        // }
        var data = result.data.VmSpecList;
        console.log("caller " + caller)
        if ( caller == "virtualmachinespecmng") {
            console.log("return get Data");
            virtualMachineSpecListCallbackSuccess(caller, data, sortType);	
            // setVirtualMachineSpecListAtServerSpec(data, sortType);
        }else if ( caller == "vmassistpopup"){
            filterSpecsByRangeCallbackSuccess(caller, data);
        }
	}).catch(error => {
		console.warn(error);
		console.log(error.response) 
        var errorMessage = error.response.data.error;
        var statusCode = error.response.status;
        commonErrorAlert(statusCode, errorMessage);
	});

}

// /lookupImages
function getCommonLookupImageList(caller, connectionName) {    
    //var url = "/setting/resources/vmimage/lookupvmimagelist"
    var url = "/setting/resources/machineimage/lookupimages"
    axios.get(url, {
        headers: {
            // 'Authorization': "{{ .apiInfo}}",
            'Content-Type': "application/json"
        },
        params: {
            connectionName: connectionName
        }
    }).then(result => {
        console.log("get Image List : ", result.data);
        
        var data = result.data.VirtualMachineImageList;
        
		// Data가져온 뒤 set할 method 호출
		if( caller == "vmimagemng"){
			console.log("return get Data")			
			lookupVmImageListCallbackSuccess(caller, data)		
		}
    // }).catch(function(error){
    //     console.log("list error : ",error);        
    // });
	}).catch(error => {
		console.warn(error);
		console.log(error.response) 
        lookupVmImageListCallbackFail(error)
	});
}

//
///ns/{nsId}/resources/fetchImages
function getCommonFetchImages(caller, connectionName) {
    var url = "/setting/resources/machineimage/fetchimages"
    axios.post(url, {
        headers: {
            'Content-Type': "application/json"
        }        
    }).then(result => {
        console.log(result);
        if(result.data.status == 200 || result.data.status == 201){
            commonAlert("Image Fetched");                            
        }else{
            commonAlert("Fail to Image Fetched");
        }
	}).catch(error => {
		console.warn(error);
		console.log(error.response) 
        var errorMessage = error.response.data.error;
        var statusCode = error.response.status;
        commonErrorAlert(statusCode, errorMessage);
	});
}


// MCIS 목록 존재여부
function getCommonMcisList(caller) {
    var url = "/operation/manages/mcismng/list"

    axios.get(url, {
        headers: {
            'Content-Type': "application/json"
        }
    }).then(result => {
        console.log("get Mcis List : ", result.data);
        
        var data = result.data.McisList;

        // if ( caller == "mainmcis") {
            console.log("return get Data");            
			getMcisListCallbackSuccess(caller, data)		
		// }
    }).catch(error => {
		console.warn(error);
		console.log(error.response) 
        getMcisListCallbackFail(error)
	});
}

function getCommonMcisList(caller) {
    var url = "/operation/manages/mcismng/list"

    axios.get(url, {
        headers: {
            'Content-Type': "application/json"
        }
    }).then(result => {
        console.log("get Mcis List : ", result.data);
        
        var data = result.data.McisList;

        // if ( caller == "mainmcis") {
            console.log("return get Data");            
			getMcisListCallbackSuccess(caller, data)		
		// }
    }).catch(error => {
		console.warn(error);
		console.log(error.response) 
        getMcisListCallbackFail(error)
	});
}

function getCommonMcksList(caller) {
    var url = "/operation/manages/mcksmng/list"

    axios.get(url, {
        headers: {
            'Content-Type': "application/json"
        }
    }).then(result => {
        console.log("get Mcks List : ", result.data);
        
        var data = result.data.McksList;

        // if ( caller == "mainmcis") {
            console.log("return get Data");            
			getMcksListCallbackSuccess(caller, data)		
		// }
    }).catch(error => {
		console.warn(error);
		console.log(error.response) 
        getMcksListCallbackFail(error)
	});
}


function getCommonVmSecurityGroupInfo(caller, securityGroupId){
    var url = "/setting/resources/securitygroup/" + securityGroupId
    
    axios.get(url,{
        headers: {
            'Content-Type': "application/json"
        }
    }).then(result => {
        console.log("get SecurityGroup List : ", result.data);
        
        var data = result.data.SecurityGroupInfo;

        // if ( caller == "mainmcis") {
            console.log("return get Data");            
			getSecurityGroupCallbackSuccess(caller, data)		
		// }
    }).catch(error => {
		console.warn(error);
		console.log(error.response) 
        getSecurityGroupCallbackFail(error)
	});
}

function getCommonVmImageInfo(caller, imageId){
    
    //var url = CommonURL+"/ns/"+NAMESPACE+"/resources/image/"+imageId
    // var apiInfo = ApiInfo
    var url = "/setting/resources/machineimage/" + imageId
    axios.get(url,{
        // headers:{
        //     'Authorization': apiInfo
        // }
    }).then(result=>{
        console.log(result);
        getCommonVmImageInfoCallbackSuccess(caller, result.data.VirtualMachineImageInfo);        
    })

}


// MCIS에 명령어 날리기
function postRemoteCommandMcis(mcisID, commandWord){
    var orgUrl = "/operation/manages/mcismng/cmd/mcis/:mcisID";
    var urlParamMap = new Map();
    urlParamMap.set(":mcisID", mcisID)
    var url = setUrlByParam(orgUrl, urlParamMap)

    console.log(" command = " + commandWord)    
    axios.post(url, {
        // headers: {
        //     'Content-Type': "application/json"
        // },
        command: commandWord        
    }).then(result => {
        console.log(result);
        if(result.data.status == 200 || result.data.status == 201){
            commonAlert("Success to Send the Command " + result.data.message);
        }else{
            commonAlert("Fail to Send the Command " + result.data.message);
        }
	}).catch(error => {
		console.warn(error);
		console.log(error.response) 
        var errorMessage = error.response.data.error;
        var statusCode = error.response.status;
        commonErrorAlert(statusCode, errorMessage);
	});
}

// VM에 명령어 날리기
function postRemoteCommandVmOfMcis(mcisID, vmID, commandWord){    
    //RemoteCommandVmOfMcis
    var orgUrl = "/operation/manages/mcismng/cmd/mcis/:mcisID/vm/:vmID";
    var urlParamMap = new Map();
    urlParamMap.set(":mcisID", mcisID)
    urlParamMap.set(":vmID", vmID)
    var url = setUrlByParam(orgUrl, urlParamMap)

    console.log(" command = " + commandWord)    
    axios.post(url, {
        // headers: {
        //     'Content-Type': "application/json"
        // },
        command: commandWord        
    }).then(result => {
        console.log(result);
        if(result.data.status == 200 || result.data.status == 201){
            commonAlert("Success to Send the Command " + result.data.message);
        }else{
            commonAlert("Fail to Send the Command " + result.data.message);
        }
	}).catch(error => {
		console.warn(error);
		console.log(error.response) 
        var errorMessage = error.response.data.error;
        var statusCode = error.response.status;
        commonErrorAlert(statusCode, errorMessage);
	});
}

// dragonfly monitoring agent 설치 및 동작여부
function checkDragonFlyMonitoringAgent(mcisID, vmID){
  return true;
}
// form 화면에서 조회에 문제가 있는 경우 표시
// 모든 form 화면 시작할 때(onLoad 시) 체크하도록
// Header.html 에 정의
function checkLoadStatus(){
    var returnMessage = $("#returnMessage").val();
    var returnStatusCode = $("#returnStatusCode").val();
    if( returnStatusCode != 200 && returnStatusCode != 201){
        commonErrorAlert(returnStatusCode, returnMessage);
    }
}
