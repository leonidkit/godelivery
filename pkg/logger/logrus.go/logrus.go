package logrus

import (
	"strings"

	"github.com/sirupsen/logrus"
)

type Logrus struct {
	*logrus.Logger
}

func formatFilePath(path string) string {
	arr := strings.Split(path, "/")
	return arr[len(arr)-1]
}

func New() *Logrus {
	lgrs := &Logrus{
		logrus.New(),
	}

	lgrs.SetReportCaller(true)
	formatter := &logrus.TextFormatter{
		TimestampFormat:        "02-01-2006 15:04:05", // the "time" field configuration
		FullTimestamp:          true,
		DisableLevelTruncation: true, // log level field configuration
	}
	lgrs.SetFormatter(formatter)

	return lgrs
}
