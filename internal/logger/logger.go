package logger

import (
	"fmt"
	"time"
)

func Debug(output ...any) {
	printWithTimestamp("[DEBUG] ", fmt.Sprint(output...))
}

func Info(output ...any) {
	printWithTimestamp("[INFO] ", fmt.Sprint(output...))
}

func Warn(output ...any) {
	printWithTimestamp("[WARN] ", fmt.Sprint(output...))
}

func Error(output ...any) {
	printWithTimestamp("[ERROR] ", fmt.Sprint(output...))
}

func printWithTimestamp(output ...any) {
	timestamp := time.Now().Format("2006-01-02 15:04:05") + ":"
	fmt.Println(timestamp, fmt.Sprint(output...))
}
