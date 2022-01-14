package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"github.com/sirupsen/logrus"

	"github.com/cuipeiyu/go-logrus-aliyun-sls/hook"
)

func main() {
	var endpoint, accessKey, accessSecret string
	var project, logstor, topic, source string
	flag.StringVar(&endpoint, "endpoint", "cn-shanghai.log.aliyuncs.com", "")
	flag.StringVar(&accessKey, "access-key", "", "")
	flag.StringVar(&accessSecret, "access-secret", "", "")
	flag.StringVar(&project, "project", "logrus-test", "")
	flag.StringVar(&logstor, "logstor", "demo", "")
	flag.StringVar(&topic, "topic", "topic1", "")
	flag.StringVar(&source, "source", "127.0.0.1", "")
	flag.Parse()

	if endpoint == "" || accessKey == "" || accessSecret == "" || project == "" || logstor == "" || topic == "" || source == "" {
		panic("invalid params")
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})

	producerConfig := producer.GetDefaultProducerConfig()
	producerConfig.IsJsonType = true
	producerConfig.Endpoint = endpoint
	producerConfig.AccessKeyID = accessKey
	producerConfig.AccessKeySecret = accessSecret

	h := hook.NewSLSHook(
		producerConfig,
		hook.SetProject(project),
		hook.SetLogstor(logstor),
		hook.SetTopic(topic),
		hook.SetSource(source),
	)
	logrus.AddHook(h)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	i := 1
	s := time.Now()
loop:
	for time.Since(s) < time.Minute {
		select {
		case sig := <-sc:
			_ = sig
			h.SafeClose()
			break loop

		default:
			logrus.WithField("index", i).Infof("demo message %d", i)
			i++
			time.Sleep(300 * time.Millisecond)
			continue loop
		}
	}

	logrus.Info("exit")
}
