package integration

import (
	"context"
	"testing"
	"net/http"
	"net/http/httptest"

	g "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	
	"github.com/stackrox/acs-fleet-manager/internal/rhacs"
	"github.com/stackrox/acs-fleet-manager/internal/rhacs/internal/api/public"
	"github.com/stackrox/acs-fleet-manager/test"
	"github.com/stackrox/acs-fleet-manager/test/mocks"
	itest "github.com/stackrox/acs-fleet-manager/internal/rhacs/test"
)

// Hack as in https://github.com/operator-framework/operator-sdk/issues/2913 and suggested in https://github.com/onsi/ginkgo/issues/9
var Testing *testing.T

func TestCentralResource(t *testing.T) {
	RegisterFailHandler(g.Fail)
	Testing = t
	g.RunSpecs(t, "API resources Suite")
}

var _ = g.Describe("API resources", func() {
	var (
		t   *testing.T
		ocmServer *httptest.Server
		teardown func()
		client *public.APIClient
		ctx context.Context
	)

	g.BeforeEach(func() {
		t = Testing
		ocmServer = mocks.NewMockConfigurableServerBuilder().Build()
		var testHelper *test.Helper
		testHelper, teardown = test.NewHelperWithHooks(t, ocmServer, nil, rhacs.ConfigProviders())

		client = itest.NewApiClient(testHelper)

		account := testHelper.NewRandAccount()
		ctx = testHelper.NewAuthenticatedContext(account, nil)
	})
	g.AfterEach(func() {
		ocmServer.Close()
		teardown()
	})

	g.Context("Central API resource", func() {
		g.Context("the database is empty, and an authorized user makes a request", func() {
			g.When("listing all centrals", func() {
				g.It("returns an empty centrals list", func() {
					centrals, response, err := client.DefaultApi.GetCentrals(ctx, nil)
					Expect(err).NotTo(HaveOccurred(), "Error listing centrals: %v", err)
					Expect(response.StatusCode).To(Equal(http.StatusOK))
					Expect(centrals.Kind).To(Equal("CentralRequestList"))
					Expect(centrals.Items).To(BeEmpty())
					Expect(centrals.NextPageCursor).To(BeEmpty())
					Expect(centrals.Size).To(BeEquivalentTo(0))
				})
			})
	
			g.When("getting a missing central", func() {
				g.It("returns not found error", func() {
					_, response, err := client.DefaultApi.GetCentralById(ctx, "missing ID")
					Expect(err).Error()
					Expect(response.StatusCode).To(Equal(http.StatusNotFound))
				})	
			})
		})
	})

	g.Context("Root API resources", func() {
		g.When("getting the version metadata", func() {
			g.It("returns the expected version and collections", func()  {
				metadata, response, err := client.DefaultApi.GetVersionMetadata(ctx)
				Expect(err).NotTo(HaveOccurred(), "Error getting version metadata: %v", err)
				Expect(response.StatusCode).To(Equal(http.StatusOK))
				Expect(metadata.Id).To(Equal("v1"))
				Expect(metadata.Kind).To(Equal("APIVersion"))
				Expect(metadata.Href).To(Equal("/api/rhacs/v1"))
				Expect(metadata.Collections).Should(ContainElements(public.ObjectReference{
					Id: "centrals",
					Kind: "CentralList",
					Href: "/api/rhacs/v1/centrals",
				}))
			})
		})

		g.When("getting the service status", func() {
			g.It("returns capacity is not reached", func() {
				status, response, err := client.DefaultApi.GetServiceStatus(ctx)
				Expect(err).NotTo(HaveOccurred(), "Error getting version metadata: %v", err)
				Expect(response.StatusCode).To(Equal(http.StatusOK))
				Expect(status.Centrals.MaxCapacityReached).To(BeFalse())
			})
		})
	})
})
