
function life_cycle2(type){
    var mcis_id = $("#mcis_id").val();
    var mcis_name = $("#mcis_name").val();
    if(!mcis_id){
        alert("Please Select MCIS!!")
        return;
    }
    var nameSpace = NAMESPACE;
    console.log("Start LifeCycle method!!!")
  
    var url = CommonURL+"/ns/"+nameSpace+"/mcis/"+mcis_id+"?action="+type
    var message = mcis_name+" "+type+ " complete!."
  

    var apiInfo = ApiInfo
    axios.get(url,{
        headers:{
            'Authorization': apiInfo
        }
    }).then(result=>{
        var status = result.status
        
        console.log("life cycle result : ",result)
        var data = result.data
        console.log("result Message : ",data.message)
        if(status == 200 || status == 201){
            
            alert(message);
            location.reload();
            //show_mcis(mcis_url,"");
        }
    })
}
function life_cycle(tag,type,mcis_id,mcis_name,vm_id,vm_name,mcis_url){
    var url = ""
    var nameSpace = NAMESPACE;
    var message = ""
    console.log("Start LifeCycle method!!!")
    
    if(tag == "mcis"){
        url = CommonURL+"/ns/"+nameSpace+"/mcis/"+mcis_id+"?action="+type
        message = mcis_name+" "+type+ " complete!."
    }else{
        url = CommonURL+"/ns/"+nameSpace+"/mcis/"+mcis_id+"/vm/"+vm_id+"?action="+type
        message = vm_name+" "+type+ " complete!."
    }

    var apiInfo = ApiInfo
    axios.get(url,{
        headers:{
            'Authorization': apiInfo
        }
    }).then(result=>{
        var status = result.status
        
        console.log("life cycle result : ",result)
        var data = result.data
        console.log("result Message : ",data.message)
        if(status == 200 || status == 201){
            
            alert(message);
            location.reload();
            //show_mcis(mcis_url,"");
        }
    })
}

function short_desc(str){
    var len = str.length;
    var result = "";
    if(len > 15){
        result = str.substr(0,15)+"...";
    }else{
        result = str;
    }

    return result;
 }

 //dashboard 에 동작하는 서버 댓수 및 다양한 정보를 뿌린다.
 function show_mcis(url, map){
   console.log("Show mcis Url : ",url)
   $("#vm_detail").hide();
   checkNS();

   var apiInfo = ApiInfo;
   console.log("apiInfo : ",apiInfo);
    axios.get(url,{
        headers:{
            'Authorization': apiInfo
        }
    }).then(result=>{
      
       console.log("Dashboard Data :",result.status);
       var data = result.data;
       console.log("func show_mcis result data : ",data)
       if(!data.mcis){
          location.href = "/MCIS/reg";
          return;
       }
       if(data.mcis.length == 0 ){
        location.href = "/MCIS/reg";
        return;
     }
       console.log("show mcis's map data : ",map);
        console.log("showmcis Data : ",data)
        var html = "";
        var mcis = data.mcis;
        var len = 0
        var mcis_cnt = 0 
        if(mcis){
            len = mcis.length;
        }
        mcis_cnt = len;
        var count = 0;
        
        var server_cnt = 0;
        var run_cnt = 0;
        var stop_cnt = 0;
        for(var i in mcis){
            var vm_len = 0
            var sta = mcis[i].status;
            var sl = sta.split("-");
            var badge = "";
            var status = sl[0].toLowerCase()
            var vms = mcis[i].vm
           
            if(vms){
               vm_len = vms.length
               server_cnt = server_cnt+vm_len;
            }
            for(var o in vms){
                if(vms[o].status == "Suspended"){
                    stop_cnt++;
                }
                if(vms[o].status == "Running"){
                    run_cnt++;
                }
            }

           console.log("mcis Status 1: ", mcis[i].status)
           console.log("mcis Status 2: ", status)
            if(status == "running"){
               badge += '<span class="badge badge-pill badge-success">'+sta+'</span>'
            }else if(status == "include" ){
               badge += '<span class="badge badge-pill badge-warning">'+sta+'</span>'
            }else if(status == "suspended"){
               badge += '<span class="badge badge-pill badge-warning">'+sta+'</span>'
            }else if(status == "terminate"){
               badge += '<span class="badge badge-pill badge-dark">'+sta+'</span>'
            }else{
               badge += '<span class="badge badge-pill badge-warning">'+sta+'</span>'
            }
            count++;
            if(count == 1){

            }
            html += '<tr id="tr_id_'+count+'" class="clickable-row">'
             +'<td class="text-center">'
             +'<div class="form-input">'
             +'<span class="input">'
             +'<input type="checkbox" class="chk" id="chk_'+count+'" value="'+mcis[i].id+'" item="'+mcis[i].name+'"><i></i></span></div>'
             +'</td>'
             +'<td>'
             +badge
             +'</td>'
             +'<td><a href="#!" onclick="show_vmList(\''+mcis[i].id+'\',\''+map+'\')" >'+mcis[i].name+'</a></td>'
             +'<td>'+vm_len+'</td>'
            //  +'<td>'+vm_len+'</td>'
            
            //  +'<td>0</td>'
             +'<td>'+short_desc(mcis[i].description)+'</td>'
             +'<td>'
             +'<button type="button" class="btn btn-icon dropdown-toggle" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">'
             +'<i class="fas fa-edit"></i>'
             +'<div class="dropdown-menu dropdown-menu-right" aria-labelledby="btnGroupDrop1">'
             +'<h6 class="dropdown-header text-center" style="background-color:#F2F4F4;;cursor:default;"><i class="fas fa-recycle"></i> LifeCycle</h6>'
                 +'<a class="dropdown-item text-right" href="#!" onclick="life_cycle(\'mcis\',\'resume\',\''+mcis[i].id+'\',\''+mcis[i].name+'\',\''+url+'\')">Resume</a>'
                 +'<a class="dropdown-item text-right" href="#!" onclick="life_cycle(\'mcis\',\'suspend\',\''+mcis[i].id+'\',\''+mcis[i].name+'\',\''+url+'\')">Suspend</a>'
                 +'<a class="dropdown-item text-right" href="#!" onclick="life_cycle(\'mcis\',\'reboot\',\''+mcis[i].id+'\',\''+mcis[i].name+'\',\''+url+'\')">Reboot</a>'
                 +'<a class="dropdown-item text-right" href="#!" onclick="life_cycle(\'mcis\',\'terminate\',\''+mcis[i].id+'\',\''+mcis[i].name+'\',\''+url+'\')">Terminate</a>'
             +'</div>'
             +'</button>'
            +'</td>'
            +'</tr>';
       }
       console.log("server_cnt:",server_cnt)
       console.log("mcis_cnt:",mcis_cnt)
       var new_str = mcis_cnt+'<small class="text-muted ml-2 mb-0"> / '+server_cnt+'</small>';
       //새로운 퍼블리싱에 들어가는 텍스트
       var new_str = '<strong>'+mcis_cnt+'</strong><span>/</span> '+server_cnt;
       $("#dash_1").append(new_str);
       $("#run_cnt").text(run_cnt);
       $("#stop_cnt").text(stop_cnt);

       $("#table_1").empty();
       $("#table_1").append(html);
    //    var infra_str = "Infra - Server (MCIS : "+mcis[0].name+")"
    //    $("#infra_mcis").text(infra_str)
      // show_card(mcis[0].id,mcis[0].name);
      if(vm_len > 0){
       show_vmList(mcis[0].id,map);
      }else{
       show_vmList("",map);
      }
    
      
       
       //fnMove("table_1");
    //    $("#mcis_id").val(mcis[0].id)
    //    $("#mcis_name").val(mcis[0].name)
    //    }).catch(function(error){
    //     console.log("show mcis error at dashboard js: ",error);
    //    });
    }).catch((error) => {
        console.warn(error);
        console.log(error.response) 
    });
}

