$(document).ready(function(){
    //OS_HW popup table scrollbar
    $('#OS_HW .btn_spec').on('click', function() {
        console.log("os_hw bpn_spec clicked ** ")
        $('#OS_HW_Spec .dtbox.scrollbar-inner').scrollbar();

        // connection 정보 set
        var esSelectedProvider = $("#es_regProvider option:selected").val();
        var esSelectedRegion = $("#es_regRegion option:selected").val();
        var esSelectedConnectionName = $("#es_regConnectionName option:selected").val();

        console.log("OS_HW_Spec_Assist click");
        if( esSelectedProvider){
            $("#assist_select_provider").val(esSelectedProvider);
        }
        if( esSelectedRegion){
            $("#assist_select_resion").val(esSelectedRegion);
        }
        if( esSelectedConnectionName){
            $("#assist_select_connectionName").val(esSelectedConnectionName);
        }

        console.log("esSelectedProvider = " + esSelectedProvider + " : " + $("#assist_select_provider").val());
        console.log("esSelectedRegion = " + esSelectedRegion + " : " + $("#assist_select_resion").val());
        console.log("esSelectedConnectionName = " + esSelectedConnectionName + " : " + $("#assist_select_connectionName").val());
    });
    //Security popup table scrollbar
    $('#Security .btn_edit').on('click', function() {
    $("#security_edit").modal();
        $('#security_edit .dtbox.scrollbar-inner').scrollbar();
    });

    // $("input[name='vmInfoType']:radio").change(function () {
    //     //라디오 버튼 값을 가져온다.
    //     var formType = this.value;
                
    // });


    // server add 버튼 클릭 시
    // $('.servers_box .server_add').click(function(){	

    //     //<div class="servers_config import_servers_config" id="importServerConfig">
    //     //<div class="servers_config new_servers_config" id="expertServerConfig">
    // });

    //Servers Expert on/off
//     var check = $(".switch .ch");
//     var $Servers = $(".servers_config");
//     var $NewServers = $(".new_servers_config");
//     var $SimpleServers = $(".simple_servers_config");
//     var simple_config_cnt = 0;
//     var expert_config_cnt = 0;
    
//     check.click(function(){
//         $(".switch span.txt_c").toggle();
//         $NewServers.removeClass("active");
//     });
   
//   //Expert add
//     $('.servers_box .server_add').click(function(){
//         $NewServers.toggleClass("active");
//       if($Servers.hasClass("active")) {
//         $Servers.toggleClass("active");
//     } else {
//         $Servers.toggleClass("active");
//     }
//     });
//     // Simple add
//   $(".servers_box .switch").change(function() {
//     if ($(".switch .ch").is(":checked")) {	
//             $('.servers_box .server_add').click(function(){	
                
//                 $NewServers.addClass("active");
//                 $SimpleServers.removeClass("active");		
//             });
//     } else {
//             $('.servers_box .server_add').click(function(){
            
//                 $NewServers.removeClass("active");
//                 $SimpleServers.addClass("active");
            
            
//             });		
//     }
//   });
});


