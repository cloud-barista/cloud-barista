// funtcion requestAjax(url, method, data){
//     console.log("Request URL : ",url)
//     var met = method.toLowerCase
//     $.ajax({
//         url : url,
//         type: method,
//         data: data

//     }).then(function(result){
//         console.log(result)
//     })
// }



function getNameSpace(){
    var url = CommonURL+"/ns"
    axios.get(url).then(result =>{
        var data = result.data.ns
        var namespace = ""
        for( var i in data){
            if(i == 0 ){
                namespace = data[i].id
            }
        }
        $("#namespace1").val(namespace);

    })
}

function fnMove(target){
    var offset = $("#" + target+"").offset();
    console.log("FnMove offset : ",offset)
    $('html, body').animate({scrollTop : offset.top}, 400);
}

function getVMStatus(vm_name, connection_name){
    var url = "/vmstatus/"+vm_name+"?connection_name="+connection_name

    $.ajax({
        url: url,
        async:false,
        type:'GET',
        success : function(res){
            var vm_status = res.Status 

        }
    })
}