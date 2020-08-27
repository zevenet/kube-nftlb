# What does the test consist of?

This test is based on configuring the "externalIPs" field.

If there are external IPs that route to one or more cluster nodes, Kubernetes Services can be exposed on those externalIPs. Traffic that ingresses into the cluster with the external IP (as destination IP), on the Service port, will be routed to one of the Service endpoints. externalIPs are not managed by Kubernetes and are the responsibility of the cluster administrator.

We can configure this field within the "Specs" section:

	spec:
	    type: ClusterIP
	    externalIPs:
              - 192.168.10.89
              - 192.168.10.90
              - 192.168.10.92

As mentioned before, the administration of these IPs depends on the administrator. If you want to have connectivity from outside your local environment, you have to manually configure your network interface with those IPs.
Locally, you do not need to configure anything, you can try to make a request to the backends without problems.

# How check the test status?

This test verifies that an extra service has been configured for each IP defined within the externalIPs field. This service will have the same backends as the original service.

If everything has gone well we should see them in the logs.
