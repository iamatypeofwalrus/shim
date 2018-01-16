package shim

// Log provides an easy way to get insight into what Lambda and API Gateway are passing and how
// Shim is handling those arguments
type Log interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}
