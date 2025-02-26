package main

import (
	"os"
	"fmt"
	"sort"
	"syscall"
	"context"
	"os/signal"
	"database/sql"
	s "strings"

	"mycs/cmd/kaosendrequest"
	"mycs/cmd/oshotproc"
	"mycs/cmd/lguproc"
	"mycs/cmd/otpatproc"
	"mycs/cmd/otplguproc"
	"mycs/cmd/otposhotproc"
	config "mycs/cmd/kaoconfig"
	databasepool "mycs/cmd/kaodatabasepool"

	"github.com/gin-gonic/gin"
	"github.com/takama/daemon"
	_ "github.com/go-sql-driver/mysql"
)

const (
	name        = "DHNServer"
	description = "대형네트웍스 카카오 발송 서버"
)

var dependencies = []string{name+".service"}

var resultTable string

type Service struct {
	daemon.Daemon
}

func (service *Service) Manage() (string, error) {

	usage := "Usage: "+name+" install | remove | start | stop | status"

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

	config.InitConfig()

	databasepool.InitDatabase()

	var rLimit syscall.Rlimit

	rLimit.Max = 50000
	rLimit.Cur = 50000

	err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)

	if err != nil {
		config.Stdlog.Println("Error Setting Rlimit ", err)
	}

	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)

	if err != nil {
		config.Stdlog.Println("Error Getting Rlimit ", err)
	}

	config.Stdlog.Println("Rlimit Final", rLimit)

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
	config.Stdlog.Println(name+" 시작")

	//모든 서비스
	allService := map[string]string{}
	allCtxC := map[string]interface{}{}

	//알림톡
	atctx, cancel := context.WithCancel(context.Background())
	go kaosendrequest.AlimtalkProc(atctx)
	allCtxC["AL"] = cancel
	allService["AL"] = "AL"

	//알림톡 재발송
	atrsctx, cancel := context.WithCancel(context.Background())
	go kaosendrequest.AlimtalkResendProc(atrsctx)
	allCtxC["ALRS"] = cancel
	allService["ALRS"] = "ALRS"

	//친구톡
	frctx, cancel := context.WithCancel(context.Background())
	go kaosendrequest.FriendtalkProc(frctx)
	allCtxC["FR"] = cancel
	allService["FR"] = "FR"

	//친구톡 재발송
	ftrsctx, cancel := context.WithCancel(context.Background())
	go kaosendrequest.FriendtalkResendProc(ftrsctx)
	allCtxC["FRRS"] = cancel
	allService["FRRS"] = "FRRS"

	if s.EqualFold(config.Conf.RESPONSE_METHOD, "polling") {

		ppctx, ppcancel := context.WithCancel(context.Background())

		go kaosendrequest.PollingProc(ppctx)

		allCtxC["pp"] = ppcancel
		allService["pp"] = "PollingProc"

		rpctx, rpcancel := context.WithCancel(context.Background())

		go kaosendrequest.ResultProc(rpctx)

		allCtxC["rp"] = rpcancel
		allService["rp"] = "PollingRP"

	}

	noshotUser := map[string]string{}
	noshotCtxC := map[string]interface{}{}

	noshotUserList, error := databasepool.DB.Query("select distinct user_id from DHN_CLIENT_LIST dcl where dcl.use_flag = 'Y' and upper(ifnull(dcl.dest, '')) = 'NOSHOT'")
	isNOshot := true
	if error != nil {
		config.Stdlog.Println("오샷 유저 select 오류 ")
		isNOshot = false
	}
	defer noshotUserList.Close()

	if isNOshot {
		var user_id sql.NullString

		for noshotUserList.Next() {
			noshotUserList.Scan(&user_id)
			ctx, cancel := context.WithCancel(context.Background())
			ctx = context.WithValue(ctx, "user_id", user_id.String)
			go oshotproc.NOshotProcess(user_id.String, ctx)

			noshotUser[user_id.String] = user_id.String
			noshotCtxC[user_id.String] = cancel

			allCtxC["NOS"+user_id.String] = cancel
			allService["NOS"+user_id.String] = user_id.String

		}
	}

	nolctx, nolcancel := context.WithCancel(context.Background())
	go oshotproc.NMSGProcess(nolctx)
	allCtxC["noshotmsg"] = nolcancel
	allService["noshotmsg"] = "NOshot MSG"

	oshotUser := map[string]string{}
	oshotCtxC := map[string]interface{}{}

	oshotUserList, error := databasepool.DB.Query("select distinct user_id from DHN_CLIENT_LIST dcl where dcl.use_flag = 'Y' and upper(ifnull(dcl.dest, '')) = 'OSHOT'")
	isOshot := true
	if error != nil {
		config.Stdlog.Println("오샷 유저 select 오류 ")
		isOshot = false
	}
	defer oshotUserList.Close()

	if isOshot {
		var user_id sql.NullString

		for oshotUserList.Next() {
			oshotUserList.Scan(&user_id)
			ctx, cancel := context.WithCancel(context.Background())
			ctx = context.WithValue(ctx, "user_id", user_id.String)
			go oshotproc.OshotProcess(user_id.String, ctx)

			oshotUser[user_id.String] = user_id.String
			oshotCtxC[user_id.String] = cancel

			allCtxC["OS"+user_id.String] = cancel
			allService["OS"+user_id.String] = user_id.String

		}
	}

	olctx, olcancel := context.WithCancel(context.Background())
	go oshotproc.LMSProcess(olctx)
	allCtxC["oshotlms"] = olcancel
	allService["oshotlms"] = "Oshot LMS"

	osctx, oscancel := context.WithCancel(context.Background())
	go oshotproc.SMSProcess(osctx)
	allCtxC["oshotsms"] = oscancel
	allService["oshotsms"] = "Oshot SMS"

	lguUser := map[string]string{}
	lguCtxC := map[string]interface{}{}

	lguUserList, error := databasepool.DB.Query("select distinct user_id from DHN_CLIENT_LIST dcl where dcl.use_flag = 'Y' and upper(ifnull(dcl.dest, '')) = 'LGU'")
	isLgu := true
	if error != nil {
		config.Stdlog.Println("Lgu 유저 select 오류 ")
		isLgu = false
	}
	defer lguUserList.Close()

	if isLgu {
		var user_id sql.NullString

		for lguUserList.Next() {
			lguUserList.Scan(&user_id)
			ctx, cancel := context.WithCancel(context.Background())
			ctx = context.WithValue(ctx, "user_id", user_id.String)
			go lguproc.LguProcess(user_id.String, ctx)

			lguUser[user_id.String] = user_id.String
			lguCtxC[user_id.String] = cancel

			allCtxC["LG"+user_id.String] = cancel
			allService["LG"+user_id.String] = user_id.String

		}
	}

	llctx, llcancel := context.WithCancel(context.Background())
	go lguproc.LMSProcess(llctx)
	allCtxC["lgulms"] = llcancel
	allService["lgulms"] = "LGU LMS"

	lsctx, lscancel := context.WithCancel(context.Background())
	go lguproc.SMSProcess(lsctx)
	allCtxC["lgusms"] = lscancel
	allService["lgusms"] = "LGU SMS"

	//OTP 영역 시작

	otpAtUser := map[string]string{}
	otpAtCtxC := map[string]interface{}{} 
	otpLguUser := map[string]string{}
	otpLguCtxC := map[string]interface{}{}
	otpOshotUser := map[string]string{}
	otpOshotCtxC := map[string]interface{}{}

	oatctx, cancel := context.WithCancel(context.Background())
	go otpatproc.AlimtalkProc(oatctx)

	otpAtUser["OTPAT"] = "OTPAT"
	otpAtCtxC["OTPAT"] = cancel

	allCtxC["OTPAT"] = cancel
	allService["OTPAT"] = "OTPAT"

	ctx, cancel := context.WithCancel(context.Background())
	go otposhotproc.OshotProcess(ctx)

	otpOshotUser["OTPOSHOT"] = "OTPOSHOT"
	otpOshotCtxC["OTPOSHOT"] = cancel

	allCtxC["OTPOSHOT"] = cancel
	allService["OTPOSHOT"] = "OTPOSHOT"

	ollctx, ollcancel := context.WithCancel(context.Background())
	go otplguproc.LMSProcess(ollctx)
	allCtxC["otplgulms"] = ollcancel
	allService["otplgulms"] = "LGU OTP LMS"

	olsctx, olscancel := context.WithCancel(context.Background())
	go otplguproc.SMSProcess(olsctx)
	allCtxC["otplgusms"] = olscancel
	allService["otplgusms"] = "LGU OTP SMS"

	oomctx, oomcancel := context.WithCancel(context.Background())
	go otposhotproc.MSGProcess(oomctx)
	allCtxC["otposhotmsg"] = oomcancel
	allService["otposhotmsg"] = "OSHOT OTP MSG"

	//OTP 영역 종료

	//API 영역 시작
	r := gin.New()
	r.Use(gin.Recovery())
	serCmd := `DHN Server API
Command :
/ostop?uid=dhn   	 	-> dhn Oshot process stop.
/orun?uid=dhn    	 	-> dhn Oshot process run.
/olist           	 	-> 실행 중인 Oshot process User List.

/nostop?uid=dhn   	 	-> dhn NOshot process stop.
/norun?uid=dhn    	 	-> dhn NOshot process run.
/nolist           	 	-> 실행 중인 NOshot process User List.

/lgstop?uid=dhn   	 	-> dhn Lgu process stop.
/lgrun?uid=dhn   	 	-> dhn Lgu process run.
/lglist           	 	-> 실행 중인 Lgu process User List.

/otpstop?uid=dhn  	 	-> dhn OTP process stop.
/otpatstop?uid=dhn   	-> dhn 알림톡 OTP process stop.
/otprun?uid=dhn&pf=XX	-> dhn OTP process run.
/otplist          	 	-> 실행 중인 OTP process User List.

/all             	 	-> DHNServer process list
/allstop         	 	-> DHNServer process stop
`
	r.GET("/", func(c *gin.Context) {
		c.String(200, serCmd)
	})

	r.GET("/ostop", func(c *gin.Context) {
		var uid string
		uid = c.Query("uid")
		temp := oshotCtxC[uid]
		if temp != nil {
			cancel := oshotCtxC[uid].(context.CancelFunc)
			cancel()
			delete(oshotCtxC, uid)
			delete(oshotUser, uid)
			delete(allService, "OS"+uid)
			delete(allCtxC, "OS"+uid)
			c.String(200, uid+" 종료 신호 전달 완료")
		} else {
			c.String(200, uid+"는 실행 중이지 않습니다.")
		}

	})

	r.GET("/orun", func(c *gin.Context) {
		var uid string
		uid = c.Query("uid")
		temp := oshotCtxC[uid]
		if temp != nil {
			c.String(200, uid+" 이미 실행 중입니다.")
		} else {
			ctx, cancel := context.WithCancel(context.Background())
			ctx = context.WithValue(ctx, "user_id", uid)
			go oshotproc.OshotProcess(uid, ctx)

			oshotCtxC[uid] = cancel
			oshotUser[uid] = uid

			allCtxC["OS"+uid] = cancel
			allService["OS"+uid] = uid

			_, err := databasepool.DB.Exec("update DHN_CLIENT_LIST set dest = 'OSHOT' where use_flag = 'Y' and user_id = ?", uid)
				
			if err != nil {
				config.Stdlog.Println(uid," /orun 오샷 DHN_CLIENT_LIST 업데이트 실패 : ", err)
			}

			c.String(200, uid+" 시작 신호 전달 완료")
		}
	})

	r.GET("/olist", func(c *gin.Context) {
		var key string
		for k := range oshotUser {
			key = key + k + "\n"
		}
		c.String(200, key)
	})

	r.GET("/nostop", func(c *gin.Context) {
		var uid string
		uid = c.Query("uid")
		temp := oshotCtxC[uid]
		if temp != nil {
			cancel := oshotCtxC[uid].(context.CancelFunc)
			cancel()
			delete(oshotCtxC, uid)
			delete(oshotUser, uid)
			delete(allService, "OS"+uid)
			delete(allCtxC, "OS"+uid)
			c.String(200, uid+" 종료 신호 전달 완료")
		} else {
			c.String(200, uid+"는 실행 중이지 않습니다.")
		}

	})

	r.GET("/norun", func(c *gin.Context) {
		var uid string
		uid = c.Query("uid")
		temp := oshotCtxC[uid]
		if temp != nil {
			c.String(200, uid+" 이미 실행 중입니다.")
		} else {
			ctx, cancel := context.WithCancel(context.Background())
			ctx = context.WithValue(ctx, "user_id", uid)
			go oshotproc.OshotProcess(uid, ctx)

			oshotCtxC[uid] = cancel
			oshotUser[uid] = uid

			allCtxC["OS"+uid] = cancel
			allService["OS"+uid] = uid

			_, err := databasepool.DB.Exec("update DHN_CLIENT_LIST set dest = 'OSHOT' where use_flag = 'Y' and user_id = ?", uid)
				
			if err != nil {
				config.Stdlog.Println(uid," /orun 오샷 DHN_CLIENT_LIST 업데이트 실패 : ", err)
			}

			c.String(200, uid+" 시작 신호 전달 완료")
		}
	})

	r.GET("/nolist", func(c *gin.Context) {
		var key string
		for k := range oshotUser {
			key = key + k + "\n"
		}
		c.String(200, key)
	})

	r.GET("/lgstop", func(c *gin.Context) {
		var uid string
		uid = c.Query("uid")
		temp := lguCtxC[uid]
		if temp != nil {
			cancel := lguCtxC[uid].(context.CancelFunc)
			cancel()
			delete(lguCtxC, uid)
			delete(lguUser, uid)

			delete(allService, "LG"+uid)
			delete(allCtxC, "LG"+uid)

			c.String(200, uid+" 종료 신호 전달 완료")
		} else {
			c.String(200, uid+"는 실행 중이지 않습니다.")
		}

	})

	r.GET("/lgrun", func(c *gin.Context) {
		var uid string
		uid = c.Query("uid")
		temp := lguCtxC[uid]
		if temp != nil {
			c.String(200, uid+" 이미 실행 중입니다.")
		} else {
			ctx, cancel := context.WithCancel(context.Background())
			ctx = context.WithValue(ctx, "user_id", uid)
			go lguproc.LguProcess(uid, ctx)

			lguCtxC[uid] = cancel
			lguUser[uid] = uid

			allCtxC["LG"+uid] = cancel
			allService["LG"+uid] = uid

			_, err := databasepool.DB.Exec("update DHN_CLIENT_LIST set dest = 'LGU' where use_flag = 'Y' and user_id = ?", uid)
				
			if err != nil {
				config.Stdlog.Println(uid," /orun 오샷 DHN_CLIENT_LIST 업데이트 실패 : ", err)
			}

			c.String(200, uid+" 시작 신호 전달 완료")
		}
	})

	r.GET("/lglist", func(c *gin.Context) {
		var key string
		for k := range lguUser {
			key = key + k + "\n"
		}
		c.String(200, key)
	})

	otpstop := func() {
		temp := otpLguCtxC["OTPLGU"]
		if temp != nil {
			cancel := otpLguCtxC["OTPLGU"].(context.CancelFunc)
			cancel()
			delete(otpLguCtxC, "OTPLGU")
			delete(otpLguUser, "OTPLGU")

			delete(allService, "OTPLGU")
			delete(allCtxC, "OTPLGU")
		}
		temp = otpOshotCtxC["OTPOSHOT"]
		if temp != nil {
			cancel := otpOshotCtxC["OTPOSHOT"].(context.CancelFunc)
			cancel()
			delete(otpOshotCtxC, "OTPOSHOT")
			delete(otpOshotUser, "OTPOSHOT")

			delete(allService, "OTPOSHOT")
			delete(allCtxC, "OTPOSHOT")
		}


	}

	r.GET("/otpstop", func(c *gin.Context) {
		uid := c.Query("uid")
		if uid == "dhn" {
			otpstop()
			c.String(200, "OTP 종료 신호 전달 완료")
		} else {
			c.String(200, "OTP 종료 불가")
		}
		
	})

	r.GET("/otpatstop", func(c *gin.Context) {
		uid := c.Query("uid")
		if uid == "dhn" {
			cancel := otpAtCtxC["OTPAT"].(context.CancelFunc)
			cancel()
			delete(otpAtCtxC, "OTPAT")
			delete(otpAtUser, "OTPAT")

			delete(allService, "OTPAT")
			delete(allCtxC, "OTPAT")
			c.String(200, "OTP 종료 신호 전달 완료")
		} else {
			c.String(200, "OTP 종료 불가")
		}
		
	})

	r.GET("/otprun", func(c *gin.Context) {
		uid := c.Query("uid")
		pf := c.Query("pf")

		if uid == "dhn" && (pf == "lgu" || pf == "nano") {
			otpstop()
			if pf == "lgu" {
				ctx, cancel := context.WithCancel(context.Background())
				go otplguproc.LguProcess(ctx)

				otpLguCtxC["OTPLGU"] = cancel
				otpLguUser["OTPLGU"] = "OTPLGU"

				allCtxC["OTPLGU"] = cancel
				allService["OTPLGU"] = "OTPLGU"

				c.String(200, "OTP LGU 시작 신호 전달 완료")
			} else if pf == "nano" {
				ctx, cancel := context.WithCancel(context.Background())
				go otposhotproc.OshotProcess(ctx)

				otpOshotUser["OTPOSHOT"] = "OTPOSHOT"
				otpOshotCtxC["OTPOSHOT"] = cancel

				allCtxC["OTPOSHOT"] = cancel
				allService["OTPOSHOT"] = "OTPOSHOT"

				c.String(200, "OTP OSHOT 시작 신호 전달 완료")
			}
			
		} else if uid == "dhn" && pf == "at" {
			ctx, cancel := context.WithCancel(context.Background())
			go otpatproc.AlimtalkProc(ctx)

			otpAtUser["OTPAT"] = "OTPAT"
			otpAtCtxC["OTPAT"] = cancel

			allCtxC["OTPAT"] = cancel
			allService["OTPAT"] = "OTPAT"

			c.String(200, "OTP 알림톡 시작 신호 전달 완료")
		} else {
			c.String(200, "OTP 시작 불가")
		}
	})

	r.GET("/otplist", func(c *gin.Context) {
		var key string
		for k := range otpLguUser {
			key = key + k + "\n"
		}

		for l := range otpOshotUser {
			key = key + l + "\n"
		}
		c.String(200, key)
	})

	r.GET("/all", func(c *gin.Context) {
		var key string
		skeys := make([]string, 0, len(allService))
		for k := range allService {
			skeys = append(skeys, k)
		}
		sort.Strings(skeys)
		for _, k := range skeys {
			key = key + k + "\n"
		}
		c.String(200, key)
	})

	r.GET("/allstop", func(c *gin.Context) {
		var key string

		for k := range allService {
			cancel := allCtxC[k].(context.CancelFunc)
			cancel()

			delete(allCtxC, k)
			delete(allService, k)

		}

		c.String(200, key)
	})

	r.Run(":" + config.Conf.SERVER_PORT)

	//API 영역 종료
}
