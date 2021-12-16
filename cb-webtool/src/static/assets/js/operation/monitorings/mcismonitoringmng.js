$(document).ready(function() {
    getCommonMcisList("mcismonitoringmng", true, '', "id")
    resizeContent();
});
$(window).resize(function() {
    resizeContent();
});
//Selected MCIS selectbox(CPUs,Memory,DiskIO,Network) width 반응형 적용
function resizeContent() {
    $(".g_list .gbox .sel").each(function(){
        var $list =  $(this),
                $label =  $list.find('label'),
                $labelWidth = $label.width(),
                $gboxWidth = $list.width(),
                $selectbox =  $list.find('.selectbox');
        $list.each(function(){
            $selectbox.css({'width':($gboxWidth-$labelWidth-20)+'px'});
        });
    });
}

$(document).ready(function(){
    // page Load 시점에 이미 가져옴.
    // checkNS();
    // var nsid = "{{ .NameSpace}}";
    // // console.log("ready Monitoring Mcis, nsid :",nsid);
    // getMcisList(nsid);
    // // show_mcis(url);
})

// <option value="{ {$item.ID} }"  selected>{ {$item.Name} }|{ {$item.Status} }|{ {$item.Description} }-->
// MCIS 목록 조회 후 화면에 Set
function getMcisListCallbackSuccess(caller, mcisList){

    // MCIS Status
    var addMcis = "";
    addMcis += '<option>Choose a Target MCIS for Monitoring</option>';
    if(!isEmpty(mcisList) && mcisList.length > 0 ){
        var initMcis = $("#init_mcis").val();
        var mcisExist = false;// monitoring할 mcis가 없을 수도 있음.
        console.log(mcisList)
        for(var i in mcisList){
            // if(i == 0 ){
            //     addMcis +='<option value="'+mcisList[i].id+'"  selected>'+mcisList[i].name+"|"+mcisList[i].status+"|"+mcisList[i].description
            // }else{
            //     addMcis +='<option value="'+mcisList[i].id+'" >'+mcisList[i].name+"|"+mcisList[i].status+"|"+mcisList[i].description
            // }
            addMcis +='<option value="'+mcisList[i]+'">'+mcisList[i] + '</option>';
            if( initMcis == mcisList[i]){
                mcisExist = true;
            }
        }
        $("#mcisList").empty()
        $("#mcisList").append(addMcis)
        if (initMcis && mcisExist) {
            console.log("initMcis = " + initMcis)
        }else{
            console.log("initMcis is not exists");
            initMcis = mcisList[0]
        }

        // id="mcisList"
        $("#mcisList").val(initMcis).prop("selected", true);
        selectMonitoringMcis(initMcis)
    }else{
        var addMcis = "";

        // $("#mcisList").append(addMcis);
    }
}

// 조회 실패시.
function getMcisListCallbackFail(caller, error){
    // List table에 no data 표시? 또는 조회 오류를 표시?
}

// function getMcisList(nsid){
//     var url = "{{ .comURL.TumbleBugURL}}"+"/ns/"+nsid+"/mcis";
//     var apiInfo = "{{ .apiInfo}}";
//     axios.get(url,{
//         headers:{
//             'Authorization': apiInfo,
//             'Content-Type' : "application/json"
//         }
//     }).then(result=>{
//         var data = result.data.mcis;
//         console.log("getMCISList data: ",data);
//         var init_mcis = '';
//         if (data.length == 0) {
//             console.warn("data is empty, response data : ", data);
//
//         } else {
//             init_mcis = data[0].id
//         }
//         console.log("init mcis : ", init_mcis);
//         var s_mcis_list = "<option value=''>Choose a Target MCIS for Monitoring</option>"
//         for(var i in data){
//             if(i == 0 ){
//                 s_mcis_list +='<option value="'+data[i].id+'"  selected>'+data[i].name+"|"+data[i].status+"|"+data[i].description
//             }else{
//                 s_mcis_list +='<option value="'+data[i].id+'" >'+data[i].name+"|"+data[i].status+"|"+data[i].description
//             }
//         }
//
//
//         $("#mcisList").empty()
//         $("#mcisList").append(s_mcis_list)
//         if (init_mcis != "") {
//             selectMonitoringMcis(init_mcis)
//         }
//     })
// }

