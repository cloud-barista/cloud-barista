
$(document).ready(function () {
	//btn_spec
	// #ID 에 .클래스명_assist
	//	대상 class명.toggleClass
	$('#OS_HW_Spec_Assist .btn_spec_assist').click(function () {
		$(".spec_select_box").toggleClass("active");

	});

	$('#OS_HW_Spec .btn_image_assist').click(function () {
		$(".spec_select_box").toggleClass("active");
	});


	// 방안 1. shown 이후 sleep 3초
	// 방안 2. z-index 변경
	// 방안 3. 지도 클릭이 필수가 아니면 priority option 선택 시 지도 div를 show. 기본은 hide
	$("#recommendVmAssist").on("shown.bs.modal", function (e) {
		console.log("shown.bs.modal")
		console.log(e)
		showMap()
	});

	$("#recommendVmAssist").on("show.bs.modal", function (e) {
		console.log("show.bs.modal")
		console.log(e)
	});
});

function sleep(ms) {
	const wakeUpTime = Date.now() + ms;
	while (Date.now() < wakeUpTime) { }
}

var JZMap;
function showMap() {
	// //
	//
	// var locationInfo = new Object();
	// locationInfo.id = "1"
	// locationInfo.name = "pin"
	// locationInfo.cloudType = "aws";
	// locationInfo.latitude = "34.3800";
	// locationInfo.longitude = "131.7000"
	// locationInfo.markerIndex = 1
	// setMap(locationInfo)

	$("#recommend_map").empty();
	sleep(2000)
	JZMap = map_init_target("recommend_map")
	addClickPin(JZMap)
}

// Map 관련 설정
function setMap(locationInfo) {
	//show_mcis2(url,JZMap);
	//function show_mcis2(url, map){
	// var JZMap = map;

	if (locationInfo == undefined) {
		var locationInfo = new Object();
		locationInfo.id = "1"
		locationInfo.name = "pin"
		locationInfo.cloudType = "aws";
		locationInfo.latitude = "34.3800";
		locationInfo.longitude = "131.7000"
		locationInfo.markerIndex = 1
	}

	console.log("setMap")
	console.log(locationInfo)
	$("#map").empty();

	var JZMap = map_init()// mcis.map.js 파일에 정의되어 있으므로 import 필요.  TODO : map click할 때 feature 에 id가 없어 tooltip 에러나고 있음. 해결필요

	//지도 그리기 관련
	var polyArr = new Array();

	var longitudeValue = locationInfo.longitude;
	var latitudeValue = locationInfo.latitude;
	console.log(longitudeValue + " : " + latitudeValue);
	if (longitudeValue && latitudeValue) {
		console.log("drawMap before")
		drawMap(JZMap, longitudeValue, latitudeValue, locationInfo)
		console.log("drawMap after")
	}
}

// // var JZMap;
// var locationMap = new Object();
// function setMap(locationInfo) {
// 	//show_mcis2(url,JZMap);
// 	//function show_mcis2(url, map){
// 	// var JZMap = map;
// 	console.log(recommendVmAssist)
// 	if (locationInfo == undefined) {
// 		var locationInfo = new Object();
// 		locationInfo.id = "1"
// 		locationInfo.name = "pin"
// 		locationInfo.cloudType = "aws";
// 		locationInfo.latitude = "34.3800";
// 		locationInfo.longitude = "131.7000"
// 		locationInfo.markerIndex = 1
// 	}
// 	alert(1)
// 	console.log("setMap")
// 	//지도 그리기 관련
// 	var polyArr = new Array();
//
// 	var longitudeValue = locationInfo.longitude;
// 	var latitudeValue = locationInfo.latitude;
// 	console.log(longitudeValue + " : " + latitudeValue);
// 	alert(2)
// 	if (longitudeValue && latitudeValue) {
// 		$("#map").empty()
// 		var locationMap = map_init()// mcis.map.js 파일에 정의되어 있으므로 import 필요.  TODO : map click할 때 feature 에 id가 없어 tooltip 에러나고 있음. 해결필요
// 		drawMap(locationMap, longitudeValue, latitudeValue, locationInfo)
// 		alert(3)
// 	}
// 	alert(4)
// 	console.log(recommendVmAssist)
// }


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
function getRegionListFilterAtAssist(provider, targetRegionObj) {
	// region 목록 filter
	selectBoxFilterByText(targetRegionObj, provider)
	$("#" + targetRegionObj + " option:eq(0)").attr("selected", "selected");
}

