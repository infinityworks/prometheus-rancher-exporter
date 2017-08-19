package config

import (
	"os"

	cfg "github.com/infinityworks/go-common/config"
)

// Config - Defines application configuration
type Config struct {
	*cfg.BaseConfig
	MetaDataURL   string
	HideSys       bool
	AgentStates   []string
	HostStates    []string
	StackStates   []string
	ServiceStates []string
	HealthStates  []string
}

// Init - Initialises Config struct with safe defaults if not present
func Init() Config {
	ac := cfg.Init()

	var url = "http://rancher-metadata/latest/"

	if os.Getenv("METADATA_URL") != "" {
		url = os.Getenv("METADATA_URL")
	}

	hide := true

	appConfig := Config{
		&ac,
		url,
		hide,
		[]string{"activating", "active", "reconnecting", "disconnected", "disconnecting", "finishing-reconnect", "reconnected"},
		[]string{"activating", "active", "deactivating", "error", "erroring", "inactive", "provisioned", "purged", "purging", "registering", "removed", "removing", "requested", "restoring", "updating_active", "updating_inactive"},
		[]string{"activating", "active", "canceled_upgrade", "canceling_upgrade", "error", "erroring", "finishing_upgrade", "removed", "removing", "requested", "restarting", "rolling_back", "updating_active", "upgraded", "upgrading"},
		[]string{"activating", "active", "canceled_upgrade", "canceling_upgrade", "deactivating", "finishing_upgrade", "inactive", "registering", "removed", "removing", "requested", "restarting", "rolling_back", "updating_active", "updating_inactive", "upgraded", "upgrading"},
		[]string{"healthy", "unhealthy"},
	}

	return appConfig
}
