package domain

type Env struct {
	SQS *SQSEnv
}

type SQSEnv struct {
	QueueName string `env:"SQS_QUEUE_NAME"`
	QueueURL  string `env:"SQS_QUEUE_URL"`
	MaxNum    int64  `env:"SQS_MAX_NUM"`
	WaitTime  int64  `env:"SQS_WAIT_TIME"`
}
