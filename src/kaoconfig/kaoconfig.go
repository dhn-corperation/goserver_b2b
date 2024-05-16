package kaoconfig

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	// "net"
	"net/http"

	ini "github.com/BurntSushi/toml"
	"github.com/go-resty/resty/v2"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
)

type Config struct {
	DB               string
	DBURL            string

	CENTER_PORT      string

	KAKAO_SENDLIMIT  int
	SERVER_PORT      string
	PROFILE_KEY      string
	RESPONSE_METHOD  string
	PHONE_MSG_FLAG   string
	CHANNEL          string
	DEBUG            string

	API_SERVER       string
	CENTER_SERVER    string
	IMAGE_SERVER     string

	NANO_IDENTI_CODE string
	PHONE_TYPE_FLAG  string
	
	OTP_MSG_FLAG     string
}

var Conf Config
var Stdlog *log.Logger
var BasePath string
var IsRunning bool = true
var ResultLimit int = 1000
var Client *resty.Client
var GoClient *http.Client = &http.Client{
	Timeout: time.Second * 30,
	Transport: &http.Transport{
		TLSHandshakeTimeout: 10 * time.Second,
	},
}
var StrSpecialMap map[string]string

func InitConfig() {
	realpath, _ := os.Executable()
	dir := filepath.Dir(realpath)
	logDir := filepath.Join(dir, "logs")
	err := createDir(logDir)
	if err != nil {
		log.Fatalf("Failed to ensure log directory: %s", err)
	}
	path := filepath.Join(logDir, "DHNServer")
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
	BasePath = dir + "/"
	Client = resty.New()

}

func readConfig() Config {
	realpath, _ := os.Executable()
	dir := filepath.Dir(realpath)
	var configfile = filepath.Join(dir, "config.ini")
	_, err := os.Stat(configfile)
	if err != nil {

		err := createConfig(configfile)
		if err != nil {
			Stdlog.Println("Config file create fail")
		}
		Stdlog.Println("config.ini 생성완료 작성을 해주세요.")

		system_exit("DHNServer")
		fmt.Println("Config file is missing : ", configfile)
	}

	var result Config
	result.PHONE_TYPE_FLAG = "N"
	_, err1 := ini.DecodeFile(configfile, &result)

	if err1 != nil {
		fmt.Println("Config file read error : ", err1)
	}

	return result
}

func InitCenterConfig() {
	realpath, _ := os.Executable()
	dir := filepath.Dir(realpath)
	logDir := filepath.Join(dir, "logs")
	err := createDir(logDir)
	if err != nil {
		log.Fatalf("Failed to ensure log directory: %s", err)
	}
	path := filepath.Join(logDir, "DHNCenter")
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
	BasePath = dir + "/"
	//Client = resty.New()

}

func readCenterConfig() Config {
	realpath, _ := os.Executable()
	dir := filepath.Dir(realpath)
	var configfile = filepath.Join(dir, "config.ini")

	_, err := os.Stat(configfile)
	if err != nil {

		err := createConfig(configfile)
		if err != nil {
			Stdlog.Println("Config file create fail")
		}
		Stdlog.Println("config.ini 생성완료 작성을 해주세요.")

		system_exit("DHNCenter")
		fmt.Println("Config file is missing : ", configfile)
	}

	var result Config
	_, err1 := ini.DecodeFile(configfile, &result)

	if err1 != nil {
		fmt.Println("Config file read error : ", err1)
	}

	return result
}

func createDir(dirName string) error {
	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	return nil
}

func createConfig(dirName string) error {
	fo, err := os.Create(dirName)
	if err != nil {
		return fmt.Errorf("Config file create fail: %w", err)
	}
	configData := []string{
		`# DB 관련`,
		`DB = "DB종류"`,
		`DBURL = "사용자:패스워드@tcp(000.000.000.000:포트번호)/데이터베이스"`,
		``,
		`# CENTER 관련`,
		`CENTER_PORT = "CENTER 포트번호"`,
		``,
		`# SERVER 관련`,
		`KAKAO_SENDLIMIT = 500`,
		`SERVER_PORT = "SERVER 포트번호"`,
		`PROFILE_KEY = "프로필키"`,
		`RESPONSE_METHOD = "push" (알림톡 전송방식 : push, polling)`,
		`PHONE_MSG_FLAG = "YES (2차 발송 플래그(1차 발송이 알림톡, 친구톡 일 시))"`,
		`CHANNEL = "채널명"`,
		`DEBUG = "N" (친구톡 데이터 확인 플래그)`,
		``,
		`# 카카오 API URL`,
		`API_SERVER = "https://bzm-api.kakao.com/"`,
		`CENTER_SERVER = "https://bzm-center.kakao.com/"`,
		`IMAGE_SERVER = "https://bzm-upload-api.kakao.com/"`,
		``,
		`# 나노 메시지`,
		`NANO_IDENTI_CODE = "302190001" (default identification code)`,
		`NANO_TYPE_FLAG = "N" (전화번호 별 분리 작업, 010이 붙은 것과 붙지 않은 것)`,
		``,
		`#추가할 설정 내용 필요에 따라 작성`,
	}

	for _, line := range configData {
		fmt.Fprintln(fo, line)
	}

	return nil
}

func system_exit(service_name string) {
	cmd := exec.Command("systemctl", "stop", service_name)
	if err := cmd.Run(); err != nil {
		Stdlog.Println(service_name+" 서비스 종료 실패:", err)
	} else {
		Stdlog.Println(service_name + " 서비스가 성공적으로 종료되었습니다.")
	}
}
