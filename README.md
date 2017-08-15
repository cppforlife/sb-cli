# sb-cli

Tiny CLI to use with [Open Service Brokers](https://openservicebrokerapi.org/).

Usage:

```
$ source .envrc
$ bin/test

$ export SB_BROKER_URL=http://...
$ export SB_BROKER_USERNAME=...
$ export SB_BROKER_PASSWORD=...

$ out/sb-cli create-service-instance test1 service1 plan1
$ out/sb-cli delete-service-instance test1 service1 plan1
```

## Todo

- list service and plans
- better auth errors
- auto detect service and plan if only one
- uaa auth
- list instances
- set context during create
