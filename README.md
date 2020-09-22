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

```console
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

## Local Configuration ⚙

We have to remove the chains that kubernetes configures by default. To achieve this we have to stop the kubelet service, add a variable to the configuration file and reactivate the service. Follow the following commands:

```console
systemctl stop kubelet.service
echo "makeIPTablesUtilChains: false" >> /var/lib/kubelet/config.yaml
nft delete table ip nat
nft delete table ip filter
nft delete table ip mangle
systemctl start kubelet.service
```
If everything has gone well the kubelet service will not create those tables again. Now you will have to apply a rule to recover the connection with your deployments:

```console
iptables -N POSTROUTING -t filter
iptables -A POSTROUTING -t nat -s 172.17.0.0/16 ! -o docker0 -j MASQUERADE
```

## Creation of a simple service ✏

In this section we are going to see the different settings that we can apply to create our service. The first thing we have to know is that it is a service and how we can create a simple one and check that it has been created correctly.

A Service is an abstraction which defines a logical set of Pods and a policy by which to access them. A Service in Kubernetes is a REST object, similar to a Pod. Like all the REST objects, you can POST a Service definition to the API server to create a new instance. The name of a Service object must be a valid DNS label name. For example:

```yaml
# service.yaml configuration creates a service
apiVersion: v1
kind: Service
metadata:
  name: my-service
  labels:
    app: front
spec:
  type: ClusterIP
  selector:
    app: front
  ports:
    - name: http
      protocol: TCP
      port: 8080
      targetPort: 80
```

This specification creates a new Service object named “my-service”, which targets TCP port 8080 on any Pod with the app=front label.

To apply this configuration and verify that our service has been created we have to use the following commands:

- Apply the configuration contained within the yaml file.

```console
kubectl apply -f service.yaml
```

- Shows all services, a service called "my-service" should appear.

```console
kubectl get services -A
```

```
NAMESPACE     NAME         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                  AGE
default       my-service   ClusterIP   IP_cluster      <none>        8080/TCP                 6m12s
```

Now we are going to check that after the creation of our service our farm has been correctly configured. To do that we need the nftlb key generated during compilation to be able to make requests to the nftlb api. The key in the '.env' file. You can copy it from there or launch the following commands from the kube-nftlb directory:

```console
NFTLB_KEY=$(grep 'NFTLB_KEY' .env | sed 's/NFTLB_KEY=//'); curl -H "Key: $NFTLB_KEY" http://localhost:5555/farms/my-service--http
```

```json
{
        "farms": [
                {
                        "name": "my-service--http",
                        "family": "ipv4",
                        "virtual-addr": "IP",
                        "virtual-ports": "8080",
                        "source-addr": "",
                        "mode": "snat",
                        "protocol": "tcp",
                        "scheduler": "rr",
                        "sched-param": "none",
                        "persistence": "none",
                        "persist-ttl": "60",
                        "helper": "none",
                        "log": "none",
                        "mark": "0x0",
                        "priority": "1",
                        "state": "up",
                        "new-rtlimit": "0",
                        "new-rtlimit-burst": "0",
                        "rst-rtlimit": "0",
                        "rst-rtlimit-burst": "0",
                        "est-connlimit": "0",
                        "tcp-strict": "off",
                        "queue": "-1",
                        "intra-connect": "on",
                        "backends": [],
                        "policies": []
                }
        ]
}
```

*The curl that we have launched returns a json with the information configured in our farms.*

## Creation and assignment of deployments ✏

In this section we will see how to create a deployment and how we can assign it to other pods (our service). But first we have to know what a deployment is.

Deployments represent a set of multiple, identical Pods with no unique identities. A Deployment runs multiple replicas of your application and automatically replaces any instances that fail or become unresponsive. In this way, Deployments help ensure that one or more instances of your application are available to serve user requests. Deployments are managed by the Kubernetes Deployment controller.

Deployments use a Pod template, which contains a specification for its Pods. The Pod specification determines how each pod should look like: what applications should run inside its containers, which volumes the Pods should mount, its labels, and more. Let's see an example:

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: lower-prio
  labels:
    app: front
spec:
  replicas: 2
  selector:
    matchLabels:
      app: front
  template:
    metadata:
      labels:
        app: front
    spec:
      containers:
      - name: nginx
        image: nginx:alpine
```
Through the "matchLabels" field we can find the pod of our service. We are going to apply our deployment and check that it has been created correctly.

