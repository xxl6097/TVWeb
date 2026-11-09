package service

import (
	"github.com/kardianos/service"
	"github.com/xxl6097/glog/pkg/z"
	"github.com/xxl6097/go-service/pkg/gs/igs"
)

// TVWebService represents the TVWeb service
type TVWebService struct {
	// TODO: Add TVBox integration when available
}

// OnConfig returns the service configuration
func (t *TVWebService) OnConfig() *service.Config {
	return &service.Config{
		Name:        "TVWebService",
		DisplayName: "TVWeb Service",
		Description: "TVWeb - Video Box Web Service",
	}
}

// OnVersion returns the service version
func (t *TVWebService) OnVersion() string {
	return "0.0.1"
}

// OnRun starts the service
func (t *TVWebService) OnRun(svc igs.Service) error {
	z.L().Info("Starting TVWeb service...")

	// TODO: Initialize and start TVBox when the dependency is available
	// For now, this is a placeholder implementation

	z.L().Info("TVWeb service started successfully")
	return nil
}

// OnStop stops the TVWeb service
func (t *TVWebService) OnStop() {
	z.L().Info("Stopping TVWeb service...")

	// TODO: Stop TVBox when implemented

	z.L().Info("TVWeb service stopped")
}

// OnShutdown handles service shutdown
func (t *TVWebService) OnShutdown() {
	z.L().Info("Shutting down TVWeb service...")
}

// GetAny returns service configuration or state
func (t *TVWebService) GetAny(key string) ([]byte, []string) {
	return nil, nil
}

// OnFinish handles service completion
func (t *TVWebService) OnFinish() {
	z.L().Info("TVWeb service finished")
}

// UnInstall uninstalls the service
func (t *TVWebService) UnInstall() {
	z.L().Info("Uninstalling TVWeb service...")
}
