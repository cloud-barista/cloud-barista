
$(document).ready(function(){
	
	//btn_spec
	// #ID 에 .클래스명_assist
	//	대상 class명.toggleClass
	$('#OS_HW_Spec_Assist .btn_spec_assist').click(function(){
		$(".spec_select_box").toggleClass("active");
	
	});

	$('#OS_HW_Spec .btn_image_assist').click(function(){
		$(".spec_select_box").toggleClass("active");
	});
});
															
function openTextFile() {
    var input = document.createElement("input");
    input.type = "file";
    input.accept = "text/plain"; // 확장자가 xxx, yyy 일때, ".xxx, .yyy"
    input.onchange = function (event) {
        processFile(event.target.files[0]);
    };
    input.click();
}

// 선택한 파일을 읽어 화면에 보여줌
function processFile(file) {
    var reader = new FileReader();
    reader.onload = function () {
		console.log(reader.result);
        $("#fileContent").val(reader.result);
    };
    //reader.readAsText(file, /* optional */ "euc-kr");
	reader.readAsText(file);
}


// function exportVmScript(vmIndex){
	
// 	var connectionNameVal = $("#p_connectionName_" + vmIndex).val();
// 	var descriptionVal = $("#p_description_" + vmIndex).val();
// 	var imageIdVal = $("#p_imageId_" + vmIndex).val();
// 	var labelVal = $("#p_label_" + vmIndex).val();
// 	var nameVal = $("#p_name_" + vmIndex).val();
// 	var securityGroupIdsVal = $("#p_securityGroupIds_" + vmIndex).val();
// 	var specIdVal = $("#p_specId_" + vmIndex).val();
// 	var sshKeyIdVal = $("#p_sshKeyId_" + vmIndex).val();
// 	var subnetIdVal = $("#p_subnetId_" + vmIndex).val();
// 	var vNetIdVal = $("#p_vNetId_" + vmIndex).val();
// 	var vmGroupSizeVal = $("#p_vmGroupSize_" + vmIndex).val();
// 	var vmUserAccountVal = $("#p_vmUserAccount_" + vmIndex).val();
// 	var vmUserPasswordVal = $("#p_vmUserPassword_" + vmIndex).val();

// 	var paramValueAppend = '"';
// 	var vmCreateScript = "";
// 	vmCreateScript += '{	';
// 	vmCreateScript += paramValueAppend + 'connectionName' + paramValueAppend + ' : ' + paramValueAppend + connectionNameVal + paramValueAppend;
// 	vmCreateScript += ',' + paramValueAppend + 'description' + paramValueAppend + ' : ' + paramValueAppend + descriptionVal + paramValueAppend;
// 	vmCreateScript += ',' + paramValueAppend + 'imageId' + paramValueAppend + ' : ' + paramValueAppend + imageIdVal + paramValueAppend;
// 	vmCreateScript += ',' + paramValueAppend + 'label' + paramValueAppend + ' : ' + paramValueAppend + labelVal + paramValueAppend;
// 	vmCreateScript += ',' + paramValueAppend + 'name' + paramValueAppend + ' : ' + paramValueAppend + nameVal + paramValueAppend;
// 	// vmCreateScript += ',securityGroupIds: ';
//     // vmCreateScript += '	' + paramValueAppend + securityGroupIdsVal + paramValueAppend;
// 	vmCreateScript += ',' + paramValueAppend + 'specId' + paramValueAppend + ' : ' + paramValueAppend + specIdVal + paramValueAppend;
// 	vmCreateScript += ',' + paramValueAppend + 'sshKeyId' + paramValueAppend + ' : ' + paramValueAppend + sshKeyIdVal + paramValueAppend;
// 	vmCreateScript += ',' + paramValueAppend + 'subnetId' + paramValueAppend + ' : ' + paramValueAppend + subnetIdVal + paramValueAppend;
// 	vmCreateScript += ',' + paramValueAppend + 'vNetId' + paramValueAppend + ' : ' + paramValueAppend + vNetIdVal + paramValueAppend;
// 	vmCreateScript += ',' + paramValueAppend + 'vmGroupSize' + paramValueAppend + ' : ' + paramValueAppend + vmGroupSizeVal + paramValueAppend;
// 	vmCreateScript += ',' + paramValueAppend + 'vmUserAccount' + paramValueAppend + ' : ' + paramValueAppend + vmUserAccountVal + paramValueAppend;
// 	vmCreateScript += ',' + paramValueAppend + 'vmUserPassword' + paramValueAppend + ' : ' + paramValueAppend + vmUserPasswordVal + paramValueAppend;
// 	vmCreateScript += '}';

	
// 	$("#exportFileName").val(nameVal);
// 	$("#vmExportScript").val(vmCreateScript);
// }

