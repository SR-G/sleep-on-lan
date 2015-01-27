package main

import (
    "io"
    "io/ioutil"
    "log"
)

var (
    Trace   *log.Logger
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger
)

func PreInitLoggers() {
    InitLoggers(ioutil.Discard, ioutil.Discard, ioutil.Discard, ioutil.Discard) 
}

func InitLoggers(traceHandle io.Writer, infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer) {
    Trace = log.New(traceHandle, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
    Info = log.New(infoHandle, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
    Warning = log.New(warningHandle, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
    Error = log.New(errorHandle, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}