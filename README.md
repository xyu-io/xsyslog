# xsyslog
the client and server implementations of syslog

## client demo
```go
// 创建发送者实例
	sysCli, err := NewSyslogClient(Config{
		Network:       "udp",
		Address:       "127.0.0.1:514",
		AppName:       "syslog",
		Formatter:     RFC5424Formatter,
		Level:         logrus.InfoLevel,
		DisableOutput: true,
	})
	if err != nil {
		panic(err)
		return
	}

	sysCli.SendWarn("waring! waring! waring!")
```
## server demo
```go
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
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
```
