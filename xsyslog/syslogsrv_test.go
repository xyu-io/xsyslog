package xsyslog

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestSysLogSrv(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	srv := NewSyslogServer(514, logrus.New(), ctx)

	msg, err := srv.RunSyslogReceiver()
	if err != nil {
		return
	}

	go func() {
		for {
			select {
			case e := <-msg:
				fmt.Printf("syslog recieved >> %+v \n", e)
			}
		}
	}()

	time.Sleep(time.Second * 15)
}

func TestSysLogSrvForever(t *testing.T) {
	ctx := context.Background()

	srv := NewSyslogServer(515, logrus.New(), ctx)

	msg, err := srv.RunSyslogReceiver()
	if err != nil {
		return
	}

	for {
		select {
		case e := <-msg:
			fmt.Printf("syslog recieved >> %+v \n", e)
		}
	}

}