function selectMonitoringMcis(mcisId){
    $("#mcis_id").val(mcisId);
    // var nsid = NAMESPACE
    // var url = "{{ .comURL.TumbleBugURL}}"+"/ns/"+nsid+"/mcis/"+mcis_id;
    // var apiInfo = "{{ .apiInfo}}";

    getCommonMcisData("mcismonitoringmng", mcisId)
    // var url = "/operation/manages/mcismng/" + mcisId
    // axios.get(url,{
    //     headers:{
    //         // 'Authorization': apiInfo,
    //         'Content-Type' : "application/json"
    //     }
    // }).then(result=>{
    //     console.log("selectMonitoringMcis result")
    //     var mcis = result.data.McisInfo
    //     console.log(mcis)
    //     var vms = mcis.vm
    //     console.log(vms)
    //     var vm_badge ="";
    //     var vm_options = "";
    //     var init_vm = vms[0].id
    //     if(vms){
    //         vm_len = vms.length
    //         for(var o in vms){
    //             var vm_status = vms[o].status
    //             vm_options +='<option value="'+vms[o].id+'">'+vms[o].name+'|'+vms[o].status+'|'+vms[o].description
    //
    //             var vmStatusIcon = "bgbox_b";
    //             if(vm_status == "Running"){
    //                 vmStatusIcon = "bgbox_b";
    //             }else if(vm_status == "include" ){
    //                 vmStatusIcon = "bgbox_g";
    //             }else if(vm_status == "Suspended"){
    //                 vmStatusIcon = "bgbox_g";
    //             }else if(vm_status == "Terminated"){
    //                 vmStatusIcon = "bgbox_r";
    //             }else{
    //                 vmStatusIcon = "bgbox_g";
    //             }
    //             vm_badge += '<li class="sel_cr ' + vmStatusIcon + '" ><a href="javascript:void(0);" onclick="selectVm(\''+mcisId+'\',\''+vms[o].id+'\')" ><span class="txt">'+vms[o].name+'</span></a></li>';
    //             console.log("vm_status : ", vm_status)
    //
    //         }
    //         var sta = mcis.status
    //         var sl = sta.split("-");
    //         var status = sl[0].toLowerCase()
    //         var mcis_badge = '';
    //         var mcisStatusIcon = "icon_running_db.png";
    //         if(status == "running"){
    //             mcisStatusIcon = 'icon_running_db.png'
    //         }else if(status == "include" ){
    //             mcisStatusIcon = 'icon_stop_db.png'
    //         }else if(status == "suspended"){
    //             mcisStatusIcon = 'icon_stop_db.png'
    //         }else if(status == "terminate"){
    //             mcisStatusIcon = 'icon_terminate_db.png'
    //         }else{
    //             mcisStatusIcon = 'icon_stop_db.png'
    //         }
    //         mcis_badge = '<img src="/assets/img/contents/' + mcisStatusIcon + '" alt="' + status + '"/> '
    //
    //         $("#mcis_info_txt").text("[ "+mcis.name+"("+mcis.id+")"+" ]");
    //         $("#monitoring_mcis_status_img").empty()
    //         $("#monitoring_mcis_status_img").append(mcis_badge)
    //         $("#vmArrList").empty();
    //         $("#vmArrList").append(vm_badge);
    //
    //         // vm list options
    //         $("#vmList").empty()
    //         $("#vmList").append(vm_options)
    //
    //         $(".ds_cont_mbox .mtbox .g_list .listbox li.sel_cr").each(function(){
    //             var $sel_list = $(this),
    //             $detail_view = $(".monitoring_view");
    //             $sel_list.off("click").click(function(){
    //                 $sel_list.addClass("active");
    //                 $sel_list.siblings().removeClass("active");
    //                 $detail_view.addClass("active");
    //                 $detail_view.siblings().removeClass("active");
    //
    //                 $sel_list.off("click").click(function(){
    //                     if( $(this).hasClass("active") ) {
    //                         $sel_list.removeClass("active");
    //                         $detail_view.removeClass("active");
    //                     } else {
    //                             $sel_list.addClass("active");
    //                             $sel_list.siblings().removeClass("active");
    //                             $detail_view.addClass("active");
    //                             $detail_view.siblings().removeClass("active");
    //                     }
    //                 });
    //             });
    //         });
    //
    //     }// end of vms if
    //
    // })
}

