# What does the test consist of?

This test is based on checking the priority on persistence fields. They are all possible situations to configure persistence in our service.

We can configure the persistence field in our farms through the annotations field of the yaml configuration file. 
The allowed parameters are listed in the nftlb api documentation: [**nftlb api documentation**](https://github.com/zevenet/nftlb)

# How can run the test

Just launch the script and the three examples are launched:

 - priority 1: Persistence is configured through annotations

 - priority 2: The session Affinity field is present but annotations have priority. The session time if the "sessionAfinityConfig" field is not present is the default.

 - priority 3: There is only the session affinity field, it collects the values ​​of those fields.

All the tests are in the same file, one below the other.
