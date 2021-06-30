jQuery.fn.center = function () {
    console.log("height");
    this.css('top', Math.max(0,(($(window).height()-$(this).outerHeight())/2) + $(window).scrollTop())+'px');
    
    console.log($(window).height() + " - " + $(this).outerHeight() + " : " + (($(window).height()-$(this).outerHeight())/2));
    return this;
}
    
// div id = Ajax_Loading 이 있어야 함.
// 요청 인터셉터
axios.interceptors.request.use(function (config) {
        console.log("axios.interceptors.request")      
        try{
        // $('#loadingContainer').css('position', 'fixed');
        $("#loadingContainer").center();
        $('#loadingContainer').show();
        // $('#loadingContainer').modal();
        }catch(e){
            console.log(e);
        }
        return config;
    }, function (error) {
        console.log("axios.interceptors.request error")
        console.log(error)
        // 에라 나면 로딩 끄기
        $('#loadingContainer').hide();
        // AjaxLoadingShow(false);
        return Promise.reject(error);
    });

// 응답 인터셉터
axios.interceptors.response.use(function (response) {
        console.log("axios.interceptors.response")
        // 응답 받으면 로딩 끄기
        $('#loadingContainer').hide();
        return response;
    }, function (error) {
        console.log("axios.interceptors.response error")
        console.log(error)
        // 응답 에러 시에도 로딩 끄기
        $('#loadingContainer').hide();
        return Promise.reject(error);
    });

function AjaxLoadingShow(isShow){
    try{
        if(isShow) {
            $('#Ajax_Loading').show();
        }else{
            $('#Ajax_Loading').hide();
        }
    }catch(e){
        alert(e);
    }
}
//========== 로딩 바 시작 =========    
// $(document).ready(function(){
//     $('#Ajax_Loading').hide(); //첫 시작시 로딩바를 숨겨준다.
//  })
//  .ajaxStart(function(){
//      $('#Ajax_Loading').show(); //모든 ajax 통신 시작시 로딩바를 보여준다.
//      //$('html').css("cursor", "wait"); //마우스 커서를 로딩 중 커서로 변경
//  })
//  .ajaxStop(function(){
//      $('#Ajax_Loading').hide(); //모든 ajax 통신 종료시 로딩바를 숨겨준다.
//      //$('html').css("cursor", "auto"); //마우스 커서를 원래대로 돌린다
//  });
//========== 로딩 바 종료 =========

// 다른 화면으로 이동 시킬 때 Loading bar 표시를 위해
function changePage(url){
    $('#loadingContainer').show();// page 이동 전 loading bar를 보여준다.
    location.href = url;
}

// 그런데 inputtype=text 를 password로 바꾸기만해도 해당기능이 동작 함.
function showPassword(passwordObjId){
    var passObj = document.getElementById(passwordObjId);
    // var passObj = $("#" + passwordObjId);
    console.log( passObj)
    console.log( " pw obj tyle " + passObj.type)
    if (passObj.type ==="password") {
        passObj.type = "text";
    }else{
        passObj.type = "password"; 
    }
}
// 문자열이 빈 경우 defaultString을 return
function nvl(str, defaultStr){         
    if(typeof str == "undefined" || str == null || str == "")
        str = defaultStr ;
     
    return str ;
}
function nvlDash(str){         
    if(typeof str == "undefined" || str == null || str == "" || str == "undefined")
        str = '-';
     
    return str ;
}

function guideAreaHide(){
	console.log("hide brfore")
	$("#guideArea").modal("hide");
	console.log("hide after")
}

// message를 표현할 alert 창
function commonAlert(alertMessage){
    console.log(alertMessage);
    // $('#alertText').text(alertMessage);
    $('#alertText').html(alertMessage);
    $("#alertArea").modal();
}
// alert창 닫기
function commonAlertClose(){
    $("#alertArea").modal("hide");
}

// 에러 메세지 alert 통일 용
function commonErrorAlert(statusCode, message){
    commonAlert("Error(" + statusCode + ") : " + message);
}

