package args

import (
	"encoding/json"
	"os"
	"strings"
)

type Parser struct {
	args       map[string]string
	envPrefix  string
	configPath string
}

func NewParser(envPrefix string) Parser {
	return Parser{
		args:      make(map[string]string),
		envPrefix: envPrefix,
	}
}

func convertNameToEnv(name string) string {
	envKey := strings.ToUpper(name)
	envKey = strings.ReplaceAll(envKey, "", "_")
	return envKey
}

func readConfigFile(configPath string) (map[string]string, error) {
	var config map[string]string = make(map[string]string)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return config, nil
}

func (p *Parser) AddArg(name string, desc string) {
	p.args[name] = ""
}



