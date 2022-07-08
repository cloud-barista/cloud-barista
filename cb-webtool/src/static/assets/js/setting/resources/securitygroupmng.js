$(document).ready(function () {
    //action register open / table view close
    // $('#RegistBox .btn_ok.register').click(function(){
    //     $(".dashboard.register_cont").toggleClass("active");
    //     $(".dashboard.server_status").removeClass("view");
    //     $(".dashboard .status_list tbody tr").removeClass("on");
    //     //ok 위치이동
    //     $('#RegistBox').on('hidden.bs.modal', function () {
    //         var offset = $("#CreateBox").offset();
    //         $("#wrap").animate({scrollTop : offset.top}, 300);
    //     })		
    // });
});
// scroll
$(document).ready(function () {
    //checkbox all
    // $("#th_chall").click(function() {
    // if ($("#th_chall").prop("checked")) {
    //     $("input[name=chk]").prop("checked", true);
    // } else {
    //     $("input[name=chk]").prop("checked", false);
    // }
    // })

    //table 스크롤바 제한
    $(window).on("load resize", function () {
        //     var vpwidth = $(window).width();
        //   if (vpwidth > 768 && vpwidth < 1800) {
        //     $(".dashboard_cont .dataTable").addClass("scrollbar-inner");
        //          $(".dataTable.scrollbar-inner").scrollbar();
        //   } else {
        //     $(".dashboard_cont .dataTable").removeClass("scrollbar-inner");
        //   }

        setTableHeightForScroll('securityGroupList', 300)

    });
});

$(document).ready(function () {
    // order_type = "name"
    // getSGList(order_type);
    // var apiInfo = "{{ .apiInfo}}";
    // getCloudOS(apiInfo,'provider');

    //firewallRegisterBox
})


// add/delete 시 area 표시
function displaySecurityGroupInfo(targetAction) {
    if (targetAction == "REG") {
        $('#securityGroupCreateBox').toggleClass("active");
        $('#securityGroupInfoBox').removeClass("view");
        $('#securityGroupListTable').removeClass("on");
        var offset = $("#securityGroupCreateBox").offset();
        // var offset = $("#" + target+"").offset();
        $("#TopWrap").animate({ scrollTop: offset.top }, 300);

        // form 초기화
        $("#regVpcName").val('')
        $("#regDescription").val('')
        $("#regCidrBlock").val('')
        $("#regSubnet").val('')
        goFocus('securityGroupCreateBox');
    } else if (targetAction == "REG_SUCCESS") {
        $('#securityGroupCreateBox').removeClass("active");
        $('#securityGroupInfoBox').removeClass("view");
        $('#securityGroupListTable').addClass("on");

        var offset = $("#securityGroupCreateBox").offset();
        $("#TopWrap").animate({ scrollTop: offset.top }, 0);

        // form 초기화
        $("#regCspSecurityGroupName").val('')
        $("#regDescription").val('')
        $("#regProvider").val('')
        $("#regConnectionName").val('')

        $("#regVNetId").val('')
        $("#regInbound").val('')
        $("#regOutbound").val('')

        getSecurityGroupList("name");
    } else if (targetAction == "DEL") {
        $('#securityGroupCreateBox').removeClass("active");
        $('#securityGroupInfoBox').addClass("view");
        $('#securityGroupListTable').removeClass("on");

        var offset = $("#securityGroupInfoBox").offset();
        $("#TopWrap").animate({ scrollTop: offset.top }, 300);

    } else if (targetAction == "DEL_SUCCESS") {
        $('#securityGroupCreateBox').removeClass("active");
        $('#securityGroupInfoBox').removeClass("view");
        $('#securityGroupListTable').addClass("on");

        var offset = $("#securityGroupInfoBox").offset();
        $("#TopWrap").animate({ scrollTop: offset.top }, 0);

        getSecurityGroupList("name");
    } else if (targetAction == "CLOSE") {
        $('#securityGroupCreateBox').removeClass("active");
        $('#securityGroupInfoBox').removeClass("view");
        $('#securityGroupListTable').addClass("on");

        var offset = $("#securityGroupInfoBox").offset();
        $("#TopWrap").animate({ scrollTop: offset.top }, 0);
    }
}

