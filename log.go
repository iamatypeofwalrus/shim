package shim

// Log is a simple logging interface that is satisfied by the standard library logger amongst other idiomatic loggers
type Log interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}
