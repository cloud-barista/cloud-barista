$(document).ready(function(){
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

    //checkbox all
    // $("#th_chall").click(function() {
    //     if ($("#th_chall").prop("checked")) {
    //         $("input[name=chk]").prop("checked", true);
    //     } else {
    //         $("input[name=chk]").prop("checked", false);
    //     }
    // })
        
    // //table 스크롤바 제한
    // $(window).on("load resize",function(){
    //         var vpwidth = $(window).width();
    //     if (vpwidth > 768 && vpwidth < 1800) {
    //         $(".dashboard_cont .dataTable").addClass("scrollbar-inner");
    //             $(".dataTable.scrollbar-inner").scrollbar();
    //     } else {
    //         $(".dashboard_cont .dataTable").removeClass("scrollbar-inner");
    //     }
    // });

    setTableHeightForScroll('serverSpecList', 300)

    $('.btn_assist').on('click', function() {
        lookupSpecList()
    });
});

$(document).ready(function () {
    // order_type = "name"
    // getVMSpecList(order_type);

    // var apiInfo = "{{ .apiInfo}}";
    // getCloudOS(apiInfo,'provider');
});

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

// 등록/상세 area 보이기 숨기기
function displayVmSpecInfo(targetAction){
    if( targetAction == "REG"){
        $('#vmSpecCreateBox').toggleClass("active");
        $('#vmSpecInfoBox').removeClass("view");
        $('#vmSpecListTable').removeClass("on");
        var offset = $("#vmSpecCreateBox").offset();
        // var offset = $("#" + target+"").offset();
    	$("#TopWrap").animate({scrollTop : offset.top}, 300);

        // form 초기화
        $("#regSpecName").val('')
        $("#regCspSpecName").val('')

    }else if ( targetAction == "REG_SUCCESS"){
        $('#vmSpecCreateBox').removeClass("active");
        $('#vmSpecInfoBox').removeClass("view");
        $('#vmSpecListTable').addClass("on");
        var html = '<option selected>Select Configname</option>';
        // form 초기화
        $("#regSpecName").val('');
        $("#regProvider").val('');	
        // $("#regConnectionName").empty();
        //$("#regConnectionName").append(html);
        $("#regConnectionName").val('');
        $("#regCspSpecName").val('');
        
        var offset = $("#vmSpecCreateBox").offset();
        $("#TopWrap").animate({scrollTop : offset.top}, 0);
        
        getVmSpecList("name");
    }else if ( targetAction == "DEL"){
        $('#vmSpecCreateBox').removeClass("active");
        $('#vmSpecInfoBox').addClass("view");
        $('#vmSpecListTable').removeClass("on");

        var offset = $("#vmSpecInfoBox").offset();
    	$("#TopWrap").animate({scrollTop : offset.top}, 300);

    }else if ( targetAction == "DEL_SUCCESS"){
        console.log("$$$$$$$$$DelSuccess$$$$$$$$$$$");
        $('#vmSpecCreateBox').removeClass("active");
        $('#vmSpecInfoBox').removeClass("view");
        $('#vmSpecListTable').addClass("on");

        var offset = $("#vmSpecInfoBox").offset();
        $("#TopWrap").animate({scrollTop : offset.top}, 0);

        getVmSpecList("name");
    }else if ( targetAction == "CLOSE"){
        $('#vmSpecCreateBox').removeClass("active");
        $('#vmSpecInfoBox').removeClass("view");
        $('#vmSpecListTable').addClass("on");

        var offset = $("#vmSpecInfoBox").offset();
        $("#TopWrap").animate({scrollTop : offset.top}, 0);
    }
}