function deleteSecurityGroup() {
    var sgId = "";
    var count = 0;

    $("input[name='chk']:checked").each(function () {
        count++;
        sgId = sgId + $(this).val() + ",";
    });
    sgId = sgId.substring(0, sgId.lastIndexOf(","));

    console.log("sgId : ", sgId);
    console.log("count : ", count);

    if (sgId == '') {
        commonAlert("삭제할 대상을 선택하세요.");
        return false;
    }

    if (count != 1) {
        commonAlert("삭제할 대상을 하나만 선택하세요.");
        return false;
    }

    //var url = CommonURL + "/ns/" + NAMESPACE + "/resources/securityGroup/" + sgId;
    var url = "/setting/resources" + "/securitygroup/del/" + sgId

    axios.delete(url, {
        headers: {
            // 'Authorization': "{{ .apiInfo}}",
            'Content-Type': "application/json"
        }
    }).then(result => {
        var data = result.data
        console.log(result);
        console.log(data);
        if (result.status == 200 || result.status == 201) {
            commonAlert(data.message);
            // location.reload(true);
            displaySecurityGroupInfo("DEL_SUCCESS")
        } else {
            commonAlert(data.error);
        }
        // }).catch(function(error){
        //     console.log("sg del error : ",error);        
        // });
    }).catch((error) => {
        console.warn(error);
        console.log(error.response)
        var errorMessage = error.response.data.error;
        var statusCode = error.response.status;
        commonErrorAlert(statusCode, errorMessage);
    });
}

// Security Group 목록 조회
function getSecurityGroupList(sortType) {
    getCommonSecurityGroupList("securitygroupmng", sortType)
}

function getSecurityGroupListCallbackFail(error) {
    var errorMessage = error.response.data.error;
    var statusCode = error.response.status;
    commonErrorAlert(statusCode, errorMessage);
}

function setSecurityGroupListAtServerImage(data, sortType) {
    var html = ""
    console.log("Data : ", data);

    if (data == null) {
        html += '<tr><td class="overlay hidden" data-th="" colspan="5">No Data</td></tr>'

        $("#sgList").empty();
        $("#sgList").append(html);

        ModalDetail()
    } else {
        if (data.length) { // null exception if not exist
            if (sortType) {
                console.log("check : ", sortType);
                data.filter(list => list.name !== "").sort((a, b) => (a[sortType] < b[sortType] ? - 1 : a[sortType] > b[sortType] ? 1 : 0)).map((item, index) => (
                    html += '<tr onclick="showSecurityGroupInfo(\'' + item.cspSecurityGroupName + '\');">'
                    + '<td class="overlay hidden column-50px" data-th="">'
                    + '<input type="hidden" id="sg_info_' + index + '" value="' + item.cspSecurityGroupName + '|' + item.connectionName + '"/>'
                    + '<input type="checkbox" name="chk" value="' + item.cspSecurityGroupName + '" id="raw_' + index + '" title="" /><label for="td_ch1"></label> <span class="ov off"></span>'
                    + '<input type="hidden" value="' + item.systemLabel + '"/>'
                    + '</td>'
                    + '<td class="btn_mtd ovm" data-th="cspSecurityGroupName">' + item.cspSecurityGroupName
                    // + '<a href="javascript:void(0);"><img src="/assets/img/contents/icon_copy.png" class="td_icon" alt=""/></a> <span class="ov"></span></td>'
                    + '</td>'
                    + '<td class="overlay hidden" data-th="connectionName">' + item.connectionName + '</td>'
                    + '<td class="overlay hidden" data-th="description">' + item.description + '</td>'
                    // + '<td class="overlay hidden column-60px" data-th=""><a href="javascript:void(0);"><img src="/assets/img/contents/icon_link.png" class="icon" alt=""/></a></td>' 
                    + '</tr>'
                ))
            } else {
                data.filter((list) => list.name !== "").map((item, index) => (
                    html += '<tr onclick="showSecurityGroupInfo(\'' + item.cspSecurityGroupName + '\');">'
                    + '<td class="overlay hidden column-50px" data-th="">'
                    + '<input type="hidden" id="sg_info_' + index + '" value="' + item.cspSecurityGroupName + '"/>'
                    + '<input type="checkbox" name="chk" value="' + item.cspSecurityGroupName + '" id="raw_' + index + '" title="" /><label for="td_ch1"></label> <span class="ov off"></span>'
                    + '<input type="hidden" value="' + item.systemLabel + '"/>'
                    + '</td>'
                    + '<td class="btn_mtd ovm" data-th="cspSecurityGroupName">' + item.cspSecurityGroupName + '<span class="ov"></span></td>'
                    + '<td class="overlay hidden" data-th="connectionName">' + item.connectionName + '</td>'
                    + '<td class="overlay hidden" data-th="description">' + item.description + '</td>'
                    // + '<td class="overlay hidden column-60px" data-th=""><a href="javascript:void(0);"><img src="/assets/img/contents/icon_link.png" class="icon" alt=""/></a></td>' 
                    + '</tr>'
                ))

            }

            $("#sgList").empty();
            $("#sgList").append(html);

            ModalDetail()
        }
    }
}

