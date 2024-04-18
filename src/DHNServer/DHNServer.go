package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"

	config "mycs/src/kaoconfig"
	databasepool "mycs/src/kaodatabasepool"

	//"kaoreqreceive"

	//"kaocenter"
	"mycs/src/kaosendrequest"
	"mycs/src/nanoproc"
	"mycs/src/oshotproc"
	"mycs/src/otpproc"

	//"strconv"
	//"time"
	s "strings"

	"github.com/gin-gonic/gin"
	"github.com/takama/daemon"

	"context"
	"sort"
	//"reflect"
)

const (
	name        = "DHNServer"
	description = "대형네트웍스 카카오 발송 서버"
)

var dependencies = []string{"DHNServer.service"}

var resultTable string

type Service struct {
	daemon.Daemon
}

func (service *Service) Manage() (string, error) {

	usage := "Usage: DHNServer install | remove | start | stop | status"

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
	config.Stdlog.Println("DHN Server 시작")

	allService := map[string]string{}
	allCtxC := map[string]interface{}{}

	alim_user_list, error := databasepool.DB.Query("select distinct user_id from DHN_CLIENT_LIST where alimtalk='Y'")
	isAlim := true
	if error != nil {
		config.Stdlog.Println("알림톡 유저 select 오류 ")
		isAlim = false
	}
	defer alim_user_list.Close()

	alimUser := map[string]string{}
	alimCtxC := map[string]interface{}{}

	if isAlim {
		var user_id sql.NullString
		for alim_user_list.Next() {

			alim_user_list.Scan(&user_id)

			ctx, cancel := context.WithCancel(context.Background())
			go kaosendrequest.AlimtalkProc(user_id.String, ctx)

			alimCtxC[user_id.String] = cancel
			alimUser[user_id.String] = user_id.String

			allCtxC["AL"+user_id.String] = cancel
			allService["AL"+user_id.String] = user_id.String

		}
	}

	ftctx, ftcancel := context.WithCancel(context.Background())

	go kaosendrequest.FriendtalkProc(ftctx)

	allCtxC["ft"] = ftcancel
	allService["ft"] = "ft"

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

	oshotUser := map[string]string{}
	oshotCtxC := map[string]interface{}{}

	if s.EqualFold(config.Conf.PHONE_MSG_FLAG, "YES") {
		oshotUserList, error := databasepool.DB.Query("select distinct user_id from DHN_CLIENT_LIST dcl  where dcl.use_flag   = 'Y' and IFNULL(dcl.dest, 'OSHOT') = 'OSHOT'")
		isOshot := true
		if error != nil {
			config.Stdlog.Println("Oshot 유저 select 오류 ")
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

			olctx, olcancel := context.WithCancel(context.Background())
			go oshotproc.LMSProcess(olctx)
			allCtxC["oshotlms"] = olcancel
			allService["oshotlms"] = "Oshot LMS"

			osctx, oscancel := context.WithCancel(context.Background())
			go oshotproc.SMSProcess(osctx)
			allCtxC["oshotsms"] = oscancel
			allService["oshotsms"] = "Oshot SMS"

		}
	}

	nanoUser := map[string]string{}
	nanoCtxC := map[string]interface{}{}

	if s.EqualFold(config.Conf.NANO_MSG_FLAG, "YES") {

		nanoUserList, error := databasepool.DB.Query("select distinct user_id from DHN_CLIENT_LIST dcl  where dcl.use_flag   = 'Y' and IFNULL(dcl.dest, 'OSHOT') = 'NANO'")
		isNano := true
		if error != nil {
			config.Stdlog.Println("Nano 유저 select 오류 ")
			isNano = false
		}
		defer nanoUserList.Close()

		if isNano {
			var user_id sql.NullString

			for nanoUserList.Next() {
				nanoUserList.Scan(&user_id)

				if s.EqualFold(config.Conf.PHONE_TYPE_FLAG, "N") { // 기본
					// ctx, cancel := context.WithCancel(context.Background())
					// ctx = context.WithValue(ctx, "user_id", user_id.String)

					// go nanoproc.NanoProcess(user_id.String, ctx)

					// nanoCtxC[user_id.String] = cancel

					// allCtxC["NN"+user_id.String] = cancel
					// allService["NN"+user_id.String] = user_id.String
				} else if s.EqualFold(config.Conf.PHONE_TYPE_FLAG, "Y") { // 콜비서

					ctxY, cancelY := context.WithCancel(context.Background())
					ctxY = context.WithValue(ctxY, "user_id", user_id.String)

					ctxN, cancelN := context.WithCancel(context.Background())
					ctxN = context.WithValue(ctxN, "user_id", user_id.String)

					go nanoproc.NanoProcess_Y(user_id.String, ctxY) // 010으로 시작하는 번호
					go nanoproc.NanoProcess_N(user_id.String, ctxN) // 010이 아닌 번호

					nanoCtxC[user_id.String+"_Y"] = cancelY
					allCtxC["NN"+user_id.String+"_Y"] = cancelY
					allService["NN"+user_id.String+"_Y"] = "NanoService Y"

					nanoCtxC[user_id.String+"_N"] = cancelN
					allCtxC["NN"+user_id.String+"_N"] = cancelN
					allService["NN"+user_id.String+"_N"] = "NanoService N"
				}

				nanoUser[user_id.String] = user_id.String

			}

			// nlctx, nlcancel := context.WithCancel(context.Background())

			// go nanoproc.NanoLMSProcess(nlctx)

			// allCtxC["nanolms"] = nlcancel
			// allService["nanolms"] = "Nano LMS"

			// nsctx, nscancel := context.WithCancel(context.Background())

			// go nanoproc.NanoSMSProcess(nsctx)

			// allCtxC["nanosms"] = nscancel
			// allService["nanosms"] = "Nano SMS"

			if s.EqualFold(config.Conf.PHONE_TYPE_FLAG, "Y") { // 콜비서

				nlctxG, nlcancelG := context.WithCancel(context.Background())

				go nanoproc.NanoLMSProcess_G(nlctxG)

				allCtxC["nanolmsG"] = nlcancelG
				allService["nanolmsG"] = "Nano LMS G"

				nsctxG, nscancelG := context.WithCancel(context.Background())

				go nanoproc.NanoSMSProcess_G(nsctxG)

				allCtxC["nanosmsG"] = nscancelG
				allService["nanosmsG"] = "Nano SMS G"

			}
		}
	}

	if s.EqualFold(config.Conf.OTP_MSG_FLAG, "YES") {
		otpctx, otpcancel := context.WithCancel(context.Background())
		go otpproc.OTPProcess(otpctx)
		allCtxC["otpproc"] = otpcancel
		allService["otpproc"] = "OTP Proc"

		otplctx, otplcancel := context.WithCancel(context.Background())
		go otpproc.OTPLMSProcess(otplctx)
		allCtxC["otplms"] = otplcancel
		allService["otplms"] = "OTP LMS"

		otpsctx, otpscancel := context.WithCancel(context.Background())
		go otpproc.OTPSMSProcess(otpsctx)
		allCtxC["otpsms"] = otpscancel
		allService["otpsms"] = "OTP SMS"
	}

	r := gin.New()
	r.Use(gin.Recovery())
	//r := gin.Default()
	serCmd := `DHN Server 
Command :
/astop?uid=dhn   -> dhn Alimtalk process stop.
/arun?uid=dhn    -> dhn Alimtalk process run.
/alist           -> 실행 중인 Alimtalk process User List.

/ostop?uid=dhn   -> dhn Oshot process stop.
/orun?uid=dhn    -> dhn Oshot process run.
/olist           -> 실행 중인 Oshot process User List.

/nstop?uid=dhn   -> dhn Nano process stop.
/nrun?uid=dhn    -> dhn Nano process run.
/nlist           -> dhn Nano process run.

/all             -> DHNServer process list
/allstop         -> DHNServer process stop
`
	r.GET("/", func(c *gin.Context) {
		//time.Sleep(30 * time.Second)
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
			c.String(200, uid+"이미 실행 중입니다.")
		} else {
			ctx, cancel := context.WithCancel(context.Background())
			ctx = context.WithValue(ctx, "user_id", uid)
			go oshotproc.OshotProcess(uid, ctx)

			oshotCtxC[uid] = cancel
			oshotUser[uid] = uid

			allCtxC["OS"+uid] = cancel
			allService["OS"+uid] = uid

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

	r.GET("/astop", func(c *gin.Context) {
		var uid string
		uid = c.Query("uid")
		temp := alimCtxC[uid]
		if temp != nil {
			cancel := alimCtxC[uid].(context.CancelFunc)
			cancel()
			delete(alimCtxC, uid)
			delete(alimUser, uid)

			delete(allService, "AL"+uid)
			delete(allCtxC, "AL"+uid)

			c.String(200, uid+" 종료 신호 전달 완료")
		} else {
			c.String(200, uid+"는 실행 중이지 않습니다.")
		}

	})

	r.GET("/arun", func(c *gin.Context) {
		var uid string
		uid = c.Query("uid")
		temp := alimCtxC[uid]
		if temp != nil {
			c.String(200, uid+"이미 실행 중입니다.")
		} else {
			ctx, cancel := context.WithCancel(context.Background())
			ctx = context.WithValue(ctx, "user_id", uid)
			go kaosendrequest.AlimtalkProc(uid, ctx)

			alimCtxC[uid] = cancel
			alimUser[uid] = uid

			allCtxC["AL"+uid] = cancel
			allService["AL"+uid] = uid

			c.String(200, uid+" 시작 신호 전달 완료")
		}
	})

	r.GET("/alist", func(c *gin.Context) {
		var key string
		for k := range alimUser {
			key = key + k + "\n"
		}
		c.String(200, key)
	})

	r.GET("/nstop", func(c *gin.Context) {
		var uid string
		uid = c.Query("uid")
		temp := nanoCtxC[uid]
		if temp != nil {
			cancel := nanoCtxC[uid].(context.CancelFunc)
			cancel()
			delete(nanoUser, uid)
			if s.EqualFold(config.Conf.PHONE_TYPE_FLAG, "N") { // 기본
				delete(nanoCtxC, uid)

				delete(allService, "NN"+uid)
				delete(allCtxC, "NN"+uid)
			} else if s.EqualFold(config.Conf.PHONE_TYPE_FLAG, "Y") { // 콜비서
				delete(nanoCtxC, uid+"_Y")
				delete(nanoCtxC, uid+"_N")

				delete(allService, "NN"+uid+"_Y")
				delete(allService, "NN"+uid+"_N")

				delete(allCtxC, "NN"+uid+"_Y")
				delete(allCtxC, "NN"+uid+"_N")
			}

			c.String(200, uid+" 종료 신호 전달 완료")
		} else {
			c.String(200, uid+"는 실행 중이지 않습니다.")
		}

	})

	r.GET("/nrun", func(c *gin.Context) {
		var uid string
		uid = c.Query("uid")
		temp := nanoCtxC[uid]
		if temp != nil {
			c.String(200, uid+"이미 실행 중입니다.")
		} else {
			nanoUser[uid] = uid

			if s.EqualFold(config.Conf.PHONE_TYPE_FLAG, "N") { // 기본
				ctx, cancel := context.WithCancel(context.Background())
				ctx = context.WithValue(ctx, "user_id", uid)

				go nanoproc.NanoProcess(uid, ctx)

				nanoCtxC[uid] = cancel

				allCtxC["NN"+uid] = cancel
				allService["NN"+uid] = uid
			} else if s.EqualFold(config.Conf.PHONE_TYPE_FLAG, "Y") { // 콜비서

				ctxY, cancelY := context.WithCancel(context.Background())
				ctxY = context.WithValue(ctxY, "user_id", uid)

				ctxN, cancelN := context.WithCancel(context.Background())
				ctxN = context.WithValue(ctxN, "user_id", uid)

				go nanoproc.NanoProcess_Y(uid, ctxY) // 010으로 시작하는 번호
				go nanoproc.NanoProcess_N(uid, ctxN) // 010이 아닌 번호

				nanoCtxC[uid+"_Y"] = cancelY
				allCtxC["NN"+uid+"_Y"] = cancelY
				allService["NN"+uid+"_Y"] = "NanoService Y"

				nanoCtxC[uid+"_N"] = cancelN
				allCtxC["NN"+uid+"_N"] = cancelN
				allService["NN"+uid+"_N"] = "NanoService N"
			}

			c.String(200, uid+" 시작 신호 전달 완료")
		}
	})

	r.GET("/nlist", func(c *gin.Context) {
		var key string
		for k := range nanoUser {
			key = key + k + "\n"
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
}
