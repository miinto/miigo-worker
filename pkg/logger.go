package interfaces

type Logger interface {
	SetMainPrefix(prefix string)
	SetTempPrefix(prefix string)
	Log(body string, level string)
}
