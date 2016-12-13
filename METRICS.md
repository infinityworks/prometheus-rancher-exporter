# Example Metrics

```
# HELP rancher_host_state State of defined host as reported by the Rancher API
# TYPE rancher_host_state gauge
rancher_host_state{name="server.example.com",rancherURL="http://1.1.1.1:8080/v2-beta",state="activating"} 0
rancher_host_state{name="server.example.com",rancherURL="http://1.1.1.1:8080/v2-beta",state="active"} 1
rancher_host_state{name="server.example.com",rancherURL="http://1.1.1.1:8080/v2-beta",state="deactivating"} 0
rancher_host_state{name="server.example.com",rancherURL="http://1.1.1.1:8080/v2-beta",state="error"} 0
rancher_host_state{name="server.example.com",rancherURL="http://1.1.1.1:8080/v2-beta",state="erroring"} 0
rancher_host_state{name="server.example.com",rancherURL="http://1.1.1.1:8080/v2-beta",state="inactive"} 0
rancher_host_state{name="server.example.com",rancherURL="http://1.1.1.1:8080/v2-beta",state="provisioned"} 0
rancher_host_state{name="server.example.com",rancherURL="http://1.1.1.1:8080/v2-beta",state="purged"} 0
rancher_host_state{name="server.example.com",rancherURL="http://1.1.1.1:8080/v2-beta",state="purging"} 0
rancher_host_state{name="server.example.com",rancherURL="http://1.1.1.1:8080/v2-beta",state="registering"} 0
rancher_host_state{name="server.example.com",rancherURL="http://1.1.1.1:8080/v2-beta",state="removed"} 0
rancher_host_state{name="server.example.com",rancherURL="http://1.1.1.1:8080/v2-beta",state="removing"} 0
rancher_host_state{name="server.example.com",rancherURL="http://1.1.1.1:8080/v2-beta",state="requested"} 0
rancher_host_state{name="server.example.com",rancherURL="http://1.1.1.1:8080/v2-beta",state="restoring"} 0
rancher_host_state{name="server.example.com",rancherURL="http://1.1.1.1:8080/v2-beta",state="updating_active"} 0
rancher_host_state{name="server.example.com",rancherURL="http://1.1.1.1:8080/v2-beta",state="updating_inactive"} 0
# HELP rancher_service_health_status HealthState of the service, as reported by the Rancher API. Either (1) or (0)
# TYPE rancher_service_health_status gauge
rancher_service_health_status{name="rocketchat",rancherURL="http://1.1.1.1:8080/v2-beta"} 0
# HELP rancher_service_scale scale of defined service as reported by Rancher
# TYPE rancher_service_scale gauge
rancher_service_scale{name="rocketchat",rancherURL="http://1.1.1.1:8080/v2-beta"} 1
# HELP rancher_service_state State of the service, as reported by the Rancher API
# TYPE rancher_service_state gauge
rancher_service_state{name="rocketchat",rancherURL="http://1.1.1.1:8080/v2-beta",state="activating"} 0
rancher_service_state{name="rocketchat",rancherURL="http://1.1.1.1:8080/v2-beta",state="active"} 0
rancher_service_state{name="rocketchat",rancherURL="http://1.1.1.1:8080/v2-beta",state="canceled_upgrade"} 0
rancher_service_state{name="rocketchat",rancherURL="http://1.1.1.1:8080/v2-beta",state="canceling_upgrade"} 0
rancher_service_state{name="rocketchat",rancherURL="http://1.1.1.1:8080/v2-beta",state="deactivasting"} 0
rancher_service_state{name="rocketchat",rancherURL="http://1.1.1.1:8080/v2-beta",state="finishing_upgrade"} 0
rancher_service_state{name="rocketchat",rancherURL="http://1.1.1.1:8080/v2-beta",state="inactive"} 1
rancher_service_state{name="rocketchat",rancherURL="http://1.1.1.1:8080/v2-beta",state="registering"} 0
rancher_service_state{name="rocketchat",rancherURL="http://1.1.1.1:8080/v2-beta",state="removed"} 0
rancher_service_state{name="rocketchat",rancherURL="http://1.1.1.1:8080/v2-beta",state="removing"} 0
rancher_service_state{name="rocketchat",rancherURL="http://1.1.1.1:8080/v2-beta",state="requested"} 0
rancher_service_state{name="rocketchat",rancherURL="http://1.1.1.1:8080/v2-beta",state="restarting"} 0
rancher_service_state{name="rocketchat",rancherURL="http://1.1.1.1:8080/v2-beta",state="rolling_back"} 0
rancher_service_state{name="rocketchat",rancherURL="http://1.1.1.1:8080/v2-beta",state="updating_active"} 0
rancher_service_state{name="rocketchat",rancherURL="http://1.1.1.1:8080/v2-beta",state="updating_inactive"} 0
rancher_service_state{name="rocketchat",rancherURL="http://1.1.1.1:8080/v2-beta",state="upgraded"} 0
rancher_service_state{name="rocketchat",rancherURL="http://1.1.1.1:8080/v2-beta",state="upgrading"} 0
# HELP rancher_stack_health_state HealthState of defined stack as reported by Rancher
# TYPE rancher_stack_health_state gauge
rancher_stack_health_state{name="rocket-chat",rancherURL="http://1.1.1.1:8080/v2-beta"} 0
# HELP rancher_stack_state State of defined stack as reported by Rancher
# TYPE rancher_stack_state gauge
rancher_stack_state{name="rocket-chat",rancherURL="http://1.1.1.1:8080/v2-beta",state="activating"} 0
rancher_stack_state{name="rocket-chat",rancherURL="http://1.1.1.1:8080/v2-beta",state="active"} 1
rancher_stack_state{name="rocket-chat",rancherURL="http://1.1.1.1:8080/v2-beta",state="canceled_upgrade"} 0
rancher_stack_state{name="rocket-chat",rancherURL="http://1.1.1.1:8080/v2-beta",state="canceling_upgrade"} 0
rancher_stack_state{name="rocket-chat",rancherURL="http://1.1.1.1:8080/v2-beta",state="error"} 0
rancher_stack_state{name="rocket-chat",rancherURL="http://1.1.1.1:8080/v2-beta",state="erroring"} 0
rancher_stack_state{name="rocket-chat",rancherURL="http://1.1.1.1:8080/v2-beta",state="finishing_upgrade"} 0
rancher_stack_state{name="rocket-chat",rancherURL="http://1.1.1.1:8080/v2-beta",state="removed"} 0
rancher_stack_state{name="rocket-chat",rancherURL="http://1.1.1.1:8080/v2-beta",state="removing"} 0
rancher_stack_state{name="rocket-chat",rancherURL="http://1.1.1.1:8080/v2-beta",state="requested"} 0
rancher_stack_state{name="rocket-chat",rancherURL="http://1.1.1.1:8080/v2-beta",state="rolling_back"} 0
rancher_stack_state{name="rocket-chat",rancherURL="http://1.1.1.1:8080/v2-beta",state="updating_active"} 0
rancher_stack_state{name="rocket-chat",rancherURL="http://1.1.1.1:8080/v2-beta",state="upgraded"} 0
rancher_stack_state{name="rocket-chat",rancherURL="http://1.1.1.1:8080/v2-beta",state="upgrading"} 0
```
