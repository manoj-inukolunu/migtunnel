package util

import "log"

func LogWithPrefix(prefix string, strToLog string) {
	log.Println("[" + prefix + "] " + strToLog)
}
