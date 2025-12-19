package application_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/jarcoal/httpmock"
	app "github.com/sevaho/theonepager/src"
	"github.com/sevaho/theonepager/src/environment"
	"github.com/sevaho/theonepager/src/pkg/resty"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var port = 30001
var testApp = fmt.Sprintf("http://localhost:%d", port)

var testConfig = `---
applications:
  1. Test:
    - name: "Test App"
      description: "Test application"
      link: https://example.com
    - name: "Test App No Icon"
      description: "Test app without icon"
      link: https://noicon.example.com
`

var _ = Context("Application", func() {
	var ctx context.Context
	var cancel context.CancelFunc
	var wg sync.WaitGroup
	var env *environment.Environment
	var configFile *os.File

	BeforeEach(func() {
		configFile, err := ioutil.TempFile("/tmp", "config_*.yaml")
		if err != nil {
			Fail(fmt.Sprintf("Failed to create temp file: %v", err))
		}

		_, err = configFile.WriteString(testConfig)
		if err != nil {
			Fail(fmt.Sprintf("Failed to write config: %v", err))
		}
		configFile.Close()

		// setup the environment
		env = environment.New()
		env.LOG_LEVEL = 4
		env.IS_DEVELOPMENT = false // Force use of embedded templates for tests
		env.CONFIG_FILE_PATH = configFile.Name()

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

		// Clean up temporary file
		if configFile != nil {
			os.Remove(configFile.Name())
		}
	})

	JustBeforeEach(func() {
		httpmock.ActivateNonDefault(resty.Client.GetClient())

		// Allow localhost requests to pass through
		httpmock.RegisterNoResponder(func(req *http.Request) (*http.Response, error) {
			host := req.URL.Hostname()
			if host == "localhost" || host == "127.0.0.1" || host == "" {
				return httpmock.InitialTransport.RoundTrip(req)
			}
			return nil, fmt.Errorf("unmocked external URL: %s", req.URL)
		})
	})
	JustAfterEach(func() {
		httpmock.DeactivateAndReset()
	})

	When("Fetching the index page", func() {
		It("should render the correct html page", func() {
			// when
			res, err := resty.Client.R().Get(testApp)

			// then
			Expect(err).To(BeNil())
			Expect(res.StatusCode()).To(Equal(200), res.String())
		})
	})

	When("Fetching the healthz endpoint", func() {
		It("should return OK status", func() {
			// when
			res, err := resty.Client.R().Get(testApp + "/healthz")

			// then
			Expect(err).To(BeNil())
			Expect(res.StatusCode()).To(Equal(200))
			Expect(res.String()).To(Equal("OK"))
		})
	})

	When("Using the icon cache proxy", func() {
		It("should cache external resources and only call external URL once", func() {
			testURL := "https://example.com/icon.png"
			mockResponseBody := []byte("fake-image-data")

			httpmock.RegisterResponder("GET", testURL,
				httpmock.NewBytesResponder(200, mockResponseBody))

			// when: doing a request to the icon cache
			res1, err := resty.Client.R().Get(testApp + "/v1/api/iconcache?link=" + testURL)
			Expect(err).To(BeNil())
			Expect(res1.StatusCode()).To(Equal(200))
			Expect(res1.Body()).To(Equal(mockResponseBody))

			// when: doing another request to the icon cache
			res2, err := resty.Client.R().Get(testApp + "/v1/api/iconcache?link=" + testURL)
			Expect(err).To(BeNil())
			Expect(res2.StatusCode()).To(Equal(200))
			Expect(res2.Body()).To(Equal(mockResponseBody))

			// then: the callcount should still be 1
			Expect(httpmock.GetCallCountInfo()["GET "+testURL]).To(Equal(1))
		})
	})

})
