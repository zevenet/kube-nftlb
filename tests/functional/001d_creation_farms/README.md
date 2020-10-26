# Description

This test consists of creating multiple farms from a Service that have two backends assigned for each farm.

## Details

Several farms can be generated from a single Service. A ServicePort corresponds to a farm. Each farm has the same backends with the same name and IP, the only thing that varies between those backends is their port.
