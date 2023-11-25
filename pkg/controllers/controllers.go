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

package controllers

import (
<<<<<<< HEAD
	"k8s.io/client-go/kubernetes"
=======
	"context"

	"github.com/aws/aws-sdk-go/aws/session"
	servicesqs "github.com/aws/aws-sdk-go/service/sqs"
	"github.com/samber/lo"
>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13
	"k8s.io/utils/clock"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/aws/karpenter-core/pkg/cloudprovider"
	"github.com/aws/karpenter-core/pkg/controllers/disruption"
	"github.com/aws/karpenter-core/pkg/controllers/disruption/orchestration"
	"github.com/aws/karpenter-core/pkg/controllers/leasegarbagecollection"
	metricsnode "github.com/aws/karpenter-core/pkg/controllers/metrics/node"
	metricsnodepool "github.com/aws/karpenter-core/pkg/controllers/metrics/nodepool"
	metricspod "github.com/aws/karpenter-core/pkg/controllers/metrics/pod"
	"github.com/aws/karpenter-core/pkg/controllers/node/termination"
	"github.com/aws/karpenter-core/pkg/controllers/node/termination/terminator"
	nodeclaimconsistency "github.com/aws/karpenter-core/pkg/controllers/nodeclaim/consistency"
	nodeclaimdisruption "github.com/aws/karpenter-core/pkg/controllers/nodeclaim/disruption"
	nodeclaimgarbagecollection "github.com/aws/karpenter-core/pkg/controllers/nodeclaim/garbagecollection"
	nodeclaimlifecycle "github.com/aws/karpenter-core/pkg/controllers/nodeclaim/lifecycle"
	nodeclaimtermination "github.com/aws/karpenter-core/pkg/controllers/nodeclaim/termination"
	nodepoolcounter "github.com/aws/karpenter-core/pkg/controllers/nodepool/counter"
	nodepoolhash "github.com/aws/karpenter-core/pkg/controllers/nodepool/hash"
	"github.com/aws/karpenter-core/pkg/controllers/provisioning"
	"github.com/aws/karpenter-core/pkg/controllers/state"
	"github.com/aws/karpenter-core/pkg/controllers/state/informer"
	"github.com/aws/karpenter-core/pkg/events"
	"github.com/aws/karpenter-core/pkg/operator/controller"
<<<<<<< HEAD
)

func NewControllers(
	clock clock.Clock,
	kubeClient client.Client,
	kubernetesInterface kubernetes.Interface,
	cluster *state.Cluster,
	recorder events.Recorder,
	cloudProvider cloudprovider.CloudProvider,
) []controller.Controller {

	p := provisioning.NewProvisioner(kubeClient, kubernetesInterface.CoreV1(), recorder, cloudProvider, cluster)
	evictionQueue := terminator.NewQueue(kubernetesInterface.CoreV1(), recorder)
	disruptionQueue := orchestration.NewQueue(kubeClient, recorder, cluster, clock, p)

	return []controller.Controller{
		p, evictionQueue, disruptionQueue,
		disruption.NewController(clock, kubeClient, p, cloudProvider, recorder, cluster, disruptionQueue),
		provisioning.NewController(kubeClient, p, recorder),
		nodepoolhash.NewNodePoolController(kubeClient),
		informer.NewDaemonSetController(kubeClient, cluster),
		informer.NewNodeController(kubeClient, cluster),
		informer.NewPodController(kubeClient, cluster),
		informer.NewNodePoolController(kubeClient, cluster),
		informer.NewNodeClaimController(kubeClient, cluster),
		termination.NewController(kubeClient, cloudProvider, terminator.NewTerminator(clock, kubeClient, evictionQueue), recorder),
		metricspod.NewController(kubeClient),
		metricsnodepool.NewController(kubeClient),
		metricsnode.NewController(cluster),
		nodepoolcounter.NewNodePoolController(kubeClient, cluster),
		nodeclaimconsistency.NewNodeClaimController(clock, kubeClient, recorder, cloudProvider),
		nodeclaimlifecycle.NewNodeClaimController(clock, kubeClient, cloudProvider, recorder),
		nodeclaimgarbagecollection.NewController(clock, kubeClient, cloudProvider),
		nodeclaimtermination.NewNodeClaimController(kubeClient, cloudProvider),
		nodeclaimdisruption.NewNodeClaimController(clock, kubeClient, cluster, cloudProvider),
		leasegarbagecollection.NewController(kubeClient),
	}
=======
	"github.com/aws/karpenter/pkg/cache"
	"github.com/aws/karpenter/pkg/cloudprovider"
	"github.com/aws/karpenter/pkg/controllers/interruption"
	nodeclaimgarbagecollection "github.com/aws/karpenter/pkg/controllers/nodeclaim/garbagecollection"
	nodeclaimtagging "github.com/aws/karpenter/pkg/controllers/nodeclaim/tagging"
	"github.com/aws/karpenter/pkg/controllers/nodeclass"
	"github.com/aws/karpenter/pkg/operator/options"
	"github.com/aws/karpenter/pkg/providers/amifamily"
	"github.com/aws/karpenter/pkg/providers/instance"
	"github.com/aws/karpenter/pkg/providers/instanceprofile"
	"github.com/aws/karpenter/pkg/providers/pricing"
	"github.com/aws/karpenter/pkg/providers/securitygroup"
	"github.com/aws/karpenter/pkg/providers/sqs"
	"github.com/aws/karpenter/pkg/providers/subnet"
)

func NewControllers(ctx context.Context, sess *session.Session, clk clock.Clock, kubeClient client.Client, recorder events.Recorder,
	unavailableOfferings *cache.UnavailableOfferings, cloudProvider *cloudprovider.CloudProvider, subnetProvider *subnet.Provider,
	securityGroupProvider *securitygroup.Provider, instanceProfileProvider *instanceprofile.Provider, instanceProvider *instance.Provider,
	pricingProvider *pricing.Provider, amiProvider *amifamily.Provider) []controller.Controller {

	controllers := []controller.Controller{
		nodeclass.NewNodeClassController(kubeClient, recorder, subnetProvider, securityGroupProvider, amiProvider, instanceProfileProvider),
		nodeclaimgarbagecollection.NewController(kubeClient, cloudProvider),
		nodeclaimtagging.NewController(kubeClient, instanceProvider),
	}
	if options.FromContext(ctx).InterruptionQueue != "" {
		controllers = append(controllers, interruption.NewController(kubeClient, clk, recorder, lo.Must(sqs.NewProvider(ctx, servicesqs.New(sess), options.FromContext(ctx).InterruptionQueue)), unavailableOfferings))
	}
	if options.FromContext(ctx).IsolatedVPC {
		logging.FromContext(ctx).Infof("assuming isolated VPC, pricing information will not be updated")
	} else {
		controllers = append(controllers, pricing.NewController(pricingProvider))
	}
	return controllers
>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13
}
