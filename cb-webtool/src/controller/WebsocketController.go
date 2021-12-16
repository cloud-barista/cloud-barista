package controller

import (
	"fmt"
	"log"

	// "io"
	// "net"
	"net/http"
	"os"

	"encoding/json"
	"strconv"
	"time"

	// "github.com/labstack/echo/v4"
	"github.com/labstack/echo"
	// "github.com/labstack/echo/v4/middleware"
	"github.com/gorilla/websocket"
	// "golang.org/x/net/websocket"

	// "github.com/cloud-barista/cb-webtool/src/service"
	echotemplate "github.com/foolin/echo-template"

	service "github.com/cloud-barista/cb-webtool/src/service"

	util "github.com/cloud-barista/cb-webtool/src/util"
)

var SpiderURL = os.Getenv("SPIDER_IP_PORT")
var TumbleBugURL = os.Getenv("TUMBLE_IP_PORT")
var DragonFlyURL = os.Getenv("DRAGONFLY_IP_PORT")
var LadyBugURL = os.Getenv("LADYBUG_IP_PORT")

var retryInterval = os.Getenv("KEEP_ALIVE_INTERVAL")
var checkInterval = 5

// Websocket 호출 Test form
func HelloForm(c echo.Context) error {
	fmt.Println("============== Websocket HelloForm ===============")

	return echotemplate.Render(c, http.StatusOK,
		"WebsocketTest", // 파일명
		map[string]interface{}{})
}

// Gorilla WebSocket 호출 Test
var (
	upgrader = websocket.Upgrader{}
)

//slice
// sort.Slice(ss, func(i, j int) bool {
// 	return ss[i].Value > ss[j].Value
// })
// sort.Ints
// sort.Float64
// sort.Strings

// t := time.Now()
// d1 := t.Add(time.Hour* 4)
// d2 := t.Add(time.Hour* -4)

