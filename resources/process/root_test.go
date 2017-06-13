package process

import (
	"context"
	"net"
	"strconv"
	"testing"

	"net/http/httptest"
	"net/url"

	"github.com/deviceio/hmapi"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type test_root_when_creating_process struct {
	suite.Suite
}

func (t *test_root_when_creating_process) Test_process_successfully_created() {
	objects := t.getTestObjects()
	defer objects.server.Close()

	resp, err := objects.client.
		Resource("/process").
		Form("create").
		AddFieldAsString("cmd", "whoami").
		Submit(context.Background())

	assert.NotNil(t.T(), resp)
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), 201, resp.StatusCode)

	process, err := objects.client.
		Resource(resp.Header.Get("Location")).
		Get(context.Background())

	assert.NotNil(t.T(), process)
	assert.Nil(t.T(), err)
}

func (t *test_root_when_creating_process) getTestObjects() (ret struct {
	root   *root
	server *httptest.Server
	client hmapi.Client
}) {
	router := mux.NewRouter()
	RegisterRoutes(router)

	ret.server = httptest.NewServer(router)

	url, _ := url.Parse(ret.server.URL)
	hoststr, portstr, _ := net.SplitHostPort(url.Host)
	port, _ := strconv.ParseInt(portstr, 10, 0)

	ret.client = hmapi.NewClient(&hmapi.ClientConfig{
		Auth:   &hmapi.AuthNone{},
		Host:   hoststr,
		Port:   int(port),
		Scheme: hmapi.HTTP,
	})

	return
}

func TestRunRootTestSuites(t *testing.T) {
	suite.Run(t, new(test_root_when_creating_process))
}
