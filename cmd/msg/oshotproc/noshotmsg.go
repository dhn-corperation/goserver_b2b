package oshotproc

import (
	"fmt"
	"sync"
	"time"
	"strconv"
	"context"
	s "strings"
	"database/sql"

	config "mycs/configs"
	databasepool "mycs/configs/databasepool"

	_ "github.com/go-sql-driver/mysql"
)

var errmsg map[string]string = map[string]string{
	"0":  "초기 입력 상태 (default)",
	"1":  "전송 요청 완료(결과수신대기)",
	"3":  "메시지 형식 오류",
	"5":  "휴대폰번호 가입자 없음(미등록)",
	"6":  "전송 성공",
	"7":  "결번(or 서비스 정지)",
	"8":  "단말기 전원 꺼짐",
	"9":  "단말기 음영지역",
	"10": "단말기내 수신메시지함 FULL로 전송 실패 (구:단말 Busy, 기타 단말문제)",
	"11": "기타 전송실패",
	"12": "번호이동",
	"13": "스팸차단 발신번호",
	"14": "스팸차단 수신번호",
	"15": "스팸차단 메시지내용",
	"16": "스팸차단 기타",
	"20": "*단말기 서비스 불가",
	"21": "단말기 서비스 일시정지",
	"22": "단말기 착신 거절",
	"23": "단말기 무응답 및 통화중 (busy)",
	"28": "단말기 MMS 미지원",
	"29": "기타 단말기 문제",
	"35": "발신/착신번호 에러",
	"36": "유효하지 않은 수신번호(망)",
	"37": "유효하지 않은 발신번호(망)",
	"38": "반송건수 한도 초과(망)",
	"39": "동일메시지 제한",
	"49": "이통사 착신전화 오류",
	"50": "이통사 컨텐츠 에러",
	"51": "이통사 전화번호 세칙 미준수 발신번호",
	"52": "이통사 발신번호 변작으로 등록된 발신번호",
	"53": "이통사 번호도용문자 차단서비스에 가입된 발신번호",
	"54": "이통사 발신번호 기타",
	"55": "이통사 전송시간 초과",
	"56": "이통사 동일 메시지 제한",
	"57": "이통사 결과 수신 시간 만료",
	"58": "이통사 리포트 미수신",
	"59": "이통사 기타",
	"60": "컨텐츠 크기 오류(초과 등)",
	"61": "잘못된 메시지 타입",
	"69": "컨텐츠 기타",
	"73": "[Agent] MMS첨부파일 오류. 상대경로체크, 지정경로체크(옵션 – 기본(N / 사용하지않음), 필요 시 설정(Y))",
	"74": "[Agent] 중복발송 차단 (동일한 수신번호와 메시지 발송 - 기본off, 설정필요)",
	"75": "[Agent] 발송 Timeout",
	"76": "[Agent] 유효하지않은 발신번호",
	"77": "[Agent] 유효하지않은 수신번호",
	"78": "[Agent] 컨텐츠 오류 (MMS파일없음 등)",
	"79": "[Agent] 기타",
	"80": "고객필터링 차단 (발신번호, 수신번호, 메시지 등)",
	"81": "080 수신거부",
	"84": "중복발송 차단",
	"86": "유효하지 않은 수신번호",
	"87": "유효하지 않은 발신번호",
	"88": "발신번호 미등록 차단",
	"89": "시스템필터링 기타",
	"90": "발송제한 시간 초과",
	"91": "식별코드 오류",
	"92": "잔액부족",
	"93": "월 발송량 초과",
	"94": "일 발송량 초과",
	"95": "초당 발송량 초과 (재전송 필요)",
	"96": "발송시스템 일시적인 부하 (재전송 필요)",
	"97": "전송 네트워크 오류 (재전송 필요)",
	"98": "외부발송시스템 장애 (재전송 필요)",
	"99": "발송시스템 장애 (재전송 필요)",
}

func NMSGProcess(ctx context.Context) {
	var wg sync.WaitGroup

	for {
		select {
			case <- ctx.Done():
		
		    config.Stdlog.Println("NOshot MSG process가 10초 후에 종료 됨.")
		    time.Sleep(10 * time.Second)
		    config.Stdlog.Println("NOshot MSG process 종료 완료")
		    return
		default:	
			var t = time.Now()

			if t.Day() <= 3 {
				wg.Add(1)
				go pre_msgProcess(&wg)
			}

			wg.Add(1)
			go msgProcess(&wg)
			wg.Wait()
		}
	}
}

