# What does the test consist of?

In this test can define in which netfilter flow in which stage you are going to print logs. The options are:

- output: log for traffic going from the host to the pods
- forward: for traffic that passes through the host. It can be between two pods or from outside to a pod.

# How check the test status?

This test allows us to find out if the "logs" field of our farm has been configured. If it is configured with the value that we have passed it, it means that it has been configured correctly.