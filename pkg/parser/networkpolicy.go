package parser

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/zevenet/kube-nftlb/pkg/types"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NetworkPolicyNamesAsPolicies
func NetworkPolicyNamesAsPolicies(nwPolicy *networkingv1.NetworkPolicy) *types.Policies {
	policies := make([]types.Policy, 0)

	for _, policyType := range []string{"input", "output"} {
		for _, listType := range []string{"whitelist", "blacklist"} {
			policies = append(policies, types.Policy{
				Name: formatPolicyName(nwPolicy.Name, policyType, listType),
			})
		}
	}

	return &types.Policies{
		Policies: policies,
	}
}

// NetworkPolicyAsFarms
func NetworkPolicyAsFarms(nwPolicy *networkingv1.NetworkPolicy) (*types.Farms, error) {
	// Get pod IP list
	//podIPList := getPodIPListFromPodSelector(nwPolicy.Namespace, &nwPolicy.Spec.PodSelector)

	return nil, nil
}

// NetworkPolicyAsPolicies
func NetworkPolicyAsPolicies(nwPolicy *networkingv1.NetworkPolicy) (*types.Policies, error) {
	// Error checking
	if len(nwPolicy.Spec.PolicyTypes) == 0 {
		return nil, fmt.Errorf("%s: Spec.PolicyTypes is empty", nwPolicy.Name)
	}

	// Make wait group, ingress and egress can be processed in parallel
	var wg sync.WaitGroup
	wg.Add(2)

	// Set default values
	policies := make([]types.Policy, 0)
	ingressTypePolicyExists, egressTypePolicyExists := containsPolicyTypes(nwPolicy)

	go func() {
		// Format ingress policies
		if ingressTypePolicyExists && nwPolicy.Spec.Ingress != nil {
			whitelistPolicy, blacklistPolicy := getIngressPolicies(nwPolicy)
			policies = append(policies, whitelistPolicy, blacklistPolicy)
		}
		wg.Done()
	}()

	go func() {
		// Format egress policies
		if egressTypePolicyExists && nwPolicy.Spec.Egress != nil {
			whitelistPolicy, blacklistPolicy := getEgressPolicies(nwPolicy)
			policies = append(policies, whitelistPolicy, blacklistPolicy)
		}
		wg.Done()
	}()

	// Release wait lock when both goroutines are done
	wg.Wait()

	return &types.Policies{
		Policies: policies,
	}, nil
}

func containsPolicyTypes(nwPolicy *networkingv1.NetworkPolicy) (ingress bool, egress bool) {
	for _, value := range nwPolicy.Spec.PolicyTypes {
		if value == networkingv1.PolicyTypeIngress {
			ingress = true
		} else if value == networkingv1.PolicyTypeEgress {
			egress = true
		}
	}
	return
}

func newPolicy(list []string, name string, policyType string, listType string) types.Policy {
	return types.Policy{
		Name:     formatPolicyName(name, policyType, listType),
		Type:     listType,
		Elements: parseList(list),
	}
}

func getLists(policyPeer networkingv1.NetworkPolicyPeer) ([]string, []string) {
	// Make ingress whitelist and blacklist
	whitelist := make([]string, 0)
	blacklist := make([]string, 0)

	// IPBlock must be checked, it's optional
	if policyPeer.IPBlock != nil {
		// Add IPBlock CIDR to whitelist (IPBlock must have a CIDR, we don't check if it's empty)
		whitelist = append(whitelist, policyPeer.IPBlock.CIDR)

		// Add Except IPs to blacklist (Except must be checked, it's optional)
		if policyPeer.IPBlock.Except != nil {
			blacklist = append(blacklist, policyPeer.IPBlock.Except...)
		}
	}

	// labelSelector must be checked, it's optional
	if policyPeer.NamespaceSelector != nil {
		// Add pod IPs selected by PodSelector (can be nil) that are inside namespaces selected by
		// NamespaceSelector (can't be nil) to whitelist; this works as a logical AND
		whitelist = append(whitelist, getPodIPListFromSelectors(policyPeer.NamespaceSelector, policyPeer.PodSelector)...)
	} else if policyPeer.PodSelector != nil {
		// Add pod IPs selected by PodSelector (can't be nil) that are inside all namespaces
		whitelist = append(whitelist, getPodIPListFromPodSelector(metav1.NamespaceAll, policyPeer.PodSelector)...)
	}

	return whitelist, blacklist
}

