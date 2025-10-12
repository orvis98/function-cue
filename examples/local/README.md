# Example manifests

You can run your function locally and test it using `crossplane render`
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
    message: 'Unready resources: configMap'
    reason: Creating
    status: "False"
    type: Ready
---
apiVersion: v1
data:
  compositeName: example-xr
kind: ConfigMap
metadata:
  annotations:
    crossplane.io/composition-resource-name: configMap
  labels:
    crossplane.io/composite: example-xr
  name: foo
  namespace: default
  ownerReferences:
  - apiVersion: example.crossplane.io/v1
    blockOwnerDeletion: true
    controller: true
    kind: XR
    name: example-xr
    uid: ""
---
apiVersion: render.crossplane.io/v1beta1
kind: Result
message: cue module executed successfully
severity: SEVERITY_NORMAL
step: run-the-template
```

```shell
# sample composition in an actual cluster
echo """apiVersion: apiextensions.crossplane.io/v1
kind: Composition
metadata:
  name: function-cue
spec:
  compositeTypeRef:
    apiVersion: example.crossplane.io/v1
    kind: XR
  mode: Pipeline
  pipeline:
  - step: run-the-template
    functionRef:
      name: function-cue
    input:
      apiVersion: template.fn.crossplane.io/v1beta1
      kind: Input
      debug: true
      asModule: true
      cueMod: |
        module: "cue.example"
        language: {
          version: "v0.12.0"
        }
        deps: {
          "cue.dev/x/k8s.io@v0": {
            v:       "v0.5.0"
            default: true
          }
        }
      script: |
        package main
        import (
          corev1 "cue.dev/x/k8s.io/api/core/v1"
        )
        request: {}
        response: desired: resources: {
          configMap: resource: corev1.#ConfigMap & {
            apiVersion: "v1"
            kind: "ConfigMap"
            data: compositeName: request.observed.composite.resource.metadata.name
          }
        }" | kubectl apply -f -
```
