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

package instanceprofile

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/mitchellh/hashstructure/v2"
	"github.com/patrickmn/go-cache"
	"github.com/samber/lo"
	v1 "k8s.io/api/core/v1"

	corev1beta1 "sigs.k8s.io/karpenter/pkg/apis/v1beta1"

	"github.com/aws/karpenter-provider-aws/pkg/apis/v1beta1"
	awserrors "github.com/aws/karpenter-provider-aws/pkg/errors"
	"github.com/aws/karpenter-provider-aws/pkg/operator/options"
)

type Provider struct {
	region string
	iamapi iamiface.IAMAPI
	ec2api ec2iface.EC2API
	cache  *cache.Cache
}

func NewProvider(region string, iamapi iamiface.IAMAPI, ec2api ec2iface.EC2API, cache *cache.Cache) *Provider {
	return &Provider{
		region: region,
		iamapi: iamapi,
		ec2api: ec2api,
		cache:  cache,
	}
}

func (p *Provider) Create(ctx context.Context, nodeClass *v1beta1.EC2NodeClass) (string, error) {
	tags := lo.Assign(nodeClass.Spec.Tags, map[string]string{
		fmt.Sprintf("kubernetes.io/cluster/%s", options.FromContext(ctx).ClusterName): "owned",
		corev1beta1.ManagedByAnnotationKey:                                            options.FromContext(ctx).ClusterName,
		v1beta1.LabelNodeClass:                                                        nodeClass.Name,
		v1.LabelTopologyRegion:                                                        p.region,
	})
	// An instance profile exists for this NodeClass with correct role
	if cachedProfile, ok := p.cache.Get(string(nodeClass.UID)); ok {
		profile := cachedProfile.(*iam.InstanceProfile)
		if ProfileHasMatchingRole(profile, nodeClass.Spec.Role) {
			return aws.StringValue(profile.InstanceProfileName), nil
		}
	}
	profilePath := GetProfilePath(ctx, nodeClass)
	profileName := GetProfileName(ctx, p.region, nodeClass, time.Now())
	var instanceProfile *iam.InstanceProfile
	// Validate if the instance profile exists and has the correct role assigned to it
	retrievedProfiles, err := p.iamapi.ListInstanceProfilesWithContext(ctx, &iam.ListInstanceProfilesInput{PathPrefix: aws.String(profilePath)})
	if err != nil {
		if !awserrors.IsNotFound(err) {
			return "", fmt.Errorf("listing instance profiles for nodeClass %s, %q", nodeClass.Name, err)
		}
		o, err := p.iamapi.CreateInstanceProfileWithContext(ctx, &iam.CreateInstanceProfileInput{
			InstanceProfileName: aws.String(profileName),
			Tags:                lo.MapToSlice(tags, func(k, v string) *iam.Tag { return &iam.Tag{Key: aws.String(k), Value: aws.String(v)} }),
			Path:                aws.String(fmt.Sprintf("%s", profileName)),
		})
		if err != nil {
			return "", fmt.Errorf("creating instance profile %q, %w", profileName, err)
		}
		instanceProfile = o.InstanceProfile
	} else {
		for _, retrievedProfile := range retrievedProfiles.InstanceProfiles {
			if ProfileHasMatchingRole(retrievedProfile, nodeClass.Spec.Role) {
				return aws.StringValue(retrievedProfile.InstanceProfileName), nil
			}
		}
		o, err := p.iamapi.CreateInstanceProfileWithContext(ctx, &iam.CreateInstanceProfileInput{
			InstanceProfileName: aws.String(profileName),
			Tags:                lo.MapToSlice(tags, func(k, v string) *iam.Tag { return &iam.Tag{Key: aws.String(k), Value: aws.String(v)} }),
			Path:                aws.String(fmt.Sprintf("%s", GetProfilePath(ctx, nodeClass))),
		})
		if err != nil {
			return "", fmt.Errorf("creating instance profile %q, %w", profileName, err)
		}
		instanceProfile = o.InstanceProfile
	}
	if _, err = p.iamapi.AddRoleToInstanceProfileWithContext(ctx, &iam.AddRoleToInstanceProfileInput{
		InstanceProfileName: aws.String(profileName),
		RoleName:            aws.String(nodeClass.Spec.Role),
	}); err != nil {
		return "", fmt.Errorf("adding role %q to instance profile %q, %w", nodeClass.Spec.Role, profileName, err)
	}
	p.cache.SetDefault(string(nodeClass.UID), instanceProfile)
	return profileName, nil
}

