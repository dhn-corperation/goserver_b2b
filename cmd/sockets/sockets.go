package sockets

import (

	config "mycs/configs"
	krr "mycs/cmd/kakao/kaoreqreceive"
	// cm "mycs/internal/commons"

	"github.com/valyala/fasthttp"
	"github.com/fasthttp/websocket"
)

var upgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize:    1024 * 1024, // 1MB
	WriteBufferSize:   1024 * 1024, // 1MB
	EnableCompression: true, // 데이터 압축
	CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
		// 모든 요청을 허용하거나 특정 도메인만 허용 가능
		return true // 보안이 필요하면 도메인 검사를 추가
	},
}

func ServeWs(ctx *fasthttp.RequestCtx) {
	clientIP := ctx.RemoteIP().String()
	config.Stdlog.Println("New WebSocket connection from : ", clientIP)

	err := upgrader.Upgrade(ctx, func(ws *websocket.Conn) {
		defer func() {
			config.Stdlog.Println("WebSocket closed from : ", clientIP)
			ws.Close()
		}()
		for {
			msgType, msg, err := ws.ReadMessage()

			// 연결이 끊기거나, 데이터 수신 중 문제가 발생
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					// 비정상 종료 (예상치 못한 종료)
					config.Stdlog.Println("WebSocket error : ", err, " / from : ", clientIP)
				} else {
					// 정상 종료 (CloseNormalClosure, CloseGoingAway)
					config.Stdlog.Println("Client disconnected from : ", clientIP)
				}
				break
			}

			if msgType == websocket.TextMessage {
				var resp []byte
				resp = krr.ReqReceiveSocket(msg)
				config.Stdlog.Println("Message : ", string(resp))
				ws.WriteMessage(msgType, resp)
			} else if msgType == websocket.BinaryMessage {

			}

			config.Stdlog.Println("Received JSON:", string(msg))
		}
	})

	if err != nil {
		config.Stdlog.Println("socket upgrader error :", err)
	}
}