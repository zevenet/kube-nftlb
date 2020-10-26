# Description

This test is based on checking the priority on the `persistence` field.

## Details

The allowed parameters are listed in the [**nftlb api documentation**](https://github.com/zevenet/nftlb).

Annotations should take preference against the `sessionAffinity` field. `persistence` field must be `srcport` as defined in the annotation.
