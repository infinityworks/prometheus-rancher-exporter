package compatibility

import (
	"strings"

	"github.com/prometheus/log"
)

// StacksURLCheck - Checks the API version for Rancher to determine the correct URL
func StacksURLCheck(rancherURL string) string {

	var stacksEndpoint string

	if strings.Contains(rancherURL, "v1") {
		log.Info("Version 1 API detected, using legacy API fields")
		stacksEndpoint = "/environments/"

	} else if strings.Contains(rancherURL, "v2") {
		log.Info("Version 2 API detected, using updated API fields")
		stacksEndpoint = "/stacks/"

	} else {
		log.Info("No known API version detected, defaulting to /stacks/")
		stacksEndpoint = "/stacks/"
	}

	return stacksEndpoint

}
