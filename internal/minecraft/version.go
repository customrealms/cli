package minecraft

type Version interface {
	String() string
	ApiVersion() string
	ServerJarType() string
	ServerJarUrl() string
	PluginJarUrl() string
}
