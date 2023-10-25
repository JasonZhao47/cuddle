package sms

type Service interface {
	Send() func()
}