//새로운 퍼블리싱에 적용할 function
function show_mcis2(url, map){
    console.log("Show mcis Url : ",url)
    $("#vm_detail").hide();
    checkNS();
 
    var apiInfo = ApiInfo;
    var JZMap = map;
    console.log("apiInfo : ",apiInfo);
    axios.get(url,{
        headers:{
            'Authorization': apiInfo
        }
    }).then(result=>{
       
        console.log("Dashboard Data :",result.status);
        var data = result.data;
        console.log("func show_mcis result data : ",data)
        if(!data.mcis){
           location.href = "/Manage/MCIS/reg";
           return;
        }
        if(data.mcis.length == 0 ){
            location.href = "/Manage/MCIS/reg";
            return;
        }
        console.log("show mcis's map data : ",map);
        console.log("showmcis Data : ",data)
        var html = "";
        var mcis = data.mcis;
        var len = 0
        var mcis_cnt = 0 
        if(mcis){
            len = mcis.length;
        }
        mcis_cnt = len;
        var count = 0;
        
        var server_cnt = 0;
         
        var html = "";
        var run_cnt = 0;
        var stop_cnt = 0;
        for(var i in mcis){
            count++;
            var vm_run_cnt = 0;
            var vm_stop_cnt = 0;
            var terminate_cnt = 0;
            var vm_len = 0
            var sta = mcis[i].status;
            var sl = sta.split("-");
            var mcis_badge = "";
            var vm_badge = "";
            var status = sl[0].toLowerCase()
            var vms = mcis[i].vm
            console.log("mcis status : ",status)
            var vm_status = "";
            if(vms){
                vm_len = vms.length
                server_cnt = server_cnt+vm_len;
            }
            //VM  상태 및 기타 생성하기
            var vm_cnt = 0
            var vm_html = "";
            //지도 그리기 관련
            var polyArr = new Array();
            for(var o in vms){
                vm_cnt++;
                var vm_status = vms[o].status
                var lat = vms[0].location.latitude
                var long = vms[0].location.longitude
                var provider = vms[0].location.cloudType

                var fromLonLat = long+" "+lat;
                if(long && lat){
                    polyArr.push(fromLonLat)
                    drawMap(JZMap,long,lat,vms[o])
                }

                console.log(lat, long, provider)
                 
                if(vms[o].status == "Suspended"){
                    stop_cnt++;
                    vm_stop_cnt++;
                }else if(vms[o].status == "Running"){
                    run_cnt++;
                    vm_run_cnt++;
                }else if(vms[o].status == "Terminated"){
                    terminate_cnt++;
                }else{
                    stop_cnt++;
                    vm_stop_cnt++;
                }

                if(vm_status == "Running"){
                    vm_badge += "shot bgbox_b"
                }else if(vm_status == "include" ){
                    vm_badge += "shot bgbox_y"
                }else if(vm_status == "Suspended"){
                    vm_badge += "shot bgbox_y"
                }else if(vm_status == "Terminated"){
                    vm_badge += "shot bgbox_r"
                }else{
                    vm_badge += "shot bgbox_g"
                }

                vm_html +='<div class="'+vm_badge+'"><a href="javascript:void(0);"><span>'+vm_cnt+'</span></a></div>'
            }
            var polygon = "";
            console.log("poly arr : ",polyArr);
            if(polyArr.length > 1){
                polygon = polyArr.join(", ")
                polygon = "POLYGON(("+polygon+"))";
            }else{
                polygon = "POLYGON(("+polyArr[0]+"))";
            }
            if(polyArr.length >1){
                drawPoligon(JZMap,polygon);
            }

            //MCIS name  / MCIS 상태
            if(status == "running"){
                mcis_badge += 'state color_b'
            }else if(status == "include" ){
                mcis_badge += 'state color_y'
            }else if(status == "suspended"){
                mcis_badge += 'state color_y'
            }else if(status == "terminate"){
                mcis_badge += 'state color_r'
            }else{
                mcis_badge += 'state color_g'
            }
            var cursor = ""
            //  if(i == mcis_cnt-1){
            //      cursor = "active"
            //  }

            html +='<div class="areabox dbinfo cursor '+cursor+'" id="mcis_areabox_'+i+'" onclick="change_mcis(\''+mcis[i].id+'\',\''+mcis[i].name+'\',\'mcis_areabox_'+i+'\')">'
                  +'<div class="box">';
            html += '<div class="top">'
                  +'<div class="txtbox">'
                  +'<div class="tit">'+mcis[i].name+'</div>'
                  +'<div class="txt"><span class="bgbox_b"></span>Available 01</div>'
                  +'</div>'
                  +'<div class="'+mcis_badge+'"></div>'
                  +'</div>';
            // 전체 인프라 갯수 및 각각의 상태에 따른 VM 갯수
            html +='<div class="numbox">infra <strong class="color_b">'+vm_cnt+'</strong>' 
                   +'<span class="line">(</span> <span class="num color_b">'+vm_run_cnt+'</span>' 
                   +'<span class="line">/</span> <span class="num color_y">'+vm_stop_cnt+'</span>' 
                   +'<span class="line">/</span> <span class="num color_r">'+terminate_cnt+'</span>'
                   +'<span class="line">)</span></div>';

            // 서버 갯수
            html += '<div class="numinfo">'
            html += '<div class="num">server'+vm_cnt+'</div>'
            html += '</div>'
            // 각각의 VM 항목들
            html +='<div class="shotbox">'
            html += vm_html;
            html +='</div></div></div>'

 
            console.log("mcis Status 1: ", mcis[i].status)
            console.log("mcis Status 2: ", status)
             
             
            if(count == 1){
 
            }
            
        }
        html +='<div class="areabox">'
             +'<div class="box">'
             +'<a href="/Manage/MCIS/reg" class="btn_add"><span>+</span></a>'
             +'</div></div>';
        console.log("server_cnt:",server_cnt)
        console.log("mcis_cnt:",mcis_cnt)
        var new_str = mcis_cnt+'<small class="text-muted ml-2 mb-0"> / '+server_cnt+'</small>';
        //새로운 퍼블리싱에 들어가는 텍스트
        var new_str = '<strong>'+mcis_cnt+'</strong><span>/</span> '+server_cnt;
        $("#dash_1").append(new_str);
        $("#run_cnt").text(run_cnt);
        $("#stop_cnt").text(stop_cnt);
 
        $("#table_1").empty();
        $("#table_1").append(html);
   
        //event 속성
        // $(".dashboard.dashboard_cont .ds_cont .dbinfo").each(function(){
        //     var $list = $(this);
        //     $list.on('click', function(){
        //           if( $(this).hasClass("active") ) {
        //               $list.removeClass("active");
        //       } else {
        //               $list.addClass("active");
        //               $list.siblings().removeClass("active");
        //       }
        //       });
        //   });

    // }).catch(function(error){
    //  console.log("show mcis error at dashboard js: ",error);
    // });
    }).catch((error) => {
        console.warn(error);
        console.log(error.response) 
    });
 }
 
 function change_mcis(id,name,target){
    var mcis_id = id
    var mcis_name = name
    var init_select_areabox = $("#init_select_areabox").val()
    $target = $("#"+target)
    if($target.hasClass("active")){
        location.href = "/Manage/MCIS/list/"+mcis_id+"/"+mcis_name
        return;
    }
    $("[id^='mcis_areabox_']").each(function(){
        var s_id = $(this).attr("id");
        if(s_id == target){
            $(this).addClass("active");
         
        }else{
            $(this).removeClass("active");
           
        }
    })
    $("#mcis_id").val(mcis_id)
    $("#mcis_name").val(mcis_name)
    
 }

 function change_mcis_db(id,name,target){
    var mcis_id = id
    var mcis_name = name
    
    //$target.addClass("active")
    $("[id^='mcis_areabox_']").each(function(){
        var s_id = $(this).attr("id");
        if(s_id == target){
            $(this).addClass("active");
         
        }else{
            $(this).removeClass("active");
           
        }
    })
    $("#mcis_id").val(mcis_id)
    $("#mcis_name").val(mcis_name)
    
 }
