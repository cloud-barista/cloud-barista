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
});

/* scroll */
$(document).ready(function(){
    //checkbox all
    // $("#th_chall").click(function() {
    //     if ($("#th_chall").prop("checked")) {
    //     $("input[name=chk]").prop("checked", true);
    //     } else {
    //     $("input[name=chk]").prop("checked", false);
    //     }
    // })
    
    // //table 스크롤바 제한
    // $(window).on("load resize",function(){
    //     var vpwidth = $(window).width();
    //     if (vpwidth > 768 && vpwidth < 1800) {
    //     $(".dashboard_cont .dataTable").addClass("scrollbar-inner");
    //         $(".dataTable.scrollbar-inner").scrollbar();
    //     } else {
    //     $(".dashboard_cont .dataTable").removeClass("scrollbar-inner");
    //     }
    // });
    setTableHeightForScroll('serverImageList', 300)

    $('.btn_assist').on('click', function() {
        lookupVmImageList()
    });
});

$(document).ready(function() {
    //hidden input box
    // $('.btn_arr').click(function(){
    //     $(this).toggleClass("up");
    // if ($(".graybox.ipbox_hidden").css("display") == "none") {
    //     $(".graybox.ipbox_hidden").show();
    //   } else {
    //     $(".graybox.ipbox_hidden").hide();
    //   }
    // });	
});
$(document).ready(function () {
    // order_type = "name"
    // getVirtualMachineImageList(order_type);
    
    // var apiInfo = "{{ .apiInfo}}";
    // getCloudOS(apiInfo,'provider');
});

// function fnMove(target) {
//     var offset = $("#" + target).offset();
//     console.log("fn move offset : ", offset)
//     $('html, body').animate({
//         scrollTop: offset.top
//     }, 400);
// }

// function goFocus(target) {
//     console.log(event)
//     event.preventDefault()

//     $("#" + target).focus();
//     fnMove(target)
// }

// 등록/상세 area 보이기 숨기기
function displayVirtualMachineImageInfo(targetAction){
    if( targetAction == "REG"){
        $('#virtualMachineImageCreateBox').toggleClass("active");
        $('#virtualMachineImageInfoBox').removeClass("view");
        $('#virtualMachineImageListTable').removeClass("on");
        var offset = $("#virtualMachineImageCreateBox").offset();
        // var offset = $("#" + target+"").offset();
    	$("#TopWrap").animate({scrollTop : offset.top}, 300);

        // form 초기화
        $("#regImageName").val('')
        $("#regCspImgId").val('')
        $("#regCspImgName").val('')
        $("#regGuestOS").val('')
        $("#regDescription").val('')

        $("#assistVmImage").css("display", "block");

    }else if ( targetAction == "REG_SUCCESS"){
        $('#virtualMachineImageCreateBox').removeClass("active");
        $('#virtualMachineImageInfoBox').removeClass("view");
        $('#virtualMachineImageListTable').addClass("on");
        
        // form 초기화
        $("#regImageName").val('')
        $("#regCspImgId").val('')
        $("#regCspImgName").val('')
        $("#regGuestOS").val('')
        $("#regProvider").val('')
        $("#regConnectionName").val('')
        $("#regDescription").val('')
        
        var offset = $("#virtualMachineImageCreateBox").offset();
        $("#TopWrap").animate({scrollTop : offset.top}, 0);
        
        getVirtualMachineImageList("name");
    }else if ( targetAction == "DEL"){
        $('#virtualMachineImageCreateBox').removeClass("active");
        $('#virtualMachineImageInfoBox').addClass("view");
        $('#virtualMachineImageListTable').removeClass("on");

        var offset = $("#virtualMachineImageInfoBox").offset();
    	$("#TopWrap").animate({scrollTop : offset.top}, 300);

    }else if ( targetAction == "DEL_SUCCESS"){
        $('#virtualMachineImageCreateBox').removeClass("active");
        $('#virtualMachineImageInfoBox').removeClass("view");
        $('#virtualMachineImageListTable').addClass("on");

        var offset = $("#virtualMachineImageInfoBox").offset();
        $("#TopWrap").animate({scrollTop : offset.top}, 0);

        getVirtualMachineImageList("name");
    }else if ( targetAction == "CLOSE"){
        $('#virtualMachineImageCreateBox').removeClass("active");
        $('#virtualMachineImageInfoBox').removeClass("view");
        $('#virtualMachineImageListTable').addClass("on");

        var offset = $("#virtualMachineImageInfoBox").offset();
        $("#TopWrap").animate({scrollTop : offset.top}, 0);
    }
}