// confirm modal창 보이기 modal창이 열릴 때 해당 창의 text 지정, close될 때 action 지정
function commonConfirmOpen(targetAction){
    console.log("commonConfirmOpen : " + targetAction)

    //  [ id , 문구]
    let confirmModalTextMap = new Map(
        [
            ["Logout", "Would you like to logout?"],
            ["Config", "Would you like to set Cloud config ?"],
            ["SDK", "Would you like to set Cloud Driver SDK ?"],
            ["Credential", "Would you like to set Credential ?"],
            ["Region", "Would you like to set Region ?"],
            ["Provider", "Would you like to set Cloud Provider ?"],

            ["MoveToConnection", "Would you like to set Cloud config ?"],
            ["DeleteCloudConnection", "Would you like to delete <br /> the Cloud connection? "],

            ["DeleteCredential", "Would you like to delete <br /> the Credential? "],
            ["DeleteDriver", "Would you like to delete <br /> the Driver? "],
            ["DeleteRegion", "Would you like to delete <br /> the Region? "],


            // ["IdPassRequired", "ID/Password required !"],    --. 이거는 confirm이 아니잖아
            ["idpwLost", "Illegal account / password 다시 입력 하시겠습니까?"],
            ["ManageNS", "Would you like to manage <br />Name Space?"],
            ["NewNS", "Would you like to add a new Name Space?"],
            ["AddNewNameSpace", "Would you like to register NameSpace <br />Resource ?"],
            ["NameSpace", "Would you like to move <br />selected NameSpace?"],
            ["ChangeNameSpace", "Would you like to move <br />selected NameSpace?"],
            ["DeleteNameSpace", "Would you like to delete <br />selected NameSpace?"],

            ["AddNewVpc", "Would you like to create a new Network <br />Resource ?"],
            ["DeleteVpc", "Are you sure to delete this Network <br />Resource ?"],

            ["AddNewSecurityGroup", "Would you like to create a new Security <br />Resource ?"],
            ["DeleteSecurityGroup", "Would you like to delete Security <br />Resource ?"],
            
            ["AddNewSshKey", "Would you like to create a new SSH key <br />Resource ?"],
            ["DeleteSshKey", "Would you like to delete SSH key <br />Resource ?"],     
            
            ["AddNewVirtualMachineImage", "Would you like to register Image <br />Resource ?"],
            ["DeleteVirtualMachineImage", "Would you like to un-register Image <br />Resource ?"],  
            ["FetchImages", "Would you like to fetch images <br /> to this NameSpace ?"],  
            
            ["AddNewVmSpec", "Would you like to register Spec <br />Resource ?"],
            ["DeleteVmSpec", "Would you like to un-register Spec <br />Resource ?"],  
            ["FetchSpecs", "Would you like to fetch Spec <br /> to this NameSpace ?"],  

            ["GotoMonitoringPerformance", "Would you like to view performance <br />for MCIS ?"],
            ["GotoMonitoringFault", "Would you like to view fault <br />for MCIS ?"],
            ["GotoMonitoringCost", "Would you like to view cost <br />for MCIS ?"],
            ["GotoMonitoringUtilize", "Would you like to view utilize <br />for MCIS ?"],

            ["McisLifeCycleReboot", "Would you like to reboot MCIS ?"],// mcis_life_cycle('reboot')
            ["McisLifeCycleSuspend", "Would you like to suspend MCIS ?"],//onclick="mcis_life_cycle('suspend')
            ["McisLifeCycleResume", "Would you like to resume MCIS ?"],//onclick="mcis_life_cycle('resume')"
            ["McisLifeCycleTerminate", "Would you like to terminate MCIS ?"],//onclick="mcis_life_cycle('terminate')
            ["McisManagement", "Would you like to manage MCIS ?"],// 해당 function 없음...
            ["MoveToMcisManagement", "Would you like to manage MCIS ?"],
            ["MoveToMcisManagementFromDashboard", "Would you like to manage MCIS ?"],
            
            ["AddNewMcis", "Would you like to create MCIS ?"],
            ["DeleteMcis", "Are you sure to delete this MCIS? "],
            ["ImportScriptOfMcis", "Would you like to import MCIS script? "],            
            ["ExportScriptOfMcis", "Would you like to export MCIS script? "],
            
            ["AddNewVmOfMcis", "Would you like to add a new VM to this MCIS ?"],

            ["VmLifeCycle", "Would you like to view Server ?"],
            ["VmLifeCycleReboot", "Would you like to reboot MCIS ?"], //onclick="vm_life_cycle('reboot')"
            ["VmLifeCycleSuspend", "Would you like to suspend MCIS ?"], // onclick="vm_life_cycle('suspend')"
            ["VmLifeCycleResume", "Would you like to resume MCIS ?"], // onclick="vm_life_cycle('resume')"
            ["VmLifeCycleTerminate", "Would you like to terminate MCIS ?"], // onclick="vm_life_cycle('terminate')"
            ["VmManagement", "Would you like to manage VM ?"], // 해당 function 없음
            ["AddNewVm", "Would you like to add VM ?"], //onclick="vm_add()"
            ["ExportVmScriptOfMcis", "Would you like to export VM script ?"], //onclick="vm_add()"
            

            ["DifferentConnection", "Do you want to set different connectionName?"],
            ["DifferentConnectionAtSecurityGroup", "Do you want to set different connectionName?"],

            ["AddMonitoringAlertPolicy", "Would you like to register Threshold ?"],
            ["DeleteMonitoringAlertPolicy", "Are you sure to delete this Threshold ?"],
            ["AddNewMcks", "Would you like to create MCKS ?"],
            ["DeleteMcks", "Are you sure to delete this MCKS? "],
            ["AddNewNodeOfMcks", "Would you like to add a new Node to this MCKS ?"],
            ["DeleteNodeOfMcks", "Would you like to delete a Node of this MCKS ?"],
            

            ["AddMonitoringAlertEventHandler", "Would you like to add<br />Monitoring Alert Event-Handler ?"],
            ["deleteMonitoringAlertEventHandler", "Are you sure to delete<br />this Monitoring Alert Event-Handler?"],
        ]
    );
    console.log(confirmModalTextMap.get(targetAction));
    try{
        // $('#modalText').text(targetText);// text아니면 html로 해볼까? 태그있는 문구가 있어서
        //$('#modalText').text(confirmModalTextMap.get(targetAction));
        $('#confirmText').html(confirmModalTextMap.get(targetAction));
        $('#confirmOkAction').val(targetAction);
        
        if( targetAction == "Region"){
            // button에 target 지정
            // data-target="#Add_Region_Register"
            // TODO : confirm 으로 물어본 뒤 OK버튼 클릭 시 targetDIV 지정하도록
        }
        $('#confirmArea').modal(); 
    }catch(e){
        console.log(e);
        alert(e);
    }
}

