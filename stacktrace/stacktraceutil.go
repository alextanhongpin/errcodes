package stacktrace

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"sort"
	"strings"

	"golang.org/x/exp/slices"
)

func reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < len(s)/2; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func formatFrame(frame runtime.Frame) string {
	return fmt.Sprintf("at %s (in %s:%d)",
		prettyFunction(frame.Function),
		prettyFile(frame.File),
		frame.Line,
	)
}

func callers(skip int) []uintptr {
	const depth = 64
	var pc [depth]uintptr
	n := runtime.Callers(skip, pc[:])
	if n == 0 {
		return nil
	}

	var pcs = pc[:n]
	return pcs
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

func mergePCs(x []uintptr, y ...uintptr) []uintptr {
	z := append(x, y...)

	sort.Slice(z, func(i, j int) bool {
		return z[i] < z[j]
	})

	return slices.Compact(z)
}

func frames(pcs []uintptr) []runtime.Frame {
	var res []runtime.Frame
	frames := runtime.CallersFrames(pcs)
	for {
		f, more := frames.Next()

		// Skip runtime and testing package.
		if strings.HasPrefix(f.Function, "runtime") ||
			strings.HasPrefix(f.Function, "testing") {
			continue
		}

		// Skip files with underscore.
		// e.g. _testmain.go
		if strings.HasPrefix(f.File, "_") {
			continue
		}

		if f.Function == "" {
			break
		}

		res = append(res, f)
		if !more {
			break
		}
	}

	return res
}
