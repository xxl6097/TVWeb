package service

import (
	"github.com/xxl6097/TVBoxWeb/pkg/tvbox"
	"github.com/xxl6097/glog/pkg/z"
	"go.uber.org/zap"
)

// TVWebService represents the TVWeb service
type TVWebService struct {
	tvbox *tvbox.TVBox
}

// Start starts the TVWeb service
func (t *TVWebService) Start() error {
	z.L().Info("Starting TVWeb service...")

	// Initialize TVBox
	var err error
	t.tvbox, err = tvbox.NewTVBox()
	if err != nil {
		z.L().Error("Failed to initialize TVBox", zap.Error(err))
		return err
	}

	// Start the TVBox
	if err := t.tvbox.Start(); err != nil {
		z.L().Error("Failed to start TVBox", zap.Error(err))
		return err
	}

	z.L().Info("TVWeb service started successfully")
	return nil
}

// Stop stops the TVWeb service
func (t *TVWebService) Stop() error {
	z.L().Info("Stopping TVWeb service...")

	if t.tvbox != nil {
		if err := t.tvbox.Stop(); err != nil {
			z.L().Error("Failed to stop TVBox", zap.Error(err))
			return err
		}
	}

	z.L().Info("TVWeb service stopped")
	return nil
}

// GetName returns the service name
func (t *TVWebService) GetName() string {
	return "TVWebService"
}
