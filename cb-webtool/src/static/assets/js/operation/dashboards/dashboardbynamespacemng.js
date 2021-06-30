$(document).ready(function(){
    jQuery('.sc_box.scrollbar-inner').scrollbar();
})

// dashboard 의 MCIS 목록에서 mcis 선택 : 색상반전, 선택한 mcis id set -> status변경에 사용
function selectMcis(id,name,target, obj){
    console.log("selectMcis")
    var mcis_id = id
    var mcis_name = name
    var init_select_areabox = $("#init_select_areabox").val()
    $target = $("#"+target)

    
    // if($target.hasClass("active")){
    //     location.href = "/Manage/MCIS/list/"+mcis_id+"/"+mcis_name
    //     return;
    // }

    // $("[id^='mcis_areabox_']").each(function(){
    //     var s_id = $(this).attr("id");
    //     console.log(s_id + ":" + target)
    //     if(s_id == target){
    //         try{
    //             var s_id = $(this).attr("id");
    //             $(this).addClass("active"); 
    //             console.log(s_id + " addClass active")
    //         }catch(e){
    //             console.log(e)
    //         }

    //     }else{
    //         $(this).removeClass("active");
    //         // console.log(s_id + "removeClass active")
    //     }
    // })
    // console.log(" active / deactive ")
    $("#mcis_id").val(mcis_id)
    $("#mcis_name").val(mcis_name)    
    console.log(" mcis_id =" + mcis_id + ", mcis_name = " + mcis_name);
 }

// callMcisLifeCycle -> McisLifeCycle -> callbackMcisLifeCycle
// confirm창을 띄울 때 mcismng와 동일한 key로 호출하므로 callback함수 이름도 같아야 한다.(util.js 참조)
function callMcisLifeCycle(type){
    var selectedCount = 0;
    // 선택된 mcis 가 있는지 체크.
    $("[id^='mcis_areabox_']").each(function(){        
        if($(this).hasClass("active")){            
            selectedCount++
            mcisLifeCycle($("#mcis_id").val(), type);//mcislifecycle.js 호출
        }
    })

    if( selectedCount == 0){
        commonAlert("Please Select MCIS!!")
    }

    /////// TODO : util.mcislifecycle.js 를 호출하도록 변경
    
}

// McisLifeCycle을 호출 한 뒤 return값 처리
function callbackMcisLifeCycle(resultStatus, resultData, type){
    var message = "MCIS "+type+ " complete!."
    if(resultStatus == 200 || resultStatus == 201){            
        commonAlert(message);
        location.reload();//완료 후 페이지를 reload -> 해당 mcis만 reload
        // 해당 mcis 조회
        // 상태 count 재설정
    }
}

function setMap(){
    //show_mcis2(url,JZMap);
    //function show_mcis2(url, map){
    // var JZMap = map;
    var JZMap = map_init()// TODO : map click할 때 feature 에 id가 없어 tooltip 에러나고 있음. 해결필요 

    //지도 그리기 관련
    var polyArr = new Array();

    // $("[id^='vmID_']").each(function(){
    $("input[name=vmID]").each(function(vmIndex, item){
        // var vmID = $(this).attr("id");
        // var vmIndex = vmID.split ("_")[1];
        var vmIDValue = $("#vmID_" + vmIndex).val();
        var vmNameValue = $("#vmName_" + vmIndex).val();
        var vmStatusValue = $("#vmStatus_" + vmIndex).val();
        var longitudeValue = $("#longitude_" + vmIndex).val();
        var latitudeValue = $("#latitude_" + vmIndex).val();

        var vms = new Object();
        vms.id = vmIDValue;
        vms.name = vmNameValue;
        vms.longitudeValue = longitudeValue;
        vms.latitudeValue = latitudeValue;
        // vms.status = vmStatusValue;
        // vms.status = vmStatusValue;

        var fromLonLat = longitudeValue+" "+latitudeValue;
        console.log(longitudeValue + " : " + latitudeValue);
        if(longitudeValue && latitudeValue){
            // polyArr.push(fromLonLat)
            drawMap(JZMap,longitudeValue,latitudeValue,vms)

            var polygon = "POLYGON(("+fromLonLat+"))";
            // drawPoligon(JZMap,fromLonLat);
            drawPoligon(JZMap,polygon);
        }
    // for(var i in mcis){
    //     for(var o in vms){
    //         vm_cnt++;
    //         var vm_status = vms[o].status
    //         var lat = vms[0].location.latitude
    //         var long = vms[0].location.longitude
    //         var provider = vms[0].location.cloudType

    //         // console.log("info : ",info)
    //         // point_feature.set('title',info.name)
    //         // point_feature.set('vm_status',info.status)
    //         // point_feature.set('vm_id',info.id)
    //         // point_feature.set('id',info.id)

    //         var fromLonLat = long+" "+lat;
    //         if(long && lat){
    //             polyArr.push(fromLonLat)
    //             drawMap(JZMap,long,lat,vms[o])
    //         }
            
    //         var polygon = "";
    //          console.log("poly arr : ",polyArr);
    //          if(polyArr.length > 1){
    //            polygon = polyArr.join(", ")
    //            polygon = "POLYGON(("+polygon+"))";
    //          }else{
    //            polygon = "POLYGON(("+polyArr[0]+"))";
    //          }
    //          if(polyArr.length >1){
    //             drawPoligon(JZMap,polygon);
    //           }
    //     }
    })
}