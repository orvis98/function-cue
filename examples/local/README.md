# Example manifests

You can run your function locally and test it using `crossplane beta render`
with these example manifests.

```shell
# Run the function locally
$ go run . --insecure --debug
```

```shell
# Then, in another terminal, call it with these example manifests
$ crossplane render examples/local/xr.yaml examples/local/composition.yaml examples/local/functions.yaml -r
---
apiVersion: example.crossplane.io/v1
kind: XR
metadata:
  name: example-xr
status:
  conditions:
  - lastTransitionTime: "2024-01-01T00:00:00Z"
    reason: Available
    status: "True"
    type: Ready
---
apiVersion: render.crossplane.io/v1beta1
kind: Result
message: cue module executed successfully
severity: SEVERITY_NORMAL
step: run-the-template
```