```console
kubectl -f apply deployment.yaml
kubectl get pods --all-namespaces
```

```
NAMESPACE     NAME                             READY   STATUS    RESTARTS   AGE
default       lower-prio-64588d8b49-jjlvm      1/1     Running   0          12s
default       lower-prio-64588d8b49-lvk92      1/1     Running   0          12s
```

Now we are going to check that after creating our deployment, our farm has the backends configured correctly. We will have as many backends configured as replicas we have specified.

```console
NFTLB_KEY=$(grep 'NFTLB_KEY' .env | sed 's/NFTLB_KEY=//')
curl -H "Key: $NFTLB_KEY" http://localhost:5555/farms/my-service--http
```

```json
{
        "farms": [
                {
                        "name": "my-service--http",
                        "family": "ipv4",
                        "virtual-addr": "IP",
                        "virtual-ports": "8080",
                        "source-addr": "",
                        "mode": "snat",
                        "protocol": "tcp",
                        "scheduler": "rr",
                        "sched-param": "none",
                        "persistence": "none",
                        "persist-ttl": "60",
                        "helper": "none",
                        "log": "none",
                        "mark": "0x0",
                        "priority": "1",
                        "state": "up",
                        "new-rtlimit": "0",
                        "new-rtlimit-burst": "0",
                        "rst-rtlimit": "0",
                        "rst-rtlimit-burst": "0",
                        "est-connlimit": "0",
                        "tcp-strict": "off",
                        "queue": "-1",
                        "intra-connect": "on",
                        "backends": [
                                {
                                        "name": "lower-prio-64588d8b49-lvk92",
                                        "ip-addr": "IP",
                                        "port": "80",
                                        "weight": "1",
                                        "priority": "1",
                                        "mark": "0x0",
                                        "est-connlimit": "0",
                                        "state": "up"
                                },
                                {
                                        "name": "lower-prio-64588d8b49-jjlvm",
                                        "ip-addr": "IP",
                                        "port": "80",
                                        "weight": "1",
                                        "priority": "1",
                                        "mark": "0x0",
                                        "est-connlimit": "0",
                                        "state": "up"
                                }
                        ],

                        "policies": []
                }
        ]
}
```

## Setting up our service 📌

We can configure our service with different settings. In general, to configure our service we will use annotations, a field used in our configuration file yaml. In a few words, annotations are a field that will allow us to enter data outside kubernetes.

Through this field we can configure our service with different values that nftlb supports. For example, we can configure the mode of our service, if our backends have persistence or change our load balancing scheduling. We are going to see all the configuration that we can add using annotations, and then we are going to see a small example of the syntax of our annotations.

### Configure Mode

We can configure how the load balancer layer 4 core is going to operate. The options are:
- **snat** the backend responds to the load balancer in order to send the response to the client.
- **dnat** the backend will respond directly to the client, load balancer has to be configured as gateway in the backend server.
- **dsr (Direct Server Return)** the client connects to the VIP, then the load balancer changes its destination MAC address for the backend MAC address.
- **stlsdnat (Stateless DNAT)** the load balancer switch destination address for the backend address and forward it to the backend as DNAT does, but it doesn’t manage any kind of connection information.

```
service.kubernetes.io/kube-nftlb-load-balancer-mode: "snat"
service.kubernetes.io/kube-nftlb-load-balancer-mode: "dnat"
service.kubernetes.io/kube-nftlb-load-balancer-mode: "dsr"
service.kubernetes.io/kube-nftlb-load-balancer-mode: "stlsdnat"
```

### Configure Persistence

We can configure the type of persistence that is used on the configured farm. This can be configured in two ways. Via annotations and with the sessionAffinity field.

Through annotations:

