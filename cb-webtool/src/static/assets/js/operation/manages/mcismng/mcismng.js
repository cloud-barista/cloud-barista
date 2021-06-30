var selectedMcis = "";
$(document).ready(function(){
    checkLoadStatus();

    jQuery('.sc_box.scrollbar-inner').scrollbar();// CP / connectin의 구름 이미지들 창이 작아졌을 때 scroll 생기도록

    // MCIS List의 상단 의 checkbox 클릭시 전체 선택하도록 
    $("#th_chall").click(function() {
        if ($("#th_chall").prop("checked")) {
            $("input[name=chk]").prop("checked", true);
        } else {
            $("input[name=chk]").prop("checked", false);
        }
    })

    // 지도 표시
    // setRegionMap();

    selectedMcisID = $("#selected_mcis_id").val();
    
    // console.log(selectedMcisID);
    // Dashboard 등에서 선택한 MCIS Mng를 하면 해당 Mcis로만 보이도록 (전체에서 filter 기능만 수행)
    if( selectedMcisID != undefined && selectedMcisID != ""){
        // mcisList filter        
        filterTable("mcisListTable", "Name", selectedMcisID);
    }

    setTableHeightForScroll("mcisListTable", 700);


    // 상세 Tab 선택시 monitoring일 때 monitoring 조회
    $('a[data-toggle="tab"]').on('shown.bs.tab', function (e) {
        var target = $(e.target).attr("href") // activated tab
        if ( target == '#McisMonitoring'){
            var selectedMcisID = $("#selected_mcis_id").val();
            var selectedVmID = $("#selected_vm_id").val();
            showVmMonitoring(selectedMcisID,selectedVmID)
        }
        
    });
});

///////////// MCIS Handling //////////////

// 등록 form으로 이동
function createNewMcis(){// Manage_MCIS_Life_Cycle_popup.html
    var targetUrl = "/operation/manage" + "/mcismng/regform"
    // location.href = "/Manage/MCIS/reg"
    // $('#loadingContainer').show();
    // location.href = url;
    changePage(targetUrl)
}

// MCIS 제어 : 선택한 VM의 상태 변경 
// callMcisLifeCycle -> util.callMcisLifeCycle -> callbackMcisLifeCycle 순으로 호출 됨
function callMcisLifeCycle(type){
    var checked_nothing = 0;
    $("[id^='td_ch_']").each(function(){
       
        if($(this).is(":checked")){
            checked_nothing++;
            console.log("checked")
            var mcisID = $(this).val()
            mcisLifeCycle(mcisID, type);           
        }else{
            console.log("checked nothing")
           
        }
    })
    if(checked_nothing == 0){
        commonAlert("Please Select MCIS!!")
        return;
    }
}

// McisLifeCycle을 호출 한 뒤 return값 처리
function callbackMcisLifeCycle(resultStatus, resultData, type){
    var message = "MCIS "+type+ " complete!."
    if(resultStatus == 200 || resultStatus == 201){            
        commonAlert(message);
        location.reload();//완료 후 페이지를 reload -> 해당 mcis만 reload
        // 해당 mcis 조회
        // 상태 count 재설정
    }else{
        commonAlert("MCIS " + type + " failed!");
    }
}

// list에 선택된 MCIS 제거. 1개씩
function deleteMCIS(){
    var checkedCount = 0;
    var mcisID = "";
    $("[id^='td_ch_']").each(function(){
       
        if($(this).is(":checked")){
            checkedCount++;
            console.log("checked")
            mcisID = $(this).val();
            // 여러개를 지울 때 호출하는 함수를 만들어 여기에서 호출
        }else{
            console.log("checked nothing")
           
        }
    })

    if(checkedCount == 0){
        commonAlert("Please Select MCIS!!")
        return;
    }else if( checkedCount > 1){
        commonAlert("Please Select One MCIS at a time")
        return;
    }

    // TODO : 삭제 호출부분 function으로 뺼까?
    var url = "/operation/manages/mcismng/" + mcisID;               
    axios.delete(url,{})
        .then(result=>{
            console.log("get  Data : ",result.data);

            var statusCode = result.data.status;
            var message = result.data.message;
            
            if( statusCode != 200 && statusCode != 201) {
                commonAlert(message +"(" + statusCode + ")");
                return;
            }else{
                commonAlert(message);
                // TODO : MCIS List 조회
                location.reload();
            }
            
        }).catch((error) => {
            console.warn(error);
            console.log(error.response)
            var errorMessage = error.response.statusText;
            var statusCode = error.response.status;
            commonErrorAlert(statusCode, errorMessage) 
        });

}
////////////// MCIS Handling end //////////////// 



////////////// VM Handling ///////////
function addNewVirtualMachine(){
    var mcis_id = $("#mcis_id").val()
    var mcis_name = $("#mcis_name").val()
    // location.href = "/Manage/MCIS/reg/"+mcis_id+"/"+mcis_name
    location.href = "/operation/manages/mcismng/regform/"+mcis_id+"/"+mcis_name;
}