function ModalDetail() {
    $(".dashboard .status_list tbody tr").each(function () {
        var $td_list = $(this),
            $status = $(".server_status"),
            $detail = $(".server_info");
        $td_list.off("click").click(function () {
            $td_list.addClass("on");
            $td_list.siblings().removeClass("on");
            $status.addClass("view");
            $status.siblings().removeClass("on");
            $(".dashboard.register_cont").removeClass("active");
            $td_list.off("click").click(function () {
                if ($(this).hasClass("on")) {
                    console.log("reg ok button click")
                    $td_list.removeClass("on");
                    $status.removeClass("view");
                    $detail.removeClass("active");
                } else {
                    $td_list.addClass("on");
                    $td_list.siblings().removeClass("on");
                    $status.addClass("view");

                    $status.siblings().removeClass("view");
                    $(".dashboard.register_cont").removeClass("active");
                }
            });
        });
    });
}

function showSecurityGroupInfo(sgName) {
    console.log("sgName showSecurityGroupInfo : ", sgName);
    //var sgName = target;

    $(".stxt").html(sgName);

    // var apiInfo = "{{ .apiInfo}}";

    // var url = CommonURL+"/ns/"+NAMESPACE+"/resources/securityGroup/"+ sgName;
    var url = "/setting/resources" + "/securitygroup/" + sgName
    console.log("security group URL : ", url)

    return axios.get(url, {
        headers: {
            // 'Authorization': apiInfo
        }

    }).then(result => {
        //var data = result.data
        console.log(result.data);
        var data = result.data.SecurityGroupInfo;
        console.log("Show Data : ", data);

        var dtlCspSecurityGroupName = data.cspSecurityGroupName;
        var dtlDescription = data.description;
        var dtlConnectionName = data.connectionName;
        var dtlvNetId = data.vNetId;

        var dtlFirewall = data.firewallRules;
        console.log("firefire : ", dtlFirewall);
        var inbounds = "";
        var outbounds = "";
        var cidrs = "";
        for (var i in dtlFirewall) {
            console.log("direc : ", dtlFirewall[i].direction);
            var firewallDirection = (dtlFirewall[i].direction).toLowerCase();
            if (firewallDirection == "inbound" || firewallDirection == "ingress") {
                inbounds += dtlFirewall[i].ipProtocol
                    + ' ' + dtlFirewall[i].fromPort + '~' + dtlFirewall[i].toPort + ' '
            } else if (firewallDirection == "outbound" || firewallDirection == "egress") {
                outbounds += dtlFirewall[i].ipProtocol
                    + ' ' + dtlFirewall[i].fromPort + '~' + dtlFirewall[i].toPort + ' '
            } else {// 정의되지 않은 항목은 inbound쪽에 몰아주기
                inbound += dtlFirewall[i].ipProtocol
                    + ' ' + dtlFirewall[i].fromPort + '~' + dtlFirewall[i].toPort + ' '
            }
            cidrs += dtlFirewall[i].cidr + ' ';
        }
        console.log(cidrs);
        // console.log(dtlvNetId);
        $('#dtlCspSecurityGroupName').empty();
        $('#dtlDescription').empty();
        $('#dtlProvider').empty();
        $('#dtlConnectionName').empty();
        $('#dtlvNetId').empty();
        $('#dtlInbound').empty();
        $('#dtlOutbound').empty();
        $('#dtlCidr').empty();

        $('#dtlCspSecurityGroupName').val(dtlCspSecurityGroupName);
        $('#dtlDescription').val(dtlDescription);
        //$('#dtlProvider').val(dtlProvider);
        $('#dtlConnectionName').val(dtlConnectionName);
        $('#dtlvNetId').val(dtlvNetId);
        $('#dtlInbound').append(inbounds);
        $('#dtlOutbound').append(outbounds);
        $('#dtlCidr').append(cidrs);

        getProviderNameByConnection(dtlConnectionName, 'dtlProvider')// provider는 connection 정보에서 가져옴
        // var providerValue = getProviderNameByConnection(dtlConnectionName, dtlProvider)
        // $('#dtlProvider').val(providerValue);
        // console.log("providerValue = " + providerValue)

        displaySecurityGroupInfo("DEL")// 상세 area 보여주고 focus이동
        // }).catch(function(error){
        //     console.log("show sg info error : ",error);        
        // });
    }).catch((error) => {
        console.warn(error);
        console.log(error.response)
        var errorMessage = error.response.data.error;
        var statusCode = error.response.status;
        commonErrorAlert(statusCode, errorMessage);
    });
}

