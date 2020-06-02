package goerror

import (
    "fmt"
    "runtime"
)

var (
    DefaultStackTraceCap      = 15
    DefaultStackTraceSkipLine = 2
)

type frame struct {
    funcName string
    line     int
    path     string
}

func (e *GoError) StackTrace() string {
    stackTrance := ""

    for _, frame := range e.frames {
        stackTrance += fmt.Sprintf("%s: %s %d \n", frame.funcName, frame.path, frame.line)
    }

    return stackTrance
}

func trace(skip int) []*frame {
    frames := make([]*frame, 0)

    i := 0
    for {
        if i >= DefaultStackTraceCap {
            break
        }

        pc, path, line, ok := runtime.Caller(skip)
        if !ok {
            break
        }

        frames = append(frames, &frame{
            funcName: runtime.FuncForPC(pc).Name(),
            line:     line,
            path:     path,
        })

        skip++
        i++
    }

    return frames
}
