package xsyslog

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	syslog "github.com/RackSec/srslog"
	"github.com/sirupsen/logrus"
)

// SyslogClient syslog 客户端
type SyslogClient struct {
	Level logrus.Level
	entry *logrus.Entry
}

type Config struct {
	Network       string           // 协议 udp, tcp, tcp+tls
	Address       string           // 地址 localhost:514
	Cert          string           // tls证书路径
	AppName       string           // log标签
	Formatter     syslog.Formatter // log 格式化模板
	DisableOutput bool             // 是否关闭控制台打印日志, 默认开启
	Level         logrus.Level     // 自定义日志等级，也可以通过SendXX指定其他等级输出
}

// NewSyslogClient 创建 NewSyslogClient 实例
func NewSyslogClient(c Config) (*SyslogClient, error) {
	hook, err := CreatHook(c.Network, c.Address, c.AppName, c.Cert, c.Formatter)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()
	logger.Hooks.Add(hook)
	if c.DisableOutput {
		logger.SetOutput(io.Discard) // 关闭控制台输出
	}
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:    true,
		DisableTimestamp: true,
	})

	// 下面这行打印error等级会覆盖hook中syslog的level
	// logger.Info("this is a syslog Info message")
	// logger.Error("this is a syslog Error message")
	// logger.Fatal("this is a syslog Fatal message")

	return &SyslogClient{Level: c.Level, entry: logrus.NewEntry(logger)}, nil
}

// CreatHook 对外暴露使用，便于支持本库不支持的其他自定义场景
func CreatHook(network, address, appName string, certPath string, formatter syslog.Formatter) (logrus.Hook, error) {
	if certPath != "" {
		hook, err := NewSyslogHookTls(address, appName, certPath, formatter)
		if err != nil {
			return hook, err
		}
		return hook, nil
	}

	hook, err := NewSyslogHook(network, address, appName, formatter)
	if err != nil {
		return hook, err
	}

	return hook, nil
}

func (s *SyslogClient) SendLog(message any) {
	msg, err := outMsg(message)
	if err != nil {
		return
	}
	s.entry.Logf(s.Level, msg)
}

func (s *SyslogClient) SendInfo(message any) {
	msg, err := outMsg(message)
	if err != nil {
		return
	}
	s.entry.Logf(logrus.InfoLevel, msg)
}

func (s *SyslogClient) SendWarn(message any) {
	msg, err := outMsg(message)
	if err != nil {
		return
	}
	s.entry.Logf(logrus.WarnLevel, msg)
}

func (s *SyslogClient) SendDebug(message any) {
	msg, err := outMsg(message)
	if err != nil {
		return
	}
	s.entry.Logf(logrus.DebugLevel, msg)
}

func (s *SyslogClient) SendError(message any) {
	msg, err := outMsg(message)
	if err != nil {
		return
	}
	s.entry.Logf(logrus.ErrorLevel, msg)
}

func (s *SyslogClient) SendFatal(message any) {
	msg, err := outMsg(message)
	if err != nil {
		return
	}
	s.entry.Logf(logrus.FatalLevel, msg)
}

func (s *SyslogClient) SendTrace(message any) {
	msg, err := outMsg(message)
	if err != nil {
		return
	}
	s.entry.Logf(logrus.TraceLevel, msg)
}

func outMsg(message any) (string, error) {
	var out = ""
	if _, ok := message.(string); ok {
		out = message.(string)
		return out, nil
	}
	bytes, err := json.Marshal(message)
	if err != nil {
		return "", errors.New(fmt.Sprint("Error marshaling message: ", message))
	}

	//reg := regexp.MustCompile(`[\n\r\t\\/\"]`)
	//out = reg.ReplaceAllString(string(bytes), "'")
	return string(bytes), nil
}
