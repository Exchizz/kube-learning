[![Create and publish a Docker image](https://github.com/Exchizz/kube-learning/actions/workflows/go.yml/badge.svg)](https://github.com/Exchizz/kube-learning/actions/workflows/go.yml)

# What is this ?

This is a tool for playing around with kubernetes probes.
https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/

This program (kube-learning) allows you to change the status of the probes via stdin by entering commands. Fx. "health: false" - this will cause the health check to fail. Kubernetes will then try to reschedule the pod where "kube-learning" is running.

The tool currently supports the following probes:

- Liveness probe
- Readyness probe

(Startup probe is not supported at the moment)

---

Run and test locally:
docker run -p 8080:8080 --rm -i -t ghcr.io/exchizz/kube-learning:master

Run and test in kubernetes:

## Installation

```
kubectl apply -f https://raw.githubusercontent.com/Exchizz/kube-learning/master/manifest/deployment.yaml
```

_NOTE_ This installation does not create service for the pod nor is an ingress object created. If you want to play around with readiness-probes, you need to expose the pod via an ingress gateway.

## Test

Create proxy from your PC to the pod running in your cluster:

```
kubectl port-forward deployment/kube-learning-deployment  8080:8080
```

In another terminal:

```
kubectl attach -i -t deployment/kube-learning-deployment
```

If you don't see any output type ? and press enter.

visit https://localhost:8080 in your browser and see the output from the pod

Output from example run:

```
time="2023-03-13T16:05:12+01:00" level=info msg=-------------------------------
time="2023-03-13T16:05:12+01:00" level=info msg="Command syntax: <cmd>:<value>"
time="2023-03-13T16:05:12+01:00" level=info msg="Server is listening on port 8080"
time="2023-03-13T16:05:12+01:00" level=info msg="Examples: "
time="2023-03-13T16:05:12+01:00" level=info msg="  Readyness: false"
time="2023-03-13T16:05:12+01:00" level=info msg="  liveness: true"
time="2023-03-13T16:05:12+01:00" level=info msg=-------------------------------
verbose: true  <- input
time="2023-03-13T16:05:17+01:00" level=info msg="Enabling debug log"
time="2023-03-13T16:05:18+01:00" level=debug msg="127.0.0.1:65512 GET /" node_name=local pod_name=default_pod_name
time="2023-03-13T16:05:19+01:00" level=debug msg="127.0.0.1:65512 GET /" node_name=local pod_name=default_pod_name
time="2023-03-13T16:05:19+01:00" level=debug msg="127.0.0.1:65512 GET /" node_name=local pod_name=default_pod_name
verbose: false  <- input
time="2023-03-13T16:05:32+01:00" level=info msg="Disabling debug log"
?  <- input
Status:
  Liveness: false
  Readyness: false
  Node name: local
  Pod name: default_pod_name
health: false  <- input
?  <- input
Status:
  Health: false
  Liveness: false
  Readyness: false
  Node name: local
  Pod name: default_pod_name
```

# Deploy to kubernetes

In order for kube-learning to know its pod-name and hostname, you need to pass that information to the pod as env vars when deploying the pod.
For more information: https://raw.githubusercontent.com/kubernetes/website/main/content/en/examples/pods/inject/dapi-envars-pod.yaml

```
      env:
        - name: NODE_NAME
          valueFrom:
            fieldRef:
              fieldPath: spec.nodeName
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
```

The following env variables allow further customization:

- `KUBELEARN_ALIVE`: sets the starting value of the liveness probe
- `KUBELEARN_READY`: sets the starting value of the readiness probe
- `KUBELEARN_DEBUG`: sets the starting value of the debug mode
- `KUBELEARN_DEBUG_FORMAT`: sets the format of the debug logs (using a [text.template](https://pkg.go.dev/text/template) template)

# Example output from /

```
curl http://localhost:8080/
03-13-2023 16:32:49|local|default_pod_name|says hello
```

_local_ is the default value of the the env-var NODE_NAME

_default_pod_name_ is the default value of the env-var POD_NAME

# Commands

```
? <enter> - shows the status of all probes and where the pod is running
verbose: true | false <enter> - enables/disable verbose logging (default is off)
liveness: true | false <enter> - sets state of liveness probe
readyness: true | false <enter> - sets state of readyness probe
```

Ideas:

- Show when sigterm/sigkill is received
