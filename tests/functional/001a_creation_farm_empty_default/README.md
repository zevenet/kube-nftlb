# Description

This test consists of creating a Service (that doesn't expose anything) and making a farm for nftlb based on that Service. The farm must be assigned a default name with some default values.

## Details

A farm name is defined by its Service name and its related ServicePort name. If the ServicePort has an empty name, it's replaced by a default one. Also, some other values have defaults, so they must be set.
