package clustermanager

import (
	"context"
	"github.com/vahidmostofi/acfg/internal/configuration"
	"io"
)

type ClusterManager interface{
	WaitAllDeploymentsAreStable(ctx context.Context)
	Deploy(ctx context.Context, reader io.Reader) error
	UpdateConfigurationsAndWait(ctx context.Context, config map[string]*configuration.Configuration) error
}
