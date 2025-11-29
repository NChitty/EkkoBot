package environment

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
)

func GetEnvironmentFileValue(key string) string {
	fileName, isPresent := os.LookupEnv(key)
	if !isPresent {
		slog.Error(fmt.Sprintf("%s environment variable is unset.", key))
		os.Exit(1)
	}
	if fileName == "" {
		slog.Error(fmt.Sprintf("%s environment variable is empty.", key))
		os.Exit(1)
	}

	discordKeyFile, err := os.Open(fileName)
	if err == nil {
		scanner := bufio.NewScanner(discordKeyFile)
		scanner.Scan()
		return scanner.Text()
	}
	os.Exit(1)
	return ""
}
