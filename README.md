# kube-nftlb

[![GoDoc](https://godoc.org/github.com/zevenet/kube-nftlb?status.svg)](https://godoc.org/github.com/zevenet/kube-nftlb)
[![Go Report Card](https://goreportcard.com/badge/github.com/zevenet/kube-nftlb)](https://goreportcard.com/report/github.com/zevenet/kube-nftlb)

##### Author: Víctor Manuel Oliver Acosta



## Description

`kube-nftlb` is a Kubernetes Pod made by two containers (`client` and `daemon`) able to communicate the Kubernetes API Server, using a Debian image with `nftlb` / `nftables` installed.

This project can request information from the API Server such as new, updated or deleted Services/Endpoints, and make rules in `nftables` accordingly.


## Software required before proceeding

* Docker
* Docker-machine
* Minikube [**v0.30.0**](https://github.com/kubernetes/minikube/releases/tag/v0.30.0) _(already started with_ `--kubernetes-version="v1.12.0"`_)_
* Golang
* `client-go`
* `nftables` and `nftlb` installed in the host or VM

... Or you can run `debian_tools_installer.sh` **as root** after a fresh Debian Testing install in a virtualized environment.


## Getting the cluster ready

**You must only do these steps if you have NOT done it before, and if you meet the specified conditions mentioned in each point.** Otherwise, you can skip this section.

* The first thing you have to do is clone the project on your machine:
```
root@pc: git clone https://github.com/zevenet/kube-nftlb
```
* This is a mandatory step if you started Minikube with `--vm-driver=none`, and you mustn't do it if that's not your case. `coredns` won't be able to resolve external hostnames unless you run this command:
```
root@pc: cd kube-nftlb
root@pc: kubectl apply -f yaml/give_internet_access_to_pods.yaml
```
* The cluster needs a `kube-nftlb` privileged rol, because in order to use `kube-nftlb` for communicating the API Server, it needs to be recognised and authenticated by the API Server. Run this command:
```
root@pc: kubectl apply -f yaml/authentication_system_level_from_pod.yaml
```


## Project test: steps to follow

1. The project will be available locally following the above steps. But first, `nftables` rules need to be monitorized in order to notice the changes that are being made. Run these commands and hide the terminal for later:
```
user@pc: su
root@pc: watch -n 1 nft list table nftlb
```

2. Open another terminal. To get inside the project directory, run these commands:
```
user@pc: su
root@pc: cd kube-nftlb
```

3. The script `build.sh` will compile `main.go` and will build a Docker container to put it inside the cluster. **Before running it, you MUST read the script. And be careful, all `nftables` rules you may have set could be flushed**. Once you have read it and adapted it to your use case, run:
```
root@pc: sh build.sh
```

4. Once the script has finished, the `kube-nftlb` Pod will be made as [DaemonSet](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/). Inside `yaml` there's a file ready for this, apply it to the cluster by running this:
```
root@pc: kubectl apply -f yaml/create_nftlb_as_daemonset.yaml
```
**Notice how rules are made in the first terminal you opened.**

5. The test will be made with a [Ghost](https://ghost.org/) instance, exposing, editing and deleting a Service. Run this command:
```
root@pc: kubectl create deployment ghost --image=ghost
```

6. The `ghost` Pod will be exposed through a Service with this command:
```
root@pc: kubectl expose deployment ghost --port=2368
```
**Notice how `ghost` rules are made in the first terminal you opened.**

7. Update the Service with this command, changing the port from 2368 to 2369, and save the file:
```
root@pc: kubectl edit service ghost
```
**Notice how `ghost` port has changed in the first terminal you opened.**

8. Delete the Service with this command:
```
root@pc: kubectl delete service ghost
```
**Notice how `ghost` rules are deleted in the first terminal you opened.**


## FAQ

* **I've done everything already, how can I stop watching `nftables` rules?**

Press `Control` + `C`.

* **I have followed the guide and I've got no errors. But, how can I delete the `kube-nftlb` Pod to test the project again from the start?**

Run this command as root:
```
root@pc: kubectl delete -f yaml/create_nftlb_as_daemonset.yaml
```

* **How can I also delete the `ghost` Pod? The guide explains how to delete its Service, but not its Pod.**

Run this command as root:
```
root@pc: kubectl delete deployment ghost
```
