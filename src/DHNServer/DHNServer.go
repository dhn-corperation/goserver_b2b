package main

import (
	"os"
	"fmt"
	"sort"
	"syscall"
	"context"
	"strconv"
	"os/signal"
	"database/sql"
	s "strings"


	"mycs/src/kaosendrequest"
	"mycs/src/nanoproc"
	"mycs/src/oshotproc"
	"mycs/src/ktproc"
	"mycs/src/lguproc"
	"mycs/src/otpatproc"
	"mycs/src/otplguproc"
	"mycs/src/otpnanoproc"
	config "mycs/src/kaoconfig"
	databasepool "mycs/src/kaodatabasepool"

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

	ktxroUser := map[string]string{}
	ktxroCtxC := map[string]interface{}{}

	ktxroUserList, error := databasepool.DB.Query("select distinct user_id from DHN_CLIENT_LIST dcl where dcl.use_flag = 'Y' and upper(ifnull(dcl.dest, '')) = 'KTXRO'")
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

	nanoUserList, error := databasepool.DB.Query("select distinct user_id, nano_tel_seperate from DHN_CLIENT_LIST dcl where dcl.use_flag = 'Y' and upper(ifnull(dcl.dest, '')) = 'NANO'")
	isNano := true
	if error != nil {
		config.Stdlog.Println("나노 유저 select 오류 ")
		isNano = false
	}
	defer nanoUserList.Close()

	if isNano {
		var user_id sql.NullString
		var nano_tel_seperate sql.NullString

		for nanoUserList.Next() {
			nanoUserList.Scan(&user_id, &nano_tel_seperate)

			if s.EqualFold(nano_tel_seperate.String, "N") { // 기본
				ctx, cancel := context.WithCancel(context.Background())
				ctx = context.WithValue(ctx, "user_id", user_id.String)

				go nanoproc.NanoProcess(user_id.String, ctx)

				nanoUser[user_id.String] = user_id.String
				nanoCtxC[user_id.String] = cancel

				allCtxC["NN"+user_id.String] = cancel
				allService["NN"+user_id.String] = user_id.String
			} else if s.EqualFold(nano_tel_seperate.String, "Y") { // 저가

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
	otpNanoUser := map[string]string{}
	otpNanoCtxC := map[string]interface{}{}

	oatctx, cancel := context.WithCancel(context.Background())
	go otpatproc.AlimtalkProc(oatctx)

	otpAtUser["OTPAT"] = "OTPAT"
	otpAtCtxC["OTPAT"] = cancel

	allCtxC["OTPAT"] = cancel
	allService["OTPAT"] = "OTPAT"

	ctx, cancel := context.WithCancel(context.Background())
	go otpnanoproc.NanoProcess(ctx)

	otpNanoUser["OTPNANO"] = "OTPNANO"
	otpNanoCtxC["OTPNANO"] = cancel

	allCtxC["OTPNANO"] = cancel
	allService["OTPNANO"] = "OTPNANO"

	ollctx, ollcancel := context.WithCancel(context.Background())
	go otplguproc.LMSProcess(ollctx)
	allCtxC["otplgulms"] = ollcancel
	allService["otplgulms"] = "LGU OTP LMS"

	olsctx, olscancel := context.WithCancel(context.Background())
	go otplguproc.SMSProcess(olsctx)
	allCtxC["otplgusms"] = olscancel
	allService["otplgusms"] = "LGU OTP SMS"

	onlctx, onlcancel := context.WithCancel(context.Background())
	go otpnanoproc.LMSProcess(onlctx)
	allCtxC["otpnanolms"] = onlcancel
	allService["otpnanolms"] = "NANO OTP LMS"

	onsctx, onscancel := context.WithCancel(context.Background())
	go otpnanoproc.SMSProcess(onsctx)
	allCtxC["otpnanosms"] = onscancel
	allService["otpnanosms"] = "NANO OTP SMS"

	//OTP 영역 종료

	//API 영역 시작
	r := gin.New()
	r.Use(gin.Recovery())
	serCmd := `DHN Server API
Command :
/ostop?uid=dhn   	 	-> dhn Oshot process stop.
/orun?uid=dhn    	 	-> dhn Oshot process run.
/olist           	 	-> 실행 중인 Oshot process User List.

/nstop?uid=dhn   	 	-> dhn Nano process stop.
/nrun?uid=dhn        	-> dhn Nano process run.
/nlist               	-> 실행 중인 Nano process User List.

/kstop?uid=dhn       	-> dhn KTXRO process stop.
/krun?uid=dhn&acc=0  	-> dhn KTXRO process run.
/klist               	-> 실행 중인 KTXRO process User List.

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

	r.GET("/nstop", func(c *gin.Context) {
		var uid string
		uid = c.Query("uid")
		temp := nanoCtxC[uid]
		if temp != nil {
			var nano_tel_seperate sql.NullString
			var nts string

			seperr := databasepool.DB.QueryRow("select distinct nano_tel_seperate from DHN_CLIENT_LIST where user_id = ?", uid).Scan(&nano_tel_seperate)
			if seperr != nil && seperr != sql.ErrNoRows {
				config.Stdlog.Println(uid," /nrun 나노 nano_tel_seperate 습득 실패 : ", seperr)
				nts = "N"
			} else {
				nts = nano_tel_seperate.String
			}

			cancel := nanoCtxC[uid].(context.CancelFunc)
			cancel()
			
			if s.EqualFold(nts, "N") { // 기본
				delete(nanoCtxC, uid)
				delete(nanoUser, uid)

				delete(allService, "NN"+uid)
				delete(allCtxC, "NN"+uid)
			} else if s.EqualFold(nts, "Y") { // 콜비서
				delete(nanoCtxC, uid+"_Y")
				delete(nanoUser, uid+"_Y")
				delete(allService, "NN"+uid+"_Y")
				delete(allCtxC, "NN"+uid+"_Y")

				delete(nanoCtxC, uid+"_N")
				delete(nanoUser, uid+"_N")
				delete(allService, "NN"+uid+"_N")
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
			c.String(200, uid+" 이미 실행 중입니다.")
		} else {
			var nano_tel_seperate sql.NullString
			var nts string

			seperr := databasepool.DB.QueryRow("select distinct nano_tel_seperate from DHN_CLIENT_LIST where user_id = ?", uid).Scan(&nano_tel_seperate)
			if seperr != nil && seperr != sql.ErrNoRows {
				config.Stdlog.Println(uid," /nrun 나노 nano_tel_seperate 습득 실패 : ", seperr)
				nts = "N"
			} else {
				nts = nano_tel_seperate.String
			}


			if s.EqualFold(nts, "N") { // 기본
				ctx, cancel := context.WithCancel(context.Background())
				ctx = context.WithValue(ctx, "user_id", uid)

				go nanoproc.NanoProcess(uid, ctx)

				nanoCtxC[uid] = cancel
				nanoUser[uid] = uid

				allCtxC["NN"+uid] = cancel
				allService["NN"+uid] = uid
			} else if s.EqualFold(nts, "Y") { // 콜비서

				ctxY, cancelY := context.WithCancel(context.Background())
				ctxY = context.WithValue(ctxY, "user_id", uid)

				ctxN, cancelN := context.WithCancel(context.Background())
				ctxN = context.WithValue(ctxN, "user_id", uid)

				go nanoproc.NanoProcess_Y(uid, ctxY) // 010으로 시작하는 번호
				go nanoproc.NanoProcess_N(uid, ctxN) // 010이 아닌 번호

				nanoCtxC[uid+"_Y"] = cancelY
				nanoUser[uid+"_Y"] = uid
				allCtxC["NN"+uid+"_Y"] = cancelY
				allService["NN"+uid+"_Y"] = "NanoService Y"

				nanoCtxC[uid+"_N"] = cancelN
				nanoUser[uid+"_N"] = uid
				allCtxC["NN"+uid+"_N"] = cancelN
				allService["NN"+uid+"_N"] = "NanoService N"
			}

			_, err := databasepool.DB.Exec("update DHN_CLIENT_LIST set dest = 'NANO' where use_flag = 'Y' and user_id = ?", uid)
				
			if err != nil {
				config.Stdlog.Println(uid," /nrun 나노 DHN_CLIENT_LIST 업데이트 실패 : ", err)
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
			c.String(200, uid+" 이미 실행 중입니다.")
		} else {
			ctx, cancel := context.WithCancel(context.Background())
			ctx = context.WithValue(ctx, "user_id", uid)
			go ktproc.KtProcess(uid, ctx, convAcc)

			ktxroCtxC[uid] = cancel
			ktxroUser[uid] = uid

			allCtxC["KTX"+uid] = cancel
			allService["KTX"+uid] = uid

			_, err := databasepool.DB.Exec("update DHN_CLIENT_LIST set dest = 'KTXRO' where use_flag = 'Y' and user_id = ?", uid)
				
			if err != nil {
				config.Stdlog.Println(uid," /krun KT크로샷 DHN_CLIENT_LIST 업데이트 실패 : ", err)
			}

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
		temp = otpNanoCtxC["OTPNANO"]
		if temp != nil {
			cancel := otpNanoCtxC["OTPNANO"].(context.CancelFunc)
			cancel()
			delete(otpNanoCtxC, "OTPNANO")
			delete(otpNanoUser, "OTPNANO")

			delete(allService, "OTPNANO")
			delete(allCtxC, "OTPNANO")
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
				go otpnanoproc.NanoProcess(ctx)

				otpNanoUser["OTPNANO"] = "OTPNANO"
				otpNanoCtxC["OTPNANO"] = cancel

				allCtxC["OTPNANO"] = cancel
				allService["OTPNANO"] = "OTPNANO"

				c.String(200, "OTP NANO 시작 신호 전달 완료")
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

		for l := range otpNanoUser {
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
