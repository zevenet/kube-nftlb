# Description

This test consists of creating a farm based on a Service that has two backends assigned.

## Details

Each backend is linked to the Service through the labels field, so they will be automatically assigned to the Service. The number of backends assigned to the Service is based on the number of replicas defined in the deployment.
