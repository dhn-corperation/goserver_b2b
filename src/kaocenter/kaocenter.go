package kaocenter

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	config "kaoconfig"
	db "kaodatabasepool"
	"net"
	"net/http"
	"net/url"
	"path/filepath"

	"mime/multipart"
	"os"

	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CheckCenterUser(user CenterUser) bool {
	sqlstr := "select count(1) as cnt from CENTERAPI_USER where user_id = '" + user.BizId + "' and user_key ='" + user.ApiKey + "' and use_yn = 'Y'"
	val, verr := db.DB.Query(sqlstr)
	if verr != nil {
		return false
	}
	defer val.Close()

	var cnt int
	val.Next()
	val.Scan(&cnt)

	if cnt > 0 {
		return true
	} else {
		return false
	}
}

var centerClient *http.Client = &http.Client{
	Timeout: time.Second * 30,
	Transport: &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	},
}

func Sender_token(c *gin.Context) {
	conf := config.Conf

	param := &ProfileInfo{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user := CenterUser{BizId: param.BizId, ApiKey: param.ApiKey}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}

	yellowId := param.YellowId
	phoneNumber := param.PhoneNumber

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/sender/token?yellowId="+yellowId+"&phoneNumber="+phoneNumber, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Category_all(c *gin.Context) {
	conf := config.Conf

	param := &ProfileInfo{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user := CenterUser{BizId: param.BizId, ApiKey: param.ApiKey}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/category/all", nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Category_(c *gin.Context) {
	conf := config.Conf

	param := &ProfileInfo{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user := CenterUser{BizId: param.BizId, ApiKey: param.ApiKey}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}

	categoryCode := param.CategoryCode
	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/category?categoryCode="+categoryCode, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Sender_Create(c *gin.Context) {
	conf := config.Conf

	param := &ProfileInfo{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user := CenterUser{BizId: param.BizId, ApiKey: param.ApiKey}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}

	token := param.Token
	phoneNumber := param.PhoneNumber

	param.Token = ""
	param.PhoneNumber = ""
	param.BizId = ""
	param.ApiKey = ""

	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v3/"+conf.PROFILE_KEY+"/sender/create", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("token", token)
	req.Header.Add("phoneNumber", phoneNumber)

	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	rr := &CenterResponse{}

	bytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bytes, rr)

	if strings.EqualFold(rr.Code, "200") { // 정상적인 요청 이라면

		temp := &ProfileResponse{}
		json.Unmarshal(bytes, temp)

		ProfileTable_Update(temp, param.BizId)
	}

	c.Data(http.StatusOK, "application/json", bytes)
}

func Sender_(c *gin.Context) {
	conf := config.Conf

	param := &ProfileInfo{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user := CenterUser{BizId: param.BizId, ApiKey: param.ApiKey}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}
	senderKey := param.SenderKey

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v3/"+conf.PROFILE_KEY+"/sender?senderKey="+senderKey, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	rr := &CenterResponse{}

	bytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bytes, rr)

	if strings.EqualFold(rr.Code, "200") { // 정상적인 요청 이라면

		temp := &ProfileResponse{}
		json.Unmarshal(bytes, temp)

		ProfileTable_Update(temp, param.BizId)
	}

	c.Data(http.StatusOK, "application/json", bytes)
}

func Sender_Use(c *gin.Context) {
	//conf := config.Conf

	param := &ProfileInfo{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user := CenterUser{BizId: param.BizId, ApiKey: param.ApiKey}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}

	prostr := `select senderkey
					,uuid
					,name
					,status
					,(case when block='0' then 'N' else 'Y' end) as block
					,(case when dormant='0' then 'N' else 'Y' end) as dormant
					,profilestatus
					,createdat
					,modifiedat
					,categorycode
					,(case when alimtalk='0' then 'N' else 'Y' end) as alimtalk
					,(case when bizchat='0' then 'N' else 'Y' end) as bizchat
					,(case when brandtalk='0' then 'N' else 'Y' end) as brandtalk
					,commitalcompanyname
					,channelkey
					,(case when businessprofile='0' then 'N' else 'Y' end) as businessprofile
					,businesstype 
					 from CENTER_PROFILE where bizId = '` + user.BizId + `'`
	pfrows, verr := db.DB.Query(prostr)
	if verr != nil {
		c.JSON(http.StatusBadRequest, verr.Error())
		return
	}
	defer pfrows.Close()

	ProfileRows := []interface{}{}

	var senderkey, uuid, name, status, block, dormant, profilestatus, createdat, modifiedat, categorycode, alimtalk, bizchat, brandtalk, commitalcompanyname, channelkey, businessprofile, businesstype sql.NullString
	var groupKey, gname, gcreatedat sql.NullString
	for pfrows.Next() {
		var profile ProfileInfoSenderS
		pfrows.Scan(&senderkey, &uuid, &name, &status, &block, &dormant, &profilestatus, &createdat, &modifiedat, &categorycode, &alimtalk, &bizchat, &brandtalk, &commitalcompanyname, &channelkey, &businessprofile, &businesstype)

		profile.SenderKey = senderkey.String
		profile.UUID = uuid.String
		profile.Name = name.String
		profile.Status = status.String
		profile.Block = block.String
		profile.Dormant = dormant.String
		profile.ProfileStatus = profilestatus.String
		profile.CreatedAt = createdat.String
		profile.ModifiedAt = modifiedat.String
		profile.CategoryCode = categorycode.String
		profile.Alimtalk = alimtalk.String
		profile.Bizchat = bizchat.String
		profile.Brandtalk = brandtalk.String
		profile.CommitalCompanyName = commitalcompanyname.String
		profile.ChannelKey = channelkey.String
		profile.BusinessProfile = businessprofile.String
		profile.BusinessType = businesstype.String

		groupstr := `select cpg.groupKey 
					,cpg.name
					,cpg.createdAt 
					from CENTER_PROFILE_GROUP_JOIN cpgj 
					inner join CENTER_PROFILE_GROUP cpg 
					on cpgj.groupKey = cpg.groupKey  
					where cpgj.senderKey  = '` + senderkey.String + `'`
		grows, gerr := db.DB.Query(groupstr)
		if gerr != nil {
			c.JSON(http.StatusBadRequest, gerr.Error())
			return
		}
		defer grows.Close()

		var GroupRows []GroupJson

		for grows.Next() {
			var grow GroupJson
			grows.Scan(&groupKey, &gname, &gcreatedat)

			grow.GroupKey = groupKey.String
			grow.Name = gname.String
			grow.CreatedAt = gcreatedat.String

			GroupRows = append(GroupRows, grow)
		}
		profile.Groups = GroupRows
		ProfileRows = append(ProfileRows, profile)
	}

	c.JSON(http.StatusOK, ProfileRows)
	//c.Data(http.StatusOK, "application/json", ProfileRows)
}

