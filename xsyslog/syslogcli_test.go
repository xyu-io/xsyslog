package xsyslog

import (
	"github.com/sirupsen/logrus"
	"testing"
)

type Data struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

func TestSysLogCli(t *testing.T) {
	// 创建发送者实例
	sysCli, err := NewSyslogClient(Config{
		Network:       "udp",
		Address:       "127.0.0.1:515",
		AppName:       "syslog",
		Formatter:     RFC5424Formatter,
		Level:         logrus.InfoLevel,
		DisableOutput: true,
	})
	if err != nil {
		panic(err)
		return
	}

	msg := Data{
		ID:      1001,
		Type:    "Test",
		Tag:     "data-self",
		Message: "hello world",
	}
	// 发送多条日志消息
	for i := 0; i < 2; i++ {
		sysCli.SendLog(msg)
	}
	sysCli.SendWarn("waring! waring! waring!")
}
