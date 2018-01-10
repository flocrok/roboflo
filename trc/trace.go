package trc

import (
	"io/ioutil"
	"io"
	"os"
	"fmt"
    "log"
	"time"
)

var (
    Trace   *log.Logger
    Info   *log.Logger
    Warning   *log.Logger
    Error   *log.Logger
    Stochlog   *log.Logger
    Tickerlog   *log.Logger
	Level int // 0 Rien, 1 Error, 2 Warning, 3 Info, 4 Trace
)

func Init(level int, monnaie string) int {
	var (
		traceHandle io.Writer = ioutil.Discard
		infoHandle io.Writer = ioutil.Discard
		warningHandle io.Writer = ioutil.Discard
		errorHandle io.Writer = ioutil.Discard
	)

	Level = level
	LogFilename := fmt.Sprintf("LogFile_%s.log",time.Now().Format("2006-01-02"))
	LogFile, err := os.OpenFile(LogFilename,os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {return 1}

	StochFilename := fmt.Sprintf("StochFile_%s.csv",time.Now().Format("2006-01-02"))
	StochFile, err := os.OpenFile(StochFilename,os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {return 1}

	TickerFilename := fmt.Sprintf("TickerFile_%s_%s.log",monnaie,time.Now().Format("2006-01-02"))
	TickerFile, err := os.OpenFile(TickerFilename,os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {return 1}

	if Level > 0 {errorHandle = LogFile}

	if Level > 1 { warningHandle = LogFile}
	if Level > 2 { infoHandle = LogFile}
	if Level > 3 { traceHandle = LogFile}

    Trace = log.New(traceHandle,
        "TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

    Info = log.New(infoHandle,
        "INFO : ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Warning = log.New(warningHandle,
        "WARN : ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Error = log.New(errorHandle,
        "ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
		
	Stochlog = log.New(StochFile,"",0)
	Tickerlog = log.New(TickerFile,"",0)
	return 0
}

func Rotate(monnaie string) int {
	
	TickerFilename := fmt.Sprintf("TickerFile_%s_%s.log",monnaie,time.Now().Format("2006-01-02"))
	TickerFile, err := os.OpenFile(TickerFilename,os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {return 1}

	Tickerlog = log.New(TickerFile,"",0)
	return 0
}