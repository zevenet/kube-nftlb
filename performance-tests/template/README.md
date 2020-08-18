# performance-tests/template

This is a directory with every basic file to test any deployment or daemonset.

## What's in this directory?

```
.
├── kubes
│   └── kube-test.yaml
├── README.md
└── testdata
    ├── deployments
    │   └── resource-test.yaml
    └── services
        └── resource-test.yaml

4 directories, 4 files
```

### **kubes/kube-test.yaml**

Daemonset + RBAC. Configuration related to the daemonset must be added to this file. The following block must be defined and added to the `kind: Daemonset` part in the file:

```yaml
metadata:
  name: kube-test
  namespace: kube-system
  labels:
    name: kube-test
    app: kube-test
spec:
  selector:
    matchLabels:
      name: kube-test
      app: kube-test
  template:
    metadata:
      labels:
        name: kube-test
        app: kube-test
```

Replace `kube-test` with the actual name.

### **testdata/deployments/resource-test.yaml**

Deployment. The following block must be defined and added to the file:

```yaml
metadata:
  name: resource-test
  labels:
    app: resource-test
spec:
  selector:
    matchLabels:
      app: resource-test
  template:
    metadata:
      labels:
        app: resource-test
```

Replace `resource-test` with the actual filename. Also, number of replicas must be specified at the end of the filename.

### **testdata/services/resource-test.yaml**

Service for that deployment. The following block must be defined and added to the file:

```yaml
metadata:
  name: resource-test
spec:
  selector:
    app: resource-test
```

Replace `resource-test` with the actual filename. Also, number of replicas must be specified at the end of the filename.
