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

//폼의 Validation을 체크함.
//<input> tag에 "required" 옵션이 추가된 항목의 값이 공백인 경우 false  그렇지 않으면 true 리턴
function chkFormValidate(formObj) {
    try{
        var objs = formObj.find("[required]");
        //alert(objs.length)

        // required 옵션이 체크된 필드 들의 값을 조회 함.(현재는 Text 필드만 가능)
        for(var i = 0; i < objs.length; i++) {
            if(objs.eq(i).val() == '') {
                alert("Please enter a value.");
                objs.eq(i).focus();
                return false;
            }
        }
        return true;
    } catch (e) {
        alert(e);
        return false;
    }
}

function getOSType(image_id){
    var url = CommonURL+"/ns/"+NAMESPACE+"/resources/image/"+image_id
    return axios.get(url).then(result=>{
        var data = result.data
        var osType = data.guestOS
        console.log("Image Data : ",data);
        return osType;
        })
}
function checkNS(){
    var url = CommonURL+"/ns";
    axios.get(url).then(result =>{
        var data = result.data.ns
       if(!data){
        alert("NameSpace가 등록되어 있지 않습니다.\n등록페이지로 이동합니다.")
        location.href ="/NS/reg";
        return;
       }else{
           return;
       }
    })

}
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
function cancel_btn(){
    if(confirm("Cancel it?")){
        history.back();
    }else{
        return;
    }
}
function close_btn(){
    if(confirm("close it?")){
        $("#transDiv").hide();
    }else{
        return;
    }
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