func (p *Provider) Delete(ctx context.Context, nodeClass *v1beta1.EC2NodeClass) error {
	profiles, err := p.iamapi.ListInstanceProfilesWithContext(ctx, &iam.ListInstanceProfilesInput{
		PathPrefix: aws.String(GetProfilePath(ctx, nodeClass)),
	})
	if err != nil {
		return awserrors.IgnoreNotFound(fmt.Errorf("getting instance profiles for nodeClass: %q, %w", nodeClass.Name, err))
	}
	for _, profile := range profiles.InstanceProfiles {
		// instance profiles should not be deleted when

		// Instance profiles can only have a single role assigned to them so this profile either has 1 or 0 roles
		// https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_use_switch-role-ec2_instance-profiles.html
		if len(profile.Roles) == 1 {
			if _, err = p.iamapi.RemoveRoleFromInstanceProfileWithContext(ctx, &iam.RemoveRoleFromInstanceProfileInput{
				InstanceProfileName: profile.InstanceProfileName,
				RoleName:            profile.Roles[0].RoleName,
			}); err != nil {
				return fmt.Errorf("removing role %q from instance profile %q, %w", profile.Roles[0].RoleName, profile.InstanceProfileName, err)
			}
		}
		//verify not in use using ec2 describe instances api
		//ec2 + describe-instances + instanceprofile
		instances, err := p.ec2api.DescribeInstancesWithContext(ctx, &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("iam-instance-profile.id"),
					Values: []*string{profile.InstanceProfileId},
				},
			},
		})
		if err == nil && instances.Reservations > 0 {
		}
		if _, err = p.iamapi.DeleteInstanceProfileWithContext(ctx, &iam.DeleteInstanceProfileInput{
			InstanceProfileName: profile.InstanceProfileName,
		}); err != nil {
			return awserrors.IgnoreNotFound(fmt.Errorf("deleting instance profile %q, %w", profile.InstanceProfileName, err))
		}
	}
	return nil
}

// GetProfileName gets the string for the profile name based on the cluster name and the NodeClass UUID and the provided time.
// The length of this string can never exceed the maximum instance profile name limit of 128 characters.
func GetProfileName(ctx context.Context, region string, nodeClass *v1beta1.EC2NodeClass, currentTime time.Time) string {
	return fmt.Sprintf("%s_%d", options.FromContext(ctx).ClusterName, lo.Must(hashstructure.Hash(fmt.Sprintf("%s%s%s", region, nodeClass.Name, currentTime.String()), hashstructure.FormatV2, nil)))
}

// GetProfilePath returns the instance profile path for a given nodeClass existing on a cluster
func GetProfilePath(ctx context.Context, nodeClass *v1beta1.EC2NodeClass) string {
	return fmt.Sprintf("/%s/%s/", options.FromContext(ctx).ClusterName, nodeClass.Name)
}

// ProfileHasMatchingRole checks if a given instance profile contains a matching role to the rolename provided
func ProfileHasMatchingRole(profile *iam.InstanceProfile, roleName string) bool {
	// Instance profiles can only have a single role assigned to them so this profile either has 1 or 0 roles
	// https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_use_switch-role-ec2_instance-profiles.html
	if len(profile.Roles) == 1 {
		return aws.StringValue(profile.Roles[0].RoleName) == roleName
	}
	return false
}
