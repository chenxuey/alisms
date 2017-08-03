#### 阿里云消息服务-短信服务


install 
> go get github.com/chenxuey/alisms

调用代码
```
EndPoint := "mns.cn-hangzhou.aliyuncs.com"
AccessId := "access"
AccessKey := "key"
TopicName := "topicName"
client := NewClient(EndPoint, AccessId, AccessKey, TopicName)
batch := NewBatchSmsAttr("标签", "SMS_10000")
var param map[string]string
param = make(map[string]string)
param["code"] = "123456"
param["product"] = "test"
batch.AddReceiver("18657421361", param)

client.SetMessageRequest(batch)
err := client.SendMessage()
res := client.GetResult()
fmt.Println(res, err)
```
