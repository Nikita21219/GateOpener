package utils

import "os"

func Debug() bool {
	debugMode := false
	if debug := os.Getenv("DEBUG"); debug == "1" || debug == "true" {
		debugMode = true
	}

	if debug := os.Getenv("debug"); debug == "1" || debug == "true" {
		debugMode = true
	}

	return debugMode
}
