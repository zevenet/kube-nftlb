# performance-tests/pool

Charts are made based on this data pool.

## Expected output

After running `generate_pool_data.sh`, this directory should look like this:

```
pool/
├── README.md
├── rules
│   ├── by-count-type
│   │   ├── kube-nftlb-replicas-test-010.txt
│   │   ├── kube-nftlb-replicas-test-050.txt
│   │   ├── kube-nftlb-replicas-test-100.txt
│   │   ├── kube-proxy-replicas-test-010.txt
│   │   ├── kube-proxy-replicas-test-050.txt
│   │   └── kube-proxy-replicas-test-100.txt
│   ├── by-endpoints-number
│   │   ├── kube-nftlb-create-service.txt
│   │   ├── kube-nftlb-delete-service.txt
│   │   ├── kube-proxy-create-service.txt
│   │   └── kube-proxy-delete-service.txt
│   └── columns
│       ├── kube-nftlb-replicas-test-010-create-service.txt
│       ├── kube-nftlb-replicas-test-010-delete-service.txt
│       ├── kube-nftlb-replicas-test-050-create-service.txt
│       ├── kube-nftlb-replicas-test-050-delete-service.txt
│       ├── kube-nftlb-replicas-test-100-create-service.txt
│       ├── kube-nftlb-replicas-test-100-delete-service.txt
│       ├── kube-proxy-replicas-test-010-create-service.txt
│       ├── kube-proxy-replicas-test-010-delete-service.txt
│       ├── kube-proxy-replicas-test-050-create-service.txt
│       ├── kube-proxy-replicas-test-050-delete-service.txt
│       ├── kube-proxy-replicas-test-100-create-service.txt
│       └── kube-proxy-replicas-test-100-delete-service.txt
└── time
    ├── by-count-type
    │   ├── kube-nftlb-replicas-test-010.txt
    │   ├── kube-nftlb-replicas-test-050.txt
    │   ├── kube-nftlb-replicas-test-100.txt
    │   ├── kube-proxy-replicas-test-010.txt
    │   ├── kube-proxy-replicas-test-050.txt
    │   └── kube-proxy-replicas-test-100.txt
    ├── by-endpoints-number
    │   ├── kube-nftlb-create-service.txt
    │   ├── kube-nftlb-delete-service.txt
    │   ├── kube-proxy-create-service.txt
    │   └── kube-proxy-delete-service.txt
    └── columns
        ├── kube-nftlb-replicas-test-010-create-service.txt
        ├── kube-nftlb-replicas-test-010-delete-service.txt
        ├── kube-nftlb-replicas-test-050-create-service.txt
        ├── kube-nftlb-replicas-test-050-delete-service.txt
        ├── kube-nftlb-replicas-test-100-create-service.txt
        ├── kube-nftlb-replicas-test-100-delete-service.txt
        ├── kube-proxy-replicas-test-010-create-service.txt
        ├── kube-proxy-replicas-test-010-delete-service.txt
        ├── kube-proxy-replicas-test-050-create-service.txt
        ├── kube-proxy-replicas-test-050-delete-service.txt
        ├── kube-proxy-replicas-test-100-create-service.txt
        └── kube-proxy-replicas-test-100-delete-service.txt

8 directories, 45 files
```
