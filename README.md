[![Create and publish a Docker image](https://github.com/Exchizz/kube-learning/actions/workflows/go.yml/badge.svg)](https://github.com/Exchizz/kube-learning/actions/workflows/go.yml)
# What is this ? 
This is a tool for playing around with kubernetes probes.
https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/


This program (kube-learning) allows you to change the status of the probes via stdin. Fx. "health: false" - this will cause the health check to fail. Kubernetes will then try to reschedule the pod where "kube-learning" is running.

The tool currently supports the following probes:
 - Health probe
 - Liveness probe 
 - Readyness probe

(Startup probe is not supported at the moment)

---

Run and test locally:
docker run -p 8081:8081 --rm -i -t ghcr.io/exchizz/kube-learning:master


Output from example run:

```
-------------------------------
Command syntax: <cmd>:<value>
Health: false
Readyness: false
liveness: true
-------------------------------
2023/03/10 13:06:50 172.17.0.1:58754 GET /
2023/03/10 13:06:50 172.17.0.1:58754 GET /favicon.ico
2023/03/10 13:06:53 172.17.0.1:58756 GET /healthz
2023/03/10 13:06:53 172.17.0.1:58756 GET /favicon.ico
health:false <- input from keyboard
cmd_health called with value: false
2023/03/10 13:07:00 172.17.0.1:58758 GET /healthz
2023/03/10 13:07:02 172.17.0.1:58758 GET /healthz
2023/03/10 13:07:03 172.17.0.1:58758 GET /healthz
health: true <- input from keyboard
```



Ideas:
 - Show when sigterm/sigkill is received
 - Command for disabling debugging (it's noisy)
 - Command for printing the status of all probes
 - Show current pod name and node name in "prompt" (for debugging when running multiple instances of kube-learning)