- **srcip** Source IP, will assign the same backend for every incoming connection depending on the source IP address only.
- **srcport** Source Port, will assign the same backend for every incoming connection depending on the source port only.
- **srcmac** Source MAC, With this option, the farm will assign the same backend for every incoming connection depending on link-layer MAC address of the packet.
- **srcip-srcport** Source IP and Source Port, will assign the same backend for every incoming connection depending on both, source IP and source port.
- **srcip-dstport** Source IP and Destination Port, will assign the same backend for every incoming connection depending on both, source IP and destination port.

```
service.kubernetes.io/kube-nftlb-load-balancer-persistence: "srcip"
service.kubernetes.io/kube-nftlb-load-balancer-persistence: "srcport"
service.kubernetes.io/kube-nftlb-load-balancer-persistence: "srcip-srcport"
service.kubernetes.io/kube-nftlb-load-balancer-persistence: "srcip-dstport"
```

Through sessionAffinity field:

- **ClientIP** The sessionAffinity ClientIP is equivalent to the "srcip" field in annotations.
- **sessionAffinityConfig** Through the "timeoutSeconds" field we can configure the stickiness timeout in seconds.

```yaml
spec:
  type: ClusterIP
  selector:
    app: front
  sessionAffinity: ClientIP
  sessionAffinityConfig:
    clientIP:
      timeoutSeconds: 10
```

By default, settings made with annotations have priority, even if the sessionAffinity field is defined. The "stickiness timeout in seconds" cannot be configured via annotations. The default value is chosen unless there is a sessionAffinity field and a sessionAffinityConfig where it will collect the value of that field.

### Configure Scheduler

We can configure the type of load balancing scheduling used to dispatch the traffic between the backends. The options are:

- **rr** does a sequential select between the backend pool, each backend will receive the same number of requests.
- **symhash** balance packets that match the same source IP and port and destination IP and port, so it could hash a connection in both ways (during inbound and outbound).
- **hash-srcip** balances packets that match the same source IP to the same backend.
- **hash-srcip-srcport** balances packets that match the same source IP and port to the same backend.

```
service.kubernetes.io/kube-nftlb-load-balancer-scheduler: "rr"
service.kubernetes.io/kube-nftlb-load-balancer-persistence: "symhash"
service.kubernetes.io/kube-nftlb-load-balancer-persistence: "hash-srcip"
service.kubernetes.io/kube-nftlb-load-balancer-persistence: "hash-srcip-srcport"
```

### Configure Helper

We can configure the helper of the layer 4 protocol to be balanced to be used. The options are:

- **none** it's the default option.
- **amanda** enabling this option, the farm will be listening for incoming UDP packets to the current virtual IP and then will parse AMANDA headers for each packet in order to be correctly distributed to the backends.
- **ftp** enabling this option, the farm will be listening for incoming TCP connections to the current virtual IP and port 21 by default, and then will parse FTP headers for each packet in order to be correctly distributed to the backends.
- **h323** enabling this option, the farm will be listening for incoming TCP and UDP packets to the current virtual IP and port.
- **irc** enabling this option, the farm will be listening for incoming TCP connections to the current virtual IP and port and then will parse IRC headers for each packet in order to be correctly distributed to the backends.
- **netbios-ns** enabling this option, the farm will be listening for incoming UDP packets to the current virtual IP and port and then will parse NETBIOS-NS headers for each packet in order to be correctly distributed to the backends.
- **pptp** enabling this option, the farm will be listening for incoming TCP connections to the current virtual IP and port and then will parse the PPTP headers for each packet in order to be correctly distributed to the backends.
- **sane** enabling this option, the farm will be listening for incoming TCP connections to the current virtual IP and port and then will parse the SANE headers for each packet in order to be correctly distributed to the backends.
- **sip** enabling this option, the farm will be listening for incoming UDP packets to the current virtual IP and port 5060 by default, and then will parse SIP headers for each packet in order to be correctly distributed to the backends.
- **snmp** enabling this option, the farm will be listening for incoming UDP packets to the current virtual IP and port and then will parse the SNMP headers for each packet in order to be correctly distributed to the backends.
- **tftp** enabling this option, the farm will be listening for incoming UDP packets to the current virtual IP and port 69 by default, and then will parse TFTP headers for each packet in order to be correctly distributed to the backends.