// confirm modal창 보이기 modal창이 열릴 때 해당 창의 text 지정, close될 때 action 지정, text 내용 전송
function commonConfirmMsgOpen(targetAction, message){
    console.log("commonConfirmMsgOpen : " + targetAction)
    
    try{
        $('#confirmText').html(message);
        $('#confirmOkAction').val(targetAction);
                
        $('#confirmArea').modal(); 
    }catch(e){
        console.log(e);
        alert(e);
    }
}

// confirm modal창에서 ok버튼 클릭시 수행할 method 지정
function commonConfirmOk(){
    //modalArea
    var targetAction = $('#confirmOkAction').val();
    if( targetAction == "Logout"){
        // Logout처리하고 index화면으로 간다. Logout ==> cookie expire
        // location.href="/logout"
        var targetUrl = "/logout"
        changePage(targetUrl)
        
    }else if ( targetAction == "MoveToConnection"){
        var targetUrl="/setting/connections/cloudconnectionconfig/mngform"
        changePage(targetUrl)
    }else if ( targetAction == "DeleteCloudConnection"){
        deleteCloudConnection();    
    }else if ( targetAction == "Config"){
        //id="Config"
        console.log("commonConfirmOk " + targetAction);
    }else if ( targetAction == "SDK"){
        //id="SDK"
        console.log("commonConfirmOk " + targetAction);
    }else if ( targetAction == "DeleteCredential"){
        deleteCredential();
    }else if ( targetAction == "DeleteDriver"){
        deleteDriver();
    }else if ( targetAction == "DeleteRegion"){
        deleteRegion();

    }else if ( targetAction == "Credential"){
        //id="Credential"
        console.log("commonConfirmOk " + targetAction);
    }else if ( targetAction == "Region"){
        //id="Region"
        console.log("commonConfirmOk " + targetAction);
    }else if ( targetAction == "Provider"){
        //id="Provider"
        console.log("commonConfirmOk " + targetAction);
    }else if ( targetAction == "required"){//-- IdPassRequired
        console.log("commonConfirmOk " + targetAction);
    }else if ( targetAction == "idpwLost"){//-- 
        console.log("commonConfirmOk " + targetAction);
    }else if ( targetAction == "ManageNS"){//-- ManageNS
        var targetUrl = "/setting/namespaces/namespace/mngform"
        changePage(targetUrl)
    }else if ( targetAction == "NewNS"){//-- NewNS
        var targetUrl = "/setting/namespaces/namespace/mngform"
        changePage(targetUrl)
    }else if ( targetAction == "ChangeNameSpace"){//-- ChangeNameSpace
        var changeNameSpaceID = $("#tempSelectedNameSpaceID").val();
        setDefaultNameSpace(changeNameSpaceID)
    }else if ( targetAction == "AddNewNameSpace"){//-- AddNewNameSpace
        displayNameSpaceInfo("REG")
        goFocus('ns_reg');// 해당 영역으로 scroll
    }else if ( targetAction == "DeleteNameSpace"){
        deleteNameSpace ()
    }else if ( targetAction == "AddNewVpc"){
        displayVNetInfo("REG")
        goFocus('vnetCreateBox');
    }else if ( targetAction == "DeleteVpc"){
        deleteVPC()
    }else if ( targetAction == "AddNewSecurityGroup"){
        displaySecurityGroupInfo("REG")
        goFocus('securityGroupCreateBox');
    }else if ( targetAction == "DeleteSecurityGroup"){
        deleteSecurityGroup()
    }else if ( targetAction == "AddNewSshKey"){
        displaySshKeyInfo("REG")
        goFocus('sshKeyCreateBox');
    }else if ( targetAction == "DeleteSshKey"){
        deleteSshKey()
    }else if ( targetAction == "AddNewVirtualMachineImage"){
        displayVirtualMachineImageInfo("REG")
        goFocus('virtualMachineImageCreateBox');
    }else if ( targetAction == "DeleteVirtualMachineImage"){
        deleteVirtualMachineImage()
    }else if ( targetAction == "FetchImages"){
        getCommonFetchImages();         
    }else if ( targetAction == "AddNewVmSpec"){
        displayVmSpecInfo("REG")
        goFocus('vmSpecCreateBox');
    }else if ( targetAction == "ExportVmScriptOfMcis"){
        vmScriptExport();
    }else if ( targetAction == "DeleteVmSpec"){
        deleteVmSpec();  
    }else if ( targetAction == "FetchSpecs"){
        var connectionName = $("#regConnectionName").val();
        putFetchSpecs(connectionName);         
    }else if ( targetAction == "GotoMonitoringPerformance"){
        // alert("모니터링으로 이동 GotoMonitoringPerformance")
        // location.href ="";//../operation/Monitoring_Mcis.html
        var targetUrl = "/operation/monitorings/mcismng/mngform"
        changePage(targetUrl)
    }else if ( targetAction == "GotoMonitoringFault"){
        // alert("모니터링으로 이동 GotoMonitoringFault")
        // location.href ="";//../operation/Monitoring_Mcis.html
        var targetUrl = "/operation/monitorings/mcismng/mngform"
        changePage(targetUrl)
    }else if ( targetAction == "GotoMonitoringCost"){
        // alert("모니터링으로 이동 GotoMonitoringCost")
        // location.href ="";//../operation/Monitoring_Mcis.html
        var targetUrl = "/operation/monitorings/mcismng/mngform"
        changePage(targetUrl)
    }else if ( targetAction == "GotoMonitoringUtilize"){
        // alert("모니터링으로 이동 GotoMonitoringUtilize")
        // location.href ="";//../operation/Monitoring_Mcis.html    
        var targetUrl = "/operation/monitorings/mcismng/mngform"
        changePage(targetUrl)
    }else if ( targetAction == "McisLifeCycleReboot"){
        callMcisLifeCycle('reboot')
    }else if ( targetAction == "McisLifeCycleSuspend"){
        callMcisLifeCycle('suspend')
    }else if ( targetAction == "McisLifeCycleResume"){
        callMcisLifeCycle('resume')
    }else if ( targetAction == "McisLifeCycleTerminate"){
        callMcisLifeCycle('terminate')
    }else if ( targetAction == "McisManagement"){
        alert("수행할 function 정의되지 않음");
    }else if ( targetAction == "MoveToMcisManagementFromDashboard"){
        var mcisID = $("#mcis_id").val();
        var targetUrl = "/operation/manages/mcismng/mngform?mcisid=" + mcisID;
        changePage(targetUrl)
    }else if ( targetAction == "MoveToMcisManagement"){
        var targetUrl = "/operation/manages/mcismng/mngform";
        changePage(targetUrl)
    }else if ( targetAction == "AddNewMcis"){
        // $('#loadingContainer').show();
        // location.href ="/operation/manages/mcis/regform/";
        var targetUrl = "/operation/manages/mcismng/regform";
        changePage(targetUrl)
    }else if ( targetAction == "DeleteMcis"){
        deleteMCIS();
        
    }else if ( targetAction == "ImportScriptOfMcis"){
        mcisScriptImport();
    }else if ( targetAction == "ExportScriptOfMcis"){
        mcisScriptExport();
    }else if ( targetAction == "VmLifeCycle"){
        alert("수행할 function 정의되지 않음");
    }else if ( targetAction == "VmLifeCycleReboot"){
        vmLifeCycle('reboot')
    }else if ( targetAction == "VmLifeCycleSuspend"){
        vmLifeCycle('suspend')
    }else if ( targetAction == "VmLifeCycleResume"){
        vmLifeCycle('resume')
    }else if ( targetAction == "VmLifeCycleTerminate"){
        vmLifeCycle('terminate')
    }else if ( targetAction == "VmManagement"){
        alert("수행할 function 정의되지 않음");
    }else if ( targetAction == "AddNewVm"){
        addNewVirtualMachine()
    }else if ( targetAction == "AddNewVmOfMcis"){
        addNewVirtualMachine()
    }else if ( targetAction == "ExportVmScriptOfMcis"){
        vmScriptExport();
    }else if ( targetAction == "--"){
        addNewVirtualMachine()
    }else if ( targetAction == "monitoringConfigPolicyConfig"){
        regMonitoringConfigPolicy()
    }else if ( targetAction == "DifferentConnection"){
        setAndClearByDifferentConnectionName();
    }else if ( targetAction == "DifferentConnectionAtSecurityGroup"){
        uncheckDifferentConnectionAtSecurityGroup();
    }else if ( targetAction == "AddMonitoringAlertPolicy"){
        addMonitoringAlertPolicy();
    }else if ( targetAction == "DeleteMonitoringAlertPolicy"){
        deleteMonitoringAlertPolicy();
    }else if ( targetAction == "AddNewMcks"){
        var targetUrl = "/operation/manages/mcksmng/regform";
        changePage(targetUrl)
    }else if ( targetAction == "AddNewNodeOfMcks"){
        addNewNode();
    }else if ( targetAction == "DeleteNodeOfMcks"){
        deleteNodeOfMcks();        
    }else if ( targetAction == "AddMonitoringAlertEventHandler"){
        addMonitoringAlertEventHandler();
    }else if ( targetAction == "deleteMonitoringAlertEventHandler"){
        deleteMonitoringAlertEventHandler();
    }else if ( targetAction == "DeleteMcks"){
        deleteMCKS();
    }else {
        alert("수행할 function 정의되지 않음 " + targetAction);
    }
    console.log("commonConfirmOk " + targetAction);
    commonConfirmClose();
}

