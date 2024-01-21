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

package options

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

<<<<<<< HEAD
<<<<<<< HEAD
	"github.com/samber/lo"
	cliflag "k8s.io/component-base/cli/flag"

	"github.com/aws/karpenter-core/pkg/apis/settings"
	"github.com/aws/karpenter-core/pkg/utils/env"
)

var (
	validLogLevels = []string{"", "debug", "info", "error"}

	Injectables = []Injectable{&Options{}}
)

type optionsKey struct{}

type FeatureGates struct {
	Drift    bool
	inputStr string
}

// Options contains all CLI flags / env vars for karpenter-core. It adheres to the options.Injectable interface.
type Options struct {
	ServiceName          string
	DisableWebhook       bool
	WebhookPort          int
	MetricsPort          int
	WebhookMetricsPort   int
	HealthProbePort      int
	KubeClientQPS        int
	KubeClientBurst      int
	EnableProfiling      bool
	EnableLeaderElection bool
	MemoryLimit          int64
	LogLevel             string
	BatchMaxDuration     time.Duration
	BatchIdleDuration    time.Duration
	FeatureGates         FeatureGates
=======
	"k8s.io/apimachinery/pkg/util/sets"

	coreoptions "github.com/aws/karpenter-core/pkg/operator/options"
	"github.com/aws/karpenter-core/pkg/utils/env"
	"github.com/aws/karpenter/pkg/apis/settings"
=======
	coreoptions "sigs.k8s.io/karpenter/pkg/operator/options"
	"sigs.k8s.io/karpenter/pkg/utils/env"
>>>>>>> 6ebba50ce424ccd5e33b3c84b4f10a8e78d54539
)

func init() {
	coreoptions.Injectables = append(coreoptions.Injectables, &Options{})
}

type optionsKey struct{}

type Options struct {
	AssumeRoleARN           string
	AssumeRoleDuration      time.Duration
	ClusterCABundle         string
	ClusterName             string
	ClusterEndpoint         string
	IsolatedVPC             bool
	VMMemoryOverheadPercent float64
	InterruptionQueue       string
	ReservedENIs            int
<<<<<<< HEAD
>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13

	setFlags map[string]bool
=======
>>>>>>> 6ebba50ce424ccd5e33b3c84b4f10a8e78d54539
}

<<<<<<< HEAD
type FlagSet struct {
	*flag.FlagSet
}

// BoolVarWithEnv defines a bool flag with a specified name, default value, usage string, and fallback environment
// variable.
func (fs *FlagSet) BoolVarWithEnv(p *bool, name string, envVar string, val bool, usage string) {
	*p = env.WithDefaultBool(envVar, val)
	fs.BoolFunc(name, usage, func(val string) error {
		if val != "true" && val != "false" {
			return fmt.Errorf("%q is not a valid value, must be true or false", val)
		}
		*p = (val) == "true"
		return nil
	})
}

func (o *Options) AddFlags(fs *FlagSet) {
	fs.StringVar(&o.ServiceName, "karpenter-service", env.WithDefaultString("KARPENTER_SERVICE", ""), "The Karpenter Service name for the dynamic webhook certificate")
	fs.BoolVarWithEnv(&o.DisableWebhook, "disable-webhook", "DISABLE_WEBHOOK", false, "Disable the admission and validation webhooks")
	fs.IntVar(&o.WebhookPort, "webhook-port", env.WithDefaultInt("WEBHOOK_PORT", 8443), "The port the webhook endpoint binds to for validation and mutation of resources")
	fs.IntVar(&o.MetricsPort, "metrics-port", env.WithDefaultInt("METRICS_PORT", 8000), "The port the metric endpoint binds to for operating metrics about the controller itself")
	fs.IntVar(&o.WebhookMetricsPort, "webhook-metrics-port", env.WithDefaultInt("WEBHOOK_METRICS_PORT", 8001), "The port the webhook metric endpoing binds to for operating metrics about the webhook")
	fs.IntVar(&o.HealthProbePort, "health-probe-port", env.WithDefaultInt("HEALTH_PROBE_PORT", 8081), "The port the health probe endpoint binds to for reporting controller health")
	fs.IntVar(&o.KubeClientQPS, "kube-client-qps", env.WithDefaultInt("KUBE_CLIENT_QPS", 200), "The smoothed rate of qps to kube-apiserver")
	fs.IntVar(&o.KubeClientBurst, "kube-client-burst", env.WithDefaultInt("KUBE_CLIENT_BURST", 300), "The maximum allowed burst of queries to the kube-apiserver")
	fs.BoolVarWithEnv(&o.EnableProfiling, "enable-profiling", "ENABLE_PROFILING", false, "Enable the profiling on the metric endpoint")
	fs.BoolVarWithEnv(&o.EnableLeaderElection, "leader-elect", "LEADER_ELECT", true, "Start leader election client and gain leadership before executing the main loop. Enable this when running replicated components for high availability.")
	fs.Int64Var(&o.MemoryLimit, "memory-limit", env.WithDefaultInt64("MEMORY_LIMIT", -1), "Memory limit on the container running the controller. The GC soft memory limit is set to 90% of this value.")
	fs.StringVar(&o.LogLevel, "log-level", env.WithDefaultString("LOG_LEVEL", ""), "Log verbosity level. Can be one of 'debug', 'info', or 'error'")
	fs.DurationVar(&o.BatchMaxDuration, "batch-max-duration", env.WithDefaultDuration("BATCH_MAX_DURATION", 10*time.Second), "The maximum length of a batch window. The longer this is, the more pods we can consider for provisioning at one time which usually results in fewer but larger nodes.")
	fs.DurationVar(&o.BatchIdleDuration, "batch-idle-duration", env.WithDefaultDuration("BATCH_IDLE_DURATION", time.Second), "The maximum amount of time with no new pending pods that if exceeded ends the current batching window. If pods arrive faster than this time, the batching window will be extended up to the maxDuration. If they arrive slower, the pods will be batched separately.")
	fs.StringVar(&o.FeatureGates.inputStr, "feature-gates", env.WithDefaultString("FEATURE_GATES", "Drift=false"), "Optional features can be enabled / disabled using feature gates. Current options are: Drift")
}