function virtualMachineSpecListCallbackSuccess(caller, data, sortType) {
// function setVirtualMachineSpecListAtServerSpec(data, sort_type) {
    var html = ""
    console.log("Caller : ", caller);
    console.log("Data : ", data);
    console.log("SortType : ", sortType);

    if (data == null) {
        html += '<tr><td class="overlay hidden" data-th="" colspan="5">No Data</td></tr>'

        $("#specList").empty()
        $("#specList").append(html)
    } else {
        if (data.length) {
            if (sortType) {
                console.log("check : ", sortType);
                data.filter(list => list.name !== "").sort((a, b) => (a[sortType] < b[sortType] ? - 1 : a[sortType] > b[sortType] ? 1 : 0)).map((item, index) => (
                    html += '<tr onclick="showVmSpecInfo(\'' + item.name + '\');">' 
                        + '<td class="overlay hidden column-50px" data-th="">' 
                        + '<input type="hidden" id="spec_info_' + index + '" value="' + item.name + '|' + item.connectionName + '|' + item.cspSpecName + '"/>' 
                        + '<input type="checkbox" name="chk" value="' + item.name + '" id="raw_'  + index + '" title="" /><label for="td_ch1"></label> <span class="ov off"></span></td>' 
                        + '<td class="btn_mtd ovm" data-th="name ">' + item.name  + '<span class="ov"></span></td>'
                        + '<td class="overlay hidden" data-th="connectionName">' + item.connectionName + '</td>' 
                        + '<td class="overlay hidden" data-th="cspSpecName">' + item.cspSpecName + '</td>'  
                        // + '<td class="overlay hidden column-60px" data-th=""><a href="javascript:void(0);"><img src="/assets/img/contents/icon_link.png" class="icon" alt=""/></a></td>' 
                        + '</tr>'
                ))
            } else {
                data.filter((list) => list.name !== "").map((item, index) => (
                    html += '<tr onclick="showVmSpecInfo(\'' + item.name + '\');">' 
                        + '<td class="overlay hidden column-50px" data-th="">' 
                        + '<input type="hidden" id="spec_info_' + index + '" value="' + item.name + '"/>' 
                        + '<input type="checkbox" name="chk" value="' + item.name + '" id="raw_'  + index + '" title="" /><label for="td_ch1"></label> <span class="ov off"></span></td>' 
                        + '<td class="btn_mtd ovm" data-th="name ">' + item.name  + '<span class="ov"></span></td>'
                        + '<td class="overlay hidden" data-th="connectionName">' + item.connectionName + '</td>' 
                        + '<td class="overlay hidden" data-th="cspSpecName">' + item.cspSpecName + '</td>'  
                        // + '<td class="overlay hidden column-60px" data-th=""><a href="javascript:void(0);"><img src="/assets/img/contents/icon_link.png" class="icon" alt=""/></a></td>' 
                        + '</tr>'
                ))
            }

            $("#specList").empty()
            $("#specList").append(html)
            console.log("setVirtualMachineImageSpecAtServerSpec completed");
        }
    }
}

function virtualMachineSpecListCallbackFail(error) {
    var errorMessage = error.response.data.error;
    var statusCode = error.response.status;
    commonErrorAlert(statusCode, errorMessage);
}

function getVmSpecList(sort_type) {
    getCommonVirtualMachineSpecList("virtualmachinespecmng", sort_type);

//     console.log(sort_type);
//     var url = "/setting/resources"+"/vmspec/list"
//     console.log("URL : ",url)
//     axios.get(url, {
//         headers: {
//             // 'Authorization': "{{ .apiInfo}}",
//             'Content-Type': "application/json"
//         }
//     }).then(result => {
//         console.log("get Spec List : ", result.data);
        
//         var data = result.data.VmSpecList;
//         var html = ""
        
//         console.log("data.length : ", data);

//         if (data == null) {
//             console.log("################여기##############");
//             html += '<tr><td class="overlay hidden" data-th="" colspan="5">No Data</td></tr>'

//             $("#specList").empty()
//             $("#specList").append(html)
//         } else {
//             if (sort_type) {
//                 console.log("check : ", sort_type);
//                 data.filter(list => list.name !== "").sort((a, b) => (a[sort_type] < b[sort_type] ? - 1 : a[sort_type] > b[sort_type] ? 1 : 0)).map((item, index) => (
//                     html += '<tr onclick="showVmSpecInfo(\'' + item.name + '\');">' 
//                         + '<td class="overlay hidden column-50px" data-th="">' 
//                         + '<input type="hidden" id="spec_info_' + index + '" value="' + item.name + '|' + item.connectionName + '|' + item.cspSpecName + '"/>' 
//                         + '<input type="checkbox" name="chk" value="' + item.name + '" id="raw_'  + index + '" title="" /><label for="td_ch1"></label> <span class="ov off"></span></td>' 
//                         + '<td class="btn_mtd ovm" data-th="name ">' + item.name  + '<span class="ov"></span></td>'
//                         + '<td class="overlay hidden" data-th="connectionName">' + item.connectionName + '</td>' 
//                         + '<td class="overlay hidden" data-th="cspSpecName">' + item.cspSpecName + '</td>'  
//                         + '<td class="overlay hidden column-60px" data-th=""><a href="javascript:void(0);"><img src="/assets/img/contents/icon_link.png" class="icon" alt=""/></a></td>' 
//                         + '</tr>'
//                 ))
//             } else {
//                 data.filter((list) => list.name !== "").map((item, index) => (
//                     html += '<tr onclick="showVmSpecInfo(\'' + item.name + '\');">' 
//                         + '<td class="overlay hidden column-50px" data-th="">' 
//                         + '<input type="hidden" id="spec_info_' + index + '" value="' + item.name + '"/>' 
//                         + '<input type="checkbox" name="chk" value="' + item.name + '" id="raw_'  + index + '" title="" /><label for="td_ch1"></label> <span class="ov off"></span></td>' 
//                         + '<td class="btn_mtd ovm" data-th="name ">' + item.name  + '<span class="ov"></span></td>'
//                         + '<td class="overlay hidden" data-th="connectionName">' + item.connectionName + '</td>' 
//                         + '<td class="overlay hidden" data-th="cspSpecName">' + item.cspSpecName + '</td>'  
//                         + '<td class="overlay hidden column-60px" data-th=""><a href="javascript:void(0);"><img src="/assets/img/contents/icon_link.png" class="icon" alt=""/></a></td>' 
//                         + '</tr>'
//                 ))
//             }
        
//             $("#specList").empty()
//             $("#specList").append(html)

//             // displayVmSpecInfo("REG_SUCCESS");
//         }
//     })
}