func getIngressPolicies(nwPolicy *networkingv1.NetworkPolicy) (types.Policy, types.Policy) {
	// Make ingress whitelist and blacklist
	whitelist := make([]string, 0)
	blacklist := make([]string, 0)

	// For every egress rule, add IPs to whitelist and blacklist
	for _, rule := range nwPolicy.Spec.Ingress {
		// Get port list
		//portList := getPortList(rule.Ports)

		for _, from := range rule.From {
			peerWhitelist, peerBlacklist := getLists(from)
			whitelist = append(whitelist, peerWhitelist...)
			blacklist = append(blacklist, peerBlacklist...)
		}
	}

	return newPolicy(whitelist, nwPolicy.Name, "input", "whitelist"), newPolicy(blacklist, nwPolicy.Name, "input", "blacklist")
}

func getEgressPolicies(nwPolicy *networkingv1.NetworkPolicy) (types.Policy, types.Policy) {
	// Make egress whitelist and blacklist
	whitelist := make([]string, 0)
	blacklist := make([]string, 0)

	for _, rule := range nwPolicy.Spec.Egress {
		// Get port list
		//portList := getPortList(rule.Ports)

		for _, to := range rule.To {
			peerWhitelist, peerBlacklist := getLists(to)
			whitelist = append(whitelist, peerWhitelist...)
			blacklist = append(blacklist, peerBlacklist...)
		}
	}

	return newPolicy(whitelist, nwPolicy.Name, "output", "whitelist"), newPolicy(blacklist, nwPolicy.Name, "output", "blacklist")
}

func parseList(list []string) []types.Element {
	elements := make([]types.Element, 0)

	// Parse IP list as an array of { "data": IP } elements
	for _, ip := range list {
		elements = append(elements, types.Element{
			Data: ip,
		})
	}

	return elements
}

func getPortList(policyPorts []networkingv1.NetworkPolicyPort) string {
	portList := make([]string, 0)

	for _, value := range policyPorts {
		portList = append(portList, value.Port.StrVal)
	}

	return strings.Join(portList, ",")
}

func getPodIPListFromSelectors(labelNamespaceSelector *metav1.LabelSelector, labelPodSelector *metav1.LabelSelector) []string {
	podIPList := make([]string, 0)

	// Fill list options
	opts := metav1.ListOptions{
		LabelSelector: labelNamespaceSelector.String(),
	}

	// Get namespaces selected by labelNamespaceSelector
	namespaces, _ := clientset.CoreV1().Namespaces().List(context.TODO(), opts)
	for _, namespace := range namespaces.Items {
		// Get pods inside this namespace that are also matched by labelPodSelector
		podIPList = append(podIPList, getPodIPListFromPodSelector(namespace.Name, labelPodSelector)...)
	}

	return podIPList
}

func getPodIPListFromPodSelector(namespace string, labelPodSelector *metav1.LabelSelector) []string {
	podIPList := make([]string, 0)

	// Fill list options
	opts := metav1.ListOptions{}
	if labelPodSelector != nil {
		opts.LabelSelector = labelPodSelector.String()
	}

	// Get selected pods in this namespace
	pods, _ := clientset.CoreV1().Pods(namespace).List(context.TODO(), opts)
	for _, pod := range pods.Items {
		// Get IPs from this pod
		for _, podIP := range pod.Status.PodIPs {
			podIPList = append(podIPList, podIP.IP)
		}
	}

	return podIPList
}

func formatPolicyName(name string, policyType string, listType string) string {
	return fmt.Sprintf("%s--%s-%s", name, policyType, listType)
}
