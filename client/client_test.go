package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var expectedClient = ServiceBusClient{
	connectionString: "AZURE_CONNECTION_STRING",
	timeout:          time.Duration(10),
}

func TestNew(t *testing.T) {
	actualClient := New("AZURE_CONNECTION_STRING", time.Duration(10))

	// Assert that both expected and actual are the same
	assert.Equal(t, &expectedClient, actualClient)
}