// Websocket 호출 및 Set sample
// 검토할 내용. 여러브라우저에서 호출 후 페이지 이동시 해당 소켓 닫히는지. 비활성화 소켓 닫는 방법.
// 특정 시점 이후의 Data만 가져올 수 있는지
// 특정 시퀀스 이후 Data만 가져올 수 있는지
// 특정 시간 이전 Data는 삭제처리
// map의 key를 현재시간의 unixtime = 숫자로 하면 가능할 것 같은데...
func HelloGorillaWebSocket(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Print("upgrade:", err)
	}
	defer ws.Close()

	t := time.Now()
	ws.SetWriteDeadline(t.Add(time.Second * 1200))
	for {
		//messageType, message, err := ws.ReadMessage()
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		var objmap map[string]interface{}
		_ = json.Unmarshal(message, &objmap)
		event := objmap["event"].(string) // event에 구분 : 요청은 open, req, close.   응답은 res로 한다.

		// 받은 Event 값에서 calltime 가지고 조회 목록에서 해당시간 이후만 가져오도록 한다.
		// callTime이 없는 경우(open) 이면 현재시간 - 2시간(default) 이후의 값을 가져오도록 한다.
		// sendData := map[string]interface{}{
		// 	"event": "res",
		// 	"data":  nil,
		// }
		switch event {
		case "open": // 화면이 처음 열렸을 때
			log.Printf("Received: %s\n", event)

			go testData(c) // test data 처리
			log.Println("is socket working : open started")
			defaultTime := t.Add(time.Minute * -5) // 기본 조회시간은 현재시간 - 5분

			socketDataMap := service.GetWebsocketMessageByProcessTime(defaultTime.UnixNano(), c)
			log.Printf("Received: %s\n", event)
			returnMessage := map[string]interface{}{
				"event":    "res",
				"messag":   socketDataMap,
				"callTime": time.Now().UnixNano(),
			}

			sendErr := ws.WriteJSON(returnMessage)
			if sendErr == nil {
				// service.SetWebsocketMessageBySend(key, true, c)	// 전송 성공여부에 대해 굳이 update가 필요한가??
			}
			log.Println("is socket working : open finished")
		case "req": // 특정시간 이후 모두 조회. 조회할 시간이 parameter로 넘어온다. // key값이 unixTime으로 되어 있으므로  string -> int64 -> unixTime -> time
			// sendData["data"] = objmap["data"]
			sCallTime := objmap["callTime"].(string)
			log.Println("is socket working : req started")
			nCallTime, nErr := strconv.ParseInt(sCallTime, 10, 64)
			if nErr != nil {
				d2 := t.Add(time.Minute * -5)
				nCallTime = d2.UnixNano()
			}

			uCallTime := time.Unix(nCallTime, 0)
			socketDataMap := service.GetWebsocketMessageByProcessTime(uCallTime.UnixNano(), c)
			log.Printf("Received: %s\n", event)
			returnMessage := map[string]interface{}{
				"event":    "res",
				"messag":   socketDataMap,
				"callTime": time.Now().UnixNano(),
			}

			sendErr := ws.WriteJSON(returnMessage)
			if sendErr == nil {
				// service.SetWebsocketMessageBySend(key, true, c)	// 전송 성공여부에 대해 굳이 update가 필요한가??
			}
			log.Println("is socket working : req finished")
			// 마지막 조회
		case "close": // page 이동시
			ws.Close()
			log.Printf("closed")
		}

		// refineSendData, err := json.Marshal(sendData)
		// err = ws.WriteMessage(mt, refineSendData)
		// if err != nil {
		// 	log.Println("write:", err)
		// 	break
		// }

		//// Write
		// log.Println("is socket working : start")

		// 안보낸 것을 보낼 때... 여부가 따로 필요 없음. 시간으로
		// hasSent := false
		// socketDataMap := service.GetWebsocketMessageBySend(hasSent, c)
		// for key, val := range socketDataMap {
		// 	if val.Send == false {
		// 		// if val.Status == "complete" && val.Send == false {
		// 		sendMessage := socketDataMap[key]
		// 		sendMessage.Event = "res"
		// 		sendErr := ws.WriteJSON(sendMessage)
		// 		if sendErr == nil {
		// 			// sendMessage.Send = true
		// 			service.SetWebsocketMessageBySend(key, true, c)

		// 		}
		// 	}
		// }
		// log.Println("is socket working : finish ")
	}
	return err
}

// func HelloGorillaWebSocket(c echo.Context) error {
// 	log.Println("HelloGorillaWebSocket")
// 	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
// 	if err != nil {
// 		return err
// 	}
// 	defer ws.Close()

// 	messageType, p, err := ws.ReadMessage()
// 	if err != nil {
// 		log.Println("ReadMessage")
// 		log.Println(err)
// 		return err
// 	}
// 	// // print out that message for clarity
// 	// fmt.Println(string(p))
// 	//err := conn.WriteMessage(websocket.TextMessage, []byte("Echo push"))
// 	if err := ws.WriteMessage(messageType, p); err != nil {
// 		log.Println("WriteMessage err")
// 		log.Println(err)
// 		return err
// 	}
// 	time.Sleep(time.Second * 10)

// 	// for문 내에서
// 	// 최초요청일 때는 특정 시간 이후 모두 가져오기
// 	// 이후로는 send=false인 것들만 가져오서 전송?

// 	// testCount := 0
// 	// taskKey := "testns" + "||" + "mcis" + "||" + "testmcis"
// 	// // ws.SetReadDeadline(time.Now().Add(30))
// 	// for {

// 	// 	// 사용 예제.
// 	// 	testCount++

// 	// 	// request
// 	// 	if testCount == 5 {
// 	// 		service.SetWebsocketMessage(util.TASKTYPE_MCIS, taskKey, "create", "request", c)
// 	// 	}

// 	// 	// ing
// 	// 	if testCount == 10 {
// 	// 		service.SetWebsocketMessage(util.TASKTYPE_MCIS, taskKey, "create", "ing", c)
// 	// 	}

// 	// 	// complete
// 	// 	if testCount == 15 {
// 	// 		service.SetWebsocketMessage(util.TASKTYPE_MCIS, taskKey, "create", "complete", c)
// 	// 	}

