$(document).ready(function () {
    order_type = "name"
    //checkbox all
    $("#th_chall").click(function () {
        if ($("#th_chall").prop("checked")) {
            $("input[name=chk]").prop("checked", true);
        } else {
            $("input[name=chk]").prop("checked", false);
        }
    })

    //table 스크롤바 제한
    $(window).on("load resize", function () {
        var vpwidth = $(window).width();
        if (vpwidth > 768 && vpwidth < 1800) {
            $(".dashboard_cont .dataTable").addClass("scrollbar-inner");
            $(".dataTable.scrollbar-inner").scrollbar();
        } else {
            $(".dashboard_cont .dataTable").removeClass("scrollbar-inner");
        }

        setTableHeightForScroll('vpcListTable', 300);
    });
});

// TODO : filter 기능, sort 기능

// $(document).ready(function () {

//     // var defaultNameSpace = "{{ .DefaultNameSpaceID }}"
//     // alert(defaultNameSpace)
//     // var nameSpaceList = "{{ .NameSpaceList }}"
//     // alert(nameSpaceList);
//     // page load시 이미 가져왔음
//     // getVpcList(order_type);
//     // getCloudOS(apiInfo,'provider');
// })                      


// function goFocus(target) {
//     console.log(event)
//     event.preventDefault();

//     $("#" + target).focus();
//     fnMove(target)
// }

// function fnMove(target) {
//     var offset = $("#" + target).offset();
//     console.log("fn move offset : ", offset);
//     $('html, body').animate({
//         scrollTop: offset.top
//     }, 400);
// }

// function goDelete() {
function deleteVPC() {
    var vNetId = "";
    var count = 0;

    $("input[name='chk']:checked").each(function () {
        count++;
        vNetId = vNetId + $(this).val() + ",";
    });
    vNetId = vNetId.substring(0, vNetId.lastIndexOf(","));

    console.log("vNetId : ", vNetId);
    console.log("count : ", count);

    if (vNetId == '') {
        commonAlert("삭제할 대상을 선택하세요.");
        return false;
    }

    if (count != 1) {
        commonAlert("삭제할 대상을 하나만 선택하세요.");
        return false;
    }

    // var url = CommonURL + "/ns/" + NAMESPACE + "/resources/vNet/" + vNetId;
    var url = "/setting/resources" + "/network/del/" + vNetId
    console.log("del vnet url : ", url);

    axios.delete(url, {
        headers: {
            // 'Authorization': "{{ .apiInfo}}",
            'Content-Type': "application/json"
        }
    }).then(result => {
        var data = result.data;
        console.log(result);
        console.log(data);
        if (result.status == 200 || result.status == 201) {
            //commonAlert("Success Delete Network")
            commonAlert(data.message)
            // location.reload(true);
            //vNetInfoBox 안보이게
            displayVNetInfo("DEL_SUCCESS")
            // getVpcList("name");
        } else {
            commonAlert(result.data.error)
        }
        // }).catch(function(error){
        //     commonAlert(error)
        //     console.log("Network delete error : ",error);        
        // });
    }).catch((error) => {
        console.warn(error);
        console.log(error.response)
        var errorMessage = error.response.data.error;
        var statusCode = error.response.status;
        commonErrorAlert(statusCode, errorMessage);
    });
}