function show_vmList(mcis_id,map){
    $("#vm_detail").hide();
    $("#chart_detail").hide();
    $("#map_detail").hide();

   
   var url = CommonURL+"/ns/"+NAMESPACE+"/mcis/"+mcis_id;
   var mcis_url = CommonURL+"/ns/"+NAMESPACE+"/mcis";
   var apiInfo = ApiInfo;
   console.log("vmList",url)
   if(mcis_id){
       //여기가 geo location 정보 가져 오는 곳
    
    console.log("vm list map info : ",map)
   
    if(map){
        $("#map").empty();
        map = map_init();
    }
    //map = map_init();
    getGeoLocationInfo(mcis_id,map);
    //여기는 차트 불러 오는 부분
       $.ajax({
           type:'GET',
           url:url,
           beforeSend : function(xhr){
            xhr.setRequestHeader("Authorization", apiInfo);
            xhr.setRequestHeader("Content-type","application/json");
        },
           
       // async:false,
          
       }).done(function(data){
        var vm = data.vm
        var mcis_name = data.name 
        $("#mcis_id").val(mcis_id)
        $("#mcis_name").val(mcis_name)
        var html = "";
        console.log("VM DATA : ",vm)
        for(var i in vm){
            var sta = vm[i].status;
            //ip 정보 가져 오기
            console.log("========get ip region info=======")
            
        
            var status = sta.toLowerCase()
            console.log("VM Status : ",status)
            var configName = vm[i].connectionName
            console.log("outer vm configName : ",configName)
            var count = 0;
            $.ajax({
                url: SpiderURL+"/connectionconfig",
                async:false,
                type:'GET',
                beforeSend : function(xhr){
                    xhr.setRequestHeader("Authorization", apiInfo);
                    xhr.setRequestHeader("Content-type","application/json");
                },
                
               

            }).done(function(data2){
                res = data2.connectionconfig
                var badge = "";
                for(var k in res){
                    // console.log(" i value is : ",i)
                    // console.log("outer config name : ",configName)
                    // console.log("Inner ConfigName : ",res[k].ConfigName)
                    if(res[k].ConfigName == vm[i].connectionName){
                        var provider = res[k].ProviderName
                        var kv_list = vm[i].cspViewVmDetail.KeyValueList
                        var archi = ""
                        for(var p in kv_list){
                            if(kv_list[p].Key == "Architecture"){
                             archi = kv_list[p].Value 
                            }
                        }
                        
                        if(status == "running"){
                            badge += '<span class="badge badge-pill badge-success">RUNNING</span>'
                        }else if(status == "suspended"){
                            badge += '<span class="badge badge-pill badge-warning">SUSPEND</span>'
                        }else if(status == "terminate"){
                            badge += '<span class="badge badge-pill badge-dark">TERMINATED</span>'
                        }else{
                            badge += '<span class="badge badge-pill badge-dark">'+status+'</span>'
                        }
                        count++;
                        if(count == 1){
            
                        }
                        html += '<tr id="tr_id_'+count+'"  class="clickable-row">'
                        +'<td class="text-center">'
                        +'<div class="form-input">'
                        +'<span class="input">'
                        +'<input type="checkbox" item="'+mcis_name+'"    mcisid="'+mcis_id+'" class="chk2" id="chk2_'+count+'" value="'+vm[i].id+'|'+mcis_id+'"><i></i></span></div>'
                        +'</td>'
                        +'<td>'
                        +badge
                        +'</td>'
                        +'<td><a href="#!" onclick="show_vm(\''+mcis_id+'\',\''+vm[i].id+'\',\''+vm[i].name+'\',\''+vm[i].imageId+'\');">'+vm[i].name+'</a></td>'
                        
                        +'<td>'+provider+'</td>'
                        +'<td>'+vm[i].region.Region+'</td>'
         
                       
                        +'<td>'+archi+'</td>'
                        +'<td>'+vm[i].publicIP+'</td>'
                        +'<td>'+short_desc(vm[i].description)+'</td>'
                        
                        +'<td>'
                        +'<button type="button" class="btn btn-icon dropdown-toggle" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">'
                        +'<i class="fas fa-edit"></i>'
                        +'<div class="dropdown-menu dropdown-menu-right" aria-labelledby="btnGroupDrop1">'
                        +'<h6 class="dropdown-header text-center" style="background-color:#F2F4F4;;cursor:default;"><i class="fas fa-recycle"></i> LifeCycle</h6>'
                            +'<a class="dropdown-item text-right" href="#!" onclick="life_cycle(\'vm\',\'resume\',\''+mcis_id+'\',\''+mcis_name+'\',\''+vm[i].id+'\',\''+vm[i].name+'\',\''+mcis_url+'\')">Resume</a>'
                            +'<a class="dropdown-item text-right" href="#!" onclick="life_cycle(\'vm\',\'suspend\',\''+mcis_id+'\',\''+mcis_name+'\',\''+vm[i].id+'\',\''+vm[i].name+'\',\''+mcis_url+'\')">Suspend</a>'
                            +'<a class="dropdown-item text-right" href="#!" onclick="life_cycle(\'vm\',\'reboot\',\''+mcis_id+'\',\''+mcis_name+'\',\''+vm[i].id+'\',\''+vm[i].name+'\',\''+mcis_url+'\')">Reboot</a>'
                            +'<a class="dropdown-item text-right" href="#!" onclick="life_cycle(\'vm\',\'terminate\',\''+mcis_id+'\',\''+mcis_name+'\',\''+vm[i].id+'\',\''+vm[i].name+'\',\''+mcis_url+'\')">Terminate</a>'
                        +'</div>'
                        +'</button>'
                        // +'<button type="button" class="btn btn-icon"  aria-haspopup="true" aria-expanded="false" onclick="agentSetup(\''+mcis_id+'\',\''+vm[i].id+'\',\''+vm[i].publicIP+'\')">'
                        // +'<i class="fas fa-desktop"></i>'
                       // +'<div class="dropdown-menu dropdown-menu-right" aria-labelledby="btnGroupDrop2">'
                           // +'<a class="dropdown-item" href="#!" onclick="life_cycle(\'vm\',\'resume\',\''+mcis_id+'\',\''+mcis_name+'\',\''+vm[i].id+'\',\''+vm[i].name+'\')">Resume</a>'
                           // +'<a class="dropdown-item" href="#!" onclick="life_cycle(\'vm\',\'suspend\',\''+mcis_id+'\',\''+mcis_name+'\',\''+vm[i].id+'\',\''+vm[i].name+'\')">Suspend</a>'
                           // +'<a class="dropdown-item" href="#!" onclick="life_cycle(\'vm\',\'reboot\',\''+mcis_id+'\',\''+mcis_name+'\',\''+vm[i].id+'\',\''+vm[i].name+'\')">Reboot</a>'
                           // +'<a class="dropdown-item" href="#!" onclick="life_cycle(\'vm\',\'terminate\',\''+mcis_id+'\',\''+mcis_name+'\',\''+vm[i].id+'\',\''+vm[i].name+'\')">Terminate</a>'
                       // +'</div>'
                       // +'</button>'
                        +'</td>'
                        +'</tr>';
                    }
                    
                    
                    }
                    $("#table_2").empty();
                    $("#table_2").append(html);
                    fnMove("table_2");

            })
            
            }
    })

   }else{
       $("#table_2").empty();
       $("#table_2").append("<td colspan='8'>Does not Exist</td>");
   }
           
   
}

