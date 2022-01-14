# logrus hook aliyun sls

## Demo

```go
import "github.com/aliyun/aliyun-log-go-sdk/producer"
import "github.com/sirupsen/logrus"
import "github.com/cuipeiyu/go-logrus-aliyun-sls/hook"

func initLogrus() {
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
}

```
