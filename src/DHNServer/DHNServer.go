package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	// "time"

	_ "github.com/go-sql-driver/mysql"

	config "mycs/src/kaoconfig"
	databasepool "mycs/src/kaodatabasepool"

	//"kaoreqreceive"

	//"kaocenter"
	"mycs/src/kaosendrequest"
	"mycs/src/nanoproc"
	"mycs/src/oshotproc"
	"mycs/src/otpproc"
	"mycs/src/ktproc"
	"mycs/src/lguproc"

	"strconv"
	s "strings"

	"github.com/gin-gonic/gin"
	"github.com/takama/daemon"

	"context"
	"sort"
	//"reflect"

)

const (
	name        = "DHNServer2"
	description = "대형네트웍스 카카오 발송 서버"
)

var dependencies = []string{"DHNServer2.service"}

var resultTable string

type Service struct {
	daemon.Daemon
}

func (service *Service) Manage() (string, error) {

	usage := "Usage: DHNServer2 install | remove | start | stop | status"

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

	_, err := databasepool.DB.Exec("update DHN_CLIENT_LIST set pre_send_type = 0, pre_update_date = null")
				
	if err != nil {
		config.Stdlog.Println("DHN_CLIENT_LIST의 pre_send_type, pre_update_date 컬럼 초기화 실패 : ", err)
	}

	allService := map[string]string{}
	allCtxC := map[string]interface{}{}

	alim_user_list, error := databasepool.DB.Query("select distinct user_id from DHN_CLIENT_LIST where use_flag = 'Y' and alimtalk='Y'")
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

	oshotUserList, error := databasepool.DB.Query("select distinct user_id from DHN_CLIENT_LIST dcl  where dcl.use_flag   = 'Y' and IFNULL(dcl.dest, '') = 'OSHOT'")
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
	}

	olctx, olcancel := context.WithCancel(context.Background())
	go oshotproc.LMSProcess(olctx)
	allCtxC["oshotlms"] = olcancel
	allService["oshotlms"] = "Oshot LMS"

	osctx, oscancel := context.WithCancel(context.Background())
	go oshotproc.SMSProcess(osctx)
	allCtxC["oshotsms"] = oscancel
	allService["oshotsms"] = "Oshot SMS"

	ktxroUser := map[string]string{}
	ktxroCtxC := map[string]interface{}{}

	ktxroUserList, error := databasepool.DB.Query("select distinct user_id from DHN_CLIENT_LIST dcl  where dcl.use_flag   = 'Y' and IFNULL(dcl.dest, '') = 'KTXRO'")
	isKtxro := true
	if error != nil {
		config.Stdlog.Println("KT 크로샷 유저 select 오류 ")
		isKtxro = false
	}
	defer ktxroUserList.Close()

	if isKtxro {
		var user_id sql.NullString

		for ktxroUserList.Next() {
			ktxroUserList.Scan(&user_id)
			ctx, cancel := context.WithCancel(context.Background())
			ctx = context.WithValue(ctx, "user_id", user_id.String)
			go ktproc.KtProcess(user_id.String, ctx, 0)

			ktxroUser[user_id.String] = user_id.String
			ktxroCtxC[user_id.String] = cancel

			allCtxC["KTX"+user_id.String] = cancel
			allService["KTX"+user_id.String] = user_id.String

		}
	}

	klctx, klcancel := context.WithCancel(context.Background())
	go ktproc.LMSProcess(klctx, 0)
	allCtxC["ktxrolms"] = klcancel
	allService["ktxrolms"] = "Ktxro LMS"

	Ksctx, Kscancel := context.WithCancel(context.Background())
	go ktproc.SMSProcess(Ksctx, 0)
	allCtxC["ktxrosms"] = Kscancel
	allService["ktxrosms"] = "Ktxro SMS"

	nanoUser := map[string]string{}
	nanoCtxC := map[string]interface{}{}

	nanoUserList, error := databasepool.DB.Query("select distinct user_id, nano_tel_seperate from DHN_CLIENT_LIST dcl  where dcl.use_flag   = 'Y' and IFNULL(dcl.dest, '') = 'NANO'")
	isNano := true
	if error != nil {
		config.Stdlog.Println("Nano 유저 select 오류 ")
		isNano = false
	}
	defer nanoUserList.Close()

	if isNano {
		var user_id sql.NullString
		var nano_tel_seperate sql.NullString

		for nanoUserList.Next() {
			nanoUserList.Scan(&user_id)
			nanoUserList.Scan(&nano_tel_seperate)

			if s.EqualFold(nano_tel_seperate.String, "N") { // 기본
				ctx, cancel := context.WithCancel(context.Background())
				ctx = context.WithValue(ctx, "user_id", user_id.String)

				go nanoproc.NanoProcess(user_id.String, ctx)

				nanoUser[user_id.String] = user_id.String
				nanoCtxC[user_id.String] = cancel

				allCtxC["NN"+user_id.String] = cancel
				allService["NN"+user_id.String] = user_id.String
			} else if s.EqualFold(nano_tel_seperate.String, "Y") { // 콜비서

				ctxY, cancelY := context.WithCancel(context.Background())
				ctxY = context.WithValue(ctxY, "user_id", user_id.String)

				ctxN, cancelN := context.WithCancel(context.Background())
				ctxN = context.WithValue(ctxN, "user_id", user_id.String)

				go nanoproc.NanoProcess_Y(user_id.String, ctxY) // 010으로 시작하는 번호
				go nanoproc.NanoProcess_N(user_id.String, ctxN) // 010이 아닌 번호

				nanoUser[user_id.String+"_Y"] = user_id.String
				nanoCtxC[user_id.String+"_Y"] = cancelY
				allCtxC["NN"+user_id.String+"_Y"] = cancelY
				allService["NN"+user_id.String+"_Y"] = "NanoService Y"

				nanoUser[user_id.String+"_N"] = user_id.String
				nanoCtxC[user_id.String+"_N"] = cancelN
				allCtxC["NN"+user_id.String+"_N"] = cancelN
				allService["NN"+user_id.String+"_N"] = "NanoService N"
			}

		}

	}

	nlctx, nlcancel := context.WithCancel(context.Background())

	go nanoproc.NanoLMSProcess(nlctx)

	allCtxC["nanolms"] = nlcancel
	allService["nanolms"] = "Nano LMS"

	nsctx, nscancel := context.WithCancel(context.Background())

	go nanoproc.NanoSMSProcess(nsctx)

	allCtxC["nanosms"] = nscancel
	allService["nanosms"] = "Nano SMS"

	nlctxG, nlcancelG := context.WithCancel(context.Background())

	go nanoproc.NanoLMSProcess_G(nlctxG)

	allCtxC["nanolmsG"] = nlcancelG
	allService["nanolmsG"] = "Nano LMS G"

	nsctxG, nscancelG := context.WithCancel(context.Background())

	go nanoproc.NanoSMSProcess_G(nsctxG)

	allCtxC["nanosmsG"] = nscancelG
	allService["nanosmsG"] = "Nano SMS G"

	lguUser := map[string]string{}
	lguCtxC := map[string]interface{}{}

	lguUserList, error := databasepool.DB.Query("select distinct user_id from DHN_CLIENT_LIST dcl  where dcl.use_flag   = 'Y' and IFNULL(dcl.dest, '') = 'LGU'")
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

	// llctx, llcancel := context.WithCancel(context.Background())
	// go lguproc.LMSProcess(llctx)
	// allCtxC["lgulms"] = olcancel
	// allService["lgulms"] = "LGU LMS"

	// lsctx, lscancel := context.WithCancel(context.Background())
	// go lguproc.SMSProcess(lsctx)
	// allCtxC["lgusms"] = lscancel
	// allService["lgusms"] = "LGU SMS"


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
/nlist           -> 실행 중인 Nano process User List.

