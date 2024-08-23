package xsyslog

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	syslog "github.com/RackSec/srslog"
	"github.com/sirupsen/logrus"
)

const (
	SecureProto = "tcp+tls"
)

var (
	RFC5424Formatter = syslog.RFC5424Formatter // 推荐使用,注意消息长度2048，超过会被截断
	RFC3164Formatter = syslog.RFC3164Formatter // 已被弃用
	UnixFormatter    = syslog.UnixFormatter
	DefaultFormatter = syslog.DefaultFormatter
)

// SyslogTlsHook to send logs via syslog.
type SyslogTlsHook struct {
	Writer *syslog.Writer
}

func NewSyslogHook(network, raddr string, tag string, formatter syslog.Formatter) (*SyslogTlsHook, error) {
	w, err := syslog.Dial(network, raddr, syslog.LOG_USER|syslog.LOG_WARNING, tag)
	if err != nil {
		return nil, err
	}
	if formatter != nil {
		w.SetFormatter(formatter)
	}

	return &SyslogTlsHook{w}, err
}

func NewSyslogHookTls(raddr string, tag string, certPath string, formatter syslog.Formatter) (*SyslogTlsHook, error) {
	serverCert, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(serverCert)
	config := tls.Config{
		RootCAs: pool,
	}
	config.InsecureSkipVerify = true // 关闭任意证书放行

	w, err := syslog.DialWithTLSConfig(SecureProto, raddr, syslog.LOG_USER|syslog.LOG_WARNING, tag, &config)
	if err != nil {
		return nil, err
	}
	if formatter != nil {
		w.SetFormatter(formatter)
	}
	return &SyslogTlsHook{w}, err
}

func (hook *SyslogTlsHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	switch entry.Level {
	case logrus.PanicLevel:
		return hook.Writer.Crit(line)
	case logrus.FatalLevel:
		return hook.Writer.Crit(line)
	case logrus.ErrorLevel:
		return hook.Writer.Err(line)
	case logrus.WarnLevel:
		return hook.Writer.Warning(line)
	case logrus.InfoLevel:
		return hook.Writer.Info(line)
	case logrus.DebugLevel:
		return hook.Writer.Debug(line)
	default:
		return nil
	}
}

func (hook *SyslogTlsHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