// Inbound / Outbound Modal 표시
function displayInOutBoundRegModal(isShow) {
    if (isShow) {
        var provider = $("#regProvider option:selected").val();
        console.log("provider " + provider)
        if (provider == "") {
            commonAlert("Please select a provider first");
            return;
        }
        setTableHeightForScroll('firewallRSlistTbl', 300);

        $("#firewallRegisterBox").modal();
        $('.dtbox.scrollbar-inner').scrollbar();
    } else {
        $("#securityGroupCreateBox").toggleClass("active");
    }
}

// function getProvider(target) {
//     console.log("getProvidergetProvider : ",target);
//     var url2 = SpiderURL+"/connectionconfig/" + target;

//     return axios.get(url2,{
//         headers:{
//             'Authorization': apiInfo
//         }

//     }).then(result=>{
//         var data2 = result.data;
//         console.log("adddd : ", data2);

//         var Provider = data2.ProviderName;

//         $('#dtlProvider').val(Provider);
//     })        
// }

// function getConnectionInfo(provider){
//     // var url = SpiderURL+"/connectionconfig";
//     var url = ""
//     console.log("provider : ",provider)
//     //var provider = $("#provider option:selected").val();
//     var html = "";
//     // var apiInfo = ApiInfo
//     axios.get(url,{
//         headers:{
//             // 'Authorization': apiInfo
//         }
//     }).then(result=>{
//         console.log('getConnectionConfig result: ',result)
//         var data = result.data.connectionconfig
//         console.log("connection data : ",data);
//         var count = 0; 
//         var configName = "";
//         var confArr = new Array();
//         for(var i in data){
//             if(provider == data[i].ProviderName){ 
//                 count++;
//                 html += '<option value="'+data[i].ConfigName+'" item="'+data[i].ProviderName+'">'+data[i].ConfigName+'</option>';
//                 configName = data[i].ConfigName
//                 confArr.push(data[i].ConfigName)

//             }
//         }
//         if(count == 0){
//             alert("해당 Provider에 등록된 Connection 정보가 없습니다.")
//                 html +='<option selected>Select Configname</option>';
//         }
//         if(confArr.length > 1){
//             configName = confArr[0];
//         }
//         $("#regConnectionName").empty();
//         $("#regConnectionName").append(html); 

//         getVnetInfo(configName);
//     })
// }

