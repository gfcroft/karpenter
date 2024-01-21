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

package apis

import (
	_ "embed"

	"github.com/samber/lo"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"

<<<<<<< HEAD
<<<<<<< HEAD
	"github.com/aws/karpenter-core/pkg/apis/settings"
	"github.com/aws/karpenter-core/pkg/apis/v1beta1"
=======
	"github.com/aws/karpenter/pkg/apis/settings"
	"github.com/aws/karpenter/pkg/apis/v1beta1"

	"github.com/samber/lo"

	"github.com/aws/karpenter-core/pkg/apis"
	coresettings "github.com/aws/karpenter-core/pkg/apis/settings"
>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13
	"github.com/aws/karpenter-core/pkg/utils/functional"
=======
	"github.com/aws/karpenter-provider-aws/pkg/apis/v1beta1"

	"github.com/samber/lo"

	"sigs.k8s.io/karpenter/pkg/apis"
	"sigs.k8s.io/karpenter/pkg/utils/functional"
>>>>>>> 6ebba50ce424ccd5e33b3c84b4f10a8e78d54539
)

var (
	// Builder includes all types within the apis package
	Builder = runtime.NewSchemeBuilder(
		v1beta1.SchemeBuilder.AddToScheme,
	)
	// AddToScheme may be used to add all resources defined in the project to a Scheme
	AddToScheme = Builder.AddToScheme
<<<<<<< HEAD
	Settings    = []settings.Injectable{&settings.Settings{}}
=======
>>>>>>> 6ebba50ce424ccd5e33b3c84b4f10a8e78d54539
)

//go:generate controller-gen crd object:headerFile="../../hack/boilerplate.go.txt" paths="./..." output:crd:artifacts:config=crds
var (
<<<<<<< HEAD
	//go:embed crds/karpenter.sh_nodepools.yaml
	NodePoolCRD []byte
	//go:embed crds/karpenter.sh_nodeclaims.yaml
	NodeClaimCRD []byte
	CRDs         = []*v1.CustomResourceDefinition{
		lo.Must(functional.Unmarshal[v1.CustomResourceDefinition](NodePoolCRD)),
		lo.Must(functional.Unmarshal[v1.CustomResourceDefinition](NodeClaimCRD)),
	}
=======
	//go:embed crds/karpenter.k8s.aws_ec2nodeclasses.yaml
	EC2NodeClassCRD []byte
	CRDs            = append(apis.CRDs,
		lo.Must(functional.Unmarshal[v1.CustomResourceDefinition](EC2NodeClassCRD)),
	)
>>>>>>> 6ebba50ce424ccd5e33b3c84b4f10a8e78d54539
)
