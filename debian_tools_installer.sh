#!/bin/sh

# This is a script that helps you set up a fresh Debian Testing system
# running in a virtualized environment. It installs everything you need
# before testing kube-nftlb. Don't forget to run this script as root.

# Recommended Debian Buster ISO:
#    https://cdimage.debian.org/debian-cd/current/amd64/iso-cd/
#    debian-10.4.0-amd64-netinst.iso

# 0. Change directory to /root/
cd


# 1. Update packages and upgrade them
apt-get update
apt-get upgrade -y


# 2. Install Docker (latest version)
apt-get install -y apt-transport-https ca-certificates curl gnupg2 software-properties-common
curl -fsSL https://download.docker.com/linux/debian/gpg | apt-key add -
add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/debian $(lsb_release -cs) stable"
apt-get update
apt-get install -y docker-ce docker-ce-cli containerd.io


# 3. Install Docker Machine (v0.16.2)
# Releases: https://github.com/docker/machine/releases/
curl -L https://github.com/docker/machine/releases/download/v0.16.2/docker-machine-`uname -s`-`uname -m` >/tmp/docker-machine &&
chmod +x /tmp/docker-machine &&
cp /tmp/docker-machine /usr/local/bin/docker-machine


# 4. Install kubectl (no hypervisor)
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | tee -a /etc/apt/sources.list.d/kubernetes.list
apt-get update
apt-get install -y kubectl


# 5. Install Minikube (latest version)
curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 && chmod +x minikube && cp minikube /usr/local/bin/ && rm minikube


# 6. Install Golang (v.1.14.2)
wget https://dl.google.com/go/go1.14.2.linux-amd64.tar.gz
tar xvfz go1.14.2.linux-amd64.tar.gz 
mv go /usr/local/go
cat << 'EOF' >> ~/.bashrc
export GOROOT=/usr/local/go
export GOPATH=$HOME/goProjects
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
EOF
source ~/.bashrc 
go version 

# 7. Install nftables
apt install -y nftables

# 8. Install conntrack
apt-get install conntrack

# 9. Start Minikube
  # if you are virtualizing, remember to leave 2 CPU
minikube start --vm-driver=none
