package wechatrequests

import (
	"bytes"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/liangboceo/yuanboot/abstractions"
	"github.com/liangboceo/yuanboot/abstractions/xlog"
	"io/ioutil"
	"net/http"
)

func PostWechatMessage(sendUrl, msg string) string {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", sendUrl, bytes.NewBuffer([]byte(msg)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("charset", "UTF-8")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	strBody := string(body)
	fmt.Println("response Body:", strBody)
	return strBody
}

func SendTxtMessage(request GrafanaAlertRequest, config abstractions.IConfiguration) string {
	tag := request.GetTag()
	logger := xlog.GetXLogger("wechat")
	js, _ := sonic.Marshal(request)
	logger.Infof("Request json: %s", string(js))
	if tag == "" {
		logger.Infof("no send")
		return ""
	}
	sendUrl := config.Get(fmt.Sprintf("alert.%s.webhook_url", tag)).(string)
	linkUrl := config.Get(fmt.Sprintf("alert.%s.link_url", tag)).(string)
	logger.Infof("request tag:%s", tag)
	logger.Infof(sendUrl)
	logger.Infof(linkUrl)

	var message *MarkdownMessage
	if request.State == "alerting" && len(request.EvalMatches) > 0 {
		message = &MarkdownMessage{
			Markdown: struct {
				Content string `json:"content" gorm:"column:content"`
			}{
				Content: "## " + request.RuleName + ",请相关同事注意。\n" +
					" > [报警信息] : " + request.Message + "\n" +
					" > [报警次数] : <font color=\"warning\">" + request.GetMetricValue() + "次</font>" + "\n" +
					" > [报警明细] : (" + linkUrl + ")\n",
			},
			Msgtype: "markdown",
		}
	}
	msg, _ := sonic.Marshal(message)
	msgStr := string(msg)
	logger.Infof("send message:%s", msgStr)

	//return sendUrl + msgStr
	return PostWechatMessage(sendUrl, msgStr)
}
