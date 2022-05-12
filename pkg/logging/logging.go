package logging

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

func Debug(format string, v ...interface{}) {
	if viper.GetBool("DEBUG") {
		format = fmt.Sprintf("[debug] %s\n", format)
		log.Output(2, fmt.Sprintf(format, v...))
	}
}

func WrapError(format string, a ...any) error {
	return fmt.Errorf(format, a...)
}

func Warning(format string, v ...interface{}) {
	format = fmt.Sprintf("WARNING: %s\n", format)
	fmt.Fprintf(os.Stderr, format, v...)
}
