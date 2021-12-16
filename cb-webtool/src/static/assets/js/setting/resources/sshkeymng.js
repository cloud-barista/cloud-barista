$(document).ready(function () {
    //action register open / table view close
    // $('#RegistBox .btn_ok.register').click(function () {
    //     $(".dashboard.register_cont").toggleClass("active");
    //     $(".dashboard.server_status").removeClass("view");
    //     $(".dashboard .status_list tbody tr").removeClass("on");
    //     //ok 위치이동
    //     $('#RegistBox').on('hidden.bs.modal', function () {
    //         var offset = $("#CreateBox").offset();
    //         $("#wrap").animate({
    //             scrollTop: offset.top
    //         }, 300);
    //     })
    // });
});
/* scroll */
$(document).ready(function () {
    //checkbox all
    // $("#th_chall").click(function () {
    //     if ($("#th_chall").prop("checked")) {
    //         $("input[name=chk]").prop("checked", true);
    //     } else {
    //         $("input[name=chk]").prop("checked", false);
    //     }
    // })

    // //table 스크롤바 제한
    // $(window).on("load resize", function () {
    //     var vpwidth = $(window).width();
    //     if (vpwidth > 768 && vpwidth < 1800) {
    //         $(".dashboard_cont .dataTable").addClass("scrollbar-inner");
    //         $(".dataTable.scrollbar-inner").scrollbar();
    //     } else {
    //         $(".dashboard_cont .dataTable").removeClass("scrollbar-inner");
    //     }
    // });

    setTableHeightForScroll('sshkeyList', 300)
});

$(document).ready(function () {
    // order_type = "name"
    // getSSHKeyList(order_type);

    // var apiInfo = "{{ .apiInfo}}";
    // getCloudOS(apiInfo,'provider');
})


// function fnMove(target) {
//     var offset = $("#" + target).offset();
//     console.log("fn move offset : ", offset)
//     $('html, body').animate({
//         scrollTop: offset.top
//     }, 400);
// }

// area 표시
function displaySshKeyInfo(targetAction) {
    if (targetAction == "REG") {
        $('#sshKeyCreateBox').toggleClass("active");
        $('#sskKeyInfoBox').removeClass("view");
        $('#sshKeyListTable').removeClass("on");
        var offset = $("#sshKeyCreateBox").offset();
        // var offset = $("#" + target+"").offset();
        $("#TopWrap").animate({ scrollTop: offset.top }, 300);

        // form 초기화
        $("#regCspSshKeyName").val('');
        //$("#regProvider").val('');
        //$("#regCregConnectionNameidrBlock").val('');
        goFocus('sshKeyCreateBox');
    } else if (targetAction == "REG_SUCCESS") {
        $('#sshKeyCreateBox').removeClass("active");
        $('#sskKeyInfoBox').removeClass("view");
        $('#sshKeyListTable').addClass("on");

        var offset = $("#sshKeyCreateBox").offset();
        $("#TopWrap").animate({ scrollTop: offset.top }, 0);

        // form 초기화
        $("#regCspSshKeyName").val('');
        $("#regProvider").val('');
        $("#regCregConnectionNameidrBlock").val('');

        getSshKeyList("name");
    } else if (targetAction == "DEL") {
        $('#sshKeyCreateBox').removeClass("active");
        $('#sskKeyInfoBox').addClass("view");
        $('#sshKeyListTable').removeClass("on");

        var offset = $("#sskKeyInfoBox").offset();
        $("#TopWrap").animate({ scrollTop: offset.top }, 300);

    } else if (targetAction == "DEL_SUCCESS") {
        $('#sshKeyCreateBox').removeClass("active");
        $('#sskKeyInfoBox').removeClass("view");
        $('#sshKeyListTable').addClass("on");

        var offset = $("#sskKeyInfoBox").offset();
        $("#TopWrap").animate({ scrollTop: offset.top }, 0);

        getSshKeyList("name");
    } else if (targetAction == "CLOSE") {
        $('#sshKeyCreateBox').removeClass("active");
        $('#sskKeyInfoBox').removeClass("view");
        $('#sshKeyListTable').addClass("on");

        var offset = $("#sskKeyInfoBox").offset();
        $("#TopWrap").animate({ scrollTop: offset.top }, 0);
    }
}

