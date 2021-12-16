package websocket

// client와 서버가 통신하면서 비동기로 처리된 data 알림용
// echo session에서 값을 저장하는 key=socketdata 이며 map안에 WebSocketMessage 객체가 들어감  ex) store.Get("socketdata"), store.Set("socketdata", socketDataMap)
// 추후 필요시 userid 추가
// Send 는 기본 false 이며 전송이 되면 true 바꾼다.
// 생성에 관해서는 생성, 중단, 재가동, 삭제 등이 life cycle이므로 lifeCycle 상태값이 다르므로 해당 상태값을 사용한다. ex) create_request -> creating -> created or create completed 등
// 처리 상태는 해당 job에서 사용하는 상태를 사용. 임시로 요청상태(request), 진행중(ing), 완료(complete), 실패(fail)

import (
	"time"
)

type WebSocketMessage struct {
	TaskType    string    `json:"taskType"`    // 처리되는 Operation 구분. ex) mcis, mcks, vm, vnet, connection ....
	TaskKey     string    `json:"taskKey"`     // session에서 값을 찾기 위한 key   namespace || tasktype || id or name   ex) ns01||mcis||testmcis01    TODO : unique할까??  taskType에 대해서는 unique 할 듯.
	LifeCycle   string    `json:"lifeCycle"`   // 처리 구분 create, suspend, resume, terminate, delete ... 등
	Status      string    `json:"status"`      // 처리상태 ( request, ing, complete ...)  tb에서 보내는 상태에 따라 달라질 수 있음.
	Send        bool      `json:"send"`        // 메세지 전송여부(UI로 Push 했는지)
	Message     string    `json:"message"`     // 전달 할 Message
	Event       string    `json:"event"`       // 화면이 열렸을 때(open), 닫혔을 때(close), 요청(req), 응답(res)
	ProcessTime time.Time `json:"processTime"` // 상태변경 시간. 마지막으로 수신한 processTime 을 UI에서 가지고 해당 시간을 넘기면 lastProcessTime 이후꺼만 조회하면 될 듯.
}

type SocketMessage struct {
	SaveTime int64            `json:"saveTime"`
	Message  WebSocketMessage `json:"message"`
}

type SocketMessageList []SocketMessage
