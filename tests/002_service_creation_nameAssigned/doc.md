# What does the test consist of?

This test consists of creating a simple service that does not have a default name assigned or any backend associated with it.

# The assignment of the name of the service on which it is based?

The name of the service is defined by the name of the service and the name of the port. See [1],[2]

>
apiVersion: v1
kind: Service
metadata:
  name: my-service <[1]
  labels:
    app: front
spec:
  type: ClusterIP
  selector:
    app: front
  ports:
    - name: http <[2]
      protocol: TCP 
      port: 8080
      targetPort: 80

In the case in which the name does not exist in the port field, it is assigned a default one (default)

# How check the test status

In this test we have to verify that our service has been created and has been assigned a correct name (my-service--http)