// 	// 	log.Println(" start to read")
// 	// 	//// Read
// 	// 	_, readmsg, err := ws.ReadMessage()
// 	// 	if err != nil {
// 	// 		log.Println("ReadMessage err ", err)
// 	// 		c.Logger().Error(err)
// 	// 		break
// 	// 	}
// 	// 	if string(readmsg) == "close" {
// 	// 		ws.Close()
// 	// 		log.Println("ws closeed ", err)
// 	// 		return err
// 	// 	}

// 	// 	//// Write
// 	// 	log.Println("is socket working : start")

// 	// 	hasSent := false
// 	// 	socketDataMap := service.GetWebsocketMessageBySend(hasSent, c)
// 	// 	for key, val := range socketDataMap {
// 	// 		if val.Send == false {
// 	// 			// if val.Status == "complete" && val.Send == false {
// 	// 			sendMessage := socketDataMap[key]
// 	// 			sendErr := ws.WriteJSON(sendMessage)
// 	// 			if sendErr == nil {
// 	// 				// sendMessage.Send = true
// 	// 				service.SetWebsocketMessageBySend(key, true, c)

// 	// 			}
// 	// 		}
// 	// 	}
// 	// 	log.Println("is socket working : finish ", testCount)

// 	// 	time.Sleep(time.Second * 10)
// 	// }
// 	return err
// }

// 사용 예제.
func testData(c echo.Context) {
	testCount := 0
	//taskKey := "testns" + "||" + "mcis" + "||" + "testmcis"
	t := time.Now()
	//RFC3339     = "2006-01-02T15:04:05Z07:00"

	taskKey := t.Format(time.RFC3339) + "||" + "testns" + "||" + "mcis" + "||" + "testmcis"
	for {
		testCount++

		// request
		if testCount == 5 {
			service.SetWebsocketMessage(util.TASK_TYPE_MCIS, taskKey, util.MCIS_LIFECYCLE_CREATE, util.TASK_STATUS_REQUEST, c)
		}

		// ing
		if testCount == 10 {
			service.SetWebsocketMessage(util.TASK_TYPE_MCIS, taskKey, "create", "ing", c)
		}

		// complete
		if testCount == 15 {
			service.SetWebsocketMessage(util.TASK_TYPE_MCIS, taskKey, "create", "complete", c)
			break
		}

		time.Sleep(time.Second * 10)
	}

}

// Websocket 한번만. 호출되면 수행 수 닫기 : client에서 여러번 호출하게... 이러면 socket의 의미가 있나?
//func WebSocketOneShot(c echo.Context) error {
//	log.Println("WebSocketOneShot")
//	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
//	if err != nil {
//		return err
//	}
//	defer ws.Close()
//
//	// go testData(c) // 임시로 data 넣기
//
//	hasSent := false
//	socketDataMap := service.GetWebsocketMessageBySend(hasSent, c)
//	for key, val := range socketDataMap {
//		if val.Send == false { // 이건 좀 불합리 할 것 같은데.... 특정 시간 이내의 data만 전송하게 ???
//			// if val.Status == "complete" && val.Send == false {
//			sendMessage := socketDataMap[key]
//			sendErr := ws.WriteJSON(sendMessage)
//			if sendErr == nil {
//				// service.SetWebsocketMessageBySend(key, true, c) // 전송여부 set이 필요없을 듯.
//			}
//		}
//	}
//
//	ws.Close()
//	log.Println("is socket working : finish ")
//	return err
//}

// Listener 에서 감지된 Data 변경을 UI 로 push : listener 구현이 어떻게 될지 모르므로 일단은 남겨 놓음.
// func GorillaWebSocketPush(c echo.Context) error {
// 	err := c.WriteMessage(ws.TextMessage, {message})
// }

type pushMessage struct {
	pushpush string
}

