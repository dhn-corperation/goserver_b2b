package oshotproc

import (
	"database/sql"
	config "kaoconfig"
	databasepool "kaodatabasepool"
	"fmt"

	"strconv"
	s "strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func SMSProcess() {
	var wg sync.WaitGroup
	
	
	var db = databasepool.DB
	var errlog = config.Stdlog
	var oshotTable [][]string
	var otable sql.NullString
	
	var OshotQuery = "select distinct a.oshot from DHN_CLIENT_LIST a where a.use_flag = 'Y' "

	OshotTable, err := db.Query(OshotQuery)

	if err != nil {
		errlog.Fatal("DHN CLIENT LIST 조회 오류 ")
	}
	defer OshotTable.Close()
	
	for OshotTable.Next() {
		OshotTable.Scan(&otable)
		oshotTable = append(oshotTable, []string{otable.String})
	}
	errlog.Println("Oshot SMS length : ", len(oshotTable))
	for {
		for _, tableName := range oshotTable {
			var t = time.Now()
			
			if t.Day() < 3 {
				wg.Add(1)
				go pre_smsProcess(&wg, tableName[0])
			}
		
			wg.Add(1)
			go smsProcess(&wg, tableName[0])
		}

		wg.Wait()

	}

}

func smsProcess(wg *sync.WaitGroup, ostable string) {

	defer wg.Done()
	var db = databasepool.DB
	var errlog = config.Stdlog

	var isProc = true
	var t = time.Now()
	var monthStr = fmt.Sprintf("%d%02d", t.Year(), t.Month())

	var SMSTable = ostable + "SMS_" + monthStr

	//db.Exec("UPDATE OShotSMS SET SendDT=now(), SendResult='6', Telecom='000' WHERE SendResult=1 and date_add(insertdt, interval 6 HOUR) < now()")
	//db.Exec("insert into " + SMSTable + " SELECT * FROM OShotSMS WHERE SendResult>1 AND SendDT is not null and telecom = '000'")
	//db.Exec("delete FROM OShotSMS WHERE SendResult>1 AND SendDT is not null and telecom = '000'")

    errmsg := map[string]string{
"0"  :  "초기 입력 상태 (default)",
"1"  :  "전송 요청 완료(결과수신대기)",
"3"  :  "메시지 형식 오류",
"5"  :  "휴대폰번호 가입자 없음(미등록)",
"6"  :  "전송 성공",
"7"  :  "결번(or 서비스 정지)",
"8"  :  "단말기 전원 꺼짐",
"9"  :  "단말기 음영지역",
"10"  : "단말기내 수신메시지함 FULL로 전송 실패 (구:단말 Busy, 기타 단말문제)",
"11"  : "기타 전송실패",
"13"  : "스팸차단 발신번호",
"14"  : "스팸차단 수신번호",
"15"  : "스팸차단 메시지내용",
"16"  : "스팸차단 기타",
"20"  : "*단말기 서비스 불가",
"21"  : "단말기 서비스 일시정지",
"22"  : "단말기 착신 거절",
"23"  : "단말기 무응답 및 통화중 (busy)",
"28"  : "단말기 MMS 미지원",
"29"  : "기타 단말기 문제",
"36"  : "유효하지 않은 수신번호(망)",
"37"  : "유효하지 않은 발신번호(망)",
"50"  : "이통사 컨텐츠 에러",
"51"  : "이통사 전화번호 세칙 미준수 발신번호",
"52"  : "이통사 발신번호 변작으로 등록된 발신번호",
"53"  : "이통사 번호도용문자 차단서비스에 가입된 발신번호",
"54"  : "이통사 발신번호 기타",
"59"  : "이통사 기타",
"60"  : "컨텐츠 크기 오류(초과 등)",
"61"  : "잘못된 메시지 타입",
"69"  : "컨텐츠 기타",
"74"  : "[Agent] 중복발송 차단 (동일한 수신번호와 메시지 발송 - 기본off, 설정필요)",
"75"  : "[Agent] 발송 Timeout",
"76"  : "[Agent] 유효하지않은 발신번호",
"77"  : "[Agent] 유효하지않은 수신번호",
"78"  : "[Agent] 컨텐츠 오류 (MMS파일없음 등)",
"79"  : "[Agent] 기타",
"80"  : "고객필터링 차단 (발신번호, 수신번호, 메시지 등)",
"81"  : "080 수신거부",
"84"  : "중복발송 차단",
"86"  : "유효하지 않은 수신번호",
"87"  : "유효하지 않은 발신번호",
"88"  : "발신번호 미등록 차단",
"89"  : "시스템필터링 기타",
"90"  : "발송제한 시간 초과",
"92"  : "잔액부족",
"93"  : "월 발송량 초과",
"94"  : "일 발송량 초과",
"95"  : "초당 발송량 초과 (재전송 필요)",
"96"  : "발송시스템 일시적인 부하 (재전송 필요)",
"97"  : "전송 네트워크 오류 (재전송 필요)",
"98"  : "외부발송시스템 장애 (재전송 필요)",
"99"  : "발송시스템 장애 (재전송 필요)",
    }
    
	//발송 6시간 지난 메세지는 응답과 상관 없이 성공 처리 함.

	var groupQuery = "select cb_msg_id, SendResult, SendDT, MsgID, telecom,userid  from " + SMSTable + " a where a.proc_flag = 'Y' "

	groupRows, err := db.Query(groupQuery)
	if err != nil {
		errlog.Println("스마트미 SMS 조회 중 오류 발생")
		errcode := err.Error()

		if s.Index(errcode, "1146") > 0 {
			db.Exec("Create Table IF NOT EXISTS " + SMSTable + " like " + ostable + "SMS")
			errlog.Println(SMSTable + " 생성 !!")

		} else {
			//errlog.Fatal(err)
		}

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
				
				db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.message_type = 'PH', dr.code = '" + errcode + "', dr.message = concat(dr.message, ',"+ val +"'), dr.remark1 = '" + telecom.String + "', dr.remark2 = '" +senddt.String+"' where  userid='" + userid.String + "' and msgid = '" + cb_msg_id.String + "'")
			} else {
				db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.message_type = 'PH', dr.code = '0000', dr.message = '', dr.remark1 = '" + tr_net + "', dr.remark2 = '" +senddt.String+"' where  userid='" + userid.String + "' and msgid = '" + cb_msg_id.String + "'")			
			}
 
			db.Exec("update " + SMSTable + " set proc_flag = 'N' where msgid = '" + msgid.String + "'") 
		}
	}
}


