package application_test

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	app "github.com/sevaho/theonepager/src"
	"github.com/sevaho/theonepager/src/environment"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var client = resty.New()
var port = 30001
var testApp = fmt.Sprintf("http://localhost:%d", port)

var _ = Context("Application", func() {
	var ctx context.Context
	var cancel context.CancelFunc
	var wg sync.WaitGroup
	var env *environment.Environment

	BeforeEach(func() {
		// setup the environment
		env = environment.New()
		env.LOG_LEVEL = 4
		env.IS_DEVELOPMENT = false // Force use of embedded templates for tests
		env.CONFIG_FILE_PATH = "config.yaml"

		// setup application
		ctx, cancel = context.WithCancel(context.Background())
		wg.Add(1)
		go func() { defer wg.Done(); app.Run(ctx, port, env) }()

		// Wait till application is ready to test
		waitForReady(ctx, time.Second*2, testApp+"/healthz")
	})
	AfterEach(func() {
		cancel()
		wg.Wait()
	})

	JustBeforeEach(func() {
		httpmock.Activate()
	})
	JustAfterEach(func() {
		httpmock.DeactivateAndReset()
	})

	When("Fetching the index page", func() {
		It("should render the correct html page", func() {
			// when
			res, err := client.R().Get(testApp)

			// then
			Expect(err).To(BeNil())
			Expect(res.StatusCode()).To(Equal(200), res.String())
		})
	})

})
