package logger

import (
	"fmt"
	"github.com/fatih/color"
)

const (
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
)

type Log interface {
	Write(bool) string
}

type InfoLog string // INFO

func (il InfoLog) Write(colored bool) string {
	if colored {
		return color.CyanString("[%s]: %s\n", INFO, il)
	}
	return fmt.Sprintf("[%s]: %s\n", INFO, il)
}

type WarnLog struct {
	severity WarnSeverity
	message  string
}

func (wl WarnLog) Write(colored bool) string {
	if colored {
		switch wl.severity {
		case WarnSevere:
			return color.HiYellowString("[%s (%s)]: %s\n", WARN, wl.severity, wl.message)
		case WarnMedium:
			return color.YellowString("[%s (%s)]: %s\n", WARN, wl.severity, wl.message)
		}
		return color.HiCyanString("[%s (%s)]: %s\n", WARN, wl.severity, wl.message)
	}
	return fmt.Sprintf("[%s (%s)]: %s\n", WARN, wl.severity, wl.message)
}

type WarnSeverity string

const (
	WarnSevere WarnSeverity = "High"
	WarnMedium              = "Medium"
	WarnLow                 = "Low"
)

type ErrorLog struct {
	error error
	code  int
}

func (el ErrorLog) Write(colored bool) string {
	if colored {
		return color.RedString("[%s (%d)]: %s \n", ERROR, el.code, el.error.Error())
	}
	return fmt.Sprintf("[%s (%d)]: %s\n", ERROR, el.code, el.error.Error())
}