function agentSetup(mcis_id,vm_id,public_ip){
   
        alert("Monitoring service on this server is turned off.");
        if(confirm("Would to enable the monitoring service?")){
            var reg_url = "/monitoring/install/agent"
            var query_param = "/"+mcis_id+"/"+vm_id+"/"+public_ip;
            console.log("agent setup query param: ",query_param);
            location.href = reg_url+query_param;
        }else{
            return;
        }
       

    

    
}
 

 function show_vm(mcis_id,vm_id,vm_name,image_id){
    checkDragonFly(mcis_id,vm_id);
    show_vmSSHInfo(mcis_id, vm_id);
    show_vmDetailList(mcis_id, vm_id);
    show_vmSpecInfo(mcis_id, vm_id);
    show_vmNetworkInfo(mcis_id, vm_id);
    show_vmSecurityGroupInfo(mcis_id, vm_id);
    
    show_images(image_id);
    $("#current_vmid").val(vm_id);
    $("#server_text").empty();
    $("#server_text").append("<strong>"+vm_name+" Server Info</strong>");
    $("#vm_detail").show();
 }

 function sel_table(targetNo,mcid){
     var $target = $("#card_"+targetNo+"");
     var html = "";
     url = CommonURL+"/ns/"+NAMESPACE+"/mcis/"+mcid
     var apiInfo = ApiInfo
    axios.get(url,{
        headers:{
            'Authorization': apiInfo
        }
    }).then(result=>{
         var data = result.data.vm
         for(var i in data){

         }
     })
     html += '<tr>'
             +'<td class="text-center">'
             +'<div class="form-input">'
             +'<span class="input">'
             +'<input type="checkbox" id=""><i></i>'
             +'</span></div>'
             +'</td>'
             +'<td>1</td>'
             +'<td><a href="">Baristar1</a></td>'
             +'<td>aws driver 1aws driver ver0.1</td>'
             +'<td>aws key 1</td>'
             +'<td>ap-northest-1</td>'
             +'<td>'
             +'<div class="custom-control custom-switch">'
             +'<input type="checkbox" class="custom-control-input" id="customSwitch1">'
             +'<label class="custom-control-label" for="customSwitch1"></label></div>'
             +'</td>'
             +'<td>'
             +'<span class="badge badge-pill badge-warning">stop</span>'
             +'</td>'
             +'<td>2019-05-05</td>'
             +'</tr>';
             
    $target.empty();         
    $target.append(html);

 }

 function deleteHandler(cl,target,){
    var url = SpiderURL+"/connectionconfig"
 }

 function mcis_delete(){
    
    var cnt = 0;
    var mcis_id = "";
    var apiInfo = ApiInfo;
    $(".chk").each(function(){
        if($(this).is(":checked")){
            //alert("chk");
            cnt++;
            mcis_id = $(this).val();        
        }
        if(cnt < 1 ){
            alert("Select Delete");
            return;
        }

        if(cnt == 1){
           console.log("mcis_id ; ",mcis_id)
            var url = CommonURL+"/ns/"+NAMESPACE+"/mcis/"+mcis_id
            
            if(confirm("Delete?")){
             axios.delete(url,{
                headers :{
                    'Content-type': 'application/json',
                    'Authorization': apiInfo,
                    }
             }).then(result=>{
                 var data = result.data
                 if(result.status == 200){
                     alert(data.message)
                     location.reload(true)
                 }
             })
            }
        }

        if(cnt >1){
            alert("It Only one Delete")
            return;
        }

    })
 }