function vmLifeCycle(type){
    var mcisID = $("#mcis_id").val();
    var vmID = $("#vm_id").val();
    var vmName = $("#vm_name").val();
    
    // var checked =""
    // $("[id^='td_ch_'").each(function(){
    //     if($(this).is(":checked")){
    //         var checked_value = $(this).val();
    //         console.log("checked value : ",checked_value)
    //     }else{
    //         console.log("체크된게 없어!!")
    //     }
    // })
    // return;
    if(!mcisID){
        commonAlert("Please Select MCIS!!")
        return;
    }
    if(!vmID){
        commonAlert("Please Select VM!!")
        return;
    }
    
    // var nameSpace = NAMESPACE;
    console.log("Start LifeCycle method!!!")

    //url = CommonURL+"/ns/"+nameSpace+"/mcis/"+mcis_id+"/vm/"+vm_id+"?action="+type
    //var url = "/operation/manage" + "/mcis/" + mcis + "/vm/" + vm_id + "/action/" + type
    var url = "/operation/manages" + "/mcismng/proc/vmlifecycle";
    // url = "http://54.248.3.145:1323/tumblebug/ns/ns-01/mcis/mz-azure-mcis/vm/mz-azure-ubuntu1804-5?action=suspend";
    // var apiInfo = "Basic ZGVmYXVsdDpkZWZhdWx0"
    
    // /////////
    // axios.get(url,{
    //         headers:{
    //                 'Authorization': apiInfo
    //             }
    //         }).then(result=>{
    //             var data = result.data
    //             console.log(data)
    //         });
    /////////
    // + mcis + "/vm/" + vm_id + "/action/" + type
    
    console.log("life cycle3 url : ",url);
   
    var message = vmName+" "+type+ " complete!."  
    var obj = {
        mcisID : mcisID,
        vmID : vmID,
        lifeCycleType : type
     }
    axios.post(url,obj,{
        headers: { 
            'Content-type': 'application/json',
            // 'Authorization': apiInfo, 
        }
    // })
    // axios.post(url,{
    //     headers: { },
        // mcisID:mcis_id,
        // vmID:vm_id,
        // vmLifeCycleType:type
    }).then(result=>{
        var status = result.status
        
        console.log("life cycle result : ",result)
        var data = result.data
        console.log("result Message : ",data.message)
        if(status == 200 || status == 201){            
            commonAlert(message);
            location.reload();// TODO 일단은 Reaoad : 해당 영역(MCIS의 VM들 status 조회)를 refresh할 수 있는 기능 필요
            //show_mcis(mcis_url,"");
        }
    })
}
///////////// VM Handling end ///////////



const config_arr = new Array();

// refresh : mcis 및 vm정보조회
// 각 mcis 별 vmstatus 목록


// List Of MCIS 클릭 시 
// mcis 테이블의 선택한 row 강조( on )
// 해당 MCIS의 VM 상태목록 보여주는 함수 호출
function clickListOfMcis(id,index){
    console.log("click view mcis id :",id)
    
    // MCIS Info 에 mcis id 표시
    $("#mcis_id").val(id);

    // MCIS Info area set
    showServerListAndStatusArea(id,index);

    makeMcisScript(index);// export를 위한 script 준비
}

const mcisInfoDataList = new Array()// test_arr : mcisInfo 1개 전체, pageLoad될 때, refresh 할때 data를 set. mcis클릭시 상세내용 보여주기용 조회

// MCIS Info area 안의 Server List / Status 내용 표시
// 해당 MCIS의 모든 VM 표시
// TODO : 클릭했을 때 서버에서 조회하는것으로 변경할 것.
function showServerListAndStatusArea(mcis_id, mcisIndex){
    
    var mcisID =  $("#mcisID" + mcisIndex).val();
    var mcisName =  $("#mcisName" + mcisIndex).val();
    var mcisDescription =  $("#mcisDescription" + mcisIndex).val();
    var mcisStatus =  $("#mcisStatus" + mcisIndex).val();
    var mcisCloudConnections = $("#mcisCloudConnections" + mcisIndex).val();
    var vmTotalCountOfMcis = $("#mcisVmTotalCount" + mcisIndex).val();
    var vms = $("#mcisVmStatusList" + mcisIndex).val();

    $("#mcis_info_txt").text("[ "+ mcisName +" ]");
    $("#mcis_server_info_status").empty();
    $("#mcis_server_info_status").append('<strong>Server List / Status</strong>  <span class="stxt">[ '+mcisName+' ]</span>  Server('+vmTotalCountOfMcis+')')

    //
    $("#mcis_info_name").val(mcisName+" / "+mcisID)
    $("#mcis_info_description").val(mcisDescription);
    // $("#mcis_info_targetStatus").val(targetStatus);
    // $("#mcis_info_targetAction").val(targetAction);
    $("#mcis_info_cloud_connection").val(mcisCloudConnections)    //
    
    $("#mcis_name").val(mcisName)

    var mcis_badge = "";
    var mcisStatusIcon = "";
    if(mcisStatus == "running"){ mcisStatusIcon = "icon_running_db.png"        
    }else if(mcisStatus == "include" ){ mcisStatusIcon = "icon_stop_db.png"
    }else if(mcisStatus == "suspended"){mcisStatusIcon = "icon_stop_db.png"
    }else if(mcisStatus == "terminate"){mcisStatusIcon = "icon_terminate_db.png"
    }else{
        mcisStatusIcon = "icon_stop_db.png"
    }
    mcis_badge = '<img src="/assets/img/contents/' + mcisStatusIcon +'" alt=""/> '
    $("#service_status_icon").empty();
    $("#service_status_icon").append(mcis_badge)

    var vm_badge = "";
    $("[id^='mcisVmID_']").each(function(){		
        var mcisVm = $(this).attr("id").split("_")
        thisMcisIndex = mcisVm[1]
        vmIndexOfMcis = mcisVm[2]

        if( thisMcisIndex == mcisIndex){

            var vmID = $("#mcisVmID_" + mcisIndex + "_" + vmIndexOfMcis).val();
            var vmName = $("#mcisVmName_" + mcisIndex + "_" + vmIndexOfMcis).val();
            var vmStatus = $("#mcisVmStatus_" + mcisIndex + "_" + vmIndexOfMcis).val();
            vmStatus = vmStatus.toLowerCase();
            
            var vmStatusIcon ="bgbox_g";            
            if(vmStatus == "running"){ 
                vmStatusIcon ="bgbox_b"
            }else if(vmStatus == "include" ){
                vmStatusIcon ="bgbox_g"
            }else if(vmStatus == "suspended"){
                vmStatusIcon ="bgbox_g"
            }else if(vmStatus == "terminated"){
                vmStatusIcon ="bgbox_r"
            }else{
                vmStatusIcon ="bgbox_g"
            }
            vm_badge += '<li class="sel_cr ' + vmStatusIcon + '"><a href="javascript:void(0);" onclick="vmDetailInfo(\''+mcisID+'\',\''+mcisName+'\',\''+vmID+'\')"><span class="txt">'+vmName+'</span></a></li>';
        }
    });
    // console.log(vm_badge);
    $("#mcis_server_info_box").empty();
    $("#mcis_server_info_box").append(vm_badge);

    // var sta = mcisStatus;
    // var sl = sta.split("-");
    // var status = sl[0].toLowerCase()
    // var vm_badge = "";
    
    // var vmList = vms.split("@") // vm목록은 @
    // console.log("vmList " + vmList);
    // // for(var x in vmList){
    // for( var x= 0; x < vmList.length; x++){
    //     var vmInfo = vmList[x].split("|") // 이름과 상태는 "|"로 구분
    //     console.log("x " + x);
    //     console.log("vmInfo " + vmInfo);

    //     vmID = vmInfo[0];
    //     vmName = vmInfo[1];

    //     vmStatus = vmInfo[1].toLowerCase();

    //     var vmStatusIcon ="bgbox_g";
        
    //     if(vmStatus == "running"){ 
    //         vmStatusIcon ="bgbox_b"
    //     }else if(vmStatus == "include" ){
    //         vmStatusIcon ="bgbox_g"
    //         // vm_badge += '<li class="sel_cr bgbox_g"><a href="javascript:void(0);" onclick="click_view_vm(\''+mcisID+'\',\''+vmID+'\')"><span class="txt">'+vmName+'</span></a></li>';
    //     }else if(vmStatus == "suspended"){
    //         vmStatusIcon ="bgbox_g"
    //         // vm_badge += '<li class="sel_cr bgbox_g"><a href="javascript:void(0);" onclick="click_view_vm(\''+mcisID+'\',\''+vmID+'\')"><span class="txt">'+vmName+'</span></a></li>';
            
    //     }else if(vmStatus == "terminated"){
    //         vmStatusIcon ="bgbox_r"
    //         // vm_badge += '<li class="sel_cr bgbox_r"><a href="javascript:void(0);" onclick="click_view_vm(\''+mcisID+'\',\''+vmID+'\')"><span class="txt">'+vmName+'</span></a></li>';
            
    //     }else{
    //         vmStatusIcon ="bgbox_g"
    //         // vm_badge += '<li class="sel_cr bgbox_g"><a href="javascript:void(0);" onclick="click_view_vm(\''+mcisID+'\',\''+vmID+'\')"><span class="txt">'+vmName+'</span></a></li>';
    //     }
    //     vm_badge += '<li class="sel_cr ' + vmStatusIcon + '"><a href="javascript:void(0);" onclick="vmDetailInfo(\''+mcisID+'\',\''+mcisName+'\',\''+vmID+'\')"><span class="txt">'+vmName+'</span></a></li>';
    //     //console.log(vm_badge);
    //     $("#mcis_server_info_box").empty();
    //     $("#mcis_server_info_box").append(vm_badge);
    // }

    //Manage MCIS Server List on/off : table을 클릭하면 해당 Row 에 active style로 보여주기
    $(".dashboard .ds_cont .area_cont .listbox li.sel_cr").each(function(){
        var $sel_list = $(this);
        var $detail = $(".server_info");
        console.log($sel_list);
        console.log($detail);
        console.log(">>>>>");
        $sel_list.off("click").click(function(){
            $sel_list.addClass("active");
            $sel_list.siblings().removeClass("active");
            $detail.addClass("active");
            $detail.siblings().removeClass("active");
            $sel_list.off("click").click(function(){
                if( $(this).hasClass("active") ) {
                    $sel_list.removeClass("active");
                    $detail.removeClass("active");
                } else {
                    $sel_list.addClass("active");
                    $sel_list.siblings().removeClass("active");
                    $detail.addClass("active");
                    $detail.siblings().removeClass("active");
                }
            });
        });
    }); 
}

