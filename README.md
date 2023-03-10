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



Ideas:
 - Show when sigterm/sigkill is received
 - Command for disabling debugging (it's noisy)
 - Command for printing the status of all probes
 - Show current pod name and node name in "prompt" (for debugging when running multiple instances of kube-learning)
