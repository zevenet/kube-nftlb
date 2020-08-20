# kube-nftlb

[![GoDev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go)](https://pkg.go.dev/github.com/zevenet/kube-nftlb?tab=overview)
[![Go report card](https://goreportcard.com/badge/github.com/zevenet/kube-nftlb)](https://goreportcard.com/report/github.com/zevenet/kube-nftlb)
![License](https://img.shields.io/github/license/zevenet/kube-nftlb)

`kube-nftlb` is a Kubernetes Daemonset able to communicate the Kubernetes API Server, based on a Debian Buster image with [`nftlb`](https://github.com/zevenet/nftlb) installed.

It can request information from the API Server such as new, updated or deleted Services/Endpoints, and make rules in `nftables` accordingly.

## Prerequisites 📋

* Docker
* Minikube
* `kubectl`
* `nftables`
* `libnftnl11`
* `conntrack`

Also, you can run `debian_tools_installer.sh` **as root** after a fresh Debian Buster install.

```console
root@debian:kube-nftlb# ./debian_tools_installer.sh
```

## Installation 🔧

```
# Clone the project
user@debian:~# git clone https://github.com/zevenet/kube-nftlb

# Change directory
user@debian:~# cd kube-nftlb

# Copy and rename .env.example to .env
user@debian:kube-nftlb# cp .env.example .env

# Generate a random password for nftlb
user@debian:kube-nftlb# NFTLB_KEY=$(base64 -w 32 /dev/urandom | tr -d /+ | head -n 1) ; sed -i "s/^NFTLB_KEY=.*$/NFTLB_KEY=$NFTLB_KEY/" .env

# Change user to root
user@debian:kube-nftlb# su

# Modify scripts permissions to grant root execute access
root@debian:kube-nftlb# chmod +x *.sh

# Build the Docker image with build.sh (prerequisites must be met before this)
root@debian:kube-nftlb# ./build.sh
```

## Deployment 🚀

1. Start Minikube without `kube-proxy` being deployed by default:
```console
root@debian:kube-nftlb# minikube start --vm-driver=none --extra-config=kubeadm.skip-phases=addon/kube-proxy
```

2. The cluster needs to apply some settings, and they are inside `yaml/`. `coredns` will be able to resolve external hostnames and `kube-nftlb` will be deployed after running this command:

```console
root@debian:kube-nftlb# kubectl apply -f yaml
```
