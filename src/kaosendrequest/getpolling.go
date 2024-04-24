package kaosendrequest

import (
	"bytes"
	"encoding/json"
	"fmt"
	kakao "mycs/src/kakaojson"
	config "mycs/src/kaoconfig"
	databasepool "mycs/src/kaodatabasepool"
	"mycs/src/kaocommon"

	"io/ioutil"
	"net"
	"net/http"

	"strconv"
	s "strings"
	"sync"

	"github.com/go-resty/resty/v2"
	"context"
	"time"
)

var polprocCnt int

func PollingProc(ctx context.Context) {
	var wg sync.WaitGroup

	for {
			select {
				case <- ctx.Done():
			
			    config.Stdlog.Println("Polling process가 20초 후에 종료 됨.")
			    time.Sleep(20 * time.Second)
			    config.Stdlog.Println("Polling process 종료 완료")
			    return
			default:	
				wg.Add(1)

				go getPollingProcess(&wg)

				wg.Wait()
			}
	}
}

func getPollingProcess(wg *sync.WaitGroup) {

	defer wg.Done()
	var conf = config.Conf
	var stdlog = config.Stdlog

	channel := make(map[string]interface{})
	channel["channel_key"] = conf.CHANNEL

	jsonData, _ := json.Marshal(channel)
	req, err := http.NewRequest("POST", conf.API_SERVER + "v3/" + conf.PROFILE_KEY + "/responseAll", bytes.NewBuffer(jsonData))
	if err != nil {
		config.Stdlog.Println("알림톡(폴링) 결과수신 에러 request 만들기 실패 ", err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("channel_key", conf.CHANNEL)

	resp, err := config.GoClient.Do(req)
	if err != nil {
		// 에러가 발생한 경우 처리
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			// 타임아웃 오류 처리
			config.Stdlog.Println("알림톡(폴링) 결과수신 타임아웃 / error : ", err.Error())
		} else {
			// 기타 오류 처리
			config.Stdlog.Println("알림톡(폴링) 결과수신 실패 / error : ", err.Error())
		}
		return
	}

	bodyData, _ := ioutil.ReadAll(resp.Body)
	statusCode := resp.StatusCode

	resp.Body.Close()

	if statusCode == 200 {
		var kakaoResp kakao.PollingResponse
		json.Unmarshal([]byte(bodyData), &kakaoResp)
        
        atValues := []kaocommon.AtPollingResColumn{}

		for i, _ := range kakaoResp.Response.Success {

			msgid := kakaoResp.Response.Success[i].Serial_number[9:len(kakaoResp.Response.Success[i].Serial_number)]
			stdlog.Println("성공 : " + msgid + " / " + kakaoResp.Response.Success[i].Received_at)

			atValues = append(atValues, kaocommon.AtPollingResColumn{
				Msgid : msgid,
				Type : "S",
			})
		}

		for i, _ := range kakaoResp.Response.Fail {
			msgid := kakaoResp.Response.Fail[i].Serial_number[9:len(kakaoResp.Response.Fail[i].Serial_number)]
			stdlog.Println("실퍠 : " + msgid + " / " + kakaoResp.Response.Fail[i].Received_at)
			
			atValues = append(atValues, kaocommon.AtPollingResColumn{
				Msgid : msgid,
				Type : "F",
			})

		}
		
		if len(atValues) > 0 {
			tx, err := databasepool.DB.Begin()
			if err != nil {
				stdlog.Println("getpolling.go / getPollingProcess / dhn_polling_result / 트랜젝션 초기화 실패 ", err)
			}
			defer tx.Rollback()

			atStmt, err := tx.Prepare("insert into dhn_polling_result values ($1, $2, now())")
			if err != nil {
				stdlog.Println("getpolling.go / getPollingProcess / dhn_polling_result / ftStmt 초기화 실패 ", err)
				return
			}

			for _, data := range atValues {
				_, err := atStmt.Exec(data.Msgid, data.Type)
				if err != nil {
					stdlog.Println("getpolling.go / getPollingProcess / dhn_polling_result / ftStmt personal Exec ", err)
				}
			}
			
			atStmt.Close()

			err = tx.Commit()
			if err != nil {
				stdlog.Println("getpolling.go / getPollingProcess / dhn_polling_result / ftStmt commit ", err)
			}

		}
		
		if kakaoResp.Response_id > 0 {
			req, err := http.NewRequest("POST", conf.API_SERVER + "v3/" + conf.PROFILE_KEY + "/response/" + strconv.Itoa(kakaoResp.Response_id) + "/complete", nil)
			if err != nil {
				stdlog.Println("알림톡(폴링) 결과수신 후처리 에러 ", err.Error())
				return
			}

			_, err := config.GoClient.Do(req)
			if err != nil {
				// 에러가 발생한 경우 처리
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					// 타임아웃 오류 처리
					stdlog.Println("알림톡(폴링) 결과수신 후처리 타임아웃 error : ", err.Error())
				} else {
					// 기타 오류 처리
					stdlog.Println("알림톡(폴링) 결과수신 후처리 실패 error : ", err.Error())
				}
				return
			}
		}
	}
}