function getConnection(){
    var apiInfo = ApiInfo;
    $.ajax({
        url: SpiderURL+"/connectionconfig",
        async:false,
        type:'GET',
        beforeSend : function(xhr){
            xhr.setRequestHeader("Authorization", apiInfo);
            xhr.setRequestHeader("Content-type","application/json");
        },
       

    }).done( function(data2){
        res = data2.connectionconfig
        console.log("connection info : ",res);
        var provider = "";
        var aws_cnt = 0;
        var gcp_cnt = 0;
        var azure_cnt = 0;
        var open_cnt = 0;
        var cloudIt_cnt = 0;
        var ali_cnt = 0;
        var cp_cnt = 0;
        var connection_cnt = 0;
        var html = "";
        for(var k in res){
            provider = res[k].ProviderName 
            connection_cnt++;
            provider = provider.toLowerCase();
            console.log("provider lowercase : ",provider);
            
            if(provider == "aws"){
                aws_cnt++;  
             
            }
            if(provider == "azure"){
                azure_cnt++;
                 
            }
            if(provider == "alibaba"){
                ali_cnt++;
              
                    
            }
            if(provider == "gcp"){
                gcp_cnt++;
            
            }
            if(provider == "cloudit"){
                cloudIt_cnt++;
              
            }
            if(provider == "openstack"){
                open_cnt++;
              
            }
        }
        
        
        if(aws_cnt > 0 ){
           
            html +='<li class="bg_b">'
                 +'<a href="#!"><span>AWS('
                 +aws_cnt
                 +')</span></a></li>';          
        }
        if(azure_cnt > 0){
            html +='<li class="bg_y">'
                 +'<a href="#!"><span>AZ('
                 +azure_cnt
                 +')</span></a></li>';       
        }
        if(ali_cnt > 0){
           
            html +='<li class="bg_r">'
                 +'<a href="#!"><span>ALI('
                 +ali_cnt
                 +')</span></a></li>';       
                
        }
        if(gcp_cnt > 0){
          
            html +='<li class="bg_g">'
            +'<a href="#!"><span>GCP('
            +gcp_cnt
            +')</span></a></li>';     
        }
        if(cloudIt_cnt > 0){
          
            html +='<li class="bg_n">'
            +'<a href="#!"><span>CLIT('
            +cloudIt_cnt
            +')</span></a></li>';  
        }
        if(open_cnt > 0){
           
            html +='<li class="bg_b">'
            +'<a href="#!"><span>OPS('
            +open_cnt
            +')</span></a></li>';  
        }

        if(aws_cnt > 1){
            aws_cnt = 1
        }
        if(azure_cnt > 1){
            azure_cnt = 1
        }
        if(ali_cnt > 1){
            ali_cnt = 1
        }
        if(open_cnt > 1){
            open_cnt = 1
        }
        if(cloudIt_cnt > 1){
            cloudIt_cnt = 1
        }
        if(gcp_cnt > 1){
            gcp_cnt = 1
        }

        cp_cnt = aws_cnt+azure_cnt+ali_cnt+open_cnt+cloudIt_cnt+gcp_cnt;
        var str = '<strong>'+cp_cnt+'</strong><span>/</span>'+connection_cnt;
        $("#dash_2").empty();
        $("#dash_2").append(str);
        $("#dash_3").empty();
        $("#dash_3").append(html);
    })
    
}