// function getVnetInfo(configName){
//     console.log("vnet : ", configName);
//     var configName = configName;
//     if(!configName){
//         configName = $("#connectionName option:selected").val();
//     }
//     var url = CommonURL+"/ns/"+NAMESPACE+"/resources/vNet";
//     var html = "";
//     var apiInfo = ApiInfo;
//     axios.get(url,{
//         headers:{
//             'Authorization': apiInfo
//         }
//     }).then(result=>{
//         data = result.data.vNet;
//         console.log("vNetwork Info : ",result);
//         console.log("vNetwork data : ",data);
//         for(var i in data){
//             if(data[i].connectionName == configName){
//                 html += '<option value="'+data[i].id+'" selected>'+data[i].cspVNetName+'('+data[i].id+')</option>'; 
//             }
//         }

//         $("#regVNetId").empty();
//         $("#regVNetId").append(html);  
//     })
// }

$(document).ready(function () {
    //Firewall RuleSet pop table scrollbar
    // var fwrsJsonList = "";

    // $('.btn_register').on('click', function() {
    //     $("#firewallRegisterBox").modal();
    //         $('.dtbox.scrollbar-inner').scrollbar();
    //     });	

});

function applyFirewallRuleSet() {
    var provider = $("#regProvider option:selected").val();
    var fromPortValue = $("input[name='fromport']").length;
    var toPortValue = $("input[name='toport']").length;
    var ipprotocolValue = $("select[name='ipprotocol']").length;
    var directionValue = $("select[name='direction']").length;
    var cidrValue = $("input[name='cidr']").length;
    var isSameDirection = true;// direction이 하나만 선택해야하는 경우 사용
    var isSameCidr = true;// cidr을 하나만 사용할 수 있는 경우 사용

    var fromPortData = new Array(fromPortValue);
    var toPortData = new Array(toPortValue);
    var ipprotocolData = new Array(ipprotocolValue);
    var directionData = new Array(directionValue);
    var cidrData = new Array(cidrValue);

    for (var i = 0; i < fromPortValue; i++) {
        fromPortData[i] = $("input[name='fromport']")[i].value;
        console.log("fromPortData" + [i] + " : ", fromPortData[i]);
    }
    for (var i = 0; i < toPortValue; i++) {
        toPortData[i] = $("input[name='toport']")[i].value;
        console.log("toPortData" + [i] + " : ", toPortData[i]);
    }
    for (var i = 0; i < ipprotocolValue; i++) {
        ipprotocolData[i] = $("select[name='ipprotocol']")[i].value;
        console.log("ipprotocolData" + [i] + " : ", ipprotocolData[i]);
    }
    var firstDirectionValue = $("select[name='direction']")[0].value;
    for (var i = 0; i < directionValue; i++) {
        directionData[i] = $("select[name='direction']")[i].value;
        console.log("directionData" + [i] + " : ", directionData[i]);
        if (firstDirectionValue != directionData[i]) {
            isSameDirection = false;
        }
    }

    var firstCidrValue = $("input[name='cidr']")[0].value;
    for (var i = 0; i < cidrValue; i++) {
        cidrData[i] = $("input[name='cidr']")[i].value;
        console.log("firstCidrValue : " + firstCidrValue + " : ", cidrData[i]);
        if (firstCidrValue != cidrData[i]) {
            isSameCidr = false;
        }
    }

    if (provider == "GCP") {
        if (isSameDirection == false) {
            commonAlert("GCP allows only one direction");
            return;
        }
        if (isSameCidr == false) {
            commonAlert("GCP CIDR Blocks must all be same");
            return;
        }

        var cidrSize = $("input[name='cidr']").length;
        if (cidrSize == 1) {
            var cidrValue = $("input[name='cidr']")[0].value;
            if (!cidrValue) {
                commonAlert("GCP requires CIDR Block");
                return;
            }
        } else {
            for (var i = 0; i < cidrValue; i++) {
                cidrData[i] = $("input[name='cidr']")[i].value;
                console.log("cidrData" + [i] + " : ", cidrData[i]);
                if (!cidrData[i]) {
                    commonAlert("GCP requires CIDR Block");
                    return;
                }
            }
        }

    }

    fwrsJsonList = new Array();

    for (var i = 0; i < fromPortValue; i++) {
        var RSData = "RSData" + i;
        var RSData = new Object();
        RSData.direction = directionData[i];
        RSData.fromPort = fromPortData[i];
        RSData.ipProtocol = ipprotocolData[i];
        RSData.toPort = toPortData[i];
        RSData.cidr = cidrData[i];

        fwrsJsonList.push(RSData);
    }

    var inbound = "";
    var outbound = "";
    for (var i in fwrsJsonList) {
        console.log(fwrsJsonList[i]);
        if (fwrsJsonList[i].direction == "inbound") {
            inbound += fwrsJsonList[i].ipProtocol
                + ' ' + fwrsJsonList[i].fromPort + '~' + fwrsJsonList[i].toPort + ' '
        } else if (fwrsJsonList[i].direction == "outbound") {
            outbound += fwrsJsonList[i].ipProtocol
                + ' ' + fwrsJsonList[i].fromPort + '~' + fwrsJsonList[i].toPort + ' '
        }
    }
    console.log("ininin : ", inbound);
    console.log("outoutout : ", outbound);

    $("#regInbound").empty();
    $("#regOutbound").empty();
    $("#regInbound").append(inbound);
    $("#regOutbound").append(outbound);

    $("#firewallRegisterBox").modal("hide");
}

