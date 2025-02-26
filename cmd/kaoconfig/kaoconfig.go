package kaoconfig

import (
	"fmt"
	"log"
	"os"
	"time"
	"os/exec"
	"net/http"
	"crypto/tls"
	"sync/atomic"
	"path/filepath"

	ini "github.com/BurntSushi/toml"
	"github.com/go-resty/resty/v2"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
)

type Config struct {
	DNS              string

	SSL_FLAG         string
	SSL_PORT		 string

	DB               string
	DBURL            string
	SENDLIMIT        int
	REALLIMIT        int

	CENTER_PORT      string
	SERVER_PORT      string
	
	PROFILE_KEY      string
	API_SERVER       string
	CENTER_SERVER    string
	IMAGE_SERVER     string
	RESPONSE_METHOD  string
	CHANNEL          string
	KISA_CODE        string
}

var Conf Config
var Stdlog *log.Logger
var BasePath string
var IsRunning bool = true
var ResultLimit int = 1000
var Client *resty.Client
var RL int32

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
	
	Client = resty.New().
		SetTimeout(100 * time.Second).
		SetTLSClientConfig(&tls.Config{MinVersion: tls.VersionTLS12}).
		SetRetryCount(3).
		SetRetryWaitTime(2 * time.Second).
		SetTransport(&http.Transport{
			MaxIdleConns:		Conf.REALLIMIT,
			MaxIdleConnsPerHost: Conf.REALLIMIT/5,
			IdleConnTimeout:	 90 * time.Second,
		})

	RL = int32(Conf.REALLIMIT)

	go func(rl int32){
		for{
			time.Sleep(1 * time.Second)
			atomic.StoreInt32(&RL, rl)
		}
	}(RL)

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
		`# DNS`,
		`DNS = "주소"`,
		``,
		`# SSL`,
		`SSL_FLAG = "Y" or "N"`,
		`SSL_PORT = "포트"`,
		``,
		`# DB`,
		`DB = "DB종류"`,
		`DBURL = "사용자:패스워드@tcp(000.000.000.000:포트번호)/데이터베이스"`,
		`SENDLIMIT = 그룹핑건수(숫자)`,
		`REALLIMIT = 초당발송건수(숫자)`,
		``,
		`# AGENT`,
		`CENTER_PORT = "센터서버 포트번호"`,
		`SERVER_PORT = "발송서버 포트번호"`,
		``,
		`# KAKAO`,
		`PROFILE_KEY = "프로필키"`,
		`API_SERVER = "https://bzm-api.kakao.com/"`,
		`CENTER_SERVER = "https://bzm-center.kakao.com/"`,
		`IMAGE_SERVER = "https://bzm-upload-api.kakao.com/"`,
		`RESPONSE_METHOD = "push"`,
		`CHANNEL = "채널명"`,
		``,
		`# SMS`,
		`KISA_CODE = "KISA코드"`,
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
