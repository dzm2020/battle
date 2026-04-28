package log

import "fmt"

func Debug(format string, a ...any) {
	format = fmt.Sprintf("[debug]	%s\n", format)
	fmt.Printf(format, a...)
}

func Info(format string, a ...any) {
	format = fmt.Sprintf("[info]	%s\n", format)
	fmt.Printf(format, a...)
}

func Error(format string, a ...any) {
	format = fmt.Sprintf("[error]	%s\n", format)
	fmt.Printf(format, a...)
}
