package kaoconfig

import (
	"fmt"
	"log"
	"os"
	"time"

	ini "github.com/BurntSushi/toml"
	"github.com/go-resty/resty/v2"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
)

type Config struct {
	DB              string
	DBURL           string
	PORT            string
	PROFILE_KEY     string
	API_SERVER      string
	CENTER_SERVER   string
	IMAGE_SERVER    string
	CHANNEL         string
	RESPONSE_METHOD string
	SENDLIMIT       int
	PHONE_MSG_FLAG  string
	OTP_MSG_FLAG    string
	DEBUG           string
}

type ErrRes struct {
	Key   string
	Value string
}

var Conf Config
var Stdlog *log.Logger
var BasePath string
var IsRunning bool = true
var ResultLimit int = 1000
var Client *resty.Client
var ErrCode map[string]ErrRes // 다우오피스 에러코드 매핑

func InitConfig() {
	path := "/root/DHNNanoServer/log/DHNNanoServer"
	//path := "./log/DHNServer"
	loc, _ := time.LoadLocation("Asia/Seoul")
	writer, err := rotatelogs.New(
		fmt.Sprintf("%s-%s.log", path, "%Y-%m-%d"),
		rotatelogs.WithLocation(loc),
		rotatelogs.WithMaxAge(-1),
		rotatelogs.WithRotationCount(7),
	)

	if err != nil {
		log.Fatalf("Failed to Initialize Log File %s", err)
	}

	log.SetOutput(writer)
	stdlog := log.New(os.Stdout, "INFO -> ", log.Ldate|log.Ltime)
	stdlog.SetOutput(writer)
	Stdlog = stdlog

	Conf = readConfig()
	errCode()
	BasePath = "/root/DHNNanoServer/"
	Client = resty.New()

}

func readConfig() Config {
	var configfile = "/root/DHNNanoServer/config.ini"
	//var configfile = "./config.ini"
	_, err := os.Stat(configfile)
	if err != nil {
		fmt.Println("Config file is missing : ", configfile)
	}

	var result Config
	_, err1 := ini.DecodeFile(configfile, &result)

	if err1 != nil {
		fmt.Println("Config file read error : ", err1)
	}

	return result
}

func InitCenterConfig() {
	path := "/root/DHNNanoCenter/log/DHNNanoCenter"
	//	path := "./log/DHNNanoCenter"
	loc, _ := time.LoadLocation("Asia/Seoul")
	writer, err := rotatelogs.New(
		fmt.Sprintf("%s-%s.log", path, "%Y-%m-%d"),
		rotatelogs.WithLocation(loc),
		rotatelogs.WithMaxAge(-1),
		rotatelogs.WithRotationCount(7),
	)

	if err != nil {
		log.Fatalf("Failed to Initialize Log File %s", err)
	}

	log.SetOutput(writer)
	stdlog := log.New(os.Stdout, "INFO -> ", log.Ldate|log.Ltime)
	stdlog.SetOutput(writer)
	Stdlog = stdlog

	Conf = readCenterConfig()
	BasePath = "/root/DHNNanoCenter/"
	//Client = resty.New()

}

func readCenterConfig() Config {
	var configfile = "/root/DHNNanoCenter/config.ini"
	//	var configfile = "./config.ini"
	_, err := os.Stat(configfile)
	if err != nil {
		fmt.Println("Config file is missing : ", configfile)
	}

	var result Config
	_, err1 := ini.DecodeFile(configfile, &result)

	if err1 != nil {
		fmt.Println("Config file read error : ", err1)
	}

	return result
}

