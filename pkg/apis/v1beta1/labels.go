/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/sets"
)

<<<<<<< HEAD
// Well known labels and resources
const (
	ArchitectureAmd64    = "amd64"
	ArchitectureArm64    = "arm64"
	CapacityTypeSpot     = "spot"
	CapacityTypeOnDemand = "on-demand"
)
=======
func init() {
	v1beta1.RestrictedLabelDomains = v1beta1.RestrictedLabelDomains.Insert(RestrictedLabelDomains...)
	v1beta1.WellKnownLabels = v1beta1.WellKnownLabels.Insert(
		LabelInstanceHypervisor,
		LabelInstanceEncryptionInTransitSupported,
		LabelInstanceCategory,
		LabelInstanceFamily,
		LabelInstanceGeneration,
		LabelInstanceSize,
		LabelInstanceLocalNVME,
		LabelInstanceCPU,
		LabelInstanceMemory,
		LabelInstanceNetworkBandwidth,
		LabelInstanceGPUName,
		LabelInstanceGPUManufacturer,
		LabelInstanceGPUCount,
		LabelInstanceGPUMemory,
		LabelInstanceAcceleratorName,
		LabelInstanceAcceleratorManufacturer,
		LabelInstanceAcceleratorCount,
		v1.LabelWindowsBuild,
	)
}
>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13

// Karpenter specific domains and labels
const (
	NodePoolLabelKey        = Group + "/nodepool"
	NodeInitializedLabelKey = Group + "/initialized"
	NodeRegisteredLabelKey  = Group + "/registered"
	CapacityTypeLabelKey    = Group + "/capacity-type"
)

// Karpenter specific annotations
const (
	DoNotDisruptAnnotationKey          = Group + "/do-not-disrupt"
	ProviderCompatabilityAnnotationKey = CompatabilityGroup + "/provider"
	ManagedByAnnotationKey             = Group + "/managed-by"
	NodePoolHashAnnotationKey          = Group + "/nodepool-hash"
)

// Karpenter specific finalizers
const (
	TerminationFinalizer = Group + "/termination"
)

