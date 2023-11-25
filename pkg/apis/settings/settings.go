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

package settings

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/multierr"
	v1 "k8s.io/api/core/v1"
	"knative.dev/pkg/configmap"

	coresettings "github.com/aws/karpenter-core/pkg/apis/settings"
)

type settingsKeyType struct{}

var ContextKey = settingsKeyType{}

var defaultSettings = &Settings{
	BatchMaxDuration:  time.Second * 10,
	BatchIdleDuration: time.Second,
	DriftEnabled:      false,
}

// +k8s:deepcopy-gen=true
type Settings struct {
	BatchMaxDuration  time.Duration
	BatchIdleDuration time.Duration
	// This feature flag is temporary and will be removed in the near future.
	DriftEnabled bool
}

func (*Settings) ConfigMap() string {
	return "karpenter-global-settings"
}

// Inject creates a Settings from the supplied ConfigMap
func (*Settings) Inject(ctx context.Context, cm *v1.ConfigMap) (context.Context, error) {
	s := defaultSettings.DeepCopy()
	if cm == nil {
		return ToContext(ctx, s), nil
	}

	if err := configmap.Parse(cm.Data,
		configmap.AsDuration("batchMaxDuration", &s.BatchMaxDuration),
		configmap.AsDuration("batchIdleDuration", &s.BatchIdleDuration),
		configmap.AsBool("featureGates.driftEnabled", &s.DriftEnabled),
	); err != nil {
		return ctx, fmt.Errorf("parsing settings, %w", err)
	}
	if err := s.Validate(); err != nil {
		return ctx, fmt.Errorf("validating settings, %w", err)
	}
	return ToContext(ctx, s), nil
}

<<<<<<< HEAD
func (in *Settings) Validate() (err error) {
	if in.BatchMaxDuration < time.Second {
		err = multierr.Append(err, fmt.Errorf("batchMaxDuration cannot be less then 1s"))
	}
	if in.BatchIdleDuration < time.Second {
		err = multierr.Append(err, fmt.Errorf("batchIdleDuration cannot be less then 1s"))
	}
	return err
=======
func (*Settings) FromContext(ctx context.Context) coresettings.Injectable {
	return FromContext(ctx)
>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13
}

func ToContext(ctx context.Context, s *Settings) context.Context {
	return context.WithValue(ctx, ContextKey, s)
}

func FromContext(ctx context.Context) *Settings {
	data := ctx.Value(ContextKey)
	if data == nil {
		panic("settings not in context")
	}
	return data.(*Settings)
}
<<<<<<< HEAD
=======

// AsTypedString passes the value at key through into the target, if it exists.
func AsTypedString[T ~string](key string, target *T) configmap.ParseFunc {
	return func(data map[string]string) error {
		if raw, ok := data[key]; ok {
			*target = T(raw)
		}
		return nil
	}
}

// AsStringMap parses a value as a JSON map of map[string]string.
func AsStringMap(key string, target *map[string]string) configmap.ParseFunc {
	return func(data map[string]string) error {
		if raw, ok := data[key]; ok && raw != "" {
			m := map[string]string{}
			if err := json.Unmarshal([]byte(raw), &m); err != nil {
				return err
			}
			*target = m
		}
		return nil
	}
}
>>>>>>> 1db74f402628818c1f6ead391cc039d2834e7e13