func Echo(conn *websocket.Conn) {
	for {
		//m := pushMessage{pushpush: "echo"}

		// err := conn.ReadJSON(&m)
		// if err != nil {
		//     fmt.Println("Error reading json.", err)
		// }

		// fmt.Printf("Got message: %#v\n", m)

		// if err = conn.WriteJSON(m); err != nil {
		// 	fmt.Println(err)
		// }

		// m := msg{}
		// err := conn.ReadJSON(&m)

		err := conn.WriteMessage(websocket.TextMessage, []byte("Echo push"))
		if err != nil {
			// c.Logger().Error(err)
			// fmt.Printf(err)
			log.Println(err)
		}
	}
}

// WebSocket 통신
// 받은 Event 값에서 calltime 가지고 조회 목록에서 해당시간 이후만 가져오도록 한다.
// callTime이 없는 경우(open) 이면 현재시간 - 2시간(default) 이후의 값을 가져오도록 한다.
func GetWebSocketData(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Print("upgrade:", err)
	}
	defer ws.Close()

	t := time.Now()
	ws.SetWriteDeadline(t.Add(time.Second * 1200))

	for {
		//messageType, message, err := ws.ReadMessage()
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		var objmap map[string]interface{}
		_ = json.Unmarshal(message, &objmap)
		event := objmap["event"].(string) // event에 구분 : 요청은 open, req, close.   응답은 res로 한다.
		log.Printf("Received: %s\n", event)
		switch event {
		case "open": // 화면이 처음 열렸을 때
			log.Println("is socket working : open started")
			defaultTime := t.Add(time.Minute * -60) // 기본 조회시간은 현재시간 - 60분

			socketDataList := service.GetWebsocketMessageByProcessTime(defaultTime.UnixNano(), c)
			log.Println(socketDataList)
			returnMessage := map[string]interface{}{
				"event":    "res",
				"message":  socketDataList,
				"callTime": time.Now().UnixNano(),
			}

			sendErr := ws.WriteJSON(returnMessage)
			if sendErr != nil {
				log.Println("ws send Err ", sendErr.Error())
			}
			log.Println("is socket working : open finished")
		case "req": // 특정시간 이후 모두 조회. 조회할 시간이 parameter로 넘어온다. // key값이 unixTime으로 되어 있으므로  string -> int64 -> unixTime -> time
			// sendData["data"] = objmap["data"]
			log.Println(objmap)
			sCallTime := objmap["callTime"].(string)
			log.Println("is socket working : req started")
			nCallTime, nErr := strconv.ParseInt(sCallTime, 10, 64)
			if nErr != nil {
				log.Println("sCallTime err  ", sCallTime)
				d2 := t.Add(time.Minute * -5)
				nCallTime = d2.UnixNano()
			}

			uCallTime := time.Unix(0, nCallTime)
			//socketDataMap := service.GetWebsocketMessageByProcessTime(d2.UnixNano(), c)
			socketDataMap := service.GetWebsocketMessageByProcessTime(uCallTime.UnixNano(), c)
			returnMessage := map[string]interface{}{
				"event":    "res",
				"message":  socketDataMap,
				"callTime": time.Now().UnixNano(),
			}

			sendErr := ws.WriteJSON(returnMessage)
			if sendErr != nil {
				log.Println("ws send Err ", sendErr.Error())
			}
			log.Println("is socket working : req finished")
			// 마지막 조회
		case "close": // page 이동시
			ws.Close()
			log.Printf("socket closed ")
		}

	}
	return err
}

// TODO :
// 1. front-end에서 상태값 요청
// 알림 icon이 있는 header(?)에서 open일 때 socket연결하고 가져온 Data가 있으면 lastProcessTime을 set, inter val 마다  lastProcessTime을 넘겨서 req 로 값을 가져옴. 가져온 게 있으면 알림에 badge로 표시. 클릭시 가져온 목록 표시.
// 2. back-end 에서 상태값 return
// 3. back-end 에서 처리 완료된 상태값을 어떻게 push 할 것인가.... -> 고민하자.
//		방안 1. 열린 화면에서 Open 시 가져오고, 상태값이 있다면 주어진 Interval 마다 변경된 것이 있는지 요청.   back-end 에서는 변경사항이 반영되면 session에 save 해놓으면 소켓에서 요청올 때마다 꺼내가게 됨.
//
