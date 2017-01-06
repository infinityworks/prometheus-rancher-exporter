package main

import "github.com/Sirupsen/logrus"

// setLogLevel - Sets the log level based on the passed argument.
func setLogLevel(level string) {
	switch level {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
		break
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
		break
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
		break
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
		break
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
		break
	}
}
