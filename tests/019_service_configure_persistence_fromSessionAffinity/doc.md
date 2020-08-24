# What does the test consist of?

This test is based on checking the persistence field. We can configure the persistence field in our farms through the annotations field of the yaml configuration file.

The allowed parameters are listed in the nftlb api documentation: [**nftlb api documentation**](https://github.com/zevenet/nftlb)

We can also configure persistence using the sessionAffinity field and sessionAffinityConfig from our service.yaml configuration file

>
	sessionAffinity: ClientIP
	sessionAffinityConfig:
    	clientIP:
      		timeoutSeconds: 10

(more info about sessionAffinity field in kubernetes doc)

# What parameters does the test take into account

By default, what we do first is check for annotations related to persistence. If found, these annotations always have priority over the "sessionAffinity" field, even if it is configured. Persistence is configured with a mode and a session time, this session time is not defined in the annotations. This field is collected through the "sessionAffinityConfig" field or if it is not configured it collects the default values.

If it does not find annotations related to persistence, collect the values ​​of the "sessionAffinity" and sessionAffinityConfig field. If there are none of the 2, nothing is configured.

# How check the test status

In our backends we have an nginx server configured. The test is based on launching several requests to the service, in a normal situation the requests are distributed among the backends. 
But when we activate persistence, only one of them responds, for this we use the marks fields.