var totalDeployServerCount = 0;
function btn_deploy(){
    var mcis_name = $("#mcis_name").val();
    var mcis_id = $("#mcis_id").val();
    if(!mcis_id){
        commonAlert("Please Select MCIS !!!!!")
        return;
    }
    totalDeployServerCount = 0;// deploy vm 개수 초기화

    console.log(Simple_Server_Config_Arr);
    if(Simple_Server_Config_Arr){// mcissimpleconfigure.js 에 const로 정의 됨.
        var vm_len = Simple_Server_Config_Arr.length;
        if( vm_len > 0){
            totalDeployServerCount += vm_len
            console.log("Simple_Server_Config_Arr length: ",vm_len);
            // var new_obj = {}
            // new_obj['vm'] = Simple_Server_Config_Arr;
            // console.log("new obj is : ",new_obj);
            // var url = "/operation/manages/mcis/:mcisID/vm/reg/proc"
            var url = "/operation/manages/mcismng/" + mcis_id +"/vm/reg/proc"

            // 한개씩 for문으로 추가
            for(var i in Simple_Server_Config_Arr){
                new_obj = Simple_Server_Config_Arr[i];
                console.log("new obj is : ",new_obj);
                try{
                    resultVmCreateMap.set(new_obj.name, "")
                    axios.post(url,new_obj,{
                        headers :{
                            },
                    }).then(result=>{
                        console.log("MCIR VM Register data : ",result);
                        console.log("Result Status : ",result.status); 

                        var statusCode = result.data.status;
                        var message = result.data.message;
                        console.log("Result Status : ",statusCode); 
                        console.log("Result message : ",message); 

                        if(result.status == 201 || result.status == 200){
                            vmCreateCallback(new_obj.name, "Success")
                        //     commonAlert("Register Success")
                        //     // location.href = "/Manage/MCIS/list";
                        //     // $('#loadingContainer').show();
                        //     // location.href = "/operation/manages/mcis/mngform/"
                        //     var targetUrl = "/operation/manages/mcis/mngform"
                        //     changePage(targetUrl)
                        }else{
                            vmCreateCallback(new_obj.name, message)    
                        //     commonAlert("Register Fail")
                        //     //location.reload(true);
                        }
                    }).catch((error) => {
                        // console.warn(error);
                        console.log(error.response)
                        var errorMessage = error.response.data.error;
                        // commonErrorAlert(statusCode, errorMessage) 
                        vmCreateCallback(new_obj.name, errorMessage)
                    })
                }finally{
                    
                }

                // post로 호출을 했으면 해당 VM의 정보는 비활성시킨 후(클릭 Evnet 안먹게)
                // 상태값을 모니터링 하여 결과 return 까지 대기.
                // return 받으면 해당 VM
            }
        }
    }

    ///////// export
    console.log(Expert_Server_Config_Arr);
    if(Expert_Server_Config_Arr){
        var vm_len = Expert_Server_Config_Arr.length;			
        console.log("Expert_Server_Config_Arr length: ",vm_len);
        if( vm_len > 0){
            totalDeployServerCount += vm_len
            // var new_obj = {}
            // new_obj['vm'] = Simple_Server_Config_Arr;
            // console.log("new obj is : ",new_obj);
            // var url = "/operation/manages/mcis/:mcisID/vm/reg/proc"
            var url = "/operation/manages/mcismng/" + mcis_id +"/vm/reg/proc"

            // 한개씩 for문으로 추가
            for(var i in Expert_Server_Config_Arr){
                new_obj = Expert_Server_Config_Arr[i];
                console.log("new obj is : ",new_obj);
                try{
                    resultVmCreateMap.set("Expert"+ i, "")
                    axios.post(url,new_obj,{
                        headers :{
                            },
                    }).then(result=>{
                        console.log("MCIR VM Register data : ",result);
                        console.log("Result Status : ",result.status); 

                        var statusCode = result.data.status;
                        var message = result.data.message;
                        console.log("Result Status : ",statusCode); 
                        console.log("Result message : ",message); 

                        if(result.status == 201 || result.status == 200){
                            vmCreateCallback("Expert"+ i, "Success")                   
                        }else{
                            vmCreateCallback("Expert"+ i, message)
                        }
                    }).catch((error) => {
                        // console.warn(error);
                        console.log(error.response)
                        var errorMessage = error.response.data.error;
                        // commonErrorAlert(statusCode, errorMessage) 
                        vmCreateCallback("Expert"+ i, errorMessage)
                    })
                }finally{
                    
                }

                // post로 호출을 했으면 해당 VM의 정보는 비활성시킨 후(클릭 Evnet 안먹게)
                // 상태값을 모니터링 하여 결과 return 까지 대기.
                // return 받으면 해당 VM
            }
        }
    }
    ///////// import
    if(Import_Server_Config_Arr){// mcissimpleconfigure.js 에 const로 정의 됨.
        // TODO : 어차피 simple/expert와 로직이 다른데... 
        // json 그대로 넘기도록
        var vm_len = Import_Server_Config_Arr.length;
        if( vm_len > 0){
            console.log("Import_Server_Config_Arr length: ",vm_len);
            totalDeployServerCount += vm_len
            // var new_obj = {}
            // new_obj['vm'] = Simple_Server_Config_Arr;
            // console.log("new obj is : ",new_obj);
            // var url = "/operation/manages/mcis/:mcisID/vm/reg/proc"
            var url = "/operation/manages/mcismng/" + mcis_id +"/vm/reg/proc"

            // 한개씩 for문으로 추가
            for(var i in Import_Server_Config_Arr){
                new_obj = Import_Server_Config_Arr[i];
                console.log("new obj is : ",new_obj);
                try{
                    resultVmCreateMap.set("Import"+ i, "")
                    axios.post(url,new_obj,{
                        headers :{
                            },
                    }).then(result=>{
                        console.log("MCIR VM Register data : ",result);
                        console.log("Result Status : ",result.status); 

                        var statusCode = result.data.status;
                        var message = result.data.message;
                        console.log("Result Status : ",statusCode); 
                        console.log("Result message : ",message); 

                        if(result.status == 201 || result.status == 200){
                            vmCreateCallback("Import"+ i, "Success")
                        //     commonAlert("Register Success")
                        //     // location.href = "/Manage/MCIS/list";
                        //     // $('#loadingContainer').show();
                        //     // location.href = "/operation/manages/mcis/mngform/"
                        //     var targetUrl = "/operation/manages/mcis/mngform"
                        //     changePage(targetUrl)
                        }else{
                            vmCreateCallback("Import"+ i, message)    
                        //     commonAlert("Register Fail")
                        //     //location.reload(true);
                        }
                    }).catch((error) => {
                        // console.warn(error);
                        console.log(error.response)
                        var errorMessage = error.response.data.error;
                        // commonErrorAlert(statusCode, errorMessage) 
                        vmCreateCallback("Import"+ i, errorMessage)
                    })
                }finally{
                    
                }

                // post로 호출을 했으면 해당 VM의 정보는 비활성시킨 후(클릭 Evnet 안먹게)
                // 상태값을 모니터링 하여 결과 return 까지 대기.
                // return 받으면 해당 VM
            }
        }
    }
}

