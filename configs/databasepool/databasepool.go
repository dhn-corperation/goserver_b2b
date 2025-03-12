package databasepool

import (
	"log"
	"time"
	"database/sql"

	config "mycs/configs"
	"mycs/internal/structs"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *sql.DB

type DBWrapper struct {
	*gorm.DB
}

func InitDatabase(flag bool) {
	db, err := sql.Open(config.Conf.DB, config.Conf.DBURL)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(1 * time.Minute)

	DB = db

	if flag {
		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn: db, // 기존 database/sql 연결 사용
		}), &gorm.Config{})

		if err != nil {
			config.Stdlog.Println("GORM 초기화 실패 : ", err)
			log.Fatal("GORM 초기화 실패 : ", err)
		} else {
			dbWrapper := &DBWrapper{gormDB}
			dbWrapper.initTables()
			initProcedure(db)
			initEvent(db)
		}
	}
}

func (db *DBWrapper) initTables(){
	config.Stdlog.Println("DHN_CLIENT_LIST 테이블 마이그레이션 시작")
	err := db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").AutoMigrate(structs.ClientList{})
	if err != nil {
		config.Stdlog.Println("DHN_CLIENT_LIST 테이블 생성 실패 err : ", err)
	}

	config.Stdlog.Println("DHN_RECEPTION 테이블 마이그레이션 시작")
	err = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").AutoMigrate(structs.Reception{})
	if err != nil {
		config.Stdlog.Println("DHN_RECEPTION 테이블 생성 실패 err : ", err)
	}

	config.Stdlog.Println("DHN_REQUEST 테이블 마이그레이션 시작")
	err = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").AutoMigrate(structs.Request{})
	if err != nil {
		config.Stdlog.Println("DHN_REQUEST 테이블 생성 실패 err : ", err)
	}

	config.Stdlog.Println("DHN_REQUEST_AT 테이블 마이그레이션 시작")
	err = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").AutoMigrate(structs.RequestAt{})
	if err != nil {
		config.Stdlog.Println("DHN_REQUEST_AT 테이블 생성 실패 err : ", err)
	}

	config.Stdlog.Println("DHN_REQUEST_RESEND 테이블 마이그레이션 시작")
	err = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").AutoMigrate(structs.RequestResend{})
	if err != nil {
		config.Stdlog.Println("DHN_REQUEST_RESEND 테이블 생성 실패 err : ", err)
	}

	config.Stdlog.Println("DHN_REQUEST_AT_RESEND 테이블 마이그레이션 시작")
	err = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").AutoMigrate(structs.RequestAtResend{})
	if err != nil {
		config.Stdlog.Println("DHN_REQUEST_AT_RESEND 테이블 생성 실패 err : ", err)
	}

	config.Stdlog.Println("DHN_RESULT 테이블 마이그레이션 시작")
	err = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").AutoMigrate(structs.Result{})
	if err != nil {
		config.Stdlog.Println("DHN_RESULT 테이블 생성 실패 err : ", err)
	}

	config.Stdlog.Println("DHN_RESULT_TEMP 테이블 마이그레이션 시작")
	err = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").AutoMigrate(structs.ResultTemp{})
	if err != nil {
		config.Stdlog.Println("DHN_RESULT_TEMP 테이블 생성 실패 err : ", err)
	}

	config.Stdlog.Println("DHN_RESULT_BK_TEMP 테이블 마이그레이션 시작")
	err = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").AutoMigrate(structs.ResultBkTemp{})
	if err != nil {
		config.Stdlog.Println("DHN_RESULT_BK_TEMP 테이블 생성 실패 err : ", err)
	}

	config.Stdlog.Println("DHN_RESULT_STA 테이블 마이그레이션 시작")
	err = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").AutoMigrate(structs.ResultSta{})
	if err != nil {
		config.Stdlog.Println("DHN_RESULT_STA 테이블 생성 실패 err : ", err)
	}

	config.Stdlog.Println("SPECIAL_CHARACTER 테이블 마이그레이션 시작")
	err = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").AutoMigrate(structs.SpecialCharacter{})
	if err != nil {
		config.Stdlog.Println("SPECIAL_CHARACTER 테이블 생성 실패 err : ", err)
	} else {
		db.Unscoped().Where("1 = 1").Delete(&structs.SpecialCharacter{})
		sc := []structs.SpecialCharacter{
			{Origin_hex_code: "C2A0", Enabled: "Y"},
			{Origin_hex_code: "E280A8", Enabled: "Y"},
			{Origin_hex_code: "E280A4", Dest_str: "·", Enabled: "Y"},
			{Origin_hex_code: "E29E9F", Dest_str: "→", Enabled: "Y"},
			{Origin_hex_code: "E29C94", Dest_str: "√", Enabled: "Y"},
			{Origin_hex_code: "E280A3", Dest_str: "·", Enabled: "Y"},
			{Origin_hex_code: "EFBBBF", Enabled: "Y"},
		}
		res := db.Create(&sc)
		if res.Error != nil {
			config.Stdlog.Println("SPECIAL_CHARACTER 데이터 삽입 실패 err : ", res.Error)
		}
	}

	config.Stdlog.Println("API_MMS_IMAGES 테이블 마이그레이션 시작")
	err = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").AutoMigrate(structs.ApiMmsImages{})
	if err != nil {
		config.Stdlog.Println("API_MMS_IMAGES 테이블 생성 실패 err : ", err)
	}
}

