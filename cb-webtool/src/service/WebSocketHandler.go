package service

import (
	//"github.com/cloud-barista/cb-webtool/src/model"
	//tbmcis "github.com/cloud-barista/cb-webtool/src/model/tumblebug/mcis"
	"log"
	"time"

	echosession "github.com/go-session/echo-session"
	"github.com/labstack/echo"

	modelsocket "github.com/cloud-barista/cb-webtool/src/model/websocket"
)

// WebSocket에 전달할 Message Set
func SetWebsocketMessage(taskType string, taskKey string, lifeCycle string, status string, c echo.Context) {
	store := echosession.FromContext(c)
	socketDataStore, isStoreOk := store.Get("socketdata")
	socketDataMap := map[string]modelsocket.WebSocketMessage{}
	if !isStoreOk {
	} else {
		socketDataMap = socketDataStore.(map[string]modelsocket.WebSocketMessage) // 없으면 생성
	}

	websocketMessage := modelsocket.WebSocketMessage{}

	websocketMessage.TaskType = taskType
	websocketMessage.TaskKey = taskKey
	websocketMessage.LifeCycle = lifeCycle
	websocketMessage.Status = status
	websocketMessage.ProcessTime = time.Now()

	socketDataMap[taskKey] = websocketMessage
	store.Set("socketdata", socketDataMap)
	store.Save()
	log.Println("setsocketdata" + taskKey + " :  " + lifeCycle + " " + status)
}

// 전송여부를 set 하는데... 시간을 기준으로 가져올 것으로 필요없을 것.
func SetWebsocketMessageBySend(taskKey string, hasSend bool, c echo.Context) {
	store := echosession.FromContext(c)
	socketDataStore, isStoreOk := store.Get("socketdata")
	if isStoreOk {
		socketDataMap := socketDataStore.(map[string]modelsocket.WebSocketMessage)
		websocketMessage := socketDataMap[taskKey]

		websocketMessage.Send = hasSend
		socketDataMap[taskKey] = websocketMessage
		store.Set("socketdata", socketDataMap)
		store.Save()
	} else {
		log.Println("SetWebsocketMessageBySend is not Ok ")
	}
}

// TaskKey에 해당하는 값 조회 : 요청, 완료 값이 return
func GetWebsocketMessageByTaskKey(taskType string, taskKey string, c echo.Context) map[int64]modelsocket.WebSocketMessage {
	store := echosession.FromContext(c)
	socketDataStore, ok := store.Get("socketdata")
	returnSocketDataMap := map[int64]modelsocket.WebSocketMessage{}
	// returnWebsocketMessage := modelsocket.WebSocketMessage{}

	if ok {
		socketDataMap := socketDataStore.(map[int64]modelsocket.WebSocketMessage)
		for key, val := range socketDataMap {
			log.Println("show socketData with key : getsocketdata ", key, val)
			if val.TaskKey == taskKey {
				returnSocketDataMap[key] = val
				log.Println("show socketData with key by send : getsocketdata ", key, val)
			}
		}
		// }
	} else {
		log.Println("socketDataStore is not Ok ")
	}
	return returnSocketDataMap
}

// 전송 상태에 따른 값 목록 조회. sendMessage==false 이면 전송 전 data목록만 :: 시간을 param으로 하므로 필요 없을 것. deprecated.
//func GetWebsocketMessageBySend(send bool, c echo.Context) map[int64]modelsocket.WebSocketMessage {
func GetWebsocketMessageBySend(send bool, c echo.Context) []modelsocket.WebSocketMessage {
	store := echosession.FromContext(c)
	socketDataStore, ok := store.Get("socketdata")
	//websocketMessageMap := map[int64]modelsocket.WebSocketMessage{}
	socketResultList := []modelsocket.WebSocketMessage{}
	if ok {
		socketDataMap := socketDataStore.(map[int64]modelsocket.WebSocketMessage)
		for key, val := range socketDataMap {
			log.Println("show socketData with key : getsocketdata ", key, val)
			socketMessage := socketDataMap[key]
			// socketMessage.CallTime = time.Now().UnixNano()
			//if val.Send == send {
			//	socketMessageMap[key] = socketMessage
			//	log.Println("show socketData with key by send : getsocketdata ", key, val)
			//}
			//socketResultList = append(socketResultList, modelsocket.SocketMessage{SaveTime: key, Message: socketMessage})
			socketResultList = append(socketResultList, socketMessage)
		}

	} else {
		log.Println("socketDataStore is not Ok ")
	}

	return socketResultList
}