function deleteVirtualMachineImage() {
    var imageId = "";
    var count = 0;

    $( "input[name='chk']:checked" ).each (function (){
        count++;
        imageId = imageId + $(this).val()+"," ;
    });
    imageId = imageId.substring(0,imageId.lastIndexOf( ","));
    
    console.log("imageId : ", imageId);
    console.log("count : ", count);

    if(imageId == ''){
        commonAlert("삭제할 대상을 선택하세요.");
        return false;
    }

    if(count != 1){
        commonAlert("삭제할 대상을 하나만 선택하세요.");
        return false;
    }

    // var url = CommonURL + "/ns/" + NAMESPACE + "/resources/image/" + imageId;
    //var u = SpiderURL + "/vmimage/" + imageId;
    var url = "/setting/resources" + "/machineimage/del/" + imageId

    axios.delete(url, {
        headers: {
            // 'Authorization': "{{ .apiInfo}}",
            'Content-Type': "application/json"
        }
    }).then(result => {
        var data = result.data;
        console.log(data);
        if (result.status == 200 || result.status == 201) {
            // commonAlert("Success Delete Image.");
            commonAlert(data.message);
            displayVirtualMachineImageInfo("DEL_SUCCESS")
            // location.reload(true);
        }else{
            commonAlert(data.error);
        }
    // }).catch(function(error){
    //     console.log("image delete error : ",error);        
    // });
    }).catch((error) => {
        console.warn(error);
        console.log(error.response)
        var errorMessage = error.response.data.error;
        var statusCode = error.response.status;
        commonErrorAlert(statusCode, errorMessage);
    });
}                                                  

function getVirtualMachineImageList(sort_type) {
    console.log("여기 : ", sort_type);
    var sortType = sort_type;
    console.log("sortType : ", sortType);

    getCommonVirtualMachineImageList("virtualmachineimagemng", sortType)
    // console.log(sort_type);
    // // var url = CommonURL + "/ns/" + NAMESPACE + "/resources/image";
    // var url = "/setting/resources" + "/machineimage/list"
    // axios.get(url, {
    //     headers: {
    //         // 'Authorization': "{{ .apiInfo}}",
    //         'Content-Type': "application/json"
    //     }
    // }).then(result => {
    //     console.log("get Image List : ", result.data);
        
    //     var data = result.data.VirtualMachineImageList;
    //     var html = ""
        
    //     if (data.length) {
    //         if (sort_type) {
    //             console.log("check : ", sort_type);
    //             data.filter(list => list.name !== "").sort((a, b) => (a[sort_type] < b[sort_type] ? - 1 : a[sort_type] > b[sort_type] ? 1 : 0)).map((item, index) => (
    //                 html += '<tr onclick="showVirtualMachinImageInfo(\'' + item.name + '\');">' 
    //                     + '<td class="overlay hidden" data-th="">' 
    //                     + '<input type="hidden" id="img_info_' + index + '" value="' + item.name + '|' + item.cspImageId + '"/>' 
    //                     + '<input type="checkbox" name="chk" value="' + item.name + '" id="raw_'  + index + '" title="" /><label for="td_ch1"></label> <span class="ov off"></span></td>' 
    //                     + '<td class="btn_mtd ovm" data-th="cspImageId ">' + item.cspImageId  + '<span class="ov"></span></td>'
    //                     + '<td class="overlay hidden" data-th="name">' + item.name + '</td>' 
    //                     + '<td class="overlay hidden" data-th="description">' + item.description + '</td>'  
    //                     + '<td class="overlay hidden" data-th=""><a href="javascript:void(0);"><img src="/assets/img/contents/icon_link.png" class="icon" alt=""/></a></td>' 
    //                     + '</tr>'
    //             ))
    //         } else {
    //             data.filter((list) => list.name !== "").map((item, index) => (
    //                 html += '<tr onclick="showVirtualMachinImageInfo(\'' + item.name + '\');">' 
    //                     + '<td class="overlay hidden" data-th="">' 
    //                     + '<input type="hidden" id="img_info_' + index + '" value="' + item.name  + '"/>'
    //                     + '<input type="checkbox" name="chk" value="' + item.name + '" id="raw_' + index + '" title="" /><label for="td_ch1"></label> <span class="ov off"></span></td>' 
    //                     + '<td class="btn_mtd ovm" data-th="cspImageId">' + item.cspImageId + '<span class="ov"></span></td>' 
    //                     + '<td class="overlay hidden" data-th="name">' + item.name + '</td>' 
    //                     + '<td class="overlay hidden" data-th="description">' + item.description + '</td>' 
    //                     + '<td class="overlay hidden" data-th=""><a href="javascript:void(0);"><img src="/assets/img/contents/icon_link.png" class="icon" alt=""/></a></td>' 
    //                     + '</tr>'
    //             ))
    //         }
        
    //         $("#imgList").empty()
    //         $("#imgList").append(html)
            
    //         ModalDetail()
    //     }
    // }).catch(function(error){
    //     console.log("list error : ",error);        
    // });
}

