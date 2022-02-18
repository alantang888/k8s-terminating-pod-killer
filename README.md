# K8S Terminating Pod Killer

Some rare case. Pod can stucking in terminating status. 
Even over `terminationGracePeriodSeconds` those pod still can stucking here. 
This app is for kill that kind of pod.

## Environment Variable
- `NAMESPACE`: Which namespace this app need to check and kill pod. Default `""` will check all namespaces
- `KILL_MINUTE`: How long time in minute after `terminationGracePeriodSeconds` should kill that pod
