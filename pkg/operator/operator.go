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

package operator

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"runtime/debug"
	"sync"
	"time"

<<<<<<< HEAD
	coordinationv1 "k8s.io/api/coordination/v1"
=======
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awsclient "github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/patrickmn/go-cache"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13

	"github.com/go-logr/zapr"
	"github.com/samber/lo"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/util/flowcontrol"
	"k8s.io/utils/clock"
	knativeinjection "knative.dev/pkg/injection"
	knativelogging "knative.dev/pkg/logging"
	"knative.dev/pkg/signals"
	"knative.dev/pkg/system"
	"knative.dev/pkg/webhook"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"

<<<<<<< HEAD
	"github.com/aws/karpenter-core/pkg/apis"
	"github.com/aws/karpenter-core/pkg/apis/v1beta1"
	"github.com/aws/karpenter-core/pkg/events"
	"github.com/aws/karpenter-core/pkg/operator/controller"
	"github.com/aws/karpenter-core/pkg/operator/injection"
	"github.com/aws/karpenter-core/pkg/operator/logging"
	"github.com/aws/karpenter-core/pkg/operator/options"
	"github.com/aws/karpenter-core/pkg/operator/scheme"
	"github.com/aws/karpenter-core/pkg/webhooks"
)

const (
	appName   = "karpenter"
	component = "controller"
)
=======
	coreapis "github.com/aws/karpenter-core/pkg/apis"
	"github.com/aws/karpenter-core/pkg/apis/v1alpha5"
	corev1beta1 "github.com/aws/karpenter-core/pkg/apis/v1beta1"
	"github.com/aws/karpenter-core/pkg/operator"
	"github.com/aws/karpenter-core/pkg/operator/scheme"
	"github.com/aws/karpenter/pkg/apis"
	awscache "github.com/aws/karpenter/pkg/cache"
	"github.com/aws/karpenter/pkg/operator/options"
	"github.com/aws/karpenter/pkg/providers/amifamily"
	"github.com/aws/karpenter/pkg/providers/instance"
	"github.com/aws/karpenter/pkg/providers/instanceprofile"
	"github.com/aws/karpenter/pkg/providers/instancetype"
	"github.com/aws/karpenter/pkg/providers/launchtemplate"
	"github.com/aws/karpenter/pkg/providers/pricing"
	"github.com/aws/karpenter/pkg/providers/securitygroup"
	"github.com/aws/karpenter/pkg/providers/subnet"
	"github.com/aws/karpenter/pkg/providers/version"
)

func init() {
	lo.Must0(apis.AddToScheme(scheme.Scheme))
	v1alpha5.NormalizedLabels = lo.Assign(v1alpha5.NormalizedLabels, map[string]string{"topology.ebs.csi.aws.com/zone": corev1.LabelTopologyZone})
	corev1beta1.NormalizedLabels = lo.Assign(corev1beta1.NormalizedLabels, map[string]string{"topology.ebs.csi.aws.com/zone": corev1.LabelTopologyZone})
	coreapis.Settings = append(coreapis.Settings, apis.Settings...)
}

