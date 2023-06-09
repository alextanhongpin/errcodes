package stacktrace

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/alextanhongpin/errcodes/stacktrace/internal"
)

const indent = "    "
const head = "Origin is:"
const tail = "Ends here:"
const body = "Caused by:"

type ErrorTrace = internal.ErrorTrace

func New(msg string, args ...any) error {
	return internal.New(msg, args...)
}

func WithStack(err error) error {
	return internal.Wrap(err, "")
}

func Wrap(err error, cause string) error {
	return internal.Wrap(err, cause)
}

func Sprint(err error, reversed bool) string {
	return sprint(err, reversed)
}

func StackTrace(err error) []Frame {
	return stacktrace(err)
}

func Unwrap(err error) ([]uintptr, map[uintptr]string) {
	return internal.Unwrap(err)
}

type Frame struct {
	ID       int    `json:"id"`
	Cause    string `json:"cause"`
	File     string `json:"file"`
	Line     int    `json:"line"`
	Function string `json:"function"`
}

func stacktrace(err error) []Frame {
	if err == nil {
		return nil
	}

	var res []Frame

	pcs, cause := Unwrap(err)
	pcs = filterFrames(pcs)

	var id int
	frames := runtime.CallersFrames(pcs)
	for {
		id++
		frame, more := frames.Next()
		if skipFrame(frame) {
			if !more {
				break
			}

			continue
		}

		msg, _ := cause[frame.PC+1]
		res = append(res, Frame{
			ID:       id,
			Cause:    msg,
			File:     frame.File,
			Function: frame.Function,
			Line:     frame.Line,
		})
		if !more {
			break
		}
	}

	return res
}

func sprint(err error, reversed bool) string {
	if err == nil {
		return ""
	}

	var sb strings.Builder

	sb.WriteString("Error:")
	sb.WriteRune(' ')
	sb.WriteString(err.Error())
	sb.WriteRune('\n')

	pcs, cause := Unwrap(err)
	pcs = filterFrames(pcs)
	pcs, cause = prettyCause(pcs, cause)
	if reversed {
		reverse(pcs)
	}

	frames := runtime.CallersFrames(pcs)
	for {
		frame, more := frames.Next()
		if skipFrame(frame) {
			if !more {
				break
			}

			continue
		}

		msg, ok := cause[frame.PC+1]
		if ok && msg != "" {
			sb.WriteString(indent)
			sb.WriteString(msg)
			sb.WriteRune('\n')
		}
		sb.WriteString(indent)
		sb.WriteString(indent)
		sb.WriteString(formatFrame(frame))
		if !more {
			break
		}

		sb.WriteRune('\n')
	}

	return sb.String()
}

func filterFrames(pcs []uintptr) []uintptr {
	var res []uintptr

	frames := runtime.CallersFrames(pcs)
	for {
		f, more := frames.Next()
		if skipFrame(f) {
			if !more {
				break
			}
			continue
		}

		res = append(res, f.PC+1)
		if !more {
			break
		}
	}

	return res
}

func skipFrame(f runtime.Frame) bool {
	// Skip empty function.
	return f.Function == "" ||
		// Skip runtime and testing package.
		strings.HasPrefix(f.Function, "runtime") ||
		strings.HasPrefix(f.Function, "testing") ||
		strings.HasPrefix(f.Function, "net") ||

		// Skip files with underscore.
		// e.g. _testmain.go
		strings.HasPrefix(f.File, "_")
}

func formatFrame(frame runtime.Frame) string {
	return fmt.Sprintf("at %s (in %s:%d)",
		prettyFunction(frame.Function),
		prettyFile(frame.File),
		frame.Line,
	)
}

func prettyFile(f string) string {
	wd, err := os.Getwd()
	if err != nil {
		return f
	}

	f = strings.TrimPrefix(f, wd)
	return strings.TrimPrefix(f, "/")
}

func prettyFunction(f string) string {
	_, file := path.Split(f)
	return file
}

func prettyCause(pcs []uintptr, cause map[uintptr]string) ([]uintptr, map[uintptr]string) {
	switch len(pcs) {
	case 0:
	case 1:
	default:
		pc := pcs[0]
		// Display the first line as "Origin is:".
		if msg, ok := cause[pc]; ok {
			cause[pc] = fmt.Sprintf("%s %s", head, msg)
		} else {
			cause[pc] = head
		}

		// Display the intermediate line as "Caused by:".
		for pc := range cause {
			if pc == pcs[0] || pc == pcs[len(pcs)-1] {
				continue
			}

			if msg, ok := cause[pc]; ok {
				cause[pc] = fmt.Sprintf("%s %s", body, msg)
			}
		}

		// Display the last line as "Ends here:".
		pc = pcs[len(pcs)-1]
		if msg, ok := cause[pc]; ok {
			cause[pc] = fmt.Sprintf("%s %s", tail, msg)
		} else {
			cause[pc] = tail
		}
	}
	return pcs, cause
}

func reverse[T any](s []T) {
	internal.Reverse(s)
}
