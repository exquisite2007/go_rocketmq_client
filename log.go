package rocketmq

import (
    "io"
    "log"
    "os"
    "io/ioutil"
)

var (
    Trace   *log.Logger
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger
    logFile  os.File
)

func InitLog(filename string,level string) {

    logFile, err := os.OpenFile(filename, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("error opening file: %v", err)
    }
    var traceHandler, infoHandler,warningHandler,errorHandler io.Writer
    switch level{
    case "TRACE":
        traceHandler,infoHandler,warningHandler,errorHandler=logFile,logFile,logFile,logFile
        break
    case "INFO":
        infoHandler,warningHandler,errorHandler=logFile,logFile,logFile
        traceHandler=ioutil.Discard
        break
    case "WARNING":
        traceHandler,infoHandler= ioutil.Discard,ioutil.Discard
        warningHandler,errorHandler=logFile,logFile
        break
        traceHandler,infoHandler,warningHandler=ioutil.Discard,ioutil.Discard,ioutil.Discard
        errorHandler=logFile
        break
    default:
        traceHandler,infoHandler,warningHandler,errorHandler=ioutil.Discard,ioutil.Discard,ioutil.Discard,ioutil.Discard
    }
    Trace = log.New(traceHandler,
        "RMQ TRACE: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Info = log.New(infoHandler,
        "RMQ INFO: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Warning = log.New(warningHandler,
        "RMQ WARNING: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Error = log.New(errorHandler,
        "RMQ ERROR: ",
        log.Ldate|log.Ltime|log.Lshortfile)
}

func CloseLogFile(){
    logFile.Close()
}