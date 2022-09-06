package upgrade

import (
	"errors"

	"github.com/golang/glog"
	"github.com/stackrox/acs-fleet-manager/fleetshard/pkg/central/reconciler"
)

// Upgrader provides methods to upgrade all central instances running on a dataplane cluster
type Upgrader struct {
}

// Upgrade starts the upgrade process
func (u *Upgrader) Upgrade(reconcilers []*reconciler.CentralReconciler, centralVersion string) {
	// all centrals get updated with a pause annotation
	// wait for reconcilers to stop reconciling
	for {
		readyForUpgrade := true
		for _, r := range reconcilers {
			if r.IsBusy() {
				readyForUpgrade = false
			}
		}

		if readyForUpgrade {
			break
		}
	}

	for _, r := range reconcilers {
		if err := r.Pause(); err != nil {
			glog.Errorf("applying pause annotations: %v", err)
			if err := u.Rollback(reconcilers); err != nil {
				glog.Errorf("rolling back pause annotations: %v", err)
			}
			break
		}
	}

	// Upgrade operator

	// remove pause annotation
	for _, r := range reconcilers {
		if err := r.Unpause(); err != nil {
			glog.Errorf("removing pause annotations: %v", err)
		}
	}
}

// Rollback reverts all operation done throughout an upgrade process
func (u *Upgrader) Rollback(reconcilers []*reconciler.CentralReconciler) error {
	return errors.New("not implemented")
}