// SshKey 목록 조회
function getSshKeyList(sort_type) {
    //var url = "{{ .comURL.SpiderURL}}" + "/connectionconfig";
    // var url = CommonURL+"/ns/"+NAMESPACE+"/resources/sshKey";
    var url = "/setting/resources" + "/sshkey/list"
    axios.get(url, {
        headers: {
            // 'Authorization': "{{ .apiInfo}}",
            'Content-Type': "application/json"
        }
    }).then(result => {
        console.log("get SSH Data : ", result.data);
        var data = result.data.SshKeyList; // exception case : if null 
        var html = ""

        if (data == null) {
            html += '<tr><td class="overlay hidden" data-th="" colspan="5">No Data</td></tr>'

            $("#sList").empty();
            $("#sList").append(html);

            ModalDetail()
        } else {
            if (data.length) { // null exception if not exist
                if (sort_type) {
                    console.log("check : ", sort_type);
                    data.filter(list => list.name !== "").sort((a, b) => (a[sort_type] < b[sort_type] ? - 1 : a[sort_type] > b[sort_type] ? 1 : 0)).map((item, index) => (
                        //html += '<tr onclick="showSshKeyInfo(\'' + item.cspSshKeyName + '\');">'
                        html += '<tr onclick="showSshKeyInfo(\'' + item.id + '\');">'
                        + '<td class="overlay hidden column-50px" data-th="">'
                        + '<input type="hidden" id="ssh_info_' + index + '" value="' + item.name + '|' + item.connectionName + '|' + item.cspSshKeyName + '"/>'
                        + '<input type="checkbox" name="chk" value="' + item.name + '" id="raw_' + index + '" title="" /><label for="td_ch1"></label> <span class="ov off"></span></td>'
                        + '<td class="btn_mtd ovm" data-th="Name">' + item.id
                        // + '<a href="javascript:void(0);"><img src="/assets/img/contents/icon_copy.png" class="td_icon" alt=""/></a> <span class="ov"></span></td>'
                        + '</td>'
                        + '<td class="overlay hidden" data-th="connectionName">' + item.connectionName + '</td>'
                        + '<td class="overlay hidden" data-th="cspSshKeyName">' + item.cspSshKeyName + '</td>'
                        // + '<td class="overlay hidden column-60px" data-th=""><a href="javascript:void(0);"><img src="/assets/img/contents/icon_link.png" class="icon" alt=""/></a></td>' 
                        + '</tr>'
                    ))
                } else {
                    data.filter((list) => list.name !== "").map((item, index) => (
                        //html += '<tr onclick="showSshKeyInfo(\'' + item.cspSshKeyName + '\');">'
                        html += '<tr onclick="showSshKeyInfo(\'' + item.id + '\');">'
                        + '<td class="overlay hidden column-50px" data-th="">'
                        + '<input type="hidden" id="ssh_info_' + index + '" value="' + item.name + '"/>'
                        + '<input type="checkbox" name="chk" value="' + item.name + '" id="raw_' + index + '" title="" /><label for="td_ch1"></label> <span class="ov off"></span></td>'
                        + '<td class="btn_mtd ovm" data-th="id">' + item.id + '<span class="ov"></span></td>'
                        + '<td class="overlay hidden" data-th="connectionName">' + item.connectionName + '</td>'
                        + '<td class="overlay hidden" data-th="cspSshKeyName">' + item.cspSshKeyName + '</td>'
                        // + '<td class="overlay hidden column-60px" data-th=""><a href="javascript:void(0);"><img src="/assets/img/contents/icon_link.png" class="icon" alt=""/></a></td>' 
                        + '</tr>'
                    ))

                }

                $("#sList").empty();
                $("#sList").append(html);

                ModalDetail()

            }
        }

        // }).catch(function(error){
        //     console.log("get sshKeyList error : ",error);        
        // });
    }).catch((error) => {
        console.warn(error);
        console.log(error.response)
        var errorMessage = error.response.data.error;
        var statusCode = error.response.status;
        commonErrorAlert(statusCode, errorMessage);
    });
}

// function goFocus(target) {

//     console.log(event)
//     event.preventDefault()

//     $("#" + target).focus();
//     fnMove(target)
// }

function deleteSshKey() {
    // function goDelete() {
    var selSshKeyId = "";
    var count = 0;

    $("input[name='chk']:checked").each(function () {
        count++;
        selSshKeyId = selSshKeyId + $(this).val() + ",";
    });
    selSshKeyId = selSshKeyId.substring(0, selSshKeyId.lastIndexOf(","));

    console.log("sshKeyId : ", selSshKeyId);
    console.log("count : ", count);

    if (selSshKeyId == '') {
        alert("삭제할 대상을 선택하세요.");
        return false;
    }

    if (count != 1) {
        alert("삭제할 대상을 하나만 선택하세요.");
        return false;
    }

    // var url = CommonURL + "/ns/" + NAMESPACE + "/resources/sshKey/" + selSshKeyId;
    var url = "" + "" + selSshKeyId;
    var url = "/setting/resources" + "/sshkey/del/" + selSshKeyId;
    axios.delete(url, {
        headers: {
            // 'Authorization': "{{ .apiInfo}}",
            'Content-Type': "application/json"
        }
    }).then(result => {
        var data = result.data;
        console.log(data);
        if (result.status == 200 || result.status == 201) {
            // commonAlert("Success Delete SSH Key.");
            commonAlert(data.message);
            // location.reload(true);

            displaySshKeyInfo("DEL_SUCCESS");
            //getSshKeyList("name");

            getSshKeyList("name");
        } else {
            commonAlert(data.error);
        }
        // }).catch(function(error){
        //     console.log("get delete error : ",error);        
        // });
    }).catch((error) => {
        console.warn(error);
        console.log(error.response)
        var errorMessage = error.response.data.error;
        var statusCode = error.response.status;
        commonErrorAlert(statusCode, errorMessage);
    });
}

