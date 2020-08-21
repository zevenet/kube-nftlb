# What does the test consist of?

This test is based on checking the dsr mode field

# Do you need any special parameters to configure?

To use DSR we have to configure the DSR mode by using annotations. In addition, the IP of the service itself must be configured in the loopback interface of the backends associated with the service (this is done automatically) and the backends have the same ports as the VIP

# What environment will be configured for the test?

We will configure a service with the DSR mode activated. Then we will create two backends and assign it to that service. Once there, we will check that the loopback type interface has been configured with the service IP.

Then we will apply a second configuration file (it is an update of it, where it goes from DSR to SNAT mode). We do this to verify that all the previous configuration that we have made on the loopback interfaces of the backends has been eliminated.

# How check the test status

To check that everything works we have to see if the interface of the backends associated with the service has the client's ip defined in its loopback interface and that after changing the type of service this configuration has disappeared