// function saveVmInfoToFile(){
// 	var fileName = $("#exportFileName").val();
// 	var exportScript = $("#vmExportScript").val();
	
// 	var element = document.createElement('a');
// 	// element.setAttribute('href', 'data:text/plain;charset=utf-8,' + encodeURIComponent(exportScript));
// 	element.setAttribute('href', 'data:text/json;charset=utf-8,' + encodeURIComponent(exportScript));
// 	// element.setAttribute('download', fileName);
// 	element.setAttribute('download', fileName + ".json");

// 	element.style.display = 'none';
// 	document.body.appendChild(element);

// 	element.click();

// 	document.body.removeChild(element);

// }

// assist에서 provider 선택시 retion filter
function getRegionListFilterAtAssist(provider, targetRegionObj){
	// region 목록 filter
	selectBoxFilterByText(targetRegionObj, provider)
	$("#" + targetRegionObj + " option:eq(0)").attr("selected", "selected");
}

// assist popup에서 조회조건에 맞는 spec을 검색
function assistFilterSpec(){
	var conditionArr = new Array();
	conditionArr.push("cost_per_hour");
	conditionArr.push("ebs_bw_Mbps");
	// conditionArr.push("evaluationScore_01");
	// conditionArr.push("evaluationScore_01");
	// conditionArr.push("evaluationScore_01");
	// conditionArr.push("evaluationScore_01");
	// conditionArr.push("evaluationScore_01");
	// conditionArr.push("evaluationScore_01");
	// conditionArr.push("evaluationScore_01");
	// conditionArr.push("evaluationScore_01");
	// conditionArr.push("evaluationScore_01");
	// conditionArr.push("evaluationScore_01");

	// conditionArr.push("gpumem_GiB");
	conditionArr.push("max_num_storage");
	// conditionArr.push("max_total_storage_TiB");
	// conditionArr.push("mem_GiB");
	// conditionArr.push("net_bw_Gbps");
	// conditionArr.push("num_core");
	// conditionArr.push("num_gpu");
	// conditionArr.push("num_storage");
	conditionArr.push("num_vCPU");
	// conditionArr.push("storage_GiB");
	
	// 
	var searchObj = {}
	searchObj['connectionName'] = $("#assist_select_connectionName").val();
	
	// var condition_CostPerHour = {}
	// condition_CostPerHour['max'] = 100
	// condition_CostPerHour['min'] = 10
	// searchObj['cost_per_hour'] = condition_CostPerHour;

	// var condition_ebsBwMbps = {}
	// condition_ebsBwMbps['max'] = Number(ebsBwMbpsMax)
	// condition_ebsBwMbps['min'] = Number(ebsBwMbpsMax)
	// searchObj['ebs_bw_Mbps'] = condition_ebsBwMbps;
	// assist_num_vCPU_min
	for( var i = 0 ; i < conditionArr.length; i++){
		var conditionMaxValue = $("#assist_" + conditionArr[i] + "_max").val();
		var conditionMinValue = $("#assist_" + conditionArr[i] + "_min").val();
		console.log("conditionMinValue=" + conditionMinValue);
		console.log("conditionMaxValue=" + conditionMaxValue);
		if( conditionMaxValue && conditionMinValue){
			var conditionParam = {};
			// conditionParam['max'] = conditionMaxValue;
			// conditionParam['min'] = conditionMinValue;
			conditionParam['max'] = Number(conditionMaxValue);
			conditionParam['min'] = Number(conditionMinValue);
			searchObj[conditionArr[i]] = conditionParam;
		}
	}	
	// console.log(searchObj);
	// axios 전송
	getCommonFilterSpecsByRange("vmassistpopup", searchObj);
	// assist_specList 에 append
}