```
service.kubernetes.io/kube-nftlb-load-balancer-helper: "none"
service.kubernetes.io/kube-nftlb-load-balancer-helper: "amanda"
service.kubernetes.io/kube-nftlb-load-balancer-helper: "ftp"
service.kubernetes.io/kube-nftlb-load-balancer-helper: "h323"
service.kubernetes.io/kube-nftlb-load-balancer-helper: "irc"
service.kubernetes.io/kube-nftlb-load-balancer-helper: "netbios-ns"
service.kubernetes.io/kube-nftlb-load-balancer-helper: "pptp"
service.kubernetes.io/kube-nftlb-load-balancer-helper: "sane"
service.kubernetes.io/kube-nftlb-load-balancer-helper: "sip"
service.kubernetes.io/kube-nftlb-load-balancer-helper: "snmp"
service.kubernetes.io/kube-nftlb-load-balancer-helper: "tftp"
```

### Configure Logs

We can define in which netfilter flow in which stage you are going to print logs. The options are:

- **none** it's the default option.
- **output** log for traffic going from the host to the pods.
- **forward** for traffic that passes through the host. It can be between two pods or from outside to a pod.

```
service.kubernetes.io/kube-nftlb-load-balancer-log: "none"
service.kubernetes.io/kube-nftlb-load-balancer-log: "output"
service.kubernetes.io/kube-nftlb-load-balancer-log: "forward"
```

### How to set up annotations

It is very simple, all we have to do is add it to the configuration file of our service. Let's see an example:

```yaml
# annotations-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: my-service
  labels:
    app: front
  annotations:
    service.kubernetes.io/kube-nftlb-load-balancer-mode: "snat"
    service.kubernetes.io/kube-nftlb-load-balancer-scheduler: "hash-srcip"
spec:
  type: ClusterIP
  selector:
    app: front
  ports:
    - name: http
      protocol: TCP
      port: 8080
      targetPort: 80
```

## Benchmarks 📊

This data can be found at `resources/` directory.

### Environment

- **Host**: Bare metal, not virtualized
- **CPU**: Intel i5 660 (4) @ 3.334GHz
- **OS**: Debian GNU/Linux 10.5 (buster) x86_64 (netinst, no DE)
- **Kernel**: Linux debian 5.7.0-0.bpo.2-amd64 #1 SMP Debian 5.7.10-1~bpo10+1 (2020-07-30) x86_64 GNU/Linux
- **Memory**: 3740MiB
- **How many times was the test repeated**: 30
- **Software**:
  - Minikube v1.12.3 (flags: `--vm-driver=none --extra-config=kubeadm.skip-phases=addon/kube-proxy`)
  - Docker v19.03.12
  - iptables v1.8.3
  - nftables v0.9.4
  - kubectl client v1.19.0, server v1.18.3
  - kube-proxy v1.18.3 (iptables mode)
  - kube-nftlb commit 4f89518e28efae24b88b5b42e8a0a4698bb17fa6

_If we're missing something, open an issue and let us know._

### Charts: time by endpoints number

[![Chart time/by-endpoints-number/create-service](./resources/time/by-endpoints-number/create-service.png)](./resources/time/by-endpoints-number/create-service.png)

[![Chart time/by-endpoints-number/delete-service](./resources/time/by-endpoints-number/delete-service.png)](./resources/time/by-endpoints-number/delete-service.png)

### Charts: rules by count type

[![Chart rules/by-count-type/replicas-test-010](./resources/rules/by-count-type/replicas-test-010.png)](./resources/rules/by-count-type/replicas-test-010.png)

[![Chart rules/by-count-type/replicas-test-050](./resources/rules/by-count-type/replicas-test-050.png)](./resources/rules/by-count-type/replicas-test-050.png)

[![Chart rules/by-count-type/replicas-test-100](./resources/rules/by-count-type/replicas-test-100.png)](./resources/rules/by-count-type/replicas-test-100.png)