function createSecurityGroup() {
    var cspSecurityGroupName = $("#regCspSecurityGroupName").val();
    var description = $("#regDescription").val();
    var provider = $("#regProvider option:selected").val();
    var connectionName = $("#regConnectionName").val();
    var vNetId = $("#regVNetId").val();

    if (!cspSecurityGroupName) {
        commonAlert("Input New Security Group Name")
        $("#regCspSshKeyName").focus()
        return;
    }

    // connection 이 GCP인 경우 체크 : inbound/outbund는 1종류만 가능, cidrBlock 필수
    // 여기에서도 check가 필요할 까? fwrsJsonList -> 돌아야 함.


    console.log("cspSecurityGroupName : ", cspSecurityGroupName);
    console.log("description : ", description);
    console.log("provider : ", provider);
    console.log("connectionName : ", connectionName);
    console.log("vNetId : ", vNetId);
    console.log("fwrsJsonList : ", fwrsJsonList);// TODO : 비어있으면 에러나므로 valid check 필요

    var apiInfo = "{{ .apiInfo}}";
    // var url = CommonURL+"/ns/"+NAMESPACE+"/resources/securityGroup"
    var url = "/setting/resources" + "/securitygroup/reg"
    console.log("Security Group Reg URL : ", url)
    var obj = {
        connectionName: connectionName,
        description: description,
        firewallRules: fwrsJsonList,
        name: cspSecurityGroupName,
        vNetId: vNetId
    }
    console.log("info connectionconfig obj Data : ", obj);
    if (cspSecurityGroupName) {
        axios.post(url, obj, {
            headers: {
                'Content-type': 'application/json',
                'Authorization': apiInfo,
            }
        }).then(result => {
            console.log("result sg : ", result);
            if (result.status == 200 || result.status == 201) {
                commonAlert("Success Create Security Group!!");
                //displaySecurityGroupInfo("REG_SUCCESS");
                //등록하고 나서 화면을 그냥 고칠 것인가?
                displaySecurityGroupInfo("REG_SUCCESS")
                //아니면 화면을 리로딩 시킬것인가?
                // location.reload();
                // $("#btn_add2").click()
                // $("#namespace").val('')
                // $("#nsDesc").val('')
            } else {
                commonAlert("Fail Create Security Group")
            }
            // }).catch(function(error){
        }).catch((error) => {
            console.warn(error);
            console.log(error.response)

            console.log("sg create error : ");

            //return c.JSON(http.StatusBadRequest, map[string]interface{}{
            // 	"message": respStatus.Message,
            // 	"status":  respStatus.StatusCode,
            // })
            // console.log(error.response.data);
            // console.log("Status", error.response.status);// map interface의 status
            // console.log("Error", error.response.data.error);// map interface의 message
            // console.log('Error', error.message);// [0]번째 param : http.StatusBadRequest

            // message는 , 로 나눈 뒤 : 를 기준으로 key,value로 파싱하면 될 것 같음.
            // console.log(error.config);
        });
    } else {
        alert("Input Security Group Name")
        $("#regCspSecurityGroupName").focus()
        return;
    }
}