// function ModalDetail(targetAction) {
//     if( targetAction == "REG_SUCCESS" ) {
//         console.log("##########VM SPEC REG_SUCCESS")
//         $(".dashboard.register_cont").removeClass("active");
//         $(".dashboard.server_status").removeClass("view");
//         $(".dashboard .status_list tbody tr").addClass("on");
        
//         var offset = $("#vmSpecCreateBox").offset();
//         $("#wrap").animate({scrollTop : offset.top}, 0);        
        
//         // 등록 폼 초기화
//         $("#regSpecName").val('');	
//         $("#regProvider").val('');	
//         $("#regConnectionName").val('');	     
//         $("#regCspSpecName").val('');             
//     }
    
//     // $(".dashboard .status_list tbody tr").each(function () {
//     //     var $td_list = $(this),
//     //         $status = $(".server_status"),
//     //         $detail = $(".server_info");
//     //     // $td_list.off("click").click(function () {
//     //     //     $td_list.addClass("on");
//     //     //     $td_list.siblings().removeClass("on");
//     //     //     $status.addClass("view");
//     //     //     $status.siblings().removeClass("on");
//     //     //     $(".dashboard.register_cont").removeClass("active");
//     //         $td_list.off("click").click(function () {
//     //             if ($(this).hasClass("on")) {
//     //                 console.log("reg ok button click")
//     //                 $td_list.removeClass("on");
//     //                 $status.removeClass("view");
//     //                 $detail.removeClass("active");
//     //             } else {
//     //                 $td_list.addClass("on");
//     //                 $td_list.siblings().removeClass("on");
//     //                 $status.addClass("view");

//     //                 $status.siblings().removeClass("view");
//     //                 $(".dashboard.register_cont").removeClass("active");
//     //             }
//     //         });
//     //     });
//     // });
// }

function showVmSpecInfo(target) {
    console.log("target showVMSpecInfo : ", target);
    // var apiInfo = "{{ .apiInfo}}";
    var vmSpecId = encodeURIComponent(target);
    
    var url = "/setting/resources"+"/vmspec/" + vmSpecId;
    console.log("URL : ",url)
    
    return axios.get(url,{
        headers:{
            // 'Authorization': apiInfo
        }
    
    }).then(result=>{
        var data = result.data.VmSpec
        console.log("Show Data : ",data);

        var dtlSpecName = data.name;
        var dtlConnectionName = data.connectionName;
        var dtlCspSpecName = data.cspSpecName;

        $("#dtlSpecName").empty();
        $("#dtlProvider").empty();
        $("#dtlConnectionName").empty();
        $("#dtlCspSpecName").empty();
        

        $("#dtlSpecName").val(dtlSpecName);
        $("#dtlConnectionName").val(dtlConnectionName);
        $("#dtlCspSpecName").val(dtlCspSpecName);

        getProviderNameByConnection(dtlConnectionName, 'dtlProvider')// provider는 connection 정보에서 가져옴

        displayVmSpecInfo("DEL")
    }) 
}


