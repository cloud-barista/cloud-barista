$(document).ready(function(){


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
//         if($Servers.hasClass("active")) {
//         $Servers.toggleClass("active");
//         } else {
//             $Servers.toggleClass("active");
//         }
//     });
//     // Simple add
//     $(".servers_box .switch").change(function() {
//         if ($(".switch .ch").is(":checked")) {	
//                 $('.servers_box .server_add').click(function(){	
                    
//                     $NewServers.addClass("active");
//                     $SimpleServers.removeClass("active");		
//                 });
//         } else {
//             $('.servers_box .server_add').click(function(){
            
//                 $NewServers.removeClass("active");
//                 $SimpleServers.addClass("active");
            
            
//             });		
//         }
//     });
});

// 서버 더하기버튼 클릭시 서버정보 입력area 보이기/숨기기
// isExpert의 체크 여부에 따라 바뀜.
// newServers 와 simpleServers가 있음.
function displayNewServerForm(){
   
    var simpleServerConfig = $("#simpleServerConfig");
    var expertServerConfig = $("#expertServerConfig");    
    var importServerConfig = $("#importServerConfig");

    if ($("#isImport").is(":checked")) {
        simpleServerConfig.removeClass("active");
        expertServerConfig.removeClass("active");
        importServerConfig.addClass("active");
    }else if ($("#isExpert").is(":checked")) {
        simpleServerConfig.removeClass("active");
        expertServerConfig.addClass("active");
        importServerConfig.removeClass("active");        
    }else{
        //simpleServerConfig        
        simpleServerConfig.addClass("active");
        expertServerConfig.removeClass("active");
        importServerConfig.removeClass("active");        
    }
}



// 서버정보 입력 area에서 'DONE'버튼 클릭시 array에 담고 form을 초기화

var TotalServerConfigArr = new Array();
// deploy 버튼 클릭시 등록한 서버목록을 배포.
// function btn_deploy(){
function deployMcis(){
    var mcis_name = $("#mcis_name").val();
    if(!mcis_name){
        commonAlert("Please Input MCIS Name!!!!!")
        return;
    }
    var mcis_desc = $("#mcis_desc").val();
    var placement_algo = $("#placement_algo").val();
    var installMonAgent = $("#installMonAgent").val();

    var new_obj = {}
    
    var vm_len = 0;

    if(IsImport){
        // ImportedMcisScript.name = mcis_name;
        // ImportedMcisScript.description = mcis_desc;
        // ImportedMcisScript.installMonAgent = installMonAgent;
        // console.log(ImportedMcisScript);
        //var theJson = jQuery.parseJSON($(this).val())
        //$("#mcisImportScriptPretty").val(fmt);	
        new_obj = $("#mcisImportScriptPretty").val();
    }else{
//         console.log(Simple_Server_Config_Arr)
        
        // mcis 생성이므로 mcisID가 없음
        new_obj['name'] = mcis_name
        new_obj['description'] = mcis_desc
        new_obj['installMonAgent'] = installMonAgent
    
        if( Simple_Server_Config_Arr ){
            vm_len = Simple_Server_Config_Arr.length;			
            for(var i in Simple_Server_Config_Arr){
                TotalServerConfigArr.push(Simple_Server_Config_Arr[i]);
            }
        }

        if( Expert_Server_Config_Arr ){
            vm_len = Expert_Server_Config_Arr.length;			
            for(var i in Expert_Server_Config_Arr){
                TotalServerConfigArr.push(Expert_Server_Config_Arr[i]);
            }
        }
        
        if(TotalServerConfigArr){
            vm_len = TotalServerConfigArr.length;			
            console.log("Server_Config_Arr length: ",vm_len);
            new_obj['vm'] = TotalServerConfigArr;
            console.log("new obj is : ",new_obj);
        }else{
            commonAlert("Please Input Servers");
            $(".simple_servers_config").addClass("active");
            $("#s_name").focus();
        }
    }

    // var url = CommonURL+"/ns/"+NAMESPACE+"/mcis";
    var url = "/operation/manages/mcismng/reg/proc"
    try{
        axios.post(url,new_obj,{
            headers :{
                'Content-type': 'application/json',
                },
        }).then(result=>{
            console.log("MCIR Register data : ",result);
            console.log("Result Status : ",result.status); 
            if(result.status == 201 || result.status == 200){
                commonAlert("Register Success")                
                var targetUrl = "/operation/manages/mcismng/mngform"
                changePage(targetUrl)
            }else{
                commonAlert("Register Fail")
                //location.reload(true);
            }
        }).catch((error) => {
            // console.warn(error);
            console.log(error.response)
            var errorMessage = error.response.data.error;
            var statusCode = error.response.status;
            commonErrorAlert(statusCode, errorMessage) 
            
        })
    }catch(error){
		commonAlert(error);
		console.log(error);
	}    
}

