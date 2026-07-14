package git

import (
	"testing"
)

func TestClientSatisfiesGitAPI(t *testing.T) {
	// This test verifies that *Client implements GitAPI
	var _ GitAPI = (*Client)(nil)
}