var (
<<<<<<< HEAD
	// RestrictedLabelDomains are either prohibited by the kubelet or reserved by karpenter
	RestrictedLabelDomains = sets.New(
		"kubernetes.io",
		"k8s.io",
		Group,
	)

	// LabelDomainExceptions are sub-domains of the RestrictedLabelDomains but allowed because
	// they are not used in a context where they may be passed as argument to kubelet.
	LabelDomainExceptions = sets.New(
		"kops.k8s.io",
		v1.LabelNamespaceSuffixNode,
		v1.LabelNamespaceNodeRestriction,
	)

	// WellKnownLabels are labels that belong to the RestrictedLabelDomains but allowed.
	// Karpenter is aware of these labels, and they can be used to further narrow down
	// the range of the corresponding values by either provisioner or pods.
	WellKnownLabels = sets.New(
		NodePoolLabelKey,
		v1.LabelTopologyZone,
		v1.LabelTopologyRegion,
		v1.LabelInstanceTypeStable,
		v1.LabelArchStable,
		v1.LabelOSStable,
		CapacityTypeLabelKey,
		v1.LabelWindowsBuild,
	)

	// RestrictedLabels are labels that should not be used
	// because they may interfere with the internal provisioning logic.
	RestrictedLabels = sets.New(
		v1.LabelHostname,
	)

	// NormalizedLabels translate aliased concepts into the controller's
	// WellKnownLabels. Pod requirements are translated for compatibility.
	NormalizedLabels = map[string]string{
		v1.LabelFailureDomainBetaZone:   v1.LabelTopologyZone,
		"beta.kubernetes.io/arch":       v1.LabelArchStable,
		"beta.kubernetes.io/os":         v1.LabelOSStable,
		v1.LabelInstanceType:            v1.LabelInstanceTypeStable,
		v1.LabelFailureDomainBetaRegion: v1.LabelTopologyRegion,
	}
=======
	AWSToKubeArchitectures = map[string]string{
		"x86_64":                  v1beta1.ArchitectureAmd64,
		v1beta1.ArchitectureArm64: v1beta1.ArchitectureArm64,
	}
	WellKnownArchitectures = sets.NewString(
		v1beta1.ArchitectureAmd64,
		v1beta1.ArchitectureArm64,
	)
	RestrictedLabelDomains = []string{
		Group,
	}
	RestrictedTagPatterns = []*regexp.Regexp{
		// Adheres to cluster name pattern matching as specified in the API spec
		// https://docs.aws.amazon.com/eks/latest/APIReference/API_CreateCluster.html
		regexp.MustCompile(`^kubernetes\.io/cluster/[0-9A-Za-z][A-Za-z0-9\-_]*$`),
		regexp.MustCompile(fmt.Sprintf("^%s$", regexp.QuoteMeta(v1alpha5.ProvisionerNameLabelKey))),
		regexp.MustCompile(fmt.Sprintf("^%s$", regexp.QuoteMeta(v1alpha5.MachineManagedByAnnotationKey))),
		regexp.MustCompile(fmt.Sprintf("^%s$", regexp.QuoteMeta(v1beta1.NodePoolLabelKey))),
		regexp.MustCompile(fmt.Sprintf("^%s$", regexp.QuoteMeta(v1beta1.ManagedByAnnotationKey))),
	}
	AMIFamilyBottlerocket                      = "Bottlerocket"
	AMIFamilyAL2                               = "AL2"
	AMIFamilyUbuntu                            = "Ubuntu"
	AMIFamilyWindows2019                       = "Windows2019"
	AMIFamilyWindows2022                       = "Windows2022"
	AMIFamilyCustom                            = "Custom"
	Windows2019                                = "2019"
	Windows2022                                = "2022"
	WindowsCore                                = "Core"
	Windows2019Build                           = "10.0.17763"
	Windows2022Build                           = "10.0.20348"
	ResourceNVIDIAGPU          v1.ResourceName = "nvidia.com/gpu"
	ResourceAMDGPU             v1.ResourceName = "amd.com/gpu"
	ResourceAWSNeuron          v1.ResourceName = "aws.amazon.com/neuron"
	ResourceHabanaGaudi        v1.ResourceName = "habana.ai/gaudi"
	ResourceAWSPodENI          v1.ResourceName = "vpc.amazonaws.com/pod-eni"
	ResourcePrivateIPv4Address v1.ResourceName = "vpc.amazonaws.com/PrivateIPv4Address"

	LabelNodeClass = Group + "/ec2nodeclass"

	LabelInstanceHypervisor                   = Group + "/instance-hypervisor"
	LabelInstanceEncryptionInTransitSupported = Group + "/instance-encryption-in-transit-supported"
	LabelInstanceCategory                     = Group + "/instance-category"
	LabelInstanceFamily                       = Group + "/instance-family"
	LabelInstanceGeneration                   = Group + "/instance-generation"
	LabelInstanceLocalNVME                    = Group + "/instance-local-nvme"
	LabelInstanceSize                         = Group + "/instance-size"
	LabelInstanceCPU                          = Group + "/instance-cpu"
	LabelInstanceMemory                       = Group + "/instance-memory"
	LabelInstanceNetworkBandwidth             = Group + "/instance-network-bandwidth"
	LabelInstanceGPUName                      = Group + "/instance-gpu-name"
	LabelInstanceGPUManufacturer              = Group + "/instance-gpu-manufacturer"
	LabelInstanceGPUCount                     = Group + "/instance-gpu-count"
	LabelInstanceGPUMemory                    = Group + "/instance-gpu-memory"
	LabelInstanceAcceleratorName              = Group + "/instance-accelerator-name"
	LabelInstanceAcceleratorManufacturer      = Group + "/instance-accelerator-manufacturer"
	LabelInstanceAcceleratorCount             = Group + "/instance-accelerator-count"
	AnnotationNodeClassHash                   = Group + "/ec2nodeclass-hash"
	AnnotationInstanceTagged                  = Group + "/tagged"
>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13
)

// IsRestrictedLabel returns an error if the label is restricted.
func IsRestrictedLabel(key string) error {
	if WellKnownLabels.Has(key) {
		return nil
	}
	if IsRestrictedNodeLabel(key) {
		return fmt.Errorf("label %s is restricted; specify a well known label: %v, or a custom label that does not use a restricted domain: %v", key, sets.List(WellKnownLabels), sets.List(RestrictedLabelDomains))
	}
	return nil
}

// IsRestrictedNodeLabel returns true if a node label should not be injected by Karpenter.
// They are either known labels that will be injected by cloud providers,
// or label domain managed by other software (e.g., kops.k8s.io managed by kOps).
func IsRestrictedNodeLabel(key string) bool {
	if WellKnownLabels.Has(key) {
		return true
	}
	labelDomain := GetLabelDomain(key)
	for exceptionLabelDomain := range LabelDomainExceptions {
		if strings.HasSuffix(labelDomain, exceptionLabelDomain) {
			return false
		}
	}
	for restrictedLabelDomain := range RestrictedLabelDomains {
		if strings.HasSuffix(labelDomain, restrictedLabelDomain) {
			return true
		}
	}
	return RestrictedLabels.Has(key)
}

func GetLabelDomain(key string) string {
	if parts := strings.SplitN(key, "/", 2); len(parts) == 2 {
		return parts[0]
	}
	return ""
}
