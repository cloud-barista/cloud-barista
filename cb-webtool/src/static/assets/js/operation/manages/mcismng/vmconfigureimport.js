

// json 객체를 textarea에 표시할 때 예쁘게
function jsonFormatter(vmInfoObj){
	// var fmt = JSON.stringify(vmInfoObj, null, "\t"); // stringify with tabs inserted at each level
	var fmt = JSON.stringify(vmInfoObj, null, 4);    // stringify with 4 spaces at each level
	$("#vmImportScriptPretty").val(fmt);	
}

// 선택한 파일을 읽어 form에 Set
function setVmInfoToForm(vmInfoObj){
	//export form
	$("#i_name").val(vmInfoObj.name);
	$("#i_description").val(vmInfoObj.description);
	$("#i_connectionName").val(vmInfoObj.connectionName);
	$("#i_imageId").val(vmInfoObj.imageId);	
	$("#i_specId").val(vmInfoObj.specId);
	$("#i_subnetId").val(vmInfoObj.subnetId);
	$("#i_vNetId").val(vmInfoObj.vNetId);
	$("#i_securityGroupIds").val(vmInfoObj.securityGroupIds);
	$("#i_sshKeyId").val(vmInfoObj.sshKeyId);
	$("#i_label").val(vmInfoObj.label);

	$("#i_vmUserAccount").val(vmInfoObj.vmUserAccount);
	$("#i_vmUserPassword").val(vmInfoObj.vmUserPassword);

	var addServerCnt = vmInfoObj.vmGroupSize == "" ? 0: vmInfoObj.vmGroupSize;
	$("#i_vm_add_cnt").val(addServerCnt);

	$("#i_vmImportScript").val(JSON.stringify(vmInfoObj));
	
}

			
const Import_Server_Config_Arr = new Array();
var import_data_cnt = 0
const importServerCloneObj = obj=>JSON.parse(JSON.stringify(obj))
function importDone_btn(){
	var import_form = $("#import_form").serializeObject()
	var server_name = import_form.name
	var server_cnt = parseInt(import_form.i_vm_add_cnt)
	console.log('server_cnt : ',server_cnt)
	var add_server_html = "";
	
	if(server_cnt > 1){
		for(var i = 1; i <= server_cnt; i++){
			var new_vm_name = server_name+"-"+i;
			var object = importServerCloneObj(import_form)
			object.name = new_vm_name
			
			add_server_html +='<li onclick="view_import(\''+import_data_cnt+'\')">'
					+'<div class="server server_on bgbox_b">'
					+'<div class="icon"></div>'
					+'<div class="txt">'+new_vm_name+'</div>'
					+'</div>'
					+'</li>';
			Import_Server_Config_Arr.push(object)
			console.log(i+"번째 import form data 입니다. : ",object);
		}
	}else{
		Import_Server_Config_Arr.push(import_form)
		add_server_html +='<li onclick="view_import(\''+import_data_cnt+'\')">'
						+'<div class="server server_on bgbox_b">'
						+'<div class="icon"></div>'
						+'<div class="txt">'+server_name+'</div>'
						+'</div>'
						+'</li>';

	}

	// Done 버튼 클릭 시 form은 비활성
	$(".import_servers_config").removeClass("active");

	// server List에 추가
	$("#mcis_server_list").prepend(add_server_html)
	console.log("import btn click and import form data : ",import_form)
	console.log("import data array : ",Import_Server_Config_Arr);
	import_data_cnt++;
	$("#import_form").each(function(){
		this.reset();
	})
}
function view_import(cnt){
	console.log('view import cnt : ',cnt);
	var select_form_data = Import_Server_Config_Arr[cnt]
	console.log('select_form_data : ', select_form_data);
	$(".simple_servers_config").removeClass("active")
	$(".expert_servers_config").removeClass("active")
	$(".import_servers_config").addClass("active")
}

