package config

type nodeConfig struct {
	Url     string
	Timeout int
}

var NodeConfig *nodeConfig

func init() {
	url := EnvRequiredString("NODE_URL")

	timeout := EnvOptionalInt("NODE_TIMEOUT", 30)

	NodeConfig = &nodeConfig{
		Url:     url,
		Timeout: timeout,
	}
}
