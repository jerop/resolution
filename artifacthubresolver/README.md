# Artifact Hub Resolver

Remote Resource Resolver of type `artifacthub`.

## Parameters

| Param Name       | Description                                                                 | Example Value  |
|------------------|-----------------------------------------------------------------------------|----------------|
| `kind`           | Either `task` or `pipeline`.                                                | `task`         |
| `name`           | The name of the task or pipeline to fetch from the hub.                     | `golang-build` |
| `version`        | Version of task or pipeline to pull in from hub. Wrap the number in quotes! | `"0.5.0"`      |

## Getting Started

### Requirements

See the [getting started instructions](https://github.com/tektoncd/resolution/tree/main/docs/getting-started.md) in 
the Tekton Resolution repo.

### Install

1. Install the Hub resolver:

```bash
$ ko apply -f ./config
```

### ResolutionRequest

Try creating a `ResolutionRequest` for an `artifacthub` entry:

```bash
$ cat <<EOF > rrtest.yaml
apiVersion: resolution.tekton.dev/v1alpha1
kind: ResolutionRequest
metadata:
  name: fetch-artifacthub-entry
  labels:
    resolution.tekton.dev/type: artifacthub
spec:
  params:
    kind: task
    name: git-clone
    version: "0.5"
    kind: task
EOF

$ kubectl apply -f ./rrtest.yaml

$ kubectl get resolutionrequest -w fetch-artifacthub-entry
```

You should shortly see the `ResolutionRequest` succeed and the content of the `git-clone` yaml base64-encoded in the
object's `status.data` field.

### TaskRun

This is an example `TaskRun` using `artifacthub` resolver:

```yaml
apiVersion: tekton.dev/v1beta1
kind: TaskRun
metadata:
  generateName: ah-curl-
spec:
  taskRef:
    resolver: artifacthub
    resource:
      - name: kind
        value: "task"
      - name: name
        value: "curl"
      - name: version
        value: "0.1.0"
  params:
    - name: url
      value: "www.google.com"
```

If you execute this `TaskRun`, it should resolve the `Task` and run successfully.
