package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/stackrox/acs-fleet-manager/internal/dinosaur/pkg/api/private"
	"github.com/stackrox/acs-fleet-manager/internal/dinosaur/pkg/api/public"
	"golang.org/x/sys/unix"
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"net/http"
	"os"
	"os/signal"
)

const ClusterID = "1234567890abcdef1234567890abcdef"

var (
	URLGetCentrals      = fmt.Sprintf("http://127.0.0.1:8000/api/rhacs/v1/agent-clusters/%s/dinosaurs", ClusterID)
	URLPutCentralStatus = fmt.Sprintf("http://127.0.0.1:8000/api/rhacs/v1/agent-clusters/%s/dinosaurs/status", ClusterID)
)

/**
- 1. setting up fleet-manager
- 2. calling API to get Centrals/Dinosaurs
- 3. Applying Dinosaurs into the Kubernetes API
- 4. Implement polling
- 5. Report status to fleet-manager
*/
func main() {
	// This is needed to make `glog` believe that the flags have already been parsed, otherwise
	// every log messages is prefixed by an error message stating the the flags haven't been
	// parsed.
	_ = flag.CommandLine.Parse([]string{})

	// Always log to stderr by default
	if err := flag.Set("logtostderr", "true"); err != nil {
		glog.Info("Unable to set logtostderr to true")
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, unix.SIGTERM)

	glog.Info("fleetshard application has been started")

	synchronize()

	sig := <-sigs
	glog.Infof("Caught %s signal", sig)
	glog.Info("fleetshard application has been stopped")
}

func synchronize() {
	ocmToken := os.Getenv("OCM_TOKEN")
	if ocmToken == "" {
		glog.Fatal("empty ocm token")
	}

	buf := bytes.Buffer{}
	r, err := http.NewRequest(http.MethodGet, URLGetCentrals, &buf)
	if err != nil {
		glog.Fatal(err)
	}

	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ocmToken))
	// TODO(???): Support pagination
	client := http.Client{}

	glog.Info("Calling the Fleet Manager to get the list of Centrals")

	resp, err := client.Do(r)
	if err != nil {
		glog.Fatal(err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Fatal(err)
	}

	glog.Infof("GOT RESPONSE: %s\n\n", string(respBody))

	list := private.ManagedDinosaurList{}
	err = json.Unmarshal(respBody, &list)
	if err != nil {
		glog.Fatal(err)
	}

	statuses := make(map[string]private.DataPlaneDinosaurStatus)
	for _, v := range list.Items {
		glog.Infof("received cluster: %s", v.Metadata.Name)
		statuses[v.Id] = private.DataPlaneDinosaurStatus{
			Conditions: []private.DataPlaneClusterUpdateStatusRequestConditions{
				{
					Type:   "Ready",
					Status: "True",
				},
			},
		}
	}

	glog.Infof("Calling to update %d statuses %q", len(statuses), URLPutCentralStatus)

	updateBody, err := json.Marshal(statuses)
	if err != nil {
		glog.Fatal(err)
	}

	bufUpdateBody := &bytes.Buffer{}
	bufUpdateBody.Write(updateBody)
	updateReq, err := http.NewRequest(http.MethodPut, URLPutCentralStatus, bufUpdateBody)
	if err != nil {
		glog.Fatal(err)
	}

	updateReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ocmToken))
	resp, err = client.Do(updateReq)
	if err != nil {
		glog.Fatal(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)

	glog.Infof(string(body))
}

type ClusterReconciliator struct {
	client kubernetes.Interface
}

func (s ClusterReconciliator) reconcileCentralDeployment(r public.CentralRequest) {
	//secret := &coreV1.Secret{}
	//key := client.ObjectKey{Namespace: r.Namespace(), Name: name}
	//if err := r.Client().Get(ctx, key, secret); err != nil {
	//	if !apiErrors.IsNotFound(err) {
	//		return errors.Wrapf(err, "checking existence of %s secret", name)
	//	}
	//	secret = nil
	//}
}
