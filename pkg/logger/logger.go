package logger

import "github.com/sirupsen/logrus"

type Interface interface {
	logrus.FieldLogger
}
