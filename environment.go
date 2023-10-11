package service

import (
	"fmt"
	"os"
	"strings"
)

func ParseEnvironment() (env map[string]string) {
	rawEnv := os.Environ()
	for _, rawVar := range rawEnv {
		envVar := strings.Split(rawVar, "=")
		env[envVar[0]] = envVar[1]
	}
	return env
}

type EnvVar struct {
	Key   string
	Value string
}

func (p *Process) Environment() (env []string) {
	for eKey, eValue := range p.Env {
		env = append(env, fmt.Sprintf("%s=%s", eKey, eValue))
	}
	return env
}

func (p *Process) AppendEnvVar(envKey, envValue string) *Process {
	p.Env[envKey] = envValue
	return p

}

func (p *Process) ParseEnvironment() (env map[string]string) {
	rawEnv := os.Environ()
	for _, rawVar := range rawEnv {
		envVar := strings.Split(rawVar, "=")
		env[envVar[0]] = envVar[1]
	}
	p.Env = env
	return p.Env
}

func (p *Process) IsRunning() bool {
	// TODO: Is there really no better way to handle this?
	return os.Getenv("_DAEMON") == "1"
}
