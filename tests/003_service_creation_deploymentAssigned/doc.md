# What does the test consist of?

This test consists of creating a service that have two backend assigned. 

# Are there any specific settings on the backends or something?

Each of the backends that we configure and are linked to the service through the labels field, in this case, app:front (see [1],[2]) will be automatically assigned to the service.

> 
	# service.yaml
	apiVersion: v1
	kind: Service
	metadata:
	  name: my-service
	  labels:
	    app: front <[1]
	...

	# deployment.yaml (aka backends)
	apiVersion: apps/v1
	kind: Deployment
	metadata:
	  name: lower-prio
	  labels:
	    app: front <[2]
	spec:
	  replicas: 2
	  selector:
	    matchLabels:
	      app: front <[2]
	  template:
		metadata:
          labels:
            app: front <[2]
    ...

The number of backends assigned to the service is based on the number of replicas defined in the yaml file of the deployment

# How check the test status

In this test we have to verify that our service has been created and has been assigned two backends.