// VM 목록에서 VM 클릭시 해당 VM의 상세정보 
function vmDetailInfo(mcisID, mcisName, vmID){
    var url = "/operation/manages/mcismng/" + mcisID + "/vm/" + vmID
    axios.get(url,{})
        .then(result=>{
            console.log("get  Data : ",result.data);

            var statusCode = result.data.status;
            var message = result.data.message;
            
            if( statusCode != 200 && statusCode != 201) {
                commonAlert(message +"(" + statusCode + ")");
                return;
            }
            var data = result.data.VmInfo;
            var connectionConfig = result.data.ConnectionConfigInfo;
            
            var vmId = data.id;
            var vmName = data.name;
            var vmStatus = data.status;
            var vmDescription = data.description;
            var vmPublicIp = data.publicIP === undefined ? "" : data.publicIP;
            var vmSshKeyID = data.sshKeyId;
            //vm info
            $("#vm_id").val(vmId);   
            $("#vm_name").val(vmName);
            console.log("vm_id " + vmId + ", vm_name " + vmName + ", vm_SshKeyID " + vmSshKeyID)

            $("#manage_mcis_popup_vm_id").val(vmId)
            $("#manage_mcis_popup_mcis_id").val(mcisID)
            $("#manage_mcis_popup_sshkey_name").val(vmSshKeyID)
            
            $("#server_info_text").text('['+vmName+'/'+mcisName+']')
            $("#server_detail_info_text").text('['+vmName+'/'+mcisName+']')

            
            var vmBadge =""
            var vmStatusIcon = "icon_running_db.png";
            if(vmStatus == "Running"){
                vmStatusIcon = "icon_running_db.png";
            }else if(vmStatus == "include"){
                vmStatusIcon = "icon_stop_db.png";
            }else if(vmStatus == "Suspended"){
                vmStatusIcon = "icon_stop_db.png";
            }else if(vmStatus == "Terminated"){
                vmStatusIcon = "icon_terminate_db.png";
            }else{
                vmStatusIcon = "icon_stop_db.png";
            }
            vmBadge = '<img src="/assets/img/contents/' + vmStatusIcon +'" alt="' + vmStatus + '"/>'

            $("#server_detail_view_server_status").val(vmStatus);
            $("#server_info_status_img").empty()
            $("#server_info_status_img").append(vmBadge)

            $("#server_info_name").val(vmName +"/"+ vmID)
            $("#server_info_desc").val(vmDescription)

            // ip information
            $("#server_info_public_ip").val(vmPublicIp)
            $("#server_detail_info_public_ip_text").text("Public IP : "+vmPublicIp)
            $("#server_info_public_dns").val(data.publicDNS)
            $("#server_info_private_ip").val(data.privateIP)
            $("#server_info_private_dns").val(data.privateDNS)

            $("#server_detail_view_public_ip").val(vmPublicIp)
            $("#server_detail_view_public_dns").val(data.publicDNS)
            $("#server_detail_view_private_ip").val(data.privateIP)
            $("#server_detail_view_private_dns").val(data.privateDNS)
        
            $("#manage_mcis_popup_public_ip").val(vmPublicIp)

            //////vm detail tab////
            var vmDetail = data.cspViewVmDetail;
            //    //cspvmdetail
            // var vmDetailKeyValueList = vmDetail.KeyValueList
            var vmDetailKeyValueList = vmDetail.keyValueList
            var architecture = "";   
            if(vmDetailKeyValueList){
                for (var keyIndex in vmDetailKeyValueList) {
                    if( vmDetailKeyValueList[keyIndex].Key == "Architecture"){// ?? 이게 뭐지?
                        architecture = vmDetailKeyValueList[keyIndex].Value  
                        break;
                    }
                }
                // architecture = vmDetailKeyValueList.filter(item => item.Key === "Architecture")
                // console.log("architecture : ",architecture.length)
                // if(architecture.length > 0){
                //     architecture = architecture[0].Value
                //     console.log("architecture2 : ",architecture)                    
                // }
                // console.log("architecture = " + architecture)
                $("#server_info_archi").val(architecture)
                $("#server_detail_view_archi").val(architecture)
            }
            //    // server spec
            // var vmSecName = data.VmSpecName            
            var vmSpecName = vmDetail.vmspecName;
            $("#server_info_vmspec_name").val(vmSpecName)
            $("#server_detail_view_server_spec").text(vmSpecName)
            //var spec_id = data.specId
            var specId = data.specId
            // set_vmSpecInfo(spec_id);// memory + cpu  : TODO : spec정보는 자주변경되는것이 아닌데.. 매번 통신할 필요있나...
           
            var startTime = vmDetail.startTime
            $("#server_info_start_time").val(startTime)
            
            
            var locationInfo = data.location;
            var cloudType = locationInfo.cloudType;
            var cspIcon = ""
            if(cloudType == "aws"){
                cspIcon = "img_logo1"
            }else if(cloudType == "azure"){
                cspIcon = "img_logo5"
            }else if(cloudType == "gcp"){
                cspIcon = "img_logo7"
            }else if(cloudType == "cloudit"){
                cspIcon = "img_logo6"
            }else if(cloudType == "openstack"){
                cspIcon = "img_logo9"
            }else if(cloudType == "alibaba"){
                cspIcon = "img_logo4"
            }else{
                csp_icon = '<img src="/assets/img/contents/img_logo1.png" alt=""/>'
            }
            $("#server_info_csp_icon").empty()
            $("#server_info_csp_icon").append('<img src="/assets/img/contents/' + cspIcon + '.png" alt=""/>')
            $("#server_connection_view_csp").val(cloudType)
            $("#manage_mcis_popup_csp").val(cloudType)

            
            var latitude = locationInfo.latitude;
            var longitude = locationInfo.longitude;
            var briefAddr = locationInfo.briefAddr;
            var nativeRegion = locationInfo.nativeRegion;
           
            if( locationInfo){
                // 지도에 표시
                // $("#map").empty();
                // map = map_init();
                // var map = map_init_target('map2')
                // console.log("map");
                // let pointInfo = new Map();
                // pointInfo.set("title", "111");
                // pointInfo.set("vm_status", "222");
                // pointInfo.set("vm_id", "333");
                // pointInfo.set("id", "444");

                // drawMap(map, longitude, latitude, pointInfo);
                // console.log("drawMap");
                // setRegionMap(locationInfo);
            }
            // region zone locate
            
            //var region = data.region.region
            var region = data.region.Region            
            // var zone = data.region.zone
            var zone = data.region.Zone
            console.log(vmDetail.iid);
            $("#server_info_region").val(briefAddr +":"+region)
            $("#server_info_zone").val(zone)
            $("#server_info_cspVMID").val("cspVMID : "+vmDetail.iid.nameId)

            $("#server_detail_view_region").val(briefAddr +" : "+region)
            $("#server_detail_view_zone").val(zone)

            $("#server_connection_view_region").val(briefAddr +"("+region+")")
            $("#server_connection_view_zone").val(zone)

            // connection name
            var connectionName = data.connectionName;
            $("#server_info_connection_name").val(connectionName)
            $("#server_connection_view_connection_name").val(connectionName)

            // credential and driver info
            console.log("connectionConfig : ",connectionConfig)
            if(connectionConfig){
                var credentialName = connectionConfig.CredentialName
                var driverName = connectionConfig.DriverName
                $("#server_connection_view_credential_name").val(credentialName)
                $("#server_connection_view_driver_name").val(driverName)
            }

            // server id / system id
            $("#server_detail_view_server_id").val(data.id)
            // systemid 를 사용할 경우 아래 꺼 사용
            $("#server_detail_view_server_id").val(vmDetail.iid.systemId)
            
            // image id
            var imageIId = vmDetail.imageIId.nameId
            var imageId = data.imageId
            getCommonVmImageInfo('mcisvmdetail', imageId) // 
            $("#server_detail_view_image_id").text(imageId+"("+imageIId+")")

            //vpc subnet
            var vpcId = vmDetail.vpcIID.nameId
            var vpcSystemId = vmDetail.vpcIID.systemId
            var subnetId = vmDetail.subnetIID.nameId
            var subnetSystemId = vmDetail.subnetIID.systemId
            var eth = vmDetail.networkInterface
            $("#server_detail_view_vpc_id").text(vpcId+"("+vpcSystemId+")")
            // set_vmVPCInfo(vpcId, subnetId);

            $("#server_detail_view_subnet_id").text(subnetId+"("+subnetSystemId+")")
            $("#server_detail_view_eth").val(eth)

            // user account
            $("#server_detail_view_access_id_pass").val(vmDetail.vmuserId +"/ *** ")
            $("#server_detail_view_user_id_pass").val(data.vmUserAccount +"/ *** ")
            $("#manage_mcis_popup_user_name").val(data.vmUserAccount)


            var append_sg = ''

            var sg_arr = vmDetail.securityGroupIIds
            console.log(sg_arr);
            if(sg_arr){
                sg_arr.map((item,index)=>{           
                    console.log("item index = " + index)
                    append_sg +='<a href="javascript:void(0);" onclick="getCommonVmSecurityGroupInfo(\'mcisvm\',\''+item.nameId+'\');"title="'+item.nameId+'" >'+item.nameId+'</a> '
                })
            }  
            console.log("append sg : ",append_sg)
            
            $("#server_detail_view_security_group").empty()
            $("#server_detail_view_security_group").append(append_sg);

            $("#server_detail_view_keypair_name").val(vmDetail.keyPairIId.nameId)
            // ... TODO : 우선 제어명령부터 처리. 나중에 해당항목 mapping하여 확인 
            ////// vm connection tab //////


            $("#selected_mcis_id").val(mcisID);
            $("#selected_vm_id").val(vmID);
            var installMonAgent = data.monAgentStatus;
            console.log("install mon agent : ",installMonAgent)
            if(installMonAgent == "installed"){
                var isWorking = checkDragonFlyMonitoringAgent(mcisID, vmID);
                if( isWorking){
                    $("#mcis_detail_info_check_monitoring").prop("checked",true)
                    $("#mcis_detail_info_check_monitoring").attr("disabled",true)
                }else{
                    $("#mcis_detail_info_check_monitoring").prop("checked",false)
                    $("#mcis_detail_info_check_monitoring").attr("disabled",false)
                }                
            }else{
                $("#mcis_detail_info_check_monitoring").prop("checked",false)
                $("#mcis_detail_info_check_monitoring").attr("disabled",false)
            }

            ////// vm mornitoring tab 으로 이동 //////            
            // install Mon agent
            // showVmMonitoring(mcisID,vmID)            
        }
    // ).catch(function(error){
    //     var statusCode = error.response.data.status;
    //     var message = error.response.data.message;
    //     commonErrorAlert(statusCode, message)        
    // });
    ).catch((error) => {
        console.warn(error);
        console.log(error.response)
        var errorMessage = error.response.statusText;
        var statusCode = error.response.status;
        commonErrorAlert(statusCode, errorMessage) 
    });

    // $("#Detail").show();// 첫번째 Detail tab 표시.
    $('[href="#Detail"]').tab('show');

    /////////////////////


//    // credential and driver info
//    console.log("config arr2 : ",config_arr)
//    console.log("connection_name :",connection_name)
//    var arr_config = config_arr
//    console.log("arr_config : ",arr_config);
//    if(arr_config){
//        var config_info = arr_config.filter(cred => cred.ConfigName === connection_name)[0]
//        console.log("inner config info : ",config_info)
//        console.log("config_info : ",config_info)
//        var credentialName = config_info.CredentialName
//        var driverName = config_info.DriverName
//        $("#server_connection_view_credential_name").val(credentialName)
//        $("#server_connection_view_driver_name").val(driverName)
//    }
  

   
   
//    // server id / system id
//    $("#server_detail_view_server_id").val(select_vm.id)
//    // systemid 를 사용할 경우 아래 꺼 사용
//    //$("#server_detail_view_server_id").val(vm_detail.IId.SystemId)
  
//    // image id
//    var imageIId = vm_detail.ImageIId.NameId
//    var imageId = select_vm.imageId
//    set_vmImageInfo(imageId)
//    $("#server_detail_view_image_id_text").text(imageId+"("+imageIId+")")

//    //vpc subnet
//    var vpcId = vm_detail.VpcIID.NameId
//    var vpcSystemId = vm_detail.VpcIID.SystemId
//    var subnetId = vm_detail.SubnetIID.NameId
//    var subnetSystemId = vm_detail.SubnetIID.SystemId
//    var eth = vm_detail.NetworkInterface
//    $("#server_detail_view_vpc_id_text").text(vpcId+"("+vpcSystemId+")")
//    set_vmVPCInfo(vpcId, subnetId);

//    $("#server_detail_view_subnet_id_text").text(subnetId+"("+subnetSystemId+")")
//    $("#server_detail_view_eth_text").val(eth)

//    // install Mon agent
//    var installMonAgent = select_vm.monAgentStatus
//    checkDragonFly(mcis_id,vm_id)

//    // device info
//    var root_device_type = vm_detail.VMBootDisk
//    var root_device = vm_detail.VMBootDisk
//    var block_device = vm_detail.VMBlockDisk
//    $("#server_detail_view_root_device_type").val(root_device_type)
//    $("#server_detail_view_root_device").val(root_device)
//    $("#server_detail_view_block_device").val(block_device)

//     // key pair info
   
//     $("#server_detail_view_keypair_name").val(vm_detail.KeyPairIId.NameId)
//     var sshkey = vm_detail.KeyPairIId.NameId
//     if(sshkey){
//        set_vmSSHInfo(sshkey)
//     }
//     // user account
//     $("#server_detail_view_access_id_pass").val(vm_detail.VMUserId +"/"+vm_detail.VMUserPasswd)
//     $("#server_detail_view_user_id_pass").val(select_vm.vmUserAccount +"/"+select_vm.vmUserPassword)
//     $("#manage_mcis_popup_user_name").val(select_vm.vmUserAccount)
    
//     // namespace 
//     var ns_id = NAMESPACE
//     $("#manage_mcis_popup_ns_id").val(ns_id)
    

//     // security Gorup
//    var append_sg = ''

//    var sg_arr = vm_detail.SecurityGroupIIds
//    if(sg_arr){
//        //여기서 호출해서 세부 값을 가져 오자
       
//        sg_arr.map((item,index)=>{
           
//            append_sg +='<a href="javascript:void(0);" onclick="set_vmSecurityGroupInfo(\''+item.NameId+'\');"title="'+item.NameId+'" >'+item.NameId+'</a> '
//        })
//    }
  
//    console.log("append sg : ",append_sg)
   
//    $("#server_detail_view_security_group").empty()
//    $("#server_detail_view_security_group").append(append_sg);

}