func Sender_Multi(c *gin.Context) {
	//conf := config.Conf

	param := &ProfileInfoMulti{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user := CenterUser{BizId: param.BizId, ApiKey: param.ApiKey}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}

	if len(param.SenderKey) > 0 {
		ProfileRows := []interface{}{}
		for i := range param.SenderKey {
			prostr := `select senderkey
					,uuid
					,name
					,status
					,(case when block='0' then 'N' else 'Y' end) as block
					,(case when dormant='0' then 'N' else 'Y' end) as dormant
					,profilestatus
					,createdat
					,modifiedat
					,categorycode
					,(case when alimtalk='0' then 'N' else 'Y' end) as alimtalk
					,(case when bizchat='0' then 'N' else 'Y' end) as bizchat
					,(case when brandtalk='0' then 'N' else 'Y' end) as brandtalk
					,commitalcompanyname
					,channelkey
					,(case when businessprofile='0' then 'N' else 'Y' end) as businessprofile
					,businesstype 
					 from CENTER_PROFILE where bizId = '` + user.BizId + `'
					and senderkey = '` + param.SenderKey[i] + `'`
			pfrows, verr := db.DB.Query(prostr)
			if verr != nil {
				c.JSON(http.StatusBadRequest, verr.Error())
				return
			}
			defer pfrows.Close()

			var senderkey, uuid, name, status, block, dormant, profilestatus, createdat, modifiedat, categorycode, alimtalk, bizchat, brandtalk, commitalcompanyname, channelkey, businessprofile, businesstype sql.NullString
			var groupKey, gname, gcreatedat sql.NullString
			for pfrows.Next() {
				var profile ProfileInfoSenderS
				pfrows.Scan(&senderkey, &uuid, &name, &status, &block, &dormant, &profilestatus, &createdat, &modifiedat, &categorycode, &alimtalk, &bizchat, &brandtalk, &commitalcompanyname, &channelkey, &businessprofile, &businesstype)

				profile.SenderKey = senderkey.String
				profile.UUID = uuid.String
				profile.Name = name.String
				profile.Status = status.String
				profile.Block = block.String
				profile.Dormant = dormant.String
				profile.ProfileStatus = profilestatus.String
				profile.CreatedAt = createdat.String
				profile.ModifiedAt = modifiedat.String
				profile.CategoryCode = categorycode.String
				profile.Alimtalk = alimtalk.String
				profile.Bizchat = bizchat.String
				profile.Brandtalk = brandtalk.String
				profile.CommitalCompanyName = commitalcompanyname.String
				profile.ChannelKey = channelkey.String
				profile.BusinessProfile = businessprofile.String
				profile.BusinessType = businesstype.String

				groupstr := `select cpg.groupKey 
					,cpg.name
					,cpg.createdAt 
					from CENTER_PROFILE_GROUP_JOIN cpgj 
					inner join CENTER_PROFILE_GROUP cpg 
					on cpgj.groupKey = cpg.groupKey  
					where cpgj.senderKey  = '` + senderkey.String + `'`
				grows, gerr := db.DB.Query(groupstr)
				if gerr != nil {
					c.JSON(http.StatusBadRequest, gerr.Error())
					return
				}
				defer grows.Close()

				var GroupRows []GroupJson

				for grows.Next() {
					var grow GroupJson
					grows.Scan(&groupKey, &gname, &gcreatedat)

					grow.GroupKey = groupKey.String
					grow.Name = gname.String
					grow.CreatedAt = gcreatedat.String

					GroupRows = append(GroupRows, grow)
				}
				profile.Groups = GroupRows
				ProfileRows = append(ProfileRows, profile)
			}

		}
		c.JSON(http.StatusOK, ProfileRows)
	} else {
		c.JSON(http.StatusOK, "[]")
	}
	//c.Data(http.StatusOK, "application/json", ProfileRows)
}

