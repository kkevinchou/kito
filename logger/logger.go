package logger

import (
	"fmt"

	"github.com/kkevinchou/ant/settings"
)

func Debug(message string) {
	if settings.LoggingLevel <= 0 {
		fmt.Println(message)
	}
}

func Debug1(message string) {
	if settings.LoggingLevel <= -1 {
		fmt.Println(message)
	}
}