func (o *Options) Parse(fs *FlagSet, args ...string) error {
=======
func (o *Options) AddFlags(fs *coreoptions.FlagSet) {
	fs.StringVar(&o.AssumeRoleARN, "assume-role-arn", env.WithDefaultString("ASSUME_ROLE_ARN", ""), "Role to assume for calling AWS services.")
	fs.DurationVar(&o.AssumeRoleDuration, "assume-role-duration", env.WithDefaultDuration("ASSUME_ROLE_DURATION", 15*time.Minute), "Duration of assumed credentials in minutes. Default value is 15 minutes. Not used unless aws.assumeRole set.")
	fs.StringVar(&o.ClusterCABundle, "cluster-ca-bundle", env.WithDefaultString("CLUSTER_CA_BUNDLE", ""), "Cluster CA bundle for nodes to use for TLS connections with the API server. If not set, this is taken from the controller's TLS configuration.")
	fs.StringVar(&o.ClusterName, "cluster-name", env.WithDefaultString("CLUSTER_NAME", ""), "[REQUIRED] The kubernetes cluster name for resource discovery.")
	fs.StringVar(&o.ClusterEndpoint, "cluster-endpoint", env.WithDefaultString("CLUSTER_ENDPOINT", ""), "The external kubernetes cluster endpoint for new nodes to connect with. If not specified, will discover the cluster endpoint using DescribeCluster API.")
	fs.BoolVarWithEnv(&o.IsolatedVPC, "isolated-vpc", "ISOLATED_VPC", false, "If true, then assume we can't reach AWS services which don't have a VPC endpoint. This also has the effect of disabling look-ups to the AWS pricing endpoint.")
	fs.Float64Var(&o.VMMemoryOverheadPercent, "vm-memory-overhead-percent", env.WithDefaultFloat64("VM_MEMORY_OVERHEAD_PERCENT", 0.075), "The VM memory overhead as a percent that will be subtracted from the total memory for all instance types.")
	fs.StringVar(&o.InterruptionQueue, "interruption-queue", env.WithDefaultString("INTERRUPTION_QUEUE", ""), "Interruption queue is disabled if not specified. Enabling interruption handling may require additional permissions on the controller service account. Additional permissions are outlined in the docs.")
	fs.IntVar(&o.ReservedENIs, "reserved-enis", env.WithDefaultInt("RESERVED_ENIS", 0), "Reserved ENIs are not included in the calculations for max-pods or kube-reserved. This is most often used in the VPC CNI custom networking setup https://docs.aws.amazon.com/eks/latest/userguide/cni-custom-network.html.")
}

func (o *Options) Parse(fs *coreoptions.FlagSet, args ...string) error {
>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}
		return fmt.Errorf("parsing flags, %w", err)
	}
<<<<<<< HEAD

