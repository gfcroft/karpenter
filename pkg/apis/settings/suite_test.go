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

package settings_test

import (
	"context"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	. "knative.dev/pkg/logging/testing"

	"github.com/aws/karpenter-core/pkg/apis/settings"
)

var ctx context.Context

func TestSettings(t *testing.T) {
	ctx = TestContextWithLogger(t)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Settings")
}

var _ = Describe("Validation", func() {
	It("should succeed to set defaults", func() {
		cm := &v1.ConfigMap{
			Data: map[string]string{},
		}
		ctx, err := (&settings.Settings{}).Inject(ctx, cm)
		Expect(err).ToNot(HaveOccurred())
		s := settings.FromContext(ctx)
		Expect(s.BatchMaxDuration).To(Equal(time.Second * 10))
		Expect(s.BatchIdleDuration).To(Equal(time.Second))
		Expect(s.DriftEnabled).To(BeFalse())
	})
	It("should succeed to set custom values", func() {
		cm := &v1.ConfigMap{
			Data: map[string]string{
				"batchMaxDuration":          "30s",
				"batchIdleDuration":         "5s",
				"featureGates.driftEnabled": "true",
			},
		}
		ctx, err := (&settings.Settings{}).Inject(ctx, cm)
		Expect(err).ToNot(HaveOccurred())
		s := settings.FromContext(ctx)
		Expect(s.BatchMaxDuration).To(Equal(time.Second * 30))
		Expect(s.BatchIdleDuration).To(Equal(time.Second * 5))
		Expect(s.DriftEnabled).To(BeTrue())
	})
	It("should fail validation when batchMaxDuration is negative", func() {
		cm := &v1.ConfigMap{
			Data: map[string]string{
<<<<<<< HEAD
				"batchMaxDuration": "-10s",
			},
		}
		_, err := (&settings.Settings{}).Inject(ctx, cm)
		Expect(err).To(HaveOccurred())
	})
	It("should fail validation when batchMaxDuration is less then 1s", func() {
=======
				"aws.clusterEndpoint":            "https://00000000000000000000000.gr7.us-west-2.eks.amazonaws.com",
				"aws.clusterName":                "my-cluster",
				"aws.defaultInstanceProfile":     "karpenter",
				"aws.enablePodENI":               "true",
				"aws.enableENILimitedPodDensity": "false",
				"aws.isolatedVPC":                "true",
				"aws.vmMemoryOverheadPercent":    "0.1",
				"aws.tags":                       `{"tag1": "value1", "tag2": "value2", "example.com/tag": "my-value"}`,
				"aws.reservedENIs":               "1",
				"aws.nodeNameConvention":         "resource-name",
			},
		}
		ctx, err := (&settings.Settings{}).Inject(ctx, cm)
		Expect(err).ToNot(HaveOccurred())
		s := settings.FromContext(ctx)
		Expect(s.DefaultInstanceProfile).To(Equal("karpenter"))
		Expect(s.EnablePodENI).To(BeTrue())
		Expect(s.EnableENILimitedPodDensity).To(BeFalse())
		Expect(s.IsolatedVPC).To(BeTrue())
		Expect(s.VMMemoryOverheadPercent).To(Equal(0.1))
		Expect(len(s.Tags)).To(Equal(3))
		Expect(s.Tags).To(HaveKeyWithValue("tag1", "value1"))
		Expect(s.Tags).To(HaveKeyWithValue("tag2", "value2"))
		Expect(s.Tags).To(HaveKeyWithValue("example.com/tag", "my-value"))
		Expect(s.ReservedENIs).To(Equal(1))
	})
	It("should succeed validation when tags contain parts of restricted domains", func() {
		cm := &v1.ConfigMap{
			Data: map[string]string{
				"aws.clusterEndpoint": "https://00000000000000000000000.gr7.us-west-2.eks.amazonaws.com",
				"aws.clusterName":     "my-cluster",
				"aws.tags":            `{"karpenter.sh/custom-key": "value1", "karpenter.sh/managed": "true", "kubernetes.io/role/key": "value2", "kubernetes.io/cluster/other-tag/hello": "value3"}`,
			},
		}
		ctx, err := (&settings.Settings{}).Inject(ctx, cm)
		Expect(err).ToNot(HaveOccurred())
		s := settings.FromContext(ctx)
		Expect(s.Tags).To(HaveKeyWithValue("karpenter.sh/custom-key", "value1"))
		Expect(s.Tags).To(HaveKeyWithValue("karpenter.sh/managed", "true"))
		Expect(s.Tags).To(HaveKeyWithValue("kubernetes.io/role/key", "value2"))
		Expect(s.Tags).To(HaveKeyWithValue("kubernetes.io/cluster/other-tag/hello", "value3"))
	})
	It("should fail validation when clusterEndpoint is invalid (not absolute)", func() {
>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13
		cm := &v1.ConfigMap{
			Data: map[string]string{
				"batchMaxDuration": "800ms",
			},
		}
		_, err := (&settings.Settings{}).Inject(ctx, cm)
		Expect(err).To(HaveOccurred())
	})
	It("should fail validation when batchMaxDuration is set to empty", func() {
		cm := &v1.ConfigMap{
			Data: map[string]string{
				"batchMaxDuration": "",
			},
		}
		_, err := (&settings.Settings{}).Inject(ctx, cm)
		Expect(err).To(HaveOccurred())
	})
	It("should fail validation when batchIdleDuration is negative", func() {
		cm := &v1.ConfigMap{
			Data: map[string]string{
				"batchIdleDuration": "-1s",
			},
		}
		_, err := (&settings.Settings{}).Inject(ctx, cm)
		Expect(err).To(HaveOccurred())
	})
	It("should fail validation when batchIdleDuration is less then 1s", func() {
		cm := &v1.ConfigMap{
			Data: map[string]string{
				"batchIdleDuration": "800ms",
			},
		}
		_, err := (&settings.Settings{}).Inject(ctx, cm)
		Expect(err).To(HaveOccurred())
	})
	It("should fail validation when batchIdleDuration is set to empty", func() {
		cm := &v1.ConfigMap{
			Data: map[string]string{
				"batchMaxDuration": "",
			},
		}
		_, err := (&settings.Settings{}).Inject(ctx, cm)
		Expect(err).To(HaveOccurred())
	})
	It("should fail validation when driftEnabled is not a valid boolean value", func() {
		cm := &v1.ConfigMap{
			Data: map[string]string{
				"featureGates.driftEnabled": "foobar",
			},
		}
		_, err := (&settings.Settings{}).Inject(ctx, cm)
		Expect(err).To(HaveOccurred())
	})
})
