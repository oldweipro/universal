package idgen

import (
	"fmt"
	"os"
)

func GenerateServiceID(name string) string {
	hostname, _ := os.Hostname()
	return fmt.Sprintf("%s-%s-%d", hostname, name, os.Getpid())
}