//confirm modal창에서 cancel 버튼 클릭시 수행할 method 지정. 그냥 창만 듣을 경우에는 commonModalClose() 호출
var rollbackObjArr = [];
function commonConfirmCancel(targetAction){
    console.log("commonConfirmCancel : " + targetAction)
    //
    if( targetAction == 'DifferentConnection'){
        // set 했던것들 초기화.
        for( var i = 0; i < rollbackObjArr.length; i++){
            $("#" + rollbackObjArr[i]).val('');
        }
    }
    commonConfirmClose();
}
// confirm modal창 닫기. setting값 초기화
function commonConfirmClose(){
    $('#confirmText').text('');
    $('#confirmOkAction').val('');
    // $('#modalArea').hide(); 
    $("#confirmArea").modal("hide");
}

//////// Prompt start ////////
// confirm modal창 보이기 modal창이 열릴 때 해당 창의 text 지정, close될 때 action 지정
function commonPromptOpen(targetAction, targetObjId){
    console.log("commonPromptOpen : " + targetAction)
    
    let promptModalTextMap = new Map(
        [
            ["FilterName", "필터링할 단어를 입력하세요"],
            ["FilterCloudProvider", "필터링할 단어를 입력하세요"],
            ["FilterDriver", "필터링할 단어를 입력하세요"],
            ["FilterCredential", "필터링할 단어를 입력하세요"],
            ["RsFltVPCName", "필터링할 단어를 입력하세요"],
            ["RsFltCIDRBlock", "필터링할 단어를 입력하세요"],
            ["RsFltSecurityGroupName", "필터링할 단어를 입력하세요"],
            ["RsFltConnectionName", "필터링할 단어를 입력하세요"],
            ["RsFltSshName", "필터링할 단어를 입력하세요"],
            ["RsFltSshConnName", "필터링할 단어를 입력하세요"],
            ["RsFltSshKeyName", "필터링할 단어를 입력하세요"],
            ["RsFltSrvImgId", "필터링할 단어를 입력하세요"],
            ["RsFltSrvImgName", "필터링할 단어를 입력하세요"],
            ["RsFltSrvSpecName", "필터링할 단어를 입력하세요"],
            ["RsFltSrvSpecConnName", "필터링할 단어를 입력하세요"],
            ["RsFltSrvCspSpecName", "필터링할 단어를 입력하세요"],
            ["NSFltName", "필터링할 단어를 입력하세요"],
            ["NSFltId", "필터링할 단어를 입력하세요"],
            ["NSFltDescription", "필터링할 단어를 입력하세요"],
            ["AlertPolicyName", "필터링할 단어를 입력하세요"],
            ["AlertPolicyMeasurement", "필터링할 단어를 입력하세요"],
            ["AlertPolicyTargetType", "필터링할 단어를 입력하세요"],
            ["AlertPolicyEventType", "필터링할 단어를 입력하세요"],
            ["Config", "Would you like to set Cloud config ?"],
            ["FilterMcisName", "필터링할 단어를 입력하세요"],
            ["FilterMcisStatus", "필터링할 단어를 입력하세요"],
            ["FilterMcisDesc", "필터링할 단어를 입력하세요"],
            ["OprMngMcksStatus", "필터링할 단어를 입력하세요"],
            ["OprMngMcksName", "필터링할 단어를 입력하세요"],
            ["OprMngMcksNetworkCni", "필터링할 단어를 입력하세요"],

            ["RemoteCommandMcis", "Please enter a command to execute"],
            ["RemoteCommandVmOfMcis", "Please enter a command to execute"],
            
        ]
    );
    console.log(promptModalTextMap.get(targetAction));
    try{
        $('#promptQuestion').html(promptModalTextMap.get(targetAction));
        $('#promptText').val('');

        $('#promptTargetObjId').val(targetObjId);
        $('#promptOkAction').val(targetAction);// Prompt입력창에서 OK버튼을 눌렀을 때 이동할 targetKey
                
        $('#promptArea').modal();
        $("#promptArea").on('shown.bs.modal', function () {
            $(this).find('#promptText').focus();
        });
    }catch(e){
        console.log(e);
        alert(e);
    }
}

