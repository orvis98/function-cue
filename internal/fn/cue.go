package fn

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	fnv1 "github.com/crossplane/function-sdk-go/proto/v1"
	"github.com/pkg/errors"
)

func (f *Cue) NewConfig(dir string, script string, moduleCue string, cacheDir string, registry string, requestVar string, debug bool) (*load.Config, error) {
	// set up variables
	if moduleCue == "" {
		moduleCue = `
module: "cue.example"
language: version: "v0.12.0"`
	}
	if cacheDir == "" {
		cacheDir = "/tmp/.cuecache"
	}
	if registry == "" {
		registry = "registry.cue.works"
	}
	if requestVar == "" {
		requestVar = "#request"
	}

	// set up CUE
	if err := os.Setenv("CUE_CACHE_DIR", filepath.Join(cacheDir, registry)); err != nil {
		return nil, err
	}
	if err := os.Setenv("CUE_REGISTRY", registry); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Join(dir, "cue.mod"), 0755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(dir, "cue.mod/module.cue"), []byte(moduleCue), 0644); err != nil {
		return nil, err
	}
	finalScript := fmt.Sprintf("%s\n%s: _", script, requestVar)
	if err := os.WriteFile(filepath.Join(dir, "script.cue"), []byte(finalScript), 0644); err != nil {
		return nil, err
	}

	if debug {
		log.Printf("[config:begin]\nCUE_CACHE_DIR: \"%s\"\nCUE_REGISTRY: \"%s\"\n[config:end]", cacheDir, registry)
		log.Printf("[module:begin]\n%s\n[module:end]", moduleCue)
		log.Printf("[script:begin]\n%s\n[script:end]", finalScript)
	}

	return &load.Config{ModuleRoot: dir, Dir: dir}, nil
}

func (f *Cue) Evaluate(req *fnv1.RunFunctionRequest, config *load.Config, requestVar string, responseVar string) (*fnv1.RunFunctionResponse, error) {
	insts := load.Instances([]string{}, config)
	if len(insts) == 0 {
		return nil, errors.New("no CUE instances found")
	}
	runtime := cuecontext.New()
	val := runtime.BuildInstance(insts[0]).FillPath(cue.ParsePath(requestVar), req)
	if val.Err() != nil {
		return nil, errors.Wrap(val.Err(), "compile cue code")
	}
	var ret fnv1.RunFunctionResponse
	if err := val.Decode(&ret); err != nil {
		return nil, errors.Wrap(err, "decode")
	}
	return &ret, nil
}
