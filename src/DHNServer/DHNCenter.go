package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"crypto/tls"
	"net/http"
	// "time"
	"runtime/debug"
	"time"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	config "mycs/src/kaoconfig"
	databasepool "mycs/src/kaodatabasepool"
	"mycs/src/kaoreqreceive"

	"mycs/src/kaocenter"

	"github.com/gin-gonic/gin"
	"github.com/takama/daemon"
	"github.com/valyala/fasthttp"
)

const (
	name        = "DHNCenter"
	description = "대형네트웍스 카카오 Center API"
)

var dependencies = []string{"DHNCenter.service"}

var resultTable string

type Service struct {
	daemon.Daemon
}

func (service *Service) Manage() (string, error) {

	usage := "Usage: DHNCenter install | remove | start | stop | status"

	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}
	resultProc()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	for {
		select {
		case killSignal := <-interrupt:
			config.Stdlog.Println("Got signal:", killSignal)
			config.Stdlog.Println("Stoping DB Conntion : ", databasepool.DB.Stats())
			defer databasepool.DB.Close()
			if killSignal == os.Interrupt {
				return "Daemon was interrupted by system signal", nil
			}
			return "Daemon was killed", nil
		}
	}
}

func loadTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return &tls.Config{Certificates: []tls.Certificate{cert}}, nil
}

