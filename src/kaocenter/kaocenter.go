package kaocenter

import (
	"bytes"
	"encoding/json"

	"fmt"
	"io"
	"io/ioutil"
	config "mycs/src/kaoconfig"
	db "mycs/src/kaodatabasepool"
	"net"
	"net/http"
	"net/url"
	"path/filepath"

	"mime/multipart"
	"os"

	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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

	yellowId := c.Query("yellowId")
	phoneNumber := c.Query("phoneNumber")

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

	categoryCode := c.Query("categoryCode")
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

	token := c.Request.Header.Get("token")
	phoneNumber := c.Request.Header.Get("phoneNumber")

	param := &SenderCreate{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
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
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Sender_(c *gin.Context) {
	conf := config.Conf

	senderKey := c.Query("senderKey")

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
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
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

	param := &SenderDelete{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
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

func Template_Create(c *gin.Context) {
	conf := config.Conf

	param := &TemplateCreate{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	jsonstr, _ := json.Marshal(param)
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
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
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
	// _ = writer.WriteField("templateAd", tc.TemplateAd)
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

func Template_(c *gin.Context) {
	conf := config.Conf

	senderKey := c.Query("senderKey")
	templateCode := c.Query("templateCode")
	senderKeyType := c.Query("senderKeyType")

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
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Template_Request(c *gin.Context) {
	conf := config.Conf

	param := &TemplateRequest{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
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

func Template_Cancel_Request(c *gin.Context) {
	conf := config.Conf

	param := &TemplateRequest{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
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
	
	fmt.Println("T U Call")
	param := &TemplateUpdate{}
	
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	jsonstr, _ := json.Marshal(param)
	fmt.Println("Json : ", string(jsonstr))
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
	bytes, _ := ioutil.ReadAll(resp.Body)
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
	// _ = writer.WriteField("newTemplateAd", tu.NewTemplateAd)
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

	param := &TemplateRequest{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
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

	param := &TemplateRequest{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
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

func Template_Delete(c *gin.Context) {
	conf := config.Conf

	param := &TemplateRequest{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
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
		"attachment":    mustOpen(config.BasePath+"upload/" + newFileName),
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

	categoryCode := c.Query("categoryCode")

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

	param := &TemplateRequest{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
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

func Group_(c *gin.Context) {
	conf := config.Conf

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
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Group_Sender(c *gin.Context) {
	conf := config.Conf

	groupKey := c.Query("groupKey")

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
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Group_Sender_Add(c *gin.Context) {
	conf := config.Conf

	param := &GroupSenderAdd{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
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
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Group_Sender_Remove(c *gin.Context) {
	conf := config.Conf

	param := &GroupSenderAdd{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
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
	bytes, _ := ioutil.ReadAll(resp.Body)
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

	senderKey := c.Query("senderKey")

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/plugin/callbackUrl/list?senderKey="+senderKey, nil)
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
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/plugin/callbackUrl/create", buff)
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
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/plugin/callbackUrl/update", buff)
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
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/plugin/callbackUrl/delete", buff)
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
		"image": mustOpen(config.BasePath+"upload/" +newFileName),
	}

	resp, err := upload(conf.IMAGE_SERVER+"v1/"+conf.PROFILE_KEY+"/image/friendtalk", param)

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

	err = c.SaveUploadedFile(file, config.BasePath+"upload/" + newFileName)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	param := map[string]io.Reader{
		"image": mustOpen(config.BasePath+"upload/" + newFileName),
	}

	resp, err := upload(conf.IMAGE_SERVER+"v1/"+conf.PROFILE_KEY+"/image/friendtalk/wide", param)

	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func AT_Image(c *gin.Context) {
	conf := config.Conf

	file, err := c.FormFile("image")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	extension := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + extension

	err = c.SaveUploadedFile(file, config.BasePath+"upload/" + newFileName)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	param := map[string]io.Reader{
		"image": mustOpen(config.BasePath+"upload/" + newFileName),
	}

	resp, err := upload(conf.IMAGE_SERVER+ "v1/"+conf.PROFILE_KEY+"/image/alimtalk/template", param)

	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func AL_Image(c *gin.Context) {
	conf := config.Conf

	file, err := c.FormFile("image")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	extension := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + extension

	err = c.SaveUploadedFile(file, config.BasePath+"upload/" + newFileName)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	param := map[string]io.Reader{
		"image": mustOpen(config.BasePath+"upload/" + newFileName),
	}

	resp, err := upload(conf.IMAGE_SERVER+ "v1/"+conf.PROFILE_KEY+"/image/alimtalk", param)

	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func MMS_Image(c *gin.Context) {
	//conf := config.Conf
	var newFileName1,newFileName2,newFileName3 string

	userID := c.PostForm("userid")
	file1, err1 := c.FormFile("image1")
	
	var startNow = time.Now()
	var group_no = fmt.Sprintf("%04d%02d%02d%02d%02d%02d%09d", startNow.Year(), startNow.Month(), startNow.Day(), startNow.Hour(), startNow.Minute(), startNow.Second(), startNow.Nanosecond())
						
	if err1 != nil {
		config.Stdlog.Println("File 1 Parameter 오류 : " , err1)
	} else {
		extension1 := filepath.Ext(file1.Filename)
		newFileName1 = config.BasePath+"upload/mms/" + uuid.New().String() + extension1
	
		err := c.SaveUploadedFile(file1, newFileName1)
		if err != nil {
			config.Stdlog.Println("File 1 저장 오류 : ", newFileName1, err)
			newFileName1 = ""
		}
	}

	file2, err2 := c.FormFile("image2")
	
	if err2 != nil {
		config.Stdlog.Println("File 2 Parameter 오류 : " , err2)
	} else {
		extension2 := filepath.Ext(file2.Filename)
		newFileName2 = config.BasePath+"upload/mms/" + uuid.New().String() + extension2
	
		err := c.SaveUploadedFile(file2, newFileName2)
		if err != nil {
			config.Stdlog.Println("File 2 저장 오류 : ", newFileName2, err)
			newFileName2 = ""
		}
	}

	file3, err3 := c.FormFile("image3")
	
	if err3 != nil {
		config.Stdlog.Println("File 3 Parameter 오류 : " , err3)
	} else {
		extension3 := filepath.Ext(file3.Filename)
		newFileName3 = config.BasePath+"upload/mms/" + uuid.New().String() + extension3
	
		err := c.SaveUploadedFile(file3, newFileName3)
		if err != nil {
			config.Stdlog.Println("File 3 저장 오류 : ", newFileName3, err)
			newFileName3 = ""
		}
	}
 
	if len(newFileName1) > 0 || len(newFileName2) > 0 || len(newFileName2) > 0  {
	
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
			"image group":group_no,
		})
	} else {
		c.JSON(http.StatusNoContent, gin.H{
			"message":"Error",
		})
	}
}


func Image_wideItemList(c *gin.Context) {
	conf := config.Conf
	config.Stdlog.Println("Call ")
	
	var newFileName1,newFileName2,newFileName3,newFileName4 string
	
	file1, err1 := c.FormFile("image_1")
	if err1 != nil {
		config.Stdlog.Println(err1.Error())
		c.String(http.StatusBadRequest, fmt.Sprintf("File 1 - get form err: %s", err1.Error()))
		return
	}

	extension := filepath.Ext(file1.Filename)
	newFileName1 = uuid.New().String() + extension

	err1 = c.SaveUploadedFile(file1, config.BasePath+"upload/" + newFileName1)
	if err1 != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("File 1 - get form err: %s", err1.Error()))
		return
	}

	file2, err2 := c.FormFile("image_2")
	if err2 == nil {
		extension := filepath.Ext(file2.Filename)
		newFileName2 = uuid.New().String() + extension
			
		err2 = c.SaveUploadedFile(file2, config.BasePath+"upload/" + newFileName2)
		if err2 != nil {
			newFileName2 = "_"
		}
	} else {
		newFileName2 = "_"
	}

	file3, err3 := c.FormFile("image_3")
	if err3 == nil {
		extension := filepath.Ext(file3.Filename)
		newFileName3 = uuid.New().String() + extension
	
		err3 = c.SaveUploadedFile(file3, config.BasePath+"upload/" + newFileName3)
		if err3 != nil {
			newFileName3 = "_"
		}
	} else {
		newFileName3 = "_"
	}

	file4, err4 := c.FormFile("image_4")
	if err4 == nil {
		extension := filepath.Ext(file4.Filename)
		newFileName4 = uuid.New().String() + extension
	
		err4 = c.SaveUploadedFile(file4, config.BasePath+"upload/" + newFileName4)
		if err4 != nil {
			newFileName4 = "_"
		}
	} else {
		newFileName4 = "_"
	}
	
	param := map[string]io.Reader{
		"image_1": mustOpen(config.BasePath+"upload/" + newFileName1),
		"image_2": mustOpen(config.BasePath+"upload/" + newFileName2),
		"image_3": mustOpen(config.BasePath+"upload/" + newFileName3),
		"image_4": mustOpen(config.BasePath+"upload/" + newFileName4),
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

	resp, err1 := upload(conf.IMAGE_SERVER+ "v1/"+conf.PROFILE_KEY+"/image/friendtalk/wideItemList", param)
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Image_carousel(c *gin.Context) {
	conf := config.Conf
	config.Stdlog.Println("Call ")
	
	var newFileName1,newFileName2,newFileName3,newFileName4,newFileName5,newFileName6 string
	
	file1, err1 := c.FormFile("image_1")
	if err1 != nil {
		config.Stdlog.Println(err1.Error())
		c.String(http.StatusBadRequest, fmt.Sprintf("File 1 - get form err: %s", err1.Error()))
		return
	}

	extension := filepath.Ext(file1.Filename)
	newFileName1 = uuid.New().String() + extension

	err1 = c.SaveUploadedFile(file1, config.BasePath+"upload/" + newFileName1)
	if err1 != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("File 1 - get form err: %s", err1.Error()))
		return
	}

	file2, err2 := c.FormFile("image_2")
	if err2 == nil {
		extension := filepath.Ext(file2.Filename)
		newFileName2 = uuid.New().String() + extension
			
		err2 = c.SaveUploadedFile(file2, config.BasePath+"upload/" + newFileName2)
		if err2 != nil {
			newFileName2 = "_"
		}
	} else {
		newFileName2 = "_"
	}

	file3, err3 := c.FormFile("image_3")
	if err3 == nil {
		extension := filepath.Ext(file3.Filename)
		newFileName3 = uuid.New().String() + extension
	
		err3 = c.SaveUploadedFile(file3, config.BasePath+"upload/" + newFileName3)
		if err3 != nil {
			newFileName3 = "_"
		}
	} else {
		newFileName3 = "_"
	}

	file4, err4 := c.FormFile("image_4")
	if err4 == nil {
		extension := filepath.Ext(file4.Filename)
		newFileName4 = uuid.New().String() + extension
	
		err4 = c.SaveUploadedFile(file4, config.BasePath+"upload/" + newFileName4)
		if err4 != nil {
			newFileName4 = "_"
		}
	} else {
		newFileName4 = "_"
	}
	
	file5, err5 := c.FormFile("image_5")
	if err5 == nil {
		extension := filepath.Ext(file5.Filename)
		newFileName5 = uuid.New().String() + extension
	
		err5 = c.SaveUploadedFile(file5, config.BasePath+"upload/" + newFileName5)
		if err5 != nil {
			newFileName5 = "_"
		}
	} else {
		newFileName5 = "_"
	}
	
	file6, err6 := c.FormFile("image_6")
	if err6 == nil {
		extension := filepath.Ext(file6.Filename)
		newFileName6 = uuid.New().String() + extension
	
		err6 = c.SaveUploadedFile(file6, config.BasePath+"upload/" + newFileName6)
		if err6 != nil {
			newFileName6 = "_"
		}
	} else {
		newFileName6 = "_"
	}
		
	param := map[string]io.Reader{
		"image_1": mustOpen(config.BasePath+"upload/" + newFileName1),
		"image_2": mustOpen(config.BasePath+"upload/" + newFileName2),
		"image_3": mustOpen(config.BasePath+"upload/" + newFileName3),
		"image_4": mustOpen(config.BasePath+"upload/" + newFileName4),
		"image_5": mustOpen(config.BasePath+"upload/" + newFileName5),
		"image_6": mustOpen(config.BasePath+"upload/" + newFileName6),
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

	resp, err1 := upload(conf.IMAGE_SERVER+ "v1/"+conf.PROFILE_KEY+"/image/friendtalk/wideItemList", param)
	bytes, _ := ioutil.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func Get_Polling_Id(c *gin.Context) {
	conf := config.Conf
	respId := c.Param("respid")

	buff := bytes.NewBuffer([]byte("{}"))
	req, err := http.NewRequest("POST", conf.API_SERVER+"/v3/"+conf.PROFILE_KEY+"/response/"+respId, buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err2 := centerClient.Do(req)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, err2.Error())
		return
	}

	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)

}

func AT_Highlight_Image(c *gin.Context) {
	conf := config.Conf

	param, err := image_Seq_Mapping(c, map[string]io.Reader{}, 0, "image")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	// file, err := c.FormFile("image")
	// if err != nil {
	// 	c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
	// 	return
	// }

	// extension := filepath.Ext(file.Filename)
	// newFileName := uuid.New().String() + extension

	// err = c.SaveUploadedFile(file, config.BasePath+"upload/" + newFileName)
	// if err != nil {
	// 	c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
	// 	return
	// }

	// param := map[string]io.Reader{
	// 	"image": mustOpen(config.BasePath+"upload/" + newFileName),
	// }

	resp, err := upload(conf.IMAGE_SERVER+"v1/"+conf.PROFILE_KEY+"/image/alimtalk/itemHighlight", param)
	if err != nil {
		config.Stdlog.Println("File upload 오류 : ", err)
	}

	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func FT_Carousel_Feed_Image(c *gin.Context) {
	conf := config.Conf

	param, err := image_Seq_Mapping(c, map[string]io.Reader{}, 10, "image")
	if err != nil {
		config.Stdlog.Println("image Mapping 오류 : ", err)
	}

	resp, _ := upload(conf.IMAGE_SERVER+"v1/"+conf.PROFILE_KEY+"/image/friendtalk/carousel", param)
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func FT_Carousel_Commerce_Image(c *gin.Context) {
	conf := config.Conf

	param, err := image_Seq_Mapping(c, map[string]io.Reader{}, 11, "image")
	if err != nil {
		config.Stdlog.Println("image Mapping 오류 : ", err)
	}

	resp, _ := upload(conf.IMAGE_SERVER+"v1/"+conf.PROFILE_KEY+"/image/friendtalk/carouselCommerce", param)
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func DM_Default_Image(c *gin.Context) {
	conf := config.Conf

	param, err := image_Seq_Mapping(c, map[string]io.Reader{}, 0, "image")
	if err != nil {
		config.Stdlog.Println("image Mapping 오류 : ", err)
	}

	resp, _ := upload(conf.IMAGE_SERVER+"v2/"+conf.PROFILE_KEY+"/image/default", param)
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func DM_Wide_Image(c *gin.Context) {
	conf := config.Conf

	param, err := image_Seq_Mapping(c, map[string]io.Reader{}, 0, "image")
	if err != nil {
		config.Stdlog.Println("image Mapping 오류 : ", err)
	}

	resp, _ := upload(conf.IMAGE_SERVER+"v2/"+conf.PROFILE_KEY+"/image/wide", param)
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func DM_Widelist_First_image(c *gin.Context) {
	conf := config.Conf

	param, err := image_Seq_Mapping(c, map[string]io.Reader{}, 0, "image")
	if err != nil {
		config.Stdlog.Println("image Mapping 오류 : ", err)
	}

	resp, _ := upload(conf.IMAGE_SERVER+"v2/"+conf.PROFILE_KEY+"/image/wideItemList/first", param)
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func DM_Widelist_Image(c *gin.Context) {
	conf := config.Conf

	param, err := image_Seq_Mapping(c, map[string]io.Reader{}, 3, "image")
	if err != nil {
		config.Stdlog.Println("image Mapping 오류 : ", err)
	}

	resp, _ := upload(conf.IMAGE_SERVER+"v2/"+conf.PROFILE_KEY+"/image/wideItemList", param)
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func DM_Carousel_Feed_Image(c *gin.Context) {
	conf := config.Conf

	param, err := image_Seq_Mapping(c, map[string]io.Reader{}, 10, "image")
	if err != nil {
		config.Stdlog.Println("image Mapping 오류 : ", err)
	}

	resp, _ := upload(conf.IMAGE_SERVER+"v2/"+conf.PROFILE_KEY+"/image/carouselFeed", param)
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func DM_Carousel_Commerce_Image(c *gin.Context) {
	conf := config.Conf

	param, err := image_Seq_Mapping(c, map[string]io.Reader{}, 11, "image")
	if err != nil {
		config.Stdlog.Println("image Mapping 오류 : ", err)
	}

	resp, _ := upload(conf.IMAGE_SERVER+"v2/"+conf.PROFILE_KEY+"/image/carouselCommerce", param)
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

// 친구톡 API
// 별첨1 - 비즈폼 업로드 요청
func Bizform_upload_(c *gin.Context) {
	conf := config.Conf

	param := &Bizform_upload{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/bizform/upload", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)

}

// 별첨2 - 친구톡 발송 가능 모수 확인
func Ft_possible_(c *gin.Context) {
	conf := config.Conf

	param := &Ft_possible{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/friendtalk/possible", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)

}

// 센터 API
// 발신 프로필 조회2 (톡 채널 키로 조회)
func Sender_channel(c *gin.Context) {
	conf := config.Conf

	talkChannelKey := c.Param("talkChannelKey")

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v3/"+conf.PROFILE_KEY+"/sender/"+talkChannelKey, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

// 최근 변경 발신 프로필 조회
func Sender_modified(c *gin.Context) {
	conf := config.Conf

	//since := c.Query("since")
	//page := c.Query("page")
	//count := c.Query("count")

	params := map[string]string{
		"since": c.Query("since"),
		"page":  c.Query("page"),
		"count": c.Query("count"),
	}

	query := c.Request.URL.Query()
	for key, value := range params {
		if value != "" {
			query.Set(key, value)
		}
	}

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v3/"+conf.PROFILE_KEY+"/sender/last_modified?"+query.Encode(), nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

// 검수요청 (파일첨부)
func Template_request_with_file(c *gin.Context) {
	conf := config.Conf

	param := map[string]io.Reader{
		"senderKey":     strings.NewReader(c.PostForm("senderKey")),
		"templateCode":  strings.NewReader(c.PostForm("templateCode")),
		"senderKeyType": strings.NewReader(c.PostForm("senderKeyType")),
		"comment":       strings.NewReader(c.PostForm("comment")),
	}

	file, err := c.FormFile("attachment")
	if err == nil { 
		extension := filepath.Ext(file.Filename)
		newFileName := uuid.New().String() + extension

		err = c.SaveUploadedFile(file, config.BasePath+"upload/"+newFileName)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}
		param["attachment"] = mustOpen(config.BasePath + "upload/" + newFileName)
	}

	resp, err := upload(conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/request_with_file", param)
	if err != nil {
		config.Stdlog.Println("File upload 오류 : ", err)
	}
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

// 검수 승인 취소
func Template_cancel_approval_(c *gin.Context) {
	conf := config.Conf

	param := &TemplateRequest{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/cancel_approval", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)

}

// 기등록된 템플릿 (타입 : BA, EX) 을 채널추가버튼 및 채널추가안내문구가 포함된 템플릿으로 전환 /template/convertAddCh
func Template_convertAddCh_(c *gin.Context) {
	conf := config.Conf

	param := &Template_convertAddCh{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/alimtalk/template/convertAddCh", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)

}

// 채널에 발신 프로필 추가
func Channel_sender_add_(c *gin.Context) {
	conf := config.Conf

	param := &Channel_sender{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/channel/sender/add", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)

}

// 채널에 발신 프로필 삭제
func Channel_sender_remove_(c *gin.Context) {
	conf := config.Conf

	param := &Channel_sender{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/channel/sender/remove", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)

}

// 알림톡, 친구톡 발송 일별 통계
func Stat_daily(c *gin.Context) {
	conf := config.Conf

	//beginDate := c.Query("beginDate")
	//endDate := c.Query("endDate")
	//productType := c.Query("productType")
	//page := c.Query("page")

	params := map[string]string{
		"beginDate":   c.Query("beginDate"),
		"endDate":     c.Query("endDate"),
		"productType": c.Query("productType"),
		"page":        c.Query("page"),
	}

	if !MissingParams(c, params) {
		return
	}

	query := c.Request.URL.Query()
	for key, value := range params {
		if value != "" {
			query.Set(key, value)
		}
	}

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/stat/daily?"+query.Encode(), nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

// 그룹 태그 생성
func GroupTag_create(c *gin.Context) {
	conf := config.Conf

	param := &Group_Tag_create{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/groupTag/create", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)

}

// 그룹 태그 조회
func GroupTag_(c *gin.Context) {
	conf := config.Conf

	//senderKey := c.Query("senderKey")
	//groupTagKey := c.Query("groupTagKey")

	params := map[string]string{
		"senderKey":   c.Query("senderKey"),
		"groupTagKey": c.Query("groupTagKey"),
	}

	if !MissingParams(c, params) {
		return
	}

	query := c.Request.URL.Query()
	for key, value := range params {
		if value != "" {
			query.Set(key, value)
		}
	}

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/groupTag?"+query.Encode(), nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

// 그룹 태그 전체 조회
func GroupTag_list(c *gin.Context) {
	conf := config.Conf

	//senderKey := c.Query("senderKey")

	params := map[string]string{
		"senderKey": c.Query("senderKey"),
	}

	if !MissingParams(c, params) {
		return
	}

	query := c.Request.URL.Query()
	for key, value := range params {
		if value != "" {
			query.Set(key, value)
		}
	}

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/groupTag/list?"+query.Encode(), nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)

}

// 그룹 태그 수정
func GroupTag_update(c *gin.Context) {
	conf := config.Conf
	param := &Group_Tag_update{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/groupTag/update", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)

}

// 그룹 태그 삭제
func GroupTag_delete(c *gin.Context) {
	conf := config.Conf
	param := &Group_Tag_delete{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/groupTag/delete", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)

}

// 광고성 메시지(다이렉트) 템플릿 등록
func Direct_template_create_(c *gin.Context) {
	conf := config.Conf
	param := &Direct_template_create{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v3/"+conf.PROFILE_KEY+"/direct/template/create", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

// 광고성메시지(다이렉트) 템플릿 조회
func Direct_template_(c *gin.Context) {
	conf := config.Conf
	code := c.Param("code")

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/direct/template/"+code, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

// 광고성메시지(다이렉트) 템플릿 수정
func Direct_template_update_(c *gin.Context) {
	conf := config.Conf
	code := c.Param("code")
	param := &Direct_template_create{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v3/"+conf.PROFILE_KEY+"/direct/template/update/"+code, buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

// 광고성메시지(다이렉트) 템플릿 삭제
func Direct_template_delete_(c *gin.Context) {
	conf := config.Conf
	code := c.Param("code")

	buff := bytes.NewBuffer([]byte(`{}`))
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v2/"+conf.PROFILE_KEY+"/direct/template/delete/"+code, buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

// 발신채널 전환
func Direct_convert_(c *gin.Context) {
	conf := config.Conf
	param := &Direct_convert{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/sender/direct/convert", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

// 발신채널 전환 상태 확인
func Direct_convert_result(c *gin.Context) {
	conf := config.Conf

	params := map[string]string{
		"senderKey": c.Query("senderKey"),
	}

	if !MissingParams(c, params) {
		return
	}

	query := c.Request.URL.Query()
	for key, value := range params {
		if value != "" {
			query.Set(key, value)
		}
	}

	req, err := http.NewRequest("GET", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/sender/direct/convert/result?"+query.Encode(), nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

// 발신채널에 연결된 비즈월렛 변경
func Direct_bizWallet_change_(c *gin.Context) {
	conf := config.Conf
	param := &Direct_bizWallet_change{}
	err := c.Bind(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	jsonstr, _ := json.Marshal(param)
	buff := bytes.NewBuffer(jsonstr)
	req, err := http.NewRequest("POST", conf.CENTER_SERVER+"api/v1/"+conf.PROFILE_KEY+"/sender/direct/bizWallet/change", buff)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := centerClient.Do(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	bytes, _ := io.ReadAll(resp.Body)
	c.Data(http.StatusOK, "application/json", bytes)
}

func image_Seq_Mapping(c *gin.Context, param map[string]io.Reader, max int, filename string) (map[string]io.Reader, error) {
	var retErr error = nil
	if max > 0 {
		for a := 1; a <= max; a++ {
			file, err := c.FormFile(filename + "_" + strconv.Itoa(a))
			newFileName := ""
			if err == nil {
				extension := filepath.Ext(file.Filename)
				newFileName = uuid.New().String() + extension
				err2 := c.SaveUploadedFile(file, config.BasePath+"upload/"+newFileName)
				if err2 != nil {
					newFileName = "_"
					retErr = err2
				}
			} else {
				newFileName = "_"
				retErr = err
			}

			if newFileName != "_" {
				param[filename+"_"+strconv.Itoa(a)] = mustOpen(config.BasePath + "upload/" + newFileName)
			}
		}
	} else {
		file, err := c.FormFile(filename)
		newFileName := ""
		if err == nil {
			extension := filepath.Ext(file.Filename)
			newFileName = uuid.New().String() + extension
			err2 := c.SaveUploadedFile(file, config.BasePath+"upload/"+newFileName)
			if err2 != nil {
				newFileName = "_"
				retErr = err2
			}
		} else {
			newFileName = "_"
			retErr = err
		}

		if newFileName != "_" {
			param[filename] = mustOpen(config.BasePath + "upload/" + newFileName)
		}
	}
	return param, retErr
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

func MissingParams(c *gin.Context, params map[string]string) bool {
	var missingParams []string

	for param, value := range params {
		if value == "" {
			missingParams = append(missingParams, param)
		}
	}

	if len(missingParams) > 0 {
		message := "필수값이 부족합니다. ( "
		for i, param := range missingParams {
			if i != 0 {
				message += ", "
			}
			message += param
		}
		message += " )"
		c.JSON(999, gin.H{
			"message": message,
		})
		return false
	}

	return true
}
