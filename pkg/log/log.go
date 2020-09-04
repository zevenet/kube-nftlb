package log

import (
	"fmt"

	"github.com/zevenet/kube-nftlb/pkg/config"
)

var logLevel = config.ClientLevelLogs

func WriteLogFuncGeneral(levelLog int, action string, resourceName string, obj interface{}) {
	// Print the logs of the add or remove actions of the funcs. Every time a service is created or deleted these logs will appear
	// In essence it prints the action (either delete or add) of the resource and prints the content of the object.
	// See watchers.go
	if checkLevelLogs(levelLog) {
		message := fmt.Sprintf("\n"+action+"%s:\n%s", resourceName, obj)
		fmt.Println(message)
	}
}

func WriteLogFuncUpdate(levelLog int, resourceName string, oldObj interface{}, newObj interface{}) {
	// Print the logs of the update actions. In this case it prints both the original data and after doing the update
	// In essence it prints the action of the resource and prints the content of the object before and after being updated
	// See watchers.go
	if checkLevelLogs(levelLog) {
		message := fmt.Sprintf("\nUpdatedd %s:\n* BEFORE: %s\n* NOW: %s", resourceName, oldObj, newObj)
		fmt.Println(message)
	}
}

func WriteLog(levelLog int, message string) {
	// Within the add function, all the information we show about the process
	if checkLevelLogs(levelLog) {
		fmt.Println(message)
	}
}

func checkLevelLogs(levelPrint int) bool {
	return levelPrint <= logLevel
}
