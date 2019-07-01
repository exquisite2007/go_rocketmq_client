package rocketmq

import (
    "io"
    "log"
)

var (
    Trace   *log.Logger
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger
)

func InitLog(
    traceHandle io.Writer,
    infoHandle io.Writer,
    warningHandle io.Writer,
    errorHandle io.Writer) {

    Trace = log.New(traceHandle,
        "RMQ TRACE: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Info = log.New(infoHandle,
        "RMQ INFO: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Warning = log.New(warningHandle,
        "RMQ WARNING: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Error = log.New(errorHandle,
        "RMQ ERROR: ",
        log.Ldate|log.Ltime|log.Lshortfile)
}