// 조회 성공 시 Monitoring Tab 표시
function showVmMonitoring(mcisID, vmID){
    $("#mcis_detail_info_check_monitoring").prop("checked",true)
    $("#mcis_detail_info_check_monitoring").attr("disabled",true)
    // $("#Monitoring_tab").show();
    //var duration = "5m"
    var duration = "30m"
    var period_type = "m"
    var metric_arr = ["cpu","memory","disk","network"];
    var statisticsCriteria = "last";
    // TODO : Analytics View 는 안보이게 
    for(var i in metric_arr){
        getVmMetric("canvas_"+i,metric_arr[i],mcisID,vmID,metric_arr[i],period_type,statisticsCriteria,duration);
    }    
    //$("#Monitoring_tab").hide();
    
 }
 

// getVMMetric 는 mcis.chart.js로 이동 

////////////////

// MCIS script export
function mcisScriptExport(){
    
    // var mcisID = $("#mcis_id").val();
    var vmID = $("#vm_id").val();
    var vmName = $("#vm_name").val();

    var checkedCount = 0;
    var mcisID = "";
    $("[id^='td_ch_']").each(function(){
       
        if($(this).is(":checked")){
            checkedCount++;
            console.log("checked")
            mcisID = $(this).val();
            // 여러개를 지울 때 호출하는 함수를 만들어 여기에서 호출
        }else{
            console.log("checked nothing")
           
        }
    })

    if(checkedCount == 0){
        commonAlert("Please Select MCIS!!")
        return;
    }else if( checkedCount > 1){
        commonAlert("Please Select One MCIS at a time")
        return;
    }

    // 위 값으로 mcisIndex, vmIndex 를 찾자
    var mcisIndex = 0;
    var vmIndex = 0;
    //mcisID{{$index}}
    console.log("mcisScriptExport start")
    $("[id^='mcisID']").each(function(){
        if( mcisID == $(this).val()){
            mcisIndex = $(this).attr("id").replace("mcisID", "")
            return false;
        }
    });
    console.log("index " + mcisIndex)
    if( $("#m_mcisExportScript_" + mcisIndex).val() == ""){
        makeMcisScript(mcisIndex);
    }
    console.log("mcisscript")
    console.log($("#m_mcisExportScript_" + mcisIndex).val())
    saveToMcisAsJsonFile(mcisIndex, vmIndex);
}

