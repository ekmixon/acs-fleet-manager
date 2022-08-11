package reconciler

import private "github.com/stackrox/acs-fleet-manager/generated/privateapi"

func readyStatus() *private.CentralStatus {
	return &private.CentralStatus{
		Conditions: []*private.Condition{
			{
				Type:   "Ready",
				Status: "True",
			},
		},
	}
}

func deletedStatus() *private.CentralStatus {
	return &private.CentralStatus{
		Conditions: []*private.Condition{
			{
				Type:   "Ready",
				Status: "False",
				Reason: "Deleted",
			},
		},
	}
}

func installingStatus() *private.CentralStatus {
	return &private.CentralStatus{
		Conditions: []*private.Condition{
			{
				Type:   "Ready",
				Status: "False",
				Reason: "Installing",
			},
		},
	}
}
