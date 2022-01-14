package hook

import (
	"fmt"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/aliyun-log-go-sdk/producer"
	"github.com/sirupsen/logrus"
)

var _ logrus.Hook = (*SLSHook)(nil)

type OptionFunc func(*Option)

type Option struct {
	project, logstor, topic, source string
}

func SetProject(name string) OptionFunc {
	return func(o *Option) {
		o.project = name
	}
}

func SetLogstor(name string) OptionFunc {
	return func(o *Option) {
		o.logstor = name
	}
}

func SetTopic(name string) OptionFunc {
	return func(o *Option) {
		o.topic = name
	}
}

func SetSource(name string) OptionFunc {
	return func(o *Option) {
		o.source = name
	}
}

func NewSLSHook(producerConfig *producer.ProducerConfig, opts ...OptionFunc) *SLSHook {
	opt := &Option{}
	if len(opts) > 0 {
		for _, fun := range opts {
			if fun != nil {
				fun(opt)
			}
		}
	}

	producer := producer.InitProducer(producerConfig)

	producer.Start()

	return &SLSHook{opt, producer}
}

type SLSHook struct {
	opt      *Option
	producer *producer.Producer
}

func (hook *SLSHook) Fire(entry *logrus.Entry) error {
	ts := uint32(entry.Time.Unix())

	contents := []*sls.LogContent{}
	for key, s := range entry.Data {
		v := fmt.Sprint(s)
		contents = append(contents, &sls.LogContent{
			Key:   &key,
			Value: &v,
		})
	}

	msgKey := "message"
	msgContent := entry.Message
	contents = append(contents, &sls.LogContent{
		Key:   &msgKey,
		Value: &msgContent,
	})

	log := &sls.Log{
		Time:     &ts,
		Contents: contents,
	}

	hook.producer.SendLog(hook.opt.project, hook.opt.logstor, hook.opt.topic, hook.opt.source, log)
	return nil
}

func (hook *SLSHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *SLSHook) SafeClose() {
	if hook.producer != nil {
		hook.producer.SafeClose()
	}
}

func (hook *SLSHook) Close(timeoutMs int64) error {
	if hook.producer != nil {
		return hook.producer.Close(timeoutMs)
	}
	return nil
}
