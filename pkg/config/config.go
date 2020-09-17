package config

import (
	"github.com/zevenet/kube-nftlb/pkg/env"
	"github.com/zevenet/kube-nftlb/pkg/types"
)

var (
	ClientCfgPath         = env.GetString("CLIENT_CFG_PATH")
	ClientLevelLogs       = types.LogLevel(env.GetInt("CLIENT_LOGS_LEVEL"))
	DockerInterfaceBridge = env.GetString("DOCKER_INTERFACE_BRIDGE")
)
