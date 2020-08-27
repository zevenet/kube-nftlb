# What does the test consist of?

We can configure the helper of the layer 4 protocol to be balanced to be used. The options are:

We can configure the helper of the layer 4 protocol to be balanced to be used. The options are:
- **none** it's the default option
- **amanda** enabling this option, the farm will be listening for incoming UDP packets to the current virtual IP and then will parse AMANDA headers for each packet in order to be correctly distributed to the backends
- **ftp** enabling this option, the farm will be listening for incoming TCP connections to the current virtual IP and port 21 by default, and then will parse FTP headers for each packet in order to be correctly distributed to the backends.
- **h323** enabling this option, the farm will be listening for incoming TCP and UDP packets to the current virtual IP and port
- **irc** enabling this option, the farm will be listening for incoming TCP connections to the current virtual IP and port and then will parse IRC headers for each packet in order to be correctly distributed to the backends
- **netbios-ns** enabling this option, the farm will be listening for incoming UDP packets to the current virtual IP and port and then will parse NETBIOS-NS headers for each packet in order to be correctly distributed to the backends
- **pptp** enabling this option, the farm will be listening for incoming TCP connections to the current virtual IP and port and then will parse the PPTP headers for each packet in order to be correctly distributed to the backends
- **sane** enabling this option, the farm will be listening for incoming TCP connections to the current virtual IP and port and then will parse the SANE headers for each packet in order to be correctly distributed to the backends.
- **sip** enabling this option, the farm will be listening for incoming UDP packets to the current virtual IP and port 5060 by default, and then will parse SIP headers for each packet in order to be correctly distributed to the backends
- **snmp** enabling this option, the farm will be listening for incoming UDP packets to the current virtual IP and port and then will parse the SNMP headers for each packet in order to be correctly distributed to the backends
- **tftp** enabling this option, the farm will be listening for incoming UDP packets to the current virtual IP and port 69 by default, and then will parse TFTP headers for each packet in order to be correctly distributed to the backends

```code
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

# How check the test status?

This test allows us to find out if the "helper" field of our farm has been configured. If it is configured with the value that we have passed it, it means that it has been configured correctly.