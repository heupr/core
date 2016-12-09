package bhattacharya

import (
	"bytes"
	"io"
	"log"
	"os"
)

type CoralReefLogger struct {
	Name     string
	CRLog    *log.Logger
	CRBuff   *bytes.Buffer
	Backtest bool
}

func CreateLog(name string, backtest bool) CoralReefLogger {
	buf := bytes.Buffer{}
	logger := log.New(&buf, name, log.Lshortfile)
	newLog := CoralReefLogger{Name: name, CRLog: logger, CRBuff: &buf, Backtest: backtest}
	return newLog
}

func (crl *CoralReefLogger) Log(input interface{}) {
	crl.CRLog.Print(input)
}

func (crl *CoralReefLogger) Flush() {
	output := ""
	if crl.Backtest {
		output = "../../data/backtests/" + crl.Name
	} else {
		output = "../../data/logs/" + crl.Name
	}
	f, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}
	io.WriteString(f, crl.CRBuff.String())
}