// vm script export
function vmScriptExport(){
    
    var mcisID = $("#mcis_id").val();
    var vmID = $("#vm_id").val();
    var vmName = $("#vm_name").val();
    console.log("vmScriptExport start")
    // 위 값으로 mcisIndex, vmIndex 를 찾자
    var mcisIndex = 0;
    var vmIndex = 0;
    $("[id^='mcisVmID_']").each(function(){
        if( vmID == $(this).val()){
            var mcisVm = $(this).attr("id").split("_")
            mcisIndex = mcisVm[1]
            vmIndex = mcisVm[2]
            return false;
        }
    });

    if(!mcisID){
        commonAlert("Please Select MCIS!!")
        return;
    }
    if(!vmID){
        commonAlert("Please Select VM!!")
        return;
    }
    
    console.log("index " + mcisIndex + " , " + vmIndex)
    if( $("#m_vmExportScript_" + mcisIndex + "_" + vmIndex).val() == ""){
        makeVmScript(mcisIndex, vmIndex);
    }
    console.log("vmscript")
    console.log($("#m_vmExportScript_" + mcisIndex + "_" + vmIndex).val())
    saveToVmAsJsonFile(mcisIndex, vmIndex);
}
// mcis를 선택하면 해당 mcis를 export할 준비를 함
// lifecycle 의 ExportScriptOfMcis 를 통해 선택한 mcis script를 file로 저장
function makeMcisScript(mcisIndex){
    var vms = 'mcisVmID_' + mcisIndex + "_";
    var vmIndex = 0;
    // vmScript 먼저 생성
    console.log("in makeMcisScript " + mcisIndex + " vms : " + vms);
    $("[id^='" + vms +"']").each(function(){
        makeVmScript(mcisIndex, vmIndex);
        vmIndex++;
    });
    console.log(" gogo mcis script");
	var mcisIDVal = $("#m_mcisID_" + mcisIndex).val();
    var mcisNameVal = $("#m_mcisName_" + mcisIndex).val();
	var mcisLabelVal = $("#m_mcisLabel_" + mcisIndex).val();
	var mcisDescriptionVal = $("#m_mcisDescription_" + mcisIndex).val();
	var mcisInstallMonAgentVal = $("#m_mcisInstallMonAgent_" + mcisIndex).val();
	

    var paramValueAppend = '"';
    var mcisCreateScript = "";
    console.log(" gogo mcis script2 ");
	mcisCreateScript += '{	';
	mcisCreateScript += paramValueAppend + 'name' + paramValueAppend + ' : ' + paramValueAppend + mcisNameVal + paramValueAppend;
	mcisCreateScript += ',' + paramValueAppend + 'description' + paramValueAppend + ' : ' + paramValueAppend + mcisDescriptionVal + paramValueAppend;
	mcisCreateScript += ',' + paramValueAppend + 'label' + paramValueAppend + ' : ' + paramValueAppend + mcisLabelVal + paramValueAppend;
	mcisCreateScript += ',' + paramValueAppend + 'installMonAgent' + paramValueAppend + ' : ' + paramValueAppend + mcisInstallMonAgentVal + paramValueAppend;
	console.log(mcisCreateScript);
    // vmScript 가져오기
    console.log("vm Size =" + vmIndex);
    mcisCreateScript += ',' + paramValueAppend + 'vm' + paramValueAppend + ':[';
    var addedVmIndex = 0;
    for( var i = 0; i < vmIndex; i++){
        var vmScript = $("#m_vmExportScript_" + mcisIndex + "_" + i).val();// 여기에 담겨있음.(위에서 먼저 호출해서 생성 해 둠)

        console.log(i );
        console.log(vmScript);

        if( vmScript == undefined) continue;// VM이 Terminated 된 경우 등에서는 vmScript가 정상적으로 생성되지 않음.

        if( addedVmIndex > 0) mcisCreateScript += ",";        
        
        mcisCreateScript += vmScript;
        addedVmIndex++;
    }
    mcisCreateScript += ']';
    mcisCreateScript += '}';

    $("#m_exportFileName_" + mcisIndex).val(mcisNameVal);
	$("#m_mcisExportScript_" + mcisIndex).val(mcisCreateScript);

    console.log("mcisCreateScript============");
    console.log(mcisCreateScript);
}

