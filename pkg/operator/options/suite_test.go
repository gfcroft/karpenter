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

package options_test

import (
	"context"
	"flag"
<<<<<<< HEAD
=======
	"fmt"
>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
	. "knative.dev/pkg/logging/testing"

<<<<<<< HEAD
	"github.com/aws/karpenter-core/pkg/apis/settings"
	"github.com/aws/karpenter-core/pkg/operator/options"
	"github.com/aws/karpenter-core/pkg/test"
)

var ctx context.Context
var fs *options.FlagSet
var opts *options.Options

func TestOptions(t *testing.T) {
=======
	coreoptions "github.com/aws/karpenter-core/pkg/operator/options"
	"github.com/aws/karpenter/pkg/apis/settings"
	"github.com/aws/karpenter/pkg/operator/options"
	"github.com/aws/karpenter/pkg/test"
)

var ctx context.Context

func TestAPIs(t *testing.T) {
>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13
	ctx = TestContextWithLogger(t)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Options")
}

var _ = Describe("Options", func() {
	var envState map[string]string
	var environmentVariables = []string{
<<<<<<< HEAD
		"KARPENTER_SERVICE",
		"DISABLE_WEBHOOK",
		"WEBHOOK_PORT",
		"METRICS_PORT",
		"WEBHOOK_METRICS_PORT",
		"HEALTH_PROBE_PORT",
		"KUBE_CLIENT_BURST",
		"ENABLE_PROFILING",
		"LEADER_ELECT",
		"MEMORY_LIMIT",
		"LOG_LEVEL",
		"BATCH_MAX_DURATION",
		"BATCH_IDLE_DURATION",
		"FEATURE_GATES",
	}

=======
		"ASSUME_ROLE_ARN",
		"ASSUME_ROLE_DURATION",
		"CLUSTER_CA_BUNDLE",
		"CLUSTER_NAME",
		"CLUSTER_ENDPOINT",
		"ISOLATED_VPC",
		"VM_MEMORY_OVERHEAD_PERCENT",
		"INTERRUPTION_QUEUE",
		"RESERVED_ENIS",
	}

	var fs *coreoptions.FlagSet
	var opts *options.Options

>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13
	BeforeEach(func() {
		envState = map[string]string{}
		for _, ev := range environmentVariables {
			val, ok := os.LookupEnv(ev)
			if ok {
				envState[ev] = val
			}
			os.Unsetenv(ev)
		}
<<<<<<< HEAD

		fs = &options.FlagSet{
=======
		fs = &coreoptions.FlagSet{
>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13
			FlagSet: flag.NewFlagSet("karpenter", flag.ContinueOnError),
		}
		opts = &options.Options{}
		opts.AddFlags(fs)
<<<<<<< HEAD
=======

		// Inject default settings
		var err error
		ctx, err = (&settings.Settings{}).Inject(ctx, nil)
		Expect(err).To(BeNil())
>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13
	})

	AfterEach(func() {
		for _, ev := range environmentVariables {
			os.Unsetenv(ev)
		}
		for ev, val := range envState {
			os.Setenv(ev, val)
		}
	})

<<<<<<< HEAD
	Context("FeatureGates", func() {
		DescribeTable(
			"should successfully parse well formed feature gate strings",
			func(str string, driftVal bool) {
				gates, err := options.ParseFeatureGates(str)
				Expect(err).To(BeNil())
				Expect(gates.Drift).To(Equal(driftVal))
			},
			Entry("basic true", "Drift=true", true),
			Entry("basic false", "Drift=false", false),
			Entry("with whitespace", "Drift\t= false", false),
			Entry("multiple values", "Hello=true,Drift=false,World=true", false),
		)
	})

	Context("Parse", func() {
		It("should use the correct default values", func() {
			err := opts.Parse(fs)
			Expect(err).To(BeNil())
			expectOptionsMatch(opts, test.Options(test.OptionsFields{
				ServiceName:          lo.ToPtr(""),
				DisableWebhook:       lo.ToPtr(false),
				WebhookPort:          lo.ToPtr(8443),
				MetricsPort:          lo.ToPtr(8000),
				WebhookMetricsPort:   lo.ToPtr(8001),
				HealthProbePort:      lo.ToPtr(8081),
				KubeClientQPS:        lo.ToPtr(200),
				KubeClientBurst:      lo.ToPtr(300),
				EnableProfiling:      lo.ToPtr(false),
				EnableLeaderElection: lo.ToPtr(true),
				MemoryLimit:          lo.ToPtr[int64](-1),
				LogLevel:             lo.ToPtr(""),
				BatchMaxDuration:     lo.ToPtr(10 * time.Second),
				BatchIdleDuration:    lo.ToPtr(time.Second),
				FeatureGates: test.FeatureGates{
					Drift: lo.ToPtr(false),
				},
			}))
		})

		It("shouldn't overwrite CLI flags with environment variables", func() {
			err := opts.Parse(
				fs,
				"--karpenter-service", "cli",
				"--disable-webhook",
				"--webhook-port", "0",
				"--metrics-port", "0",
				"--webhook-metrics-port", "0",
				"--health-probe-port", "0",
				"--kube-client-qps", "0",
				"--kube-client-burst", "0",
				"--enable-profiling",
				"--leader-elect=false",
				"--memory-limit", "0",
				"--log-level", "debug",
				"--batch-max-duration", "5s",
				"--batch-idle-duration", "5s",
				"--feature-gates", "Drift=true",
			)
			Expect(err).To(BeNil())
			expectOptionsMatch(opts, test.Options(test.OptionsFields{
				ServiceName:          lo.ToPtr("cli"),
				DisableWebhook:       lo.ToPtr(true),
				WebhookPort:          lo.ToPtr(0),
				MetricsPort:          lo.ToPtr(0),
				WebhookMetricsPort:   lo.ToPtr(0),
				HealthProbePort:      lo.ToPtr(0),
				KubeClientQPS:        lo.ToPtr(0),
				KubeClientBurst:      lo.ToPtr(0),
				EnableProfiling:      lo.ToPtr(true),
				EnableLeaderElection: lo.ToPtr(false),
				MemoryLimit:          lo.ToPtr[int64](0),
				LogLevel:             lo.ToPtr("debug"),
				BatchMaxDuration:     lo.ToPtr(5 * time.Second),
				BatchIdleDuration:    lo.ToPtr(5 * time.Second),
				FeatureGates: test.FeatureGates{
					Drift: lo.ToPtr(true),
				},
			}))
		})

		It("should use environment variables when CLI flags aren't set", func() {
			os.Setenv("KARPENTER_SERVICE", "env")
			os.Setenv("DISABLE_WEBHOOK", "true")
			os.Setenv("WEBHOOK_PORT", "0")
			os.Setenv("METRICS_PORT", "0")
			os.Setenv("WEBHOOK_METRICS_PORT", "0")
			os.Setenv("HEALTH_PROBE_PORT", "0")
			os.Setenv("KUBE_CLIENT_QPS", "0")
			os.Setenv("KUBE_CLIENT_BURST", "0")
			os.Setenv("ENABLE_PROFILING", "true")
			os.Setenv("LEADER_ELECT", "false")
			os.Setenv("MEMORY_LIMIT", "0")
			os.Setenv("LOG_LEVEL", "debug")
			os.Setenv("BATCH_MAX_DURATION", "5s")
			os.Setenv("BATCH_IDLE_DURATION", "5s")
			os.Setenv("FEATURE_GATES", "Drift=true")
			fs = &options.FlagSet{
				FlagSet: flag.NewFlagSet("karpenter", flag.ContinueOnError),
			}
			opts.AddFlags(fs)
			err := opts.Parse(fs)
			Expect(err).To(BeNil())
			expectOptionsMatch(opts, test.Options(test.OptionsFields{
				ServiceName:          lo.ToPtr("env"),
				DisableWebhook:       lo.ToPtr(true),
				WebhookPort:          lo.ToPtr(0),
				MetricsPort:          lo.ToPtr(0),
				WebhookMetricsPort:   lo.ToPtr(0),
				HealthProbePort:      lo.ToPtr(0),
				KubeClientQPS:        lo.ToPtr(0),
				KubeClientBurst:      lo.ToPtr(0),
				EnableProfiling:      lo.ToPtr(true),
				EnableLeaderElection: lo.ToPtr(false),
				MemoryLimit:          lo.ToPtr[int64](0),
				LogLevel:             lo.ToPtr("debug"),
				BatchMaxDuration:     lo.ToPtr(5 * time.Second),
				BatchIdleDuration:    lo.ToPtr(5 * time.Second),
				FeatureGates: test.FeatureGates{
					Drift: lo.ToPtr(true),
				},
			}))
		})

		It("should correctly merge CLI flags and environment variables", func() {
			os.Setenv("WEBHOOK_PORT", "0")
			os.Setenv("METRICS_PORT", "0")
			os.Setenv("WEBHOOK_METRICS_PORT", "0")
			os.Setenv("HEALTH_PROBE_PORT", "0")
			os.Setenv("KUBE_CLIENT_QPS", "0")
			os.Setenv("KUBE_CLIENT_BURST", "0")
			os.Setenv("ENABLE_PROFILING", "true")
			os.Setenv("LEADER_ELECT", "false")
			os.Setenv("MEMORY_LIMIT", "0")
			os.Setenv("LOG_LEVEL", "debug")
			os.Setenv("BATCH_MAX_DURATION", "5s")
			os.Setenv("BATCH_IDLE_DURATION", "5s")
			os.Setenv("FEATURE_GATES", "Drift=true")
			fs = &options.FlagSet{
				FlagSet: flag.NewFlagSet("karpenter", flag.ContinueOnError),
			}
			opts.AddFlags(fs)
			err := opts.Parse(
				fs,
				"--karpenter-service", "cli",
				"--disable-webhook",
			)
			Expect(err).To(BeNil())
			expectOptionsMatch(opts, test.Options(test.OptionsFields{
				ServiceName:          lo.ToPtr("cli"),
				DisableWebhook:       lo.ToPtr(true),
				WebhookPort:          lo.ToPtr(0),
				MetricsPort:          lo.ToPtr(0),
				WebhookMetricsPort:   lo.ToPtr(0),
				HealthProbePort:      lo.ToPtr(0),
				KubeClientQPS:        lo.ToPtr(0),
				KubeClientBurst:      lo.ToPtr(0),
				EnableProfiling:      lo.ToPtr(true),
				EnableLeaderElection: lo.ToPtr(false),
				MemoryLimit:          lo.ToPtr[int64](0),
				LogLevel:             lo.ToPtr("debug"),
				BatchMaxDuration:     lo.ToPtr(5 * time.Second),
				BatchIdleDuration:    lo.ToPtr(5 * time.Second),
				FeatureGates: test.FeatureGates{
					Drift: lo.ToPtr(true),
				},
			}))
		})

		DescribeTable(
			"should correctly parse boolean values",
			func(arg string, expected bool) {
				err := opts.Parse(fs, arg)
				Expect(err).ToNot(HaveOccurred())
				Expect(opts.DisableWebhook).To(Equal(expected))
			},
			Entry("explicit true", "--disable-webhook=true", true),
			Entry("explicit false", "--disable-webhook=false", false),
			Entry("implicit true", "--disable-webhook", true),
			Entry("implicit false", "", false),
		)
	})

	Context("Merge", func() {
		BeforeEach(func() {
			ctx = settings.ToContext(ctx, &settings.Settings{
				BatchMaxDuration:  50 * time.Second,
				BatchIdleDuration: 50 * time.Second,
				DriftEnabled:      true,
			})
		})

		It("shouldn't overwrite BatchMaxDuration when specified by CLI", func() {
			err := opts.Parse(fs, "--batch-max-duration", "1s")
			Expect(err).To(BeNil())
			opts.MergeSettings(ctx)
			Expect(opts.BatchMaxDuration).To(Equal(time.Second))
		})
		It("shouldn't overwrite BatchIdleDuration when specified by CLI", func() {
			err := opts.Parse(fs, "--batch-idle-duration", "1s")
			Expect(err).To(BeNil())
			opts.MergeSettings(ctx)
			Expect(opts.BatchIdleDuration).To(Equal(time.Second))
		})
		It("shouldn't overwrite FeatureGates.Drift when specified by CLI", func() {
			err := opts.Parse(fs, "--feature-gates", "Drift=false")
			Expect(err).To(BeNil())
			opts.MergeSettings(ctx)
			Expect(opts.FeatureGates.Drift).To(BeFalse())
		})
		It("should use values from settings when not specified", func() {
			err := opts.Parse(fs, "--batch-max-duration", "1s", "--feature-gates", "Drift=false")
			Expect(err).To(BeNil())
			opts.MergeSettings(ctx)
			Expect(opts.BatchIdleDuration).To(Equal(50 * time.Second))
=======
	Context("Merging", func() {
		It("shouldn't overwrite options when all are set", func() {
			err := opts.Parse(
				fs,
				"--assume-role-arn", "options-cluster-role",
				"--assume-role-duration", "20m",
				"--cluster-ca-bundle", "options-bundle",
				"--cluster-name", "options-cluster",
				"--cluster-endpoint", "https://options-cluster",
				"--isolated-vpc",
				"--vm-memory-overhead-percent", "0.1",
				"--interruption-queue", "options-cluster",
				"--reserved-enis", "10",
			)
			Expect(err).ToNot(HaveOccurred())
			ctx = settings.ToContext(ctx, &settings.Settings{
				AssumeRoleARN:           "settings-cluster-role",
				AssumeRoleDuration:      time.Minute * 22,
				ClusterCABundle:         "settings-bundle",
				ClusterName:             "settings-cluster",
				ClusterEndpoint:         "https://settings-cluster",
				IsolatedVPC:             true,
				VMMemoryOverheadPercent: 0.05,
				InterruptionQueueName:   "settings-cluster",
				ReservedENIs:            8,
			})
			opts.MergeSettings(ctx)
			expectOptionsEqual(opts, test.Options(test.OptionsFields{
				AssumeRoleARN:           lo.ToPtr("options-cluster-role"),
				AssumeRoleDuration:      lo.ToPtr(20 * time.Minute),
				ClusterCABundle:         lo.ToPtr("options-bundle"),
				ClusterName:             lo.ToPtr("options-cluster"),
				ClusterEndpoint:         lo.ToPtr("https://options-cluster"),
				IsolatedVPC:             lo.ToPtr(true),
				VMMemoryOverheadPercent: lo.ToPtr[float64](0.1),
				InterruptionQueue:       lo.ToPtr("options-cluster"),
				ReservedENIs:            lo.ToPtr(10),
			}))

		})
		It("should overwrite options when none are set", func() {
			err := opts.Parse(fs)
			Expect(err).ToNot(HaveOccurred())
			ctx = settings.ToContext(ctx, &settings.Settings{
				AssumeRoleARN:           "settings-cluster-role",
				AssumeRoleDuration:      time.Minute * 22,
				ClusterCABundle:         "settings-bundle",
				ClusterName:             "settings-cluster",
				ClusterEndpoint:         "https://settings-cluster",
				IsolatedVPC:             true,
				VMMemoryOverheadPercent: 0.05,
				InterruptionQueueName:   "settings-cluster",
				ReservedENIs:            8,
			})
			opts.MergeSettings(ctx)
			expectOptionsEqual(opts, test.Options(test.OptionsFields{
				AssumeRoleARN:           lo.ToPtr("settings-cluster-role"),
				AssumeRoleDuration:      lo.ToPtr(22 * time.Minute),
				ClusterCABundle:         lo.ToPtr("settings-bundle"),
				ClusterName:             lo.ToPtr("settings-cluster"),
				ClusterEndpoint:         lo.ToPtr("https://settings-cluster"),
				IsolatedVPC:             lo.ToPtr(true),
				VMMemoryOverheadPercent: lo.ToPtr[float64](0.05),
				InterruptionQueue:       lo.ToPtr("settings-cluster"),
				ReservedENIs:            lo.ToPtr(8),
			}))

		})
		It("should correctly merge options and settings when mixed", func() {
			err := opts.Parse(
				fs,
				"--assume-role-arn", "options-cluster-role",
				"--cluster-ca-bundle", "options-bundle",
				"--cluster-name", "options-cluster",
				"--cluster-endpoint", "https://options-cluster",
				"--interruption-queue", "options-cluster",
			)
			Expect(err).ToNot(HaveOccurred())
			ctx = settings.ToContext(ctx, &settings.Settings{
				AssumeRoleARN:           "settings-cluster-role",
				AssumeRoleDuration:      time.Minute * 20,
				ClusterCABundle:         "settings-bundle",
				ClusterName:             "settings-cluster",
				ClusterEndpoint:         "https://settings-cluster",
				IsolatedVPC:             true,
				VMMemoryOverheadPercent: 0.1,
				InterruptionQueueName:   "settings-cluster",
				ReservedENIs:            10,
			})
			opts.MergeSettings(ctx)
			expectOptionsEqual(opts, test.Options(test.OptionsFields{
				AssumeRoleARN:           lo.ToPtr("options-cluster-role"),
				AssumeRoleDuration:      lo.ToPtr(20 * time.Minute),
				ClusterCABundle:         lo.ToPtr("options-bundle"),
				ClusterName:             lo.ToPtr("options-cluster"),
				ClusterEndpoint:         lo.ToPtr("https://options-cluster"),
				IsolatedVPC:             lo.ToPtr(true),
				VMMemoryOverheadPercent: lo.ToPtr[float64](0.1),
				InterruptionQueue:       lo.ToPtr("options-cluster"),
				ReservedENIs:            lo.ToPtr(10),
			}))
		})

		It("should correctly fallback to env vars when CLI flags aren't set", func() {
			os.Setenv("ASSUME_ROLE_ARN", "env-role")
			os.Setenv("ASSUME_ROLE_DURATION", "20m")
			os.Setenv("CLUSTER_CA_BUNDLE", "env-bundle")
			os.Setenv("CLUSTER_NAME", "env-cluster")
			os.Setenv("CLUSTER_ENDPOINT", "https://env-cluster")
			os.Setenv("ISOLATED_VPC", "true")
			os.Setenv("VM_MEMORY_OVERHEAD_PERCENT", "0.1")
			os.Setenv("INTERRUPTION_QUEUE", "env-cluster")
			os.Setenv("RESERVED_ENIS", "10")
			fs = &coreoptions.FlagSet{
				FlagSet: flag.NewFlagSet("karpenter", flag.ContinueOnError),
			}
			opts.AddFlags(fs)
			err := opts.Parse(fs)
			Expect(err).ToNot(HaveOccurred())
			expectOptionsEqual(opts, test.Options(test.OptionsFields{
				AssumeRoleARN:           lo.ToPtr("env-role"),
				AssumeRoleDuration:      lo.ToPtr(20 * time.Minute),
				ClusterCABundle:         lo.ToPtr("env-bundle"),
				ClusterName:             lo.ToPtr("env-cluster"),
				ClusterEndpoint:         lo.ToPtr("https://env-cluster"),
				IsolatedVPC:             lo.ToPtr(true),
				VMMemoryOverheadPercent: lo.ToPtr[float64](0.1),
				InterruptionQueue:       lo.ToPtr("env-cluster"),
				ReservedENIs:            lo.ToPtr(10),
			}))
>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13
		})
	})

	Context("Validation", func() {
<<<<<<< HEAD
		DescribeTable(
			"should parse valid log levels successfully",
			func(level string) {
				err := opts.Parse(fs, "--log-level", level)
				Expect(err).To(BeNil())
			},
			Entry("empty string", ""),
			Entry("debug", "debug"),
			Entry("info", "info"),
			Entry("error", "error"),
		)
		It("should error with an invalid log level", func() {
			err := opts.Parse(fs, "--log-level", "hello")
			Expect(err).ToNot(BeNil())
=======
		It("should fail when cluster name is not set", func() {
			err := opts.Parse(fs)
			// Overwrite ClusterName since it is commonly set by environment variables in dev environments
			opts.ClusterName = ""
			Expect(err).ToNot(HaveOccurred())
			Expect(func() {
				opts.MergeSettings(ctx)
				fmt.Printf("%#v", opts)
			}).To(Panic())
		})
		It("should fail when assume role duration is less than 15 minutes", func() {
			err := opts.Parse(fs, "--assume-role-duration", "1s")
			Expect(err).To(HaveOccurred())
		})
		It("should fail when clusterEndpoint is invalid (not absolute)", func() {
			err := opts.Parse(fs, "--cluster-endpoint", "00000000000000000000000.gr7.us-west-2.eks.amazonaws.com")
			Expect(err).To(HaveOccurred())
		})
		It("should fail when vmMemoryOverheadPercent is negative", func() {
			err := opts.Parse(fs, "--vm-memory-overhead-percent", "-0.01")
			Expect(err).To(HaveOccurred())
		})
		It("should fail when reservedENIs is negative", func() {
			err := opts.Parse(fs, "--reserved-enis", "-1")
			Expect(err).To(HaveOccurred())
>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13
		})
	})
})

<<<<<<< HEAD
func expectOptionsMatch(optsA, optsB *options.Options) {
	GinkgoHelper()
	if optsA == nil && optsB == nil {
		return
	}
	Expect(optsA).ToNot(BeNil())
	Expect(optsB).ToNot(BeNil())
	Expect(optsA.ServiceName).To(Equal(optsB.ServiceName))
	Expect(optsA.DisableWebhook).To(Equal(optsB.DisableWebhook))
	Expect(optsA.WebhookPort).To(Equal(optsB.WebhookPort))
	Expect(optsA.MetricsPort).To(Equal(optsB.MetricsPort))
	Expect(optsA.WebhookMetricsPort).To(Equal(optsB.WebhookMetricsPort))
	Expect(optsA.HealthProbePort).To(Equal(optsB.HealthProbePort))
	Expect(optsA.KubeClientQPS).To(Equal(optsB.KubeClientQPS))
	Expect(optsA.KubeClientBurst).To(Equal(optsB.KubeClientBurst))
	Expect(optsA.EnableProfiling).To(Equal(optsB.EnableProfiling))
	Expect(optsA.EnableLeaderElection).To(Equal(optsB.EnableLeaderElection))
	Expect(optsA.MemoryLimit).To(Equal(optsB.MemoryLimit))
	Expect(optsA.LogLevel).To(Equal(optsB.LogLevel))
	Expect(optsA.BatchMaxDuration).To(Equal(optsB.BatchMaxDuration))
	Expect(optsA.BatchIdleDuration).To(Equal(optsB.BatchIdleDuration))
	Expect(optsA.FeatureGates.Drift).To(Equal(optsB.FeatureGates.Drift))
=======
func expectOptionsEqual(optsA *options.Options, optsB *options.Options) {
	GinkgoHelper()
	Expect(optsA.AssumeRoleARN).To(Equal(optsB.AssumeRoleARN))
	Expect(optsA.AssumeRoleDuration).To(Equal(optsB.AssumeRoleDuration))
	Expect(optsA.ClusterCABundle).To(Equal(optsB.ClusterCABundle))
	Expect(optsA.ClusterName).To(Equal(optsB.ClusterName))
	Expect(optsA.ClusterEndpoint).To(Equal(optsB.ClusterEndpoint))
	Expect(optsA.IsolatedVPC).To(Equal(optsB.IsolatedVPC))
	Expect(optsA.VMMemoryOverheadPercent).To(Equal(optsB.VMMemoryOverheadPercent))
	Expect(optsA.InterruptionQueue).To(Equal(optsB.InterruptionQueue))
	Expect(optsA.ReservedENIs).To(Equal(optsB.ReservedENIs))
>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13
}
