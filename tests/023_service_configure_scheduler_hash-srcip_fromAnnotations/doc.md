# What does the test consist of?

This test is based on configuring the "scheduler" field through the use of annotations.

The "scheduler" field supports the following values:

- rr
- symhash
- hash


If the field is of type "hash" it has a series of extra combinations. This is because there are no specific annotations to read the "sched-param" field. 
For this reason, the following combinations are enabled to complement the hash-type sheduler.

- hash-srcip
- hash-dstip
- hash-srcport
- hash-dstport 
- hash-srcmac 
- hash-dstmac

Annotations example to configure the scheduler field:

- service.kubernetes.io/kube-nftlb-load-balancer-scheduler: "hash-srcip"

For more information consult [**nftlb api documentation**](https://github.com/zevenet/nftlb)
