package fakelogger

import "github.com/sirupsen/logrus"

type Fake struct {
	*logrus.Logger
}

func New() *Fake {
	return &Fake{
		Logger: logrus.New(),
	}
}

func (l *Fake) Errorf(format string, args ...interface{}) {
	return
}

func (l *Fake) Error(args ...interface{}) {
	return
}
