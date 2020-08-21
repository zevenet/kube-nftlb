# What does the test consist of?

This test consists of creating a multiples services that have each one two backend assigned. 

# What is the creation of various services based on?

It is based on the fact that we can define several services from the same yaml configuration file.
For this we take into account the name of the port field. We create a service for each name of the port field that you have defined. See [1], [2]

> 
	# service.yaml
	apiVersion: v1
	kind: Service
	metadata:
	  name: my-service
	  labels:
	    app: front
	spec:
	  type: ClusterIP
	  selector:
	    app: front
	  ports:
        - name: http <[1]
          protocol: TCP
          port: 8080
          targetPort: 80
        - name: https <[2]
          protocol: TCP
          port: 8181
          targetPort: 81

If we remember what the naming is like and taking this test as an example, we will create two services, one called "my-service--http" and the other called "my-service--https". 
The name of the service followed after the port name as an identifier.

# And when creating multiple services, what about the backends?

Each service has the same backends with the same name and ip, the only thing that varies between them is its port that corresponds to the targetPort of the service.

# How check the test status

In this test we have to verify that our services has been created and has been assigned two backends each one.
