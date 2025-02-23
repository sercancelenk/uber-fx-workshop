package service

type Kafkathena interface {
	Consume(msg string) error
}