func main() {

	config.InitCenterConfig()

	databasepool.InitDatabase()

	srv, err := daemon.New(name, description, daemon.SystemDaemon, dependencies...)
	if err != nil {
		config.Stdlog.Println("Error: ", err)
		os.Exit(1)
	}
	service := &Service{srv}
	status, err := service.Manage()
	if err != nil {
		config.Stdlog.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}

func resultProc() {
	config.Stdlog.Println("DHN Center API 시작")

	// 알림톡 / 친구톡 동시 호출 시 http client 호출 오류가 발생하여
	// AlimtalkProc 에서 순차적으로 알림톡 / 친구톡 호출 하도록 수정 함.
	//go kaosendrequest.AlimtalkProc()

	//go kaosendrequest.FriendtalkProc()

	//go kaosendrequest.PollingProc()
	go kaoreqreceive.TempCopyProc()

	// 	// 친구톡 와이드 아이템 리스트 이미지 업로드 요청
	// 	// POST /v1/{partner_key}/image/friendtalk/wideItemList
	// 	r.POST("/ft/wideItemList", kaocenter.Image_wideItemList)

	// 	// 위에꺼랑 코드가 동일함
	// 	r.POST("/ft/carousel", kaocenter.Image_carousel)

	// 	// =============================추가=====================================
	// 	// 알림톡 API
	// 	// 결과 응답 아이디로 조회 (polling)
	// 	// POST /v3/{partner_key}/response/{response_id}
	// 	r.POST("/al/response/:respid", kaocenter.Get_Polling_Id)

	// 	// 업로드 API
	// 	// 알림톡 하이라이트 이미지 업로드 요청
	// 	// POST /v1/{partner_key}/image/alimtalk/itemHighlight
	// 	r.POST("/at/image/itemhighlight", kaocenter.AT_Highlight_Image)

	// 	// 친구톡 캐러셀 피드 이미지 업로드 요청
	// 	// POST /v1/{partner_key}/image/friendtalk/carousel
	// 	r.POST("/ft/image/carousel", kaocenter.FT_Carousel_Feed_Image)

	// 	// 친구톡 캐러셀 커머스 이미지 업로드 요청
	// 	// POST /v1/{partner_key}/image/friendtalk/carouselCommerce
	// 	r.POST("/ft/image/carouselcommerce", kaocenter.FT_Carousel_Commerce_Image)

	// 	// 광고성 메시지 (다이렉트) 이미지 업로드 요청
	// 	// POST /v2/{partner_key}/image/default
	// 	r.POST("/dm/image/default", kaocenter.DM_Default_Image)

	// 	// 광고성메시지(다이렉트) 와이드 이미지 업로드 요청
	// 	// POST /v2/{partner_key}/image/wide
	// 	r.POST("/dm/image/wide", kaocenter.DM_Wide_Image)

	// 	// 광고성메시지(다이렉트) 와이드 리스트 첫번째 이미지 업로드 요청
	// 	// POST /v2/{partner_key}/image/wideItemList/first
	// 	r.POST("/dm/image/wideItemList/first", kaocenter.DM_Widelist_First_image)

	// 	// 광고성메시지(다이렉트) 와이드 리스트 이미지 업로드 요청
	// 	// POST /v2/{partner_key}/image/wideItemList
	// 	r.POST("/dm/image/wideItemList", kaocenter.DM_Widelist_Image)

	// 	// 광고성메시지(다이렉트) 캐러셀 피드 이미지 업로드 요청
	// 	// POST /v2/{partner_key}/image/carouselFeed
	// 	r.POST("/dm/image/carouselFeed", kaocenter.DM_Carousel_Feed_Image)

	// 	// 광고성메시지(다이렉트) 캐러셀 커머스 이미지 업로드 요청
	// 	// POST /v2/{partner_key}/image/carouselCommerce
	// 	r.POST("/dm/image/carouselcommerce", kaocenter.DM_Carousel_Commerce_Image)

	// 	// 친구톡 API 별첨
	// 	// 비즈폼 업로드 요청
	// 	// POST /api/v1/{partner_key}/bizform/upload
	// 	r.POST("/bizform/upload", kaocenter.Bizform_upload_)

	// 	// 친구톡 발송 가능 모수 확인 API
	// 	// POST /api/v1/{partner_key}/friendtalk/possible
	// 	r.POST("/ft/possible", kaocenter.Ft_possible_)

	// 	// 센터 API
	// 	// 발신 프로필 조회2 (톡 채널 키로 조회)
	// 	// GET /api/v3/{partner_key}/sender/{talkChannelKey}
	// 	r.GET("/sender/channel/:talkChannelKey", kaocenter.Sender_channel)

	// 	// 최근 변경 발신 프로필 조회
	// 	// GET /api/v3/{partner_key}/sender/last_modified
	// 	r.GET("/sender/last_modified", kaocenter.Sender_modified)

	// 	// 검수요청 (파일첨부)
	// 	// POST /api/v2/{partner_key}/alimtalk/template/request_with_file
	// 	r.POST("/template/request_with_file", kaocenter.Template_request_with_file)

	// 	// 검수 승인 취소
	// 	// POST /api/v2/{partner_key}/alimtalk/template/cancel_approval
	// 	r.POST("/template/cancel_approval", kaocenter.Template_cancel_approval_)

	// 	// 기등록된 템플릿 (타입 : BA, EX) 을 채널추가버튼 및 채널추가안내문구가 포함된 템플릿으로 전환
	// 	// POST /api/v2/{partner_key}/alimtalk/template/convertAddCh
	// 	r.POST("/template/convertAddCh", kaocenter.Template_convertAddCh_)

	// 	// 채널에 발신 프로필 추가
	// 	// POST /api/v2/{partner_key}/channel/sender/add
	// 	r.POST("/channel/sender/add", kaocenter.Channel_sender_add_)

	// 	// 채널에 발신 프로필 삭제
	// 	// POST /api/v2/{partner_key}/channel/sender/remove
	// 	r.POST("/channel/sender/remove", kaocenter.Channel_sender_remove_)

	// 	// 알림톡, 친구톡 발송 일별 통계
	// 	// GET /api/v1/{hubPartner_key}/stat/daily
	// 	r.GET("/stat/daily", kaocenter.Stat_daily)

	// 	// 그룹 태그 생성
	// 	// POST /api/v1/{hubPartner_key}/groupTag/create
	// 	r.POST("/groupTag/create", kaocenter.GroupTag_create)

	// 	// 그룹 태그 조회 (단건)
	// 	// GET /api/v1/{hubPartner_key}/groupTag
	// 	r.GET("/groupTag", kaocenter.GroupTag_)

	// 	// 그룹 태그 전체 목록 조회
	// 	// GET /api/v1/{hubPartner_key}/groupTag/list
	// 	r.GET("/groupTag/list", kaocenter.GroupTag_list)

	// 	// 그룹 태그 수정
	// 	// POST /api/v1/{hubPartner_key}/groupTag/update
	// 	r.POST("/groupTag/update", kaocenter.GroupTag_update)

	// 	// 그룹 태그 삭제
	// 	// POST /api/v1/{hubPartner_key}/groupTag/delete
	// 	r.POST("/groupTag/delete", kaocenter.GroupTag_delete)

	// 	// 광고성메시지(다이렉트) 템플릿 등록
	// 	// POST /api/v3/{partner_key}/direct/template/create
	// 	r.POST("/dm/template/create", kaocenter.Direct_template_create_)

	// 	// 광고성메시지(다이렉트) 템플릿 조회
	// 	// GET /api/v2/{hubPartnerKey}/direct/template/{code}
	// 	r.GET("/dm/template/:code", kaocenter.Direct_template_)

	// 	// 광고성메시지(다이렉트) 템플릿 수정
	// 	// POST /api/v3/{partner_key}/direct/template/update/{code}
	// 	r.POST("/dm/template/update/:code", kaocenter.Direct_template_update_)

	// 	// 광고성메시지(다이렉트) 템플릿 삭제
	// 	// POST /api/v2/{hubPartnerKey}/direct/template/delete/{code}
	// 	r.POST("/dm/template/delete/:code", kaocenter.Direct_template_delete_)

	// 	// 광고성메시지(다이렉트) 발신 프로필 API 발신 채널 전환
	// 	// POST /api/v1/:hubPartnerKey/sender/direct/convert
	// 	r.POST("/sender/dm/convert", kaocenter.Direct_convert_)

	// 	// 광고성메시지(다이렉트) 발신 프로필 API 발신 채널 전환 상태확인
	// 	// GET /api/v1/:hubPartnerKey/sender/direct/convert/result
	// 	r.GET("/sender/dm/convert/result", kaocenter.Direct_convert_result)

	// 	// 발신채널에 연결된 비즈월렛 변경
	// 	// POST /api/v1/:hubPartnerKey/sender/direct/bizWallet/change
	// 	r.POST("/dm/bizWallet/change", kaocenter.Direct_bizWallet_change_)

	// 	// 다이렉트 메시지 API -> 기본형 API
	// 	// 단건 메시지 전송 요청
	// 	r.POST("/dm/send/basic", test)

	// 	// 다건 메시지 전송 요청
	// 	r.POST("/dm/send/basic/batch", test)

	// 	// 발송 결과 확인 - 1
	// 	r.GET("/dm/basic/response/request", test)

	// 	// 발송 결과 요청 - 2
	// 	r.GET("/dm/basic/response/message", test)

	// 	// 다이렉트 메시지 API -> 자유형 API
	// 	// 단건 메시지 전송 요청
	// 	r.POST("/dm/send/freestyle", test)

	// 	// 다건 메시지 전송 요청
	// 	r.POST("/dm/send/freestyle/batch", test)

	// 	// 발송 결과 확인 - 1
	// 	r.GET("/dm/freestyle/response/request", test)

	// 	// 발송 결과 요청 - 2
	// 	r.GET("/dm/freestyle/response/message", test)

	// 	r.GET("/get_crypto", kaocenter.Get_crypto)

	// 	// SSL 사용 시 --- 시작
	// 	// certFile := "etc/letsencrypt/live/dhntest.dhn.kr/fullchain.pem"
	// 	// keyFile := "etc/letsencrypt/live/dhntest.dhn.kr/privkey.pem"

	// 	// tlsConfig, err := loadTLSConfig(certFile, keyFile)
	// 	// if err != nil {
	// 	// 	config.Stdlog.Println("SSL 인증 실패 err : ", err)
	// 	// 	return
	// 	// }

	// 	// server := &http.Server{
	// 	// 	Addr: ":443",
	// 	// 	Handler: r,
	// 	// 	TLSConfig: tlsConfig,
	// 	// }

	// 	// go func() {
	// 	// 	for {
	// 	// 		time.Sleep(24 * time.Hour)
	// 	// 		config.Stdlog.Println("자동 SSL 인증 갱신 시작")
	// 	// 		newTLSConfig, err := loadTLSConfig(certFile, keyFile)
	// 	// 		if err != nil {
	// 	// 			config.Stdlog.Println("자동 SSL 인증 갱신 실패 err : ", err)
	// 	// 			continue
	// 	// 		}
	// 	// 		server.TLSConfig = newTLSConfig
	// 	// 		config.Stdlog.Println("자동 SSL 인증 갱신 성공")
	// 	// 	}
	// 	// }()

	// 	// err = server.ListenAndServeTLS(certFile, keyFile)
	// 	// if err != nil {
	// 	// 	config.Stdlog.Println("서버 실행 실패")
	// 	// }
	// 	// SSL 사용 시 --- 끝

	// 	// SSL 미사용 시 --- 시작
	// 	r.Run(":" + config.Conf.CENTER_PORT)
	// 	// SSL 미사용 시 --- 끝
	// }()

	

	testHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/":
			ctx.SetStatusCode(fasthttp.StatusOK)
			ctx.SetBodyString("Center Server : "+config.Conf.CENTER_SERVER+",   "+"Image Server : "+config.Conf.IMAGE_SERVER+"\n")
		case "/req":
			kaoreqreceive.ReqReceive(ctx)
		case "/result":
			kaoreqreceive.Resultreq(ctx)
		case "/sresult":
			kaoreqreceive.SearchResultReq(ctx)
		case "/sender/token":
			// 카카오톡 채널 인증 토큰 요청
			// GET /api/v1/{partner_key}/sender/token
			kaocenter.Sender_token(ctx)
		case "/category/all":
			// 발신프로필 카테고리 전체 조회
			// GET /api/v1/{partner_key}/category/all
			kaocenter.Category_all(ctx)
		case "/sender/category":
			// 발신프로필 카테고리 조회
			// GET /api/v1/{partner_key}/category
			kaocenter.Category_(ctx)
		case "/sender/create":
			// 발신프로필 등록
			// POST /api/v3/{partner_key}/sender/create
			kaocenter.Sender_Create(ctx)
		case "/sender":
			// 발신프로필 조회1
			// GET /api/v3/{partner_key}/sender
			kaocenter.Sender_(ctx)
		case "/sender/delete":
			// 발신프로필 삭제
			// POST /api/v1/{partner_key}/sender/delete
			kaocenter.Sender_Delete(ctx)
		case "/sender/recover":
			// 미사용 프로필 휴면 해제
			// POST /api/v1/{partner_key}/sender/recover
			kaocenter.Sender_Recover(ctx)
		case "/template/create":
			// 템플릿 등록
			// POST /api/v2/{partner_key}/alimtalk/template/create
			kaocenter.Template_Create(ctx)
		case "/template":
			// 템플릿 조회
			// GET /api/v2/{partner_key}/alimtalk/template
			kaocenter.Template_(ctx)
		case "/template/request":
			// 검수 요청
			// POST /api/v2/{partner_key}/alimtalk/template/request
			kaocenter.Template_Request(ctx)
		case "/template/cancel_request":
			// 검수 요청 취소
			// POST /api/v2/{partner_key}/alimtalk/template/cancel_request
			kaocenter.Template_Cancel_Request(ctx)
		case "/template/update":
			// 템플릿 수정
			// POST /api/v2/{partner_key}/alimtalk/template/update
			kaocenter.Template_Update(ctx)
		case "/template/stop":
			// 템플릿 사용 중지
			// POST /api/v2/{partner_key}/alimtalk/template/stop
			kaocenter.Template_Stop(ctx)
		case "/template/reuse":
			// 템플릿 사용 중지 해제
			// POST /api/v2/{partner_key}/alimtalk/template/reuse
			kaocenter.Template_Reuse(ctx)
		case "/template/delete":
			// 템플릿 삭제
			// POST /api/v2/{partner_key}/alimtalk/template/delete
			kaocenter.Template_Delete(ctx)
		case "/template/last_modified":
			// 최근 변경 템플릿 조회
			// GET /api/v3/{partner_key}/alimtalk/template/last_modified
			kaocenter.Template_Last_Modified(ctx)
		case "/template/category/all":
			// 템플릿 카테고리 전체 조회
			// GET /api/v2/{partner_key}/alimtalk/template/category/all
			kaocenter.Template_Category_all(ctx)
		case "/template/category":
			// 템플릿 카테고리 조회
			// GET /api/v2/{partner_key}/alimtalk/template/category
			kaocenter.Template_Category_(ctx)
		case "/template/dormant/release":
			// 템플릿 휴면 해제
			// POST /api/v2/{partner_key}/alimtalk/template/dormant/release
			kaocenter.Template_Dormant_Release(ctx)
		case "/group":
			// 발신 프로필 그룹 조회
			// GET /api/v1/{partner_key}/group
			kaocenter.Group_(ctx)
		case "/group/sender":
			// 그룹에 포함된 발신 프로필 조회
			// GET /api/v3/{partner_key}/group/sender
			kaocenter.Group_Sender(ctx)
		case "/group/sender/add":
			// 그룹에 발신 프로필 추가
			// POST /api/v1/{partner_key}/group/sender/add
			kaocenter.Group_Sender_Add(ctx)
		case "/group/sender/remove":
			// 그룹에 발신 프로필 삭제
			// POST /api/v1/{partner_key}/group/sender/remove
			kaocenter.Group_Sender_Remove(ctx)
		case "/channel/create":
			// 채널 생성
			// POST /api/v2/{partner_key}/channel/create
			kaocenter.Channel_Create_(ctx)
		case "/channel/all":
			// 전체 채널 조회
			// GET /api/v2/{partner_key}/channel/all
			kaocenter.Channel_all(ctx)
		case "/channel":
			// 채널 상세 조회
			// GET /api/v2/{partner_key}/channel
			kaocenter.Channel_(ctx)
		case "/channel/update":
			// 채널 수정
			// POST /api/v2/{partner_key}/channel/update
			kaocenter.Channel_Update_(ctx)
		case "/channel/senders":
			// 채널에 발신 프로필 할당
			// POST /api/v2/{partner_key}/channel/senders
			kaocenter.Channel_Senders_(ctx)
		case "/channel/delete":
			// 채널 삭제
			// POST /api/v2/{partner_key}/channel/delete
			kaocenter.Channel_Delete_(ctx)
		case "/plugin/callbackUrl/list":
			// 플러그인 콜백 URL 조회
			// GET /api/v1/{partner_key}/plugin/callbackUrl/list
			kaocenter.Plugin_CallbackUrls_List(ctx)
		case "/plugin/callbackUrl/create":
			// 플러그인 콜백 URL 등록
			// POST /api/v1/{partner_key}/plugin/callbackUrl/create
			kaocenter.Plugin_callbackUrl_Create(ctx)
		case "/plugin/callbackUrl/update":
			// 플러그인 콜백 URL 수정
			// POST /api/v1/{partner_key}/plugin/callbackUrl/update
			kaocenter.Plugin_callbackUrl_Update(ctx)
		case "/plugin/callbackUrl/delete":
			// 플러그인 콜백 URL 삭제
			// POST /api/v1/{partner_key}/plugin/callbackUrl/delete
			kaocenter.Plugin_callbackUrl_Delete(ctx)
		case "/ft/image":
			// 친구톡 이미지 업로드 요청
			// POST /v1/{partner_key}/image/friendtalk
			kaocenter.FT_Upload(ctx)
		case "/ft/wide/image":
			// 친구톡 와이드 이미지 업로드 요청
			// POST /v1/{partner_key}/image/friendtalk/wide
			kaocenter.FT_Wide_Upload(ctx)
		case "/at/image":
			// 알림톡 템플릿 등록용 이미지 업로드 요청
			// POST /v1/{partner_key}/image/alimtalk/template
			kaocenter.AT_Image(ctx)
		case "/mms/image":
			kaocenter.MMS_Image(ctx)
		default:
			ctx.Error("Not Support", fasthttp.StatusNotFound)
		}
	}

	if err := fasthttp.ListenAndServe(":3033", testHandler); err != nil {
		config.Stdlog.Println("fasthttp 실행 실패")
	}
}

func customRecovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if r := recover(); r != nil {
                // panic 로그 기록
                config.Stdlog.Println("Recovered from panic : ", r)
                config.Stdlog.Println("Stack trace: ", string(debug.Stack()))
                
                // 500 Internal Server Error 반환
                c.JSON(http.StatusInternalServerError, gin.H{
                    "code": "error",
                    "message": "panic",
                })
                c.Abort() // 미들웨어 체인의 나머지를 중단
            }
        }()
        c.Next() // 다음 미들웨어 또는 핸들러로 넘김
    }
}

func statusDatabase() gin.HandlerFunc {
	return func(c *gin.Context){
		for {
			db, err := sql.Open(config.Conf.DB, config.Conf.DBURL)

			if err != nil {
				config.Stdlog.Println("DB 연결 실패 5초 후 재시도합니다. err : ", err)
				time.Sleep(5 * time.Second) // 5초 후 재시도
				continue
			}

			// Ping으로 연결 상태 확인
			if err := db.Ping(); err != nil {
				config.Stdlog.Println("DB 핑 신호 없음 err : ", err)
				time.Sleep(5 * time.Second) // 5초 후 재시도
				continue
			} else {
				db.SetMaxIdleConns(50)
				db.SetMaxOpenConns(50)
				databasepool.DB = db
				config.Stdlog.Println("DB 재할당 완료")
			}

			c.Next()
		}
	}
}

func test(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "nokay",
	})
}