function commonPromptOk(){
    var targetAction = $('#promptOkAction').val();
    var targetObjId = $('#promptTargetObjId').val();
    var targetValue = $('#promptText').val();

    console.log("promptOkAction : " + targetAction)
    if( targetAction == 'FilterName'){// Name이라는 Column을 Filtering
        if( targetValue ){
            filterTable(targetObjId, "Name", targetValue)
        }        
    }else if( targetAction == 'FilterCloudProvider'){// Name이라는 Column을 Filtering
        if( targetValue ){
            filterTable(targetObjId, "Cloud Provider", targetValue)
        }       
    }else if( targetAction == 'FilterDriver'){// Name이라는 Column을 Filtering
        if( targetValue ){
            filterTable(targetObjId, "Driver", targetValue)
        }       
    }else if( targetAction == 'FilterCredential'){// Name이라는 Column을 Filtering
        if( targetValue ){
            filterTable(targetObjId, "Credential", targetValue)
        }       
    }else if( targetAction == 'RsFltVPCName'){// Name이라는 Column을 Filtering
        if( targetValue ){
            filterTable(targetObjId, "VPC Name", targetValue)
        }       
    }else if( targetAction == 'RsFltCIDRBlock'){// Name이라는 Column을 Filtering
        if( targetValue ){
            filterTable(targetObjId, "CIDR Block", targetValue)
        }       
    }else if( targetAction == 'RsFltSecurityGroupName'){// Name이라는 Column을 Filtering
        if( targetValue ){
            filterTable(targetObjId, "SecurityGroup Name", targetValue)
        }       
    }else if( targetAction == 'RsFltConnectionName'){// Name이라는 Column을 Filtering
        if( targetValue ){
            filterTable(targetObjId, "Connection Name", targetValue)
        }       
    }else if( targetAction == 'RsFltSshName'){// Name이라는 Column을 Filtering
        if( targetValue ){
            filterTable(targetObjId, "Name", targetValue)
        }       
    }else if( targetAction == 'RsFltSshConnName'){// Name이라는 Column을 Filtering
        if( targetValue ){
            filterTable(targetObjId, "Connection Name", targetValue)
        }       
    }else if( targetAction == 'RsFltSshKeyName'){// Name이라는 Column을 Filtering
        if( targetValue ){
            filterTable(targetObjId, "SSH KEY Name", targetValue)
        }       
    }else if( targetAction == 'RsFltSrvImgId'){// Name이라는 Column을 Filtering
        if( targetValue ){
            filterTable(targetObjId, "Image ID", targetValue)
        }       
    }else if( targetAction == 'RsFltSrvImgName'){// Name이라는 Column을 Filtering
        if( targetValue ){
            filterTable(targetObjId, "Image Name", targetValue)
        }       
    }else if( targetAction == 'RsFltSrvSpecName'){// Name이라는 Column을 Filtering
        if( targetValue ){
            filterTable(targetObjId, "Name", targetValue)
        }       
    }else if( targetAction == 'RsFltSrvSpecConnName'){// Name이라는 Column을 Filtering
        if( targetValue ){
            filterTable(targetObjId, "Connection Name", targetValue)
        }       
    }else if( targetAction == 'RsFltSrvCspSpecName'){// Name이라는 Column을 Filtering
        if( targetValue ){
            filterTable(targetObjId, "CSP Spec Name", targetValue)
        }       
    }else if( targetAction == 'NSFltName'){// Name이라는 Column을 Filtering
        if( targetValue ){
            filterTable(targetObjId, "Name", targetValue)
        }       
    }else if( targetAction == 'NSFltId'){// Name이라는 Column을 Filtering
        if( targetValue ){
            filterTable(targetObjId, "ID", targetValue)
        }       
    }else if( targetAction == 'NSFltDescription'){// Name이라는 Column을 Filtering
        if( targetValue ){
            filterTable(targetObjId, "description", targetValue)
        }       
    }else if( targetAction == 'AlertPolicyName'){// Name이라는 Column을 Filtering
        if( targetValue ){
            filterTable(targetObjId, "Name", targetValue)
        }       
    }else if( targetAction == 'AlertPolicyMeasurement'){// Name이라는 Column을 Filtering
        if( targetValue ){
            filterTable(targetObjId, "Measurement", targetValue)
        }       
    }else if( targetAction == 'AlertPolicyTargetType'){// Name이라는 Column을 Filtering
        if( targetValue ){
            filterTable(targetObjId, "Target Type", targetValue)
        }       
    }else if( targetAction == 'AlertPolicyEventType'){// Name이라는 Column을 Filtering
        if( targetValue ){
            filterTable(targetObjId, "Alert Event Type", targetValue)
        }       
    }else if( targetAction == 'FilterMcisName'){// Name이라는 Column을 Filtering
        if( targetValue ){
            filterTable(targetObjId, "Name", targetValue)
        }       
    }else if( targetAction == 'FilterMcisStatus'){// Status이라는 Column을 Filtering
        if( targetValue ){
            filterTable(targetObjId, "Status", targetValue)
        }       
    }else if( targetAction == 'FilterMcisDesc'){// Description이라는 Column을 Filtering
        if( targetValue ){
            filterTable(targetObjId, "Description", targetValue)
        }       
    }else if( targetAction == 'OprMngMcksStatus'){// Description이라는 Column을 Filtering
        console.log("OprMngMcksStatus");
        if( targetValue ){
            filterTable(targetObjId, "Status", targetValue)
        }       
    }else if( targetAction == 'OprMngMcksName'){// Description이라는 Column을 Filtering
        if( targetValue ){
            filterTable(targetObjId, "Name", targetValue)
        }       
    }else if( targetAction == 'OprMngMcksNetworkCni'){// Description이라는 Column을 Filtering
        if( targetValue ){
            filterTable(targetObjId, "NetworkCni", targetValue)
        } 
    }else if( targetAction == 'RemoteCommandMcis'){
        if( targetValue ){
            remoteCommandMcis(targetValue);
            //postRemoteCommandMcis(targetValue);
        }
    }else if( targetAction == 'RemoteCommandVmOfMcis'){
        if( targetValue ){
            remoteCommandVmMcis(targetValue);
        }
    }

   
    commonPromptClose();
}

