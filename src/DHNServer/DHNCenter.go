package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"

	config "kaoconfig"
	databasepool "kaodatabasepool"
	"kaoreqreceive"

	"kaocenter"
	//"kaosendrequest"
	//"strconv"
	//"time"

	"github.com/gin-gonic/gin"
	"github.com/takama/daemon"
)

const (
	name        = "DHNNanoCenter"
	description = "대형네트웍스 카카오 Center API"
)

var dependencies = []string{"DHNNanoCenter.service"}

var resultTable string

type Service struct {
	daemon.Daemon
}

func (service *Service) Manage() (string, error) {

	usage := "Usage: DHNNanoCenter install | remove | start | stop | status"

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
	r := gin.New()
	r.Use(gin.Recovery())
	//r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		//time.Sleep(30 * time.Second)
		c.String(200, "Center Server : "+config.Conf.CENTER_SERVER+",   "+"Image Server : "+config.Conf.IMAGE_SERVER)
	})

	r.POST("/req", kaoreqreceive.ReqReceive)

	r.POST("/result", kaoreqreceive.Resultreq)

	r.POST("/profile/token", kaocenter.Sender_token)

	r.POST("/profile/category/all", kaocenter.Category_all)

	r.POST("/profile/category", kaocenter.Category_)

	r.POST("/profile/create", kaocenter.Sender_Create)

	r.POST("/profile", kaocenter.Sender_)

	r.POST("/profile/use", kaocenter.Sender_Use)

	r.POST("/profile/multi", kaocenter.Sender_Multi)

	r.POST("/sender/delete", kaocenter.Sender_Delete)

	r.POST("/profile/recover", kaocenter.Sender_Recover)

	r.POST("/template/add", kaocenter.Template_Create)

	r.POST("/template/codeCheck", kaocenter.Template_Code_Check)

	r.POST("/template/create_with_image", kaocenter.Template_Create_Image)

	r.POST("/template/list", kaocenter.TemplateList)

	r.POST("/template/detail", kaocenter.Template_)

	r.POST("/template/request", kaocenter.Template_Request)

	r.POST("/template/request_with_file", kaocenter.Template_Request_WF)

	r.POST("/template/cancel_request", kaocenter.Template_Cancel_Request)

	r.POST("/template/update", kaocenter.Template_Update)

	r.POST("/template/update_with_image", kaocenter.Template_Update_Image)

	r.POST("/template/stop", kaocenter.Template_Stop)

	r.POST("/template/reuse", kaocenter.Template_Reuse)

	r.POST("/template/cancel_approval", kaocenter.Template_CA)

	r.POST("/template/delete", kaocenter.Template_Delete)

	r.GET("/template/last_modified", kaocenter.Template_Last_Modified)

	r.POST("/template/comment", kaocenter.Template_Comment)

	r.POST("/template/comment_file", kaocenter.Template_Comment_File)

	r.POST("/template/category/all", kaocenter.Template_Category_all)

	r.POST("/template/category", kaocenter.Template_Category_)

	r.POST("/template/category/update", kaocenter.Template_Category_Update)

	r.POST("/template/release", kaocenter.Template_Dormant_Release)

	r.POST("/template/convertAddCh", kaocenter.Template_Convert_AC)

	r.POST("/group", kaocenter.Group_)

	r.POST("/group/profile", kaocenter.Group_Sender)

	r.POST("/group/profile/add", kaocenter.Group_Sender_Add)

	r.POST("/group/profile/delete", kaocenter.Group_Sender_Remove)

	r.POST("/channel/create", kaocenter.Channel_Create_)

	r.GET("/channel/all", kaocenter.Channel_all)

	r.GET("/channel", kaocenter.Channel_)

	r.POST("/channel/update", kaocenter.Channel_Update_)

	r.POST("/channel/senders", kaocenter.Channel_Senders_)

	r.POST("/channel/delete", kaocenter.Channel_Delete_)

	r.POST("/plugin/callbackUrl/list", kaocenter.Plugin_CallbackUrls_List)

	r.POST("/plugin/callbackUrl/create", kaocenter.Plugin_callbackUrl_Create)

	r.POST("/plugin/callbackUrl/update", kaocenter.Plugin_callbackUrl_Update)

	r.POST("/plugin/callbackUrl/delete", kaocenter.Plugin_callbackUrl_Delete)

	r.POST("/image/upload", kaocenter.FT_Upload)

	r.POST("/ft/wide/image", kaocenter.FT_Wide_Upload)

	r.POST("/image/alimtalk/template", kaocenter.AT_Image)

	r.POST("/image/alimtalk", kaocenter.AL_Image)

	r.POST("/image/alimtalk/itemHighlight", kaocenter.AL_Image_HL)

	r.POST("/mms/image", kaocenter.MMS_Image)

	r.POST("/friendinfo", kaoreqreceive.FriendInforeq)

	r.POST("/image/friendtalk/wideItemList", kaocenter.Image_wideItemList)

	r.POST("/image/friendtalk/carousel", kaocenter.Image_carousel)

	r.POST("/image/friendtalk/carouselCommerce", kaocenter.Image_carousel_C)

	r.Run(":" + config.Conf.PORT)
}