func Sender_Delete(c *gin.Context) {
	conf := config.Conf

	param := &SenderDelete{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/sender/delete", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Sender_Recover(c *gin.Context) {
	conf := config.Conf

	param := &ProfileInfo{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user := CenterUser{BizId: param.BizId, ApiKey: param.ApiKey}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}

	param.ApiKey = ""
	param.BizId = ""
	param.CategoryCode = ""
	param.ChannelKey = ""
	param.PhoneNumber = ""
	param.Token = ""
	param.YellowId = ""

	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/sender/recover", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func ProfileTable_Update(temp *ProfileResponse, bizId string) {
	sqlstr := "select count(1) as cnt from CENTER_PROFILE where bizId = '" + bizId + "' and senderKey = '" + temp.Profile.SenderKey + "'"
	//fmt.Println("DB Error", sqlstr)
	val, verr := db.DB.Query(sqlstr)
	if verr != nil {
		fmt.Println("DB Error", sqlstr)
	}
	defer val.Close()

	var cnt int
	val.Next()
	val.Scan(&cnt)

	if cnt == 0 {
		// Template 저장 하기
		insstr := `INSERT INTO CENTER_PROFILE (bizId,senderKey,uuid,name,status,block,dormant,profileStatus,createdAt,modifiedAt,categoryCode,alimtalk,bizchat,brandtalk,commitalCompanyName,channelKey,businessProfile,businessType,profileSpamLevel,profileMessageSpamLevel,clearBlockUrl) 
					VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`
		insValues := []interface{}{}
		insValues = append(insValues, bizId)
		insValues = append(insValues, temp.Profile.SenderKey)
		insValues = append(insValues, temp.Profile.UUID)
		insValues = append(insValues, temp.Profile.Name)
		insValues = append(insValues, temp.Profile.Status)
		insValues = append(insValues, temp.Profile.Block)
		insValues = append(insValues, temp.Profile.Dormant)
		insValues = append(insValues, temp.Profile.ProfileStatus)
		insValues = append(insValues, temp.Profile.CreatedAt)
		insValues = append(insValues, temp.Profile.ModifiedAt)
		insValues = append(insValues, temp.Profile.CategoryCode)
		insValues = append(insValues, temp.Profile.Alimtalk)
		insValues = append(insValues, temp.Profile.Bizchat)
		insValues = append(insValues, temp.Profile.Brandtalk)
		insValues = append(insValues, temp.Profile.CommitalCompanyName)
		insValues = append(insValues, temp.Profile.ChannelKey)
		insValues = append(insValues, temp.Profile.BusinessProfile)
		insValues = append(insValues, temp.Profile.BusinessType)
		insValues = append(insValues, temp.Profile.ProfileSpamLevel)
		insValues = append(insValues, temp.Profile.ProfileMessageSpamLevel)
		insValues = append(insValues, temp.Profile.ClearBlockURL)

		_, dberr := db.DB.Exec(insstr, insValues...)

		if dberr != nil {
			fmt.Println("Template Insert 처리 중 오류 발생 " + dberr.Error())
		}
	} else {
		// Template 저장 하기
		updstr := `update CENTER_PROFILE
						  set uuid = ? ,	 
name = ? ,	 
status = ? ,	 
block = ? ,	 
dormant = ? ,	 
profileStatus = ? ,	 
createdAt = ? ,	 
modifiedAt = ? ,	 
categoryCode = ? ,	 
alimtalk = ? ,	 
bizchat = ? , 
brandtalk = ? ,	 
commitalCompanyName = ? ,	 
channelKey = ? ,	 
businessProfile = ? ,	 
businessType = ? ,	 
profileSpamLevel = ? ,	 
profileMessageSpamLevel = ? ,	 
clearBlockUrl = ?  
					    where bizId = ?
					      and senderKey = ? `
		updValues := []interface{}{}
		updValues = append(updValues, temp.Profile.UUID)
		updValues = append(updValues, temp.Profile.Name)
		updValues = append(updValues, temp.Profile.Status)
		updValues = append(updValues, temp.Profile.Block)
		updValues = append(updValues, temp.Profile.Dormant)
		updValues = append(updValues, temp.Profile.ProfileStatus)
		updValues = append(updValues, temp.Profile.CreatedAt)
		updValues = append(updValues, temp.Profile.ModifiedAt)
		updValues = append(updValues, temp.Profile.CategoryCode)
		updValues = append(updValues, temp.Profile.Alimtalk)
		updValues = append(updValues, temp.Profile.Bizchat)
		updValues = append(updValues, temp.Profile.Brandtalk)
		updValues = append(updValues, temp.Profile.CommitalCompanyName)
		updValues = append(updValues, temp.Profile.ChannelKey)
		updValues = append(updValues, temp.Profile.BusinessProfile)
		updValues = append(updValues, temp.Profile.BusinessType)
		updValues = append(updValues, temp.Profile.ProfileSpamLevel)
		updValues = append(updValues, temp.Profile.ProfileMessageSpamLevel)
		updValues = append(updValues, temp.Profile.ClearBlockURL)
		updValues = append(updValues, bizId)
		updValues = append(updValues, temp.Profile.SenderKey)

		_, dberr := db.DB.Exec(updstr, updValues...)

		if dberr != nil {
			fmt.Println("Template Insert 처리 중 오류 발생 " + dberr.Error())
		}
	}
}

func TemplateTable_Update(temp *TemplateResponse, bizId string) {
	sqlstr := "select count(1) as cnt from CENTER_TEMPLATE where bizId = '" + bizId + "' and senderKey ='" + temp.Template.SenderKey + "' and templateCode = '" + temp.Template.TemplateCode + "'"
	fmt.Println("DB Error", sqlstr)
	val, verr := db.DB.Query(sqlstr)
	if verr != nil {
		fmt.Println("DB Error", sqlstr)
	}
	defer val.Close()

	var cnt int
	val.Next()
	val.Scan(&cnt)

	if cnt == 0 {
		// Template 저장 하기
		insstr := `insert into CENTER_TEMPLATE(bizId, senderKey, senderKeyType, templateStatus, templateCode,templateName,categoryCode,createdAt,modifiedAt)
			                               values(?, ?, ?, ?, ?, ? , ?, ?, ?)`
		insValues := []interface{}{}
		insValues = append(insValues, bizId)
		insValues = append(insValues, temp.Template.SenderKey)
		insValues = append(insValues, temp.Template.SenderKeyType)
		insValues = append(insValues, temp.Template.InspectionStatus)
		insValues = append(insValues, temp.Template.TemplateCode)
		insValues = append(insValues, temp.Template.TemplateName)
		insValues = append(insValues, temp.Template.CategoryCode)
		insValues = append(insValues, temp.Template.CreatedAt)
		insValues = append(insValues, temp.Template.ModifiedAt)

		_, dberr := db.DB.Exec(insstr, insValues...)

		if dberr != nil {
			fmt.Println("Template Insert 처리 중 오류 발생 " + dberr.Error())
		}
	} else {
		// Template 저장 하기
		updstr := `update CENTER_TEMPLATE
						  set senderKeyType = ?, 
					  	      templateStatus =?, 
						      templateName =?,
						      categoryCode =?,
						      createdAt= ?,
						      modifiedAt =?
					    where bizId = ?
					      and senderKey = ?
						  and templateCode =?`
		updValues := []interface{}{}
		updValues = append(updValues, temp.Template.SenderKeyType)
		updValues = append(updValues, temp.Template.InspectionStatus)
		updValues = append(updValues, temp.Template.TemplateName)
		updValues = append(updValues, temp.Template.CategoryCode)
		updValues = append(updValues, temp.Template.CreatedAt)
		updValues = append(updValues, temp.Template.ModifiedAt)
		updValues = append(updValues, bizId)
		updValues = append(updValues, temp.Template.SenderKey)
		updValues = append(updValues, temp.Template.TemplateCode)

		_, dberr := db.DB.Exec(updstr, updValues...)

		if dberr != nil {
			fmt.Println("Template Insert 처리 중 오류 발생 " + dberr.Error())
		}
	}
}

func Template_Create(c *gin.Context) {
	conf := config.Conf

	param := &TemplateCreateJson{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	user := CenterUser{BizId: param.BizId, ApiKey: param.ApiKey}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
	} else {
		// param.ApiKey = ""
		// param.BizId = ""
		jsonstr, _ := json.Marshal(param)
		//fmt.Println(string(jsonstr))
		buff := bytes.NewBuffer(jsonstr)
		req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/create", buff)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		req.Header.Add("Content-Type", "application/json")
		//client := &http.Client{}
		resp, err := centerClient.Do(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		rr := &CenterResponse{}

		bytes, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(bytes, rr)

		if strings.EqualFold(rr.Code, "200") { // 정상적인 요청 이라면

			temp := &TemplateResponse{}
			json.Unmarshal(bytes, temp)

			TemplateTable_Update(temp, param.BizId)

			// 요청 결과 Return
			c.Data(http.StatusOK, "application/json", bytes)
		} else {

			c.Data(http.StatusForbidden, "application/json", bytes)
		}
	}
}

func Template_Code_Check(c *gin.Context) {

	param := &TemplateInfo{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user := CenterUser{BizId: param.BizId, ApiKey: param.ApiKey}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}

	_chk, _ := regexp.Compile(`[a-zA-Z0-9\_\-]+`)
	_chkstr := _chk.ReplaceAllString(param.TemplateCode, "")
	if len(_chkstr) > 0 {
		var respc = &CenterResponse{}
		respc.Code = "0100"
		respc.Message = "허용되지 않은 문자 사용"
		c.JSON(http.StatusOK, respc)
		return
	}

	if len(param.TemplateCode) > 30 || len(param.TemplateCode) < 1 {
		var respc = &CenterResponse{}
		respc.Code = "0200"
		respc.Message = "Template Code 길이 오류( 1 ~ 30자 )"
		c.JSON(http.StatusOK, respc)
		return
	}

	sqlstr := "select count(1) as cnt from CENTER_TEMPLATE where bizId = '" + param.BizId + "' and senderKey ='" + param.SenderKey + "' and templateCode = '" + param.TemplateCode + "'"
	//fmt.Print(sqlstr)
	val, verr := db.DB.Query(sqlstr)
	if verr != nil {
		var respc = &CenterResponse{}
		respc.Code = "0999"
		respc.Message = "기타 오류"
		c.JSON(http.StatusOK, respc)
		return
	}
	defer val.Close()

	var cnt int
	val.Next()
	val.Scan(&cnt)

	if cnt > 0 {
		var respc = &CenterResponse{}
		respc.Code = "0300"
		respc.Message = "이미 사용중인 Template Code"
		c.JSON(http.StatusOK, respc)
	} else {
		var respc = &CenterResponse{}
		respc.Code = "0000"
		respc.Message = ""
		c.JSON(http.StatusOK, respc)
	}
}

func Template_Create_Image(c *gin.Context) {
	conf := config.Conf

	file, err := c.FormFile("image")
	if err != nil {
		//fmt.Println(err.Error())
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	extension := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + extension

	err = c.SaveUploadedFile(file, config.BasePath+"upload/"+newFileName)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	var tc TemplateCreate
	json.Unmarshal([]byte(c.PostForm("json")), &tc)

	cfile, _ := os.Open(config.BasePath + "upload/" + newFileName)
	defer cfile.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, perr := writer.CreateFormFile("image", filepath.Base(config.BasePath+"upload/"+newFileName))
	if perr != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", perr.Error()))
		return
	}

	_, err = io.Copy(part, cfile)

	_ = writer.WriteField("senderKey", tc.SenderKey)
	_ = writer.WriteField("templateCode", tc.TemplateCode)
	_ = writer.WriteField("templateName", tc.TemplateName)
	_ = writer.WriteField("templateContent", tc.TemplateContent)
	_ = writer.WriteField("templateMessageType", tc.TemplateMessageType)
	_ = writer.WriteField("templateExtra", tc.TemplateExtra)
	_ = writer.WriteField("templateAd", tc.TemplateAd)
	_ = writer.WriteField("templateEmphasizeType", tc.TemplateEmphasizeType)
	_ = writer.WriteField("senderKeyType", tc.SenderKeyType)
	_ = writer.WriteField("categoryCode", tc.CategoryCode)
	_ = writer.WriteField("securityFlag", strconv.FormatBool(tc.SecurityFlag))

	for key, r := range tc.Buttons {
		_ = writer.WriteField("buttons["+strconv.Itoa(key)+"].name", r.Name)
		_ = writer.WriteField("buttons["+strconv.Itoa(key)+"].linkType", r.LinkType)
		_ = writer.WriteField("buttons["+strconv.Itoa(key)+"].ordering", strconv.Itoa(r.Ordering))
		_ = writer.WriteField("buttons["+strconv.Itoa(key)+"].linkMo", r.LinkMo)
		_ = writer.WriteField("buttons["+strconv.Itoa(key)+"].linkPc", r.LinkPc)
		_ = writer.WriteField("buttons["+strconv.Itoa(key)+"].linkAnd", r.LinkAnd)
		_ = writer.WriteField("buttons["+strconv.Itoa(key)+"].linkIos", r.LinkIos)
	}

	for key, r := range tc.QuickReplies {
		_ = writer.WriteField("quickReplies["+strconv.Itoa(key)+"].name", r.Name)
		_ = writer.WriteField("quickReplies["+strconv.Itoa(key)+"].linkType", r.LinkType)
		_ = writer.WriteField("quickReplies["+strconv.Itoa(key)+"].linkMo", r.LinkMo)
		_ = writer.WriteField("quickReplies["+strconv.Itoa(key)+"].linkPc", r.LinkPc)
		_ = writer.WriteField("quickReplies["+strconv.Itoa(key)+"].linkAnd", r.LinkAnd)
		_ = writer.WriteField("quickReplies["+strconv.Itoa(key)+"].linkIos", r.LinkIos)
	}

	err = writer.Close()

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/create_with_image", body)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func TemplateList(c *gin.Context) {

	//conf := config.Conf

	param := &TemplateInfo{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user := CenterUser{BizId: param.BizId, ApiKey: param.ApiKey}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}
	tmpstr := `select senderKey,
	senderKeyType,
	templateCode,
	templateName,
	categoryCode,
	createdAt,
	modifiedAt,
	templateStatus
from
	CENTER_TEMPLATE ct
where
	ct.senderKey = '` + param.SenderKey + `'
and ct.bizId = '` + param.BizId + `'`

	if len(param.SenderKeyType) > 0 {
		tmpstr = tmpstr + ` and ct.senderKeyType = '` + param.SenderKeyType + `'`
	}
	if len(param.TemplateStatus) > 0 {
		tmpstr = tmpstr + ` and ct.templateStatus = '` + param.TemplateStatus + `'`
	}
	if len(param.Keyword) > 0 {
		tmpstr = tmpstr + ` and ct.templateName like '%` + param.Keyword + `%'`
	}
	if len(param.StartDate) > 0 {
		tmpstr = tmpstr + ` and regexp_replace( ct.createdAt, '[^0-9]', '') >= '` + param.StartDate + `'`
	}
	if len(param.EndDate) > 0 {
		tmpstr = tmpstr + ` and regexp_replace( ct.createdAt, '[^0-9]', '') <= '` + param.EndDate + `'`
	}

	trows, terr := db.DB.Query(tmpstr)

	if terr != nil {
		c.JSON(http.StatusBadRequest, terr.Error())
		return
	}

	defer trows.Close()

	var senderKey, senderKeyType, templateCode, templateName, categoryCode, createdAt, modifiedAt, templateStatus sql.NullString

	var tmplist []TemplateListJson
	var templist TemplateList_

	for trows.Next() {
		trows.Scan(&senderKey, &senderKeyType, &templateCode, &templateName, &categoryCode, &createdAt, &modifiedAt, &templateStatus)
		var t TemplateListJson
		t.SenderKey = senderKey.String
		t.SenderKeyType = senderKeyType.String
		t.TemplateCode = templateCode.String
		t.TemplateName = templateName.String
		t.CategoryCode = categoryCode.String
		t.CreatedAt = createdAt.String
		t.ModifiedAt = modifiedAt.String
		t.Status = templateStatus.String

		tmplist = append(tmplist, t)
	}
	if len(tmplist) > 0 {
		templist.TotalCount = len(tmplist)
		templist.TotalPage = 1
		templist.CurrentPage = 1
		templist.Data = tmplist

	} else {
		templist.Message = "조건에 맞는 Template 이 없습니다."
	}
	c.JSON(http.StatusOK, templist)
}

func Template_(c *gin.Context) {

	conf := config.Conf

	param := &TemplateInfo{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user := CenterUser{BizId: param.BizId, ApiKey: param.ApiKey}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}

	senderKey := param.SenderKey
	templateCode := param.TemplateCode
	senderKeyType := param.SenderKeyType

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template?senderKey="+senderKey+"&templateCode="+url.QueryEscape(templateCode)+"&senderKeyType="+senderKeyType, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	req.Header.Add("Accept-Charset", "utf-8")

	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	rr := &CenterResponse{}

	bytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bytes, rr)

	if strings.EqualFold(rr.Code, "200") { // Template 조회가 정상이면 DB 내역에 수정

		temp := &TemplateResponse{}
		json.Unmarshal(bytes, temp)

		TemplateTable_Update(temp, param.BizId)

	}

	c.Data(http.StatusOK, "application/json", bytes)
}

func Template_Request(c *gin.Context) {
	conf := config.Conf

	param := &TemplateInfo{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user := CenterUser{BizId: param.BizId, ApiKey: param.ApiKey}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}
	param.ApiKey = ""
	param.BizId = ""

	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/request", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Template_Request_WF(c *gin.Context) {
	conf := config.Conf

	user := CenterUser{BizId: c.PostForm("bizId"), ApiKey: c.PostForm("apiKey")}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}

	file, err := c.FormFile("attachment")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	extension := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + extension

	err = c.SaveUploadedFile(file, config.BasePath+"upload/"+newFileName)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	param := map[string]io.Reader{
		"attachment":    mustOpen(config.BasePath + "upload/" + newFileName),
		"senderKey":     strings.NewReader(c.PostForm("senderKey")),
		"templateCode":  strings.NewReader(c.PostForm("templateCode")),
		"senderKeyType": strings.NewReader(c.PostForm("senderKeyType")),
		"comment":       strings.NewReader(c.PostForm("comment")),
	}

	resp, err := upload(conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/request_with_file", param)

	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Template_Cancel_Request(c *gin.Context) {
	conf := config.Conf

	param := &TemplateInfo{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//user := CenterUser{BizId: param.BizId, ApiKey: param.Message}
	user := CenterUser{BizId: param.BizId, ApiKey: param.ApiKey}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}

	param.BizId = ""
	param.ApiKey = ""

	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/cancel_request", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Template_Update(c *gin.Context) {
	conf := config.Conf

	param := &TemplateUpdateJson{}

	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user := CenterUser{BizId: param.BizId, ApiKey: param.ApiKey}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}

	newP := &TemplateUpdateKakao{}
	newP.SenderKey = param.SenderKey
	newP.SenderKeyType = param.SenderKeyType
	newP.NewSenderKey = param.SenderKey
	newP.TemplateCode = param.TemplateCode
	newP.NewTemplateCode = param.NewTemplateCode
	newP.NewTemplateName = param.TemplateName
	newP.NewTemplateMessageType = param.TemplateMessageType
	newP.NewTemplateEmphasizeType = param.TemplateEmphasizeType
	newP.NewTemplateContent = param.TemplateContent
	newP.NewTemplateExtra = param.TemplateExtra
	newP.NewTemplateImageName = param.TemplateImageName
	newP.NewTemplateImageURL = param.TemplateImageURL
	newP.NewTemplateTitle = param.TemplateTitle
	newP.NewTemplateSubtitle = param.TemplateSubtitle
	newP.NewTemplateHeader = param.TemplateHeader
	newP.NewTemplateItemHighlight = param.TemplateItemHighlight
	newP.NewTemplateItem = param.TemplateItem
	newP.NewCategoryCode = param.CategoryCode
	newP.SecurityFlag = param.SecurityFlag
	newP.Buttons = param.Buttons
	newP.QuickReplies = param.QuickReplies

	jsonstr, _ := json.Marshal(newP)
	//fmt.Println("Json : ", string(jsonstr))
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/update", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	rr := &CenterResponse{}

	bytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bytes, rr)

	if strings.EqualFold(rr.Code, "200") { // Template 조회가 정상이면 DB 내역에 수정

		temp := &TemplateResponse{}
		json.Unmarshal(bytes, temp)

		TemplateTable_Update(temp, param.BizId)
	}

	c.Data(http.StatusOK, "application/json", bytes)

}

