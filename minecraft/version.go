package minecraft

type Version interface {
	String() string
	ApiVersion() string
	ServerJarUrl() string
}