// Spec Range 조회 성공
function filterSpecsByRangeCallbackSuccess(caller, data){
	console.log(data)
    console.log("caller = " + caller + ", " + data.length)

    var html = ""
    var vmSpecList = data;
    // cost_per_hour
    // ebs_bw_Mbps
    // evaluationScore_01
    // evaluationStatus
    // gpumem_GiB
    // max_num_storage
    // max_total_storage_TiB
    // mem_GiB
    // net_bw_Gbps
    // num_core
    // num_gpu
    // num_storage
    // num_vCPU
    // storage_GiB
    $("#register_box").modal()    
    if(data.length){
        vmSpecList.forEach(function(item, index) {     
            html +='<tr onclick="setAssistSpecId(\''+item.id+'\', \''+item.name+'\', \''+item.cspSpecName+'\', \''+item.connectionName+'\')">'
            +'<td class="btn_mtd" data-th="spec ID">'+item.id+'<span class="ov off"></span></td>'
            +'<td class="overlay hidden" data-th="spec Name">'+item.name+'</td>'
            +'<td class="overlay hidden" data-th="csp spec Name">'+item.cspSpecName+'</td>'
            +'<td class="overlay hidden" data-th="connection name">'+item.connectionName+'</td>'
            +'<td class="overlay hidden" data-th="os type">'+item.os_type+'</td>'
            +'<td class="overlay hidden" data-th="Cpu / core / mem / disk">CPU : '+item.num_vCPU+'<br>Core : ' + item.num_core + '<br>Disk : ' + item.storage_GiB + '</td>'
            +'<td class="overlay hidden" data-th="description">'+item.description+'</td>'
            +'</tr>'
        })
		$("#assist_specList").empty()
    	$("#assist_specList").append(html)
    }else{
		commonAlert("No result Found")
	}

    

}
// Spec Range 조회 실패
function filterSpecsByRangeCallbackFail(){
	commonAlert("Failt to Search Specs")
}

// table에서 spec 선택시 hidden으로 set
function setAssistSpecId(speID, specName, cspSpecName, connectionName){
	console.log(speID + ":" + specName + ":" + cspSpecName + ":" + connectionName)
    $("#assist_vmSpec_id").val(speID);
    $("#assist_vmSpec_specName").val(specName);
    $("#assist_vmSpec_cspSpecName").val(cspSpecName);
    $("#assist_vmSpec_connectionName").val(connectionName);
	$("#assist_vmSpec_info").val(speID + "|" + specName + "|" + connectionName + "|" + cspSpecName);

}

// apply버튼 클릭시
function applyAssistSpec(){
    var selectedSpecID = $("#assist_vmSpec_id").val();	
    if( selectedSpecID){
//<tr onclick="setValueToFormObj('es_imageList', 'tab_vmImage', 'vmImage', '{{$vmInageIndex}}', 'e_imageId');">
		// $("#tab_vmSpecInfo")
		var selectedConnectionName = $("#assist_vmSpec_connectionName").val();
		var selectedSpecInfo = $("#assist_vmSpec_info").val();
		console.log(selectedSpecInfo);
		$("#tab_vmSpecInfo").val(selectedSpecInfo);
		$("#tab_vmSpecConnectionName").val(selectedConnectionName);
		$("#e_specId" ).val(selectedSpecID);
		
		var esSelectedConnectionName = $("#es_regConnectionName option:selected").val()
		if( esSelectedConnectionName == ""){// 선택한 connectionName이 없으면 set
		$("#es_regConnectionName").val(selectedConnectionName);
		}
		$("#e_connectionName").val(selectedConnectionName);
    }

    // 초기화
    $("#assist_select_provider").val('');
	$("#assist_select_resion").val('');
	$("#assist_select_connectionName").val('');

	$("#assist_vmSpec_id").val("");
    $("#assist_vmSpec_specName").val("");
    $("#assist_vmSpec_cspSpecName").val("");
    $("#assist_vmSpec_connectionName").val("");
	$("#assist_vmSpec_info").val("");


    $("#OS_HW_Spec_Assist").modal("hide")    
}


function setSpecValueToFormObj(selectedId, selectedSpecName, cspSpecName, selectedConnectionName){
    var econnectionName = $("#e_connectionName").val();
    if(econnectionName != "" && econnectionName != selectedConnectionName){
      $("#t_connectionName").val(selectedConnectionName);// confirm을 통해서 form에 set 되므로 임시(t_connectionName)로 저장.
      commonConfirmOpen("DifferentConnection");
    }else{
      var esSelectedConnectionName = $("#es_regConnectionName option:selected").val()
      if( esSelectedConnectionName == ""){// 선택한 connectionName이 없으면 set
        $("#es_regConnectionName").val(selectedConnectionName);
      }

      $("#e_connectionName").val(selectedConnectionName);
      $("#e_imageId" + targetObjId).val(selectedId);
      
      //<input type="hidden" name="vmImage_info" id="vmImage_info_{{$vmInageIndex}}" value="{{$vmImageItem.ID}}|{{$vmImageItem.Name}}|{{$vmImageItem.ConnectionName}}|{{$vmImageItem.CspImageId}}|{{$vmImageItem.CspImageName}}|{{$vmImageItem.GuestOS}}|{{$vmImageItem.Description}}"/>
      $("#tab_vmImageInfo").val(selectedId + "|" + selectedSpecName + "|" + selectedConnectionName + "|"  + cspSpecName);      
    }  
  }	