function getVpcList(sort_type) {
    console.log(sort_type);
    // var url = CommonURL + "/ns/" + NAMESPACE + "/resources/vNet";
    //var currentNameSpace = $('$topboxDefaultNameSpaceID').val()
    // defaultNameSpace 기준으로 가져온다. (server단 session에서 가져오므로 변경하려면 현재 namesapce를 바꾸고 호출하면 됨)
    var url = "/setting/resources/network/list";
    axios.get(url, {
        headers: {
            'Content-Type': "application/json"
        }
    }).then(result => {
        console.log("get VPC List : ", result.data);
        // var data = result.data.vNet;
        var data = result.data.VNetList;

        var html = ""
        var cnt = 0;

        if (data == null) {
            html += '<tr><td class="overlay hidden" data-th="" colspan="5">No Data</td></tr>'

            $("#vpcList").empty()
            $("#vpcList").append(html)

            ModalDetail()
        } else {
            if (data.length) {
                if (sort_type) {
                    cnt++;
                    console.log("check : ", sort_type);
                    data.filter(list => list.Name !== "").sort((a, b) => (a[sort_type] < b[sort_type] ? - 1 : a[sort_type] > b[sort_type] ? 1 : 0)).map((item, index) => (
                        html += addVNetRow(item, index)
                        // html += '<tr onclick="showVNetInfo(\'' + item.Name + '\');">' 
                        //     + '<td class="overlay hidden" data-th="">' 
                        //     + '<input type="hidden" id="sg_info_' + index + '" value="' + item.Name + '|' + item.CidrBlock + '"/>' 
                        //     + '<input type="checkbox" name="chk" value="' + item.Name + '" id="raw_'  + index + '" title="" /><label for="td_ch1"></label> <span class="ov off"></span></td>' 
                        //     + '<td class="btn_mtd ovm" data-th="name">' + item.Name + '</td>'
                        //     + '<td class="overlay hidden" data-th="cidrBlock">' + item.CidrBlock + '</td>' 
                        //     + '<td class="overlay hidden" data-th="description">' + item.Description + '</td>'  
                        //     + '<td class="overlay hidden" data-th=""><a href="javascript:void(0);"><img src="/assets/img/contents/icon_link.png" class="icon" alt=""/></a></td>' 
                        //     + '</tr>'
                    ))
                } else {
                    data.filter((list) => list.Name !== "").map((item, index) => (
                        html += addVNetRow(item, index)
                        // html += '<tr onclick="showVNetInfo(\'' + item.Name + '\');">' 
                        //     + '<td class="overlay hidden" data-th="">' 
                        //     + '<input type="hidden" id="sg_info_' + index + '" value="' + item.Name  + '"/>'
                        //     + '<input type="checkbox" name="chk" value="' + item.Name + '" id="raw_' + index + '" title="" /><label for="td_ch1"></label> <span class="ov off"></span></td>' 
                        //     + '<td class="btn_mtd ovm" data-th="name">' + item.Name + '<span class="ov"></span></td>' 
                        //     + '<td class="overlay hidden" data-th="cidrBlock">' + item.CidrBlock + '</td>' 
                        //     + '<td class="overlay hidden" data-th="description">' + item.Description + '</td>' 
                        //     + '<td class="overlay hidden" data-th=""><a href="javascript:void(0);"><img src="/assets/img/contents/icon_link.png" class="icon" alt=""/></a></td>' 
                        //     + '</tr>'
                    ))
                }

                $("#vpcList").empty()
                $("#vpcList").append(html)

                ModalDetail()
            }
        }


        // }).catch(function(error){
        //     console.log("Network list error : ",error);        
        // });
    }).catch((error) => {
        console.warn(error);
        console.log(error.response)
        var errorMessage = error.response.data.error;
        var statusCode = error.response.status;
        commonErrorAlert(statusCode, errorMessage);
    });
}