function commonPromptClose(){
    $('#promptQuestion').text('');
    $('#promptText').text('');
    $('#promptOkAction').val('');
    $("#promptArea").modal("hide");
}
//////// Prompt end //////////
// provider에 등록된 connection을 selectbox에 표시
function getConnectionListForSelectbox(provider, targetSelectBoxID){
    
    var data = new Array();
    var url = "/setting/connections/cloudconnectionconfig/" + "list"
    console.log("provider : ",provider)
    var html = "";
    axios.get(url,{
        headers:{
            // 'Authorization': apiInfo
        }
    }).then(result=>{
        console.log('getConnectionConfig result: ',result)
        data = result.data.ConnectionConfig
        console.log("set data array " + data.length);
        
        console.log("connection data : ",data);
        var count = 0; 
        var configName = "";
        var confArr = new Array();
        html +='<option selected>Select Configname</option>';
        for(var i in data){
            if(provider == data[i].ProviderName){ 
                count++;
                html += '<option value="'+data[i].ConfigName+'" item="'+data[i].ProviderName+'">'+data[i].ConfigName+'</option>';
                configName = data[i].ConfigName
                confArr.push(data[i].ConfigName)                
            }
        }
        if(count == 0){
            commonAlert("해당 Provider에 등록된 Connection 정보가 없습니다.")            
        }
        console.log("targetSelectBoxID = " + targetSelectBoxID)
        $("#" + targetSelectBoxID).empty();
        $("#" + targetSelectBoxID).append(html);

        if(confArr.length > 1){
            configName = confArr[0];
            console.log("chage value")
            // 0번째 자동으로 선택하여 vNetID목록 갱신
            // $("#" + targetSelectBoxID + " option[value=" + configName + "]").prop('selected', 'selected').change();
            $("#" + targetSelectBoxID + " option[value=" + configName + "]").prop('selected', true).change();         
        }
        // getVnetInfoListForSelectbox(configName);
    // }).catch(function(error){
    //     console.log("Network data error : ",error);        
    // });
    }).catch((error) => {
        console.warn(error);
        console.log(error.response)
        var errorMessage = error.response.data.error;
        var statusCode = error.response.status;
        commonErrorAlert(statusCode, errorMessage) 
    });
}

