package commands

import (
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.ibm.com/blockchaindb/server/pkg/logger"
)

func TestMint(t *testing.T) {
	demoDir, err := ioutil.TempDir("/tmp", "cars-demo-test")
	require.NoError(t, err)
	defer os.RemoveAll(demoDir)

	err = Generate(demoDir)
	require.NoError(t, err)

	testServer, _, err := setupTestServer(t, demoDir)
	require.NoError(t, err)
	defer func() {
		if testServer != nil {
			_ = testServer.Stop()
		}
	}()
	require.NoError(t, err)
	err = testServer.Start()
	require.NoError(t, err)

	serverPort, err := testServer.Port()
	require.NoError(t, err)

	serverUrl, err := url.Parse("http://127.0.0.1:" + serverPort)
	require.NoError(t, err)

	c := &logger.Config{
		Level:         "debug",
		OutputPath:    []string{"stdout"},
		ErrOutputPath: []string{"stderr"},
		Encoding:      "console",
		Name:          "bcdb-client",
	}
	logger, err := logger.New(c)

	err = Init(demoDir, serverUrl, logger)
	require.NoError(t, err)

	out, err := MintRequest(demoDir, "dealer", "Test.Car.1", logger)
	require.NoError(t, err)
	require.Contains(t, out, "MintRequest: committed")

	index := strings.Index(out, "Key:")
	mintRequestKey := strings.TrimSpace(out[index+4:])
	require.True(t, strings.HasPrefix(mintRequestKey, "mint-request~"))

	out, err = MintApprove(demoDir, "dmv", mintRequestKey, logger)
	require.NoError(t, err)
	require.Contains(t, out, "MintApprove: committed")

	index = strings.Index(out, "Key:")
	carKey := strings.TrimSpace(out[index+4:])
	require.True(t, strings.HasPrefix(carKey, "car~"))
}
