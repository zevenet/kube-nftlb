package log

import (
	"fmt"

	"github.com/zevenet/kube-nftlb/pkg/config"
	"github.com/zevenet/kube-nftlb/pkg/types"
)

var logLevel = config.ClientLevelLogs

// WriteLog prints to stdout every message below the log level.
func WriteLog(externalLogLevel types.LogLevel, message string) {
	if externalLogLevel <= logLevel {
		fmt.Println(message)
	}
}
