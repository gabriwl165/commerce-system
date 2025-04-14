package env

import (
	"errors"
	"os"
	"strings"
)

type EnvManager struct {
	envs map[string]any
}

func (e *EnvManager) Read(key string) (any, error) {
	value, ok := e.envs[key]
	if !ok {
		return nil, errors.New("Key Not Found")
	}
	return value, nil
}

func (e *EnvManager) Write(key string, value string) {
	e.envs[key] = value
}

func Init() *EnvManager {
	envManager := &EnvManager{
		envs: make(map[string]any),
	}
	for _, env := range os.Environ() {
		values := strings.Split(env, "=")
		envManager.Write(values[0], values[1])
	}
	return envManager
}
