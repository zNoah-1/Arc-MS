package task

import (
	"fmt"
	"time"

	"github.com/zNoah-1/Arc-MS/internal/logger"
)

func InitInactiveCheck(serverList *[]map[string]any) {
	go checkForInactiveServers(serverList)
}

func checkForInactiveServers(serverList *[]map[string]any) {
	for {
		logger.Info("Checking for inactive servers...")

		unixTime := time.Now().Unix()
		//logger.Debug("Time (epoch): ", time.Now().Unix())

		i := 0
		for i < len(*serverList) {
			serverData := (*serverList)[i]
			lastUpdate := serverData["lastUpdate"].(int64)

			if unixTime-lastUpdate > 3600 { //3600 = 1h
				logger.Info(fmt.Sprint("Removing server ", serverData["ip"], " for inactivity"))
				*serverList = append((*serverList)[:i], (*serverList)[i+1:]...)

				i--
			} else if unixTime-lastUpdate > 1800 { //1800 = 30 min
				if serverData["disabled"].(bool) {
					break
				}

				logger.Info(fmt.Sprint("Disabling server ", serverData["ip"], " for inactivity"))
				serverData["disabled"] = true
			}

			i++
		}

		time.Sleep(60 * time.Second)
	}
}