// assist popup에서 조회조건에 맞는 spec을 검색
function assistFilterSpec() {
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
	for (var i = 0; i < conditionArr.length; i++) {
		var conditionMaxValue = $("#assist_" + conditionArr[i] + "_max").val();
		var conditionMinValue = $("#assist_" + conditionArr[i] + "_min").val();
		console.log("conditionMinValue=" + conditionMinValue);
		console.log("conditionMaxValue=" + conditionMaxValue);
		if (conditionMaxValue && conditionMinValue) {
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
function filterSpecsByRangeCallbackSuccess(caller, data) {
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
	if (data.length) {
		vmSpecList.forEach(function (item, index) {
			html += '<tr onclick="setAssistSpecId(\'' + item.id + '\', \'' + item.name + '\', \'' + item.cspSpecName + '\', \'' + item.connectionName + '\')">'
				+ '<td class="btn_mtd" data-th="spec ID">' + item.id + '<span class="ov off"></span></td>'
				+ '<td class="overlay hidden" data-th="spec Name">' + item.name + '</td>'
				+ '<td class="overlay hidden" data-th="csp spec Name">' + item.cspSpecName + '</td>'
				+ '<td class="overlay hidden" data-th="connection name">' + item.connectionName + '</td>'
				+ '<td class="overlay hidden" data-th="os type">' + item.os_type + '</td>'
				+ '<td class="overlay hidden" data-th="Cpu / core / mem / disk">CPU : ' + item.num_vCPU + '<br>Core : ' + item.num_core + '<br>Disk : ' + item.storage_GiB + '</td>'
				+ '<td class="overlay hidden" data-th="description">' + item.description + '</td>'
				+ '</tr>'
		})
		$("#assist_specList").empty()
		$("#assist_specList").append(html)
	} else {
		commonAlert("No result Found")
	}



}
// Spec Range 조회 실패
function filterSpecsByRangeCallbackFail() {
	commonAlert("Failt to Search Specs")
}

// table에서 spec 선택시 hidden으로 set
function setAssistSpecId(speID, specName, cspSpecName, connectionName) {
	console.log(speID + ":" + specName + ":" + cspSpecName + ":" + connectionName)
	$("#assist_vmSpec_id").val(speID);
	$("#assist_vmSpec_specName").val(specName);
	$("#assist_vmSpec_cspSpecName").val(cspSpecName);
	$("#assist_vmSpec_connectionName").val(connectionName);
	$("#assist_vmSpec_info").val(speID + "|" + specName + "|" + connectionName + "|" + cspSpecName);

}

// apply버튼 클릭시
function applyAssistSpec() {
	var selectedSpecID = $("#assist_vmSpec_id").val();
	if (selectedSpecID) {
		//<tr onclick="setValueToFormObj('es_imageList', 'tab_vmImage', 'vmImage', '{{$vmInageIndex}}', 'e_imageId');">
		// $("#tab_vmSpecInfo")
		var selectedConnectionName = $("#assist_vmSpec_connectionName").val();
		var selectedCspSpecName = $("#assist_vmSpec_cspSpecName").val();
		var selectedSpecInfo = $("#assist_vmSpec_info").val();
		console.log(selectedSpecInfo);
		$("#tab_vmSpecInfo").val(selectedSpecInfo);
		$("#tab_vmSpec_cspSpecName").val(selectedCspSpecName);
		$("#tab_vmSpecConnectionName").val(selectedConnectionName);
		$("#e_specId").val(selectedSpecID);

		var esSelectedConnectionName = $("#es_regConnectionName option:selected").val()
		if (esSelectedConnectionName == "") {// 선택한 connectionName이 없으면 set
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


function setSpecValueToFormObj(selectedId, selectedSpecName, cspSpecName, selectedConnectionName) {
	var econnectionName = $("#e_connectionName").val();
	if (econnectionName != "" && econnectionName != selectedConnectionName) {
		$("#t_connectionName").val(selectedConnectionName);// confirm을 통해서 form에 set 되므로 임시(t_connectionName)로 저장.
		commonConfirmOpen("DifferentConnection");
	} else {
		var esSelectedConnectionName = $("#es_regConnectionName option:selected").val()
		if (esSelectedConnectionName == "") {// 선택한 connectionName이 없으면 set
			$("#es_regConnectionName").val(selectedConnectionName);
		}

		$("#e_connectionName").val(selectedConnectionName);
		$("#e_imageId" + targetObjId).val(selectedId);

		//<input type="hidden" name="vmImage_info" id="vmImage_info_{{$vmInageIndex}}" value="{{$vmImageItem.ID}}|{{$vmImageItem.Name}}|{{$vmImageItem.ConnectionName}}|{{$vmImageItem.CspImageId}}|{{$vmImageItem.CspImageName}}|{{$vmImageItem.GuestOS}}|{{$vmImageItem.Description}}"/>
		$("#tab_vmImageInfo").val(selectedId + "|" + selectedSpecName + "|" + selectedConnectionName + "|" + cspSpecName);
	}
}

// EnterKey입력 시 해당 값, keyword 들이 있는 object id, 구분자(caller)
function searchAssistNetworkByEnter(event, caller) {
	if (event.keyCode === 13) {
		searchNetworkByKeyword(caller);
	}
}

//
function searchNetworkByKeyword(caller) {
	var keyword = "";
	var keywordObjId = "";
	if (caller == "searchNetworkAssistAtReg") {
		keyword = $("#keywordAssistNetwork").val();
		keywordObjId = "searchAssistNetworkKeywords";
		// network api에 connection으로 filter하는 기능이 없으므로
		//totalNetworkListByNamespace : page Load시 가져온 network List가 있으므로 해당 목록을 Filter한다.
	}

	// connection

	//
	if (!keyword) {
		commonAlert("At least a keyword required");
		return;
	}
	var addKeyword = '<div class="keyword" name="keyword_' + caller + '">' + keyword.trim() + '<button class="btn_del_image" onclick="delSearchKeyword(event, \'' + caller + '\')"></button></div>';

	$("#" + keywordObjId).append(addKeyword);
	var keywords = new Array();// 기존에 있는 keyword에 받은 keyword 추가하여 filter적용
	$("[name='keyword_" + caller + "']").each(function (idx, ele) {
		keywords.push($(this).text());
	});

	//getCommonSearchVmImageList(keywords, caller);
	filterNetworkList(keywords, caller)
}

// EnterKey입력 시 해당 값, keyword 들이 있는 object id, 구분자(caller)
function searchAssistSecurityGroupByEnter(event, caller) {
	if (event.keyCode === 13) {
		searchSecurityGroupByKeyword(caller);
	}
}

//
function searchSecurityGroupByKeyword(caller) {
	var keyword = "";
	var keywordObjId = "";
	console.log(caller)
	if (caller == "searchSecurityGroupAssistAtReg") {
		keyword = $("#keywordAssistSecurityGroup").val();
		keywordObjId = "searchAssistNetworkKeywords";
		// securityGroup api에 connection으로 filter하는 기능이 없으므로
		//totalSecurityGroupListByNamespace : page Load시 가져온 securityGroup List가 있으므로 해당 목록을 Filter한다.
	}

	// connection

	//
	if (!keyword) {
		commonAlert("At least a keyword required");
		return;
	}
	var addKeyword = '<div class="keyword" name="keyword_' + caller + '">' + keyword.trim() + '<button class="btn_del_image" onclick="delSearchKeyword(event, \'' + caller + '\')"></button></div>';

	$("#" + keywordObjId).append(addKeyword);
	var keywords = new Array();// 기존에 있는 keyword에 받은 keyword 추가하여 filter적용
	$("[name='keyword_" + caller + "']").each(function (idx, ele) {
		keywords.push($(this).text());
	});

	filterSecurityGroupList(keywords, caller)
}

// EnterKey입력 시 해당 값, keyword 들이 있는 object id, 구분자(caller)
function searchAssistSshKeyByEnter(event, caller) {
	if (event.keyCode === 13) {
		searchSshKeyByKeyword(caller);
	}
}

//
function searchSshKeyByKeyword(caller) {
	var keyword = "";
	var keywordObjId = "";
	if (caller == "searchSshKeyAssistAtReg") {
		keyword = $("#keywordAssistSshKey").val();
		keywordObjId = "searchAssistSshKeyKeywords";
		// network api에 connection으로 filter하는 기능이 없으므로
		//totalSshKeyListByNamespace : page Load시 가져온 sshKey List가 있으므로 해당 목록을 Filter한다.
	}

	// connection

	//
	if (!keyword) {
		commonAlert("At least a keyword required");
		return;
	}
	var addKeyword = '<div class="keyword" name="keyword_' + caller + '">' + keyword.trim() + '<button class="btn_del_image" onclick="delSearchKeyword(event, \'' + caller + '\')"></button></div>';

	$("#" + keywordObjId).append(addKeyword);
	var keywords = new Array();// 기존에 있는 keyword에 받은 keyword 추가하여 filter적용
	$("[name='keyword_" + caller + "']").each(function (idx, ele) {
		keywords.push($(this).text());
	});

	//getCommonSearchVmImageList(keywords, caller);
	filterSshKeyList(keywords, caller)
}

// EnterKey입력 시 해당 값, keyword 들이 있는 object id, 구분자(caller)
function searchAssistImageByEnter(event, caller) {
	if (event.keyCode === 13) {
		// searchKeyword(keyword, caller);
		searchVmImageByKeyword(caller);
		// searchKeyword($(this).val(), caller)
	}
}

//
function searchVmImageByKeyword(caller) {
	var keyword = "";
	var keywordObjId = "";
	if (caller == "searchVmImageAssistAtReg") {
		keyword = $("#keywordAssistImage").val();
		keywordObjId = "searchAssistImageKeywords";
	}

	// if (!keyword) {
	// 	commonAlert("At least a keyword required");
	// 	return;
	// }
	if (keyword != "") {
		var addKeyword = '<div class="keyword" name="keyword_' + caller + '">' + keyword.trim() + '<button class="btn_del_image" onclick="delSearchKeyword(event, \'' + caller + '\')"></button></div>';
	}

	$("#" + keywordObjId).append(addKeyword);
	var keywords = new Array();// 기존에 있는 keyword에 받은 keyword 추가하여 filter적용
	$("[name='keyword_" + caller + "']").each(function (idx, ele) {
		keywords.push($(this).text());
	});

	getCommonSearchVmImageList(keywords, caller);
}

// Assist Spec filter Search버튼 클릭시
function searchSpecsByRange(caller) {
	// var specFilter = new Object();

	var assistSpecConnectionNameVal = $("#assistSpecConnectionName option:selected").val();
	if (caller == 'searchVmSpecAssistAtReg') {

	}
	// if (assistSpecConnectionNameVal) {
	//     specFilter.connectionName = assistSpecConnectionNameVal
	// }

	// storage
	// var storageMin = $("#assist_num_storage_min").val();
	// var storageMax = $("#assist_num_storage_max").val();
	// var storageObj = new Object();
	// storageObj.min = Number(storageMin)
	// storageObj.max = Number(storageMax)

	// Core
	// var coreMin = $("#assist_num_core_min").val();
	// var coreMax = $("#assist_num_core_max").val();
	// var coreObj = new Object();
	// coreObj.min = Number(coreMin)
	// coreObj.max = Number(coreMax)

	// specFilter.numCore = { "min": coreMin, "max": coreMax };

	// vCPU
	var vCpuMin = $("#assist_num_vCPU_min").val();
	var vCpuMax = $("#assist_num_vCPU_max").val();
	var vCpuObj = new Object();
	vCpuObj.min = Number(vCpuMin)
	vCpuObj.max = Number(vCpuMax)
	// specFilter.numvCPU = { "min": vCpuMin, "max": vCpuMax };

	// memory
	var memGiBMin = $("#assist_num_memory_min").val();
	var memGiBMax = $("#assist_num_memory_max").val();
	var memGiBObj = new Object();
	memGiBObj.min = Number(vCpuMin)
	memGiBObj.max = Number(vCpuMax)

	var specFilter = {
		connectionName: assistSpecConnectionNameVal,
		// maxTotalStorageTiB: storageObj,
		// numCore: coreObj,
		numvCPU: vCpuObj,
		memGib: memGiBObj,
	}
	getCommonFilterVmSpecListByRange(specFilter, caller)

	// ID             string `json:"id"`
	// Name           string `json:"name"`
	// Description    string `json:"description"`
	// ConnectionName string `json:"connectionName"`
	// CspSpecName    string `json:"cspSpecName"`
	// OsType         string `json:"osType"`
	//
	// CostPerHour Range `json:"costPerHour"`
	// EbsBwMbps   Range `json:"ebsBwMbps"`
	//
	// EvaluationScore01 Range  `json:"evaluationScore01"`
	// EvaluationScore02 Range  `json:"evaluationScore02"`
	// EvaluationScore03 Range  `json:"evaluationScore03"`
	// EvaluationScore04 Range  `json:"evaluationScore04"`
	// EvaluationScore05 Range  `json:"evaluationScore05"`
	// EvaluationScore06 Range  `json:"evaluationScore06"`
	// EvaluationScore07 Range  `json:"evaluationScore07"`
	// EvaluationScore08 Range  `json:"evaluationScore08"`
	// EvaluationScore09 Range  `json:"evaluationScore09"`
	// EvaluationScore10 Range  `json:"evaluationScore10"`
	// EvaluationStatus  string `json:"evaluationStatus"`
	//
	// GpuModel string `json:"gpuModel"`
	// GpuP2p   string `json:"gpuP2p"`
	//
	// MaxNumStorage      Range `json:"maxNumStorage"`
	// MaxTotalStorageTiB Range `json:"maxTotalStorageTiB"`
	// MemGiB             Range `json:"memGiB"`
	//
	// NetBwGbps  Range `json:"netBwGbps"`
	// NumCore    Range `json:"numCore"`
	// NumGpu     Range `json:"numGpu"`
	// NumStorage Range `json:"numStorage"`
	// NumVCPU    Range `json:"numvCPU"`
	// StorageGiB Range `json:"storageGiB"`
}

function getRecommendVmInfo() {
	var max_cpu = $("#num_vCPU_max").val()
	var min_cpu = $("#num_vCPU_min").val()
	var max_mem = $("#num_memory_max").val()
	var min_mem = $("#num_memory_min").val()
	var max_cost = $("#num_cost_max").val()
	var min_cost = $("#num_cost_min").val()
	var limit = $("#recommendVmLimit").val()
	var lon = $("#longitude").val()
	var lat = $("#latitude").val()

	console.log(" lon " + lon + ", lat " + lat)
	if (lon == "" || lat == "") {
		commonAlert(" 지도에서 위치를 선택하세요 ")
		return;
	}
	var url = "/operation/manages/mcismng/mcisrecommendvm/list"
	var obj = {
		"filter": {
			"policy": [
				{
					"condition": [
						{
							"operand": max_cpu,
							"operator": "<="
						},
						{
							"operand": min_cpu,
							"operator": ">="
						}
					],
					"metric": "cpu"
				},
				{
					"condition": [
						{
							"operand": max_mem,
							"operator": "<="
						},
						{
							"operand": min_mem,
							"operator": ">="
						}
					],
					"metric": "memory"
				},
				{
					"condition": [
						{
							"operand": max_cost,
							"operator": "<="
						},
						{
							"operand": min_cost,
							"operator": ">="
						}
					],
					"metric": "cost"
				}
			]
		},
		"limit": limit,
		"priority": {
			"policy": [
				{
					"metric": "location",
					"parameter": [
						{
							"key": "coordinateClose",
							"val": [
								lon + "/" + lat
							]
						}
					],
					"weight": "0.3"
				}
			]
		}
	}
	axios.post(url, obj, {
		headers: {
			'Content-type': 'application/json',
		}
	}).then(result => {
		console.log("result spec : ", result);
		var statusCode = result.data.status;
		if (statusCode == 200 || statusCode == 201) {

			console.log("recommend vm result: ", result.data);
			recommendVmSpecListCallbackSuccess(result.data.VmSpecList)

		} else {
			var message = result.data.message;
			commonAlert("Fail Create Spec : " + message + "(" + statusCode + ")");

		}

	}).catch((error) => {
		console.warn(error);
		console.log(error.response)
		var errorMessage = error.response.data.error;
		var statusCode = error.response.status;
		commonErrorAlert(statusCode, errorMessage);
	});
}

function recommendVmSpecListCallbackSuccess(data) {
	var html = ""
	if (data == null || data.length == 0) {
		html += '<tr><td class="overlay hidden" data-th="" colspan="5">No Data</td></tr>'

		$("#assistRecommendSpecList").empty()
		$("#assistRecommendSpecList").append(html)
	} else {
		if (data.length) {

			data.map((item, index) => (
				html += '<tr onclick="getConnectionConfigCandidateInfo(' + index + ');">'
				+ '     <input type="hidden" id="recommendVmAssist_id_' + index + '" value="' + item.id + '"/>'
				+ '     <input type="hidden" id="recommendVmAssist_provider_' + index + '" value="' + item.providerName + '"/>'
				+ '     <input type="hidden" id="recommendVmAssist_connectionName_' + index + '" value="' + item.connectionName + '"/>'
				+ '     <input type="hidden" id="recommendVmAssist_name_' + index + '" value="' + item.name + '"/>'
				+ '     <input type="hidden" id="recommendVmAssist_cspSpec_' + index + '" value="' + item.cspSpecName + '"/>'
				+ '<td class="overlay hidden" data-th="provider">' + item.providerName + '</td>'
				+ '<td class="overlay hidden" data-th="region">' + item.regionName + '</td>'
				+ '<td class="btn_mtd ovm" data-th="name ">' + item.name + '<span class="ov"></span></td>'
				+ '<td class="overlay hidden" data-th="cspSpec">' + item.cspSpecName + '</td>'
				+ '<td class="overlay hidden" data-th="price">' + item.costPerHour + '</td>'
				+ '<td class="btn_mtd ovm" data-th="mem ">' + item.memGiB + '<span class="ov"></span></td>'
				+ '<td class="overlay hidden" data-th="vcpu">' + item.numvCPU + '</td>'
				+ '<td class="overlay hidden" data-th="evaluationScore01">' + item.evaluationScore01 + '</td>'
				+ '</tr>'
			))


			$("#assistRecommendSpecList").empty()
			$("#assistRecommendSpecList").append(html)
			console.log("setRecommendVmSpec completed");
		}
	}
}

// mcisDynamicCheckRequest -> 해당 spec에 대해 가능한 connection 구하기
function getConnectionConfigCandidateInfo(index) {
	$("#assistSelectedIndex").val(index);
	var specName = $("#recommendVmAssist_name_" + index).val()
	var cspSpecName = $("#recommendVmAssist_cspSpec_" + index).val()
	console.log(specName);
	//var specName = "aws-ap-northeast-1-t2-micro"
	console.log(specName);
	var url = "/operation/manages/mcismng/mcisdynamiccheck/list"
	var obj = {
		"commonSpec": [specName]
	}
	axios.post(url, obj, {
		headers: {
			'Content-type': 'application/json',
		}
	}).then(result => {
		console.log("result connection : ", result);
		var statusCode = result.data.status;
		if (statusCode == 200 || statusCode == 201) {

			console.log("connection result: ", result.data);
			var connectionInfo = result.data.mcisDynamicInfo.reqCheck[0]
			var connectionCandidates = connectionInfo.connectionConfigCandidates
			//if (connectionCandidates.length > 1) {
			selectConnectionConfig(connectionCandidates, cspSpecName, connectionInfo.region.providerName)
			//}
		} else {
			var message = result.data.message;
			commonAlert("Get Connection List Failed : " + message + "(" + statusCode + ")");

		}

	}).catch((error) => {
		console.warn(error);
		console.log(error.response)
		var errorMessage = error.response.data.error;
		var statusCode = error.response.status;
		commonErrorAlert(statusCode, errorMessage);
	});
}

// connection 후보 보여주기
// 가져온 connection 목록과 일치하는 spec 정보 보여주기
// page Load 시 이미 해당 namespace의 전체 목록을 가져 옴.
function selectConnectionConfig(connections, selectedCspSpecName, selectedProvider) {
	// assistConnectionList
	console.log("selected csp spec :", selectedCspSpecName)
	console.log("vmSpecList: ", totalVmSpecListByNamespace)
	var displayItemsCount = 0
	var html = ""
	connections.forEach(function (candidateConnectionName, connectionIndex) {
		var specExist = false
		totalVmSpecListByNamespace.forEach(function (vSpecItem, vSpecIndex) {
			if (candidateConnectionName == vSpecItem.connectionName) {
				displayItemsCount++
				if (selectedCspSpecName == vSpecItem.cspSpecName) {
					spec = vSpecItem.name
					html += '<tr id="connectionAssist_tr_' + vSpecIndex + '" onclick="setConnectionAndSpec(' + vSpecIndex + ');">'
						+ '     <input type="hidden" id="connectionAssist_provider_' + vSpecIndex + '" value="' + selectedProvider + '"/>'
						+ '     <input type="hidden" id="connectionAssist_specName_' + vSpecIndex + '" value="' + vSpecItem.name + '"/>'
						+ '     <input type="hidden" id="connectionAssist_cspSpecName_' + vSpecIndex + '" value="' + vSpecItem.cspSpecName + '"/>'
						+ '     <input type="hidden" id="connectionAssist_connection_' + vSpecIndex + '" value="' + vSpecItem.connectionName + '"/>'
						+ '<td class="overlay hidden" data-th="connection">' + vSpecItem.connectionName + '</td>'
						+ '<td class="overlay hidden" data-th="spec">' + spec + ' </td>'
						+ '</tr>'
					specExist = true
				}
			}
		})

		//  spec이 존재하지 않으면 spec 등록 버튼 생성
		if (!specExist) {
			displayItemsCount++
			var specButton = "<button name='' value='' class='btn_apply btn_co btn_cr_g' onclick=registerSpecOnClick('" + candidateConnectionName + "','" + selectedCspSpecName + "','" + selectedProvider + "')><span>spec 등록</span></button>"
			html += '<tr>'
				+ '<td class="overlay hidden" data-th="connection">' + candidateConnectionName + '</td>'
				+ '<td class="overlay hidden" data-th="spec">' + specButton + ' </td>'
				+ '</tr>'
		}


	});
	$("#assistConnectionList").empty()
	$("#assistConnectionList").append(html)

	if (displayItemsCount == 0) {

		commonAlert("해당 spec을 생성할 수 있는 connection이 없습니다.")

	} else {
		showConnectionAssistPopup()
	}
	///// TODO : ApplyButton Click 시
	///// - expert 모드인 경우에는 applyAssistValidCheck(caller) 에서 처리하면 됨
	///// - simple 모드에서는 비슷한 function 추가 필요. : changeConnectionInfo(configName) 으로 새로 가져와서 셋 한 뒤에
	/////   선택한 connection으로 set


	//row.style.display = '';

	// var html = ""
	// if (data.length) {
	// data.map((item, index) => (
	// 	html += '<tr onclick="setConnectionAndSpec(' + index + ');">'
	// 	+ '     <input type="hidden" id="connectionAssist_name_' + index + '" value="' + item + '"/>'
	// 	+ '<td class="overlay hidden" data-th="connection">' + item + '</td>'
	// 	+ '<td class="overlay hidden" data-th="spec"> aws-test-spec-t2-micro </td>'
	// 	+ '</tr>'
	// ))
	// $("#assistConnectionList").empty()
	// $("#assistConnectionList").append(html)
	// console.log("setConnectionList completed");

	//
	// }

	// if (html != "") {
	// 	showConnectionAssistPopup()
	// }
}

function registerSpecOnClick(regCandidateConnection, selectedCspSpecName, regProvider) {
	$("#t_regProvider").val(regProvider)
	$("#t_regRecommendConn").val(regCandidateConnection)
	$("#t_regRecommendCspSpec").val(selectedCspSpecName)
	commonPromptOpen("RegisterRecommendSpec")
}

// 앞서 setting한 connection과 선택한 connection이 같으면 그대로 set
// 다르면 바꿀건지 물어보고 새로운 connection으로 set 
function setConnectionAndSpec(index) {
	selectedProvider = $("#connectionAssist_provider_" + index).val()
	selectedConnection = $("#connectionAssist_connection_" + index).val()
	selectedSpecName = $("#connectionAssist_specName_" + index).val()
	regConnection = $("#ss_regConnectionName").val()
	console.log("regConnection: ", regConnection);
	console.log(selectedProvider);

	$("#t_regProvider").val(selectedProvider)
	$("#t_regConnectionName").val(selectedConnection)
	$("#t_spec").val(selectedSpecName)

	// Connection이 다르면 바꿀건지 물어봄
	console.log("change conn");
	if (regConnection != "" && selectedConnection != regConnection) {
		commonConfirmOpen("ChangeConnection")
	} else {
		changeCloudConnection()
	}

	$("#connectionAssist").modal("hide");
	$("#recommendVmAssist").modal("hide");
}

// selct box option 세팅
function setConnectionsForOptions(connectionList, selctedProvider) {
	var html = ""
	connectionList.forEach(function (connItem, connIndex) {
		if (selctedProvider == connItem.ProviderName) {
			html += '<option value="' + connItem.ConfigName + '">' + connItem.ConfigName + '</option>'
		}
	})
	$("#ss_regConnectionName").empty()
	$("#ss_regConnectionName").append(html)
}

// selct box option 세팅
function setResourcesForOptions(resourceType, resourceList, selectedConnetion) {
	var html = ""
	var resourceObj = ""
	html += '<option value=""> Select ' + resourceType + '</option>'
	resourceList.forEach(function (resourceItem, resourceIndex) {
		if (selectedConnetion == resourceItem.connectionName) {
			if (resourceType == "Spec") {
				html += '<option value="' + resourceItem.id + '">' + resourceItem.name + '(' + resourceItem.cspSpecName + ')</option>'
			} else if (resourceType == "SSH Key") {
				html += '<option value="' + resourceItem.id + '">' + resourceItem.cspSshKeyName + '(' + resourceItem.id + ')</option>'
			} else {
				html += '<option value="' + resourceItem.id + '">' + resourceItem.name + '(' + resourceItem.id + ')</option>'
			}
		}
	})

	if (resourceType == "Spec") {
		resourceObj = "ss_spec"
	} else if (resourceType == "OS Platform") {
		resourceObj = "ss_imageId"
	} else if (resourceType == "SSH Key") {
		resourceObj = "ss_sshKey"
	}

	$("#" + resourceObj).empty()
	$("#" + resourceObj).append(html)

}

// commonConfirmOpen("ChangeConnection")에서 ok했을 때 실행
function changeCloudConnection() {
	var selectedProvider = $("#t_regProvider").val()
	var selectedConnection = $("#t_regConnectionName").val()
	var selectedSpecName = $("#t_spec").val()

	// provider setting
	$("#ss_regProvider").val(selectedProvider)
	console.log("change cloud spec:", totalVmSpecListByNamespace);
	console.log("change cloud image:", totalImageListByNamespace);
	console.log("change cloud sshkey:", totalSshKeyListByNamespace);

	// filtering
	setConnectionsForOptions(totalCloudConnectionList, selectedProvider)
	setResourcesForOptions("Spec", totalVmSpecListByNamespace, selectedConnection)
	setResourcesForOptions("OS Platform", totalImageListByNamespace, selectedConnection)
	setResourcesForOptions("SSH Key", totalSshKeyListByNamespace, selectedConnection)

	// security group, vnet setting
	getSecurityInfo(selectedConnection);
	getVnetInfo(selectedConnection);

	// connection & spec setting
	$("#ss_regConnectionName").val(selectedConnection)
	$("#ss_spec").val(selectedSpecName)

	$("#t_regProvider").val("")
	$("#t_regConnectionName").val("")
	$("#t_spec").val("")
}

// mcis 이름을 입력했는지 확인 
function checkMcisNameExist() {
	var mcisName = $("#mcis_name").val()
	console.log(mcisName);
	if (mcisName) {
		commonConfirmOpen("AddNewMcisDynamic")
	} else {
		commonPromptOpen("AddNewMcisDynamic")
	}
}