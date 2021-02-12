package interfaces

type Logger interface {
	SetMainPrefix(prefix string)
	SetTempPrefix(prefix string)
	LogLimited(body string)
	LogVerbose(body string)
}
