package controllers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBranchNameOnRemote(t *testing.T) {
	assert.Equal(t, "main", branchNameOnRemote("refs/heads/main"))
	assert.Equal(t, "main", branchNameOnRemote("refs/remotes/origin/main"))
	assert.Equal(t, "feature/x", branchNameOnRemote("refs/remotes/origin/feature/x"))
}
