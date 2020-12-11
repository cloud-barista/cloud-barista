function life_cycle(tag,type, mcis_id,mcis_name,vm_id,vm_name){
    var url = ""
    var nameSpace = NAMESPACE;
    var message = ""
    
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
        console.log("result Message : ",result.data.message)
        if(status == 200){
            alert(message);
            location.reload(true);
        }
     })
 }

//  function suspend(tag,mcis_id,mcis_name,vm_id,vm_name){
//     var url = ""
//     var nameSpace = NAMESPACE;
//     var message = ""
    
//     if(tag == "mcis"){
//         url ="/ns/"+nameSpace+"/mcis/"+mcis_id+"?action=suspend"
//         message = mcis_name+" suspend complete!."
//      }else{
//         url ="/ns/"+nameSpace+"/mcis/"+mcis_id+"/vm/"+vm_id+"?action=suspend"
//         message = vm_name+" suspend complete!."
//      }

//      var apiInfo = ApiInfo
    // axios.get(url,{
    //     headers:{
    //         'Authorization': apiInfo
    //     }
    // })then(result=>{
//         var status = result.status
//         console.log("result Message : ",result.data.message)
//         if(status == 200){
//             alert(message);
//             location.reload(true);
//         }
//      })
//  }

//  function reboot(tag,mcis_id,mcis_name,vm_id,vm_name){
//     var url = ""
//     var nameSpace = NAMESPACE;
//     var message = ""
    
//     if(tag == "mcis"){
//         url ="/ns/"+nameSpace+"/mcis/"+mcis_id+"?action=reboot"
//         message = mcis_name+" reboot complete!."
//      }else{
//         url ="/ns/"+nameSpace+"/mcis/"+mcis_id+"/vm/"+vm_id+"?action=reboot"
//         message = vm_name+" reboot complete!."
//      }

//      var apiInfo = ApiInfo
    // axios.get(url,{
    //     headers:{
    //         'Authorization': apiInfo
    //     }
    // })then(result=>{
//         var status = result.status
//         console.log("result Message : ",result.data.message)
//         if(status == 200){
//             alert(message);
//             location.reload(true);
//         }
//      })
//  }

//  function terminate(tag,mcis_id,mcis_name,vm_id,vm_name){
//     var url = ""
//     var nameSpace = NAMESPACE;
//     var message = ""
    
//     if(tag == "mcis"){
//         url ="/ns/"+nameSpace+"/mcis/"+mcis_id+"?action=terminate"
//         message = mcis_name+" terminate complete!."
//      }else{
//         url ="/ns/"+nameSpace+"/mcis/"+mcis_id+"/vm/"+vm_id+"?action=terminate"
//         message = vm_name+" terminate complete!."
//      }

//      var apiInfo = ApiInfo
    // axios.get(url,{
    //     headers:{
    //         'Authorization': apiInfo
    //     }
    // })then(result=>{
//         var status = result.status
//         console.log("result Message : ",result.data.message)
//         if(status == 200){
//             alert(message);
//             location.reload(true);
//         }
//      })
//  }