// VNet목록에 Item 추가
function addVNetRow(item, index) {
    console.log("addVnetRow " + index);
    console.log(item)
    var html = ""
    html += '<tr onclick="showVNetInfo(\'' + item.name + '\');">'
        + '<td class="overlay hidden column-50px" data-th="">'
        + '<input type="hidden" id="sg_info_' + index + '" value="' + item.name + '"/>'
        + '<input type="checkbox" name="chk" value="' + item.name + '" id="raw_' + index + '" title="" /><label for="td_ch1"></label> <span class="ov off"></span></td>'
        + '<td class="btn_mtd ovm" data-th="name">' + item.name + '<span class="ov"></span></td>'
        + '<td class="overlay hidden" data-th="cidrBlock">' + item.cidrBlock + '</td>'
        + '<td class="overlay hidden" data-th="description">' + item.description + '</td>'
        // + '<td class="overlay hidden column-60px" data-th=""><a href="javascript:void(0);"><img src="/assets/img/contents/icon_link.png" class="icon" alt=""/></a></td>' 
        + '</tr>'
    return html;
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

function displayVNetInfo(targetAction) {
    if (targetAction == "REG") {
        $('#vnetCreateBox').toggleClass("active");
        $('#vNetInfoBox').removeClass("view");
        $('#vNetListTable').removeClass("on");
        var offset = $("#vnetCreateBox").offset();
        // var offset = $("#" + target+"").offset();
        $("#TopWrap").animate({ scrollTop: offset.top }, 300);

        // form 초기화
        $("#regVpcName").val('')
        $("#regDescription").val('')
        $("#regCidrBlock").val('')
        $("#regSubnet").val('')
        goFocus('vnetCreateBox');
    } else if (targetAction == "REG_SUCCESS") {
        $('#vnetCreateBox').removeClass("active");
        $('#vNetInfoBox').removeClass("view");
        $('#vNetListTable').addClass("on");

        var offset = $("#vnetCreateBox").offset();
        $("#TopWrap").animate({ scrollTop: offset.top }, 0);

        // form 초기화
        $("#regVpcName").val('')
        $("#regDescription").val('')
        $("#regCidrBlock").val('')
        $("#regSubnet").val('')
        getVpcList("name");
    } else if (targetAction == "DEL") {
        $('#vnetCreateBox').removeClass("active");
        $('#vNetInfoBox').addClass("view");
        $('#vNetListTable').removeClass("on");

        var offset = $("#vNetInfoBox").offset();
        $("#TopWrap").animate({ scrollTop: offset.top }, 300);

    } else if (targetAction == "DEL_SUCCESS") {
        $('#vnetCreateBox').removeClass("active");
        $('#vNetInfoBox').removeClass("view");
        $('#vNetListTable').addClass("on");

        var offset = $("#vNetInfoBox").offset();
        $("#TopWrap").animate({ scrollTop: offset.top }, 0);

        getVpcList("name");
    } else if (targetAction == "CLOSE") {
        $('#vnetCreateBox').removeClass("active");
        $('#vNetInfoBox').removeClass("view");
        $('#vNetListTable').addClass("on");

        var offset = $("#vNetInfoBox").offset();
        $("#TopWrap").animate({ scrollTop: offset.top }, 0);
    }


    //CreateBox
    // $('#RegistBox .btn_ok.register').click(function(){
    // 		$(".dashboard.register_cont").toggleClass("active");
    // 		$(".dashboard.server_status").removeClass("view");
    // 		$(".dashboard .status_list tbody tr").removeClass("on");
    // 		//ok 위치이동
    // 		$('#RegistBox').on('hidden.bs.modal', function () {
    // 			var offset = $("#CreateBox").offset();
    // 			$("#wrap").animate({scrollTop : offset.top}, 300);
    // 		})		
}

// provider에 등록 된 connection 목록 표시
function getConnectionInfo(provider) {
    // var url = SpiderURL+"/connectionconfig";
    var url = "/setting/connections/cloudconnectionconfig/" + "list"
    // console.log("provider : ",provider)
    // var provider = $("#provider option:selected").val();
    var html = "";
    // var apiInfo = ApiInfo
    axios.get(url, {
        headers: {
            // 'Authorization': apiInfo
        }
    }).then(result => {
        console.log('getConnectionConfig result: ', result)
        // var data = result.data.connectionconfig
        var data = result.data.ConnectionConfig
        console.log("connection data : ", data);
        var count = 0;
        var configName = "";
        var confArr = new Array();
        for (var i in data) {
            if (provider == data[i].ProviderName) {
                count++;
                html += '<option value="' + data[i].ConfigName + '" item="' + data[i].ProviderName + '">' + data[i].ConfigName + '</option>';
                configName = data[i].ConfigName
                confArr.push(data[i].ConfigName)
            }
        }
        if (count == 0) {
            commonAlert("해당 Provider에 등록된 Connection 정보가 없습니다.")
            html += '<option selected>Select Configname</option>';
        }
        if (confArr.length > 1) {
            configName = confArr[0];
        }
        $("#regConnectionName").empty();
        $("#regConnectionName").append(html);

        // }).catch(function(error){
        //     console.log("Network data error : ",error);        
        // });
    }).catch((error) => {
        console.warn(error);
        console.log(error.response)
    });
}

// 팝업의 subnet을 set
function applySubnet() {
    var subnetNameValue = $("input[name='reg_subnetName']").length;
    var subnetCIDRBlockValue = $("input[name='reg_subnetCidrBlock']").length;

    var subnetNameData = new Array(subnetNameValue);
    var subnetCIDRBlockData = new Array(subnetCIDRBlockValue);

    for (var i = 0; i < subnetNameValue; i++) {
        subnetNameData[i] = $("input[name='reg_subnetName']")[i].value;
        console.log("subnetNameData" + [i] + " : ", subnetNameData[i]);
    }
    for (var i = 0; i < subnetCIDRBlockValue; i++) {
        subnetCIDRBlockData[i] = $("input[name='reg_subnetCidrBlock']")[i].value;
        console.log("subnetCIDRBlockData" + [i] + " : ", subnetCIDRBlockData[i]);
    }

    subnetJsonList = new Array();//subnet 저장할 array. 전역으로 선언

    for (var i = 0; i < subnetNameValue; i++) {
        var SNData = "SNData" + i;
        var SNData = new Object();
        SNData.name = subnetNameData[i];
        SNData.ipv4_CIDR = subnetCIDRBlockData[i];
        subnetJsonList.push(SNData);
    }

    var infoshow = "";
    for (var i in subnetJsonList) {
        infoshow += subnetJsonList[i].name + " (" + subnetJsonList[i].ipv4_CIDR + ") ";
    }
    $("#regSubnet").empty();
    $("#regSubnet").val(infoshow);
    $("#subnetRegisterBox").modal("hide");
}

function createVNet() {
    var vpcName = $("#regVpcName").val();
    var description = $("#regDescription").val();
    var connectionName = $("#regConnectionName").val();
    var cidrBlock = $("#regCidrBlock").val();
    if (!vpcName) {
        commonAlert("Input New VPC Name")
        $("#regVpcName").focus()
        return;
    }
    console.log(subnetJsonList);

    // var apiInfo = "{{ .apiInfo}}";
    // var url = CommonURL+"/ns/"+NAMESPACE+"/resources/vNet"
    var url = "/setting/resources" + "/network/reg"
    console.log("vNet Reg URL : ", url)
    var obj = {
        CidrBlock: cidrBlock,
        ConnectionName: connectionName,
        Description: description,
        Name: vpcName,
        SubnetInfoList: subnetJsonList
    }
    console.log("info vNet obj Data : ", obj);

    if (vpcName) {
        axios.post(url, obj, {
            headers: {
                'Content-type': 'application/json',
                // 'Authorization': apiInfo,
            }
        }).then(result => {
            console.log("result vNet : ", result);
            var data = result.data;
            console.log(data);
            if (data.status == 200 || data.status == 201) {
                commonAlert("Success Create Network(VPC)!!")
                //등록하고 나서 화면을 그냥 고칠 것인가?
                // 등록 성공시 등록한 객체가 들어 옴. 일단 기존 List에 추가하는 것으로?
                // var data = result.data;
                // console.log(data);
                // var html = addVNetRow(data)
                // $("#vpcList").append(html)

                displayVNetInfo("REG_SUCCESS")


                //아니면 화면을 리로딩 시킬것인가?
                // location.reload();
                // $("#btn_add2").click()
                // $("#namespace").val('')
                // $("#nsDesc").val('')
            } else {
                commonAlert("Fail Create Network(VPC) " + data.message)
            }
            // }).catch(function(error){
            //     var data = error.data;
            //         console.log(data);
            //     console.log(error);        
            //     commonAlert("Network create error : ",error)            
            // });
        }).catch((error) => {
            // console.warn(error);
            console.log(error.response)
            var errorMessage = error.response.data.error;
            var statusCode = error.response.status;
            commonErrorAlert(statusCode, errorMessage)
        });
    } else {
        commonAlert("Input VPC Name")
        $("#regVpcName").focus()
        return;
    }
}

// 선택한 vNet의 상세정보 : 이미 가져왔는데 다시 가져올 필요있나?? vNetID
function showVNetInfo(vpcName) {
    console.log("showVNetInfo : ", vpcName);
    // var apiInfo = "{{ .apiInfo}}";
    // var vNetId = encodeURIComponent(vNetName);
    // $('.stxt').html(vpcName);
    $('#networkVpcName').text(vpcName)

    // var url = CommonURL+"/ns/"+NAMESPACE+"/resources/vNet/"+ vNetId;
    var url = "/setting/resources" + "/network/" + encodeURIComponent(vpcName);
    console.log("vnet detail URL : ", url)

    return axios.get(url, {
        // headers:{
        //     'Authorization': apiInfo
        // }
    }).then(result => {
        console.log(result);
        console.log(result.data);
        var data = result.data.VNetInfo
        console.log("Show Data : ", data);

        var dtlVpcName = data.name;
        var dtlDescription = data.description;
        var dtlConnectionName = data.connectionName;
        var dtlCidrBlock = data.cidrBlock;
        var dtlSubnet = "";

        var subList = data.subnetInfoList;
        for (var i in subList) {
            // dtlSubnet += subList[i].iid.nameId + " (" + subList[i].ipv4_CIDR + ")";
            dtlSubnet += subList[i].id + " (" + subList[i].ipv4_CIDR + ")";
        }
        console.log("dtlSubnet : ", dtlSubnet);

        $("#dtlVpcName").empty();
        $("#dtlDescription").empty();
        $("#dtlProvider").empty();
        $("#dtlConnectionName").empty();
        $("#dtlCidrBlock").empty();
        $("#dtlSubnet").empty();

        $("#dtlVpcName").val(dtlVpcName);
        $("#dtlDescription").val(dtlDescription);
        $("#dtlConnectionName").val(dtlConnectionName);
        $("#dtlCidrBlock").val(dtlCidrBlock);
        $("#dtlSubnet").val(dtlSubnet);

        if (dtlConnectionName == '' || dtlConnectionName == undefined) {
            commonAlert("ConnectionName is empty")
        } else {
            // getProvider(dtlConnectionName);
            // var providerValue = getProviderNameByConnection(dtlConnectionName);
            getProviderNameByConnection(dtlConnectionName, 'dtlProvider')// provider는 connection 정보에서 가져옴
            //$("#dtlProvider").val(providerValue);
        }

    }).catch(function (error) {
        console.log("Network detail error : ", error);
    });
}

// // 특정 connection 정보에서 Privider set
// function getProvider(connectionName) {
//     console.log("getProvider  : ",connectionName);
//     // var url = SpiderURL+"/connectionconfig/" + target;
//     var url = "/setting/connections"+"/cloudconnectionconfig/" + connectionName;
//     return axios.get(url,{
//         // headers:{
//         //     'Authorization': apiInfo
//         // }    
//     }).then(result=>{
//         var data = result.data;
//         console.log(data)
//         console.log(data.ConnectionConfig)
//         var provider = data.ConnectionConfig.ProviderName;
//         //var Provider = data.ConnectionConfig.providerName;
//         console.log(provider)
//         $("#dtlProvider").val(provider);
//     }).catch(function(error){
//         console.log("Network getProvider error : ",error);        
//     });
// }

function displaySubnetRegModal(isShow) {
    if (isShow) {
        $("#subnetRegisterBox").modal();
        $('.dtbox.scrollbar-inner').scrollbar();
    } else {
        $("#vnetCreateBox").toggleClass("active");
    }
}
// $(document).ready(function() {
//     var subnetJsonList = "";
//     //Subnet pop table scrollbar
//       $('.btn_register').on('click', function() {
//         $("#register_box").modal();
//         $('.dtbox.scrollbar-inner').scrollbar();
//     });	

//     /*
//     $('.register_cont .btn_cancel').click(function(){
//         $(".dashboard.register_cont").toggleClass("active");
//     });
//     */
// });


// subnet popup

// $(document).on("click","button[name=btn_add]",function(){
//     var addStaffText = '<tr class="ip" name="tr_Input">'
//         + '<td class="btn_mtd" data-th="subnet Name"><input type="text" name="" value="" placeholder="" class="pline" title="" /> <span class="ov up" name="td_ov"]></span></td>'
//         + '<td class="overlay" data-th="cidrBlock"><input type="text" name="" value="" placeholder="" class="pline" title="" /></td>'
//         + '<td class="overlay">'
//         + '<button class="btn btn_add" name="btn_add" value="">add</button>'
//         + '<button class="btn btn_del" name="delSubnet" value="">del</button>'
//         + '</td>'
//         + '</tr>';
//     var trHtml = $( "tr[name=tr_Input]:last" );
//     trHtml.after(addStaffText);
// });

// $('.dataTable .btn.btn_add').on("click", function() {
//         trHtml.after(addStaffText);
// });
var subnetJsonList = "";//저장시 subnet목록을 담을 array 
var addStaffText = '<tr name="tr_Input">'
    + '<td class="btn_mtd column-40percent" data-th="subnet Name"><input type="text" id="regSubnetName" name="reg_subnetName" value="" placeholder="" class="pline" title="" /> <span class="ov up" name="td_ov"]></span></td>'
    + '<td class="overlay" data-th="cidrBlock"><input type="text" id="regSubnetCidrBlock" name="reg_subnetCidrBlock" value="" placeholder="" class="pline" title="" /></td>'
    + '<td class="overlay column-100px">'
    + '<button class="btn btn_add" name="addSubnet" value="">add</button>'
    + '<button class="btn btn_del" name="delSubnet" value="">del</button>'
    + '</td>'
    + '</tr>';

$(document).on("click", "button[name=addSubnet]", function () {
    console.log("add subnet clicked")
    var subnetNameValue = $("input[name='reg_subnetName']").length;
    var trHtml = $("tr[name=tr_Input]:last");
    trHtml.after(addStaffText);
});
$(document).on("click", "button[name=delSubnet]", function () {
    console.log("del subnet clicked")
    var trHtml = $(this).parent().parent();
    trHtml.remove();
});

// $(document).on("click","span[name=td_ov]",function(){
//     var trHtml = $(this).parent().parent();
//     trHtml.find(".btn_mtd").toggleClass("over");
//     trHtml.find(".overlay").toggleClass("hidden");
// });