[![Create and publish a Docker image](https://github.com/Exchizz/kube-learning/actions/workflows/go.yml/badge.svg)](https://github.com/Exchizz/kube-learning/actions/workflows/go.yml)
# What is this ? 
This is a tool for playing around with kubernetes probes.
https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/


This program (kube-learning) allows you to change the status of the probes via stdin by entering commands. Fx. "health: false" - this will cause the health check to fail. Kubernetes will then try to reschedule the pod where "kube-learning" is running.

The tool currently supports the following probes:
 - Health probe
 - Liveness probe 
 - Readyness probe

(Startup probe is not supported at the moment)

---

Run and test locally:
docker run -p 8080:8080 --rm -i -t ghcr.io/exchizz/kube-learning:master


Output from example run:

```
time="2023-03-13T16:05:12+01:00" level=info msg=-------------------------------
time="2023-03-13T16:05:12+01:00" level=info msg="Command syntax: <cmd>:<value>"   
time="2023-03-13T16:05:12+01:00" level=info msg="Server is listening on port 8080"
time="2023-03-13T16:05:12+01:00" level=info msg="Examples: "
time="2023-03-13T16:05:12+01:00" level=info msg="  Health: false"
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
  Health: true
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
# Example output from /
```
curl http://localhost:8080/
03-13-2023 16:32:49|local|default_pod_name|says hello
```
*local* is the default value of the the env-var NODE_NAME

*default_pod_name* is the default value of the env-var POD_NAME

# Commands
```
? <enter> - shows the status of all probes and where the pod is running
verbose: true | false <enter> - enables/disable verbose logging (default is off)
health: true | false <enter> - sets state of health probe
liveness: true | false <enter> - sets state of liveness probe
readyness: true | false <enter> - sets state of readyness probe
```

Ideas:
 - Show when sigterm/sigkill is received
 - [x] Command for disabling debugging (it's noisy)
 - [x] Command for printing the status of all probes
 - [x] Show current pod name and node name in "prompt" (for debugging when running multiple instances of kube-learning)