// Import / Export Modal 표시
function btn_ImportExport() {
    // export할 VM을 선택한 후 export 버튼 누르라고...
    $("#VmImportExport").modal();
    $('#VmImportExport .dtbox.scrollbar-inner').scrollbar();
}

// vm 생성 결과 표시
// 여러개의 vm이 생성될 수 있으므로 각각 결과를 표시
var resultVmCreateMap = new Map();
function vmCreateCallback(resultVmKey, resultStatus){
    resultVmCreateMap.set(resultVmKey, resultStatus)
    var resultText = "";
    for (let key of resultVmCreateMap.keys()) { 
        console.log("vmCreateresult " + key + " : " + resultVmCreateMap.get(resultVmKey) );
        resultText += key + " = " + resultVmCreateMap.get(resultVmKey) + ","
        totalDeployServerCount--
    }

    $("#serverRegistResult").text(resultText);

    if( resultStatus != "Success"){
        // add된 항목 제거 해야 함.

        // array는 초기화
        Simple_Server_Config_Arr.length = 0;
        simple_data_cnt = 0
        // TODO : expert 추가하면 주석 제거할 것
        Expert_Server_Config_Arr.length = 0;
        expert_data_cnt = 0
        Import_Server_Config_Arr.length = 0;
        import_data_cnt = 0
    }

    if( totalDeployServerCount == 0){
        //getVmList();
        commonAlert($("#serverRegistResult").text());
        var targetUrl = "/operation/manages/mcismng/mngform"
        changePage(targetUrl)
    }
}

// 현재 mcis의 vm 목록 조회 : deploy후 상태볼 때 사용
function getVmList(){
    var mcis_id = $("#mcis_id").val();
    
    
    // /operation/manages/mcis/:mcisID
    var url = "/operation/manages/mcismng/" + mcis_id 
    axios.get(url,{})
    .then(result=>{
        console.log("MCIR VM Register data : ",result);
        console.log("Result Status : ",result.status); 

        var statusCode = result.data.status;
        var message = result.data.message;
        //
        console.log("Result Status : ",statusCode); 
        console.log("Result message : ",message); 


        if(result.status == 201 || result.status == 200){
            var mcis = result.data.McisInfo
            console.log(mcis)

            
            var vms = mcis.vm
            if(vms){
                vm_len = vms.length

                $("#mcis_server_list *").remove();
                var appendLi = "";

                for(var o in vms){
                    var vm_status = vms[o].status
                    var vm_name = vms[o].name

                    console.log(o + "번째 " + vm_name + " : " + vm_status)
                    // mcis_server_list 밑의 li들을 1개빼고 삭제. 
                    // 가져온 vm list 를 add? (1개는 더하기 버튼이므로)                    
                                    
                    
                    appendLi = appendLi + '<li>';
                    appendLi = appendLi + '<div class="server server_on bgbox_g">';
                    appendLi = appendLi + '<div class="icon"></div>';
                    appendLi = appendLi + '<div class="txt">' + vm_name + '</div>';
                    appendLi = appendLi + '</li>';

                    appendLi = appendLi + '</li>';                
                    
                }
                appendLi = appendLi + '<li>';
                appendLi = appendLi + '<div class="server server_add" onClick="displayNewServerForm()">';
                appendLi = appendLi + '</div>';
                appendLi = appendLi + '</li>';

                $("#mcis_server_list").append(appendLi);

                // commonAlert("VM 목록 조회 완료")
                $("#serverRegistResult").text("VM 목록 조회 완료");
            }
        }
    }).catch((error) => {
        // console.warn(error);
        console.log(error.response)
        var errorMessage = error.response.data.error;
    })
}
