package shared

import "log"

func ErrorHandler() {
	if err := recover(); err != nil {
		log.Fatalf("error occured unexpectedly: %s", err)
	}
}