// connection에 등록된 vnet List를 selectbox에 표시
function getVnetInfoListForSelectbox(configName, targetSelectBoxID){
    console.log("vnet : ", configName);
    
    var url = "/setting/resources" + "/network/list"
    var html = "";
    axios.get(url,{
        headers:{
            // 'Authorization': apiInfo
        }
    }).then(result=>{
        data = result.data.VNetList;
        console.log("vNetwork Info : ",result);
        console.log("vNetwork data : ",data);
        var count = 0; 
        for(var i in data){
            count++;
            if(data[i].connectionName == configName){
                html += '<option value="'+data[i].id+'" selected>'+data[i].cspVNetName+'('+data[i].id+')</option>'; 
            }
        }

        if( count == 0){
            commonAlert("해당 Provider에 등록된 Connection 정보가 없습니다.")
                html +='<option selected>Select Configname</option>';
        }
    
        $("#" + targetSelectBoxID).empty();
        $("#" + targetSelectBoxID).append(html);  
    })
}

function getProviderNameByConnection(configName, targetObjID){
    console.log("configName : ", configName);
    
    var url = "/setting/connections" + "/cloudconnectionconfig/" + configName
    axios.get(url,{
        headers:{
            // 'Authorization': apiInfo
        }
    }).then(result=>{
        data = result.data.ConnectionConfig;
        console.log("connection data : ",data);
        var providerName = data.ProviderName
        console.log("providerName : ",providerName);
        $("#" + targetObjID).val(providerName);
        
    })
}

function getRegionListByProviderForSelectbox(provider, targetObjID){
    console.log("getRegionListByProviderForSelectbox : ", provider);
    
    var url = "/setting/connections" + "/region/" + configName
    axios.get(url,{
        headers:{
            // 'Authorization': apiInfo
        }
    }).then(result=>{
        data = result.data.ConnectionConfig;
        console.log("connection data : ",data);
        var providerName = data.ProviderName
        console.log("providerName : ",providerName);
        $("#" + targetObjID).val(providerName);
        
    })
}


function isEmpty(str){
	if(typeof str == "undefined" || str == null || str == "")
		return true;
	else
		return false ;
}


// table의 column별로 sortType을 달리 가져간다.
// TODO : sortType을 바꾸고 table정렬을 바로 할 것인지? sort action을 통해 정렬을 할 것인지..
// tableId : 대상 tableID
// columnName : 정렬하려는 column의 text
// sorType을 찾아 사용하고, 사용한 뒤에는 반대되는 것을 넣음.
// - changerSortType 에서 정렬방식 선택, 정렬할 column의 index 찾기 -> table 에서 정렬할 column index에 해당 하는 column을 sort
var tableSortTypeletMap = new Map();
function getSortType(tableId, columnName){
    var sortType = tableSortTypeletMap.get(tableId + "|" + columnName);
    if(!sortType){
        sortType = "asc" // default
        tableSortTypeletMap.set(tableId + "|" + columnName, sortType);
    }

    var returnSortType = (sortType === 'asc') ? 'desc':'asc';
    tableSortTypeletMap.set(tableId + "|" + columnName, returnSortType);
    return returnSortType// 현재 set 된 sortType의 반대를 return    
}


