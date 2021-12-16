


$(document).ready(function () {

	// const getData = async() => {
	// 	const {data:result} = await axios.get("server adress");
	// 	return result;
	// }
	console.log("start!! " + Date.now());
	// namespace  caller, isCallback, targetObjId, searchOption
	getCommonNameSpaceList("mainnamespace", true, '', "id")
	console.log("getCommonNameSpaceList!! " + Date.now());
	// connection
	// credential
	getCommonCredentialList("maincredential", "id");
	console.log("getCommonCredentialList!! " + Date.now());
	// region
	getCommonRegionList("mainregion", "id");
	console.log("getCommonRegionList!! " + Date.now());
	// driver
	getCommonDriverList("maindriver", "id");
	console.log("getCommonDriverList!! " + Date.now());
	// resource
	// network(vnet)
	getCommonNetworkList("mainnetwork", "id")
	console.log("getCommonNetworkList!! " + Date.now());
	// securitygroup
	getCommonSecurityGroupList("mainsecuritygroup", "", "id")
	console.log("getCommonSecurityGroupList!! " + Date.now());
	// sshkey
	getCommonSshKeyList("mainsshkey", "id")
	console.log("getCommonSshKeyList!! " + Date.now());
	
	// image
	getCommonVirtualMachineImageList("mainimage", "", "id")

	// spec
	getCommonVirtualMachineSpecList("mainspec", "", "id")
	
	getCommonMcisList("mainmcis", true, "", "id")

	getCommonMcksList("mainmcks", "id")
	//$("#guideArea").modal();
});                   

let guideMap = new Map();
// TODO : 가져온 결과로 어떻게 처리할 것인지
function processMap(caller){
	console.log("GUIDE---------- " + caller)
	console.log(guideMap)
	try{
	var keyValue = guideMap.get(caller);
	if( keyValue > 0){
		$("#goto" + caller).html("")
		$("#goto" + caller).html("생성완료")
	}else{
		console.log("-- goto" + caller)
		document.getElementById("goto" + caller).style.display = "";
		if( caller.indexOf("credential") > -1
			|| caller.indexOf("driver") > -1
			|| caller.indexOf("region") > -1
			|| caller.indexOf("mcis") > -1 ){
			console.log("guide area modal ")
			$("#guideArea").modal();
		}
	}
}catch(e){console.log(e)}
	console.log("goto" + caller)
	// guideMap.forEach( (value, key, map) => {
	// 	// alert(`${key}: ${value}`); // cucumber: 500 ...
	// 	console.log(key + " : " + value);
	// 	if( key == "namespace" && value == 0 ){
	// 		$("#guideArea").modal();
	// 		return;
	// 	}

	// 	if( key == "credential" && value == 0 ){
	// 		document.getElementById("gotoCredential").style.display = "";
	// 		$("#guideArea").modal();
	// 		return;
	// 	}else{
	// 		$("#gotoCredential").html("생성완료")
	// 		document.getElementById("gotoCredential").style.display = "";
	// 		return;
	// 	}

	// 	if( key == "region" && value == 0 ){
	// 		document.getElementById("gotoRegion").style.display = "";
	// 		$("#guideArea").modal();
	// 		return;
	// 	}else{
	// 		$("#gotoRegion").text("생성완료")
	// 		document.getElementById("gotoRegion").style.display = "";
	// 		return;
	// 	}

	// 	if( key == "driver" && value == 0 ){
	// 		document.getElementById("gotoDriver").style.display = "";
	// 		$("#guideArea").modal();
	// 		return;
	// 	}else{
	// 		$("#gotoDriver").text("생성완료")
	// 		document.getElementById("gotoDriver").style.display = "";
	// 		return;
	// 	}

	// 	if( key == "network" && value == 0 ){
	// 		document.getElementById("gotoNetwork").style.display = "";
	// 		$("#guideArea").modal();
	// 		return;
	// 	}else{
	// 		$("#gotoNetwork").text("생성완료")
	// 		document.getElementById("gotoNetwork").style.display = "";
	// 		return;
	// 	}

	// 	if( key == "securitygroup" && value == 0 ){
	// 		document.getElementById("gotoSecurity").style.display = "";
	// 		$("#guideArea").modal();
	// 		return;
	// 	}else{
	// 		$("#gotoSecurity").text("생성완료")
	// 		document.getElementById("gotoSecurity").style.display = "";
	// 		return;
	// 	}

	// 	if( key == "sshkey" && value == 0 ){
	// 		document.getElementById("gotoSshKey").style.display = "";
	// 		$("#guideArea").modal();
	// 		return;
	// 	}else{
	// 		$("#gotoSshKey").text("생성완료")
	// 		document.getElementById("gotoSshKey").style.display = "";
	// 		return;
	// 	}

	// 	if( key == "image" && value == 0 ){
	// 		document.getElementById("gotoServerImage").style.display = "";
	// 		$("#guideArea").modal();
	// 		return;
	// 	}else{
	// 		$("#gotoServerImage").text("생성완료")
	// 		document.getElementById("gotoServerImage").style.display = "";
	// 		return;
	// 	}

	// 	if( key == "spec" && value == 0 ){
	// 		document.getElementById("gotoServerSpec").style.display = "";
	// 		$("#guideArea").modal();
	// 		return;
	// 	}else{
	// 		$("#gotoServerSpec").text("생성완료")
	// 		document.getElementById("gotoServerSpec").style.display = "";
	// 		return;
	// 	}
	// });
}