function createVmSpec() {
    var specId = $("#regSpecName").val();
    var specName = $("#regSpecName").val();
    var connectionName = $("#regConnectionName").val();
    var cspSpecName = $("#regCspSpecName").val();
    
    if (!specName) {
        alert("Input New Spec Name")
        $("#regSpecName").focus()
        return;
    }

    // var apiInfo = "{{ .apiInfo}}";
    var url = "/setting/resources"+"/vmspec/reg"
    console.log("URL : ",url)
    var obj = {
        id: specId,
        name: specName,
        connectionName: connectionName,
        cspSpecName: cspSpecName
    }
    console.log("info image obj Data : ", obj);
    
    if (specName) {
        axios.post(url, obj, {
            headers: {
                'Content-type': 'application/json',
                // 'Authorization': apiInfo,
            }
        }).then(result => {
            console.log("result spec : ", result);
            var statusCode = result.data.status;
            if( statusCode == 200 || statusCode == 201) {
            // if (result.status == 200 || result.status == 201) {
                commonAlert("Success Create Image!!")
                //등록하고 나서 화면을 그냥 고칠 것인가?
                displayVmSpecInfo("REG_SUCCESS");
                //getVmSpecList("name");
                //아니면 화면을 리로딩 시킬것인가?
                // location.reload();
                // $("#btn_add2").click()
                // $("#namespace").val('')
                // $("#nsDesc").val('')
            } else {
                var message = result.data.message;
                commonAlert("Fail Create Spec : " + message +"(" + statusCode + ")");
                // TODO : 이 화면에서 오류날 항목은 CSP Spec Name이 없을 떄이긴 한데.... 중복일때는 알려주는데 ts.micro3(없는 spec)일 때는 어떤오류인지...
            }
        // }).catch(function(error){
        //     console.log("get create error : ");
        //     console.log(error);
        //     commonAlert(error);// TODO : error처리하자.
        // });
        }).catch((error) => {
            console.warn(error);
            console.log(error.response)
            var errorMessage = error.response.data.error;
            var statusCode = error.response.status;
            commonErrorAlert(statusCode, errorMessage);
        });
    } else {
        commonlert("Input Spec Name")
        $("#regSpecName").focus()
        return;
    }
}

function deleteVmSpec() {
    var selSpecId = "";
    var count = 0;

    $( "input[name='chk']:checked" ).each (function (){
        count++;
        selSpecId = selSpecId + $(this).val()+"," ;
    });
    selSpecId = selSpecId.substring(0,selSpecId.lastIndexOf( ","));
    
    console.log("specId : ", selSpecId);
    console.log("count : ", count);

    if(selSpecId == ''){
        alert("삭제할 대상을 선택하세요.");
        return false;
    }

    if(count != 1){
        alert("삭제할 대상을 하나만 선택하세요.");
        return false;
    }
    
    var url = "/setting/resources"+"/vmspec/del/" + selSpecId;
    console.log("URL : ",url)
    axios.delete(url, {
        headers: {
            // 'Authorization': "{{ .apiInfo}}",
            'Content-Type': "application/json"
        }
    }).then(result => {
        var data = result.data;
        console.log(data);
        // if (result.status == 200 || result.status == 201) {
        var statusCode = result.data.status;
        if( statusCode == 200 || statusCode == 201) {
            // commonAlert("Success Delete Spec.");
            commonAlert(data.message);
            // location.reload(true);
            getVmSpecList("name");
            
            displayVmSpecInfo("DEL_SUCCESS")
        } else {
            var message = data.message;
            commonAlert("Fail Create Spec : " + message +"(" + statusCode + ")");
            // TODO : 이 화면에서 오류날 항목은 CSP Spec Name이 없을 떄이긴 한데.... 중복일때는 알려주는데 ts.micro3(없는 spec)일 때는 어떤오류인지...
        }    
    }).catch((error) => {
        console.warn(error);
        console.log(error.response)
        var errorMessage = error.response.data.error;
        var statusCode = error.response.status;
        commonErrorAlert(statusCode, errorMessage);
    });
}                                                  

