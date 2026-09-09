package integration_test

import (
	"testing"

	"abibby.com/mangadb/test"
)

func TestIntegration(t *testing.T) {
	test.Kernel(t).
		GetJSON("/api/user").
		AssertStatusOK().
		AssertJSONString(`{
			"users": []
		}`)
}