// vm을 선택하면 해당 vm을 export할 준비를 함
function makeVmScript(mcisIndex, vmIndex){
	console.log("in makeVmScript" +  mcisIndex + " : " + vmIndex)
	var connectionNameVal = $("#m_vmConnectionName_" + mcisIndex + "_" + vmIndex).val();
	var descriptionVal = $("#m_vmDescription_" + mcisIndex + "_" + vmIndex).val();
	var imageIdVal = $("#m_vmImageId_" + mcisIndex + "_" + vmIndex).val();
	var labelVal = $("#m_vmLabel_" + mcisIndex + "_" + vmIndex).val();
	var nameVal = $("#m_vmName_" + mcisIndex + "_" + vmIndex).val();
	var securityGroupIdsVal = $("#m_vmSecurityGroupIds_" + mcisIndex + "_" + vmIndex).val();
	var specIdVal = $("#m_vmSpecId_" + mcisIndex + "_" + vmIndex).val();
	var sshKeyIdVal = $("#m_vmSshKeyId_" + mcisIndex + "_" + vmIndex).val();
	var subnetIdVal = $("#m_vmSubnetId_" + mcisIndex + "_" + vmIndex).val();
	var vNetIdVal = $("#m_vmVnetId_" + mcisIndex + "_" + vmIndex).val();
	var vmGroupSizeVal = $("#m_vmGroupSize_" + mcisIndex + "_" + vmIndex).val();
	var vmUserAccountVal = $("#m_vmUserAccount_" + mcisIndex + "_" + vmIndex).val();
	var vmUserPasswordVal = $("#m_vmUserPassword_" + mcisIndex + "_" + vmIndex).val();

	var paramValueAppend = '"';
	var vmCreateScript = "";
	vmCreateScript += '{	';
	vmCreateScript += paramValueAppend + 'connectionName' + paramValueAppend + ' : ' + paramValueAppend + connectionNameVal + paramValueAppend;
	vmCreateScript += ',' + paramValueAppend + 'description' + paramValueAppend + ' : ' + paramValueAppend + descriptionVal + paramValueAppend;
	vmCreateScript += ',' + paramValueAppend + 'imageId' + paramValueAppend + ' : ' + paramValueAppend + imageIdVal + paramValueAppend;
	vmCreateScript += ',' + paramValueAppend + 'label' + paramValueAppend + ' : ' + paramValueAppend + labelVal + paramValueAppend;
	vmCreateScript += ',' + paramValueAppend + 'name' + paramValueAppend + ' : ' + paramValueAppend + nameVal + paramValueAppend;
	
    //vmCreateScript += ',' + paramValueAppend + 'securityGroupIds' + paramValueAppend + ' : ' + paramValueAppend + securityGroupIdsVal + paramValueAppend;

    console.log(securityGroupIdsVal);
    var sgVal = securityGroupIdsVal.replace(/\[/gi, "");
    sgVal = sgVal.replace(/\]/gi, "");
    var sgArr = sgVal.split(",");

    vmCreateScript += ',' + paramValueAppend + 'securityGroupIds' + paramValueAppend + ' : [';
    
    for( var i = 0; i < sgArr.length; i++){
        if( i > 0) vmCreateScript += ','
        vmCreateScript += paramValueAppend + sgArr[i] + paramValueAppend;
        console.log("securityGroupIdsVal [" + i + "] =" + sgArr[i]);
    }
    vmCreateScript += ']';
    
	vmCreateScript += ',' + paramValueAppend + 'specId' + paramValueAppend + ' : ' + paramValueAppend + specIdVal + paramValueAppend;
	vmCreateScript += ',' + paramValueAppend + 'sshKeyId' + paramValueAppend + ' : ' + paramValueAppend + sshKeyIdVal + paramValueAppend;
	vmCreateScript += ',' + paramValueAppend + 'subnetId' + paramValueAppend + ' : ' + paramValueAppend + subnetIdVal + paramValueAppend;
	vmCreateScript += ',' + paramValueAppend + 'vNetId' + paramValueAppend + ' : ' + paramValueAppend + vNetIdVal + paramValueAppend;
	vmCreateScript += ',' + paramValueAppend + 'vmGroupSize' + paramValueAppend + ' : ' + paramValueAppend + vmGroupSizeVal + paramValueAppend;
	vmCreateScript += ',' + paramValueAppend + 'vmUserAccount' + paramValueAppend + ' : ' + paramValueAppend + vmUserAccountVal + paramValueAppend;
	vmCreateScript += ',' + paramValueAppend + 'vmUserPassword' + paramValueAppend + ' : ' + paramValueAppend + vmUserPasswordVal + paramValueAppend;
	vmCreateScript += '}';

	
	$("#m_exportFileName_" + mcisIndex + "_"+ vmIndex).val(nameVal);
	$("#m_vmExportScript_" + mcisIndex + "_"+vmIndex).val(vmCreateScript);
    console.log("vmCreateScript============" + mcisIndex + ":" + vmIndex);
    console.log(vmCreateScript);
}