/*var f = getOSType("IMAGE-aws-developer").then(data=>{
    console.log("axios inner data : ",data)
});
console.log("axios return value : ",f);
*/
 function mcis_reg(){
    
    var cnt = 0;
    var mcis_id = "";
    $(".chk").each(function(){
        if($(this).is(":checked")){
            //alert("chk");
            cnt++;
            mcis_id = $(this).val();
            mcis_name = $(this).attr("item");

        }
        if(cnt < 1 ){
            alert("Select Regist");
            return;
        }

        if(cnt == 1){
           console.log("mcis_id ; ",mcis_id)
            var url = "/MCIS/reg/"+mcis_id+"/"+mcis_name
            
            if(confirm("Register?")){
                location.href = url;
            }
        }

        if(cnt >1){
            alert("Only one Regist")
            return;
        }

    })
 }

 function vm_delete(){
    
    var cnt = 0;
    var vm_id = "";
    var mcis_id ="";
    var apiInfo = ApiInfo;
    $(".chk").each(function(){
        if($(this).is(":checked")){
            //alert("chk");
            cnt++;
            id = $(this).val(); 
            idArr = id.split ("|")  
            vm_id = idArr[0]
            mcis_id = idArr[1]    
        }
        if(cnt < 1 ){
            alert("Select Delete.");
            return;
        }

        if(cnt == 1){
           console.log("mcis_id ; ",vm_id)
            var url = CommonURL+"/ns/"+NAMESPACE+"/mcis/"+mcis_id+"/vm/"+vm_id
            
            if(confirm("Delete?")){
             axios.delete(url,{
                headers :{
                    'Content-type': 'application/json',
                    'Authorization': apiInfo,
                    }
             }).then(result=>{
                 var data = result.data
                 if(result.status == 200){
                     alert(data.message)
                     location.reload(true)
                 }
             })
            }
        }

        if(cnt >1){
            alert("Only one Delete")
            return;
        }

    })
 }

 function getProvider(connectionInfo){
     url = SpiderURL+"/connectionconfig"
     var apiInfo = ApiInfo
    axios.get(url,{
        headers:{
            'Authorization': apiInfo
        }
    }).then(result=>{
         var data = result.data.connectionconfig

         for(var i in data){
             if(connetionInfo == data[i].ConfigName){}
         }
     })
 }

 function show_vmDetailList(mcis_id, vm_id){
     url = CommonURL+"/ns/"+NAMESPACE+"/mcis/"+mcis_id+"/vm/"+vm_id
     var apiInfo = ApiInfo
    axios.get(url,{
        headers:{
            'Authorization': apiInfo
        }
    }).then(result=>{
         var data = result.data;
         var publicIP = data.publicIP;
         $("#current_publicIP").val(publicIP);
         var html = ""
         $.ajax({
            url: SpiderURL+"/connectionconfig",
            async:false,
            type:'GET',
            

        }).done(function(data2){
            res = data2.connectionconfig
            var provider = "";
            for(var k in res){
                if(res[k].ConfigName == data.connectionName){
                    provider = res[k].ProviderName
                    console.log("Inner Provider : ",provider)
                }
            }

            html += '<tr>'
                    +'<th scope="colgroup"rowspan="10" class="text-center">Infra - Server</th>'

                    +'<th scope="colgroup" class="text-right">Server ID</th>'
                    +'<td  colspan="1">'+data.id+'</td>'
                    
                    
                    +'<th scope="colgroup" class="text-right">Cloud Provider</th>'
                    +'<td colspan="1">'+provider+'</td>'
                    +'</tr>'


                    +'<tr>'
                    // +'<th scope="colgroup" class="text-right">CP VMID</th>'
                    // +'<td  colspan="1">'+data.id+'</td>'
                   
                    +'<th scope="colgroup" class="text-right">Region</th>'
                    +'<td  colspan="1" >'+data.region.Region+'</td>'
                    +'<th scope="colgroup" class="text-right">Zone</th>'
                    +'<td  colspan="1">'+data.region.Zone+'</td>'
                    +'</tr>'

                    
                    +'<tr>'
                    +'<th scope="colgroup" class="text-right">Public IP</th>'
                    +'<td  colspan="1">'+data.publicIP+'</td>'
                    
                    +'<th scope="colgroup" class="text-right">Public DNS</th>'
                    +'<td  colspan="1">'+data.publicDNS+'</td>'
                    +'</tr>'

                    +'<tr>'
                    +'<th scope="colgroup" class="text-right">Private IP</th>'
                    +'<td colspan="1">'+data.privateIP+'</td>'
                    
                    +'<th scope="colgroup" class="text-right">Private DNS</th>'
                    +'<td colspan="1">'+data.privateDNS+'</td>'
                    +'</tr>'

                    +'<tr>'
                    +'<th scope="colgroup" class="text-right">Server Status</th>'
                    +'<td colspan="3">'+data.status+'</td>'
                    +'</tr>';

              
            $("#vm").empty();
            $("#vm").append(html);
            fnMove("vm_detail");

        })
       
            
         
     })

 }
 function vm_reg(){
    
    var cnt = 0;
    var mcis_id = "";
    var mcis_name = "";
    
    mcis_id = $("#mcis_id").val()
    mcis_name = $("#mcis_name").val()
    var url = "/MCIS/reg/"+mcis_id+"/"+mcis_name
    console.log("vm reg url : ",url)
    if(confirm("Add Server?")){
        location.href = url;
    }

 }