func initProcedure(db *sql.DB) {
	if !procedureExists(db, "result_back_proc") {
		createStoredProcedure(db, resultBackProc)
	}
	if !procedureExists(db, "remove_ws") {
		createStoredProcedure(db, removeWs)
	}
	if !procedureExists(db, "result_sata_proc") {
		createStoredProcedure(db, resultSataProc)
	}
}

func initEvent(db *sql.DB) {
	if !eventExists(db, "evt_result_back_proc") {
		createEvent(db, evtResultBackProc)
	}
	if !eventExists(db, "evt_result_sata_proc") {
		createEvent(db, evtResultSataProc)
	}
	if !eventExists(db, "evt_remove_reception") {
		createEvent(db, evtRemoveReception)
	}
}

// 프로시저 확인 함수
func procedureExists(db *sql.DB, procedureName string) bool {
	var count int
	query := `SELECT COUNT(*) FROM information_schema.ROUTINES WHERE ROUTINE_NAME = ? AND (ROUTINE_TYPE = 'PROCEDURE' or ROUTINE_TYPE = 'FUNCTION') AND ROUTINE_SCHEMA = DATABASE()`
	err := db.QueryRow(query, procedureName).Scan(&count)
	if err != nil {
		config.Stdlog.Println("프로시저 확인 실패:", err)
		return false
	}
	return count > 0
}

// 이벤트 확인 함수
func eventExists(db *sql.DB, eventName string) bool {
	var count int
	query := `SELECT COUNT(*) FROM information_schema.EVENTS WHERE EVENT_NAME = ? AND EVENT_SCHEMA = DATABASE()`
	err := db.QueryRow(query, eventName).Scan(&count)
	if err != nil {
		config.Stdlog.Println("이벤트 확인 실패:", err)
		return false
	}
	return count > 0
}

// 프로시저 생성 함수
func createStoredProcedure(db *sql.DB, ddl string) {
	_, err := db.Exec(ddl)
	if err != nil {
		config.Stdlog.Println("프로시저 생성 실패:", err)
	} else {
		config.Stdlog.Println("프로시저 생성 완료!")
	}
}

// 이벤트 생성 함수
func createEvent(db *sql.DB, ddl string) {
	_, err := db.Exec(ddl)
	if err != nil {
		config.Stdlog.Println("이벤트 생성 실패:", err)
	} else {
		config.Stdlog.Println("이벤트 생성 완료!")
	}
}