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

package integration_test

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/samber/lo"

	coretest "sigs.k8s.io/karpenter/pkg/test"

	awserrors "github.com/aws/karpenter-provider-aws/pkg/errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*
	TODO GW 2 - this will all need to be updated (maybe just env.GetInstanceProfileName) to
		- expect/find/grab the instance profile with the given identifying metadata
		- check that latest isn't deleted and neither is old one when not older than res window
		- check that is does get deleted when appropriate
*/
var _ = Describe("InstanceProfile Generation", func() {
	It("should generate the InstanceProfile when setting the role", func() {
		pod := coretest.Pod()
		env.ExpectCreated(nodePool, nodeClass, pod)
		env.EventuallyExpectHealthy(pod)
		node := env.ExpectCreatedNodeCount("==", 1)[0]

		instance := env.GetInstance(node.Name)
		Expect(instance.IamInstanceProfile).ToNot(BeNil())
		Expect(lo.FromPtr(instance.IamInstanceProfile.Arn)).To(ContainSubstring(nodeClass.Status.InstanceProfile))

		instanceProfile := env.EventuallyExpectInstanceProfileExists(env.GetInstanceProfileName(nodeClass))
		Expect(instanceProfile.Roles).To(HaveLen(1))
		Expect(lo.FromPtr(instanceProfile.Roles[0].RoleName)).To(Equal(nodeClass.Spec.Role))
	})
	It("should remove the generated InstanceProfile when deleting the NodeClass", func() {
		pod := coretest.Pod()
		env.ExpectCreated(nodePool, nodeClass, pod)
		env.EventuallyExpectHealthy(pod)
		env.ExpectCreatedNodeCount("==", 1)

		env.ExpectDeleted(nodePool, nodeClass)
		Eventually(func(g Gomega) {
			_, err := env.IAMAPI.GetInstanceProfileWithContext(env.Context, &iam.GetInstanceProfileInput{
				InstanceProfileName: aws.String(env.GetInstanceProfileName(nodeClass)),
			})
			g.Expect(awserrors.IsNotFound(err)).To(BeTrue())
		}).Should(Succeed())
	})
	It("should use the unmanaged instance profile", func() {
		instanceProfileName := fmt.Sprintf("KarpenterNodeInstanceProfile-%s", env.ClusterName)
		roleName := fmt.Sprintf("KarpenterNodeRole-%s", env.ClusterName)
		env.ExpectInstanceProfileCreated(instanceProfileName, roleName)
		DeferCleanup(func() {
			env.ExpectInstanceProfileDeleted(instanceProfileName, roleName)
		})

		pod := coretest.Pod()
		nodeClass.Spec.Role = ""
		nodeClass.Spec.InstanceProfile = lo.ToPtr(fmt.Sprintf("KarpenterNodeInstanceProfile-%s", env.ClusterName))
		env.ExpectCreated(nodePool, nodeClass, pod)
		env.EventuallyExpectHealthy(pod)
		node := env.ExpectCreatedNodeCount("==", 1)[0]

		instance := env.GetInstance(node.Name)
		Expect(instance.IamInstanceProfile).ToNot(BeNil())
		Expect(lo.FromPtr(instance.IamInstanceProfile.Arn)).To(ContainSubstring(nodeClass.Status.InstanceProfile))
	})
})