function getCommonMcisDataCallbackSuccess(caller, mcisInfo){

        console.log(mcisInfo)
    var mcisId = mcisInfo.id
    var vms = mcisInfo.vm
    console.log(vms)
    var vm_badge ="";
    var vm_options = "";

    if(vms){
        // var init_vm = vms[0].id
        vm_len = vms.length
        for(var o in vms){
            var vm_status = vms[o].status
            vm_options +='<option value="'+vms[o].id+'">'+vms[o].name+'|'+vms[o].status+'|'+vms[o].description

            var vmStatusIcon = "bgbox_b";
            if(vm_status == "Running"){
                vmStatusIcon = "bgbox_b";
            }else if(vm_status == "include" ){
                vmStatusIcon = "bgbox_g";
            }else if(vm_status == "Suspended"){
                vmStatusIcon = "bgbox_g";
            }else if(vm_status == "Terminated"){
                vmStatusIcon = "bgbox_r";
            }else{
                vmStatusIcon = "bgbox_g";
            }
            vm_badge += '<li class="sel_cr ' + vmStatusIcon + '" ><a href="javascript:void(0);" onclick="selectVm(\''+mcisId+'\',\''+vms[o].id+'\')" ><span class="txt">'+vms[o].name+'</span></a></li>';
            console.log("vm_status : ", vm_status)

        }
        var sta = mcisInfo.status
        var sl = sta.split("-");
        var status = sl[0].toLowerCase()
        var mcis_badge = '';
        var mcisStatusIcon = "icon_running_db.png";
        if(status == "running"){
            mcisStatusIcon = 'icon_running_db.png'
        }else if(status == "include" ){
            mcisStatusIcon = 'icon_stop_db.png'
        }else if(status == "suspended"){
            mcisStatusIcon = 'icon_stop_db.png'
        }else if(status == "terminate"){
            mcisStatusIcon = 'icon_terminate_db.png'
        }else{
            mcisStatusIcon = 'icon_stop_db.png'
        }
        mcis_badge = '<img src="/assets/img/contents/' + mcisStatusIcon + '" alt="' + status + '"/> '

        $("#mcis_info_txt").text("[ "+mcisInfo.name+"("+mcisInfo.id+")"+" ]");
        $("#monitoring_mcis_status_img").empty()
        $("#monitoring_mcis_status_img").append(mcis_badge)
        $("#vmArrList").empty();
        $("#vmArrList").append(vm_badge);

        // vm list options
        $("#vmList").empty()
        $("#vmList").append(vm_options)

        $(".ds_cont_mbox .mtbox .g_list .listbox li.sel_cr").each(function(){
            var $sel_list = $(this),
                $detail_view = $(".monitoring_view");
            $sel_list.off("click").click(function(){
                $sel_list.addClass("active");
                $sel_list.siblings().removeClass("active");
                $detail_view.addClass("active");
                $detail_view.siblings().removeClass("active");

                $sel_list.off("click").click(function(){
                    if( $(this).hasClass("active") ) {
                        $sel_list.removeClass("active");
                        $detail_view.removeClass("active");
                    } else {
                        $sel_list.addClass("active");
                        $sel_list.siblings().removeClass("active");
                        $detail_view.addClass("active");
                        $detail_view.siblings().removeClass("active");
                    }
                });
            });
        });

    }// end of vms if
}

// vm 선택시 해당 vm의 monitoring 조회
function selectVm(mcis_id,vm_id){
    $('#vm_id').val(vm_id);
    var input_duration = $("#input_duration").val();
    var duration_type = $("#duration_type").val();
    var duration = input_duration+duration_type
    var period_type = $("#vm_period").val();
    var metric = $("#select_metric").val();
    showMonitoring(mcis_id,vm_id,metric,period_type,duration);
}

function btn_view_click(){
    var sel_history = $("#sel_history").val();
    var vm_id = $("#vm_id").val();
    var mcis_id =$("#mcis_id").val();

    var input_duration = $("#input_duration").val();
    var duration_type = $("#duration_type").val();
    var duration = input_duration+duration_type
    var period_type = $("#vm_period").val();
    var metric = $("#select_metric").val();

    showMonitoring(mcis_id,vm_id,metric,period_type,duration);
}


function ModalDetail(){
    $(".dashboard .status_list tbody tr").each(function(){
    var $td_list = $(this),
            $status = $(".server_status"),
            $detail = $(".server_info");
    $td_list.off("click").click(function(){
            $td_list.addClass("on");
            $td_list.siblings().removeClass("on");
            $status.addClass("view");
            $status.siblings().removeClass("on");
            $(".dashboard.register_cont").removeClass("active");
        $td_list.off("click").click(function(){
                if( $(this).hasClass("on") ) {
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
