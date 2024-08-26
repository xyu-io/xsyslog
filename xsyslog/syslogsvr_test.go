package xsyslog

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestSysLogSrv(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)

	svr := NewSyslogServer(514, logrus.New(), ctx)

	msg, err := svr.RunSyslogReceiver()
	if err != nil {
		return
	}

	go func() {
		for {
			select {
			case e := <-msg:
				fmt.Printf("syslog recieved >> %+v \n", e)
			default:
				//avoid panic with goroutine asleep
			}
		}
	}()

	time.Sleep(time.Second * 25)
}

func TestSysLogSrvForever(t *testing.T) {
	ctx := context.Background()

	svr := NewSyslogServer(515, logrus.New(), ctx)

	msg, err := svr.RunSyslogReceiver()
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
