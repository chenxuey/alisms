package alisms

import "encoding/xml"

const (
	MNS_VERSION_HEADER = "x-mns-version"
	MNS_XML_NAMESPACE  = "http://mns.aliyuncs.com/doc/v1/"
	MNS_VERSION        = "2015-06-06"
	MNS                = "MNS"
	AUTHORIZATION      = "Authorization"
	CONTENT_LENGTH     = "Content-Length"
	CONTENT_TYPE       = "Content-Type"
	GMT_DATE_FORMAT    = "Mon, 2 Jan 2006 15:04:05 GMT"
)

// 发送 xml struct
type MessageSms struct {
	XMLName           xml.Name `xml:"Message"`
	Xmlns             string   `xml:"xmlns,attr"`
	MessageBody       string
	MessageAttributes MessageAttributes
}

type MessageAttributes struct {
	DirectSMS string
}

type DirectSMS struct {
	Type         string
	FreeSignName string
	TemplateCode string
	Receiver     string
	SmsParams    string
}

// 接收 xml struct success
type MessageReq struct {
	XMLName        xml.Name `xml:"Message"`
	Xmlns          string   `xml:"xmlns,attr"`
	MessageId      string
	MessageBodyMD5 string
}

// 接收 xml struct fail
type MessageErrReq struct {
	XMLName   xml.Name `xml:"Error"`
	Xmlns     string   `xml:"xmlns,attr"`
	Code      string
	Message   string
	RequestID string
	HostId    string
}

type MsgRequest struct {
	Status     bool
	StatusCode int
	MessageReq
	MessageErrReq
}
