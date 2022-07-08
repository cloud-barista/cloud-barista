var selectedMcis = "";// 선택 된 mcisId :// 다른 화면 등에서 mcisID를 넘겨주는 경우 사용
$(document).ready(function () {
    AjaxLoadingShow(true);
    checkLoadStatus();

    jQuery('.sc_box.scrollbar-inner').scrollbar();// CP / connectin의 구름 이미지들 창이 작아졌을 때 scroll 생기도록

    // MCIS List의 상단 의 checkbox 클릭시 전체 선택하도록
    $("#th_chall").click(function () {
        if ($("#th_chall").prop("checked")) {
            $("input[name=chk]").prop("checked", true);
        } else {
            $("input[name=chk]").prop("checked", false);
        }
    })

    $("#mcisListTable").each(function () {
        var $sel_list = $(this);
        var $detail = $(".server_info");
        $status = $(".server_status")


        $sel_list.off("click").click(function () {
            console.log("sel_list")
            console.log($sel_list);
            console.log("this---")
            console.log($(this))
        });
    });



    // selectedMcisID = $("#selected_mcis_id").val();

    // console.log(selectedMcisID);
    // Dashboard 등에서 선택한 MCIS Mng를 하면 해당 Mcis로만 보이도록 (전체에서 filter 기능만 수행)
    // if (selectedMcisID != undefined && selectedMcisID != "") {
    //     // mcisList filter
    //     // filterTable("mcisListTable", "Name", selectedMcisID);
    //     clickListOfMcis(selectedMcisID);
    // }

    setTableHeightForScroll("mcisListTable", 700);


    // 상세 Tab 선택시 tab=monitoring일 때 monitoring 조회
    $('a[data-toggle="tab"]').on('shown.bs.tab', function (e) {
        var target = $(e.target).attr("href") // activated tab
        if (target == '#McisMonitoring') {
            var selectedMcisID = $("#selected_mcis_id").val();
            var selectedVmID = $("#selected_vm_id").val();
            showVmMonitoring(selectedMcisID, selectedVmID)
        } else if (target == '#McisConnection') {
            // 지도 표시
            var pointInfo = new Object();
            pointInfo.id = "1"
            pointInfo.name = $("#server_connection_view_region").val();
            pointInfo.cloudType = $("#server_connection_view_csp").val();;//
            pointInfo.latitude = $("#server_location_latitude").val();;//
            pointInfo.longitude = $("#server_location_longitude").val();//
            pointInfo.markerIndex = 1
            setMap(pointInfo);
        }
    });

    $('#alertResultArea').on('hidden.bs.modal', function () {// bootstrap 3 또는 4
        console.log("alertResultArea")
        refreshMcisStatusData();
    })

    // "NameSpaceList":      nsList,
    // "CloudOSList":                   cloudOsList,
    // "RegionList":                    regionInfoList,
    // "CloudConnectionConfigInfoList": cloudConnectionConfigInfoList,
    // "VMImageList":                   virtualMachineImageInfoList,
    // "VMSpecList":                    vmSpecInfoList,
    // "VNetList":                      vNetInfoList,
    // "SecurityGroupList":             securityGroupInfoList,

    //function getCommonNameSpaceList(caller, isCallback, targetObjId)
    getCommonNameSpaceList("mcismng", false, "menuNameSpaceList")// 왼쪽 메뉴 표시용
    // connection
    getCommonCloudConnectionList("mcismng", true)
    // // region
    // getCommonRegionList("mcismng");

    // // resource
    // // network(vnet)
    // getCommonNetworkList("mcismng")
    // // securitygroup
    // getCommonSecurityGroupList("mcismng")
    // // sshkey
    // getCommonSshKeyList("mcismng")

    // // image
    // getCommonVirtualMachineImageList("mcismng")

    // // spec
    // getCommonVirtualMachineSpecList("mcismng")

    // //
    getCommonMcisList("mcismngready", true, "", "simple")
});

// server info 영역 보이기/숨기기
function showServerInfoArea() {
    // $status = $(".server_status")// mcis의 server 목록
    // var $detail = $(".server_info");// server 상세정보  -> 여기에서는 바로 보여줄 필요 없는데....
    // // mcis list table에 선택한 mcis tr을 active, 다른 것들은 active 제거
    // $sel_list.addClass("active");
    // $sel_list.siblings().removeClass("active");
    //
    // // server info
    // $status.addClass("view");
    // $sel_list.off("click").click(function () {
    //     if ($(this).hasClass("active")) {
    //         $sel_list.removeClass("active");
    //         $detail.removeClass("active");
    //         $status.removeClass("view");
    //     } else {
    //         $sel_list.addClass("active");
    //         $sel_list.siblings().removeClass("active");
    //         // $detail.addClass("active");
    //         // $detail.siblings().removeClass("active");
    //         $status.addClass("view");
    //     }
    // });
}

var TotalCloudConnectionList = new Array();
////////// 화면 Load 시 가져온 값들을 set하는 function들
function getCloudConnectionListCallbackSuccess(caller, connectionConfigList, sortType) {
    var totalProviderCount = 0;
    var totalConnectionConfigCount = 0;
    var providerConnectionMap = new Map();
    if (!isEmpty(connectionConfigList) && connectionConfigList.length) {
        TotalCloudConnectionList = connectionConfigList;
        totalConnectionConfigCount = connectionConfigList.length;

        var providerArr = new Array();
        for (var itemIndex in connectionConfigList) {
            var aConnectionConfig = connectionConfigList[itemIndex]
            // console.log(aConnectionConfig)
            if (providerConnectionMap.has(aConnectionConfig.ProviderName)) {
                providerConnectionMap.set(aConnectionConfig.ProviderName, providerConnectionMap.get(aConnectionConfig.ProviderName) + 1)
            } else {
                providerConnectionMap.set(aConnectionConfig.ProviderName, 1)
            }

        }
        totalProviderCount = providerArr.length
    }
    // console.log(totalConnectionConfigCount + " : " + providerConnectionMap.size)
    $("#connectionCount").text(totalConnectionConfigCount)
    $("#providerCount").text(providerConnectionMap.size)


    for (let item of providerConnectionMap) {
        console.log(item[0] + ' , ' + item[1]);
        $("#cpConnectionDetail").append('<li class="bg_etc bg_' + item[0].toLowerCase() + '"><a href="javascript:void(0);"><span class="conn_cnt">' + item[1] + '<div class="conn_tit">' + item[0] + '</div></span></a></li>');
    }
}

