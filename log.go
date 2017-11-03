package main

import "github.com/Sirupsen/logrus"

// setLogLevel - Sets the log level based on the passed argument.
func setLogLevel(level string) {
	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
		break
	case "info":
		log.SetLevel(logrus.InfoLevel)
		break
	case "warn":
		log.SetLevel(logrus.WarnLevel)
		break
	case "fatal":
		log.SetLevel(logrus.FatalLevel)
		break
	case "panic":
		log.SetLevel(logrus.PanicLevel)
		break
	}
}
