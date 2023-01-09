package main

import "runtime"

type TraceInfo struct {
	file     string
	line     int
	funcName string
}

func trace() TraceInfo {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return TraceInfo{"?", 0, "?"}
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return TraceInfo{file, line, "?"}
	}

	return TraceInfo{file, line, fn.Name()}
}
