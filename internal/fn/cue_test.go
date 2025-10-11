package fn

import (
	"os"
	"testing"

	"cuelang.org/go/cue/load"
	fnv1 "github.com/crossplane/function-sdk-go/proto/v1"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	script := ""
	f, err := New(Options{})
	require.NoError(t, err)
	tmpDir, err := os.MkdirTemp("", "test-function-cue-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	cfg, err := f.NewConfig(tmpDir, script, "", "", "", "", false)
	require.NoError(t, err)
	insts := load.Instances([]string{}, cfg)
	require.Len(t, insts, 1)
}

func TestEvaluate(t *testing.T) {
	script := `
	package main
	response: desired: resources: {}
	`
	var req *fnv1.RunFunctionRequest
	f, err := New(Options{})
	require.NoError(t, err)
	tmpDir, err := os.MkdirTemp("", "test-function-cue-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)
	cfg, err := f.NewConfig(tmpDir, script, "", "", "", "", false)
	require.NoError(t, err)
	_, err = f.Evaluate(req, cfg, "", "")
	require.NoError(t, err)
}