function setVirtualMachineImageListAtServerImage(data, sortType){
    var html = "";
    console.log(data);// TODO : no Data 일 때 (ex. 모든 이미지 삭제 등) data.length에서 오류 남.
    console.log("here Set : ", sortType);
    
    if (data == null) {
        html += '<tr><td class="overlay hidden" data-th="" colspan="5">No Data</td></tr>'

        $("#imgList").empty()
        $("#imgList").append(html)
        //console.log("setVirtualMachineImageListAtServerImage completed");
        ModalDetail()
    } else {
        if (data.length) {
            if (sortType) {
                console.log("check : ", sortType);
                data.filter(list => list.name !== "").sort((a, b) => (a[sortType] < b[sortType] ? - 1 : a[sortType] > b[sortType] ? 1 : 0)).map((item, index) => (
                    html += '<tr onclick="showVirtualMachinImageInfo(\'' + item.name + '\');">' 
                        + '<td class="overlay hidden column-50px" data-th="">' 
                        + '<input type="hidden" id="img_info_' + index + '" value="' + item.name + '|' + item.cspImageId + '"/>' 
                        + '<input type="checkbox" name="chk" value="' + item.name + '" id="raw_'  + index + '" title="" /><label for="td_ch1"></label> <span class="ov off"></span></td>' 
                        + '<td class="btn_mtd ovm" data-th="connectionName ">' + item.connectionName  + '<span class="ov"></span></td>'
                        + '<td class="btn_mtd ovm" data-th="cspImageId ">' + item.cspImageId  + '<span class="ov"></span></td>'
                        + '<td class="overlay hidden" data-th="name">' + item.name + '</td>' 
                        + '<td class="overlay hidden" data-th="description">' + item.description + '</td>'  
                        // + '<td class="overlay hidden column-60px" data-th=""><a href="javascript:void(0);"><img src="/assets/img/contents/icon_link.png" class="icon" alt=""/></a></td>' 
                        + '</tr>'
                ))
            } else {
                data.filter((list) => list.name !== "").map((item, index) => (
                    html += '<tr onclick="showVirtualMachinImageInfo(\'' + item.name + '\');">' 
                        + '<td class="overlay hidden column-50px" data-th="">' 
                        + '<input type="hidden" id="img_info_' + index + '" value="' + item.name  + '"/>'
                        + '<input type="checkbox" name="chk" value="' + item.name + '" id="raw_' + index + '" title="" /><label for="td_ch1"></label> <span class="ov off"></span></td>' 
                        + '<td class="btn_mtd ovm" data-th="cspImageId">' + item.cspImageId + '<span class="ov"></span></td>' 
                        + '<td class="overlay hidden" data-th="name">' + item.name + '</td>' 
                        + '<td class="overlay hidden" data-th="description">' + item.description + '</td>' 
                        // + '<td class="overlay hidden column-60px" data-th=""><a href="javascript:void(0);"><img src="/assets/img/contents/icon_link.png" class="icon" alt=""/></a></td>' 
                        + '</tr>'
                ))
            }
        
            $("#imgList").empty()
            $("#imgList").append(html)
            console.log("setVirtualMachineImageListAtServerImage completed");
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

// function getConnectionInfo(provider){
//     var url = SpiderURL+"/connectionconfig";
//     console.log("provider : ",provider)
//     //var provider = $("#provider option:selected").val();
//     var html = "";
//     var apiInfo = ApiInfo
//     axios.get(url,{
//         headers:{
//             'Authorization': apiInfo
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

//     })
// }

function createVirtualMachineImage() {
    var imgId = $("#regImageName").val();
    var imgName = $("#regImageName").val();
    var cspImgId = $("#regCspImgId").val();
    var guestOS = $("#regGuestOS").val();
    var connectionName = $("#regConnectionName").val();
    var description = $("#regDescription").val();
    
    var cspImgName = "";
    if(!cspImgName) {
        $("#regCspImgName").val();
    }
    
    console.log("check obj : " + imgId + ", " + imgName + ", " + cspImgId + ", " + cspImgName + ", " + guestOS + ", " + connectionName + ", " + description);
    
    if (!imgName) {
        commonAlert("Input New Image Name")
        $("#regImageName").focus()
        return;
    }

    // var apiInfo = "{{ .apiInfo}}";
    // var url = CommonURL+"/ns/"+NAMESPACE+"/resources/image?action=registerWithInfo"
    var url = "/setting/resources" + "/machineimage/reg"
    console.log("Image Reg URL : ",url)
    var obj = {
        connectionName: connectionName,
        cspImageId: cspImgId,
        cspImageName: cspImgName,
        description: description,
        guestOS: guestOS,
        id: imgId,
        name: imgName
    }
    console.log("info image obj Data : ", obj);
    
    if (imgName) {
        axios.post(url, obj, {
            headers: {
                'Content-type': 'application/json',
                // 'Authorization': apiInfo,
            }
        }).then(result => {
            console.log("result image : ", result);
            if (result.status == 200 || result.status == 201) {
                commonAlert("Success Create Image!!")
                //등록하고 나서 화면을 그냥 고칠 것인가?
                displayVirtualMachineImageInfo("REG_SUCCESS")
                //아니면 화면을 리로딩 시킬것인가?
                // location.reload();
                // $("#btn_add2").click()
                // $("#namespace").val('')
                // $("#nsDesc").val('')
            } else {
                commonAlert("Fail Create Image)")
            }
        // }).catch(function(error){
        //     console.log("image create error : ",error);        
        // });
        }).catch((error) => {
            console.warn(error);
            console.log(error.response)
            var errorMessage = error.response.data.error;
            var statusCode = error.response.status;
            commonErrorAlert(statusCode, errorMessage);
        });
    } else {
        commonAlert("Input Image Name")
        $("#regImageName").focus()
        return;
    }
}

// vmImage의 상세정보 표시.
function showVirtualMachinImageInfo(target) {
    console.log("target showInfo : ", target);
    // var apiInfo = "{{ .apiInfo}}";
    //var imageId = encodeURIComponent(target);
    $('.stxt').html(target);
    $("#assistVmImage").css("display", "none");
    
    // var url = CommonURL+"/ns/"+NAMESPACE+"/resources/image/"+ imageId;
    // var url = "/setting/resources/machineimage/" + imageId
    var url = "/setting/resources/machineimage/" + target
    console.log("image detail URL : ",url)

    return axios.get(url,{
        headers:{
            // 'Authorization': apiInfo
        }
    
    }).then(result=>{
        var data = result.data.VirtualMachineImageInfo;
        console.log("Show Data : ",data);
        var dtlImageName = data.name;
        var dtlConnectionName = data.connectionName;
        var dtlImageId = data.id;
        var dtlGuestOS = data.guestOS;
        var dtlDescription = data.description;
        var dtlCspImageName = data.cspImageName;
        var dtlCspImageId = data.cspImageId;

        $("#dtlImageName").empty();
        $("#dtlProvider").empty();
        $("#dtlConnectionName").empty();
        $("#dtlImageId").empty();
        $("#dtlGuestOS").empty();
        $("#dtlDescription").empty();
        $("#dtlCspImageName").empty();
        $("#dtlCspImageId").empty();

        $("#dtlImageName").val(dtlImageName);
        $("#dtlConnectionName").val(dtlConnectionName);
        $("#dtlImageId").val(dtlImageId);
        $("#dtlGuestOS").val(dtlGuestOS);
        $("#dtlDescription").val(dtlDescription);
        $("#dtlCspImageName").val(dtlCspImageName);
        $("#dtlCspImageId").val(dtlCspImageId);

        // getProvider(dtlConnectionName);
        getProviderNameByConnection(dtlConnectionName, 'dtlProvider')// provider는 connection 정보에서 가져옴
    // }).catch(function(error){
    //     console.log("image data error : ",error);        
    // });
    }).catch((error) => {
        console.warn(error);
        console.log(error.response)
        var errorMessage = error.response.data.error;
        var statusCode = error.response.status;
        commonErrorAlert(statusCode, errorMessage);
    });
}

// function getProvider(target) {
//     console.log("getProvidergetProvider : ",target);
//     var url = SpiderURL+"/connectionconfig/" + target;
        
//     return axios.get(url,{
//         headers:{
//             'Authorization': apiInfo
//         }
    
//     }).then(result=>{
//         var data = result.data;
        
//         var Provider = data.ProviderName;

//         $("#dtlProvider").val(Provider);
//     })        
// }

// connection에 등록된 spec목록 조회(공통함수 호출)
function lookupVmImageList(){
    $("#assistVmImageList").empty()
    // connection과 상관없이 조회 가능
    var connectionName = $("#regConnectionName").val();
    if( !connectionName){
        commonAlert("connection name required")
        return;
    }

    $("#imageAssist").modal();
    $('.dtbox.scrollbar-inner').scrollbar();

    getCommonLookupImageList("vmimagemng", connectionName);
}
// 성공 callback
function lookupVmImageListCallbackSuccess(caller, data){
    var html="";
    if (data == null) {
        html += '<tr><td class="overlay hidden" data-th="" colspan="5">No Data</td></tr>'

        $("#assistVmImageList").empty()
        $("#assistVmImageList").append(html)
    } else {
        
        $.each(data, function(index, item){
            console.log('index:' + index + ' / ' + 'item:' + item);
            console.log(item);
            keyValueMap = item.keyValueList;
            // console.log(keyValueMap);
            var mapValue = ""
            var mapName = "";
            keyValueMap.map( (mapObj, mapIndex) => {
                if( mapObj.Key == "Name"){
                    mapName = mapObj.Value
                    return
                }
                // console.log("mapIndex = " + mapIndex);
                // console.log(mapObj);
                // console.log(mapIndex);
                // mapValue += mapObj.Key + " : " + mapObj.Value + " <br/>";
            });

            var imageName = ""
            var iid = item.iid;
            var iidNameID = "";
            var iidSystemID = "";

            if( item.name == undefined || item.name == ""){
                imageName = mapName;
            }else{
                imageName = item.name;
            }
            if(iid){
                iidNameID = (iid.NameId == undefined || item.NameId == "" ? "" :  iid.NameId);
                iidSystemID = (iid.SystemId == undefined || item.SystemId == "" ? "" :  iid.SystemId);
            }
//cspImageNameID, cspImageName, cspImageGuestOS
            html += '<tr onclick="setCspVmImageInfo(\''+iidNameID+'\',\''+imageName+'\',\''+item.guestOS+'\');">'            
                + '<td class="overlay hidden" data-th="name">' + imageName + '</td>' 
                + '<td class="btn_mtd ovm" data-th="status ">' + item.status  + '<span class="ov"></span></td>'
                + '<td class="btn_mtd ovm" data-th="guestOS ">' + item.guestOS  + '<span class="ov"></span></td>'
                + '<td class="overlay hidden" data-th="vcpc">' + iidNameID + '</td>' 
                + '<td class="overlay hidden" data-th="gpu">'  + iidSystemID + '</td>' 
                + '</tr>'
        });
        
        
        $("#assistVmImageList").empty()
        $("#assistVmImageList").append(html)
        $("#lookupVmImageCount").text(data.length);
    }
}


// 조회 실패
function lookupVmImageListCallbackFail(error){
    var errorMessage = error.response.data.error;
    var statusCode = error.response.status;
    commonErrorAlert(statusCode, errorMessage);
}

// popup에서 main의 txtbox로 set
function setCspVmImageInfo(cspImageNameID, cspImageName, cspImageGuestOS){
    $("#regCspImgId").val(cspImageNameID);
    $("#regCspImgName").val(cspImageName);
    $("#regGuestOS").val(cspImageGuestOS);

    $("#imageAssist").modal("hide");
}

// connection 정보가 바뀌면 image 정보도 초기화 시킨다.
function clearCspImageInfo(){
    $("#regCspImgId").val();
    $("#regCspImgName").val();
    $("#regGuestOS").val();
}
