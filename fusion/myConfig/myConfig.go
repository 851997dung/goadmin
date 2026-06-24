package myConfig

type MyConfig struct {
	Center      serverUrl
	BackUp      serverUrl
	Api         map[string]string
	Language    string
	Email       Email
	SecretKey   string
	ServerGroup map[int]string
	Hotfix      HotfixConfig
}

type HotfixConfig struct {
	ProjectDir    string
	ServerAddr    string
	BuildGoModule string
}

type serverUrl struct {
	Host string
	Port string
}

type Email struct {
	EmailAccount           string
	EmailAuthorizationCode string
	EmailServerHost        string
	EmailServerPort        string
}
