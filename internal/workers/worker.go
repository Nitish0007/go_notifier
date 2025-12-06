package workers

type Worker interface {
	Consume()
	ConsumeRetry()
}
