package types

import "time"

type Config struct {
	ClientCfgPath         string
	ClientStartDelayTime  time.Duration
	ClientLevelLogs       int
	DockerInterfaceBridge string
}
