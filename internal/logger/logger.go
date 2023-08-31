package logger

import (
	"io"
	"log"
)

var InfoLogger = log.New(log.Default().Writer(), "INFO | ", log.Ldate|log.Ltime|log.LUTC)
var ErrorLogger = log.New(log.Default().Writer(), "ERRO | ", log.Ldate|log.Ltime|log.LUTC)
var DebugLogger = log.New(io.Discard, "DBUG | ", log.Default().Flags())

func EnableDebugLogging() {
	InfoLogger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.LUTC)
	ErrorLogger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.LUTC)
	DebugLogger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.LUTC)
	DebugLogger.SetOutput(log.Default().Writer())
	DebugLogger.Println("Debug logging enabled")
}