func msgProcess(wg *sync.WaitGroup) {
	defer wg.Done()
	defer func(){
		if r := recover(); r != nil {
			config.Stdlog.Println("NOshot msgProcess panic 발생 원인 : ", r)
			if err, ok := r.(error); ok {
				if s.Contains(err.Error(), "connection refused") {
					for {
						config.Stdlog.Println("NOshot msgProcess send ping to DB")
						err := databasepool.DB.Ping()
						if err == nil {
							break
						}
						time.Sleep(10 * time.Second)
					}
				}
			}
		}
	}()
	
	var db = databasepool.DB
	var errlog = config.Stdlog

	var isProc = true
	var t = time.Now()
	var monthStr = fmt.Sprintf("%d%02d", t.Year(), t.Month())

	var MMSTable = "OShotMSG_" + monthStr

	var groupQuery = "select etc1, SendResult, SendDT, MsgID, telecom, etc2  from " + MMSTable + " where etc3 is null "

	groupRows, err := db.Query(groupQuery)
	if err != nil {
		//errlog.Println("스마트미 MMS 조회 중 오류 발생", groupQuery)
		errcode := err.Error()
		errlog.Println("NOShot MSG 조회 중 오류 발생", groupQuery, errcode)

		if s.Index(errcode, "1146") > 0 {
			db.Exec("Create Table IF NOT EXISTS " + MMSTable + " like OShotMSG")
			errlog.Println(MMSTable + " 생성 !!")
		}
		time.Sleep(10 * time.Second)
		isProc = false
		return
	}
	defer groupRows.Close()

	if isProc {

		for groupRows.Next() {
			var cb_msg_id, sendresult, senddt, msgid, telecom, userid sql.NullString

			groupRows.Scan(&cb_msg_id, &sendresult, &senddt, &msgid, &telecom, &userid)

			tr_net := "ETC"

			if s.EqualFold(telecom.String, "011") {
				tr_net = "SKT"
			} else if s.EqualFold(telecom.String, "016") {
				tr_net = "KTF"
			} else if s.EqualFold(telecom.String, "019") {
				tr_net = "LGT"
			}

			if !s.EqualFold(sendresult.String, "6") {
				numcode, _ := strconv.Atoi(sendresult.String)
				var errcode = fmt.Sprintf("%d%03d", 7, numcode)

				val, exists := errmsg[sendresult.String]
				if !exists {
					val = "기타 오류"
				}

				db.Exec("update DHN_RESULT dr set dr.message_type = 'PH', dr.result = 'Y', dr.code = '" + errcode + "', dr.message = concat(dr.message, '," + val + "'), dr.remark1 = '" + telecom.String + "', dr.remark2 = '" + senddt.String + "' where userid='" + userid.String + "' and msgid = '" + cb_msg_id.String + "'")
			} else {
				db.Exec("update DHN_RESULT dr set dr.message_type = 'PH', dr.result = 'Y', dr.code = '0000', dr.message = '', dr.remark1 = '" + tr_net + "', dr.remark2 = '" + senddt.String + "' where userid='" + userid.String + "' and  msgid = '" + cb_msg_id.String + "'")
			}

			db.Exec("update " + MMSTable + " set etc3 = 'Y' where MsgID = '" + msgid.String + "'")
		}
	}

}

func pre_msgProcess(wg *sync.WaitGroup) {
	defer wg.Done()
	defer func(){
		if r := recover(); r != nil {
			config.Stdlog.Println("NOShot msgProcess panic 발생 원인 : ", r)
			if err, ok := r.(error); ok {
				if s.Contains(err.Error(), "connection refused") {
					for {
						config.Stdlog.Println("NOShot msgProcess send ping to DB")
						err := databasepool.DB.Ping()
						if err == nil {
							break
						}
						time.Sleep(10 * time.Second)
					}
				}
			}
		}
	}()
	
	var db = databasepool.DB
	var errlog = config.Stdlog

	var isProc = true
	var t = time.Now().Add(time.Hour * -96)
	var monthStr = fmt.Sprintf("%d%02d", t.Year(), t.Month())

	var MMSTable = "OShotMSG_" + monthStr

	var groupQuery = "select etc1, SendResult, SendDT, MsgID, telecom, etc2  from " + MMSTable + " where etc3 is null "

	groupRows, err := db.Query(groupQuery)
	if err != nil {
		errcode := err.Error()
		errlog.Println("NOShot MSG 조회 중 오류 발생", groupQuery, errcode)

		if s.Index(errcode, "1146") > 0 {
			db.Exec("Create Table IF NOT EXISTS " + MMSTable + " like OShotMSG")
			errlog.Println(MMSTable + " 생성 !!")
		}
		time.Sleep(10 * time.Second)
		isProc = false
		return
	}
	defer groupRows.Close()

	if isProc {

		for groupRows.Next() {
			var cb_msg_id, sendresult, senddt, msgid, telecom, userid sql.NullString

			groupRows.Scan(&cb_msg_id, &sendresult, &senddt, &msgid, &telecom, &userid)

			if !s.EqualFold(sendresult.String, "6") {
				numcode, _ := strconv.Atoi(sendresult.String)
				var errcode = fmt.Sprintf("%d%03d", 7, numcode)

				val, exists := errmsg[sendresult.String]
				if !exists {
					val = "기타 오류"
				}

				db.Exec("update DHN_RESULT dr set dr.message_type = 'PH', dr.result = 'Y', dr.code = '" + errcode + "', dr.message = concat(dr.message, '," + val + "'), dr.remark1 = '" + telecom.String + "', dr.remark2 = '" + senddt.String + "' where userid='" + userid.String + "' and  msgid = '" + cb_msg_id.String + "'")
			} else {
				db.Exec("update DHN_RESULT dr set dr.message_type = 'PH', dr.result = 'Y', dr.code = '0000', dr.message = '', dr.remark1 = '" + telecom.String + "', dr.remark2 = '" + senddt.String + "' where  userid='" + userid.String + "' and msgid = '" + cb_msg_id.String + "'")
			}

			db.Exec("update " + MMSTable + " set etc3 = 'Y' where MsgID = '" + msgid.String + "'")
		}
	}

}
