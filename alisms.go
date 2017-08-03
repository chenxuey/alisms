package alisms

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	Client         *http.Request
	EndPoint       string
	AccessId       string
	AccessKey      string
	MessageRequest *MessageRequest
	Date           string
	TopicName      string
	Sign           string
	MsgRequest     MsgRequest
}

func NewClient(endPoint, accessId, accessKey, topic string) *Client {
	var client Client
	client.AccessId = accessId
	client.EndPoint = endPoint
	client.AccessKey = accessKey
	client.GetTopicRef(topic)
	return &client
}

func (c *Client) GetTopicRef(topic string) *Client {
	c.TopicName = topic
	return c
}

func (c *Client) getDate() string {
	return c.Date
}

func (c *Client) SetDate() {
	tl, _ := time.LoadLocation("GMT")
	c.Date = time.Now().In(tl).Format(GMT_DATE_FORMAT)
}

func (c *Client) a(topic string) {

}
func (c *Client) setTopic() string {
	return "/topics/" + c.TopicName + "/messages"
}

// 生成签名
func (c *Client) SetAttribute() *Client {
	c.SetDate()
	signStr := "POST\n\n" + "text/xml" + "\n" + c.Date + "\n" + MNS_VERSION_HEADER + ":" + MNS_VERSION + "\n" + c.setTopic()
	mac := hmac.New(sha1.New, []byte(c.AccessKey))
	mac.Write([]byte(signStr))
	c.Sign = base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return c
}

func (c *Client) addBatch() {

}

// 发送
func (c *Client) SendMessage() error {
	// header 头
	c.SetRequestHeader()
	cl := &http.Client{}
	resp, err := cl.Do(c.Client)
	if err != nil {
		return err
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return c.setResult(resp.StatusCode, responseData)
}

func (c *Client) GetResult() MsgRequest {
	return c.MsgRequest
}

func (c *Client) setResult(statusCode int, responseBody []byte) error {
	var msgRes MsgRequest
	if statusCode == 201 {
		var sucMessages MessageReq
		err := xml.Unmarshal(responseBody, &sucMessages)
		if err != nil {
			return err
		}
		msgRes.Status = true
		msgRes.StatusCode = statusCode
		msgRes.MessageReq = sucMessages
	} else {
		var errMessages MessageErrReq
		err := xml.Unmarshal(responseBody, &errMessages)
		if err != nil {
			return err
		}
		msgRes.Status = false
		msgRes.StatusCode = statusCode
		msgRes.MessageErrReq = errMessages
	}
	c.MsgRequest = msgRes
	return nil
}

func (c *Client) SetMessageRequest(batch *BatchSmsAttr) {
	msg := NewMessageRequest()
	msg.SetMessageAttr(batch)
	msg.SetMessage("sms message")
	msg.WriteXML()
	c.MessageRequest = msg

}

func (c *Client) GetUrl() string {
	if string(c.EndPoint[len(c.EndPoint)-1]) == "/" {
		return string(c.EndPoint[0:len(c.EndPoint)-1]) + c.setTopic()
	}
	return c.EndPoint + c.setTopic()
}

// 设置 header
func (c *Client) SetRequestHeader() {
	c.SetAttribute()
	aliUrl := c.GetUrl()
	reqest, _ := http.NewRequest("POST", aliUrl, bytes.NewBuffer([]byte(c.MessageRequest.MessageXML)))
	u, _ := url.Parse(aliUrl)
	reqest.Header.Add("Host", u.Host)
	reqest.Header.Add("Date", c.getDate())
	reqest.Header.Add("Content-Length", strconv.Itoa(len(c.MessageRequest.MessageXML)))
	reqest.Header.Add("Content-Type", "text/xml")
	reqest.Header.Add(MNS_VERSION_HEADER, MNS_VERSION)
	reqest.Header.Set(AUTHORIZATION, "MNS "+c.AccessId+":"+c.Sign)
	c.Client = reqest
}

func (c *Client) GetHeader() map[string][]string {
	return c.Client.Header
}

// 构造发送内容
type BatchSmsAttr struct {
	FreeSignName string                       // 答名
	TemplateCode string                       // 模板
	SmsParams    map[string]map[string]string //  内容变量
	Receiver     []string
	DirectSms    string
}

func NewBatchSmsAttr(freeSignName, tmpCode string) *BatchSmsAttr {
	var batch BatchSmsAttr
	batch.FreeSignName = freeSignName
	batch.TemplateCode = tmpCode
	return &batch
}

// 添加发送模板
func (bsa *BatchSmsAttr) AddReceiver(phone string, params map[string]string) {
	pro := make(map[string]map[string]string)
	pro[phone] = params
	bsa.SmsParams = pro
	bsa.Receiver = append(bsa.Receiver, phone)
}

func (bsa *BatchSmsAttr) SetFreeSignName(signName string) {
	bsa.FreeSignName = signName
}

func (bsa *BatchSmsAttr) SetTemplateCode(tmpCode string) {
	bsa.TemplateCode = tmpCode
}

// 生成 direct sms
func (bsa *BatchSmsAttr) WriteDirectSms() error {
	var sms DirectSMS
	sms.FreeSignName = bsa.FreeSignName
	sms.Type = "multiContent"
	sms.TemplateCode = bsa.TemplateCode
	sms.Receiver = strings.Join(bsa.Receiver, ",")
	jsData, err := json.Marshal(bsa.SmsParams)
	if err != nil {
		return err
	}
	sms.SmsParams = string(jsData)
	dirs, err := json.Marshal(&sms)
	bsa.DirectSms = string(dirs)
	return err
}

// 构造发送体 xml
type MessageRequest struct {
	MessageBody string
	MessageAttr *BatchSmsAttr
	MessageXML  string
}

func NewMessageRequest() *MessageRequest {
	var msgReq MessageRequest
	return &msgReq
}

func (mr *MessageRequest) SetMessageAttr(batch *BatchSmsAttr) {
	mr.MessageAttr = batch
}

func (mr *MessageRequest) SetMessage(message string) {
	mr.MessageBody = message
	mr.WriteXML()
}

func (mr *MessageRequest) WriteXML() error {
	mr.MessageAttr.WriteDirectSms()
	var rps MessageSms
	rps.Xmlns = MNS_XML_NAMESPACE
	rps.MessageBody = mr.MessageBody
	rps.MessageAttributes.DirectSMS = mr.MessageAttr.DirectSms

	data, err := xml.Marshal(&rps)
	mr.MessageXML = string(data)
	return err
}
