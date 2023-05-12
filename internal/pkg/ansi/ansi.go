package ansi

import "fmt"

const (
	// Styles
	NORMAL      = "\x1b[0m"
	BOLD        = "\x1b[1m"
	FAINT       = "\x1b[2m"
	ITALIC      = "\x1b[3m"
	UNDERLINE   = "\x1b[4m"
	BLINK_SLOW  = "\x1b[5m"
	BLINK_RAPID = "\x1b[6m"
	INVERSE     = "\x1b[7m"
	CONCEAL     = "\x1b[8m"
	CROSSED_OUT = "\x1b[9m"

	// Text colors.
	BLACK   = "\x1b[30m"
	RED     = "\x1b[31m"
	GREEN   = "\x1b[32m"
	YELLOW  = "\x1b[33m"
	BLUE    = "\x1b[34m"
	MAGENTA = "\x1b[35m"
	CYAN    = "\x1b[36m"
	WHITE   = "\x1b[37m"

	// Background colors.
	BG_BLACK   = "\x1b[40m"
	BG_RED     = "\x1b[41m"
	BG_GREEN   = "\x1b[42m"
	BG_YELLOW  = "\x1b[43m"
	BG_BLUE    = "\x1b[44m"
	BG_MAGENTA = "\x1b[45m"
	BG_CYAN    = "\x1b[46m"
	BG_WHITE   = "\x1b[47m"
	// Resets
	NOSTYLE     = "\x1b[0m"
	NOUNDERLINE = "\x1b[24m"
	NOINVERSE   = "\x1b[27m"
	NOCOLOR     = "\x1b[39m"
)

func Success(a ...interface{}) {
	fmt.Print(BG_GREEN)
	fmt.Print(a...)
	fmt.Print(NOSTYLE)
}

func Warn(a ...interface{}) {
	fmt.Print(BG_YELLOW)
	fmt.Print(a...)
	fmt.Print(NOSTYLE)
}

func Error(a ...interface{}) {
	fmt.Print(BG_RED)
	fmt.Print(a...)
	fmt.Print(NOSTYLE)
}

func Successln(a ...interface{}) {
	fmt.Print(BG_GREEN)
	fmt.Print(a...)
	fmt.Println(NOSTYLE)
}

func Warnln(a ...interface{}) {
	fmt.Print(BG_YELLOW)
	fmt.Print(a...)
	fmt.Println(NOSTYLE)
}

func Errorln(a ...interface{}) {
	fmt.Print(BG_RED)
	fmt.Print(a...)
	fmt.Println(NOSTYLE)
}
