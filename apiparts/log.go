package epicrm_apiparts

//go:generate stringer -type=Severity

import (
	"log"
)

type Severity int

const (
	Debug Severity = iota
	Info
	Error
	Fatal
)

func Log(severity Severity, service string, activity string, fmt string, args ...interface{}) {
	if(severity >= Fatal) {
		log.Fatalf("%s: %s: %s: " + fmt, severity, service, activity, args)
	} else {
		log.Printf("%s: %s: %s: " + fmt, severity, service, activity, args)
	}
}

func LogError(service string, activity string, err error) {
	Log(Error, service, activity, "%s", err.Error())
}
