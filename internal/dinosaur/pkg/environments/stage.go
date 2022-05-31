package environments

import "github.com/stackrox/acs-fleet-manager/pkg/environments"

func NewStageEnvLoader() environments.EnvLoader {
	return environments.SimpleEnvLoader{
		"ocm-base-url":                         "https://api.stage.openshift.com",
		"ams-base-url":                         "https://api.stage.openshift.com",
		"enable-ocm-mock":                      "false",
		"enable-deny-list":                     "true",
		"max-allowed-instances":                "1",
		"sso-base-url":                         "https://identity.api.stage.openshift.com",
		"enable-dinosaur-external-certificate": "true",
		"cluster-compute-machine-type":         "m5.2xlarge",
	}
}
