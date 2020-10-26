# Description

This test is based on adding `addresses` to the farm if `spec.externalIPs` is not empty.

## Details

This test verifies that extra addresses have been configured for each IP defined within the externalIPs field.

If there are external IPs that route to one or more cluster nodes, Kubernetes Services can be exposed on those externalIPs. Traffic that ingresses into the cluster with the external IP (as destination IP), on the Service port, will be routed to one of the Service endpoints.