// json 파일로 저장
function saveToMcisAsJsonFile(mcisIndex){
	var fileName = "MCIS_" + $("#m_exportFileName_" + mcisIndex).val();
	var exportScript = $("#m_mcisExportScript_" + mcisIndex).val();
	
    saveFileProcess(fileName, exportScript);
}
function saveToVmAsJsonFile(mcisIndex, vmIndex){
	var fileName = "VM_" + $("#m_exportFileName_" + mcisIndex + "_" + vmIndex).val();
	var exportScript = $("#m_vmExportScript_" + mcisIndex + "_"  + vmIndex).val();
													
    saveFileProcess(fileName, exportScript);
}

// 파일명, script대로 파일 생성
function saveFileProcess(fileName, exportScript){

	var element = document.createElement('a');
	// element.setAttribute('href', 'data:text/plain;charset=utf-8,' + encodeURIComponent(exportScript));
	element.setAttribute('href', 'data:text/json;charset=utf-8,' + encodeURIComponent(exportScript));
	// element.setAttribute('download', fileName);
	element.setAttribute('download', fileName + ".json");

	element.style.display = 'none';
	document.body.appendChild(element);

	element.click();

	document.body.removeChild(element);
}

function getSecurityGroupCallbackSuccess(caller, data){
    var html = ""
    var firewallRules = data.firewallRules
    
    $("#register_box").modal()    
    firewallRules.map(item=>(        html +='<tr>'
                +'<td class="btn_mtd" data-th="fromPort">'+item.fromPort+' <span class="ov off"></span></td>'
                +'<td class="overlay hidden" data-th="toPort">'+item.toPort+'</td>'
                +'<td class="overlay hidden" data-th="toProtocol">'+item.ipProtocol+'</td>'
                +'<td class="overlay hidden " data-th="direction">'+item.direction+'</td>'
                +'</tr>'
    ))
    $("#manage_mcis_popup_sg").empty()
    $("#manage_mcis_popup_sg").append(html)
}