function show_vmSpecInfo(mcis_id, vm_id){
    var url = CommonURL+"/ns/"+NAMESPACE+"/mcis/"+mcis_id+"/vm/"+vm_id
    var apiInfo = ApiInfo
    axios.get(url,{
        headers:{
            'Authorization': apiInfo
        }
    }).then(result=>{
        var data = result.data
        var html = ""
        var url2 = CommonURL+"/ns/"+NAMESPACE+"/resources/spec"
        var spec_id = data.specId
        $.ajax({
           url: url2,
           async:false,
           type:'GET',
           beforeSend : function(xhr){
            xhr.setRequestHeader("Authorization", apiInfo);
            xhr.setRequestHeader("Content-type","application/json");
        },
           

       }).done( function(result){
        var res = result.spec
       
        for(var k in res){
            if(res[k].id == spec_id){
             html += '<tr>'
                    +'<tr>'
                    +'<th scope="colgroup" rowspan="4" class="text-right"><i class="fas fa-server"></i>Server Spec</th>'
                    +'<th scope="colgroup" class="text-right">vCPU</th>'
                    +'<td colspan="1">'+res[k].num_vCPU+' vcpu</td>'
                  
                    +'<th scope="colgroup" class="text-right">Memory(Ghz)</th>'
                    +'<td  colspan="1">'+res[k].mem_GiB+' GiB</td>'
                    +'</tr>'
                    +'<tr>'
                    +'<th scope="colgroup" class="text-right">Disk(GB)</th>'
                    +'<td colspan="1">'+res[k].storage_GiB+' GiB</th>'
                    +'<th scope="colgroup" class="text-right">Cost($) / Hour </th>'
                    +'<td colspan="1">'+res[k].cost_per_hour+'</td>'
                    +'</tr>'
                    +'<tr>'
                    +'<th scope="colgroup" class="text-right">OS Type</th>'
                    +'<td  colspan="3">'+res[k].os_type+'</td>'
                    +'</tr>'
            }
        } 
        $("#vm_spec").empty();
        $("#vm_spec").append(html);

    })
      
           
        
    })

}