func Template_Update_Image(c *gin.Context) {
	conf := config.Conf

	file, err := c.FormFile("image")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	extension := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + extension

	err = c.SaveUploadedFile(file, config.BasePath+"upload/"+newFileName)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	var tu TemplateUpdate
	json.Unmarshal([]byte(c.PostForm("json")), &tu)

	cfile, _ := os.Open(config.BasePath + "upload/" + newFileName)
	defer cfile.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, perr := writer.CreateFormFile("image", filepath.Base(config.BasePath+"upload/"+newFileName))
	if perr != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", perr.Error()))
		return
	}

	_, err = io.Copy(part, cfile)

	_ = writer.WriteField("senderKey", tu.SenderKey)
	_ = writer.WriteField("templateCode", tu.TemplateCode)
	_ = writer.WriteField("senderKeyType", tu.SenderKeyType)
	_ = writer.WriteField("newSenderKey", tu.NewSenderKey)
	_ = writer.WriteField("newTemplateCode", tu.NewTemplateCode)
	_ = writer.WriteField("newTemplateName", tu.NewTemplateName)
	_ = writer.WriteField("newTemplateContent", tu.NewTemplateContent)
	_ = writer.WriteField("newTemplateMessageType", tu.NewTemplateMessageType)
	_ = writer.WriteField("newTemplateExtra", tu.NewTemplateExtra)
	_ = writer.WriteField("newTemplateAd", tu.NewTemplateAd)
	_ = writer.WriteField("newTemplateEmphasizeType", tu.NewTemplateEmphasizeType)
	_ = writer.WriteField("newSenderKeyType", tu.NewSenderKeyType)
	_ = writer.WriteField("newCategoryCode", tu.NewCategoryCode)
	_ = writer.WriteField("securityFlag", strconv.FormatBool(tu.SecurityFlag))

	for key, r := range tu.Buttons {
		_ = writer.WriteField("buttons["+strconv.Itoa(key)+"].name", r.Name)
		_ = writer.WriteField("buttons["+strconv.Itoa(key)+"].linkType", r.LinkType)
		_ = writer.WriteField("buttons["+strconv.Itoa(key)+"].ordering", strconv.Itoa(r.Ordering))
		_ = writer.WriteField("buttons["+strconv.Itoa(key)+"].linkMo", r.LinkMo)
		_ = writer.WriteField("buttons["+strconv.Itoa(key)+"].linkPc", r.LinkPc)
		_ = writer.WriteField("buttons["+strconv.Itoa(key)+"].linkAnd", r.LinkAnd)
		_ = writer.WriteField("buttons["+strconv.Itoa(key)+"].linkIos", r.LinkIos)
	}

	for key, r := range tu.QuickReplies {
		_ = writer.WriteField("quickReplies["+strconv.Itoa(key)+"].name", r.Name)
		_ = writer.WriteField("quickReplies["+strconv.Itoa(key)+"].linkType", r.LinkType)
		_ = writer.WriteField("quickReplies["+strconv.Itoa(key)+"].linkMo", r.LinkMo)
		_ = writer.WriteField("quickReplies["+strconv.Itoa(key)+"].linkPc", r.LinkPc)
		_ = writer.WriteField("quickReplies["+strconv.Itoa(key)+"].linkAnd", r.LinkAnd)
		_ = writer.WriteField("quickReplies["+strconv.Itoa(key)+"].linkIos", r.LinkIos)
	}

	err = writer.Close()

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/update_with_image", body)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Template_Stop(c *gin.Context) {
	conf := config.Conf

	param := &TemplateInfo{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//user := CenterUser{BizId: param.BizId, ApiKey: param.Message}
	user := CenterUser{BizId: param.BizId, ApiKey: param.ApiKey}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}

	param.BizId = ""
	param.ApiKey = ""

	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/stop", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Template_Reuse(c *gin.Context) {
	conf := config.Conf

	param := &TemplateInfo{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//user := CenterUser{BizId: param.BizId, ApiKey: param.Message}
	user := CenterUser{BizId: param.BizId, ApiKey: param.ApiKey}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}

	param.BizId = ""
	param.ApiKey = ""

	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/reuse", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Template_CA(c *gin.Context) {
	conf := config.Conf

	param := &TemplateInfo{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//user := CenterUser{BizId: param.BizId, ApiKey: param.Message}
	user := CenterUser{BizId: param.BizId, ApiKey: param.ApiKey}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}

	param.BizId = ""
	param.ApiKey = ""

	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/cancel_approval", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Template_Delete(c *gin.Context) {
	conf := config.Conf

	param := &TemplateInfo{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user := CenterUser{BizId: param.BizId, ApiKey: param.ApiKey}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}

	param.ApiKey = ""
	param.BizId = ""

	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/delete", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Template_Last_Modified(c *gin.Context) {
	conf := config.Conf

	senderKey := c.Query("senderKey")
	senderKeyType := c.Query("senderKeyType")
	since := c.Query("since")
	page := c.Query("page")
	count := c.Query("count")

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/last_modified?senderKey="+senderKey+"&senderKeyType="+senderKeyType+"&since="+since+"&page="+page+"&count="+count, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Template_Comment(c *gin.Context) {
	conf := config.Conf

	param := &TemplateComment{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/comment", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Template_Comment_File(c *gin.Context) {
	conf := config.Conf

	file, err := c.FormFile("attachment")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	extension := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + extension

	err = c.SaveUploadedFile(file, config.BasePath+"upload/"+newFileName)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	param := map[string]io.Reader{
		"attachment":    mustOpen(config.BasePath + "upload/" + newFileName),
		"senderKey":     strings.NewReader(c.PostForm("senderKey")),
		"templateCode":  strings.NewReader(c.PostForm("templateCode")),
		"senderKeyType": strings.NewReader(c.PostForm("senderKeyType")),
		"comment":       strings.NewReader(c.PostForm("comment")),
	}

	resp, err := upload(conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/comment_file", param)

	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Template_Category_all(c *gin.Context) {
	conf := config.Conf

	param := &TemplateInfo{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//user := CenterUser{BizId: param.BizId, ApiKey: param.Message}
	user := CenterUser{BizId: param.BizId, ApiKey: param.ApiKey}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/category/all", nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Template_Category_(c *gin.Context) {
	conf := config.Conf

	param := &TemplateInfo{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user := CenterUser{BizId: param.BizId, ApiKey: param.ApiKey}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}

	categoryCode := param.CategoryCode

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/category?categoryCode="+categoryCode, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Template_Category_Update(c *gin.Context) {
	conf := config.Conf

	param := &TemplateCategoryUpdate{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/category/update", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Template_Dormant_Release(c *gin.Context) {
	conf := config.Conf

	param := &TemplateInfo{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//user := CenterUser{BizId: param.BizId, ApiKey: param.Message}
	user := CenterUser{BizId: param.BizId, ApiKey: param.ApiKey}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}

	param.BizId = ""
	param.ApiKey = ""

	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/dormant/release", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Template_Convert_AC(c *gin.Context) {
	conf := config.Conf

	param := &TemplateInfo{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	//user := CenterUser{BizId: param.BizId, ApiKey: param.Message}
	user := CenterUser{BizId: param.BizId, ApiKey: param.ApiKey}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}

	param.BizId = ""
	param.ApiKey = ""

	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/convertAddCh", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func GroupTable_Update(temp *GroupResponse, bizId string) {

	for i, _ := range temp.Group {

		sqlstr := "select count(1) as cnt from CENTER_PROFILE_GROUP where bizId = '" + bizId + "' and groupKey = '" + temp.Group[i].GroupKey + "'"
		//fmt.Println("DB Error", sqlstr)
		val, verr := db.DB.Query(sqlstr)
		if verr != nil {
			fmt.Println("DB Error", sqlstr)
		}
		defer val.Close()

		var cnt int
		val.Next()
		val.Scan(&cnt)

		if cnt == 0 {
			// Template 저장 하기
			insstr := `INSERT INTO CENTER_PROFILE_GROUP (bizId,groupKey, name, createdAt ) 
					VALUES(?,?,?,? );`
			insValues := []interface{}{}
			insValues = append(insValues, bizId)
			insValues = append(insValues, temp.Group[i].GroupKey)
			insValues = append(insValues, temp.Group[i].Name)
			insValues = append(insValues, temp.Group[i].CreatedAt)

			_, dberr := db.DB.Exec(insstr, insValues...)

			if dberr != nil {
				fmt.Println("Template Insert 처리 중 오류 발생 " + dberr.Error())
			}
		} else {
			// Template 저장 하기
			updstr := `update CENTER_PROFILE_GROUP
						  set  name = ? ,	  
createdAt = ?  
					    where bizId = ?
					      and groupKey = ? `
			updValues := []interface{}{}
			updValues = append(updValues, temp.Group[i].Name)
			updValues = append(updValues, temp.Group[i].CreatedAt)
			updValues = append(updValues, bizId)
			updValues = append(updValues, temp.Group[i].GroupKey)

			_, dberr := db.DB.Exec(updstr, updValues...)

			if dberr != nil {
				fmt.Println("Template Insert 처리 중 오류 발생 " + dberr.Error())
			}
		}
	}
}

func GroupProfile_Update(temp *GroupProfileResponse, bizId string, groupKey string) {

	for i, _ := range temp.Profile {
		insstr := `insert ignore into CENTER_PROFILE_GROUP_JOIN(bizId, senderKey, groupKey) values(?,?, ?)`

		insValues := []interface{}{}
		insValues = append(insValues, bizId)
		insValues = append(insValues, temp.Profile[i].SenderKey)
		insValues = append(insValues, groupKey)

		_, dberr := db.DB.Exec(insstr, insValues...)

		if dberr != nil {
			fmt.Println("Template Insert 처리 중 오류 발생 " + dberr.Error())
		}
	}
}

func Group_(c *gin.Context) {
	conf := config.Conf

	param := &GroupInfo{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user := CenterUser{BizId: param.BizId, ApiKey: param.ApiKey}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/group", nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	rr := &CenterResponse{}

	bytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bytes, rr)

	if strings.EqualFold(rr.Code, "200") { // 정상적인 요청 이라면

		temp := &GroupResponse{}
		json.Unmarshal(bytes, temp)

		go GroupTable_Update(temp, param.BizId)
	}
	c.Data(http.StatusOK, "application/json", bytes)
}

func Group_Sender(c *gin.Context) {
	conf := config.Conf

	param := &GroupInfo{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user := CenterUser{BizId: param.BizId, ApiKey: param.ApiKey}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}

	groupKey := param.GroupKey

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v3/"+conf.PROFILE_KEY+"/group/sender?groupKey="+groupKey, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	rr := &CenterResponse{}

	bytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bytes, rr)

	if strings.EqualFold(rr.Code, "200") { // 정상적인 요청 이라면

		temp := &GroupProfileResponse{}
		json.Unmarshal(bytes, temp)

		go GroupProfile_Update(temp, param.BizId, param.GroupKey)
	}
	c.Data(http.StatusOK, "application/json", bytes)
}

func Group_Sender_Add(c *gin.Context) {
	conf := config.Conf

	param := &GroupInfo{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user := CenterUser{BizId: param.BizId, ApiKey: param.ApiKey}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}
	//param.ApiKey = ""
	//param.BizId = ""

	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/group/sender/add", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	rr := &CenterResponse{}

	bytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bytes, rr)

	if strings.EqualFold(rr.Code, "200") { // 정상적인 요청 이라면
		insstr := `insert ignore into CENTER_PROFILE_GROUP_JOIN(bizId, senderKey, groupKey) values(?,?, ?)`

		insValues := []interface{}{}
		insValues = append(insValues, param.BizId)
		insValues = append(insValues, param.SenderKey)
		insValues = append(insValues, param.GroupKey)

		db.DB.Exec(insstr, insValues...)
	}
	c.Data(http.StatusOK, "application/json", bytes)
}

func Group_Sender_Remove(c *gin.Context) {
	conf := config.Conf

	param := &GroupInfo{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user := CenterUser{BizId: param.BizId, ApiKey: param.ApiKey}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}

	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/group/sender/remove", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	rr := &CenterResponse{}

	bytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bytes, rr)

	if strings.EqualFold(rr.Code, "200") { // 정상적인 요청 이라면
		delstr := `delete from CENTER_PROFILE_GROUP_JOIN where bizId = ? and senderKey = ? and groupKey = ?`

		delValues := []interface{}{}
		delValues = append(delValues, param.BizId)
		delValues = append(delValues, param.SenderKey)
		delValues = append(delValues, param.GroupKey)

		db.DB.Exec(delstr, delValues...)
	}
	c.Data(http.StatusOK, "application/json", bytes)
}

func Channel_Create_(c *gin.Context) {
	conf := config.Conf

	param := &ChannelCreate{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/channel/create", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Channel_all(c *gin.Context) {
	conf := config.Conf

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/channel/all", nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Channel_(c *gin.Context) {
	conf := config.Conf

	channelKey := c.Query("channelKey")

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/channel?channelKey="+channelKey, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Channel_Update_(c *gin.Context) {
	conf := config.Conf

	param := &ChannelCreate{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/channel/update", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Channel_Senders_(c *gin.Context) {
	conf := config.Conf

	param := &ChannelSenders{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/channel/senders", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Channel_Delete_(c *gin.Context) {
	conf := config.Conf

	param := &ChannelDelete{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/channel/delete", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Plugin_CallbackUrls_List(c *gin.Context) {
	conf := config.Conf

	param := &PluginCallbacnUrlList{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/plugin/callbackUrl/list?senderKey="+param.SenderKey, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Plugin_callbackUrl_Create(c *gin.Context) {
	conf := config.Conf

	param := &PluginCallbackUrlCreate{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/plugin/callbackUrl/create", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Plugin_callbackUrl_Update(c *gin.Context) {
	conf := config.Conf

	param := &PluginCallbackUrlUpdate{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/plugin/callbackUrl/update", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Plugin_callbackUrl_Delete(c *gin.Context) {
	conf := config.Conf

	param := &PluginCallbackUrlDelete{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/plugin/callbackUrl/delete", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	//client := &http.Client{}
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func FT_Upload(c *gin.Context) {
	conf := config.Conf

	user := CenterUser{BizId: c.PostForm("bizId"), ApiKey: c.PostForm("apiKey")}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}
	imgType := c.PostForm("imageType")
	file, err := c.FormFile("image")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	extension := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + extension

	err = c.SaveUploadedFile(file, config.BasePath+"upload/"+newFileName)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	param := map[string]io.Reader{
		"image": mustOpen(config.BasePath + "upload/" + newFileName),
	}

	var resp *http.Response
	var errup error

	if strings.EqualFold(imgType, "W") {
		resp, errup = upload(conf.IMAGE_SERVER+"v1/"+conf.PROFILE_KEY+"/image/friendtalk/wide", param)
	} else {
		resp, errup = upload(conf.IMAGE_SERVER+"v1/"+conf.PROFILE_KEY+"/image/friendtalk", param)
	}

	if errup != nil {
		c.JSON(http.StatusBadRequest, " 기타오류")
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func FT_Wide_Upload(c *gin.Context) {
	conf := config.Conf

	file, err := c.FormFile("image")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	extension := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + extension

	err = c.SaveUploadedFile(file, config.BasePath+"upload/"+newFileName)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	param := map[string]io.Reader{
		"image": mustOpen(config.BasePath + "upload/" + newFileName),
	}

	resp, err := upload(conf.IMAGE_SERVER+"v1/"+conf.PROFILE_KEY+"/image/friendtalk/wide", param)

	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func AT_Image(c *gin.Context) {
	conf := config.Conf

	user := CenterUser{BizId: c.PostForm("bizId"), ApiKey: c.PostForm("apiKey")}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	extension := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + extension

	err = c.SaveUploadedFile(file, config.BasePath+"upload/"+newFileName)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	param := map[string]io.Reader{
		"image": mustOpen(config.BasePath + "upload/" + newFileName),
	}

	resp, err := upload(conf.IMAGE_SERVER+"v1/"+conf.PROFILE_KEY+"/image/alimtalk/template", param)

	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func AL_Image(c *gin.Context) {
	conf := config.Conf

	user := CenterUser{BizId: c.PostForm("bizId"), ApiKey: c.PostForm("apiKey")}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	extension := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + extension

	err = c.SaveUploadedFile(file, config.BasePath+"upload/"+newFileName)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	param := map[string]io.Reader{
		"image": mustOpen(config.BasePath + "upload/" + newFileName),
	}

	resp, err := upload(conf.IMAGE_SERVER+"v1/"+conf.PROFILE_KEY+"/image/alimtalk", param)

	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func AL_Image_HL(c *gin.Context) {
	conf := config.Conf

	user := CenterUser{BizId: c.PostForm("bizId"), ApiKey: c.PostForm("apiKey")}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	extension := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + extension

	err = c.SaveUploadedFile(file, config.BasePath+"upload/"+newFileName)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	param := map[string]io.Reader{
		"image": mustOpen(config.BasePath + "upload/" + newFileName),
	}

	resp, err := upload(conf.IMAGE_SERVER+"v1/"+conf.PROFILE_KEY+"/image/alimtalk/itemHighlight", param)

	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func MMS_Image(c *gin.Context) {
	//conf := config.Conf
	var newFileName1, newFileName2, newFileName3 string

	userID := c.PostForm("userid")
	file1, err1 := c.FormFile("image1")

	var startNow = time.Now()
	var group_no = fmt.Sprintf("%04d%02d%02d%02d%02d%02d%09d", startNow.Year(), startNow.Month(), startNow.Day(), startNow.Hour(), startNow.Minute(), startNow.Second(), startNow.Nanosecond())

	if err1 != nil {
		config.Stdlog.Println("File 1 Parameter 오류 : ", err1)
	} else {
		extension1 := filepath.Ext(file1.Filename)
		newFileName1 = config.BasePath + "upload/mms/" + uuid.New().String() + extension1

		err := c.SaveUploadedFile(file1, newFileName1)
		if err != nil {
			config.Stdlog.Println("File 1 저장 오류 : ", newFileName1, err)
			newFileName1 = ""
		}
	}

	file2, err2 := c.FormFile("image2")

	if err2 != nil {
		config.Stdlog.Println("File 2 Parameter 오류 : ", err2)
	} else {
		extension2 := filepath.Ext(file2.Filename)
		newFileName2 = config.BasePath + "upload/mms/" + uuid.New().String() + extension2

		err := c.SaveUploadedFile(file2, newFileName2)
		if err != nil {
			config.Stdlog.Println("File 2 저장 오류 : ", newFileName2, err)
			newFileName2 = ""
		}
	}

	file3, err3 := c.FormFile("image3")

	if err3 != nil {
		config.Stdlog.Println("File 3 Parameter 오류 : ", err3)
	} else {
		extension3 := filepath.Ext(file3.Filename)
		newFileName3 = config.BasePath + "upload/mms/" + uuid.New().String() + extension3

		err := c.SaveUploadedFile(file3, newFileName3)
		if err != nil {
			config.Stdlog.Println("File 3 저장 오류 : ", newFileName3, err)
			newFileName3 = ""
		}
	}

	if len(newFileName1) > 0 || len(newFileName2) > 0 || len(newFileName2) > 0 {

		mmsinsQuery := `insert IGNORE into api_mms_images(
  user_id,
  mms_id,             
  origin1_path,
  origin2_path,
  origin3_path,
  file1_path,
  file2_path,
  file3_path   
) values %s
	`
		mmsinsStrs := []string{}
		mmsinsValues := []interface{}{}

		mmsinsStrs = append(mmsinsStrs, "(?,?,null,null,null,?,?,?)")
		mmsinsValues = append(mmsinsValues, userID)
		mmsinsValues = append(mmsinsValues, group_no)
		mmsinsValues = append(mmsinsValues, newFileName1)
		mmsinsValues = append(mmsinsValues, newFileName2)
		mmsinsValues = append(mmsinsValues, newFileName3)

		if len(mmsinsStrs) >= 1 {
			stmt := fmt.Sprintf(mmsinsQuery, strings.Join(mmsinsStrs, ","))
			_, err := db.DB.Exec(stmt, mmsinsValues...)

			if err != nil {
				config.Stdlog.Println("API MMS Insert 처리 중 오류 발생 " + err.Error())
			}

			mmsinsStrs = nil
			mmsinsValues = nil
		}

		c.JSON(http.StatusOK, gin.H{
			"image group": group_no,
		})
	} else {
		c.JSON(http.StatusNoContent, gin.H{
			"message": "Error",
		})
	}
}

func Image_wideItemList(c *gin.Context) {
	conf := config.Conf

	user := CenterUser{BizId: c.PostForm("bizId"), ApiKey: c.PostForm("apiKey")}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}

	var newFileName1, newFileName2, newFileName3, newFileName4 string

	file1, err1 := c.FormFile("imageList[0].image")
	if err1 != nil {
		config.Stdlog.Println(err1.Error())
		c.String(http.StatusBadRequest, fmt.Sprintf("File 1 - get form err: %s", err1.Error()))
		return
	}

	extension := filepath.Ext(file1.Filename)
	newFileName1 = uuid.New().String() + extension

	err1 = c.SaveUploadedFile(file1, config.BasePath+"upload/"+newFileName1)
	if err1 != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("File 1 - get form err: %s", err1.Error()))
		return
	}

	file2, err2 := c.FormFile("imageList[1].image")
	if err2 == nil {
		extension := filepath.Ext(file2.Filename)
		newFileName2 = uuid.New().String() + extension

		err2 = c.SaveUploadedFile(file2, config.BasePath+"upload/"+newFileName2)
		if err2 != nil {
			newFileName2 = "_"
		}
	} else {
		newFileName2 = "_"
	}

	file3, err3 := c.FormFile("imageList[2].image")
	if err3 == nil {
		extension := filepath.Ext(file3.Filename)
		newFileName3 = uuid.New().String() + extension

		err3 = c.SaveUploadedFile(file3, config.BasePath+"upload/"+newFileName3)
		if err3 != nil {
			newFileName3 = "_"
		}
	} else {
		newFileName3 = "_"
	}

	file4, err4 := c.FormFile("imageList[3].image")
	if err4 == nil {
		extension := filepath.Ext(file4.Filename)
		newFileName4 = uuid.New().String() + extension

		err4 = c.SaveUploadedFile(file4, config.BasePath+"upload/"+newFileName4)
		if err4 != nil {
			newFileName4 = "_"
		}
	} else {
		newFileName4 = "_"
	}

	param := map[string]io.Reader{
		"image_1": mustOpen(config.BasePath + "upload/" + newFileName1),
		"image_2": mustOpen(config.BasePath + "upload/" + newFileName2),
		"image_3": mustOpen(config.BasePath + "upload/" + newFileName3),
		"image_4": mustOpen(config.BasePath + "upload/" + newFileName4),
	}

	if newFileName4 == "_" {
		delete(param, "image_4")
	}

	if newFileName3 == "_" {
		delete(param, "image_3")
	}

	if newFileName2 == "_" {
		delete(param, "image_2")
	}

	resp, err1 := upload(conf.IMAGE_SERVER+"v1/"+conf.PROFILE_KEY+"/image/friendtalk/wideItemList", param)
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Image_carousel(c *gin.Context) {
	conf := config.Conf

	user := CenterUser{BizId: c.PostForm("bizId"), ApiKey: c.PostForm("apiKey")}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}

	var newFileName1, newFileName2, newFileName3, newFileName4, newFileName5, newFileName6 string

	file1, err1 := c.FormFile("imageList[0].image")
	if err1 != nil {
		config.Stdlog.Println(err1.Error())
		c.String(http.StatusBadRequest, fmt.Sprintf("File 1 - get form err: %s", err1.Error()))
		return
	}

	extension := filepath.Ext(file1.Filename)
	newFileName1 = uuid.New().String() + extension

	err1 = c.SaveUploadedFile(file1, config.BasePath+"upload/"+newFileName1)
	if err1 != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("File 1 - get form err: %s", err1.Error()))
		return
	}

	file2, err2 := c.FormFile("imageList[1].image")
	if err2 == nil {
		extension := filepath.Ext(file2.Filename)
		newFileName2 = uuid.New().String() + extension

		err2 = c.SaveUploadedFile(file2, config.BasePath+"upload/"+newFileName2)
		if err2 != nil {
			newFileName2 = "_"
		}
	} else {
		newFileName2 = "_"
	}

	file3, err3 := c.FormFile("imageList[2].image")
	if err3 == nil {
		extension := filepath.Ext(file3.Filename)
		newFileName3 = uuid.New().String() + extension

		err3 = c.SaveUploadedFile(file3, config.BasePath+"upload/"+newFileName3)
		if err3 != nil {
			newFileName3 = "_"
		}
	} else {
		newFileName3 = "_"
	}

	file4, err4 := c.FormFile("imageList[3].image")
	if err4 == nil {
		extension := filepath.Ext(file4.Filename)
		newFileName4 = uuid.New().String() + extension

		err4 = c.SaveUploadedFile(file4, config.BasePath+"upload/"+newFileName4)
		if err4 != nil {
			newFileName4 = "_"
		}
	} else {
		newFileName4 = "_"
	}

	file5, err5 := c.FormFile("imageList[4].image")
	if err5 == nil {
		extension := filepath.Ext(file5.Filename)
		newFileName5 = uuid.New().String() + extension

		err5 = c.SaveUploadedFile(file5, config.BasePath+"upload/"+newFileName5)
		if err5 != nil {
			newFileName5 = "_"
		}
	} else {
		newFileName5 = "_"
	}

	file6, err6 := c.FormFile("imageList[5].image")
	if err6 == nil {
		extension := filepath.Ext(file6.Filename)
		newFileName6 = uuid.New().String() + extension

		err6 = c.SaveUploadedFile(file6, config.BasePath+"upload/"+newFileName6)
		if err6 != nil {
			newFileName6 = "_"
		}
	} else {
		newFileName6 = "_"
	}

	param := map[string]io.Reader{
		"image_1": mustOpen(config.BasePath + "upload/" + newFileName1),
		"image_2": mustOpen(config.BasePath + "upload/" + newFileName2),
		"image_3": mustOpen(config.BasePath + "upload/" + newFileName3),
		"image_4": mustOpen(config.BasePath + "upload/" + newFileName4),
		"image_5": mustOpen(config.BasePath + "upload/" + newFileName5),
		"image_6": mustOpen(config.BasePath + "upload/" + newFileName6),
	}
	if newFileName6 == "_" {
		delete(param, "image_6")
	}

	if newFileName5 == "_" {
		delete(param, "image_5")
	}

	if newFileName4 == "_" {
		delete(param, "image_4")
	}

	if newFileName3 == "_" {
		delete(param, "image_3")
	}

	if newFileName2 == "_" {
		delete(param, "image_2")
	}

	resp, err1 := upload(conf.IMAGE_SERVER+"v1/"+conf.PROFILE_KEY+"/image/friendtalk/carousel", param)
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Image_carousel_C(c *gin.Context) {
	conf := config.Conf

	user := CenterUser{BizId: c.PostForm("bizId"), ApiKey: c.PostForm("apiKey")}
	if !CheckCenterUser(user) {
		c.JSON(http.StatusForbidden, "접근 권한이 없습니다.")
		return
	}

	var newFileName1, newFileName2, newFileName3, newFileName4, newFileName5, newFileName6 string

	file1, err1 := c.FormFile("imageList[0].image")
	if err1 != nil {
		config.Stdlog.Println(err1.Error())
		c.String(http.StatusBadRequest, fmt.Sprintf("File 1 - get form err: %s", err1.Error()))
		return
	}

	extension := filepath.Ext(file1.Filename)
	newFileName1 = uuid.New().String() + extension

	err1 = c.SaveUploadedFile(file1, config.BasePath+"upload/"+newFileName1)
	if err1 != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("File 1 - get form err: %s", err1.Error()))
		return
	}

	file2, err2 := c.FormFile("imageList[1].image")
	if err2 == nil {
		extension := filepath.Ext(file2.Filename)
		newFileName2 = uuid.New().String() + extension

		err2 = c.SaveUploadedFile(file2, config.BasePath+"upload/"+newFileName2)
		if err2 != nil {
			newFileName2 = "_"
		}
	} else {
		newFileName2 = "_"
	}

	file3, err3 := c.FormFile("imageList[2].image")
	if err3 == nil {
		extension := filepath.Ext(file3.Filename)
		newFileName3 = uuid.New().String() + extension

		err3 = c.SaveUploadedFile(file3, config.BasePath+"upload/"+newFileName3)
		if err3 != nil {
			newFileName3 = "_"
		}
	} else {
		newFileName3 = "_"
	}

	file4, err4 := c.FormFile("imageList[3].image")
	if err4 == nil {
		extension := filepath.Ext(file4.Filename)
		newFileName4 = uuid.New().String() + extension

		err4 = c.SaveUploadedFile(file4, config.BasePath+"upload/"+newFileName4)
		if err4 != nil {
			newFileName4 = "_"
		}
	} else {
		newFileName4 = "_"
	}

	file5, err5 := c.FormFile("imageList[4].image")
	if err5 == nil {
		extension := filepath.Ext(file5.Filename)
		newFileName5 = uuid.New().String() + extension

		err5 = c.SaveUploadedFile(file5, config.BasePath+"upload/"+newFileName5)
		if err5 != nil {
			newFileName5 = "_"
		}
	} else {
		newFileName5 = "_"
	}

	file6, err6 := c.FormFile("imageList[5].image")
	if err6 == nil {
		extension := filepath.Ext(file6.Filename)
		newFileName6 = uuid.New().String() + extension

		err6 = c.SaveUploadedFile(file6, config.BasePath+"upload/"+newFileName6)
		if err6 != nil {
			newFileName6 = "_"
		}
	} else {
		newFileName6 = "_"
	}

	param := map[string]io.Reader{
		"image_1": mustOpen(config.BasePath + "upload/" + newFileName1),
		"image_2": mustOpen(config.BasePath + "upload/" + newFileName2),
		"image_3": mustOpen(config.BasePath + "upload/" + newFileName3),
		"image_4": mustOpen(config.BasePath + "upload/" + newFileName4),
		"image_5": mustOpen(config.BasePath + "upload/" + newFileName5),
		"image_6": mustOpen(config.BasePath + "upload/" + newFileName6),
	}
	if newFileName6 == "_" {
		delete(param, "image_6")
	}

	if newFileName5 == "_" {
		delete(param, "image_5")
	}

	if newFileName4 == "_" {
		delete(param, "image_4")
	}

	if newFileName3 == "_" {
		delete(param, "image_3")
	}

	if newFileName2 == "_" {
		delete(param, "image_2")
	}

	resp, err1 := upload(conf.IMAGE_SERVER+"v1/"+conf.PROFILE_KEY+"/image/friendtalk/carousel", param)
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func upload(url string, values map[string]io.Reader) (*http.Response, error) {

	var buff bytes.Buffer
	w := multipart.NewWriter(&buff)

	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}

		if x, ok := r.(*os.File); ok {
			fw, _ = w.CreateFormFile(key, x.Name())
		} else {

			fw, _ = w.CreateFormField(key)
		}
		_, err := io.Copy(fw, r)

		if err != nil {
			return nil, err
		}

	}

	w.Close()

	req, err := http.NewRequest("POST", url, &buff)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	//client := &http.Client{}

	resp, err := centerClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func mustOpen(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		//pwd, _ := os.Getwd()
		//fmt.Println("PWD: ", pwd)
		return nil
	}
	return r
}