// tr에 정의된 column 이름으로 해당 column의 index를 찾는다.
// table 밑의 첫번째 tr에서 해당 이름을 찾음.
function getTableColumnIndex(tableId, columnName){
    var tableObj = $('#' + tableId);
    // console.log(tableObj)
    var checkSort = true;
    var rows = tableObj[0].rows;
    var columns = rows[0].cells// 첫번째 tr 
    console.log(columns);
    var columnIndex = 0;
    for (var i = 0; i < columns.length; i++) {
        var columnText = columns[i];
        console.log(columnName + ":" + columnText.innerHTML)
        if( columnName == columnText.innerText){
            columnIndex = i;
            break;
        }        
    }
    return columnIndex;
}

// table tag에 id를 줘야 한다. columnName의 첫번째 tr 아래에 있는 cell(th, td)의 text
function tableSort(tableId, columnName){
    var sortTargetColumnIndex = getTableColumnIndex(tableId, columnName)
    var sortType = getSortType(tableId, columnName);
    console.log(tableId + " : " + columnName + " : " + sortTargetColumnIndex + " : " + sortType)
    var tableObj = $('#' + tableId);
    // console.log(tableObj)
    var checkSort = true;
    var rows = tableObj[0].rows;
    // console.log(rows);
    while (checkSort) { // 현재와 다음만 비교하기때문에 위치변경되면 다시 정렬해준다.
        checkSort = false;

        for (var i = 1; i < (rows.length - 1); i++) {
            console.log("***** " + sortTargetColumnIndex + ", " + sortType)
            console.log(rows[i].cells[sortTargetColumnIndex].innerText);
            var fCell = rows[i].cells[sortTargetColumnIndex].innerText.toUpperCase();
            var sCell = rows[i + 1].cells[sortTargetColumnIndex].innerText.toUpperCase();

            var row = rows[i];

            // 오름차순<->내림차순 ( 이부분이 이해 잘안됬는데 오름차순이면 >, 내림차순이면 < 이고 if문의 내용은 동일하다 )
            if ( (sortType == 'asc' && fCell > sCell) || 
                    (sortType == 'desc' && fCell < sCell) ) {

                row.parentNode.insertBefore(row.nextSibling, row);
                checkSort = true;
            }
        }
    }    
}


// todo : fintering을 하려면 keyword를 입력 받아야 하는데???
// filter 항목에서 column을 선택하면 popup으로 keyword를 입력받아 filterTable()을 실행하게 하면 될 까?
// 1. 대상 table에 ID가 있어야 함.
// 2. filter > 대상 칼럼을 선택 시 > txt 입력창이 떠서 keyword를 입력하면 해당 내용으로 filtering
// 3. 입력 단어가 ALL 이면 모두 보여준다.
function filterTable(tableId, filterColumnName, filterKeyword){
    var filterTargetColumnIndex = getTableColumnIndex(tableId, filterColumnName)
    console.log("filterTargetColumnIndex=" + filterTargetColumnIndex);
    var filter = filterKeyword.toUpperCase();
	console.log("filter=" + filter);
    
    //var tableObj = $('#' + tableId);
	var tableObj = document.getElementById(tableId);
	var trObj = tableObj.getElementsByTagName("tr");
    //var rows = tableObj[0].rows;
	console.log(trObj.length);
    // Loop through all table rows, and hide those who don't match the search query
    // 찾은 column을 기준으로 fintering한다.
    for (i = 1; i < trObj.length; i++) {
		console.log(trObj[i]);
        var tdTag = trObj[i].getElementsByTagName("td")[filterTargetColumnIndex];
        console.log(tdTag);
        if (tdTag) {
            txtValue = tdTag.textContent || tdTag.innerText;
            console.log(txtValue + " = " + tdTag.textContent + " || " + tdTag.innerText);
            if(filter == "ALL") {
                trObj[i].style.display = "";			
			} else if (txtValue.toUpperCase().indexOf(filter) > -1) {
				trObj[i].style.display = "";
            }else {
				trObj[i].style.display = "none";
			}
        }
    }
}

// table에서 hidden으로 설정된 obj를 기준으로 filterling. 보이고 안보이고
function filterTableByHiddenColumn(tableId, hiddenColumnName, filterKeyword){

    var keyword = filterKeyword.toUpperCase();
	console.log("filter=" + keyword);

    var trs = $('#' + tableId + ' tr');
    console.log(trs);
    //for (var i = 1; i < $('#' + tableId + ' tr').size(); i++) {
    for (var i = 1; i < trs.size(); i++) {
        //var hiddenval = trs.eq(i).find('input:hidden[name="vmImageInfo"]').val();
        var hiddenval = trs.eq(i).find('input:hidden[name="' + hiddenColumnName + '"]').val();
        // console.log("hiddenval " + hiddenval);

        if(keyword == "ALL") {
            trs.eq(i).css("display", "");
        }else if (hiddenval.toUpperCase().indexOf(keyword) > -1) {
            trs.eq(i).css("display", "");
        }else {
            trs.eq(i).css("display", "none");
        }
    }    
}