package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"

	"github.com/golang/glog"
	"github.com/stackrox/acs-fleet-manager/probe/config"
	"github.com/stackrox/acs-fleet-manager/probe/pkg/metrics"
	"github.com/stackrox/acs-fleet-manager/probe/pkg/runtime"
	"golang.org/x/sys/unix"
	// "github.com/spf13/cobra"
)

func main() {
	// This is needed to make `glog` believe that the flags have already been parsed, otherwise
	// every log messages is prefixed by an error message stating the the flags haven't been
	// parsed.
	_ = flag.CommandLine.Parse([]string{})

	// Always log to stderr by default, required for glog.
	if err := flag.Set("logtostderr", "true"); err != nil {
		glog.Info("Unable to set logtostderr to true.")
	}

	// TODO: CLI commands for continuous and one-shot runs.
	// rootCmd := &cobra.Command{
	// 	Use:  "probe",
	// 	Long: "probe is a service that verifies the availability of ACS fleet manager.",
	// }

	config, err := config.GetConfig()
	if err != nil {
		glog.Fatalf("Failed to load configuration: %v", err)
	}

	glog.Infof("Probe service started.")

	metricsServer := metrics.NewMetricsServer(config.MetricsAddress)
	defer cleanupMetricsServer(metricsServer)
	go func() {
		if err := metricsServer.ListenAndServe(); err != nil {
			glog.Errorf("Failed to serve metrics: %v", err)
		}
	}()

	runtime, err := runtime.NewRuntime(config)
	if err != nil {
		glog.Fatal(err)
	}
	defer cleanupRuntime(runtime)
	go func() {
		err := runtime.Run()
		if err != nil {
			glog.Fatal(err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, unix.SIGTERM)

	sig := <-sigs
	glog.Infof("Caught %s signal.", sig)
}

func cleanupRuntime(runtime *runtime.Runtime) {
	err := runtime.Stop()
	if err != nil {
		glog.Fatal(err)
	}
	glog.Info("Probe service has been stopped.")
}

func cleanupMetricsServer(metricsServer *http.Server) {
	if err := metricsServer.Close(); err != nil {
		glog.Errorf("Failed to close metrics server: %v", err)
	}
}