func pre_smsProcess(wg *sync.WaitGroup, ostable string) {

	defer wg.Done()
	var db = databasepool.DB
	var errlog = config.Stdlog

	var isProc = true
	var t = time.Now().Add(time.Hour * -96)
	var monthStr = fmt.Sprintf("%d%02d", t.Year(), t.Month())

	var SMSTable = ostable + "SMS_" + monthStr

	//db.Exec("UPDATE OShotSMS SET SendDT=now(), SendResult='6', Telecom='000' WHERE SendResult=1 and date_add(insertdt, interval 6 HOUR) < now()")
	//db.Exec("insert into " + SMSTable + " SELECT * FROM OShotSMS WHERE SendResult>1 AND SendDT is not null and telecom = '000'")
	//db.Exec("delete FROM OShotSMS WHERE SendResult>1 AND SendDT is not null and telecom = '000'")

    errmsg := map[string]string{
"0"  :  "초기 입력 상태 (default)",
"1"  :  "전송 요청 완료(결과수신대기)",
"3"  :  "메시지 형식 오류",
"5"  :  "휴대폰번호 가입자 없음(미등록)",
"6"  :  "전송 성공",
"7"  :  "결번(or 서비스 정지)",
"8"  :  "단말기 전원 꺼짐",
"9"  :  "단말기 음영지역",
"10"  : "단말기내 수신메시지함 FULL로 전송 실패 (구:단말 Busy, 기타 단말문제)",
"11"  : "기타 전송실패",
"13"  : "스팸차단 발신번호",
"14"  : "스팸차단 수신번호",
"15"  : "스팸차단 메시지내용",
"16"  : "스팸차단 기타",
"20"  : "*단말기 서비스 불가",
"21"  : "단말기 서비스 일시정지",
"22"  : "단말기 착신 거절",
"23"  : "단말기 무응답 및 통화중 (busy)",
"28"  : "단말기 MMS 미지원",
"29"  : "기타 단말기 문제",
"36"  : "유효하지 않은 수신번호(망)",
"37"  : "유효하지 않은 발신번호(망)",
"50"  : "이통사 컨텐츠 에러",
"51"  : "이통사 전화번호 세칙 미준수 발신번호",
"52"  : "이통사 발신번호 변작으로 등록된 발신번호",
"53"  : "이통사 번호도용문자 차단서비스에 가입된 발신번호",
"54"  : "이통사 발신번호 기타",
"59"  : "이통사 기타",
"60"  : "컨텐츠 크기 오류(초과 등)",
"61"  : "잘못된 메시지 타입",
"69"  : "컨텐츠 기타",
"74"  : "[Agent] 중복발송 차단 (동일한 수신번호와 메시지 발송 - 기본off, 설정필요)",
"75"  : "[Agent] 발송 Timeout",
"76"  : "[Agent] 유효하지않은 발신번호",
"77"  : "[Agent] 유효하지않은 수신번호",
"78"  : "[Agent] 컨텐츠 오류 (MMS파일없음 등)",
"79"  : "[Agent] 기타",
"80"  : "고객필터링 차단 (발신번호, 수신번호, 메시지 등)",
"81"  : "080 수신거부",
"84"  : "중복발송 차단",
"86"  : "유효하지 않은 수신번호",
"87"  : "유효하지 않은 발신번호",
"88"  : "발신번호 미등록 차단",
"89"  : "시스템필터링 기타",
"90"  : "발송제한 시간 초과",
"92"  : "잔액부족",
"93"  : "월 발송량 초과",
"94"  : "일 발송량 초과",
"95"  : "초당 발송량 초과 (재전송 필요)",
"96"  : "발송시스템 일시적인 부하 (재전송 필요)",
"97"  : "전송 네트워크 오류 (재전송 필요)",
"98"  : "외부발송시스템 장애 (재전송 필요)",
"99"  : "발송시스템 장애 (재전송 필요)",
    }
    
	//발송 6시간 지난 메세지는 응답과 상관 없이 성공 처리 함.

	var groupQuery = "select cb_msg_id, SendResult, SendDT, MsgID, telecom, userid from " + SMSTable + " a where a.proc_flag = 'Y' "

	groupRows, err := db.Query(groupQuery)
	if err != nil {
		errlog.Println("스마트미 SMS 조회 중 오류 발생")
		errcode := err.Error()

		if s.Index(errcode, "1146") > 0 {
			db.Exec("Create Table IF NOT EXISTS " + SMSTable + " like " + ostable + "SMS")
			errlog.Println(SMSTable + " 생성 !!")

		} else {
			//errlog.Fatal(err)
		}

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
				
				db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.message_type = 'PH', dr.code = '" + errcode + "', dr.message = concat(dr.message, ',"+ val +"'), dr.remark1 = '" + telecom.String + "', dr.remark2 = '" +senddt.String+"' where  userid='" + userid.String + "' and msgid = '" + cb_msg_id.String + "'")
			} else {
				db.Exec("update DHN_RESULT dr set dr.result = 'Y', dr.message_type = 'PH', dr.code = '0000', dr.message = '', dr.remark1 = '" + telecom.String + "', dr.remark2 = '" +senddt.String+"' where  userid='" + userid.String + "' and msgid = '" + cb_msg_id.String + "'")			
			}
 
			db.Exec("update " + SMSTable + " set proc_flag = 'N' where msgid = '" + msgid.String + "'") 
		}
	}
}
