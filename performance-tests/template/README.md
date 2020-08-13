# How-to example

This is a template with every basic file to test any deployment or daemonset.

## What's in this directory?

```
.
├── expected-rule-count.sh
├── kubes
│   └── kube-test.yaml
├── README.md
└── testdata
    ├── deployments
    │   └── resource-test.yaml
    ├── expected-rule-count
    │   └── kube-test
    │       └── resource-test.txt
    └── services
        └── resource-test.yaml

6 directories, 6 files
```

### **expected-rule-count.sh**

Run this file with your 'kube-test.yaml' file to test it against the actual resources located at `testdata/deployments` and `testdata/services`.

```console
# Copy the file to the project root
root@debian:kubernetes-rules-test# cp example/expected-rule-count.sh .

# Give it execute permissions
root@debian:kubernetes-rules-test# chmod +x expected-rule-count.sh

# Run it passing your kube-test file as the first parameter (this is an example)
root@debian:kubernetes-rules-test# ./expected-rule-count.sh ./kubes/kube-test.yaml
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

Replace `kube-test` with the actual name. Also, the filename must be same the name.

### **testdata/deployments/resource-test.yaml**

Deployment file. The following block must be defined and added to the file:

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

Replace `resource-test` with the actual filename.

### **testdata/services/resource-test.yaml**

Service file. The following block must be defined and added to the file:

```yaml
metadata:
  name: resource-test
spec:
  selector:
    app: resource-test
```

Replace `resource-test` with the actual filename.

### **testdata/expected-rule-count/kube-test/resource-test.txt**

Expected rule count to be applied every service is applied. It's needed to get how much time it takes to process services. The following block must be defined and added to the file:

```
create-service: 0
delete-service: 0
delete-deployment: 0
```

Replace the numbers with the actual rule count for every line. These results from `expected-rule-count.sh` are the most important ones, the rest (`create-kube`, `create-deployment`) will be ignored.