/kstop?uid=dhn       -> dhn KTXRO process stop.
/krun?uid=dhn&acc=0  -> dhn KTXRO process run.
/klist               -> 실행 중인 KTXRO process User List.

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

			// _, err := databasepool.DB.Exec("update DHN_CLIENT_LIST set dest = NULL, pre_send_type = 1, pre_update_date = ? where user_id = ?", time.Now().Format("2006-01-02 15:04:05"), uid)
				
			// if err != nil {
			// 	config.Stdlog.Println(uid," /ostop 오샷 DHN_CLIENT_LIST 업데이트 실패 : ", err)
			// }
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

			// _, err := databasepool.DB.Exec("update DHN_CLIENT_LIST set dest = 'OSHOT' where user_id = ?", uid)
				
			// if err != nil {
			// 	config.Stdlog.Println(uid," /orun 오샷 DHN_CLIENT_LIST 업데이트 실패 : ", err)
			// }

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

			// _, err := databasepool.DB.Exec("update DHN_CLIENT_LIST set alimtalk = 'N' where user_id = ?", uid)
				
			// if err != nil {
			// 	config.Stdlog.Println(uid," /astop 알림톡 DHN_CLIENT_LIST 업데이트 실패 : ", err)
			// }

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

			// _, err := databasepool.DB.Exec("update DHN_CLIENT_LIST set alimtalk = 'Y' where user_id = ?", uid)
				
			// if err != nil {
			// 	config.Stdlog.Println(uid," /arun 알림톡 DHN_CLIENT_LIST 업데이트 실패 : ", err)
			// }

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
			
			if s.EqualFold(config.Conf.PHONE_TYPE_FLAG, "N") { // 기본
				delete(nanoCtxC, uid)
				delete(nanoUser, uid)

				delete(allService, "NN"+uid)
				delete(allCtxC, "NN"+uid)
			} else if s.EqualFold(config.Conf.PHONE_TYPE_FLAG, "Y") { // 콜비서
				delete(nanoUser, uid+"_Y")
				delete(nanoUser, uid+"_N")

				delete(nanoCtxC, uid+"_Y")
				delete(nanoCtxC, uid+"_N")

				delete(allService, "NN"+uid+"_Y")
				delete(allService, "NN"+uid+"_N")

				delete(allCtxC, "NN"+uid+"_Y")
				delete(allCtxC, "NN"+uid+"_N")
			}

			// _, err := databasepool.DB.Exec("update DHN_CLIENT_LIST set dest = NULL, pre_send_type = 2, pre_update_date = ? where user_id = ?", time.Now().Format("2006-01-02 15:04:05"), uid)
				
			// if err != nil {
			// 	config.Stdlog.Println(uid," /nstop 나노 DHN_CLIENT_LIST 업데이트 실패 : ", err)
			// }

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

			// _, err := databasepool.DB.Exec("update DHN_CLIENT_LIST set dest = 'NANO' where user_id = ?", uid)
				
			// if err != nil {
			// 	config.Stdlog.Println(uid," /nrun 나노 DHN_CLIENT_LIST 업데이트 실패 : ", err)
			// }

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

	r.GET("/kstop", func(c *gin.Context) {
		var uid string
		uid = c.Query("uid")
		temp := ktxroCtxC[uid]
		if temp != nil {
			cancel := ktxroCtxC[uid].(context.CancelFunc)
			cancel()
			delete(ktxroCtxC, uid)
			delete(ktxroUser, uid)
			delete(allService, "KTX"+uid)
			delete(allCtxC, "KTX"+uid)
			c.String(200, uid+" 종료 신호 전달 완료")

			// _, err := databasepool.DB.Exec("update DHN_CLIENT_LIST set dest = NULL, pre_send_type = 1, pre_update_date = ? where user_id = ?", time.Now().Format("2006-01-02 15:04:05"), uid)
				
			// if err != nil {
			// 	config.Stdlog.Println(uid," /kstop KT크로샷 DHN_CLIENT_LIST 업데이트 실패 : ", err)
			// }
		} else {
			c.String(200, uid+"는 실행 중이지 않습니다.")
		}

	})

	r.GET("/krun", func(c *gin.Context) {
		var uid string
		uid = c.Query("uid")
		acc := c.Query("acc")
		convAcc, err := strconv.Atoi(acc)
		if err != nil {
			c.String(200, uid+" 에러입니다. err : ", err, "  /  acc : ", acc)
			return
		}
		temp := ktxroCtxC[uid]
		if temp != nil {
			c.String(200, uid+"이미 실행 중입니다.")
		} else {
			ctx, cancel := context.WithCancel(context.Background())
			ctx = context.WithValue(ctx, "user_id", uid)
			go ktproc.KtProcess(uid, ctx, convAcc)

			ktxroCtxC[uid] = cancel
			ktxroUser[uid] = uid

			allCtxC["KTX"+uid] = cancel
			allService["KTX"+uid] = uid

			// _, err := databasepool.DB.Exec("update DHN_CLIENT_LIST set dest = 'KTXRO' where user_id = ?", uid)
				
			// if err != nil {
			// 	config.Stdlog.Println(uid," /krun KT크로샷 DHN_CLIENT_LIST 업데이트 실패 : ", err)
			// }

			c.String(200, uid+" 시작 신호 전달 완료")
		}
	})

	r.GET("/klist", func(c *gin.Context) {
		var key string
		for k := range ktxroUser {
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