// connection에 등록된 spec목록 조회(공통함수 호출)
function lookupSpecList(){
    $("#assistSpecList").empty()
    var connectionName = $("#regConnectionName").val();
    if( !connectionName){
        commonAlert("connection name required")
        return;
    }

    $("#specAssist").modal();
    $('.dtbox.scrollbar-inner').scrollbar();

    getCommonLookupSpecList("vmspecmng", connectionName);
}
// 성공 callback
function lookupSpecListCallbackSuccess(caller, data){
    var html="";
    if (data == null) {
        html += '<tr><td class="overlay hidden" data-th="" colspan="5">No Data</td></tr>'

        $("#assistSpecList").empty()
        $("#assistSpecList").append(html)
    } else {
        
        // data.filter((list) => list.name !== "").map((item, index) => (
           
        //     html += '<tr onclick="showVmSpecInfo(\'' + item.name + '\');">' 
        //         + '<td class="overlay hidden column-50px" data-th="name">' + item.name + '</td>' 
        //         + '<td class="btn_mtd ovm" data-th="region ">' + item.region  + '<span class="ov"></span></td>'
        //         + '<td class="overlay hidden" data-th="mem">' + item.mem + '</td>' 
        //         + '<td class="overlay hidden" data-th="info">' 
        //         // + item.keyValueList(keyValueMap => ( keyValueMap.InstanceType))
        //         + '</td>'  
        //         + '</tr>'
               
        // ));
        // data.map((item, index) => (
           
        //     html += '<tr onclick="showVmSpecInfo(\'' + item.name + '\');">' 
        //         + '<td class="overlay hidden column-50px" data-th="name">' + item.name + '</td>' 
        //         + '<td class="btn_mtd ovm" data-th="region ">' + item.region  + '<span class="ov"></span></td>'
        //         + '<td class="overlay hidden" data-th="mem">' + item.mem + '</td>' 
        //         + '<td class="overlay hidden" data-th="info">' 
        //         // + item.keyValueList(keyValueMap => ( keyValueMap.InstanceType))
        //         + '</td>'  
        //         + '</tr>'
               
        // ));

        $.each(data, function(index, item){
            console.log('index:' + index + ' / ' + 'item:' + item);
            console.log(item);
            // keyValueMap = item.keyValueList;
            // console.log(keyValueMap);
            // var mapValue = ""
            // keyValueMap.map( (mapObj, mapIndex) => {
            //     // console.log("mapIndex = " + mapIndex);
            //     // console.log(mapObj);
            //     // console.log(mapIndex);
            //     mapValue += mapObj.Key + " : " + mapObj.Value + " <br/>";
            // });
            var vpc = item.vcpc;
            var vcpcValue = "";
            if(vpc){
                vcpcValue = 'Clock : ' + vpc.clock + '<br/> count :' + vpc.count
            }
            var gpu = item.gpu;
            var gpuValue = "";
            if(gpu){
                gpuValue += 'count : ' + (gpu.count == undefined ? "" :  gpu.count);
                gpuValue += '<br/> mem :' + (gpu.mem == undefined ? "" :  gpu.mem);
                gpuValue += '<br/> mfr :' + (gpu.mfr == undefined ? "" :  gpu.mfr);
                gpuValue += '<br/> model :' + (gpu.model == undefined ? "" :  gpu.model);
            }

            html += '<tr onclick="setCspSpecName(\'' + item.name + '\');">' 
                + '<td class="overlay hidden" data-th="region">' + item.region + '</td>' 
                + '<td class="btn_mtd ovm" data-th="name ">' + item.name  + '<span class="ov"></span></td>'
                + '<td class="btn_mtd ovm" data-th="mem ">' + item.mem  + '<span class="ov"></span></td>'
                + '<td class="overlay hidden" data-th="vcpc">' + vcpcValue + '</td>' 
                + '<td class="overlay hidden" data-th="gpu">'  + gpuValue + '</td>' 
                + '</tr>'
        });
        
        
        $("#assistSpecList").empty()
        $("#assistSpecList").append(html)
        $("#lookupSpecCount").text(data.length);
        // displayVmSpecInfo("REG_SUCCESS");
    }
}
// popup에서 main의 txtbox로 specName set
function setCspSpecName(cspSpecName){
    $("#regCspSpecName").val(cspSpecName);
    $("#specAssist").modal("hide");
}

// 조회 실패
function lookupSpecListCallbackFail(error){
    var errorMessage = error.response.data.error;
    var statusCode = error.response.status;
    commonErrorAlert(statusCode, errorMessage);
}