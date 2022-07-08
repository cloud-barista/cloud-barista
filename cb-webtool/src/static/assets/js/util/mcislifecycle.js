
// MCIS 제어 : 선택한 MCIS내 vm들의 상태 변경 
// Dashboard 와 MCIS Manage 에서 같이 쓰므로
// callAAA -> mcisLifeCycle 호출 -> callBackAAA로 결과값전달
function mcisLifeCycle(mcisID, type) {

    var url = "/operation/manages/mcismng/proc/mcislifecycle";

    console.log("life cycle3 url : ", url);
    var message = "MCIS " + type + " complete!."
    var namespaceID = $('#topboxDefaultNameSpaceID').val();
    axios.post(url, {
        headers: {},
        namespaceID: namespaceID,
        mcisID: mcisID,
        queryParams: ["action=" + type, "force=false"]
    }).then(result => {
        console.log("mcisLifeCycle result : ", result);
        var status = result.status
        var data = result.data
        callbackMcisLifeCycle(status, data, type)
        // console.log("life cycle result : ",result)
        // console.log("result Message : ",data.message)
        // if(status == 200 || status == 201){

        //     alert(message);
        //     location.reload();
        //     //show_mcis(mcis_url,"");
        // }else{
        //     alert(status)
        //     return;
        // }
        // }).catch(function(error){
        //     // console.log(" display error : ",error);
        //     console.log(error.response.data);
        //     console.log(error.response.status);
        //     // console.log(error.response.headers); 
        //     var status = error.response.status;
        //     var data =  error.response.data

        //     callbackMcisLifeCycle(status, data, type)
        // });
    }).catch((error) => {
        console.warn(error);
        console.log(error.response)
        // var errorMessage = error.response.data.error;
        // commonErrorAlert(statusCode, errorMessage) 
        var errorMessage = error.response.data.error;
        var statusCode = error.response.status;
        commonErrorAlert(statusCode, errorMessage);
    });
}
////////////// MCIS Handling end //////////////// 