// MCIS 목록 조회 후 화면에 Set
function getMcisListCallbackSuccess(caller, mcisList) {
    totalMcisListObj = mcisList;
    console.log("total mcis:", totalMcisListObj);
    setToTalMcisStatus();// mcis상태 표시 를 위해 필요
    setTotalVmStatus();// mcis 의 vm들 상태표시 를 위해 필요
    setTotalConnection();// Mcis의 provider별 connection 표시를 위해 필요

    displayMcisListTable();
    // console.log("caller = " + caller)
    if (caller == "mcismngready") {
        var selectedMcisID = $("#selected_mcis_id").val();
        console.log(" selectedMcisID for filter " + selectedMcisID)
        // filterTable("mcisListTable", "Name", selectedMcisID);
        if (checkEmptyString(selectedMcisID)) {
            clickListOfMcis(selectedMcisID);
        }

    }
    AjaxLoadingShow(false);
    // // MCIS Status
    // var totalMcisCnt = 0;
    // var mcisStatusCountMap = new Map();
    // mcisStatusCountMap.set("running", 0);
    // mcisStatusCountMap.set("stop", 0);  // partial 도 stop으로 보고있음.
    // mcisStatusCountMap.set("terminate", 0);
    //
    // var totalServerCnt = 0;
    // var totalVmStatusCountMap = new Map();
    // totalVmStatusCountMap.set("running", 0);
    // totalVmStatusCountMap.set("stop", 0);  // partial 도 stop으로 보고있음.
    // totalVmStatusCountMap.set("terminate", 0);
    //
    // if(!isEmpty(mcisList) && mcisList.length > 0 ){
    //     //totalMcisCnt = mcisList.length;
    //     var addMcis = "";
    //     for(var mcisIndex in mcisList){
    //         var aMcis = mcisList[mcisIndex]
    //         var mcisStatus = aMcis.status
    //         var mcisTargetStatus = aMcis.targetStatus
    //         var mcisTargetAction = aMcis.targetAction
    //         var mcisProviderNames = "";//MCIS에 사용 된 provider
    //         var totalVmCountOfMcis = 0;//MCIS의 VM 갯 수
    //         var mcisDispStatus = getMcisStatusDisp(mcisStatus);// 화면 표시용 status
    //         // mcis status
    //         try{
    //             // console.log(aMcis)
    //             if( mcisStatus != ""){// mcis status 가 없는 경우는 skip
    //                 if( mcisStatusCountMap.has(mcisDispStatus) ){
    //                     mcisStatusCountMap.set(mcisDispStatus, mcisStatusCountMap.get(mcisDispStatus) + 1)
    //                 }
    //                 totalMcisCnt++;
    //             }else{
    //                 continue;// status가 없으면 mcks 일 수 있으므로 mcis에서는 count 제외
    //             }
    //         }catch(e){
    //             console.log("mcis status error")
    //         }
    //
    //         // vm status
    //         try{
    //             var vmListOfMcis = aMcis.vm;// array
    //
    //
    //             var vmStatusCountMap = new Map();
    //             vmStatusCountMap.set("running", 0);
    //             vmStatusCountMap.set("stop", 0);  // partial 도 stop으로 보고있음.
    //             vmStatusCountMap.set("terminate", 0);
    //
    //             var vmCloudConnectionMap = new Map();
    //             console.log(vmListOfMcis)
    //             if (typeof vmListOfMcis !== 'undefined' && vmListOfMcis.length > 0) {
    //                 for(var vmIndex in vmListOfMcis){
    //                     var aVm = vmListOfMcis[vmIndex];
    //                     var vmDispStatus = getVmStatusDisp(aVm.status);
    //                     totalVmCountOfMcis++;
    //                     console.log(vmDispStatus)
    //                     if( vmStatusCountMap.has(vmDispStatus) ){
    //                         vmStatusCountMap.set(vmDispStatus, vmStatusCountMap.get(vmDispStatus) + 1)// mcis내 count
    //                         totalVmStatusCountMap.set(vmDispStatus, totalVmStatusCountMap.get(vmDispStatus) + 1)// 전체 count
    //                     }
    //                     totalServerCnt++;
    //
    //                     // connections
    //                     var location = aVm.location;
    //                     if( !isEmpty(location) ){
    //                         var cloudType = location.cloudType;
    //                         if( vmCloudConnectionMap.has(cloudType) ){
    //                             vmCloudConnectionMap.set(cloudType, vmCloudConnectionMap.get(cloudType) + 1)
    //                         }else{
    //                             vmCloudConnectionMap.set(cloudType, 0)
    //                         }
    //                     }
    //                 }
    //             }// end of vm list
    //
    //             // console.log(vmCloudConnectionMap);
    //             vmCloudConnectionMap.forEach((value, key) => {
    //                 mcisProviderNames += key + " ";
    //             });
    //             console.log("mcisProviderNames=" + mcisProviderNames);
    //         }catch(e){
    //             console.log("vm status error")
    //         }
    //
    //
    //
    //         // List of Mcis table
    //         try{
    //
    //
    //             // var displayMcisStatus = "";// icon_running, icon_stop, icon_terminate
    //             // if( mcisStatus.toLowerCase().indexOf("running") ){
    //             //     displayMcisStatus = "running"
    //             // }else if ( mcisStatus.toLowerCase().indexOf("suspend")){
    //             //     displayMcisStatus = "stop"
    //             // }else if ( mcisStatus.toLowerCase().indexOf("terminate")){
    //             //     displayMcisStatus = "terminate"
    //             // }else {
    //             //     displayMcisStatus = "terminate"
    //             // }
    //
    //
    //
    //             addMcis += '<tr onclick="clickListOfMcis(\'' + aMcis.id + '\', ' + mcisIndex +' );" id="server_info_tr_' + mcisIndex +'" item="' + aMcis.id +'|' + mcisIndex + '">'
    //
    //             addMcis += '<td class="overlay hidden td_left column-14percent" data-th="Status">'
    //             addMcis += '<img id="mcisInfo_mcisStatus_icon" src="/assets/img/contents/icon_' + mcisDispStatus +'.png" class="icon" alt=""/><div id="mcisInfo_mcisstatus_' + mcisIndex + '">' + mcisStatus + '</div><span class="ov off"></span>'
    //             addMcis += '</td>'
    //             addMcis += '<td class="btn_mtd ovm column-14percent" data-th="Name"><div id="mcisInfo_mcisName_' + mcisIndex + '">' + aMcis.name + '</div><span class="ov"></span></td>'
    //             addMcis += '<td class="overlay hidden column-14percent" data-th="Cloud Connection"><div id="mcisInfo_mcisProviderNames">' + mcisProviderNames + '</div></td>'
    //
    //             addMcis += '<td class="overlay hidden column-14percent" data-th="Total Infras"><div id="mcisInfo_totalVmCountOfMcis">' + totalVmCountOfMcis + '</div></td>'
    //
    //             addMcis += '<td class="overlay hidden column-14percent" data-th="# of Servers">'
    //             addMcis += '<span class="bar" ></span> <span title="running" id="mcisInfo_vmstatus_running_' + mcisIndex + '_' + vmIndex + '">' + vmStatusCountMap.get('running') + '</span>'
    //             addMcis += '<span class="bar" >/</span> <span title="stop" id="mcisInfo_vmstatus_stop' + mcisIndex + '_' + vmIndex + '">' + vmStatusCountMap.get('stop') + '</span>'
    //             addMcis += '<span class="bar" >/</span> <span title="terminate" id="mcisInfo_vmstatus_terminate' + mcisIndex + '_' + vmIndex + '">' + vmStatusCountMap.get('terminate') + '</span>'
    //             addMcis += '</td>'
    //
    //             addMcis += '<td class="overlay hidden" data-th="Description"><div id="mcisInfo_mcisDescription' + mcisIndex + '">' + aMcis.description + '</div></td>'
    //
    //             addMcis += '<td class="overlay hidden column-60px"  data-th="">'
    //             addMcis += '<input type="checkbox" name="chk" value="' + aMcis.id + '" id="td_ch_' + mcisIndex + '" title="" />'
    //             addMcis += '<label for="td_ch_' + mcisIndex + '"></label>'
    //             addMcis += '</td>'
    //
    //             // MCIS 기본정보 hidden : 클릭시 보여주기 위해
    //             addMcis += '<input type="hidden" name="mcisID" value="' + aMcis.id + '" id="mcisID' + mcisIndex + '"/>'
    //             addMcis += '<input type="hidden" name="mcisName" value="' + aMcis.name + '" id="mcisName' + mcisIndex + '"/>'
    //             addMcis += '<input type="hidden" name="mcisStatus" value="' + mcisStatus + '" id="mcisStatus' + mcisIndex + '"/>'
    //             addMcis += '<input type="hidden" name="mcisTargetStatus" value="' + mcisTargetStatus + '" id="mcisTargetStatus' + mcisIndex + '"/>'
    //             addMcis += '<input type="hidden" name="mcisTargetAction" value="' + mcisTargetAction + '" id="mcisTargetAction' + mcisIndex + '"/>'
    //             addMcis += '<input type="hidden" name="mcisDescription" value="' + aMcis.description + '" id="mcisDescription' + mcisIndex + '"/>'
    //             addMcis += '<input type="hidden" name="mcisCloudConnections" value="' + mcisProviderNames + '" id="mcisCloudConnections' + mcisIndex + '"/>'
    //             addMcis += '<input type="hidden" name="mcisVmTotalCount" value="' + totalServerCnt + '" id="mcisVmTotalCount' + mcisIndex + '"/>'
    //
    //
    //             addMcis += '<input type="hidden" name="m_exportFileName" id="m_exportFileName_' + mcisIndex + '" value="" />'
    //             addMcis += '<input type="hidden" name="m_mcisExportScript" id="m_mcisExportScript_' + mcisIndex + '" value="" />'
    //
    //             // export 용
    //
    //             if (typeof vmListOfMcis !== 'undefined' && vmListOfMcis.length > 0) {
    //                 for(var vmIndex in vmListOfMcis){
    //                     var aVm = vmListOfMcis[vmIndex];
    //                     addMcis += '<input type="hidden" name="vmID" id="mcisVmID_' + mcisIndex + '_' + vmIndex + '" value="' + aVm.id + '" />'
    //                     addMcis += '<input type="hidden" name="vmID" id="mcisVmName_' + mcisIndex + '_' + vmIndex + '" value="' + aVm.name + '" />'
    //                     addMcis += '<input type="hidden" name="vmStatus" id="mcisVmStatus_' + mcisIndex + '_' + vmIndex + '" value="' + aVm.status + '" />'
    //
    //                     // vm export 용 m_ 는 mcis의 첫글자 m
    //                     addMcis += '<input type="hidden" name="m_vmConnectionName" id="m_vmConnectionName_' + mcisIndex + '_' + vmIndex + '" value="' + aVm.connectionName + '" />'
    //                     addMcis += '<input type="hidden" name="m_vmDescription" id="m_vmDescription_' + mcisIndex + '_' + vmIndex + '" value="' + aVm.description + '" />'
    //                     addMcis += '<input type="hidden" name="m_vmImageId" id="m_vmImageId_' + mcisIndex + '_' + vmIndex + '" value="' + aVm.imageId + '" />'
    //                     addMcis += '<input type="hidden" name="m_vmLabel" id="m_vmLabel_' + mcisIndex + '_' + vmIndex + '" value="' + aVm.label + '" />'
    //                     addMcis += '<input type="hidden" name="m_vmName" id="m_vmName_' + mcisIndex + '_' + vmIndex + '" value="' + aVm.name + '" />'
    //                     addMcis += '<input type="hidden" name="m_vmSecurityGroupIds" id="m_vmSecurityGroupIds_' + mcisIndex + '_' + vmIndex + '" value="' + aVm.securityGroupIds + '" />'
    //                     addMcis += '<input type="hidden" name="m_vmSpecId" id="m_vmSpecId_' + mcisIndex + '_' + vmIndex + '" value="' + aVm.specId + '" />'
    //                     addMcis += '<input type="hidden" name="m_vmSshKeyId" id="m_vmSshKeyId_' + mcisIndex + '_' + vmIndex + '" value="' + aVm.sshKeyId + '" />'
    //                     addMcis += '<input type="hidden" name="m_vmSubnetId" id="m_vmSubnetId_' + mcisIndex + '_' + vmIndex + '" value="' + aVm.subnetId + '" />'
    //                     addMcis += '<input type="hidden" name="m_vmVnetId" id="m_vmVnetId_' + mcisIndex + '_' + vmIndex + '" value="' + aVm.vNetId + '" />'
    //                     addMcis += '<input type="hidden" name="m_vmGroupSize" id="m_vmGroupSize_' + mcisIndex + '_' + vmIndex + '" value="0" />' // vm 생성 시 동일한 것을 몇 개 만들 것인가이며 생성 param에만 있음.조회결과에는 없음.
    //                     addMcis += '<input type="hidden" name="m_vmUserAccount" id="m_vmUserAccount_' + mcisIndex + '_' + vmIndex + '" value="' + aVm.vmUserAccount + '" />'
    //                     addMcis += '<input type="hidden" name="m_vmUserPassword" id="m_vmUserPassword_' + mcisIndex + '_' + vmIndex + '" value="' + aVm.vmUserPassword + '" />'
    //
    //                     addMcis += '<input type="hidden" name="m_vmExportFileName" id="m_vmExportFileName_' + mcisIndex + '_' + vmIndex + '" value="" />'
    //                     addMcis += '<input type="hidden" name="m_vmExportScript" id="m_vmExportScript_' + mcisIndex + '_' + vmIndex + '" value="" />'
    //
    //                 }
    //             }
    //
    //             // mcis export
    //             addMcis += '<input type="hidden" name="m_mcisDescription" id="m_mcisDescription_' + mcisIndex + '" value="' + aMcis.description + '" />'
    //             addMcis += '<input type="hidden" name="m_mcisID" id="m_mcisID_' + mcisIndex + '" value="' + aMcis.id + '" />'
    //             addMcis += '<input type="hidden" name="m_mcisInstallMonAgent" id="m_mcisInstallMonAgent_' + mcisIndex + '" value="' + aMcis.installMonAgent + '" />'
    //             addMcis += '<input type="hidden" name="m_mcisLabel" id="m_mcisLabel_' + mcisIndex + '" value="' + aMcis.label + '" />'
    //             addMcis += '<input type="hidden" name="m_mcisName" id="m_mcisName_' + mcisIndex + '" value="' + aMcis.name + '" />'
    //             addMcis += '<input type="hidden" name="m_mcisStatus" id="m_mcisStatus_' + mcisIndex + '" value="' + aMcis.status + '" />'
    //
    //
    //             addMcis += '<input type="hidden" name="m_mcisExportFileName" id="m_exportFileName_' + mcisIndex + '" value="" />'
    //             addMcis += '<input type="hidden" name="m_mcisExportScript" id="m_mcisExportScript_' + mcisIndex + '" value="" />'
    //
    //
    //
    //             addMcis += '</tr>'
    //
    //         }catch(e){
    //             console.log("list of mcis error")
    //             console.log(e)
    //         }
    //     }// end of mcis loop
    //     $("#mcisList").empty();
    //     $("#mcisList").append(addMcis);
    //
    //     $("#total_mcis").text(totalMcisCnt);
    //     $("#mcis_status_running").text(mcisStatusCountMap.get("running"));
    //     $("#mcis_status_stopped").text(mcisStatusCountMap.get("stop"));
    //     $("#mcis_status_terminated").text(mcisStatusCountMap.get("terminate"));
    //
    //     $("#total_vm").text(totalServerCnt);
    //     $("#vm_status_running").text(totalVmStatusCountMap.get("running"));
    //     $("#vm_status_stopped").text(totalVmStatusCountMap.get("stop"));
    //     $("#vm_status_terminated").text(totalVmStatusCountMap.get("terminate"));
    // }else{
    //     var addMcis = "";
    //     addMcis += '<tr>'
    //     addMcis += '<td class="overlay hidden" data-th="" colspan="8">No Data</td>'
    //     addMcis += '</tr>'
    //     $("#mcisList").empty();
    //     $("#mcisList").append(addMcis);
    // }
}

// 조회 실패시.
function getMcisListCallbackFail(caller, error) {
    // List table에 no data 표시? 또는 조회 오류를 표시?
    var addMcis = "";
    addMcis += '<tr>'
    addMcis += '<td class="overlay hidden" data-th="" colspan="8">No Data</td>'
    addMcis += '</tr>'
    $("#mcisList").empty();
    $("#mcisList").append(addMcis);
}


