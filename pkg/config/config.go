package config

import (
	"github.com/zevenet/kube-nftlb/pkg/env"
)

var (
	ClientCfgPath         = env.GetString("CLIENT_CFG_PATH")
	ClientLevelLogs       = env.GetInt("CLIENT_LOGS_LEVEL")
	ClientStartDelayTime  = env.GetTime("CLIENT_START_DELAY_TIME")
	DockerInterfaceBridge = env.GetString("DOCKER_INTERFACE_BRIDGE")
)