function showSshKeyInfo(sshKeyId) {
    console.log("target showSshKeyInfo : ", sshKeyId);
    // var sshKeyId = target;
    // var apiInfo = "{{ .apiInfo}}";
    // var url = CommonURL+"/ns/"+NAMESPACE+"/resources/sshKey/"+ sshKeyId;
    var url = "/setting/resources" + "/sshkey/" + sshKeyId;
    console.log("ssh key URL : ", url)

    return axios.get(url, {
        headers: {
            // 'Authorization': apiInfo
        }

    }).then(result => {
        var data = result.data.SshKeyInfo
        console.log("Show Data : ", data);

        var dtlCspSshKeyName = data.cspSshKeyName;
        var dtlDescription = data.description;
        var dtlUserID = data.userID;
        var dtlConnectionName = data.connectionName;
        var dtlPublicKey = data.publicKey;
        var dtlPrivateKey = data.privateKey;
        var dtlFingerprint = data.fingerprint;


        $('#dtlCspSshKeyName').empty();
        $('#dtlDescription').empty();
        $('#dtlUserID').empty();
        $('#dtlConnectionName').empty();
        $('#dtlPublicKey').empty();
        $('#dtlPrivateKey').empty();
        $('#dtlFingerprint').empty();

        $('#dtlCspSshKeyName').val(dtlCspSshKeyName);
        $('#dtlDescription').val(dtlDescription);
        $('#dtlUserID').val(dtlUserID);
        $('#dtlConnectionName').val(dtlConnectionName);
        $('#dtlPublicKey').val(dtlPublicKey);
        $('#dtlPrivateKey').val(dtlPrivateKey);
        $('#dtlFingerprint').val(dtlFingerprint);
        // }).catch(function(error){
        //     console.log("get sshKey error : ",error);        
        // });
    }).catch((error) => {
        console.warn(error);
        console.log(error.response)
        var errorMessage = error.response.data.error;
        var statusCode = error.response.status;
        commonErrorAlert(statusCode, errorMessage);
    });
}

function createSSHKey() {
    var cspSshKeyName = $("#regCspSshKeyName").val()
    var connectionName = $("#regConnectionName").val()

    console.log("info param cspSshKeyName : ", cspSshKeyName);
    console.log("info param connectionName : ", connectionName);

    if (!cspSshKeyName) {
        alert("Input New SSH Key Name")
        $("#regCspSshKeyName").focus()
        return;
    }
    if (!connectionName) {
        alert("Input Connection Name")
        $("#regConnectionName").focus()
        return;
    }

    // var apiInfo = "{{ .apiInfo}}";
    // var url = CommonURL+"/ns/"+NAMESPACE+"/resources/sshKey"
    var url = "" + "";
    var url = "/setting/resources" + "/sshkey/reg"
    console.log("ssh key URL : ", url)
    var obj = {
        name: cspSshKeyName,
        connectionName: connectionName
    }
    console.log("info connectionconfig obj Data : ", obj);
    if (cspSshKeyName) {
        axios.post(url, obj, {
            headers: {
                'Content-type': 'application/json',
                // 'Authorization': apiInfo,
            }
        }).then(result => {
            console.log(result);
            if (result.status == 200 || result.status == 201) {
                commonAlert("Success Create SSH Key")
                //등록하고 나서 화면을 그냥 고칠 것인가?
                displaySshKeyInfo("REG_SUCCESS");
                //getSshKeyList("name");
                //아니면 화면을 리로딩 시킬것인가?
                // location.reload();
                // $("#btn_add2").click()
                // $("#namespace").val('')
                // $("#nsDesc").val('')
            } else {
                commonAlert("Fail Create SSH Key")
            }
            // }).catch(function(error){
            //     console.log("get create error : ",error);        
            // });
        }).catch((error) => {
            console.warn(error);
            console.log(error.response)
            var errorMessage = error.response.statusText;
            var statusCode = error.response.status;
            commonErrorAlert(statusCode, errorMessage);
        });
    } else {
        commonAlert("Input SSH Key Name")
        $("#regCspSshKeyName").focus()
        return;
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
//     // var url = SpiderURL+"/connectionconfig";
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