// 다우오피스 에러코드 매핑
func errCode() {
	ErrCode = make(map[string]ErrRes)
	ErrCode["0000"] = ErrRes{Key: "7000", Value: "전달"}
	ErrCode["1001"] = ErrRes{Key: "7101", Value: "카카오 형식 오류"}
	ErrCode["1003"] = ErrRes{Key: "7103", Value: "Sender key (발신 프로필키) 유효하지 않음"}
	ErrCode["1006"] = ErrRes{Key: "7106", Value: "삭제된 Sender key (발신 프로필키)"}
	ErrCode["1007"] = ErrRes{Key: "7107", Value: "차단 상태 Sender key (발신 프로필키)"}
	ErrCode["1021"] = ErrRes{Key: "7108", Value: "차단 상태 카카오톡 채널 (카카오톡 채널 운영툴에서 확인)"}
	ErrCode["1022"] = ErrRes{Key: "7109", Value: "닫힌 상태 카카오톡 채널 (카카오톡 채널 운영툴에서 확인)"}
	ErrCode["1023"] = ErrRes{Key: "7110", Value: "삭제된 카카오톡 채널 (카카오톡 채널 운영툴에서 확인)"}
	ErrCode["1024"] = ErrRes{Key: "7111", Value: "삭제 대기 상태의 카카오톡 채널 (카카오톡 채널 운영툴에서 확인)"}
	ErrCode["1014"] = ErrRes{Key: "7112", Value: "유효하지 않은 사업자번호"}
	ErrCode["1025"] = ErrRes{Key: "7125", Value: "메시지 차단 상태의 카카오톡 채널 (카카오톡 채널 운영툴에서 확인)"}
	ErrCode["1030"] = ErrRes{Key: "7130", Value: "잘못된 파라미터 요청"}
	ErrCode["2003"] = ErrRes{Key: "7203", Value: "친구톡 전송 시 친구대상 아님"}
	ErrCode["3016"] = ErrRes{Key: "7204", Value: "템플릿 불일치"}
	ErrCode["2006"] = ErrRes{Key: "7206", Value: "시리얼넘버 형식 불일치"}
	ErrCode["9998"] = ErrRes{Key: "7300", Value: "기타에러"}
	ErrCode["3005"] = ErrRes{Key: "7305", Value: "성공 불확실 (30 일 이내 수신 가능)"}
	ErrCode["9999"] = ErrRes{Key: "7306", Value: "카카오 시스템 오류"}
	ErrCode["3008"] = ErrRes{Key: "7308", Value: "전화번호 오류"}
	ErrCode["3013"] = ErrRes{Key: "7311", Value: "메시지가 존재하지 않음"}
	ErrCode["3014"] = ErrRes{Key: "7314", Value: "메시지 길이 초과"}
	ErrCode["3014"] = ErrRes{Key: "7315", Value: "템플릿 없음"}
	ErrCode["3018"] = ErrRes{Key: "7318", Value: "메시지를 전송할 수 없음"}
	ErrCode["3022"] = ErrRes{Key: "7322", Value: "메시지 발송 불가 시간"}
	ErrCode["3024"] = ErrRes{Key: "7324", Value: "재전송 메시지 존재하지 않음"}
	ErrCode["3025"] = ErrRes{Key: "7325", Value: "변수 글자수 제한 초과"}
	ErrCode["3026"] = ErrRes{Key: "7326", Value: "상담/봇 전환 버튼 extra, event 글자수 제한 초과"}
	ErrCode["3027"] = ErrRes{Key: "7327", Value: "버튼 내용이 템플릿과 일치하지 않음"}
	ErrCode["3028"] = ErrRes{Key: "7328", Value: "메시지 강조 표기 타이틀이 템플릿과 일치하지 않음"}
	ErrCode["3029"] = ErrRes{Key: "7329", Value: "메시지 강조 표기 타이틀 길이 제한 초과 (50 자)"}
	ErrCode["3030"] = ErrRes{Key: "7330", Value: "메시지 타입과 템플릿 강조유형이 일치하지 않음"}
	ErrCode["3031"] = ErrRes{Key: "7331", Value: "헤더가 템플릿과 일치하지 않음"}
	ErrCode["3032"] = ErrRes{Key: "7332", Value: "헤더 길이 제한 초과(16 자)"}
	ErrCode["3033"] = ErrRes{Key: "7333", Value: "아이템 하이라이트가 템플릿과 일치하지 않음"}
	ErrCode["3034"] = ErrRes{Key: "7334", Value: "아이템 하이라이트 타이틀 길이 제한 초과 (이미지 없는 경우 30 자, 이미지 있는 경우 21 자)"}
	ErrCode["3035"] = ErrRes{Key: "7335", Value: "아이템 하이라이트 디스크립션 길이 제한 초과 (이미지 없는 경우 19 자, 이미지 있는 경우 13 자)"}
	ErrCode["3036"] = ErrRes{Key: "7336", Value: "아이템 리스트가 템플릿과 일치하지 않음"}
	ErrCode["3037"] = ErrRes{Key: "7337", Value: "아이템 리스트의 아이템의 디스크립션 길이 제한 초과(23 자)"}
	ErrCode["3038"] = ErrRes{Key: "7338", Value: "아이템 요약정보가 템플릿과 일치하지 않음"}
	ErrCode["3039"] = ErrRes{Key: "7339", Value: "아이템 요약정보의 디스크립션 길이 제한 초과(14 자)"}
	ErrCode["3040"] = ErrRes{Key: "7340", Value: "아이템 요약정보의 디스크립션에 허용되지 않은 문자 포함 (통화 기호/코드, 숫자, 콤마, 소수점, 공백을 제외한 문자 포함)"}
}