// 특정 시점 이후의 data만 추출
//func GetWebsocketMessageByProcessTime(beginTime time.Time, c echo.Context) map[int64]modelsocket.WebSocketMessage {
func GetWebsocketMessageByProcessTime(beginTime int64, c echo.Context) []modelsocket.WebSocketMessage {
	store := echosession.FromContext(c)
	socketDataStore, ok := store.Get("socketdata")
	//websocketMessageMap := []modelsocket.WebSocketMessage{}
	socketResultList := []modelsocket.WebSocketMessage{}
	if ok {
		socketDataMap := socketDataStore.(map[int64]modelsocket.WebSocketMessage)
		for key, val := range socketDataMap {
			log.Println("show socketData with key : getsocketdata ", key, val)
			websocketMessage := socketDataMap[key]
			// websocketMessage.CallTime = time.Now().UnixNano() //
			//log.Println( websocketMessage.ProcessTime )
			//websocketMessagelog.Println("beginTime : ", beginTime )
			//log.Println("is after : ", websocketMessage.ProcessTime.After(beginTime) )
			//if websocketMessage.ProcessTime.After(beginTime) {
			log.Println(" key  ", key)
			log.Println(" beginTime ", beginTime)
			log.Println(" key > beginTime ", key > beginTime)
			if key > beginTime {
				//websocketMessageMap[key] = websocketMessage
				log.Println("show socketData with key by ProcessTime : getsocketdata ", key, val)
				socketResultList = append(socketResultList, websocketMessage)
			}
		}
	} else {
		log.Println("socketDataStore is not Ok ")
	}

	return socketResultList
}

// taskType : mcis/vm/mcks ...
// lifecycle : create, suspend, resume. ....
// taskKey :
// status : requested, processing, failed, completed
// eccossion에 socketdata 에 추가. key는 timestamp인데 unixNanoTime(int64) 사용
func StoreWebsocketMessage(taskType string, taskKey string, lifeCycle string, requestStatus string, c echo.Context) {
	store := echosession.FromContext(c)
	socketDataStore, isStoreOk := store.Get("socketdata")
	socketDataMap := map[int64]modelsocket.WebSocketMessage{}
	if !isStoreOk {
	} else {
		socketDataMap = socketDataStore.(map[int64]modelsocket.WebSocketMessage) // 없으면 생성
	}

	websocketMessage := modelsocket.WebSocketMessage{}

	websocketMessage.TaskType = taskType
	websocketMessage.TaskKey = taskKey
	websocketMessage.LifeCycle = lifeCycle
	websocketMessage.Status = requestStatus
	websocketMessage.ProcessTime = time.Now()

	socketDataMap[time.Now().UnixNano()] = websocketMessage
	store.Set("socketdata", socketDataMap)
	store.Save()
	log.Println("setsocketdata" + taskKey + " :  " + lifeCycle + " " + requestStatus)
}

// 일정 시간이 지난 data는 제거. : 0이면 기본값(24), 0보다 크면 음수로 바꾸어 계산.
func ClearWebsocketMessage(expireHour int, c echo.Context) {
	store := echosession.FromContext(c)

	if expireHour == 0 {
		expireHour = 24
	} else if expireHour > 0 {
		expireHour = -1 * expireHour
	}
	t := time.Now()
	d2 := t.Add(time.Hour * time.Duration(expireHour)) // expire 시간이 지난 것들은 삭제

	renewSocketDataMap := GetWebsocketMessageByProcessTime(d2.UnixNano(), c)
	//renewSocketDataMap := GetWebsocketMessageByProcessTime(d2, c)

	store.Set("socketdata", renewSocketDataMap)
	store.Save()
	log.Println("renew setsocketdata before :", d2)
}