<<<<<<< HEAD
	if !lo.Contains(validLogLevels, o.LogLevel) {
		return fmt.Errorf("validating cli flags / env vars, invalid log level %q", o.LogLevel)
	}
	gates, err := ParseFeatureGates(o.FeatureGates.inputStr)
	if err != nil {
		return fmt.Errorf("parsing feature gates, %w", err)
	}
	o.FeatureGates = gates

	o.setFlags = map[string]bool{}
	fs.VisitAll(func(f *flag.Flag) {
		// NOTE: This assumes all CLI flags can be transformed into their corresponding environment variable. If a cli
		// flag / env var pair does not follow this pattern, this will break.
		envName := strings.ReplaceAll(strings.ToUpper(f.Name), "-", "_")
		_, ok := os.LookupEnv(envName)
		o.setFlags[f.Name] = ok
	})
	fs.Visit(func(f *flag.Flag) {
		o.setFlags[f.Name] = true
	})

=======
	// Check if each option has been set. This is a little brute force and better options might exist,
	// but this only needs to be here for one version
	o.setFlags = map[string]bool{}
	cliFlags := sets.New[string]()
	fs.Visit(func(f *flag.Flag) {
		cliFlags.Insert(f.Name)
	})
	fs.VisitAll(func(f *flag.Flag) {
		envName := strings.ReplaceAll(strings.ToUpper(f.Name), "-", "_")
		_, ok := os.LookupEnv(envName)
		o.setFlags[f.Name] = ok || cliFlags.Has(f.Name)
	})

	if err := o.Validate(); err != nil {
		return fmt.Errorf("validating options, %w", err)
	}

>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13
=======
	if err := o.Validate(); err != nil {
		return fmt.Errorf("validating options, %w", err)
	}
>>>>>>> 6ebba50ce424ccd5e33b3c84b4f10a8e78d54539
	return nil
}

func (o *Options) ToContext(ctx context.Context) context.Context {
	return ToContext(ctx, o)
}

<<<<<<< HEAD
func (o *Options) MergeSettings(ctx context.Context) {
	s := settings.FromContext(ctx)
<<<<<<< HEAD
	if !o.setFlags["batch-max-duration"] {
		o.BatchMaxDuration = s.BatchMaxDuration
	}
	if !o.setFlags["batch-idle-duration"] {
		o.BatchIdleDuration = s.BatchIdleDuration
	}
	if !o.setFlags["feature-gates"] {
		o.FeatureGates.Drift = s.DriftEnabled
	}
}

func ParseFeatureGates(gateStr string) (FeatureGates, error) {
	gateMap := map[string]bool{}
	gates := FeatureGates{}

	// Parses feature gates with the upstream mechanism. This is meant to be used with flag directly but this enables
	// simple merging with environment vars.
	if err := cliflag.NewMapStringBool(&gateMap).Set(gateStr); err != nil {
		return gates, err
	}
	if val, ok := gateMap["Drift"]; ok {
		gates.Drift = val
	}

	return gates, nil
=======
	mergeField(&o.AssumeRoleARN, s.AssumeRoleARN, o.setFlags["assume-role-arn"])
	mergeField(&o.AssumeRoleDuration, s.AssumeRoleDuration, o.setFlags["assume-role-duration"])
	mergeField(&o.ClusterCABundle, s.ClusterCABundle, o.setFlags["cluster-ca-bundle"])
	mergeField(&o.ClusterName, s.ClusterName, o.setFlags["cluster-name"])
	mergeField(&o.ClusterEndpoint, s.ClusterEndpoint, o.setFlags["cluster-endpoint"])
	mergeField(&o.IsolatedVPC, s.IsolatedVPC, o.setFlags["isolated-vpc"])
	mergeField(&o.VMMemoryOverheadPercent, s.VMMemoryOverheadPercent, o.setFlags["vm-memory-overhead-percent"])
	mergeField(&o.InterruptionQueue, s.InterruptionQueueName, o.setFlags["interruption-queue"])
	mergeField(&o.ReservedENIs, s.ReservedENIs, o.setFlags["reserved-enis"])
	if err := o.validateRequiredFields(); err != nil {
		panic(fmt.Errorf("checking required fields, %w", err))
	}
>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13
}

=======
>>>>>>> 6ebba50ce424ccd5e33b3c84b4f10a8e78d54539
func ToContext(ctx context.Context, opts *Options) context.Context {
	return context.WithValue(ctx, optionsKey{}, opts)
}

func FromContext(ctx context.Context) *Options {
	retval := ctx.Value(optionsKey{})
	if retval == nil {
<<<<<<< HEAD
		// This is a developer error if this happens, so we should panic
		panic("options doesn't exist in context")
	}
	return retval.(*Options)
}
=======
		return nil
	}
	return retval.(*Options)
}
<<<<<<< HEAD

// Note: Separated out to help with cyclomatic complexity check
func mergeField[T any](dest *T, src T, isDestSet bool) {
	if !isDestSet {
		*dest = src
	}
}
>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13
=======
>>>>>>> 6ebba50ce424ccd5e33b3c84b4f10a8e78d54539
