package network

import (
	"testing"

	. "github.com/onsi/gomega"

	hyperv1 "github.com/openshift/hypershift/api/hypershift/v1beta1"

	operatorv1 "github.com/openshift/api/operator/v1"

	"k8s.io/utils/ptr"
)

func TestReconcileDefaultIngressController(t *testing.T) {
	vxlanPort := kubevirtDefaultVXLANPort
	genevePort := kubevirtDefaultGenevePort
	v4InternalSubnet := kubevirtDefaultV4InternalSubnet

	fakePort := uint32(11111)
	testsCases := []struct {
		name                string
		inputNetwork        *operatorv1.Network
		inputNetworkType    hyperv1.NetworkType
		inputPlatformType   hyperv1.PlatformType
		disableMultiNetwork bool
		hcp                 *hyperv1.HostedControlPlane
		expectedNetwork     *operatorv1.Network
	}{
		{
			name:                "KubeVirt with OVNKubernetes uses unique default geneve port",
			inputNetwork:        NetworkOperator(),
			inputNetworkType:    hyperv1.OVNKubernetes,
			inputPlatformType:   hyperv1.KubevirtPlatform,
			disableMultiNetwork: false,
			hcp:                 nil,
			expectedNetwork: &operatorv1.Network{
				ObjectMeta: NetworkOperator().ObjectMeta,
				Spec: operatorv1.NetworkSpec{
					OperatorSpec: operatorv1.OperatorSpec{
						ManagementState: "Managed",
					},
					DefaultNetwork: operatorv1.DefaultNetworkDefinition{
						OVNKubernetesConfig: &operatorv1.OVNKubernetesConfig{
							GenevePort:       &genevePort,
							V4InternalSubnet: v4InternalSubnet,
						},
					},
				},
			},
		},
		{
			name:                "KubeVirt with OpenShiftSDN uses unique default vxlan port",
			inputNetwork:        NetworkOperator(),
			inputNetworkType:    hyperv1.OpenShiftSDN,
			inputPlatformType:   hyperv1.KubevirtPlatform,
			disableMultiNetwork: false,
			hcp:                 nil,
			expectedNetwork: &operatorv1.Network{
				ObjectMeta: NetworkOperator().ObjectMeta,
				Spec: operatorv1.NetworkSpec{
					OperatorSpec: operatorv1.OperatorSpec{
						ManagementState: "Managed",
					},
					DefaultNetwork: operatorv1.DefaultNetworkDefinition{
						OpenShiftSDNConfig: &operatorv1.OpenShiftSDNConfig{
							VXLANPort: &vxlanPort,
						},
					},
				},
			},
		},
		{
			name: "KubeVirt with OpenShiftSDN does not set port when vxlan port already exists",
			inputNetwork: &operatorv1.Network{
				ObjectMeta: NetworkOperator().ObjectMeta,
				Spec: operatorv1.NetworkSpec{
					OperatorSpec: operatorv1.OperatorSpec{
						ManagementState: "Managed",
					},
					DefaultNetwork: operatorv1.DefaultNetworkDefinition{
						OpenShiftSDNConfig: &operatorv1.OpenShiftSDNConfig{
							VXLANPort: &fakePort,
						},
					},
				},
			},
			inputNetworkType:    hyperv1.OpenShiftSDN,
			inputPlatformType:   hyperv1.KubevirtPlatform,
			disableMultiNetwork: false,
			hcp:                 nil,
			expectedNetwork: &operatorv1.Network{
				ObjectMeta: NetworkOperator().ObjectMeta,
				Spec: operatorv1.NetworkSpec{
					OperatorSpec: operatorv1.OperatorSpec{
						ManagementState: "Managed",
					},
					DefaultNetwork: operatorv1.DefaultNetworkDefinition{
						OpenShiftSDNConfig: &operatorv1.OpenShiftSDNConfig{
							VXLANPort: &fakePort,
						},
					},
				},
			},
		},
		{
			name: "KubeVirt with OVNKubernetes when geneve port already exists",
			inputNetwork: &operatorv1.Network{
				ObjectMeta: NetworkOperator().ObjectMeta,
				Spec: operatorv1.NetworkSpec{
					OperatorSpec: operatorv1.OperatorSpec{
						ManagementState: "Managed",
					},
					DefaultNetwork: operatorv1.DefaultNetworkDefinition{
						OVNKubernetesConfig: &operatorv1.OVNKubernetesConfig{
							GenevePort:       &fakePort,
							V4InternalSubnet: kubevirtDefaultV4InternalSubnet,
						},
					},
				},
			},
			inputNetworkType:    hyperv1.OVNKubernetes,
			inputPlatformType:   hyperv1.KubevirtPlatform,
			disableMultiNetwork: false,
			hcp:                 nil,
			expectedNetwork: &operatorv1.Network{
				ObjectMeta: NetworkOperator().ObjectMeta,
				Spec: operatorv1.NetworkSpec{
					OperatorSpec: operatorv1.OperatorSpec{
						ManagementState: "Managed",
					},
					DefaultNetwork: operatorv1.DefaultNetworkDefinition{
						OVNKubernetesConfig: &operatorv1.OVNKubernetesConfig{
							GenevePort:       &fakePort,
							V4InternalSubnet: kubevirtDefaultV4InternalSubnet,
						},
					},
				},
			},
		},
		{
			name: "KubeVirt with OVNKubernetes when v4InternalSubnet already exists",
			inputNetwork: &operatorv1.Network{
				ObjectMeta: NetworkOperator().ObjectMeta,
				Spec: operatorv1.NetworkSpec{
					OperatorSpec: operatorv1.OperatorSpec{
						ManagementState: "Managed",
					},
					DefaultNetwork: operatorv1.DefaultNetworkDefinition{
						OVNKubernetesConfig: &operatorv1.OVNKubernetesConfig{
							V4InternalSubnet: "100.66.0.0/16",
							GenevePort:       &genevePort,
						},
					},
				},
			},
			inputNetworkType:    hyperv1.OVNKubernetes,
			inputPlatformType:   hyperv1.KubevirtPlatform,
			disableMultiNetwork: false,
			hcp:                 nil,
			expectedNetwork: &operatorv1.Network{
				ObjectMeta: NetworkOperator().ObjectMeta,
				Spec: operatorv1.NetworkSpec{
					OperatorSpec: operatorv1.OperatorSpec{
						ManagementState: "Managed",
					},
					DefaultNetwork: operatorv1.DefaultNetworkDefinition{
						OVNKubernetesConfig: &operatorv1.OVNKubernetesConfig{
							V4InternalSubnet: "100.66.0.0/16",
							GenevePort:       &genevePort,
						},
					},
				},
			},
		},

		{
			name:                "KubeVirt with non SDN network does not set unique vxlan port",
			inputNetwork:        NetworkOperator(),
			inputNetworkType:    "fake",
			inputPlatformType:   hyperv1.KubevirtPlatform,
			disableMultiNetwork: false,
			hcp:                 nil,
			expectedNetwork: &operatorv1.Network{
				ObjectMeta: NetworkOperator().ObjectMeta,
				Spec: operatorv1.NetworkSpec{
					OperatorSpec: operatorv1.OperatorSpec{
						ManagementState: "Managed",
					},
				},
			},
		},
		{
			name:                "AWS with SDN does not set unique vxlan port",
			inputNetwork:        NetworkOperator(),
			inputNetworkType:    hyperv1.OpenShiftSDN,
			inputPlatformType:   hyperv1.AWSPlatform,
			disableMultiNetwork: false,
			hcp:                 nil,
			expectedNetwork: &operatorv1.Network{
				ObjectMeta: NetworkOperator().ObjectMeta,
				Spec: operatorv1.NetworkSpec{
					OperatorSpec: operatorv1.OperatorSpec{
						ManagementState: "Managed",
					},
				},
			},
		},
		{
			name:                "DisableMultiNetwork sets disableMultiNetwork to true",
			inputNetwork:        NetworkOperator(),
			inputNetworkType:    hyperv1.Other,
			inputPlatformType:   hyperv1.AWSPlatform,
			disableMultiNetwork: true,
			hcp:                 nil,
			expectedNetwork: &operatorv1.Network{
				ObjectMeta: NetworkOperator().ObjectMeta,
				Spec: operatorv1.NetworkSpec{
					OperatorSpec: operatorv1.OperatorSpec{
						ManagementState: "Managed",
					},
					DisableMultiNetwork: ptr.To(true),
				},
			},
		},
		{
			name:                "DisableMultiNetwork false does not set disableMultiNetwork",
			inputNetwork:        NetworkOperator(),
			inputNetworkType:    hyperv1.Other,
			inputPlatformType:   hyperv1.AWSPlatform,
			disableMultiNetwork: false,
			hcp:                 nil,
			expectedNetwork: &operatorv1.Network{
				ObjectMeta: NetworkOperator().ObjectMeta,
				Spec: operatorv1.NetworkSpec{
					OperatorSpec: operatorv1.OperatorSpec{
						ManagementState: "Managed",
					},
				},
			},
		},
		{
			name:                "None with SDN does not set unique vxlan port",
			inputNetwork:        NetworkOperator(),
			inputNetworkType:    hyperv1.OpenShiftSDN,
			inputPlatformType:   hyperv1.NonePlatform,
			disableMultiNetwork: false,
			hcp:                 nil,
			expectedNetwork: &operatorv1.Network{
				ObjectMeta: NetworkOperator().ObjectMeta,
				Spec: operatorv1.NetworkSpec{
					OperatorSpec: operatorv1.OperatorSpec{
						ManagementState: "Managed",
					},
				},
			},
		},
		{
			name:                "IBM with SDN does not set unique vxlan port",
			inputNetwork:        NetworkOperator(),
			inputNetworkType:    hyperv1.OpenShiftSDN,
			inputPlatformType:   hyperv1.IBMCloudPlatform,
			disableMultiNetwork: false,
			hcp:                 nil,
			expectedNetwork: &operatorv1.Network{
				ObjectMeta: NetworkOperator().ObjectMeta,
				Spec: operatorv1.NetworkSpec{
					OperatorSpec: operatorv1.OperatorSpec{
						ManagementState: "Managed",
					},
				},
			},
		},
		{
			name:                "Azure with SDN does not set unique vxlan port",
			inputNetwork:        NetworkOperator(),
			inputNetworkType:    hyperv1.OpenShiftSDN,
			inputPlatformType:   hyperv1.AzurePlatform,
			disableMultiNetwork: false,
			hcp:                 nil,
			expectedNetwork: &operatorv1.Network{
				ObjectMeta: NetworkOperator().ObjectMeta,
				Spec: operatorv1.NetworkSpec{
					OperatorSpec: operatorv1.OperatorSpec{
						ManagementState: "Managed",
					},
				},
			},
		},
		{
			name:                "Agent with SDN does not set unique vxlan port",
			inputNetwork:        NetworkOperator(),
			inputNetworkType:    hyperv1.OpenShiftSDN,
			inputPlatformType:   hyperv1.AgentPlatform,
			disableMultiNetwork: false,
			hcp:                 nil,
			expectedNetwork: &operatorv1.Network{
				ObjectMeta: NetworkOperator().ObjectMeta,
				Spec: operatorv1.NetworkSpec{
					OperatorSpec: operatorv1.OperatorSpec{
						ManagementState: "Managed",
					},
				},
			},
		},
		{
			name:                "OVN with IPv4 internal subnets",
			inputNetwork:        NetworkOperator(),
			inputNetworkType:    hyperv1.OVNKubernetes,
			inputPlatformType:   hyperv1.AWSPlatform,
			disableMultiNetwork: false,
			hcp: &hyperv1.HostedControlPlane{
				Spec: hyperv1.HostedControlPlaneSpec{
					Networking: hyperv1.ClusterNetworking{
						NetworkType: hyperv1.OVNKubernetes,
						OVN: &hyperv1.OVNNetworking{
							IPv4: &hyperv1.OVNIPv4Config{
								InternalTransitSwitchSubnet: "100.90.0.0/16",
								InternalJoinSubnet:          "100.91.0.0/16",
							},
						},
					},
				},
			},
			expectedNetwork: &operatorv1.Network{
				ObjectMeta: NetworkOperator().ObjectMeta,
				Spec: operatorv1.NetworkSpec{
					OperatorSpec: operatorv1.OperatorSpec{
						ManagementState: "Managed",
					},
					DefaultNetwork: operatorv1.DefaultNetworkDefinition{
						Type: "OVNKubernetes",
						OVNKubernetesConfig: &operatorv1.OVNKubernetesConfig{
							IPv4: &operatorv1.IPv4OVNKubernetesConfig{
								InternalTransitSwitchSubnet: "100.90.0.0/16",
								InternalJoinSubnet:          "100.91.0.0/16",
							},
						},
					},
				},
			},
		},
		{
			name:                "OVN with IPv6 internal subnets",
			inputNetwork:        NetworkOperator(),
			inputNetworkType:    hyperv1.OVNKubernetes,
			inputPlatformType:   hyperv1.AWSPlatform,
			disableMultiNetwork: false,
			hcp: &hyperv1.HostedControlPlane{
				Spec: hyperv1.HostedControlPlaneSpec{
					Networking: hyperv1.ClusterNetworking{
						NetworkType: hyperv1.OVNKubernetes,
						OVN: &hyperv1.OVNNetworking{
							IPv6: &hyperv1.OVNIPv6Config{
								InternalTransitSwitchSubnet: "fd99::/64",
								InternalJoinSubnet:          "fd9a::/64",
							},
						},
					},
				},
			},
			expectedNetwork: &operatorv1.Network{
				ObjectMeta: NetworkOperator().ObjectMeta,
				Spec: operatorv1.NetworkSpec{
					OperatorSpec: operatorv1.OperatorSpec{
						ManagementState: "Managed",
					},
					DefaultNetwork: operatorv1.DefaultNetworkDefinition{
						Type: "OVNKubernetes",
						OVNKubernetesConfig: &operatorv1.OVNKubernetesConfig{
							IPv6: &operatorv1.IPv6OVNKubernetesConfig{
								InternalTransitSwitchSubnet: "fd99::/64",
								InternalJoinSubnet:          "fd9a::/64",
							},
						},
					},
				},
			},
		},
	}
	for _, tc := range testsCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGomegaWithT(t)
			ReconcileNetworkOperator(tc.inputNetwork, tc.inputNetworkType, tc.inputPlatformType, tc.disableMultiNetwork, tc.hcp)
			g.Expect(tc.inputNetwork).To(BeEquivalentTo(tc.expectedNetwork))
		})
	}
}
