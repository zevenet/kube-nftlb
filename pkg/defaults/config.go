package defaults

import (
	"github.com/zevenet/kube-nftlb/pkg/env"
	"github.com/zevenet/kube-nftlb/pkg/types"
)

var (
	clientCfgPath         = env.GetString("CLIENT_CFG_PATH")
	clientLevelLogs       = env.GetInt("CLIENT_LOGS_LEVEL")
	clientStartDelayTime  = env.GetTime("CLIENT_START_DELAY_TIME")
	dockerInterfaceBridge = env.GetString("DOCKER_INTERFACE_BRIDGE")
)

// GetCfg
func GetCfg() *types.Config {
	return &types.Config{
		ClientCfgPath:         clientCfgPath,
		ClientLevelLogs:       clientLevelLogs,
		ClientStartDelayTime:  clientStartDelayTime,
		DockerInterfaceBridge: dockerInterfaceBridge,
	}
}