function setMcisData(mcisInfo, mcisIndex) {
    var mcisStatus = mcisInfo.status
    var mcisTargetStatus = mcisInfo.targetStatus
    var mcisTargetAction = mcisInfo.targetAction
    var mcisProviderNames = "";//MCIS에 사용 된 provider
    var totalVmCountOfMcis = 0;//MCIS의 VM 갯 수
    var mcisDispStatus = getMcisStatusDisp(mcisStatus);// 화면 표시용 status
    var mcisStatusImg = "/assets/img/contents/icon_" + mcisDispStatus + ".png"

    var totalServerCnt = 0;
    var totalVmStatusCountMap = new Map();
    totalVmStatusCountMap.set("running", 0);
    totalVmStatusCountMap.set("stop", 0);  // partial 도 stop으로 보고있음.
    totalVmStatusCountMap.set("terminate", 0);

    var vmListOfMcis = mcisInfo.vm;// array


    var vmStatusCountMap = new Map();
    vmStatusCountMap.set("running", 0);
    vmStatusCountMap.set("stop", 0);  // partial 도 stop으로 보고있음.
    vmStatusCountMap.set("terminate", 0);

    var vmCloudConnectionMap = new Map();
    console.log(vmListOfMcis)
    if (typeof vmListOfMcis !== 'undefined' && vmListOfMcis.length > 0) {
        for (var vmIndex in vmListOfMcis) {
            var aVm = vmListOfMcis[vmIndex];
            var vmDispStatus = getVmStatusDisp(aVm.status);
            totalVmCountOfMcis++;
            console.log(vmDispStatus)
            if (vmStatusCountMap.has(vmDispStatus)) {
                vmStatusCountMap.set(vmDispStatus, vmStatusCountMap.get(vmDispStatus) + 1)// mcis내 count
                totalVmStatusCountMap.set(vmDispStatus, totalVmStatusCountMap.get(vmDispStatus) + 1)// 전체 count
            }
            totalServerCnt++;

            // connections
            var location = aVm.location;
            if (!isEmpty(location)) {
                var cloudType = location.cloudType;
                if (vmCloudConnectionMap.has(cloudType)) {
                    vmCloudConnectionMap.set(cloudType, vmCloudConnectionMap.get(cloudType) + 1)
                } else {
                    vmCloudConnectionMap.set(cloudType, 0)
                }
            }
        }
    }// end of vm list

    vmCloudConnectionMap.forEach((value, key) => {
        mcisProviderNames += key + " ";
    });
    console.log("mcisProviderNames=" + mcisProviderNames)
    // table
    $("#mcisInfo_mcisStatus_icon_" + mcisIndex).attr("src", mcisStatusImg);
    $("#mcisInfo_mcisName_" + mcisIndex).text(mcisInfo.name);
    $("#mcisInfo_mcisProviderNames_" + mcisIndex).text(mcisProviderNames);
    $("#mcisInfo_totalVmCountOfMcis_" + mcisIndex).text(totalVmCountOfMcis)

    $("#mcisInfo_vmstatus_running_" + mcisIndex + '_' + vmIndex).text(vmStatusCountMap.get('running'))
    $("#mcisInfo_vmstatus_stop_" + mcisIndex + '_' + vmIndex).text(vmStatusCountMap.get('stop'))
    $("#mcisInfo_vmstatus_terminate_" + mcisIndex + '_' + vmIndex).text(vmStatusCountMap.get('terminate'))

    $("#mcisInfo_mcisDescription_" + mcisIndex).text(mcisInfo.description)

    $("#mcisID" + mcisIndex).text(mcisInfo.id)
    $("#mcisName" + mcisIndex).text(mcisInfo.name)
    $("#mcisStatus" + mcisIndex).text(mcisStatus)
    $("#mcisTargetStatus" + mcisIndex).text(mcisTargetStatus)
    $("#mcisTargetAction" + mcisIndex).text(mcisTargetAction)
    $("#mcisDescription" + mcisIndex).text(mcisInfo.description)
    $("#mcisCloudConnections" + mcisIndex).text(mcisProviderNames)
    $("#mcisVmTotalCount" + mcisIndex).text(totalServerCnt)
    console.log("all set " + mcisStatusImg)

    setMcisServerInfoBox(mcisIndex, vmIndex)
}
///////////// MCIS Handling //////////////

// 등록 form으로 이동
function createNewMcis() {// Manage_MCIS_Life_Cycle_popup.html
    var targetUrl = "/operation/manage" + "/mcismng/regform"
    // location.href = "/Manage/MCIS/reg"
    // $('#loadingContainer').show();
    // location.href = url;
    changePage(targetUrl)
}

function changeLifeCycle(type) {
    var checked_nothing = 0;
    var mcisIndex = 0;
    var isMcks = false;
    $("[id^='td_ch_']").each(function () {

        if ($(this).is(":checked")) {
            checked_nothing++;
            aMcisData = totalMcisListObj[mcisIndex];
            console.log(aMcisData);
            var systemLabel = aMcisData.systemLabel;
            if (systemLabel) {
                systemLabel = systemLabel.toLowerCase();
                if (systemLabel.indexOf("mcks") > -1) {
                    isMcks = true;
                    return false;
                }
            }
        } else {
            console.log("checked nothing")
        }
        mcisIndex++;
    })
    if (isMcks) {
        commonAlert("MCKS life cycle cannot be changed");
        return;
    }

    console.log("checked_nothing " + checked_nothing)
    if (checked_nothing == 0) {
        commonAlert("Please Select MCIS!!")
        return;
    }
    commonConfirmOpen(type);
}
// MCIS 제어 : 선택한 VM의 상태 변경
// changeLifeCycle -> callMcisLifeCycle -> util.callMcisLifeCycle -> callbackMcisLifeCycle 순으로 호출 됨
function callMcisLifeCycle(type) {
    var checked_nothing = 0;
    $("[id^='td_ch_']").each(function () {

        if ($(this).is(":checked")) {
            checked_nothing++;
            console.log("checked")
            var mcisID = $(this).val();
            mcisLifeCycle(mcisID, type);
        } else {
            console.log("checked nothing")

        }
    })
    if (checked_nothing == 0) {
        commonAlert("Please Select MCIS!!")
        return;
    }
}

// McisLifeCycle을 호출 한 뒤 return값 처리
function callbackMcisLifeCycle(resultStatus, resultData, type) {
    var message = "MCIS " + type + " requested!"
    if (resultStatus == 200 || resultStatus == 201) {
        // commonAlert(message);
        console.log("callbackMcisLifeCycle" + message);
        commonResultAlert(message);
        //location.reload();//완료 후 페이지를 reload -> 해당 mcis만 reload

        // 해당 mcis 조회
        // 상태 count 재설정
    } else {

        commonAlert("MCIS " + type + " failed!");
    }
}

