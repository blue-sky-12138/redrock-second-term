package sky

const (
	debugMode = iota
)

var skyMode = debugMode

func isDebugging() bool {
	return skyMode == debugMode
}
