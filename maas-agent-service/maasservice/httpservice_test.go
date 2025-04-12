package maasservice

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/netcracker/qubership-core-maas-agent/maas-agent-service/v2/httputils"
	"github.com/netcracker/qubership-core-maas-agent/maas-agent-service/v2/model"

	"github.com/netcracker/qubership-core-lib-go/v3/configloader"
	_assert "github.com/stretchr/testify/assert"
)

var (
	mockServer   *httptest.Server
	mockResponse string

	ctx                 context.Context
	httpService         MaaSService
	errorRequestCounter int
)

func setUp() {
	configloader.InitWithSourcesArray([]*configloader.PropertySource{configloader.EnvPropertySource()})
	ctx = context.Background()

	mockResponse = "dummy"
	errorRequestCounter = 0
	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != "username" || pass != "password" {
			w.WriteHeader(http.StatusForbidden)
		}

		// emulate request errors
		if errorRequestCounter > 0 {
			errorRequestCounter--
			w.WriteHeader(http.StatusInternalServerError)
		}

		fmt.Fprintln(w, mockResponse)
	}))

	mockServerUrl, _ := url.ParseRequestURI(mockServer.URL)
	httpService = MaaSService{
		CpAddr:    mockServerUrl,
		TmAddr:    mockServerUrl,
		MaasAddr:  mockServer.URL,
		Namespace: "maas-test",
		BasicRequestCreator: func(method, url string) *httputils.HttpRequest {
			return httputils.Req(method, url, model.AuthCredentials{"username", "password"}.AuthHeaderProvider)
		},
		M2MRequestCreator: func(method, url string) *httputils.HttpRequest {
			// fake M2M only for test
			return httputils.Req(method, url, model.AuthCredentials{"username", "password"}.AuthHeaderProvider)
		},
	}
}

func tearDown() {
	mockServer.Close()
}

func TestMain(m *testing.M) {
	setUp()
	exitCode := m.Run()
	tearDown()
	os.Exit(exitCode)
}

func TestApiHttpHandler_SendActiveTenants(t *testing.T) {
	assert := _assert.New(t)

	err := httpService.SendActiveTenants(ctx, "tenants")
	assert.NoError(err)
}

func TestSynchronizeTenantsToMaaS(t *testing.T) {
	//assert := _assert.New(t)
	errorRequestCounter = 1
	httpService.SynchronizeTenantsToMaaS(ctx)
	//assert.NoError(err)
}

func TestApiHttpHandler_SendCpVersionsToMaas(t *testing.T) {
	assert := _assert.New(t)

	err := httpService.SendCpVersionsToMaas(ctx, model.CpVersionsDto{})
	assert.NoError(err)
}
