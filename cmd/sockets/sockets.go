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
		// TODO 추후 적용
		// origin := string(ctx.Request.Header.Peek("Origin"))

		// // 허용할 Origin 목록 (DNS & IP)
		// allowedOrigins := map[string]bool{
		// 	// "https://example.com": true,
		// 	"http://192.168.1.100": true,  // 특정 IP 허용
		// 	"http://127.0.0.1": true,       // 로컬 테스트 허용
		// }

		// if allowedOrigins[origin] {
		// 	return true
		// }

		// config.Stdlog.Println("Blocked WebSocket connection from disallowed origin:", origin)

		// return false

		// DESCRIPTION 모든 경로 허용
		return true
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

			// DESCRIPTION 연결이 끊기거나, 데이터 수신 중 문제가 발생
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					// DESCRIPTION 비정상 종료 (예상치 못한 종료)
					config.Stdlog.Println("WebSocket error : ", err, " / from : ", clientIP)
				} else {
					// DESCRIPTION 정상 종료 (CloseNormalClosure, CloseGoingAway)
					config.Stdlog.Println("Client disconnected from : ", clientIP)
				}
				break
			}

			if msgType == websocket.TextMessage {
				var resp []byte
				resp = krr.ReqReceiveSocket(msg, clientIP)
				// config.Stdlog.Println("Message : ", string(resp))
				ws.WriteMessage(msgType, resp)
			} else if msgType == websocket.BinaryMessage {
				
			}

			// config.Stdlog.Println("Received JSON:", string(msg))
		}
	})

	if err != nil {
		config.Stdlog.Println("socket upgrader error :", err)
	}
}