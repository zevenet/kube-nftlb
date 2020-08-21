# What does the test consist of?

This test consists of creating a service that have nodeport mode enabled.

# What the nodeport mode is based on?

A NodePort service is the most primitive way to get external traffic directly to your service. NodePort, as the name implies, opens a specific port on all the Nodes and any traffic that is sent to this port is forwarded to the service.

That means, from your external machine you will be able to access the service through the nodeport port

# Are there any special settings?

We just have to configure the type of our service to "nodePort" mode and add a port to it within the port range of the nodeport field (30000-32767). See [1],[2]

> 
	# service.yaml
	# Yaml Service
	apiVersion: v1
	kind: Service
	metadata:
	  name: my-service
	  labels:
	    app: front
	spec:
	  type: NodePort <[1]
	  selector:
	    app: front
	  ports:
        - name: http
          protocol: TCP
          port: 8080
          targetPort: 80
          nodePort: 32490 <[2]

When we finish configuring it, it will create our service and an additional one called "my-service--nodePort". 
The difference between both is that the nodeport service does not have virtual-addr and its port is the one defined in the nodePort field

# How check the test status

It is enough to make a request from outside our local environment to the ip of our machine (where's kubernetes environment is installed) followed by the nodeport port.

If there is connectivity, a message from the server http should appear. If this is not the case it means something went wrong.

PD: In order to run the test you need to configure some parameters in the script. Some of them are the IP of the external client and an ssh key to establish a connection with the client. More information in the script