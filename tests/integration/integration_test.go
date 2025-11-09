package integration_test

import (
	"testing"

	"github.com/abibby/icbmdb/test"
)

func TestIntegration(t *testing.T) {
	test.Kernel(t).
		GetJSON("/api/user").
		AssertStatusOK().
		AssertJSONString(`{
			"users": []
		}`)
}
