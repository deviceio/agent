package process

import (
	"context"
	"net"
	"net/http/httptest"
	"net/url"
	"runtime"
	"strconv"
	"testing"

	"os/exec"

	"io/ioutil"

	"github.com/deviceio/hmapi"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type test_process_when_calling_stdout struct {
	suite.Suite
}

func (t *test_process_when_calling_stdout) Test_full_stdout_recieved() {
	objects := t.getTestObjects()
	defer objects.server.Close()

	var cmdstr string
	var args []string

	if runtime.GOOS == "windows" {
		cmdstr = "cmd"
		args = []string{
			"/c",
			"dir",
			"c:\\windows",
		}
	} else {
		cmdstr = "dir"
		args = []string{
			"/bin",
		}
	}

	cmd := exec.Command(cmdstr, args...)

	output, err := cmd.Output()

	if err != nil {
		t.T().Fatalf("error running command: %v", err.Error())
	}

	form := objects.client.
		Resource("/process").
		Form("create").
		AddFieldAsString("cmd", cmdstr)

	for _, arg := range args {
		form.AddFieldAsString("arg", arg)
	}

	resp, err := form.Submit(context.Background())

	assert.NotNil(t.T(), resp)
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), 201, resp.StatusCode)

	stdout, err := objects.client.
		Resource(resp.Header.Get("Location")).
		Link("stdout").
		Get(context.Background())

	assert.NotNil(t.T(), stdout)
	assert.Nil(t.T(), err)

	start, err := objects.client.
		Resource(resp.Header.Get("Location")).
		Form("start").
		Submit(context.Background())

	assert.NotNil(t.T(), start)
	assert.Nil(t.T(), err)

	stdoutdata, err := ioutil.ReadAll(stdout.Body)

	assert.NotNil(t.T(), stdoutdata)
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), len(output), len(stdoutdata))
	assert.Equal(t.T(), output, stdoutdata)

}

func (t *test_process_when_calling_stdout) getTestObjects() (ret struct {
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

func TestRunProcessTestSuite(t *testing.T) {
	suite.Run(t, new(test_process_when_calling_stdout))
}
