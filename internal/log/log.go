package log

import (
	"fmt"
	"log"
)

func Logf(component, format string, args ...any) {
	log.Printf("%s %s", pad(component), fmt.Sprintf(format, args...))
}

func Panicf(component, format string, args ...any) {
	log.Panicf("%s %s", pad(component), fmt.Sprintf(format, args...))
}

func Fatalf(component, format string, args ...any) {
	log.Fatalf("%s %s", pad(component), fmt.Sprintf(format, args...))
}

func pad(component string) string {
	return fmt.Sprintf("%-6.6s:", component)
}
