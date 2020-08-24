package logs

import (
	"fmt"

	"github.com/zevenet/kube-nftlb/pkg/config"
)

func PrintLogChannelFuncGeneral(levelLog int, action string, resourceName string, obj interface{}, logChannel chan string) chan string {
	// Print the logs of the add or remove actions of the funcs. Every time a service is created or deleted these logs will appear
	// In essence it prints the action (either delete or add) of the resource and prints the content of the object.
	// See watchers.go
	if checkLevelLogs(levelLog) == true {
		logChannel <- fmt.Sprintf("\n"+action+"%s:\n%s", resourceName, obj)
	}
	return logChannel
}

func PrintLogChannelFuncUpdate(levelLog int, action string, resourceName string, oldObj interface{}, newObj interface{}, logChannel chan string) chan string {
	// Print the logs of the update actions. In this case it prints both the original data and after doing the update
	// In essence it prints the action of the resource and prints the content of the object before and after being updated
	// See watchers.go
	if checkLevelLogs(levelLog) == true {
		logChannel <- fmt.Sprintf(action, resourceName, oldObj, newObj)
	}
	return logChannel
}

func PrintLogChannel(levelLog int, message string, logChannel chan string) {
	// Within the add function, all the information we show about the process
	if checkLevelLogs(levelLog) == true {
		logChannel <- message
	}
}

func checkLevelLogs(levelPrint int) bool {
	// Reads the level of logs established within the parameterizable values
	// If the level of logs is less than or equal to the established one, it returns true and gives the green light to print logs.
	logLevel := config.ClientLevelLogs
	if levelPrint <= logLevel {
		return true
	}
	return false
}
