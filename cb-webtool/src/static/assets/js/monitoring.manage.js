// MCIS Control
function life_cycle(tag,type,mcis_id,mcis_name,vm_id,vm_name){
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
        if(status == 200){
            
            alert(message);
            location.reload(true);
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


 function show_mcis(url){
    console.log("Show mcis Url : ",url)
   var html = "";
   var apiInfo = ApiInfo
    axios.get(url,{
        headers:{
            'Authorization': apiInfo
        }
    }).then(result=>{
       var data = result.data;
       if(!data.mcis){
          location.href = "/MCIS/reg";
          return;
       }
        console.log("showmcis Data : ",data)
        var html = "";
        var mcis = data.mcis;
        var len = mcis.length;
        var count = 0;
        
        
        for(var i in mcis){
            var sta = mcis[i].status;
            var badge = "";
            var status = sta.toLowerCase()
            var vms = mcis[i].vm
            var vm_len = 0
           
            if(vms){
               vm_len = vms.length
            }
            

           console.log("mcis Status 1: ", mcis[i].status)
           console.log("mcis Status 2: ", status)
            if(status == "running"){
               badge += '<span class="badge badge-pill badge-success">RUNNING</span>'
            }else if(status == "include-notdefinedstatus" ){
               badge += '<span class="badge badge-pill badge-warning">WARNING</span>'
            }else if(status == "suspended"){
               badge += '<span class="badge badge-pill badge-warning">SUSPEND</span>'
            }else if(status == "terminate"){
               badge += '<span class="badge badge-pill badge-dark">TERMINATED</span>'
            }else{
               badge += '<span class="badge badge-pill badge-warning">'+status+'</span>'
            }
            count++;
            if(count == 1){

            }
            html += '<tr id="tr_id_'+count+'" >'
             +'<td class="text-center">'
             +'<div class="form-input">'
             +'<span class="input">'
             +'<input type="checkbox" class="chk" id="chk_'+count+'" value="'+mcis[i].id+'" item="'+mcis[i].name+'"><i></i></span></div>'
             +'</td>'
             +'<td>'
             +badge
             +'</td>'
             +'<td><a href="#!" onclick="show_vmList(\''+mcis[i].id+'\')">'+mcis[i].name+'</a></td>'
             +'<td>'+vm_len+'</td>'
             +'<td>'+vm_len+'</td>'
             +'<td>0</td>'
             +'<td>'+short_desc(mcis[i].description)+'</td>'
             +'<td>'
             +'<button type="button" class="btn btn-icon dropdown-toggle" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">'
             +'<i class="fas fa-edit"></i>'
             +'<div class="dropdown-menu dropdown-menu-right" aria-labelledby="btnGroupDrop1">'
                 +'<a class="dropdown-item" href="#!" onclick="life_cycle(\'mcis\',\'resume\',\''+mcis[i].id+'\',\''+mcis[i].name+'\')">Resume</a>'
                 +'<a class="dropdown-item" href="#!" onclick="life_cycle(\'mcis\',\'suspend\',\''+mcis[i].id+'\',\''+mcis[i].name+'\')">Suspend</a>'
                 +'<a class="dropdown-item" href="#!" onclick="life_cycle(\'mcis\',\'reboot\',\''+mcis[i].id+'\',\''+mcis[i].name+'\')">Reboot</a>'
                 +'<a class="dropdown-item" href="#!" onclick="life_cycle(\'mcis\',\'terminate\',\''+mcis[i].id+'\',\''+mcis[i].name+'\')">Terminate</a>'
             +'</div>'
             +'</button>'
            +'</td>'
            +'</tr>';
       }
       
       $("#table_1").empty();
       $("#table_1").append(html);
       console.log("VM LEN  :" ,vm_len);
      // show_card(mcis[0].id);
       
      if(vm_len > 0){
       show_vmList(mcis[0].id);
      }else{
       show_vmList("");
      }
       
      
       
       //fnMove("table_1");
       $("#mcis_id").val(mcis[0].id)
       $("#mcis_name").val(mcis[0].name)
   });
}
function show_vmList(mcis_id){
    var apiInfo = ApiInfo;
    var url = CommonURL+"/ns/"+NAMESPACE+"/mcis/"+mcis_id;
    console.log("MCIS Mangement mcisID : ",mcis_id);
    if(mcis_id){
        $.ajax({
            type:'GET',
            url:url,
            beforeSend : function(xhr){
                xhr.setRequestHeader("Authorization", apiInfo);
                xhr.setRequestHeader("Content-type","application/json");
            },
        // async:false,
            success:function(data){
                var vm = data.vm
                var mcis_name = data.name 
                $("#mcis_id").val(mcis_id)
                $("#mcis_name").val(mcis_name)
                var html = "";
                console.log("VM DATA : ",vm)
                for(var i in vm){
                    var sta = vm[i].status;
                    
                
                    var status = sta.toLowerCase()
                    console.log("VM Status : ",status)
                    var configName = vm[i].connectionName
                    console.log("outer vm configName2 : ",configName)
                    var count = 0;
                    console.log("Spider URL : ",SpiderURL)
                    $.ajax({
                        url: SpiderURL+"/connectionconfig",
                        async:false,
                        type:'GET',
                        beforeSend : function(xhr){
                            xhr.setRequestHeader("Authorization", apiInfo);
                            xhr.setRequestHeader("Content-type","application/json");
                        },
                        
                        success : function(data2){
                            var badge = "";
                           
                            res = data2.connectionconfig
                            for(var k in res){
                                // console.log(" i value is : ",i)
                                // console.log("outer config name : ",configName)
                                // console.log("Inner ConfigName : ",res[k].ConfigName)
                                if(res[k].ConfigName == vm[i].connectionName){
                                    var provider = res[k].ProviderName
                                    console.log("Provider : ",provider);
                                    
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
                                    html += '<tr id="tr_id_'+count+'" >'
                                    +'<td class="text-center">'
                                    +'<div class="form-input">'
                                    +'<span class="input">'
                                    +'<input type="checkbox" item="'+mcis_name+'"    mcisid="'+mcis_id+'" class="chk2" id="chk2_'+count+'" value="'+vm[i].id+'|'+mcis_id+'"><i></i></span></div>'
                                    +'</td>'
                                    +'<td>'
                                    +badge
                                    +'</td>'
                                    +'<td><a href="#!" onclick="showMonitoring(\''+mcis_id+'\',\''+vm[i].id+'\');">'+vm[i].name+'</a></td>'
                        
                                    +'<td>'+provider+'</td>'
                                    +'<td>'+vm[i].region.Region+'</td>'
                                    +'<td>'+vm[i].connectionName+'</td>'
                                    +'<td>OS Type</td>'
                                    +'<td>'+vm[i].publicIP+'</td>'
                                    +'<td>'+short_desc(vm[i].description)+'</td>'
                                    +'<td>'
                                    +'<button type="button" class="btn btn-icon dropdown-toggle" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">'
                                    +'<i class="fas fa-edit"></i>'
                                    +'<div class="dropdown-menu dropdown-menu-right" aria-labelledby="btnGroupDrop1">'
                                        +'<a class="dropdown-item" href="#!" onclick="life_cycle(\'vm\',\'resume\',\''+mcis_id+'\',\''+mcis_name+'\',\''+vm[i].id+'\',\''+vm[i].name+'\')">Resume</a>'
                                        +'<a class="dropdown-item" href="#!" onclick="life_cycle(\'vm\',\'suspend\',\''+mcis_id+'\',\''+mcis_name+'\',\''+vm[i].id+'\',\''+vm[i].name+'\')">Suspend</a>'
                                        +'<a class="dropdown-item" href="#!" onclick="life_cycle(\'vm\',\'reboot\',\''+mcis_id+'\',\''+mcis_name+'\',\''+vm[i].id+'\',\''+vm[i].name+'\')">Reboot</a>'
                                        +'<a class="dropdown-item" href="#!" onclick="life_cycle(\'vm\',\'terminate\',\''+mcis_id+'\',\''+mcis_name+'\',\''+vm[i].id+'\',\''+vm[i].name+'\')">Terminate</a>'
                                    +'</div>'
                                    +'</button>'
                                    +'</td>'
                                    +'</tr>';
                                }
                                
                                
                                }
                                $("#table_2").empty();
                                $("#table_2").append(html);
                                $("#vm_detail").hide();
                                fnMove("table_2");

                        }

                    })
                    
                    }
            }
        })
    }else{
        $("#table_2").empty();
        $("#table_2").append("<td colspan='9'>Does not Exist</td>");
    }
            
    
 }
  
  function show_card(mcis_id){
      var url = CommonURL+"/ns/"+NAMESPACE+"/mcis/"+mcis_id;
      var html = "";
     var apiInfo = ApiInfo
    axios.get(url,{
        headers:{
            'Authorization': apiInfo
        }
    }).then(result=>{
         var data = result.data
         console.log("show card data : ",result)
         var vm_cnt = data.vm
         if(vm_cnt){
             vm_cnt = vm_cnt.length;
         }else{
             vm_cnt = 0;
         }
         
         
             html += '<div class="col-xl-12 col-lg-12">'
                     +'<div class="card card-stats mb-12 mb-xl-0">'
                     +'<div class="card-body">'
                     +'<div class="row">'
                     +'<div class="col">'
                     +'<h5 class="card-title text-uppercase text-muted mb-0">'+data.name+'</h5>'
                     +'<span class="h2 font-weight-bold mb-0">350,897</span>'
                     +'</div>'
                     +'<div class="col-auto">'
                     +'<div class="icon icon-shape bg-danger text-white rounded-circle shadow">'
                     //+'<i class="fas fa-chart-bar"></i>'
                     +vm_cnt
                     +'</div>'
                     +'</div>'
                     +'</div>'
                     +'<p class="mt-3 mb-0 text-muted text-sm">'
                     +'<span class="text-success mr-2"><i class="fa fa-arrow-up"></i> 3.48%</span>'
                     +'<span class="text-nowrap">Since last month</span>'
                     +'</p>'
                     +'</div>'
                     +'</div>'
                     +'</div>';
         
         $("#card").empty()
         $("#card").append(html)
         if(vm_cnt == 0){
             show_vmList("")
         }else{
             show_vmList(mcis_id)
         }
         
        
     })
  }
 function show_vm(mcis_id,vm_id){
     var url = CommonURL+"/ns/"+NAMESPACE+"/mcis/"+mcis_id+"/vm/"+vm_id;
     var html = "";
    var apiInfo = ApiInfo
    axios.get(url,{
        headers:{
            'Authorization': apiInfo
        }
    }).then(result=>{
        var data = result.data
        console.log("show card result : ",result)
                   
            html += '<div class="col-xl-12 col-lg-12">'
                    +'<div class="card card-stats mb-12 mb-xl-0">'
                    +'<div class="card-body">'
                    +'<div class="row">'
                    +'<div class="col">'
                    +'<h5 class="card-title text-uppercase text-muted mb-0">'+data.name+'</h5>'
                    +'<span class="h2 font-weight-bold mb-0">350,897</span>'
                    +'</div>'
                    +'<div class="col-auto">'
                    +'<div class="icon icon-shape bg-danger text-white rounded-circle shadow">'
                    //+'<i class="fas fa-chart-bar"></i>'
                    +"vm_cnt"
                    +'</div>'
                    +'</div>'
                    +'</div>'
                    +'<p class="mt-3 mb-0 text-muted text-sm">'
                    +'<span class="text-success mr-2"><i class="fa fa-arrow-up"></i> 3.48%</span>'
                    +'<span class="text-nowrap">Since last month</span>'
                    +'</p>'
                    +'</div>'
                    +'</div>'
                    +'</div>';
        
        $("#card").empty()
        $("#card").append(html)
       
    })
 }

 function show_vm(mcis_id,vm_id){
    show_vmDetailList(mcis_id, vm_id);
    show_vmSpecInfo(mcis_id, vm_id);
    show_vmNetworkInfo(mcis_id, vm_id);
    show_vmSecurityGroupInfo(mcis_id, vm_id);
    show_vmSSHInfo(mcis_id, vm_id);
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


    })
    if(cnt < 1 ){
        alert("등록할 대상을 선택해 주세요.");
        return;
    }

    if(cnt == 1){
       console.log("mcis_id ; ",mcis_id)
        var url = "/MCIS/reg/"+mcis_id+"/"+mcis_name
        
        if(confirm("등록하시겠습니까?")){
            location.href = url;
        }
    }

    if(cnt >1){
        alert("한개씩만 등록 가능합니다.")
        return;
    }
 }
 function vm_delete(){
    
    var cnt = 0;
    var vm_id = "";
    var mcis_id ="";
    var apiInfo = ApiInfo;
    $(".chk2").each(function(){
        if($(this).is(":checked")){
            //alert("chk");
            cnt++;
            id = $(this).val(); 
            idArr = id.split ("|")  
            vm_id = idArr[0]
            mcis_id = idArr[1]    
        }
    })
    if(cnt < 1 ){
        alert("삭제할 대상을 선택해 주세요.");
        return;
    }

    if(cnt == 1){
       console.log("mcis_id ; ",vm_id)
        var url = CommonURL+"/ns/"+NAMESPACE+"/mcis/"+mcis_id+"/vm/"+vm_id
        
        if(confirm("삭제하시겠습니까?")){
         axios.delete(url,{
            headers :{
                'Content-type': 'application/json',
                'Authorization': apiInfo,
                }
         }).then(result=>{
             var data = result.data
             console.log(result);
             if(result.status == 200){
                 alert(data.message)
                 location.reload(true)
             }
         })
        }
    }

    if(cnt >1){
        alert("한개씩만 삭제 가능합니다.")
        return;
    }
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
 function mappingMetric(obj){
    var name = obj.name
    var columnArr = obj.columns
    var valuesArr = obj.values
    var valuesCnt = valuesArr.length
    var objArr = new Array();
    for(var i in  valuesArr){
       var newObject = {}
        for(var k in valuesArr[i]){
            var key = columnArr[k]
            var value = valuesArr[i][k]
            newObject[key] = value
        }
        objArr.push(newObject)
    }
    console.log("Mapping Metric : ",objArr);
    return objArr
}

function show_vmDetailList(mcis_id, vm_id){
    url = CommonURL+"/ns/"+NAMESPACE+"/mcis/"+mcis_id+"/vm/"+vm_id
    var apiInfo = ApiInfo
    axios.get(url,{
        headers:{
            'Authorization': apiInfo
        }
    }).then(result=>{
        var data = result.data
        console.log("show vmDetail List data : ",data)
        var html = ""
        $.ajax({
           url: SpiderURL+"/connectionconfig",
           async:false,
           type:'GET',
           beforeSend : function(xhr){
            xhr.setRequestHeader("Authorization", apiInfo);
            xhr.setRequestHeader("Content-type","application/json");
        },
           success : function(data2){
               res = data2.connectionconfig
               var provider = "";
               for(var k in res){
                   if(res[k].ConfigName == data.connectionName){
                       provider = res[k].ProviderName
                       console.log("Inner Provider : ",provider)
                   }
               }
               html += '<tr>'
                   +'<th scope="colgroup"rowspan="10">Infra - Server</th>'

                   +'<th scope="colgroup">Server ID</th>'
                   +'<td  colspan="3">'+data.id+'</td>'
                   +'</tr>'

                   +'<tr>'

                   +'<th scope="colgroup">Cloud Provider</th>'
                   +'<td colspan="3">'+provider+'</td>'

                   
                   +'</tr>'

                   +'<tr>'
                   +'<th scope="colgroup">CP VMID</th>'
                   +'<td  colspan="3">'+data.id+'</td>'
                   +'</tr>'

                   +'<tr>'
                   +'<th scope="colgroup">Region</th>'
                   +'<td  colspan="3">'+data.region.Region+'</td>'
                   +'</tr>'

                   
                   +'<tr>'
                   +'<th scope="colgroup">Public IP</th>'
                   +'<td  colspan="3">'+data.publicIP+'</td>'
                   +'</tr>'

                   +'<tr>'
                   +'<th scope="colgroup">Public DNS</th>'
                   +'<td  colspan="3">'+data.publicDNS+'</td>'
                   +'</tr>'

                   +'<tr>'
                   +'<th scope="colgroup">Private IP1234</th>'
                   +'<td colspan="3">'+data.privateIP+'</td>'
                   +'</tr>'

                   +'<tr>'
                   +'<th scope="colgroup">Private DNS</th>'
                   +'<td colspan="3">'+data.privateDNS+'</td>'
                   +'</tr>'

                   +'<tr>'
                   +'<th scope="colgroup">Server Status</th>'
                   +'<td colspan="3">'+data.status+'</td>'
                   +'</tr>'
                 
               $("#vm").empty();
               $("#vm").append(html);
               fnMove("vm_detail");

           }

       })
      
           
        
    })

}

function show_vmDetailInfo(mcis_id, vm_id){
   var url = CommonURL+"/ns/"+NAMESPACE+"/mcis/"+mcis_id+"/vm/"+vm_id
   var apiInfo = ApiInfo
    axios.get(url,{
        headers:{
            'Authorization': apiInfo
        }
    }).then(result=>{
       var data = result.data
       var html = ""
       $.ajax({
          url: SpiderURL+"/connectionconfig",
          async:false,
          type:'GET',
          success : function(data){
              
              var provider = "";
              res = data.connectionconfig
              for(var k in res){
                  if(res[k].ConfigName == data.connectionName){
                      provider = res[k].ProviderName
                      console.log("Inner Provider : ",provider)
                  }
              }
              html += '<tr>'
                  +'<th scope="colgroup"rowspan="6">Resource-VM</th>'
                  +'<th scope="colgroup">cloud Provider</th>'
                  +'<td colspan="3">'+provider+'</td>'
                  +'</tr>'
                  +'<tr>'
                  +'<th scope="colgroup">VM ID</th>'
                  +'<td  colspan="3">'+data.id+'</td>'
                  +'</tr>'
                  +'<tr>'
                  +'<th scope="colgroup">Region</th>'
                  +'<td  colspan="3">'+data.region.Region+'</td>'
                  +'</tr>'
                  +'<tr>'
                  +'<th scope="colgroup">Zone</th>'
                  +'<td  colspan="3">'+data.region.Zone+'</td>'
                  +'</tr>'
                  +'<tr>'
                  +'<th scope="colgroup">PublicIP</th>'
                  +'<td  colspan="3">'+data.publicIP+'</td>'
                  +'</tr>'
                  +'<tr>'
                  +'<th scope="colgroup">PrivateIP</th>'
                  +'<td colspan="3">'+data.privateIP+'</td>'
                  +'</tr>'
                  +'</tbody>'
                  +'<tbody>'
                  +'<tr>'
                  +'<th scope="colgroup" rowspan="3">VM Meta</th>'
                  +'<th scope="colgroup">VM ID</th>'
                  +'<td colspan="3">'+data.cspViewVmDetail.Id+'</td>'
                  +'</tr>'
                  +'<tr>'
                  +'<th scope="colgroup">VM NAME</th>'
                  +'<td  colspan="3">'+data.cspViewVmDetail.Name+'</td>'
                  +'</tr>'
                  

                
              $("#vm").empty();
              $("#vm").append(html);

          }

      })
     
          
       
   })

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
        console.log("show vmSpecInfo Data : ",data)
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
           success : function(result){
               var res = result.spec
              console.log("spec data from tumble : ",res)
               for(var k in res){
                   if(res[k].id == spec_id){
                    html += '<tr>'
                          
                           +'<th scope="colgroup" rowspan="5">Server Spec</th>'
                           +'<th scope="colgroup">vCPU</th>'
                           +'<td colspan="3">'+res[k].num_vCPU+'vcpu</td>'
                           +'</tr>'
                           +'<tr>'
                           +'<th scope="colgroup">Memory(Ghz)</th>'
                           +'<td  colspan="3">'+res[k].mem_GiB+'GiB</td>'
                           +'</tr>'
                           +'<tr>'
                           +'<th scope="colgroup">Disk (GB)</th>'
                           +'<td colspan="3">'+res[k].storage_GiB+'GiB</th>'
                           +'</tr>'   
                           +'<tr>'
                           +'<th scope="colgroup">Cost($) / Hour </th>'
                           +'<td colspan="3">'+res[k].cost_per_hour+'</td>'
                           +'</tr>'
                           +'<tr>'
                           +'<th scope="colgroup">OsType</th>'
                           +'<td  colspan="3">'+res[k].os_type+'</td>'
                           +'</tr>'
                   }
               } 
               $("#vm_spec").empty();
               $("#vm_spec").append(html);

           }

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
          success : function(result){
              var res = result.network
             
              for(var k in res){
                  if(res[k].id == spec_id){
                   html += '<tr>'
                          +'<th scope="colgroup" rowspan="3">vNetwork</th>'
                          +'<th scope="colgroup">NetworkID</th>'
                          +'<td colspan="3">'+res[k].cspNetworkId+'</td>'
                          +'</tr>'
                          +'<tr>'
                          +'<th scope="colgroup">Network Name</th>'
                          +'<td  colspan="3">'+res[k].cspNetworkName+'</td>'
                          +'</tr>'
                          +'<tr>'
                          +'<th scope="colgroup">Cidr Block</th>'
                          +'<td colspan="3">'+res[k].cidrBlock+'</th>'
                          +'</tr>'
                         
                  }
              } 
              $("#vm_vnetwork").empty();
              $("#vm_vnetwork").append(html);

          }

      })
     
          
       
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
       var html = ""
       // var url2 = "/ns/"+NAMESPACE+"/resources/securityGroup"
       var spec_id = data.securityGroupIds
       var cnt = spec_id.length
       html += '<tr>'
            +'<th scope="colgroup" colspan="'+cnt+'">SecurityGroup</th>'
            +'<th scope="colgroup" colspan="'+cnt+'">SecurityGroupID</th>'
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
          success : function(result){
              var res = result.sshKey
             
              for(var k in res){
                  if(res[k].id == spec_id){
                   html += '<tr>'
                          +'<th scope="colgroup" rowspan="3">SSH KEY</th>'
                          +'<th scope="colgroup">SSH Key ID</th>'
                          +'<td colspan="3">'+res[k].id+'</td>'
                          +'</tr>'
                          +'<tr>'
                          +'<th scope="colgroup">Key Name</th>'
                          +'<td  colspan="3">'+res[k].cspSshKeyName+'</td>'
                          +'</tr>'
                          +'<tr>'
                          +'<th scope="colgroup">Description</th>'
                          +'<td colspan="3">'+res[k].description+'</th>'
                          +'</tr>'
                         
                  }
              } 
              $("#sshKey").empty();
              $("#sshKey").append(html);

          }

      })
     
          
       
   })

}