function show_vmNetworkInfo(mcis_id, vm_id){
    var url = CommonURL+"/ns/"+NAMESPACE+"/mcis/"+mcis_id+"/vm/"+vm_id
    var apiInfo = ApiInfo
    axios.get(url,{
        headers:{
            'Authorization': apiInfo
        }
    }).then(result=>{
        var data = result.data
        var html = ""
        var url2 = CommonURL+"/ns/"+NAMESPACE+"/resources/vNet"
        var spec_id = data.vNetId
        $.ajax({
           url: url2,
           async:false,
           type:'GET',
           beforeSend : function(xhr){
            xhr.setRequestHeader("Authorization", apiInfo);
            xhr.setRequestHeader("Content-type","application/json");
        },
           

       }).done(function(result){
        var res = result.vNet
       console.log("Network Info : ",result)
        for(var k in res){
            if(res[k].id == spec_id){
             var subnetInfoList = res[k].subnetInfoList
             var subnetArr = new Array()
             var str = ""
             if(subnetInfoList){
                 for(var o in subnetInfoList){
                      subnetArr.push(subnetInfoList[o].IPv4_CIDR)
                 }
                 str = subnetArr.join(",")
             }
             console.log("Subnet str : ",str)
             html += '<tr>'
                    +'<th scope="colgroup" rowspan="5" class="text-right"><i class="fas fa-network-wired"></i>Network</th>'
                    +'<th scope="colgroup" class="text-right">Network Name</th>'
                    +'<td  colspan="1">'+res[k].cspVNetName+'</td>'
                    +'<th scope="colgroup" class="text-right">Network ID</th>'
                    +'<td colspan="1">'+res[k].cspVNetId+'</td>'
                    
                    +'</tr>'
                    +'<tr>'
                    +'<th scope="colgroup" class="text-right">Cidr Block</th>'
                    +'<td colspan="3">'+res[k].cidrBlock+'</th>'
                    +'</tr>'
                    +'<tr>'
                    +'<th scope="colgroup" class="text-right">Subnet</th>'
                    +'<td colspan="3">'+str+'</th>'
                    +'</tr>'
                 //    +'<tr>'
                 //    +'<th scope="colgroup">Interface</th>'
                 //    +'<td colspan="3">'+res[k].cidrBlock+'</th>'
                 //    +'</tr>'
                   
            }
        } 
        console.log("vnetwork html : ",html)
        $("#vm_vnetwork").empty();
        $("#vm_vnetwork").append(html);

    })
      
           
        
    })

}

function show_images(image_id){
    var url = CommonURL+"/ns/"+NAMESPACE+"/resources/image/"+image_id
    var apiInfo = ApiInfo
    axios.get(url,{
        headers:{
            'Authorization': apiInfo
        }
    }).then(result=>{
        var data = result.data
        console.log("Image Data : ",data);
        var html = ""
            
        html += '<tr>'
                +'<th scope="colgroup" rowspan="5" class="text-right"><i class="fas fa-compact-disc"></i>Image</th>'
                +'<th scope="colgroup" class="text-right">Image Name</th>'
                +'<td  colspan="1">'+data.name+'</td>'
                +'<th scope="colgroup" class="text-right">Image ID</th>'
                +'<td colspan="1">'+data.id+'</td>'
                
                +'</tr>'
                +'<tr>'
                +'<th scope="colgroup" class="text-right">Guest OS</th>'
                +'<td colspan="1">'+data.guestOS+'</th>'
                
                +'<th scope="colgroup" class="text-right">Description</th>'
                +'<td colspan="1">'+data.description+'</th>'
                +'</tr>'
            
                          
             
             
               $("#vm_image").empty();
               $("#vm_image").append(html);

           })

    

}

function show_vmSecurityGroupInfo(mcis_id, vm_id){
    var url = CommonURL+"/ns/"+NAMESPACE+"/mcis/"+mcis_id+"/vm/"+vm_id
    var apiInfo = ApiInfo
    axios.get(url,{
        headers:{
            'Authorization': apiInfo
        }
    }).then(result=>{
        var data = result.data
        console.log("Security Group : ",data);
        var html = ""
        // var url2 = "/ns/"+NAMESPACE+"/resources/securityGroup"
        var spec_id = data.securityGroupIds
        var cnt = spec_id.length
        html += '<tr>'
             +'<th scope="colgroup" colspan="'+cnt+' "class="text-right"><i class="fas fa-shield-alt"></i>SecurityGroup</th>'
             +'<th scope="colgroup" colspan="'+cnt+'" class="text-right">SecurityGroup ID</th>'
        for(var i in spec_id){
            if( i == 0){
                html +='<td colspan="3">'+spec_id[i]+'</td></tr>'
            }else{
                html +='<tr><td colspan="3">'+spec_id[i]+'</td></tr>'
            }
        }
        

        $("#vm_sg").empty();
        $("#vm_sg").append(html);

                
        
    })

}



function show_vmSSHInfo(mcis_id, vm_id){
    var url = CommonURL+"/ns/"+NAMESPACE+"/mcis/"+mcis_id+"/vm/"+vm_id
    var apiInfo = ApiInfo
    axios.get(url,{
        headers:{
            'Authorization': apiInfo
        }
    }).then(result=>{

        var data = result.data
        var html = ""
        var url2 = CommonURL+"/ns/"+NAMESPACE+"/resources/sshKey"
        var spec_id = data.sshKeyId
       
        $.ajax({
           url: url2,
           async:false,
           type:'GET',
           beforeSend : function(xhr){
            xhr.setRequestHeader("Authorization", apiInfo);
            xhr.setRequestHeader("Content-type","application/json");
        },
          

       }).done(function(result){
        var res = result.sshKey
       console.log("sshKey info :",res);
        for(var k in res){
            if(res[k].id == spec_id){
             html += '<tr>'
                    +'<th scope="colgroup" rowspan="3" class="text-right"><i class="fas fa-key"></i>Access(SSH Key)</th>'
                    +'<th scope="colgroup" class="text-right">Key Name</th>'
                    +'<td  colspan="1">'+res[k].cspSshKeyName+'</td>'
                    +'<th scope="colgroup" class="text-right">SSH Key ID</th>'
                    +'<td colspan="1">'+res[k].id+'</td>'
                  
                    +'</tr>'
                    +'<tr>'
                    +'<th scope="colgroup" class="text-right">Description</th>'
                    +'<td colspan="3">'+res[k].description+'</th>'
                    +'</tr>'
                   
            }
        } 
        $("#sshKey").empty();
        $("#sshKey").append(html);

    })
      
           
        
    })

}