function getNameSpaceListCallbackSuccess(caller, data){
	console.log(data);
	if ( data == null || data == undefined || data == "null"){
		guideMap.set(caller, 0)
	}else{// 아직 data가 1건도 없을 수 있음
		if( data.length > 0){
			guideMap.set(caller, 1)
		}
	}
	// console.log(data);
	processMap(caller);
}
function getNameSpaceListCallbackFail(caller, error){
	guideMap.set(caller, 0)
	processMap(caller);
}

function getCredentialListCallbackSuccess(caller, data){
	console.log(data);
	if ( data == null || data == undefined || data == "null"){
		guideMap.set(caller, 0)
	}else{// 아직 data가 1건도 없을 수 있음
		if( data.length > 0){
			guideMap.set(caller, 1)
		}
	}
	// console.log(data);
	processMap(caller);
}
function getCredentialListCallbackFail(caller, error){
	guideMap.set(caller, 0)
	processMap(caller);
}

function getRegionListCallbackSuccess(caller, data){
	console.log(data);
	if ( data == null || data == undefined || data == "null"){
		guideMap.set(caller, 0)
	}else{// 아직 data가 1건도 없을 수 있음
		if( data.length > 0){
			guideMap.set(caller, 1)
		}
	}
	// console.log(data);
	processMap(caller);
}
function getRegionListCallbackFail(caller, error){
	guideMap.set(caller, 0)
	processMap(caller);
}

function getDriverListCallbackSuccess(caller, data){
	console.log(data);
	if ( data == null || data == undefined || data == "null"){
		guideMap.set(caller, 0)		
	}else{// 아직 data가 1건도 없을 수 있음
		if( data.length > 0){
			guideMap.set(caller, 1)
		}
	}
	// console.log(data);
	processMap(caller);
}
function getDriverListCallbackFail(caller, error){
	guideMap.set(caller, 0)
	processMap(caller);
}

function getNetworkListCallbackSuccess(caller, data){
	console.log(data);
	if ( data == null || data == undefined || data == "null"){
		guideMap.set(caller, 0)
	}else{// 아직 data가 1건도 없을 수 있음
		if( data.length > 0){
			guideMap.set(caller, 1)
		}
	}
	// console.log(data);
	processMap(caller);
}
function getNetworkListCallbackFail(caller, error){
	guideMap.set(caller, 0)
	processMap(caller);
}

function getSecurityGroupListCallbackSuccess(caller, data){
	console.log(data);
	if ( data == null || data == undefined || data == "null"){
		guideMap.set(caller, 0)
	}else{// 아직 data가 1건도 없을 수 있음
		if( data.length > 0){
			guideMap.set(caller, 1)
		}
	}
	// console.log(data);
	processMap(caller);
}
function getSecurityGroupListCallbackFail(caller, error){
	guideMap.set(caller, 0)
	processMap(caller);
}

function getSshKeyListCallbackSuccess(caller, data){
	console.log(data);
	if ( data == null || data == undefined || data == "null"){
		guideMap.set(caller, 0)
	}else{// 아직 data가 1건도 없을 수 있음
		if( data.length > 0){
			guideMap.set(caller, 1)
		}
	}
	// console.log(data);
	processMap(caller);
}
function getSshKeyListCallbackFail(caller, error){
	guideMap.set(caller, 0)
	processMap(caller);
}

//getCommonVirtualMachineImageList
function getImageListCallbackSuccess(caller, data){
	console.log(data);
	if ( data == null || data == undefined || data == "null"){
		guideMap.set(caller, 0)
	}else{// 아직 data가 1건도 없을 수 있음
		if( data.length > 0){
			guideMap.set(caller, 1)
		}
	}
	// console.log(data);
	processMap(caller);
}

function getImageListCallbackFail(caller, error){	
	guideMap.set(caller, 0)
	processMap(caller);
}

function getSpecListCallbackSuccess(caller, data){
	console.log(data);
	if ( data == null || data == undefined || data == "null"){
		guideMap.set(caller, 0)
	}else{// 아직 data가 1건도 없을 수 있음
		if( data.length > 0){
			guideMap.set(caller, 1)
		}
	}
	// console.log(data);
	processMap(caller);
}

function getSpecListCallbackFail(caller, error){	
	guideMap.set(caller, 0)
	processMap(caller);
}

function getMcisListCallbackSuccess(caller, data){
	console.log(data);
	if ( data == null || data == undefined || data == "null"){
		guideMap.set(caller, 0)
	}else{// 아직 data가 1건도 없을 수 있음
		if( data.length > 0){
			guideMap.set(caller, 1)
		}
	}
	// console.log(data);
	processMap(caller);
}

function getMcisListCallbackFail(caller, error){
	guideMap.set(caller, 0)
	processMap(caller);
}


function getMcksListCallbackSuccess(caller, data){
	console.log("getMcksListCallbackSuccess--" + caller);
	console.log(data);
	if ( data == null || data == undefined || data == "null"){
		guideMap.set(caller, 0)
	}else{// 아직 data가 1건도 없을 수 있음
		if( data.length > 0){
			guideMap.set(caller, 1)
		}
	}
	// console.log(data);
	processMap(caller);
}

function getMcksListCallbackFail(caller, error){	
	guideMap.set(caller, 0)
	processMap(caller);
}