// MCIS Create 와 VM Create의 function이름이 같음
function displayMcisImportServerFormByImport(importType){
    // var $SimpleServers = $("#simpleServerConfig");
    // var $ExpertServers = $("#expertServerConfig");
    // var $ImportServers = $("#importServerConfig");
    // var check = $(".switch .ch").is(":checked");
    // console.log("check=" + check);
    // if( check){
    //     $SimpleServers.removeClass("active");
    //     $ExpertServers.removeClass("active");            
    //     $ImportServers.addClass("active");

    //     importMCISInfoFromFile();// import창 띄우기 
    // }

    var addMcisByScriptArea = $("#addMcisByScript");
    var addVmListArea = $("#addVmList");
    var mcisInfoboxArea = $("#mcisInfobox");
    
    if( importType){
        addMcisByScriptArea.css("display", "block");
        addVmListArea.css("display", "none");
        mcisInfoboxArea.css("display", "none");
        // addMcisByScriptArea.addClass("active");
        // addVmListArea.removeClass("active");

        importMCISInfoFromFile();// import창 띄우기 
    }else{
        $("#mcisImportScriptPretty").val("");	
        addMcisByScriptArea.css("display", "none");
        addVmListArea.css("display", "block");
        mcisInfoboxArea.css("display", "block");
    }
    IsImport = importType;// 전역으로 set
}

// mcis export한 파일 선택하여 읽기
function importMCISInfoFromFile() {
    var input = document.createElement("input");
    input.type = "file";
    // input.accept = "text/plain"; // 확장자가 xxx, yyy 일때, ".xxx, .yyy"
	input.accept = ".json";
    input.onchange = function (event) {
        importMcisFileProcess(event.target.files[0]);
    };
    input.click();
}

// 선택한 MCIS 파일을 읽어 화면에 보여줌
var ImportedMcisScript = "";
var IsImport = false;
function importMcisFileProcess(file) {
	try{
		var reader = new FileReader();
		reader.onload = function () {
			console.log(reader.result);
			console.log("---1")
			// $("#fileContent").val(reader.result);
			
			var jsonStr = JSON.stringify(reader.result)
			console.log(JSON.stringify(jsonStr));

			ImportedMcisScript= $.parseJSON(reader.result);

			setMcisInfoToForm(ImportedMcisScript);
			mcisTojsonFormatter(ImportedMcisScript)
			
		};
		//reader.readAsText(file, /* optional */ "euc-kr");
		reader.readAsText(file);
	}catch(error){
		commonAlert("File Load Failed");
		console.log(error);
	}
}

// json 객체를 textarea에 표시할 때 예쁘게
function mcisTojsonFormatter(mcisInfoObj){
	var fmt = JSON.stringify(mcisInfoObj, null, 4);    // stringify with 4 spaces at each level
	$("#mcisImportScriptPretty").val(fmt);	
}

// 화면의 mcis 의 값이 변경되면 적용
function setMcisValue(mcisObjId, mcisObjValue){
    if( IsImport ){
        if( mcisObjId == "mcis_name"){
            ImportedMcisScript.name = mcisObjValue;            
        }else if( mcisObjId == "mcis_desc"){
            ImportedMcisScript.description = mcisObjValue;
        }else if( mcisObjId == "installMonAgent"){            
            ImportedMcisScript.installMonAgent = mcisObjValue;
        }
        mcisTojsonFormatter(ImportedMcisScript)
    }
}
// 선택한 파일을 읽어 form에 Set
function setMcisInfoToForm(mcisInfoObj){
    	//export form
	$("#mcis_name").val(mcisInfoObj.name);
	$("#mcis_desc").val(mcisInfoObj.description);
	// $("#label").val(mcisInfoObj.label);
    $("#installMonAgent").val(mcisInfoObj.installMonAgent);
	
	// // 수신한 obj를 바로 deploy로 던질까?
    // var url = "/operation/manages/mcismng/reg/proc"
    // try{        
    //     axios.post(url,mcisInfoObj,{
    //         headers :{
    //             'Content-type': 'application/json',
    //             // 'Authorization': apiInfo,
    //             },
    //     }).then(result=>{
    //         console.log("MCIR Register data : ",result);
    //         console.log("Result Status : ",result.status); 
    //         if(result.status == 201 || result.status == 200){
    //             commonAlert("Register Success")
    //             // location.href = "/Manage/MCIS/list";
    //             // $('#loadingContainer').show();
    //             // location.href = "/operation/manages/mcismng/mngform/"
    //             var targetUrl = "/operation/manages/mcismng/mngform"
    //             changePage(targetUrl)
    //         }else{
    //             commonAlert("Register Fail")
    //             //location.reload(true);
    //         }
    //     }).catch((error) => {
    //         // console.warn(error);
    //         console.log(error.response)
    //         var errorMessage = error.response.data.error;
    //         var statusCode = error.response.status;
    //         commonErrorAlert(statusCode, errorMessage) 
            
    //     })
    // }finally{
    //     // AjaxLoadingShow(false);
    // }  

	
}


// $(document).ready(function() {
//     //OS_HW popup table scrollbar
// //   $('#OS_HW .btn_spec').on('click', function() {
// //         $('#OS_HW_Spec .dtbox.scrollbar-inner').scrollbar();
// //     });
// //     //Security popup table scrollbar
// //   $('#Security .btn_edit').on('click', function() {
// //     $("#security_edit").modal();
// //         $('#security_edit .dtbox.scrollbar-inner').scrollbar();
// //     });
// });