// mcis삭제 전 mcks 포함여부 체크
function deleteCheckMCIS(type) {
    var checkedCount = 0;
    var mcisID = "";
    var mcisIndex = 0;
    var isMcks = false;

    $("[id^='td_ch_']").each(function () {

        if ($(this).is(":checked")) {
            checkedCount++;
            console.log("checked")
            mcisID = $(this).val();
            aMcisData = totalMcisListObj[mcisIndex];
            console.log(aMcisData);
            var systemLabel = aMcisData.systemLabel;
            if (systemLabel) {
                systemLabel = systemLabel.toLowerCase();
                if (systemLabel.indexOf("mcks") > -1) {
                    isMcks = true;
                    return false;
                }
            }
        } else {
            console.log("checked nothing")

        }
        mcisIndex++;
    });
    if (checkedCount == 0) {
        commonAlert("Please Select MCIS!!")
        return;
    } else if (checkedCount > 1) {
        commonAlert("Please Select One MCIS at a time")
        return;
    }

    if (isMcks) {
        commonAlert("MCKS cannot be deleted");
        return;
    }
    commonConfirmOpen(type)
}
// list에 선택된 MCIS 제거. 1개씩
function deleteMCIS() {
    var checkedCount = 0;
    var mcisID = "";
    $("[id^='td_ch_']").each(function () {

        if ($(this).is(":checked")) {
            checkedCount++;
            console.log("checked")
            mcisID = $(this).val();
            // 여러개를 지울 때 호출하는 함수를 만들어 여기에서 호출
        } else {
            console.log("checked nothing")

        }
    })

    if (checkedCount == 0) {
        commonAlert("Please Select MCIS!!")
        return;
    } else if (checkedCount > 1) {
        commonAlert("Please Select One MCIS at a time")
        return;
    }

    // TODO : 삭제 호출부분 function으로 뺼까?
    var url = "/operation/manages/mcismng/" + mcisID + "?option=force";
    axios.delete(url, {})
        .then(result => {
            console.log("get  Data : ", result.data);

            var statusCode = result.data.status;
            var message = result.data.message;

            if (statusCode != 200 && statusCode != 201) {
                commonAlert(message + "(" + statusCode + ")");
                return;
            } else {
                // commonAlert(message);
                // TODO : MCIS List 조회
                // location.reload();
                commonResultAlert(message);
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
function addNewVirtualMachine() {
    var mcis_id = $("#mcis_id").val()
    var mcis_name = $("#mcis_name").val()
    // location.href = "/Manage/MCIS/reg/"+mcis_id+"/"+mcis_name
    location.href = "/operation/manages/mcismng/regform/" + mcis_id + "/" + mcis_name;
}

function vmLifeCycle(type) {
    var mcisID = $("#mcis_id").val();
    var vmID = $("#vm_id").val();
    var vmName = $("#vm_name").val();


    // MCIS에서 MCKS관련 리소스는 handling하지 않음.
    var aMcis = new Object();
    var isMcks = false;
    for (var mcisIndex in totalMcisListObj) {
        var tempMcis = totalMcisListObj[mcisIndex]
        if (mcisID == tempMcis.id) {
            aMcis = tempMcis;
            var systemLabel = tempMcis.systemLabel;
            if (systemLabel) {
                systemLabel = systemLabel.toLowerCase();
                if (systemLabel.indexOf("mcks")) {
                    // commonAlert("MCKS life cycle cannot be changed")
                    // return;
                    isMcks = true;
                    break;
                }
            }
        }
    }// end of mcis loop

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
    if (!mcisID) {
        commonAlert("Please Select MCIS!!")
        return;
    }
    if (!vmID) {
        commonAlert("Please Select VM!!")
        return;
    }
    if (isMcks) {

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

    console.log("life cycle3 url : ", url);

    var message = vmName + " " + type + " requested!."
    var obj = {
        mcisID: mcisID,
        vmID: vmID,
        lifeCycleType: type
    }
    axios.post(url, obj, {
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
    }).then(result => {
        var status = result.status

        console.log("life cycle result : ", result)
        var data = result.data
        console.log("result Message : ", data.message)
        if (status == 200 || status == 201) {
            commonResultAlert(message);
            // commonAlert(message);
            // location.reload();// TODO 일단은 Reaoad : 해당 영역(MCIS의 VM들 status 조회)를 refresh할 수 있는 기능 필요
            //show_mcis(mcis_url,"");
        }
    })
}

function getVmOfMcisDataCallbackSuccess(caller, vmData, mcisID, vmID) {

}
///////////// VM Handling end ///////////



const config_arr = new Array();

// refresh : mcis 및 vm정보조회
// 각 mcis 별 vmstatus 목록


// List Of MCIS 클릭 시
// mcis 테이블의 선택한 row 강조( on )
// 해당 MCIS의 VM 상태목록 보여주는 함수 호출
function clickListOfMcis(mcisID) {
    console.log("click view mcis id :", mcisID)
    if (mcisID != "") {
        // MCIS Info 에 mcis id 표시
        $("#mcis_id").val(mcisID);
        $("#selected_mcis_id").val(mcisID);
        // $("#selected_mcis_index").val(mcisIndex);

        //클릭 시 mcisinfo로 포커스 이동


        getCommonMcisData("refreshmcisdata", mcisID)

        // MCIS Info area set
        //showServerListAndStatusArea(mcisID,mcisIndex);
        //displayMcisInfoArea(totalMcisListObj[mcisIndex]);


        //makeMcisScript(index);// export를 위한 script 준비 -> Export 실행할 때 가져오는 것으로 변경( MCIS정보는 option=simple로 가져오므로)
    }
}

const mcisInfoDataList = new Array()// test_arr : mcisInfo 1개 전체, pageLoad될 때, refresh 할때 data를 set. mcis클릭시 상세내용 보여주기용 조회

// MCIS Info area 안의 Server List / Status 내용 표시
// 해당 MCIS의 모든 VM 표시
// TODO : 클릭했을 때 서버에서 조회하는것으로 변경할 것.
function showServerListAndStatusArea(mcis_id, mcisIndex) {

    var mcisID = $("#mcisID" + mcisIndex).val();
    var mcisName = $("#mcisName" + mcisIndex).val();
    var mcisDescription = $("#mcisDescription" + mcisIndex).val();
    var mcisStatus = $("#mcisStatus" + mcisIndex).val();
    var mcisTargetStatus = $("#mcisTargetStatus" + mcisIndex).val();
    var mcisTargetAction = $("#mcisTargetAction" + mcisIndex).val();
    var mcisDispStatus = getMcisStatusDisp(mcisStatus);
    var mcisCloudConnections = $("#mcisCloudConnections" + mcisIndex).val();
    var vmTotalCountOfMcis = $("#mcisVmTotalCount" + mcisIndex).val();
    var vms = $("#mcisVmStatusList" + mcisIndex).val();

    $("#mcis_info_txt").text("[ " + mcisName + " ]");
    $("#mcis_server_info_status").empty();
    $("#mcis_server_info_status").append('<strong>Server List / Status</strong>  <span class="stxt">[ ' + mcisName + ' ]</span>  Server(' + vmTotalCountOfMcis + ')')

    //
    $("#mcis_info_name").val(mcisName + " / " + mcisID);
    $("#mcis_info_id").val(mcisID);
    $("#mcis_info_description").val(mcisDescription);
    $("#mcis_info_targetStatus").val(mcisTargetStatus);
    $("#mcis_info_targetAction").val(mcisTargetAction);
    $("#mcis_info_cloud_connection").val(mcisCloudConnections)    //

    $("#mcis_name").val(mcisName)

    var mcis_badge = "";
    var mcisStatusIcon = getMcisStatusIcon(mcisDispStatus);
    // var mcisStatusIcon = "";
    // if(mcisDispStatus == "running"){ mcisStatusIcon = "icon_running_db.png"
    // }else if(mcisDispStatus == "include" ){ mcisStatusIcon = "icon_stop_db.png"
    // }else if(mcisDispStatus == "suspended"){mcisStatusIcon = "icon_stop_db.png"
    // }else if(mcisDispStatus == "terminate"){mcisStatusIcon = "icon_terminate_db.png"
    // }else{
    //     mcisStatusIcon = "icon_stop_db.png"
    // }

    // server info
    mcis_badge = '<img src="/assets/img/contents/' + mcisStatusIcon + '" alt=""/> '
    $("#service_status_icon").empty();
    $("#service_status_icon").append(mcis_badge)

    setMcisServerInfoBox(mcisIndex, -1);
    // var vm_badge = "";
    // $("[id^='mcisVmID_']").each(function(){
    //     var mcisVm = $(this).attr("id").split("_")
    //     thisMcisIndex = mcisVm[1]
    //     vmIndexOfMcis = mcisVm[2]
    //
    //     if( thisMcisIndex == mcisIndex){
    //
    //         var vmID = $("#mcisVmID_" + mcisIndex + "_" + vmIndexOfMcis).val();
    //         var vmName = $("#mcisVmName_" + mcisIndex + "_" + vmIndexOfMcis).val();
    //         var vmStatus = $("#mcisVmStatus_" + mcisIndex + "_" + vmIndexOfMcis).val();
    //         vmStatus = vmStatus.toLowerCase();
    //         var vmDispStatus = getVmStatusDisp(vmStatus);
    //         var vmStatusIcon ="bgbox_g";
    //         if(vmDispStatus == "running"){
    //             vmStatusIcon ="bgbox_b"
    //         }else if(vmDispStatus == "include" ){
    //             vmStatusIcon ="bgbox_g"
    //         }else if(vmDispStatus == "suspended"){
    //             vmStatusIcon ="bgbox_g"
    //         }else if(vmDispStatus == "terminated"){
    //             vmStatusIcon ="bgbox_r"
    //         }else{
    //             vmStatusIcon ="bgbox_r"
    //         }
    //         vm_badge += '<li class="sel_cr ' + vmStatusIcon + '"><a href="javascript:void(0);" onclick="vmDetailInfo(\''+mcisID+'\',\''+mcisName+'\',\''+vmID+'\')"><span class="txt">'+vmName+'</span></a></li>';
    //     }
    // });
    // // console.log(vm_badge);
    // $("#mcis_server_info_box").empty();
    // $("#mcis_server_info_box").append(vm_badge);

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

}

// 현재 mcisIndex와 선택한 vm의 index(미선택은 -1)
function setMcisServerInfoBox(mcisIndex, vmIndex) {
    console.log(" mcisIndex=" + mcisIndex + ", vmIndex=" + vmIndex)
    if (Number(vmIndex) == -1) {
        var mcisID = $("#mcisID" + mcisIndex).val();
        var mcisName = $("#mcisName" + mcisIndex).val();
        var vm_badge = "";
        $("[id^='mcisVmID_']").each(function () {
            var mcisVm = $(this).attr("id").split("_")
            thisMcisIndex = mcisVm[1]
            vmIndexOfMcis = mcisVm[2]

            if (thisMcisIndex == mcisIndex) {
                var vmID = $("#mcisVmID_" + mcisIndex + "_" + vmIndexOfMcis).val();
                var vmName = $("#mcisVmName_" + mcisIndex + "_" + vmIndexOfMcis).val();
                var vmStatus = $("#mcisVmStatus_" + mcisIndex + "_" + vmIndexOfMcis).val();
                vmStatus = vmStatus.toLowerCase();
                var vmDispStatus = getVmStatusDisp(vmStatus);
                var vmStatusIcon = getVmStatusIcon(vmDispStatus);
                // var vmStatusIcon = "bgbox_g";
                // if (vmDispStatus == "running") {
                //     vmStatusIcon = "bgbox_b"
                // } else if (vmDispStatus == "include") {
                //     vmStatusIcon = "bgbox_g"
                // } else if (vmDispStatus == "suspended") {
                //     vmStatusIcon = "bgbox_g"
                // } else if (vmDispStatus == "terminated") {
                //     vmStatusIcon = "bgbox_r"
                // } else {
                //     vmStatusIcon = "bgbox_r"
                // }
                vm_badge += '<li class="sel_cr ' + vmStatusIcon + '" id="vmOfMcis_\' + mcisIndex + \'_\' + vmIndex + "><a href="javascript:void(0);" onclick="vmDetailInfo(\'' + mcisID + '\',\'' + mcisName + '\',\'' + vmID + '\')"><span class="txt">' + vmName + '</span></a></li>';
            }
        });
        // console.log(vm_badge);
        $("#mcis_server_info_box").empty();
        $("#mcis_server_info_box").append(vm_badge);

    } else {
        console.log(" change class")
        var mcisID = $("#mcisID" + mcisIndex).val();
        var mcisName = $("#mcisName" + mcisIndex).val();
        var vmID = $("#mcisVmID_" + mcisIndex + "_" + vmIndex).val();
        var vmName = $("#mcisVmName_" + mcisIndex + "_" + vmIndexOfMcis).val();
        var vmLi = $("#vmOfMcis_" + mcisIndex + "_" + vmIndex);
        var vmStatus = $("#mcisVmStatus_" + mcisIndex + "_" + vmIndex).val();
        var vmDispStatus = getVmStatusDisp(vmStatus);
        var vmStatusIcon = getVmStatusIcon(vmDispStatus);
        console.log("li attrt class " + vmStatusIcon)
        vmLi.attr('class', 'sel_cr ' + vmStatusIcon);

        var vm_badge = '<li class="sel_cr ' + vmStatusIcon + '" id="vmOfMcis_\' + mcisIndex + \'_\' + vmIndex + "><a href="javascript:void(0);" onclick="vmDetailInfo(\'' + mcisID + '\',\'' + mcisName + '\',\'' + vmID + '\')"><span class="txt">' + vmName + '</span></a></li>';
        $("#mcis_server_info_box").empty();
        $("#mcis_server_info_box").append(vm_badge);
    }
}

// 특정 MCIS 만 refresh
// vmList의 refresh버튼 클릭 시
function refreshMcisStatusData() {
    var mcisID = $("#mcis_id").val();
    console.log("refreshMcisData=" + mcisID)
    if (mcisID) {
        getCommonMcisStatusData("refreshmcisdata", mcisID);//
    }
}

// 1개 MCIS 상태 변경
function getCommonMcisStatusDataCallbackSuccess(caller, mcisStatusInfo) {
    console.log("caller " + caller)
    console.log(mcisStatusInfo)
    var mcisID = mcisStatusInfo.id;

    // 받아온 값이 없는 경우는 mcis삭제이므로 table부터 다시그려야 함.
    if (!mcisID) {
        // mcis목록 조회
        getCommonMcisList("mcismngready", true, "", "simple")
        $(".server_status").removeClass("view");
        $(".server_info").removeClass("active");
        return;
    }

    var aMcis = new Object();
    for (var mcisIndex in totalMcisListObj) {
        var tempMcis = totalMcisListObj[mcisIndex]
        if (mcisID == tempMcis.id) {
            aMcis = tempMcis;
            break;
        }
    }// end of mcis loop

    // set mcis
    console.log(aMcis.status + " : " + mcisStatusInfo.status)
    aMcis.status = mcisStatusInfo.status;
    aMcis.targetAction = mcisStatusInfo.targetAction;
    aMcis.targetStatus = mcisStatusInfo.targetStatus;

    // set vm list of mcis
    vmList = aMcis.vm;
    if (vmList) {
        aMcis.vm = mcisStatusInfo.vm
    } else {
        var vmList = mcisStatusInfo.vm;
        // 가져온 상태값
        vmList.forEach(function (vmItem, vmIndex) {
            var isExist = false;
            // 기존 상태값
            for (var aIndex in mcisStatusInfo.vm) {
                var aVm = aMcis.vm[aIndex];
                if (vmItem.id == aVm.id) {
                    aVm.status = vmItem.status;
                    aVm.installMonAgent = vmItem.installMonAgent;
                    aVm.systemLabel = vmItem.systemLabel
                    aVm.targetAction = vmItem.targetAction;
                    aVm.targetStatus = vmItem.targetStatus;
                    break;
                }
            }
        })

    }

    updateMcisData(aMcis, mcisID);
}
// MCIS의 상태값만 변경
function refreshMcisStatus(mcisID) {
    //
}

// 특정 Vm의 상태만 변경
function refreshVmStatus(mcisID, vmID) {

}

// Mcis Info영역 변경

// Server Info 변경 : 선택된 vm재조회, 상단의  MCIS Info의 해당 vm status변경 반영, 상단의 MCIS List의 요약에서 status 변경
function refreshVmList() {
    var mcisID = $("#selected_mcis_id").val();
    console.log("refreshVmList " + mcisID);
    getCommonMcisData("refreshmcisdata", mcisID)

    // MCIS Data를 가져와 Set한 이후 MCIS List table에 반영, MCIS Info에 반영

}

// refresh button 클릭 시 해당 vm 정보 조회
function refreshVmDetailInfo() {
    var mcisID = $("#selected_mcis_id").val();
    var vmID = $("#selected_vm_id").val();
    vmDetailInfo(mcisID, '', vmID)
}

// VM 목록에서 VM 클릭시 해당 VM의 상세정보
function vmDetailInfo(mcisID, mcisName, vmID) {
    // var url = "/operation/manages/mcismng/" + mcisID + "/vm/" + vmID
    // axios.get(url, {})
    //     .then(result => {
    //         console.log("get  Data : ", result.data);

    // var statusCode = result.data.status;
    // var message = result.data.message;
    //
    // if (statusCode != 200 && statusCode != 201) {
    //     commonAlert(message + "(" + statusCode + ")");
    //     return;
    // }
    console.log("vmDetailInfo " + vmID)
    var aMcis = new Object();
    for (var mcisIndex in totalMcisListObj) {
        var tempMcis = totalMcisListObj[mcisIndex]
        if (mcisID == tempMcis.id) {
            aMcis = tempMcis;
            break;
        }
    }// end of mcis loop
    console.log(aMcis);
    var vmList = aMcis.vm;
    var vmExist = false
    var data = new Object();
    for (var vmIndex in vmList) {
        var aVm = vmList[vmIndex]
        if (vmID == aVm.id) {
            //aVm = vmData;
            data = aVm;
            vmExist = true
            break;
        }
    }
    if (!vmExist) {
        console.log("vm is not exist");
        console.log(vmList)
    }
    // var data = result.data.VmInfo;
    // var connectionConfig = result.data.ConnectionConfigInfo;
    console.log("selected Vm");
    console.log(data);
    var vmId = data.id;
    var vmName = data.name;
    var vmStatus = data.status;
    var vmDescription = data.description;
    var vmPublicIp = data.publicIP == undefined ? "" : data.publicIP;
    var vmSshKeyID = data.sshKeyId;
    //vm info
    $("#vm_id").val(vmId);
    $("#vm_name").val(vmName);
    console.log("vm_id " + vmId + ", vm_name " + vmName + ", vm_SshKeyID " + vmSshKeyID)

    $("#manage_mcis_popup_vm_id").val(vmId)
    $("#manage_mcis_popup_mcis_id").val(mcisID)
    $("#manage_mcis_popup_sshkey_name").val(vmSshKeyID)

    $("#server_info_text").text('[' + vmName + '/' + mcisName + ']')
    $("#server_detail_info_text").text('[' + vmName + '/' + mcisName + ']')


    var vmBadge = ""
    var vmStatusIcon = getVmStatusIcon(vmStatus)
    // var vmStatusIcon = "icon_running_db.png";
    // if(vmStatus == "Running"){
    //     vmStatusIcon = "icon_running_db.png";
    // }else if(vmStatus == "include"){
    //     vmStatusIcon = "icon_stop_db.png";
    // }else if(vmStatus == "Suspended"){
    //     vmStatusIcon = "icon_stop_db.png";
    // }else if(vmStatus == "Terminated"){
    //     vmStatusIcon = "icon_terminate_db.png";
    // }else{
    //     vmStatusIcon = "icon_stop_db.png";
    // }
    //var vmBadge = '<img src="/assets/img/contents/' + vmStatusIcon +'" alt="' + vmStatus + '"/>'

    $("#server_detail_view_server_status").val(vmStatus);
    // $("#server_info_status_img").empty()
    // $("#server_info_status_img").append(vmBadge)
    $("#server_info_status_icon_img\n").attr("src", vmStatusIcon);

    $("#server_info_name").val(vmName + "/" + vmID)
    $("#server_info_desc").val(vmDescription)

    // ip information
    $("#server_info_public_ip").val(vmPublicIp)
    $("#server_detail_info_public_ip_text").text("Public IP : " + vmPublicIp)
    $("#server_info_public_dns").val(data.publicDNS)
    $("#server_info_private_ip").val(data.privateIP)
    $("#server_info_private_dns").val(data.privateDNS)

    $("#server_detail_view_public_ip").val(vmPublicIp)
    $("#server_detail_view_public_dns").val(data.publicDNS)
    $("#server_detail_view_private_ip").val(data.privateIP)
    $("#server_detail_view_private_dns").val(data.privateDNS)

    $("#manage_mcis_popup_public_ip").val(vmPublicIp)

    // connection tab
    var locationInfo = data.location;
    var cloudType = locationInfo.cloudType;
    $("#server_info_csp_icon").empty()
    $("#server_info_csp_icon").append('<img src="/assets/img/contents/img_logo_' + cloudType + '.png" alt=""/>')
    $("#server_connection_view_csp").val(cloudType)
    $("#manage_mcis_popup_csp").val(cloudType)


    var latitude = locationInfo.latitude;
    var longitude = locationInfo.longitude;
    var briefAddr = locationInfo.briefAddr;
    var nativeRegion = locationInfo.nativeRegion;

    if (locationInfo) {
        $("#server_location_latitude").val(latitude)
        $("#server_location_longitude").val(longitude)

        // // 지도 표시
        // var pointInfo = new Object();
        // pointInfo.id = "1"
        // pointInfo.name = nativeRegion;
        // pointInfo.cloudType = cloudType;//server_connection_view_csp
        // pointInfo.latitude = latitude;//server_location_latitude
        // pointInfo.longitude = longitude//server_location_longitude
        // pointInfo.markerIndex = 1
        // setMap(pointInfo);
    }
    // region zone locate

    //var region = data.region.region
    var region = data.region.region
    // var zone = data.region.zone
    var zone = data.region.zone
    // console.log(vmDetail.iid);
    $("#server_info_region").val(briefAddr + ":" + region)
    $("#server_info_zone").val(zone)
    // $("#server_info_cspVMID").val("cspVMID : " + vmDetail.iid.nameId)// vmDetail에 있는 항목으로 vmDetail에서 처리

    $("#server_detail_view_region").val(briefAddr + " : " + region)
    $("#server_detail_view_zone").val(zone)

    $("#server_connection_view_region").val(briefAddr + "(" + region + ")")
    $("#server_connection_view_zone").val(zone)

    // connection name
    var connectionName = data.connectionName;
    $("#server_info_connection_name").val(connectionName)
    $("#server_connection_view_connection_name").val(connectionName)

    // credential and driver info : 가져오지 않음으로
    // console.log("connectionConfig : ", connectionConfig)
    // if (connectionConfig) {
    //     var credentialName = connectionConfig.CredentialName
    //     var driverName = connectionConfig.DriverName
    //     $("#server_connection_view_credential_name").val(credentialName)
    //     $("#server_connection_view_driver_name").val(driverName)
    // }
    for (cIndex in TotalCloudConnectionList) { // TODO : connection의 driver와 connection set.
        var connInfo = TotalCloudConnectionList[cIndex];
        if (connectionName == connInfo.ConfigName) {
            $("#server_connection_view_credential_name").val(connInfo.CredentialName)
            $("#server_connection_view_driver_name").val(connInfo.DriverName)
            break;
        }
    }

    //////vm detail tab////
    var vmDetail = data.cspViewVmDetail;
    if (vmDetail) {
        //    //cspvmdetail
        // var vmDetailKeyValueList = vmDetail.KeyValueList
        var vmDetailKeyValueList = vmDetail.keyValueList
        var architecture = "";
        if (vmDetailKeyValueList) {
            for (var keyIndex in vmDetailKeyValueList) {
                if (vmDetailKeyValueList[keyIndex].key == "Architecture") {// ?? 이게 뭐지?
                    architecture = vmDetailKeyValueList[keyIndex].value
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

        // server id / system id
        $("#server_detail_view_server_id").val(data.id)
        // systemid 를 사용할 경우 아래 꺼 사용
        $("#server_detail_view_server_id").val(vmDetail.iid.systemId)

        // image id
        var imageIId = vmDetail.imageIId.nameId
        var imageId = data.imageId
        getCommonVmImageInfo('mcisvmdetail', imageId) //
        $("#server_detail_view_image_id").text(imageId + "(" + imageIId + ")")

        //vpc subnet
        var vpcId = vmDetail.vpcIID.nameId
        var vpcSystemId = vmDetail.vpcIID.systemId
        var subnetId = vmDetail.subnetIID.nameId
        var subnetSystemId = vmDetail.subnetIID.systemId
        var eth = vmDetail.networkInterface
        $("#server_detail_view_vpc_id").text(vpcId + "(" + vpcSystemId + ")")
        // set_vmVPCInfo(vpcId, subnetId);

        $("#server_detail_view_subnet_id").text(subnetId + "(" + subnetSystemId + ")")
        $("#server_detail_view_eth").val(eth)

        // user account
        $("#server_detail_view_access_id_pass").val(vmDetail.vmuserId + "/ *** ")
        $("#server_detail_view_user_id_pass").val(data.vmUserAccount + "/ *** ")
        $("#manage_mcis_popup_user_name").val(data.vmUserAccount)


        var append_sg = ''

        var sg_arr = vmDetail.securityGroupIIds
        console.log(sg_arr);
        if (sg_arr) {
            sg_arr.map((item, index) => {
                console.log("item index = " + index)
                append_sg += '<a href="javascript:void(0);" onclick="getCommonVmSecurityGroupInfo(\'mcisvm\',\'' + item.nameId + '\');"title="' + item.nameId + '" >' + item.nameId + '</a> '
            })
        }
        console.log("append sg : ", append_sg)

        $("#server_detail_view_security_group").empty()
        $("#server_detail_view_security_group").append(append_sg);

        $("#server_detail_view_keypair_name").val(vmDetail.keyPairIId.nameId)

        console.log(vmDetail.iid);
        $("#server_info_cspVMID").val("cspVMID : " + vmDetail.iid.nameId)

        // ... TODO : 우선 제어명령부터 처리. 나중에 해당항목 mapping하여 확인
    }// end of vm detail
    ////// vm connection tab //////


    $("#selected_mcis_id").val(mcisID);
    $("#selected_vm_id").val(vmID);
    var installMonAgent = data.monAgentStatus;
    console.log("install mon agent : ", installMonAgent)
    if (installMonAgent == "installed") {
        var isWorking = checkDragonFlyMonitoringAgent(mcisID, vmID);
        if (isWorking) {
            $("#mcis_detail_info_check_monitoring").prop("checked", true)
            $("#mcis_detail_info_check_monitoring").attr("disabled", true)
        } else {
            $("#mcis_detail_info_check_monitoring").prop("checked", false)
            $("#mcis_detail_info_check_monitoring").attr("disabled", false)
        }
    } else {
        $("#mcis_detail_info_check_monitoring").prop("checked", false)
        $("#mcis_detail_info_check_monitoring").attr("disabled", false)
    }

    ////// vm mornitoring tab 으로 이동 //////
    // install Mon agent
    // showVmMonitoring(mcisID,vmID)

    // script 생성
    console.log("vmScriptExport start")
    // 위 값으로 mcisIndex, vmIndex 를 찾자
    // var mcisIndex = 0;
    // var vmIndex = 0;
    // $("[id^='mcisVmID_']").each(function(){
    //     if( vmID == $(this).val()){
    //         var mcisVm = $(this).attr("id").split("_")
    //         mcisIndex = mcisVm[1]
    //         vmIndex = mcisVm[2]
    //         return false;
    //     }
    // });

    var vmCreateScript = JSON.stringify(data)
    // $("#m_exportFileName_" + mcisIndex + "_"+ vmIndex).val(vmName);
    // $("#m_vmExportScript_" + mcisIndex + "_"+vmIndex).val(vmCreateScript);
    $("#exportFileName").val(mcisID + "_" + vmID);
    $("#exportScript").val(vmCreateScript);
    // }
    // ).catch(function(error){
    //     var statusCode = error.response.data.status;
    //     var message = error.response.data.message;
    //     commonErrorAlert(statusCode, message)
    // });
    // ).catch((error) => {
    //     console.warn(error);
    //     console.log(error.response)
    //     var errorMessage = error.response.statusText;
    //     var statusCode = error.response.status;
    //     commonErrorAlert(statusCode, errorMessage)
    // });

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
var vmChartArr = new Array();
function showVmMonitoring(mcisID, vmID) {
    $("#mcis_detail_info_check_monitoring").prop("checked", true)
    $("#mcis_detail_info_check_monitoring").attr("disabled", true)
    // $("#Monitoring_tab").show();
    //var duration = "5m"
    var duration = "30m"
    var period_type = "m"
    var metric_arr = ["cpu", "memory", "disk", "network"];
    var statisticsCriteria = "last";
    // TODO : Analytics View 는 안보이게
    for (var i in metric_arr) {
        var vmChart;
        vmChart = getVmMetric(vmChart, "canvas_" + i, metric_arr[i], mcisID, vmID, metric_arr[i], period_type, statisticsCriteria, duration);
        vmChartArr.push(vmChart);
    }
    //$("#Monitoring_tab").hide();

}


// getVMMetric 는 mcis.chart.js로 이동


////////////////

// MCIS script export
function mcisScriptExport() {

    // var mcisID = $("#mcis_id").val();
    var vmID = $("#vm_id").val();
    var vmName = $("#vm_name").val();

    var checkedCount = 0;
    var mcisID = "";
    $("[id^='td_ch_']").each(function () {

        if ($(this).is(":checked")) {
            checkedCount++;
            console.log("checked")
            mcisID = $(this).val();
            // 여러개를 지울 때 호출하는 함수를 만들어 여기에서 호출
        } else {
            console.log("checked nothing")

        }
    })

    if (checkedCount == 0) {
        commonAlert("Please Select MCIS!!")
        return;
    } else if (checkedCount > 1) {
        commonAlert("Please Select One MCIS at a time")
        return;
    }

    // 위 값으로 mcisIndex, vmIndex 를 찾자
    var mcisIndex = 0;
    var vmIndex = 0;
    //mcisID{{$index}}
    console.log("mcisScriptExport start")
    $("[id^='mcisID']").each(function () {
        if (mcisID == $(this).val()) {
            mcisIndex = $(this).attr("id").replace("mcisID", "")
            return false;
        }
    });

    // MCIS정보를 가져온 뒤 save하자.
    getCommonMcisData('mcisexport', mcisID);

    // console.log("index " + mcisIndex)
    // if( $("#m_mcisExportScript_" + mcisIndex).val() == ""){
    //     makeMcisScript(mcisIndex);
    // }
    // console.log("mcisscript")
    // console.log($("#m_mcisExportScript_" + mcisIndex).val())
    // saveToMcisAsJsonFile(mcisIndex, vmIndex);
}

// vm script export
function vmScriptExport() {

    var mcisID = $("#mcis_id").val();
    var vmID = $("#vm_id").val();
    var vmName = $("#vm_name").val();
    console.log("vmScriptExport start")
    // 위 값으로 mcisIndex, vmIndex 를 찾자
    var mcisIndex = 0;
    var vmIndex = 0;
    $("[id^='mcisVmID_']").each(function () {
        if (vmID == $(this).val()) {
            var mcisVm = $(this).attr("id").split("_")
            mcisIndex = mcisVm[1]
            vmIndex = mcisVm[2]
            return false;
        }
    });

    if (!mcisID) {
        commonAlert("Please Select MCIS!!")
        return;
    }
    if (!vmID) {
        commonAlert("Please Select VM!!")
        return;
    }

    // console.log("index " + mcisIndex + " , " + vmIndex)
    // if( $("#m_vmExportScript_" + mcisIndex + "_" + vmIndex).val() == ""){
    //     makeVmScript(mcisIndex, vmIndex);
    // }
    console.log("vmscript " + mcisIndex + " : " + vmIndex)
    console.log($("#m_vmExportScript_" + mcisIndex + "_" + vmIndex).val())
    saveToVmAsJsonFile(mcisIndex, vmIndex);
}

// mcis를 선택하면 해당 mcis를 export할 준비를 함
// lifecycle 의 ExportScriptOfMcis 를 통해 선택한 mcis script를 file로 저장
function makeMcisScript(mcisIndex) {
    var vms = 'mcisVmID_' + mcisIndex + "_";
    var vmIndex = 0;
    // // vmScript 먼저 생성
    // console.log("in makeMcisScript " + mcisIndex + " vms : " + vms);
    // $("[id^='" + vms +"']").each(function(){
    //     makeVmScript(mcisIndex, vmIndex);
    //     vmIndex++;
    // });
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
    for (var i = 0; i < vmIndex; i++) {
        var vmScript = $("#m_vmExportScript_" + mcisIndex + "_" + i).val();// 여기에 담겨있음.(위에서 먼저 호출해서 생성 해 둠)

        console.log(i);
        console.log(vmScript);

        if (vmScript == undefined) continue;// VM이 Terminated 된 경우 등에서는 vmScript가 정상적으로 생성되지 않음.

        if (addedVmIndex > 0) mcisCreateScript += ",";

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
function makeVmScript(mcisIndex, vmIndex) {
    console.log("in makeVmScript" + mcisIndex + " : " + vmIndex)
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

    for (var i = 0; i < sgArr.length; i++) {
        if (i > 0) vmCreateScript += ','
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


    $("#m_exportFileName_" + mcisIndex + "_" + vmIndex).val(nameVal);
    $("#m_vmExportScript_" + mcisIndex + "_" + vmIndex).val(vmCreateScript);
    console.log("vmCreateScript============" + mcisIndex + ":" + vmIndex);
    console.log(vmCreateScript);
}

// json 파일로 저장
function saveToMcisAsJsonFile(mcisIndex) {
    // var fileName = "MCIS_" + $("#m_exportFileName_" + mcisIndex).val();
    // var exportScript = $("#m_mcisExportScript_" + mcisIndex).val();
    var fileName = $("#exportFileName").val();
    var exportScript = $("#exportScript").val();

    saveFileProcess(fileName, exportScript);
}
function saveToVmAsJsonFile(mcisIndex, vmIndex) {
    // var fileName = "VM_" + $("#m_exportFileName_" + mcisIndex + "_" + vmIndex).val();
    // var exportScript = $("#m_vmExportScript_" + mcisIndex + "_"  + vmIndex).val();
    var fileName = $("#exportFileName").val();
    var exportScript = $("#exportScript").val();

    saveFileProcess(fileName, exportScript);
}

// 파일명, script대로 파일 생성
function saveFileProcess(fileName, exportScript) {

    var element = document.createElement('a');
    // element.setAttribute('href', 'data:text/plain;charset=utf-8,' + encodeURIComponent(exportScript));
    element.setAttribute('href', 'data:text/json;charset=utf-8,' + encodeURIComponent(exportScript));
    // element.setAttribute('download', fileName);
    element.setAttribute('download', fileName + ".json");

    element.style.display = 'none';
    document.body.appendChild(element);

    element.click();

    document.body.removeChild(element);

    $("#exportFileName").val('');
    $("#exportScript").val('');

}

function getSecurityGroupCallbackSuccess(caller, data) {
    var html = ""
    var firewallRules = data.firewallRules

    $("#register_box").modal()
    firewallRules.map(item => (html += '<tr>'
        + '<td class="btn_mtd" data-th="fromPort">' + item.fromPort + ' <span class="ov off"></span></td>'
        + '<td class="overlay hidden" data-th="toPort">' + item.toPort + '</td>'
        + '<td class="overlay hidden" data-th="toProtocol">' + item.ipProtocol + '</td>'
        + '<td class="overlay hidden " data-th="direction">' + item.direction + '</td>'
        + '</tr>'
    ))
    $("#manage_mcis_popup_sg").empty()
    $("#manage_mcis_popup_sg").append(html)
}

function getSecurityGroupCallbackFail(error) {

}

function getCommonMcisDataCallbackSuccess(caller, data, mcisID) {
    if (caller == "mcisexport") {

        var mcisIndex = 0;
        // console.log("mcisScriptExport start " + mcisID)
        // console.log(data);
        $("[id^='mcisID']").each(function () {
            if (mcisID == $(this).val()) {
                mcisIndex = $(this).attr("id").replace("mcisID", "")
                return false;
            }
        });
        //
        // var mcisNameVal = $("#m_mcisName_" + mcisIndex).val();
        var mcisCreateScript = JSON.stringify(data)
        console.log(mcisCreateScript)
        // $("#m_exportFileName_" + mcisIndex).val(mcisNameVal);
        // $("#m_mcisExportScript_" + mcisIndex).val(mcisCreateScript);
        $("#exportFileName").val(data.id);
        $("#exportScript").val(mcisCreateScript);

        saveToMcisAsJsonFile(mcisIndex);
    } else if (caller == "refreshmcisdata") {
        console.log(" gogo ")
        // var mcisID = $("#selected_mcis_id").val();
        // var mcisIndex = $("#selected_mcis_index").val();
        // console.log("setMcisData " + mcisIndex)
        // setMcisData(data, mcisIndex);
        // console.log("clock List of Mcis ")
        // clickListOfMcis(mcisID,mcisIndex);

        updateMcisData(data, mcisID);
    }
}
function getMcisDataCallbackFail(error) {

}

// 지도에 marker로 region 표시.
// Map 관련 설정
var locationMap = new Object();
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
    //지도 그리기 관련
    var polyArr = new Array();

    var longitudeValue = locationInfo.longitude;
    var latitudeValue = locationInfo.latitude;
    console.log(longitudeValue + " : " + latitudeValue);
    if (longitudeValue && latitudeValue) {
        $("#map").empty()
        var locationMap = map_init()// mcis.map.js 파일에 정의되어 있으므로 import 필요.  TODO : map click할 때 feature 에 id가 없어 tooltip 에러나고 있음. 해결필요
        drawMap(locationMap, longitudeValue, latitudeValue, locationInfo)
    }
}



function getCommonVmImageInfoCallbackSuccess(caller, imageInfo) {
    // var imageInfo = data;
    var html = ""
    console.log("image info : ", imageInfo)
    html += '<a href="javascript:void(0);" title="' + imageInfo.cspImageName + '">' + imageInfo.id + '</a>'
        + '<div class="bb_info">Image Name : ' + imageInfo.name + ', GuestOS:' + imageInfo.guestOS + '</div>'

    $("#server_detail_view_image_id").empty();
    $("#server_detail_view_image_id").append(html);
    $("#server_info_os").val(imageInfo.guestOS);
    $("#server_detail_view_os").val(imageInfo.guestOS);
    bubble_box();
}

function getCommonVmImageInfoCallbackFail(caller, data) {
    // -- fail 나더라도 그냥 넘어감.
}

function bubble_box() {
    $(".bubble_box .box").each(function () {
        var $list = $(this);
        var bubble = $list.find('.bb_info');
        var menuTime;
        $list.mouseenter(function () {
            bubble.fadeIn(300);
            clearTimeout(menuTime);
        }).mouseleave(function () {
            clearTimeout(menuTime);
            menuTime = setTimeout(mTime, 100);
        });
        function mTime() {
            bubble.stop().fadeOut(100);
        }
    });
}

function remoteCommandMcis(commandWord) {
    // mcis가 선택되어 있어야 하고
    var checked_nothing = 0;
    $("[id^='td_ch_']").each(function () {

        if ($(this).is(":checked")) {
            checked_nothing++;
            console.log("checked")
            var mcisID = $(this).val()
            postRemoteCommandMcis(mcisID, commandWord);
        } else {
            console.log("checked nothing")

        }
    })
    if (checked_nothing == 0) {
        commonAlert("Please Select MCIS!!")
        return;
    }

}


function remoteCommandVmMcis(commandWord) {
    // VM 선택되어 있어야
    var mcisID = $("#mcis_id").val();
    var vmID = $("#vm_id").val();
    var vmName = $("#vm_name").val();
    console.log("remoteCommandVmMcis start")
    // 위 값으로 mcisIndex, vmIndex 를 찾자
    var mcisIndex = 0;
    var vmIndex = 0;
    $("[id^='mcisVmID_']").each(function () {
        if (vmID == $(this).val()) {
            var mcisVm = $(this).attr("id").split("_")
            mcisIndex = mcisVm[1]
            vmIndex = mcisVm[2]
            return false;
        }
    });

    if (!mcisID) {
        commonAlert("Please Select MCIS!!")
        return;
    }
    if (!vmID) {
        commonAlert("Please Select VM!!")
        return;
    }
    console.log(" commandWord = " + commandWord);
    if (!commandWord) {
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

/////////////  set : data설정, caculate : 연산, display : 화면 표시
var totalMcisListObj = new Object();// 모든 MCIS정보를 다 넣고 갱신하여 빼 쓰도록
var totalMcisStatusMap = new Map();
function setToTalMcisStatus() {
    try {
        for (var mcisIndex in totalMcisListObj) {
            var aMcis = totalMcisListObj[mcisIndex]
            var aMcisStatusCountMap = calculateMcisStatusCount(aMcis);
            totalMcisStatusMap.set(aMcis.id, aMcisStatusCountMap)
        }
    } catch (e) {
        console.log("mcis status error")
    }
    displayMcisStatusArea();
}
// // 해당 mcis에서 상태값들을 count : 1개 mcis의 상태는 1개만 있으므로 running, stop, terminate 중 1개만 1, 나머지는 0
// function calculateMcisStatusCount(mcisData){
//     console.log("calculateMcisStatusCount")
//     console.log(mcisData)
//     var mcisStatusCountMap = new Map();
//     mcisStatusCountMap.set("running", 0);
//     mcisStatusCountMap.set("stop", 0);  // partial 도 stop으로 보고있음.
//     mcisStatusCountMap.set("terminate", 0);
//     try{
//         var mcisStatus = mcisData.status
//         var mcisDispStatus = getMcisStatusDisp(mcisStatus);// 화면 표시용 status
//
//         if( mcisStatus != ""){// mcis status 가 없는 경우는 skip
//             if( mcisStatusCountMap.has(mcisDispStatus) ){
//                 mcisStatusCountMap.set(mcisDispStatus, mcisStatusCountMap.get(mcisDispStatus) + 1)
//             }
//         }
//     }catch(e){
//         console.log("mcis status error")
//     }
//     console.log(mcisStatusCountMap);
//     return mcisStatusCountMap;
// }
// 화면 표시
function displayMcisStatusArea() {
    var sumMcisCnt = 0;
    var sumMcisRunningCnt = 0;
    var sumMcisStopCnt = 0;
    var sumMcisTerminateCnt = 0;
    totalMcisStatusMap.forEach((value, key) => {
        var statusRunning = value.get("running");
        var statusStop = value.get("stop");
        var statusTerminate = value.get("terminate")
        sumMcisRunningCnt += statusRunning
        sumMcisStopCnt += statusStop
        sumMcisTerminateCnt += statusTerminate
    });
    sumMcisCnt = sumMcisRunningCnt + sumMcisStopCnt + sumMcisTerminateCnt
    $("#total_mcis").text(sumMcisCnt);
    $("#mcis_status_running").text(sumMcisRunningCnt);
    $("#mcis_status_stopped").text(sumMcisStopCnt);
    $("#mcis_status_terminated").text(sumMcisTerminateCnt);

}
var totalVmStatusMap = new Map();
// Mcis 목록에서 vmStatus만 처리 : 화면표시는 display function에서
function setTotalVmStatus() {
    try {
        for (var mcisIndex in totalMcisListObj) {
            var aMcis = totalMcisListObj[mcisIndex]
            var vmStatusCountMap = calculateVmStatusCount(aMcis);
            totalVmStatusMap.set(aMcis.id, vmStatusCountMap)
        }
    } catch (e) {
        console.log("mcis status error")
    }
    displayVmStatusArea();
}
// mcis의 vm들 상태값 변경
function setVmStatus(mcisIndex, mcisIs) {

}

// 1개 mcis 아래의 vm 들의 status만 계산
// function calculateVmStatusCount(vmList){
//     console.log("calculateVmStatusCount")
//     console.log(vmList)
//     var sumVmCnt = 0;
//     var vmStatusCountMap = new Map();
//     vmStatusCountMap.set("running", 0);
//     vmStatusCountMap.set("stop", 0);  // partial 도 stop으로 보고있음.
//     vmStatusCountMap.set("terminate", 0);
//     try{
//         for(var vmIndex in vmList) {
//             var aVm = vmList[vmIndex];
//             var vmStatus = aVm.status;
//             var vmDispStatus = getVmStatusDisp(vmStatus);
//
//             if (vmStatus != "") {// vm status 가 없는 경우는 skip
//                 if (vmStatusCountMap.has(vmDispStatus)) {
//                     vmStatusCountMap.set(vmDispStatus, vmStatusCountMap.get(vmDispStatus) + 1)
//                 }
//             }
//         }
//     }catch(e){
//         console.log("mcis status error")
//     }
//     return vmStatusCountMap;
// }

// 화면 표시
function displayVmStatusArea() {
    var sumVmCnt = 0;
    var sumVmRunningCnt = 0;
    var sumVmStopCnt = 0;
    var sumVmTerminateCnt = 0;
    totalVmStatusMap.forEach((value, key) => {
        var statusRunning = value.get("running");
        var statusStop = value.get("stop");
        var statusTerminate = value.get("terminate")
        sumVmRunningCnt += statusRunning
        sumVmStopCnt += statusStop
        sumVmTerminateCnt += statusTerminate
    });
    sumVmCnt = sumVmRunningCnt + sumVmStopCnt + sumVmTerminateCnt
    $("#total_vm").text(sumVmCnt);
    $("#vm_status_running").text(sumVmRunningCnt);
    $("#vm_status_stopped").text(sumVmStopCnt);
    $("#vm_status_terminated").text(sumVmTerminateCnt);
}
var totalCloudConnectionMap = new Map();
function setTotalConnection() {
    try {
        console.log("setTotalConnection")
        for (var mcisIndex in totalMcisListObj) {
            var aMcis = totalMcisListObj[mcisIndex]
            // var cloudConnectionCountMap = ;
            totalCloudConnectionMap.set(aMcis.id, calculateConnectionCount(aMcis.vm))
        }
    } catch (e) {
        console.log("mcis status error")
    }
    displayConnectionCountArea();
}

function displayConnectionCountArea() {
    // mcis별 합계이므로 total을 구해서 표시해야 함.
    var sumCloudConnectionMap = new Map();
    totalCloudConnectionMap.forEach((value, key) => {

    });

    // 합계를 화면에 표시
}

// MCIS의 proder들 표시
function getProviderNamesOfMcis(mcisID) {
    var mcisProviderNames = "";
    var vmCloudConnectionMap = totalCloudConnectionMap.get(mcisID);
    if (vmCloudConnectionMap) {
        vmCloudConnectionMap.forEach((value, key) => {
            mcisProviderNames += "<img class='provider_icon' src='/assets/img/contents/img_logo_" + key + ".png' alt='" + key + "'/>";
        });
    }
    return mcisProviderNames;
}

function getMCISInfoProviderNames(mcisID) {
    var mcisProviderNames = "";
    var vmCloudConnectionMap = totalCloudConnectionMap.get(mcisID);
    if (vmCloudConnectionMap) {
        vmCloudConnectionMap.forEach((value, key) => {
            mcisProviderNames += "<img class='provider_icon_info' src='/assets/img/contents/img_logo_" + key + ".png' alt='" + key + "'/>";
        });
    }
    return mcisProviderNames;
}

function displayMcisListTable() {

    // elementNames
    // id="server_info_tr_" + mcisIndex             // tr
    // id="mcisInfo_mcisStatus_icon_" + mcisIndex   // icon
    // id="mcisInfo_mcisstatus_" + mcisIndex
    // id="mcisInfo_mcisName_" + mcisIndex
    // id="mcisInfo_mcisProviderNames_" + mcisIndex
    // id="mcisInfo_totalVmCountOfMcis_" + mcisIndex
    // id="mcisInfo_vmstatus_running_" + mcisIndex
    // id="mcisInfo_vmstatus_stop_" + mcisIndex
    // id="mcisInfo_vmstatus_terminate_" + mcisIndex
    // id="mcisInfo_mcisDescription_" + mcisIndex
    // id="td_ch_" + mcisIndex                      // checkbox


    if (!isEmpty(totalMcisListObj) && totalMcisListObj.length > 0) {
        //totalMcisCnt = mcisList.length;
        var addMcis = "";
        for (var mcisIndex in totalMcisListObj) {
            var aMcis = totalMcisListObj[mcisIndex]
            if (aMcis.id != "") {
                addMcis += setMcisListTableRow(aMcis, mcisIndex);
            }
        }// end of mcis loop
        $("#mcisList").empty();
        $("#mcisList").append(addMcis);

    } else {
        var addMcis = "";
        addMcis += '<tr>'
        addMcis += '<td class="overlay hidden" data-th="" colspan="8">No Data</td>'
        addMcis += '</tr>'
        $("#mcisList").empty();
        $("#mcisList").append(addMcis);
    }
}

function setMcisListTableRow(aMcisData, mcisIndex) {
    var mcisTableRow = "";
    var mcisStatus = aMcisData.status
    var vmCloudConnectionMap = totalCloudConnectionMap.get(aMcisData.id);
    var mcisProviderNamesa = getProviderNamesOfMcis(aMcisData.id);//MCIS에 사용 된 provider
    var mcisDispStatus = getMcisStatusDisp(mcisStatus);// 화면 표시용 status

    var vmStatusCountMap = totalVmStatusMap.get(aMcisData.id);
    var totalVmCountOfMcis = vmStatusCountMap.get('running') + vmStatusCountMap.get('stop') + vmStatusCountMap.get('terminate');
    console.log(vmStatusCountMap);

    // elementNames
    // id="server_info_tr_" + mcisIndex             // tr
    // id="mcisInfo_mcisStatus_icon_" + mcisIndex   // icon
    // id="mcisInfo_mcisstatus_" + mcisIndex
    // id="mcisInfo_mcisName_" + mcisIndex
    // id="mcisInfo_mcisProviderNames_" + mcisIndex
    // id="mcisInfo_totalVmCountOfMcis_" + mcisIndex
    // id="mcisInfo_vmstatus_running_" + mcisIndex
    // id="mcisInfo_vmstatus_stop_" + mcisIndex
    // id="mcisInfo_vmstatus_terminate_" + mcisIndex
    // id="mcisInfo_mcisDescription_" + mcisIndex
    // id="td_ch_" + mcisIndex                      // checkbox
    console.log("label = " + aMcisData.systemLabel);
    // List of Mcis table
    try {
        mcisTableRow += '<tr onclick="clickListOfMcis(\'' + aMcisData.id + '\');" id="server_info_tr_' + mcisIndex + '" item="' + aMcisData.id + '|' + mcisIndex + '">'

        mcisTableRow += '<td class="overlay hidden td_left column-20percent" data-th="Status">'
        mcisTableRow += '   <div style="display: flex; align-items: center;"><img id="mcisInfo_mcisStatus_icon_' + mcisIndex + '" src="/assets/img/contents/icon_' + mcisDispStatus + '.png" class="icon" alt=""/><span id="mcisInfo_mcisstatus_' + mcisIndex + '">' + mcisStatus + '</span><span class="ov off"></span></div>'
        mcisTableRow += '</td>'
        mcisTableRow += '<td class="btn_mtd ovm column-20percent" data-th="Name">'
        mcisTableRow += '   <div style="display: inline-flex; align-items: center;">'

        mcisTableRow += '       <div id="mcisInfo_mcisName_' + mcisIndex + '">' + aMcisData.name + '</div>'
        if (aMcisData.systemLabel === "Managed by MCKS") {
            mcisTableRow += '   <img class="mcks-icon" src="/assets/img/contents/icon_kubernetes_gray.png"/>'
        }
        mcisTableRow += '       <span class="ov"></span>'
        mcisTableRow += '   </div>'
        mcisTableRow += '</td>'
        mcisTableRow += '<td class="overlay hidden column-14percent" data-th="Cloud Connection">'
        mcisTableRow += '    <div id="mcisInfo_mcisProviderNames_' + mcisIndex + '">'

        if (vmCloudConnectionMap) {
            vmCloudConnectionMap.forEach((value, key) => {
                // + mcisProviderNames +
                mcisTableRow += '    <img class="provider_icon" src="/assets/img/contents/img_logo_' + key + '.png" alt="' + key + '"/>';

            })
        }
        mcisTableRow += '   </div>'
        mcisTableRow += '</td>'

        mcisTableRow += '<td class="overlay hidden column-18percent" data-th="Total Infras"><div id="mcisInfo_totalVmCountOfMcis_' + mcisIndex + '">' + totalVmCountOfMcis + '</div></td>'

        mcisTableRow += '<td class="overlay hidden column-18percent" data-th="# of Servers">'
        mcisTableRow += '<span class="bar" ></span> <span title="running" id="mcisInfo_vmstatus_running_' + mcisIndex + '">' + vmStatusCountMap.get('running') + '</span>'
        mcisTableRow += '<span class="bar" >/</span> <span title="stop" id="mcisInfo_vmstatus_stop_' + mcisIndex + '">' + vmStatusCountMap.get('stop') + '</span>'
        mcisTableRow += '<span class="bar" >/</span> <span title="terminate" id="mcisInfo_vmstatus_terminate_' + mcisIndex + '">' + vmStatusCountMap.get('terminate') + '</span>'
        mcisTableRow += '</td>'

        // mcisTableRow += '<td class="overlay hidden" data-th="Description"><div id="mcisInfo_mcisDescription_' + mcisIndex + '">' + aMcisData.description + '</div></td>'

        mcisTableRow += '<td class="overlay hidden column-60px"  data-th="">'
        mcisTableRow += '<input type="checkbox" name="chk" value="' + aMcisData.id + '" id="td_ch_' + mcisIndex + '" title="" />'
        mcisTableRow += '<label for="td_ch_' + mcisIndex + '"></label>'
        mcisTableRow += '</td>'

        mcisTableRow += '</tr>'

    } catch (e) {
        console.log("list of mcis error")
        console.log(e)

        mcisTableRow = '<tr>'
        mcisTableRow += '<td class="overlay hidden" data-th="" colspan="8">No Data</td>'
        mcisTableRow += '</tr>'
    }
    return mcisTableRow;
}

// MCIS List table의 1개 Row Update
function updateMcisListTableRow(aMcisData, mcisIndex) {
    console.log("updateMcisListTableRow " + mcisIndex);
    console.log(aMcisData);
    var mcisStatus = aMcisData.status
    var mcisProviderNames = getProviderNamesOfMcis(aMcisData.id);//MCIS에 사용 된 provider
    var mcisDispStatus = getMcisStatusDisp(mcisStatus);// 화면 표시용 status

    var vmStatusCountMap = totalVmStatusMap.get(aMcisData.id);
    var mcisStatusImg = "/assets/img/contents/icon_" + mcisDispStatus + ".png"

    if (vmStatusCountMap) {

        var sumVmCountRunning = vmStatusCountMap.get("running")
        var sumVmCountStop = vmStatusCountMap.get("stop")
        var sumVmCountTerminate = vmStatusCountMap.get("terminate")
        var sumVmCount = sumVmCountRunning + sumVmCountStop + sumVmCountTerminate

        // id="server_info_tr_" + mcisIndex             // tr   -> 변경없음
        // id="mcisInfo_mcisStatus_icon_" + mcisIndex   // icon
        $("#mcisInfo_mcisStatus_icon_" + mcisIndex).attr("src", mcisStatusImg);

        // id="mcisInfo_mcisstatus_" + mcisIndex
        $("#mcisInfo_mcisstatus_" + mcisIndex).text(mcisStatus)
        // id="mcisInfo_mcisName_" + mcisIndex
        $("#mcisInfo_mcisName_" + mcisIndex).text(aMcisData.name)
        // id="mcisInfo_mcisProviderNames_" + mcisIndex

        $("#mcisInfo_mcisProviderNames_" + mcisIndex).empty()
        $("#mcisInfo_mcisProviderNames_" + mcisIndex).append(mcisProviderNames)
        // id="mcisInfo_totalVmCountOfMcis_" + mcisIndex
        $("#mcisInfo_totalVmCountOfMcis_" + mcisIndex).text(sumVmCount)
        // id="mcisInfo_vmstatus_running_" + mcisIndex
        $("#mcisInfo_vmstatus_running_" + mcisIndex).text(sumVmCountRunning)
        // id="mcisInfo_vmstatus_stop_" + mcisIndex
        $("#mcisInfo_vmstatus_stop_" + mcisIndex).text(sumVmCountStop)
        // id="mcisInfo_vmstatus_terminate_" + mcisIndex
        $("#mcisInfo_vmstatus_terminate_" + mcisIndex).text(sumVmCountTerminate)
        // id="mcisInfo_mcisDescription_" + mcisIndex
        $("#mcisInfo_mcisDescription_" + mcisIndex).text(sumVmCount)
        // id="td_ch_" + mcisIndex                      // checkbox -> 변경없음
    }
}

function displayMcisInfoArea(mcisData) {

    var mcisID = mcisData.id;
    var mcisName = mcisData.name;
    var mcisDescription = mcisData.description;
    var mcisStatus = mcisData.status;
    var mcisDispStatus = getMcisStatusDisp(mcisStatus);
    var mcisStatusIcon = getMcisStatusIcon(mcisDispStatus);
    var mcisTargetStatus = mcisData.targetStatus;
    var mcisTargetAction = mcisData.targetAction;
    var mcisProviderNames = getMCISInfoProviderNames(mcisData.id);//MCIS에 사용 된 provider

    var vmStatusCountMap = totalVmStatusMap.get(mcisData.id);
    if (vmStatusCountMap) {
        var sumVmCountOfMcis = vmStatusCountMap.get('running') + vmStatusCountMap.get('stop') + vmStatusCountMap.get('terminate');


        $("#mcis_info_txt").text("[ " + mcisName + " ]");

        $("#mcis_server_info_status").empty();
        $("#mcis_server_info_status").append('<strong>Server List / Status</strong>  <span class="stxt">[ ' + mcisName + ' ]</span>  Server(' + sumVmCountOfMcis + ')')

        // mcis info area
        clearMcisInfoAreaData();// 이전에 set된 data가 있으면 clear시킨다.
        // 이미 보여주고 있으면, 상세정보는 숨긴다.
        console.log("displayMcisInfoArea show " + $(".server_status").hasClass("view"))
        if ($(".server_status").hasClass("view")) {
            $(".server_info").removeClass("active");
        } else {
            $(".server_status").addClass("view");
            var offset = $("#mcis_info_box").position();
            console.log("position", offset.top);

            $("#TopWrap").animate({ scrollTop: offset.top * 1.3 }, 300);
        }

        // 초기화 후 set해야하나?

        $("#server_info_status_icon_img").attr("src", mcisStatusIcon);
        $("#mcis_info_name").val(mcisName + " / " + mcisID);
        $("#mcis_info_id").val(mcisID);
        $("#mcis_info_description").val(mcisDescription);
        $("#mcis_info_targetStatus").val(mcisTargetStatus);
        $("#mcis_info_targetAction").val(mcisTargetAction);
        $("#mcis_info_cloud_connection").empty();
        $("#mcis_info_cloud_connection").append(mcisProviderNames);    //

        $("#mcis_name").val(mcisName)


        // id="mcis_info_txt"   // <span class="stxt" id="mcis_info_txt">[ ]</span>
        // id="service_status_icon"
        // id="service_status_icon_img" // <img src="/assets/img/contents/icon_running_db.png" alt="">
        // id="mcis_info_name"  // input type=text"
        // id="mcis_info_description" // input type="text
        // id="mcis_info_targetStatus" // input type=text"
        // id="mcis_info_targetAction" // input type=text"
        // id="mcis_info_cloud_connection" // input type=text"

        // deply Algorithm

        // id="mcis_server_info_status" // div

        // <ul id="mcis_server_info_box">
        displayServerStatusList(mcisID, mcisData.vm)
    }
}

// Mcis Info Area를 data를 초기화 한다.
function clearMcisInfoAreaData() {
    $("#server_info_status_icon_img").attr("src", '');
    $("#mcis_info_name").val('');
    $("#mcis_info_id").val('');
    $("#mcis_info_description").val('');
    $("#mcis_info_targetStatus").val('');
    $("#mcis_info_targetAction").val('');
    $("#mcis_info_cloud_connection").empty();

    $("#mcis_name").val('')
}

// vm 상태별 icon으로 표시
function displayServerStatusList(mcisID, vmList) {
    console.log("displayServerStatusList")
    var mcisName = mcisID;// TODO : vmDetailInfo() 에서 mcisName을 가져올 수 있도록 보완 필요
    var vmLi = "";
    vmList.sort();
    for (var vmIndex in vmList) {
        var aVm = vmList[vmIndex]

        var vmID = aVm.id;
        var vmName = aVm.name;
        var vmStatus = aVm.status;
        var vmDispStatus = getVmStatusDisp(vmStatus);
        var vmStatusClass = getVmStatusClass(vmDispStatus)


        vmLi += '<li id="server_status_icon_' + vmID + '" class="sel_cr ' + vmStatusClass + '"><a href="javascript:void(0);" onclick="vmDetailInfo(\'' + mcisID + '\',\'' + mcisName + '\',\'' + vmID + '\')"><span class="txt">' + vmName + '</span></a></li>';
    }// end of mcis loop
    // console.log(vmLi)
    $("#mcis_server_info_box").empty();
    $("#mcis_server_info_box").append(vmLi);
    //Manage MCIS Server List on/off
    $(".dashboard .ds_cont .area_cont .listbox li.sel_cr").each(function () {
        console.log("sel_cr");
        var $sel_list = $(this),
            $detail = $(".server_info");
        $sel_list.off("click").click(function () {
            $sel_list.addClass("active");
            $sel_list.siblings().removeClass("active");
            $detail.addClass("active");
            $detail.siblings().removeClass("active");
            $sel_list.off("click").click(function () {
                if ($(this).hasClass("active")) {
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

// update McisData
function updateMcisData(aMcisData, mcisID) {
    var mcisIndex = -1;
    for (var mIndex in totalMcisListObj) {
        var aMcis = totalMcisListObj[mIndex]
        if (aMcisData.id == aMcis.id) {
            totalMcisListObj[mIndex] = aMcisData;
            mcisIndex = mIndex;
            break;
        }
    }
    if (mcisIndex > -1) {
        totalMcisListObj.push(aMcisData);
    }

    setToTalMcisStatus();// mcis상태 표시 를 위해 필요
    setTotalVmStatus();// mcis 의 vm들 상태표시 를 위해 필요
    setTotalConnection();// Mcis의 provider별 connection 표시를 위해 필요

    updateMcisListTableRow(aMcisData, mcisIndex);

    $("#mcis_id").val(mcisID);
    $("#selected_mcis_id").val(mcisID);
    $("#selected_mcis_index").val(mcisIndex);

    // MCIS Info area set
    //showServerListAndStatusArea(mcisID,mcisIndex);
    displayMcisInfoArea(aMcisData);
}

// update VM Data
function updateVmData(vmData) {
    // server detailInfo area

    // detail tab

    // connection tab

    // monitoring tab
}

// data 조회완료 시 해당 vmData set.
function setServerDetailInfoArea(mcisID, vmID, vmData) {
    // mcisList -> mcis -> vmList -> vm
    var vmList = new Object()
    var mcisName = mcisID;
    for (var mcisIndex in totalMcisListObj) {
        var aMcis = totalMcisListObj[mcisIndex]
        if (mcisID != "" && mcisID == aMcis.id) {
            vmList = aMcis.vm;
            mcisName = aMcis.name;
            break;
        }
    }

    for (var vmIndex in vmList) {
        var aVm = vmList[vmIndex]
        if (vmData.id == aVm.id) {
            //aVm = vmData;
            vmList[vmIndex] = vmData;
            break;
        }
    }

    displayServerDetailInfoArea(mcisID, vmData);
}

// 해당 area가 나타나면서 set된 data표시
function displayServerDetailInfoArea(mcisID, mcisName, vmData) {
    var vmID = vmData.id;
    var vmName = vmData.name;
    var vmStatus = vmData.status;
    var vmDispStatus = getVmStatusDisp(vmStatus);
    var vmStatusIcon = getVmStatusIcon(vmDispStatus);

    // id="server_detail_info_text" span
    // id="server_detail_info_public_ip_text" div
    // id="mcis_detail_info_check_monitoring" checkbox

    var installMonAgent = vmData.monAgentStatus;
    console.log("install mon agent : ", installMonAgent)
    // if(installMonAgent == "installed"){
    //     // var isWorking = checkDragonFlyMonitoringAgent(mcisID, vmID);
    //     var isWorking = false;
    //     if( isWorking){
    //         $("#mcis_detail_info_check_monitoring").prop("checked",true)
    //         $("#mcis_detail_info_check_monitoring").attr("disabled",true)
    //     }else{
    //         $("#mcis_detail_info_check_monitoring").prop("checked",false)
    //         $("#mcis_detail_info_check_monitoring").attr("disabled",false)
    //     }
    // }else{
    //     $("#mcis_detail_info_check_monitoring").prop("checked",false)
    //     $("#mcis_detail_info_check_monitoring").attr("disabled",false)
    // }

    $("#server_info_text").text('[' + vmName + '/' + mcisName + ']')
    $("#server_info_status_icon_img").attr("src", vmStatusIcon);


    $("#server_detail_info_text").text('[' + vmName + '/' + mcisName + ']')
    $("#server_detail_info_public_ip_text").text("Public IP : " + vmPublicIp)


    $("#server_detail_view_server_status").val(vmStatus);// detail tab


}
