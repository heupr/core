package bhattacharya

import (
	"bytes"
	"io"
	"log"
	"os"
)

type CoralReefLogger struct {
	Name   string
	CRLog  *log.Logger
	CRBuff *bytes.Buffer
}

func CreateLog(name string) CoralReefLogger {
	buf := bytes.Buffer{}
	logger := log.New(&buf, name, log.Lshortfile)
	newLog := CoralReefLogger{Name: name, CRLog: logger, CRBuff: &buf}
	return newLog
}

func (crl *CoralReefLogger) Log(input interface{}) {
	crl.CRLog.Print(input)
}

func (crl *CoralReefLogger) Flush() {
	f, err := os.OpenFile(crl.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		panic(err)
	}
	io.WriteString(f, crl.CRBuff.String())
}
