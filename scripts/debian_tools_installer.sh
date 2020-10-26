#!/usr/bin/env bash

# This is a script that helps you set up a fresh Debian Buster system
# running in a virtualized environment. It installs everything you need
# before deploying kube-nftlb. Don't forget to run this script as root.

# Recommended Debian Buster ISO:
#    https://cdimage.debian.org/debian-cd/current/amd64/iso-cd/
#    debian-10.x.y-amd64-netinst.iso

# Update packages, upgrade them and install essential tools
apt-get update
apt-get upgrade -y
apt-get install -y apt-transport-https ca-certificates curl gnupg2 software-properties-common wget

# Install Docker (latest version)
curl -fsSL https://download.docker.com/linux/debian/gpg | apt-key add -
add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/debian $(lsb_release -cs) stable"
apt-get update
apt-get install -y docker-ce docker-ce-cli containerd.io

# Install kubectl (latest version, no hypervisor)
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | tee -a /etc/apt/sources.list.d/kubernetes.list
apt-get update
apt-get install -y kubectl

# Install Minikube (latest version)
curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
chmod +x minikube
mv minikube /usr/local/bin/

# Install nftables and several necessary libraries (before that, add the zevenet repository with its corresponding key)
echo "deb [arch=amd64] http://repo.zevenet.com/ce/v5 buster main" | tee -a /etc/apt/sources.list
wget -O - http://repo.zevenet.com/zevenet.com.gpg.key | apt-key add -
apt-get update
apt-get install -y libnftnl11 nftables conntrack