function getSecurityGroupCallbackFail(error){

}

// 지도에 marker로 region 표시.  
// TODO : default로 지도 표시한 뒤 location만 받아서 marker만 추가하도록 변경필요
function setRegionMap(locationInfo){
//         var lat            = 38.1300;
//         var lon            = -78.4500;
// //     var zoom           = 1;

//     lat = 37.413294;
//     lon = 126.734086;// 서울
    console.log(locationInfo);
    
    var latitude =  37.413294;;
    var longitude = 126.734086;// 서울
    var briefAddr = "seoul region"
    var nativeRegion = "east asia";
    console.log("location is ===")
    console.log(locationInfo)
    if( locationInfo){
        latitude = locationInfo.latitude;
        longitude = locationInfo.longitude;
        briefAddr = locationInfo.briefAddr;
        nativeRegion = locationInfo.nativeRegion;
    }
    console.log("latitude= " + latitude + ", longitude = " + longitude)
    const iconFeature = new ol.Feature({
        geometry: new ol.geom.Point(ol.proj.fromLonLat([longitude, latitude])),
        // geometry: new ol.geom.Point(ol.proj.fromLonLat([-2, 53])),
        name: nativeRegion,
      });
      
      const map = new ol.Map({
        target: 'regionMap',
        layers: [
          new ol.layer.Tile({
            source: new ol.source.OSM(),
          }),
          new ol.layer.Vector({
            source: new ol.source.Vector({
              features: [iconFeature]
            }),
            style: new ol.style.Style({
              image: new ol.style.Icon({
                anchor: [0.5, 46],
                anchorXUnits: 'fraction',
                anchorYUnits: 'pixels',
                src: '/assets/img/marker/black.png'
              })
            })
          })
        ],
        view: new ol.View({
          center: ol.proj.fromLonLat([longitude, latitude]),
        //   zoom: 1    //전 세계 표시
        zoom: 4
        // zoom: 6
        })
      });

    //   $("#regionMap").css("display", "block");
}


function getCommonVmImageInfoCallbackSuccess(caller, imageInfo){
    // var imageInfo = data;
    var html = ""
    console.log("image info : ",imageInfo)
    html +='<a href="javascript:void(0);" title="'+imageInfo.cspImageName+'">'+imageInfo.id+'</a>'
          +'<div class="bb_info">Image Name : '+imageInfo.name+', GuestOS:'+imageInfo.guestOS+'</div>'
   
    $("#server_detail_view_image_id").empty();
    $("#server_detail_view_image_id").append(html);
    $("#server_info_os").val(imageInfo.guestOS);
    $("#server_detail_view_os").val(imageInfo.guestOS);
    bubble_box();
}

function getCommonVmImageInfoCallbackFail(caller, data){
    // -- fail 나더라도 그냥 넘어감.
}

function bubble_box(){
    $(".bubble_box .box").each(function(){
		var $list = $(this);
		var bubble =  $list.find('.bb_info');
		var menuTime;
		$list.mouseenter(function(){
			bubble.fadeIn(300);
			clearTimeout(menuTime);
		}).mouseleave(function(){
			clearTimeout(menuTime);
    	menuTime = setTimeout(mTime, 100);
		});
		function mTime() {
	    bubble.stop().fadeOut(100);
	  }
	});
}

function remoteCommandMcis(commandWord){    
    // mcis가 선택되어 있어야 하고
    var checked_nothing = 0;
    $("[id^='td_ch_']").each(function(){
       
        if($(this).is(":checked")){
            checked_nothing++;
            console.log("checked")
            var mcisID = $(this).val()            
            postRemoteCommandMcis(mcisID, commandWord);
        }else{
            console.log("checked nothing")
           
        }
    })
    if(checked_nothing == 0){
        commonAlert("Please Select MCIS!!")
        return;
    }

}


function remoteCommandVmMcis(commandWord){    
    // VM 선택되어 있어야     
    var mcisID = $("#mcis_id").val();
    var vmID = $("#vm_id").val();
    var vmName = $("#vm_name").val();
    console.log("remoteCommandVmMcis start")
    // 위 값으로 mcisIndex, vmIndex 를 찾자
    var mcisIndex = 0;
    var vmIndex = 0;
    $("[id^='mcisVmID_']").each(function(){
        if( vmID == $(this).val()){
            var mcisVm = $(this).attr("id").split("_")
            mcisIndex = mcisVm[1]
            vmIndex = mcisVm[2]
            return false;
        }
    });

    if(!mcisID){
        commonAlert("Please Select MCIS!!")
        return;
    }
    if(!vmID){
        commonAlert("Please Select VM!!")
        return;
    }
    console.log(" commandWord = " + commandWord);
    if(!commandWord){
        commonAlert("Please type command!!")
        return;
    }

    // $("#manage_mcis_popup_sshkey_name").val(vmSshKeyID)
    // $("#server_info_public_ip").val(vmPublicIp)
    // $("#server_detail_view_public_ip").val(vmPublicIp)
    // $("#server_info_connection_name").val(connectionName)

    // $("#server_detail_view_access_id_pass").val(vmDetail.vmuserId +"/ *** ")
    // $("#server_detail_view_user_id_pass").val(data.vmUserAccount +"/ *** ")
    // $("#manage_mcis_popup_user_name").val(data.vmUserAccount)

    // var publicIp = $("#server_info_public_ip").val();
    // var accessId = $("#manage_mcis_popup_user_name").val();
    // var sshKeyId = $("#manage_mcis_popup_sshkey_name").val();

    postRemoteCommandVmOfMcis(mcisID, vmID, commandWord);

}