var fwrsJsonList = "";// firewallRuleSet 담을 array
function getStaffText() {
    var addStaffText =
        '<tr class="ip" name="tr_Input">' +
        '<td class="btn_mtd column-16percent" data-th="fromPort"><input type="text" name="fromport" value="" placeholder="" class="pline" title="" /> <span class="ov up" name="td_ov"]></span></td>' +
        '<td class="overlay column-16percent" data-th="toPort"><input type="text" name="toport" value="" placeholder="" class="pline" title="" /></td>' +
        '<td class="overlay column-16percent" data-th="ipProtocol">' +
        '<select class="selectbox white pline" name="ipprotocol">' +
        '<option value="tcp">TCP</option>' +
        '<option value="udp">UDP</option>' +
        '<option value="*">ALL</option>' +
        '</select>' +
        '</td>' +
        '<td class="overlay" data-th="direction">' +
        '<select class="selectbox white pline" name="direction">' +
        '<option value="inbound">Inbound</option>' +
        '<option value="outbound">Outbound</option>' +
        '</select>' +
        '</td>' +
        '<td class="btn_mtd column-16percent" data-th="cidr">' +
        '<input type="text" value="" name="cidr" placeholder="" class="pline" title="" /> ' +
        '<span class="ov off"></span>' +
        '</td>' +
        '<td class="overlay column-80px">' +
        '<button class="btn btn_add" name="btn_add" value="">add</button>' +
        '<button class="btn btn_del" name="delInput" value="">del</button>' +
        '</td>' +
        '</tr>';
    return addStaffText;
}

//table row add
$(document).on("click", "button[name=btn_add]", function () {
    // var addStaffText = 
    // '<tr class="ip" name="tr_Input">'+
    //     '<td class="btn_mtd" data-th="fromPort"><input type="text" name="fromport" value="" placeholder="" class="pline" title="" /> <span class="ov up" name="td_ov"]></span></td>'+
    //     '<td class="overlay" data-th="toPort"><input type="text" name="toport" value="" placeholder="" class="pline" title="" /></td>'+
    //     '<td class="overlay" data-th="ipProtocol">'+
    //             '<select class="selectbox white pline" name="ipprotocol">'+
    //                 '<option value="tcp">TCP</option>'+
    //                 '<option value="udp">UDP</option>'+
    //             '</select>'+
    //     '</td>'+
    //     '<td class="overlay" data-th="direction">'+
    //             '<select class="selectbox white pline" name="direction">'+
    //                 '<option value="inbound">Inbound</option>'+
    //                 '<option value="outbound">Outbound</option>'+
    //             '</select>'+
    //     '</td>'+
    //     '<td class="overlay">'+
    //         '<button class="btn btn_add" name="btn_add" value="">add</button>'+
    //         '<button class="btn btn_del" name="delInput" value="">del</button>'+
    //     '</td>'+
    // '</tr>';
    var trHtml = $("tr[name=tr_Input]:last");
    // trHtml.after(addStaffText);
    trHtml.after(getStaffText());
});
// $('.dataTable .btn.btn_add').on("click", function() {
//     // trHtml.after(addStaffText);
//     var trHtml = $( "tr[name=tr_Input]:last" );
//     trHtml.after(getStaffText());
// });
$(document).on("click", "button[name=delInput]", function () {
    var trHtml = $(this).parent().parent();
    trHtml.remove();
});

$(document).on("click", "span[name=td_ov]", function () {
    var trHtml = $(this).parent().parent();
    trHtml.find(".btn_mtd").toggleClass("over");
    trHtml.find(".overlay").toggleClass("hidden");
});