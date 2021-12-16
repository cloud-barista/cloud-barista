html template controller

main.go 에서 선언한 method 구현.
필요한 서비스들을 호출하여 로직을 수행한 뒤 return

. 화면이 있는 경우 : template에서 master file 경로를 설정하여 return
. 화면이 없는 경우 : json으로 return

return 시 구현해야 하는 parameter : message, status
model.WebStatus{StatusCode: 500, Message: err.Error()}

if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {
log.Println(" respStatus ", respStatus)
return c.JSON(http.StatusBadRequest, map[string]interface{}{
"message": respStatus.Message,
"status": respStatus.StatusCode,
})
}

-> 최종 변경 : return시 error로 send, error code 도 return받는 respStatus
if respStatus.StatusCode != 200 && respStatus.StatusCode != 201 {

    return c.JSON(respStatus.StatusCode, map[string]interface{}{
        "error":  respStatus.Message,
        "status": respStatus.StatusCode,
    })

}


참고 : 
net.Dial을 이용하여 Keep Alive Check를 각 Framework로 하려했으나, IP_PORT 의 응답만 되므로 굳이 구현 필요가 없음.
(Dial로 연결된 후 server로부터 응답을 받을 때까지 conn.read()에서 대기)
WebsocketController.go 에 구현되어 있으나 사용하지 않음.