// Operator is injected into the AWS CloudProvider's factories
type Operator struct {
	*operator.Operator
>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13

// Version is the karpenter app version injected during compilation
// when using the Makefile
var Version = "unspecified"

type Operator struct {
	manager.Manager

	KubernetesInterface kubernetes.Interface
	EventRecorder       events.Recorder
	Clock               clock.Clock

	webhooks []knativeinjection.ControllerConstructor
}

// NewOperator instantiates a controller manager or panics
func NewOperator() (context.Context, *Operator) {
	// Root Context
	ctx := signals.NewContext()
	ctx = knativeinjection.WithNamespaceScope(ctx, system.Namespace())

	// Options
	ctx = injection.WithOptionsOrDie(ctx, options.Injectables...)

	// Make the Karpenter binary aware of the container memory limit
	// https://pkg.go.dev/runtime/debug#SetMemoryLimit
	if options.FromContext(ctx).MemoryLimit > 0 {
		newLimit := int64(float64(options.FromContext(ctx).MemoryLimit) * 0.9)
		debug.SetMemoryLimit(newLimit)
	}

<<<<<<< HEAD
	// Webhook
	ctx = webhook.WithOptions(ctx, webhook.Options{
		Port:        options.FromContext(ctx).WebhookPort,
		ServiceName: options.FromContext(ctx).ServiceName,
		SecretName:  fmt.Sprintf("%s-cert", options.FromContext(ctx).ServiceName),
		GracePeriod: 5 * time.Second,
	})

	// Client Config
	config := controllerruntime.GetConfigOrDie()
	config.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(float32(options.FromContext(ctx).KubeClientQPS), options.FromContext(ctx).KubeClientBurst)
	config.UserAgent = fmt.Sprintf("%s/%s", appName, Version)

	// Client
	kubernetesInterface := kubernetes.NewForConfigOrDie(config)

	// Inject settings from the ConfigMap(s) into the context
	ctx = injection.WithSettingsOrDie(ctx, kubernetesInterface, apis.Settings...)

	// Temporarily merge settings into options until configmap is removed
	// Note: injectables are pointer to those already in context
	for _, o := range options.Injectables {
		o.MergeSettings(ctx)
=======
	if assumeRoleARN := options.FromContext(ctx).AssumeRoleARN; assumeRoleARN != "" {
		config.Credentials = stscreds.NewCredentials(session.Must(session.NewSession()), assumeRoleARN,
			func(provider *stscreds.AssumeRoleProvider) { setDurationAndExpiry(ctx, provider) })
>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13
	}

	// Logging
	logger := logging.NewLogger(ctx, component, kubernetesInterface)
	ctx = knativelogging.WithLogger(ctx, logger)
	logging.ConfigureGlobalLoggers(ctx)

	knativelogging.FromContext(ctx).With("version", Version).Debugf("discovered karpenter version")

	// Manager
	mgrOpts := controllerruntime.Options{
		Logger:                     logging.IgnoreDebugEvents(zapr.NewLogger(logger.Desugar())),
		LeaderElection:             options.FromContext(ctx).EnableLeaderElection,
		LeaderElectionID:           "karpenter-leader-election",
		LeaderElectionResourceLock: resourcelock.LeasesResourceLock,
		LeaderElectionNamespace:    system.Namespace(),
		Scheme:                     scheme.Scheme,
		Metrics: server.Options{
			BindAddress: fmt.Sprintf(":%d", options.FromContext(ctx).MetricsPort),
		},
		HealthProbeBindAddress: fmt.Sprintf(":%d", options.FromContext(ctx).HealthProbePort),
		BaseContext: func() context.Context {
			ctx := context.Background()
			ctx = knativelogging.WithLogger(ctx, logger)
			ctx = injection.WithSettingsOrDie(ctx, kubernetesInterface, apis.Settings...)
			ctx = injection.WithOptionsOrDie(ctx, options.Injectables...)
			for _, o := range options.Injectables {
				o.MergeSettings(ctx)
			}
			return ctx
		},
		Cache: cache.Options{
			ByObject: map[client.Object]cache.ByObject{
				&coordinationv1.Lease{}: {
					Field: fields.SelectorFromSet(fields.Set{"metadata.namespace": "kube-node-lease"}),
				},
			},
		},
	}
	if options.FromContext(ctx).EnableProfiling {
		// TODO @joinnis: Investigate the mgrOpts.PprofBindAddress that would allow native support for pprof
		// On initial look, it seems like this native pprof doesn't support some of the routes that we have here
		// like "/debug/pprof/heap" or "/debug/pprof/block"
		mgrOpts.Metrics.ExtraHandlers = lo.Assign(mgrOpts.Metrics.ExtraHandlers, map[string]http.Handler{
			"/debug/pprof/":             http.HandlerFunc(pprof.Index),
			"/debug/pprof/cmdline":      http.HandlerFunc(pprof.Cmdline),
			"/debug/pprof/profile":      http.HandlerFunc(pprof.Profile),
			"/debug/pprof/symbol":       http.HandlerFunc(pprof.Symbol),
			"/debug/pprof/trace":        http.HandlerFunc(pprof.Trace),
			"/debug/pprof/allocs":       pprof.Handler("allocs"),
			"/debug/pprof/heap":         pprof.Handler("heap"),
			"/debug/pprof/block":        pprof.Handler("block"),
			"/debug/pprof/goroutine":    pprof.Handler("goroutine"),
			"/debug/pprof/threadcreate": pprof.Handler("threadcreate"),
		})
	}
	mgr, err := controllerruntime.NewManager(config, mgrOpts)
	mgr = lo.Must(mgr, err, "failed to setup manager")
	lo.Must0(mgr.GetFieldIndexer().IndexField(ctx, &v1.Pod{}, "spec.nodeName", func(o client.Object) []string {
		return []string{o.(*v1.Pod).Spec.NodeName}
	}), "failed to setup pod indexer")
	lo.Must0(mgr.GetFieldIndexer().IndexField(ctx, &v1.Node{}, "spec.providerID", func(o client.Object) []string {
		return []string{o.(*v1.Node).Spec.ProviderID}
	}), "failed to setup node provider id indexer")
	lo.Must0(mgr.GetFieldIndexer().IndexField(ctx, &v1beta1.NodeClaim{}, "status.providerID", func(o client.Object) []string {
		return []string{o.(*v1beta1.NodeClaim).Status.ProviderID}
	}), "failed to setup nodeclaim provider id indexer")

	lo.Must0(mgr.AddReadyzCheck("manager", func(req *http.Request) error {
		return lo.Ternary(mgr.GetCache().WaitForCacheSync(req.Context()), nil, fmt.Errorf("failed to sync caches"))
	}))
	lo.Must0(mgr.AddHealthzCheck("healthz", healthz.Ping))
	lo.Must0(mgr.AddReadyzCheck("readyz", healthz.Ping))

	lo.Must0(operator.Manager.GetFieldIndexer().IndexField(ctx, &corev1beta1.NodeClaim{}, "spec.nodeClassRef.name", func(o client.Object) []string {
		nc := o.(*corev1beta1.NodeClaim)
		if nc.Spec.NodeClassRef == nil {
			return []string{}
		}
		return []string{nc.Spec.NodeClassRef.Name}
	}), "failed to setup nodeclaim indexer")

	return ctx, &Operator{
		Manager:             mgr,
		KubernetesInterface: kubernetesInterface,
		EventRecorder:       events.NewRecorder(mgr.GetEventRecorderFor(appName)),
		Clock:               clock.RealClock{},
	}
}

<<<<<<< HEAD
func (o *Operator) WithControllers(ctx context.Context, controllers ...controller.Controller) *Operator {
	for _, c := range controllers {
		lo.Must0(c.Builder(ctx, o.Manager).Complete(c))
	}
	return o
=======
// withUserAgent adds a karpenter specific user-agent string to AWS session
func withUserAgent(sess *session.Session) *session.Session {
	userAgent := fmt.Sprintf("karpenter.sh-%s", operator.Version)
	sess.Handlers.Build.PushBack(request.MakeAddToUserAgentFreeFormHandler(userAgent))
	return sess
>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13
}

func (o *Operator) WithWebhooks(ctx context.Context, ctors ...knativeinjection.ControllerConstructor) *Operator {
	if !options.FromContext(ctx).DisableWebhook {
		o.webhooks = append(o.webhooks, ctors...)
		lo.Must0(o.Manager.AddReadyzCheck("webhooks", webhooks.HealthProbe(ctx)))
		lo.Must0(o.Manager.AddHealthzCheck("webhooks", webhooks.HealthProbe(ctx)))
	}
	return o
}

<<<<<<< HEAD
func (o *Operator) Start(ctx context.Context) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		lo.Must0(o.Manager.Start(ctx))
	}()
	if options.FromContext(ctx).DisableWebhook {
		knativelogging.FromContext(ctx).Infof("webhook disabled")
	} else {
		wg.Add(1)
		go func() {
			defer wg.Done()
			webhooks.Start(ctx, o.GetConfig(), o.KubernetesInterface, o.webhooks...)
		}()
	}
	wg.Wait()
=======
func ResolveClusterEndpoint(ctx context.Context, eksAPI eksiface.EKSAPI) (string, error) {
	clusterEndpointFromOptions := options.FromContext(ctx).ClusterEndpoint
	if clusterEndpointFromOptions != "" {
		return clusterEndpointFromOptions, nil // cluster endpoint is explicitly set
	}
	out, err := eksAPI.DescribeClusterWithContext(ctx, &eks.DescribeClusterInput{
		Name: aws.String(options.FromContext(ctx).ClusterName),
	})
	if err != nil {
		return "", fmt.Errorf("failed to resolve cluster endpoint, %w", err)
	}
	return *out.Cluster.Endpoint, nil
}

func getCABundle(ctx context.Context, restConfig *rest.Config) (*string, error) {
	// Discover CA Bundle from the REST client. We could alternatively
	// have used the simpler client-go InClusterConfig() method.
	// However, that only works when Karpenter is running as a Pod
	// within the same cluster it's managing.
	if caBundle := options.FromContext(ctx).ClusterCABundle; caBundle != "" {
		return lo.ToPtr(caBundle), nil
	}
	transportConfig, err := restConfig.TransportConfig()
	if err != nil {
		return nil, fmt.Errorf("discovering caBundle, loading transport config, %w", err)
	}
	_, err = transport.TLSConfigFor(transportConfig) // fills in CAData!
	if err != nil {
		return nil, fmt.Errorf("discovering caBundle, loading TLS config, %w", err)
	}
	return ptr.String(base64.StdEncoding.EncodeToString(transportConfig.TLS.CAData)), nil
}

func kubeDNSIP(ctx context.Context, kubernetesInterface kubernetes.Interface) (net.IP, error) {
	if kubernetesInterface == nil {
		return nil, fmt.Errorf("no K8s client provided")
	}
	dnsService, err := kubernetesInterface.CoreV1().Services("kube-system").Get(ctx, "kube-dns", metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	kubeDNSIP := net.ParseIP(dnsService.Spec.ClusterIP)
	if kubeDNSIP == nil {
		return nil, fmt.Errorf("parsing cluster IP")
	}
	return kubeDNSIP, nil
}

func setDurationAndExpiry(ctx context.Context, provider *stscreds.AssumeRoleProvider) {
	provider.Duration = options.FromContext(ctx).AssumeRoleDuration
	provider.ExpiryWindow = time.Duration(10) * time.Second
>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13
}
