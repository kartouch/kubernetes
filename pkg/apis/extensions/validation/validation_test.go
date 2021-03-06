/*
Copyright 2014 The Kubernetes Authors All rights reserved.

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

package validation

import (
	"fmt"
	"strings"
	"testing"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/util/intstr"
)

func TestValidateHorizontalPodAutoscaler(t *testing.T) {
	successCases := []extensions.HorizontalPodAutoscaler{
		{
			ObjectMeta: api.ObjectMeta{
				Name:      "myautoscaler",
				Namespace: api.NamespaceDefault,
			},
			Spec: extensions.HorizontalPodAutoscalerSpec{
				ScaleRef: extensions.SubresourceReference{
					Kind:        "ReplicationController",
					Name:        "myrc",
					Subresource: "scale",
				},
				MinReplicas:    newInt(1),
				MaxReplicas:    5,
				CPUUtilization: &extensions.CPUTargetUtilization{TargetPercentage: 70},
			},
		},
		{
			ObjectMeta: api.ObjectMeta{
				Name:      "myautoscaler",
				Namespace: api.NamespaceDefault,
			},
			Spec: extensions.HorizontalPodAutoscalerSpec{
				ScaleRef: extensions.SubresourceReference{
					Kind:        "ReplicationController",
					Name:        "myrc",
					Subresource: "scale",
				},
				MinReplicas: newInt(1),
				MaxReplicas: 5,
			},
		},
	}
	for _, successCase := range successCases {
		if errs := ValidateHorizontalPodAutoscaler(&successCase); len(errs) != 0 {
			t.Errorf("expected success: %v", errs)
		}
	}

	errorCases := []struct {
		horizontalPodAutoscaler extensions.HorizontalPodAutoscaler
		msg                     string
	}{
		{
			horizontalPodAutoscaler: extensions.HorizontalPodAutoscaler{
				ObjectMeta: api.ObjectMeta{Name: "myautoscaler", Namespace: api.NamespaceDefault},
				Spec: extensions.HorizontalPodAutoscalerSpec{
					ScaleRef:       extensions.SubresourceReference{Name: "myrc", Subresource: "scale"},
					MinReplicas:    newInt(1),
					MaxReplicas:    5,
					CPUUtilization: &extensions.CPUTargetUtilization{TargetPercentage: 70},
				},
			},
			msg: "scaleRef.kind: Required",
		},
		{
			horizontalPodAutoscaler: extensions.HorizontalPodAutoscaler{
				ObjectMeta: api.ObjectMeta{Name: "myautoscaler", Namespace: api.NamespaceDefault},
				Spec: extensions.HorizontalPodAutoscalerSpec{
					ScaleRef:       extensions.SubresourceReference{Kind: "..", Name: "myrc", Subresource: "scale"},
					MinReplicas:    newInt(1),
					MaxReplicas:    5,
					CPUUtilization: &extensions.CPUTargetUtilization{TargetPercentage: 70},
				},
			},
			msg: "scaleRef.kind: Invalid",
		},
		{
			horizontalPodAutoscaler: extensions.HorizontalPodAutoscaler{
				ObjectMeta: api.ObjectMeta{Name: "myautoscaler", Namespace: api.NamespaceDefault},
				Spec: extensions.HorizontalPodAutoscalerSpec{
					ScaleRef:       extensions.SubresourceReference{Kind: "ReplicationController", Subresource: "scale"},
					MinReplicas:    newInt(1),
					MaxReplicas:    5,
					CPUUtilization: &extensions.CPUTargetUtilization{TargetPercentage: 70},
				},
			},
			msg: "scaleRef.name: Required",
		},
		{
			horizontalPodAutoscaler: extensions.HorizontalPodAutoscaler{
				ObjectMeta: api.ObjectMeta{Name: "myautoscaler", Namespace: api.NamespaceDefault},
				Spec: extensions.HorizontalPodAutoscalerSpec{
					ScaleRef:       extensions.SubresourceReference{Kind: "ReplicationController", Name: "..", Subresource: "scale"},
					MinReplicas:    newInt(1),
					MaxReplicas:    5,
					CPUUtilization: &extensions.CPUTargetUtilization{TargetPercentage: 70},
				},
			},
			msg: "scaleRef.name: Invalid",
		},
		{
			horizontalPodAutoscaler: extensions.HorizontalPodAutoscaler{
				ObjectMeta: api.ObjectMeta{Name: "myautoscaler", Namespace: api.NamespaceDefault},
				Spec: extensions.HorizontalPodAutoscalerSpec{
					ScaleRef:       extensions.SubresourceReference{Kind: "ReplicationController", Name: "myrc", Subresource: ""},
					MinReplicas:    newInt(1),
					MaxReplicas:    5,
					CPUUtilization: &extensions.CPUTargetUtilization{TargetPercentage: 70},
				},
			},
			msg: "scaleRef.subresource: Required",
		},
		{
			horizontalPodAutoscaler: extensions.HorizontalPodAutoscaler{
				ObjectMeta: api.ObjectMeta{Name: "myautoscaler", Namespace: api.NamespaceDefault},
				Spec: extensions.HorizontalPodAutoscalerSpec{
					ScaleRef:       extensions.SubresourceReference{Kind: "ReplicationController", Name: "myrc", Subresource: ".."},
					MinReplicas:    newInt(1),
					MaxReplicas:    5,
					CPUUtilization: &extensions.CPUTargetUtilization{TargetPercentage: 70},
				},
			},
			msg: "scaleRef.subresource: Invalid",
		},
		{
			horizontalPodAutoscaler: extensions.HorizontalPodAutoscaler{
				ObjectMeta: api.ObjectMeta{Name: "myautoscaler", Namespace: api.NamespaceDefault},
				Spec: extensions.HorizontalPodAutoscalerSpec{
					ScaleRef:       extensions.SubresourceReference{Kind: "ReplicationController", Name: "myrc", Subresource: "randomsubresource"},
					MinReplicas:    newInt(1),
					MaxReplicas:    5,
					CPUUtilization: &extensions.CPUTargetUtilization{TargetPercentage: 70},
				},
			},
			msg: "scaleRef.subresource: Unsupported",
		},
		{
			horizontalPodAutoscaler: extensions.HorizontalPodAutoscaler{
				ObjectMeta: api.ObjectMeta{
					Name:      "myautoscaler",
					Namespace: api.NamespaceDefault,
				},
				Spec: extensions.HorizontalPodAutoscalerSpec{
					ScaleRef: extensions.SubresourceReference{
						Subresource: "scale",
					},
					MinReplicas: newInt(-1),
					MaxReplicas: 5,
				},
			},
			msg: "must be greater than 0",
		},
		{
			horizontalPodAutoscaler: extensions.HorizontalPodAutoscaler{
				ObjectMeta: api.ObjectMeta{
					Name:      "myautoscaler",
					Namespace: api.NamespaceDefault,
				},
				Spec: extensions.HorizontalPodAutoscalerSpec{
					ScaleRef: extensions.SubresourceReference{
						Subresource: "scale",
					},
					MinReplicas: newInt(7),
					MaxReplicas: 5,
				},
			},
			msg: "must be greater than or equal to `minReplicas`",
		},
		{
			horizontalPodAutoscaler: extensions.HorizontalPodAutoscaler{
				ObjectMeta: api.ObjectMeta{
					Name:      "myautoscaler",
					Namespace: api.NamespaceDefault,
				},
				Spec: extensions.HorizontalPodAutoscalerSpec{
					ScaleRef: extensions.SubresourceReference{
						Subresource: "scale",
					},
					MinReplicas:    newInt(1),
					MaxReplicas:    5,
					CPUUtilization: &extensions.CPUTargetUtilization{TargetPercentage: -70},
				},
			},
			msg: "must be greater than 0",
		},
	}

	for _, c := range errorCases {
		errs := ValidateHorizontalPodAutoscaler(&c.horizontalPodAutoscaler)
		if len(errs) == 0 {
			t.Errorf("expected failure for %q", c.msg)
		} else if !strings.Contains(errs[0].Error(), c.msg) {
			t.Errorf("unexpected error: %q, expected: %q", errs[0], c.msg)
		}
	}
}

func TestValidateDaemonSetStatusUpdate(t *testing.T) {
	type dsUpdateTest struct {
		old    extensions.DaemonSet
		update extensions.DaemonSet
	}

	successCases := []dsUpdateTest{
		{
			old: extensions.DaemonSet{
				ObjectMeta: api.ObjectMeta{Name: "abc", Namespace: api.NamespaceDefault},
				Status: extensions.DaemonSetStatus{
					CurrentNumberScheduled: 1,
					NumberMisscheduled:     2,
					DesiredNumberScheduled: 3,
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: api.ObjectMeta{Name: "abc", Namespace: api.NamespaceDefault},
				Status: extensions.DaemonSetStatus{
					CurrentNumberScheduled: 1,
					NumberMisscheduled:     1,
					DesiredNumberScheduled: 3,
				},
			},
		},
	}

	for _, successCase := range successCases {
		successCase.old.ObjectMeta.ResourceVersion = "1"
		successCase.update.ObjectMeta.ResourceVersion = "1"
		if errs := ValidateDaemonSetStatusUpdate(&successCase.update, &successCase.old); len(errs) != 0 {
			t.Errorf("expected success: %v", errs)
		}
	}

	errorCases := map[string]dsUpdateTest{
		"negative values": {
			old: extensions.DaemonSet{
				ObjectMeta: api.ObjectMeta{Name: "abc", Namespace: api.NamespaceDefault},
				Status: extensions.DaemonSetStatus{
					CurrentNumberScheduled: 1,
					NumberMisscheduled:     2,
					DesiredNumberScheduled: 3,
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: api.ObjectMeta{Name: "abc", Namespace: api.NamespaceDefault},
				Status: extensions.DaemonSetStatus{
					CurrentNumberScheduled: -1,
					NumberMisscheduled:     -1,
					DesiredNumberScheduled: -3,
				},
			},
		},
	}

	for testName, errorCase := range errorCases {
		if errs := ValidateDaemonSetStatusUpdate(&errorCase.old, &errorCase.update); len(errs) == 0 {
			t.Errorf("expected failure: %s", testName)
		}
	}
}

func TestValidateDaemonSetUpdate(t *testing.T) {
	validSelector := map[string]string{"a": "b"}
	validSelector2 := map[string]string{"c": "d"}
	invalidSelector := map[string]string{"NoUppercaseOrSpecialCharsLike=Equals": "b"}

	validUpdateStrategy := extensions.DaemonSetUpdateStrategy{
		Type: extensions.RollingUpdateDaemonSetStrategyType,
		RollingUpdate: &extensions.RollingUpdateDaemonSet{
			MaxUnavailable: intstr.FromInt(1),
		},
	}

	validPodSpecAbc := api.PodSpec{
		RestartPolicy: api.RestartPolicyAlways,
		DNSPolicy:     api.DNSClusterFirst,
		Containers:    []api.Container{{Name: "abc", Image: "image", ImagePullPolicy: "IfNotPresent"}},
	}
	validPodSpecDef := api.PodSpec{
		RestartPolicy: api.RestartPolicyAlways,
		DNSPolicy:     api.DNSClusterFirst,
		Containers:    []api.Container{{Name: "def", Image: "image", ImagePullPolicy: "IfNotPresent"}},
	}
	validPodSpecNodeSelector := api.PodSpec{
		NodeSelector:  validSelector,
		NodeName:      "xyz",
		RestartPolicy: api.RestartPolicyAlways,
		DNSPolicy:     api.DNSClusterFirst,
		Containers:    []api.Container{{Name: "abc", Image: "image", ImagePullPolicy: "IfNotPresent"}},
	}
	validPodSpecVolume := api.PodSpec{
		Volumes:       []api.Volume{{Name: "gcepd", VolumeSource: api.VolumeSource{GCEPersistentDisk: &api.GCEPersistentDiskVolumeSource{PDName: "my-PD", FSType: "ext4", Partition: 1, ReadOnly: false}}}},
		RestartPolicy: api.RestartPolicyAlways,
		DNSPolicy:     api.DNSClusterFirst,
		Containers:    []api.Container{{Name: "abc", Image: "image", ImagePullPolicy: "IfNotPresent"}},
	}

	validPodTemplateAbc := api.PodTemplate{
		Template: api.PodTemplateSpec{
			ObjectMeta: api.ObjectMeta{
				Labels: validSelector,
			},
			Spec: validPodSpecAbc,
		},
	}
	validPodTemplateNodeSelector := api.PodTemplate{
		Template: api.PodTemplateSpec{
			ObjectMeta: api.ObjectMeta{
				Labels: validSelector,
			},
			Spec: validPodSpecNodeSelector,
		},
	}
	validPodTemplateAbc2 := api.PodTemplate{
		Template: api.PodTemplateSpec{
			ObjectMeta: api.ObjectMeta{
				Labels: validSelector2,
			},
			Spec: validPodSpecAbc,
		},
	}
	validPodTemplateDef := api.PodTemplate{
		Template: api.PodTemplateSpec{
			ObjectMeta: api.ObjectMeta{
				Labels: validSelector2,
			},
			Spec: validPodSpecDef,
		},
	}
	invalidPodTemplate := api.PodTemplate{
		Template: api.PodTemplateSpec{
			Spec: api.PodSpec{
				RestartPolicy: api.RestartPolicyAlways,
				DNSPolicy:     api.DNSClusterFirst,
			},
			ObjectMeta: api.ObjectMeta{
				Labels: invalidSelector,
			},
		},
	}
	readWriteVolumePodTemplate := api.PodTemplate{
		Template: api.PodTemplateSpec{
			ObjectMeta: api.ObjectMeta{
				Labels: validSelector,
			},
			Spec: validPodSpecVolume,
		},
	}

	type dsUpdateTest struct {
		old    extensions.DaemonSet
		update extensions.DaemonSet
	}
	successCases := []dsUpdateTest{
		{
			old: extensions.DaemonSet{
				ObjectMeta: api.ObjectMeta{Name: "abc", Namespace: api.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:       &extensions.LabelSelector{MatchLabels: validSelector},
					Template:       &validPodTemplateAbc.Template,
					UpdateStrategy: validUpdateStrategy,
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: api.ObjectMeta{Name: "abc", Namespace: api.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:       &extensions.LabelSelector{MatchLabels: validSelector},
					Template:       &validPodTemplateAbc.Template,
					UpdateStrategy: validUpdateStrategy,
				},
			},
		},
		{
			old: extensions.DaemonSet{
				ObjectMeta: api.ObjectMeta{Name: "abc", Namespace: api.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:       &extensions.LabelSelector{MatchLabels: validSelector},
					Template:       &validPodTemplateAbc.Template,
					UpdateStrategy: validUpdateStrategy,
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: api.ObjectMeta{Name: "abc", Namespace: api.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:       &extensions.LabelSelector{MatchLabels: validSelector2},
					Template:       &validPodTemplateAbc2.Template,
					UpdateStrategy: validUpdateStrategy,
				},
			},
		},
		{
			old: extensions.DaemonSet{
				ObjectMeta: api.ObjectMeta{Name: "abc", Namespace: api.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:       &extensions.LabelSelector{MatchLabels: validSelector},
					Template:       &validPodTemplateAbc.Template,
					UpdateStrategy: validUpdateStrategy,
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: api.ObjectMeta{Name: "abc", Namespace: api.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:       &extensions.LabelSelector{MatchLabels: validSelector},
					Template:       &validPodTemplateNodeSelector.Template,
					UpdateStrategy: validUpdateStrategy,
				},
			},
		},
	}
	for _, successCase := range successCases {
		successCase.old.ObjectMeta.ResourceVersion = "1"
		successCase.update.ObjectMeta.ResourceVersion = "1"
		if errs := ValidateDaemonSetUpdate(&successCase.update, &successCase.old); len(errs) != 0 {
			t.Errorf("expected success: %v", errs)
		}
	}
	errorCases := map[string]dsUpdateTest{
		"change daemon name": {
			old: extensions.DaemonSet{
				ObjectMeta: api.ObjectMeta{Name: "", Namespace: api.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:       &extensions.LabelSelector{MatchLabels: validSelector},
					Template:       &validPodTemplateAbc.Template,
					UpdateStrategy: validUpdateStrategy,
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: api.ObjectMeta{Name: "abc", Namespace: api.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:       &extensions.LabelSelector{MatchLabels: validSelector},
					Template:       &validPodTemplateAbc.Template,
					UpdateStrategy: validUpdateStrategy,
				},
			},
		},
		"invalid selector": {
			old: extensions.DaemonSet{
				ObjectMeta: api.ObjectMeta{Name: "", Namespace: api.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:       &extensions.LabelSelector{MatchLabels: validSelector},
					Template:       &validPodTemplateAbc.Template,
					UpdateStrategy: validUpdateStrategy,
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: api.ObjectMeta{Name: "abc", Namespace: api.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:       &extensions.LabelSelector{MatchLabels: invalidSelector},
					Template:       &validPodTemplateAbc.Template,
					UpdateStrategy: validUpdateStrategy,
				},
			},
		},
		"invalid pod": {
			old: extensions.DaemonSet{
				ObjectMeta: api.ObjectMeta{Name: "", Namespace: api.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:       &extensions.LabelSelector{MatchLabels: validSelector},
					Template:       &validPodTemplateAbc.Template,
					UpdateStrategy: validUpdateStrategy,
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: api.ObjectMeta{Name: "abc", Namespace: api.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:       &extensions.LabelSelector{MatchLabels: validSelector},
					Template:       &invalidPodTemplate.Template,
					UpdateStrategy: validUpdateStrategy,
				},
			},
		},
		"change container image": {
			old: extensions.DaemonSet{
				ObjectMeta: api.ObjectMeta{Name: "", Namespace: api.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:       &extensions.LabelSelector{MatchLabels: validSelector},
					Template:       &validPodTemplateAbc.Template,
					UpdateStrategy: validUpdateStrategy,
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: api.ObjectMeta{Name: "abc", Namespace: api.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:       &extensions.LabelSelector{MatchLabels: validSelector},
					Template:       &validPodTemplateDef.Template,
					UpdateStrategy: validUpdateStrategy,
				},
			},
		},
		"read-write volume": {
			old: extensions.DaemonSet{
				ObjectMeta: api.ObjectMeta{Name: "", Namespace: api.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:       &extensions.LabelSelector{MatchLabels: validSelector},
					Template:       &validPodTemplateAbc.Template,
					UpdateStrategy: validUpdateStrategy,
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: api.ObjectMeta{Name: "abc", Namespace: api.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:       &extensions.LabelSelector{MatchLabels: validSelector},
					Template:       &readWriteVolumePodTemplate.Template,
					UpdateStrategy: validUpdateStrategy,
				},
			},
		},
		"invalid update strategy": {
			old: extensions.DaemonSet{
				ObjectMeta: api.ObjectMeta{Name: "", Namespace: api.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector:       &extensions.LabelSelector{MatchLabels: validSelector},
					Template:       &validPodTemplateAbc.Template,
					UpdateStrategy: validUpdateStrategy,
				},
			},
			update: extensions.DaemonSet{
				ObjectMeta: api.ObjectMeta{Name: "abc", Namespace: api.NamespaceDefault},
				Spec: extensions.DaemonSetSpec{
					Selector: &extensions.LabelSelector{MatchLabels: invalidSelector},
					Template: &validPodTemplateAbc.Template,
					UpdateStrategy: extensions.DaemonSetUpdateStrategy{
						Type:          extensions.RollingUpdateDaemonSetStrategyType,
						RollingUpdate: nil,
					},
				},
			},
		},
	}
	for testName, errorCase := range errorCases {
		if errs := ValidateDaemonSetUpdate(&errorCase.update, &errorCase.old); len(errs) == 0 {
			t.Errorf("expected failure: %s", testName)
		}
	}
}

func TestValidateDaemonSet(t *testing.T) {
	validSelector := map[string]string{"a": "b"}
	validUpdateStrategy := extensions.DaemonSetUpdateStrategy{
		Type: extensions.RollingUpdateDaemonSetStrategyType,
		RollingUpdate: &extensions.RollingUpdateDaemonSet{
			MaxUnavailable: intstr.FromInt(1),
		},
	}
	validPodTemplate := api.PodTemplate{
		Template: api.PodTemplateSpec{
			ObjectMeta: api.ObjectMeta{
				Labels: validSelector,
			},
			Spec: api.PodSpec{
				RestartPolicy: api.RestartPolicyAlways,
				DNSPolicy:     api.DNSClusterFirst,
				Containers:    []api.Container{{Name: "abc", Image: "image", ImagePullPolicy: "IfNotPresent"}},
			},
		},
	}
	invalidSelector := map[string]string{"NoUppercaseOrSpecialCharsLike=Equals": "b"}
	invalidPodTemplate := api.PodTemplate{
		Template: api.PodTemplateSpec{
			Spec: api.PodSpec{
				RestartPolicy: api.RestartPolicyAlways,
				DNSPolicy:     api.DNSClusterFirst,
			},
			ObjectMeta: api.ObjectMeta{
				Labels: invalidSelector,
			},
		},
	}
	successCases := []extensions.DaemonSet{
		{
			ObjectMeta: api.ObjectMeta{Name: "abc", Namespace: api.NamespaceDefault},
			Spec: extensions.DaemonSetSpec{
				Selector:       &extensions.LabelSelector{MatchLabels: validSelector},
				Template:       &validPodTemplate.Template,
				UpdateStrategy: validUpdateStrategy,
			},
		},
		{
			ObjectMeta: api.ObjectMeta{Name: "abc-123", Namespace: api.NamespaceDefault},
			Spec: extensions.DaemonSetSpec{
				Selector:       &extensions.LabelSelector{MatchLabels: validSelector},
				Template:       &validPodTemplate.Template,
				UpdateStrategy: validUpdateStrategy,
			},
		},
	}
	for _, successCase := range successCases {
		if errs := ValidateDaemonSet(&successCase); len(errs) != 0 {
			t.Errorf("expected success: %v", errs)
		}
	}

	errorCases := map[string]extensions.DaemonSet{
		"zero-length ID": {
			ObjectMeta: api.ObjectMeta{Name: "", Namespace: api.NamespaceDefault},
			Spec: extensions.DaemonSetSpec{
				Selector: &extensions.LabelSelector{MatchLabels: validSelector},
				Template: &validPodTemplate.Template,
			},
		},
		"missing-namespace": {
			ObjectMeta: api.ObjectMeta{Name: "abc-123"},
			Spec: extensions.DaemonSetSpec{
				Selector: &extensions.LabelSelector{MatchLabels: validSelector},
				Template: &validPodTemplate.Template,
			},
		},
		"empty selector": {
			ObjectMeta: api.ObjectMeta{Name: "abc", Namespace: api.NamespaceDefault},
			Spec: extensions.DaemonSetSpec{
				Template: &validPodTemplate.Template,
			},
		},
		"selector_doesnt_match": {
			ObjectMeta: api.ObjectMeta{Name: "abc", Namespace: api.NamespaceDefault},
			Spec: extensions.DaemonSetSpec{
				Selector: &extensions.LabelSelector{MatchLabels: map[string]string{"foo": "bar"}},
				Template: &validPodTemplate.Template,
			},
		},
		"invalid template": {
			ObjectMeta: api.ObjectMeta{Name: "abc", Namespace: api.NamespaceDefault},
			Spec: extensions.DaemonSetSpec{
				Selector: &extensions.LabelSelector{MatchLabels: validSelector},
			},
		},
		"invalid_label": {
			ObjectMeta: api.ObjectMeta{
				Name:      "abc-123",
				Namespace: api.NamespaceDefault,
				Labels: map[string]string{
					"NoUppercaseOrSpecialCharsLike=Equals": "bar",
				},
			},
			Spec: extensions.DaemonSetSpec{
				Selector: &extensions.LabelSelector{MatchLabels: validSelector},
				Template: &validPodTemplate.Template,
			},
		},
		"invalid_label 2": {
			ObjectMeta: api.ObjectMeta{
				Name:      "abc-123",
				Namespace: api.NamespaceDefault,
				Labels: map[string]string{
					"NoUppercaseOrSpecialCharsLike=Equals": "bar",
				},
			},
			Spec: extensions.DaemonSetSpec{
				Template: &invalidPodTemplate.Template,
			},
		},
		"invalid_annotation": {
			ObjectMeta: api.ObjectMeta{
				Name:      "abc-123",
				Namespace: api.NamespaceDefault,
				Annotations: map[string]string{
					"NoUppercaseOrSpecialCharsLike=Equals": "bar",
				},
			},
			Spec: extensions.DaemonSetSpec{
				Selector: &extensions.LabelSelector{MatchLabels: validSelector},
				Template: &validPodTemplate.Template,
			},
		},
		"invalid restart policy 1": {
			ObjectMeta: api.ObjectMeta{
				Name:      "abc-123",
				Namespace: api.NamespaceDefault,
			},
			Spec: extensions.DaemonSetSpec{
				Selector: &extensions.LabelSelector{MatchLabels: validSelector},
				Template: &api.PodTemplateSpec{
					Spec: api.PodSpec{
						RestartPolicy: api.RestartPolicyOnFailure,
						DNSPolicy:     api.DNSClusterFirst,
						Containers:    []api.Container{{Name: "ctr", Image: "image", ImagePullPolicy: "IfNotPresent"}},
					},
					ObjectMeta: api.ObjectMeta{
						Labels: validSelector,
					},
				},
			},
		},
		"invalid restart policy 2": {
			ObjectMeta: api.ObjectMeta{
				Name:      "abc-123",
				Namespace: api.NamespaceDefault,
			},
			Spec: extensions.DaemonSetSpec{
				Selector: &extensions.LabelSelector{MatchLabels: validSelector},
				Template: &api.PodTemplateSpec{
					Spec: api.PodSpec{
						RestartPolicy: api.RestartPolicyNever,
						DNSPolicy:     api.DNSClusterFirst,
						Containers:    []api.Container{{Name: "ctr", Image: "image", ImagePullPolicy: "IfNotPresent"}},
					},
					ObjectMeta: api.ObjectMeta{
						Labels: validSelector,
					},
				},
			},
		},
		"invalid update strategy - Type is not RollingUpdate": {
			ObjectMeta: api.ObjectMeta{
				Name:      "abc-123",
				Namespace: api.NamespaceDefault,
			},
			Spec: extensions.DaemonSetSpec{
				Selector: &extensions.LabelSelector{MatchLabels: validSelector},
				Template: &validPodTemplate.Template,
				UpdateStrategy: extensions.DaemonSetUpdateStrategy{
					Type: "",
					RollingUpdate: &extensions.RollingUpdateDaemonSet{
						MaxUnavailable: intstr.FromInt(1),
					},
				},
			},
		},
		"invalid update strategy - RollingUpdate field is nil": {
			ObjectMeta: api.ObjectMeta{
				Name:      "abc-123",
				Namespace: api.NamespaceDefault,
			},
			Spec: extensions.DaemonSetSpec{
				Selector: &extensions.LabelSelector{MatchLabels: validSelector},
				Template: &validPodTemplate.Template,
				UpdateStrategy: extensions.DaemonSetUpdateStrategy{
					Type:          extensions.RollingUpdateDaemonSetStrategyType,
					RollingUpdate: nil,
				},
			},
		},
		"invalid update strategy - MaxUnavailable is 0": {
			ObjectMeta: api.ObjectMeta{
				Name:      "abc-123",
				Namespace: api.NamespaceDefault,
			},
			Spec: extensions.DaemonSetSpec{
				Selector: &extensions.LabelSelector{MatchLabels: validSelector},
				Template: &validPodTemplate.Template,
				UpdateStrategy: extensions.DaemonSetUpdateStrategy{
					Type: extensions.RollingUpdateDaemonSetStrategyType,
					RollingUpdate: &extensions.RollingUpdateDaemonSet{
						MaxUnavailable:  intstr.FromInt(0),
						MinReadySeconds: 1,
					},
				},
			},
		},
		"invalid update strategy - MaxUnavailable is greater than 100%": {
			ObjectMeta: api.ObjectMeta{
				Name:      "abc-123",
				Namespace: api.NamespaceDefault,
			},
			Spec: extensions.DaemonSetSpec{
				Selector: &extensions.LabelSelector{MatchLabels: validSelector},
				Template: &validPodTemplate.Template,
				UpdateStrategy: extensions.DaemonSetUpdateStrategy{
					Type: extensions.RollingUpdateDaemonSetStrategyType,
					RollingUpdate: &extensions.RollingUpdateDaemonSet{
						MaxUnavailable:  intstr.FromString("150%"),
						MinReadySeconds: 1,
					},
				},
			},
		},
		"invalid update strategy - MaxUnavailable is negative": {
			ObjectMeta: api.ObjectMeta{
				Name:      "abc-123",
				Namespace: api.NamespaceDefault,
			},
			Spec: extensions.DaemonSetSpec{
				Selector: &extensions.LabelSelector{MatchLabels: validSelector},
				Template: &validPodTemplate.Template,
				UpdateStrategy: extensions.DaemonSetUpdateStrategy{
					Type: extensions.RollingUpdateDaemonSetStrategyType,
					RollingUpdate: &extensions.RollingUpdateDaemonSet{
						MaxUnavailable:  intstr.FromInt(-1),
						MinReadySeconds: 0,
					},
				},
			},
		},
		"invalid update strategy - MinReadySeconds is negative": {
			ObjectMeta: api.ObjectMeta{
				Name:      "abc-123",
				Namespace: api.NamespaceDefault,
			},
			Spec: extensions.DaemonSetSpec{
				Selector: &extensions.LabelSelector{MatchLabels: validSelector},
				Template: &validPodTemplate.Template,
				UpdateStrategy: extensions.DaemonSetUpdateStrategy{
					Type: extensions.RollingUpdateDaemonSetStrategyType,
					RollingUpdate: &extensions.RollingUpdateDaemonSet{
						MaxUnavailable:  intstr.FromInt(-1),
						MinReadySeconds: -1,
					},
				},
			},
		},
	}
	for k, v := range errorCases {
		errs := ValidateDaemonSet(&v)
		if len(errs) == 0 {
			t.Errorf("expected failure for %s", k)
		}
		for i := range errs {
			field := errs[i].Field
			if !strings.HasPrefix(field, "spec.template.") &&
				!strings.HasPrefix(field, "spec.updateStrategy") &&
				field != "metadata.name" &&
				field != "metadata.namespace" &&
				field != "spec.selector" &&
				field != "spec.template" &&
				field != "GCEPersistentDisk.ReadOnly" &&
				field != "spec.template.labels" &&
				field != "metadata.annotations" &&
				field != "metadata.labels" {
				t.Errorf("%s: missing prefix for: %v", k, errs[i])
			}
		}
	}
}

func validDeployment() *extensions.Deployment {
	return &extensions.Deployment{
		ObjectMeta: api.ObjectMeta{
			Name:      "abc",
			Namespace: api.NamespaceDefault,
		},
		Spec: extensions.DeploymentSpec{
			Selector: map[string]string{
				"name": "abc",
			},
			Template: api.PodTemplateSpec{
				ObjectMeta: api.ObjectMeta{
					Name:      "abc",
					Namespace: api.NamespaceDefault,
					Labels: map[string]string{
						"name": "abc",
					},
				},
				Spec: api.PodSpec{
					RestartPolicy: api.RestartPolicyAlways,
					DNSPolicy:     api.DNSDefault,
					Containers: []api.Container{
						{
							Name:            "nginx",
							Image:           "image",
							ImagePullPolicy: api.PullNever,
						},
					},
				},
			},
			UniqueLabelKey: "my-label",
		},
	}
}

func TestValidateDeployment(t *testing.T) {
	successCases := []*extensions.Deployment{
		validDeployment(),
	}
	for _, successCase := range successCases {
		if errs := ValidateDeployment(successCase); len(errs) != 0 {
			t.Errorf("expected success: %v", errs)
		}
	}

	errorCases := map[string]*extensions.Deployment{}
	errorCases["metadata.name: Required value"] = &extensions.Deployment{
		ObjectMeta: api.ObjectMeta{
			Namespace: api.NamespaceDefault,
		},
	}
	// selector should match the labels in pod template.
	invalidSelectorDeployment := validDeployment()
	invalidSelectorDeployment.Spec.Selector = map[string]string{
		"name": "def",
	}
	errorCases["`selector` does not match template `labels`"] = invalidSelectorDeployment

	// RestartPolicy should be always.
	invalidRestartPolicyDeployment := validDeployment()
	invalidRestartPolicyDeployment.Spec.Template.Spec.RestartPolicy = api.RestartPolicyNever
	errorCases["Unsupported value: \"Never\""] = invalidRestartPolicyDeployment

	// invalid unique label key.
	invalidUniqueLabelDeployment := validDeployment()
	invalidUniqueLabelDeployment.Spec.UniqueLabelKey = "abc/def/ghi"
	errorCases["spec.uniqueLabel: Invalid value"] = invalidUniqueLabelDeployment

	// rollingUpdate should be nil for recreate.
	invalidRecreateDeployment := validDeployment()
	invalidRecreateDeployment.Spec.Strategy = extensions.DeploymentStrategy{
		Type:          extensions.RecreateDeploymentStrategyType,
		RollingUpdate: &extensions.RollingUpdateDeployment{},
	}
	errorCases["may not be specified when strategy `type` is 'Recreate'"] = invalidRecreateDeployment

	// MaxSurge should be in the form of 20%.
	invalidMaxSurgeDeployment := validDeployment()
	invalidMaxSurgeDeployment.Spec.Strategy = extensions.DeploymentStrategy{
		Type: extensions.RollingUpdateDeploymentStrategyType,
		RollingUpdate: &extensions.RollingUpdateDeployment{
			MaxSurge: intstr.FromString("20Percent"),
		},
	}
	errorCases["must be an integer or percentage"] = invalidMaxSurgeDeployment

	// MaxSurge and MaxUnavailable cannot both be zero.
	invalidRollingUpdateDeployment := validDeployment()
	invalidRollingUpdateDeployment.Spec.Strategy = extensions.DeploymentStrategy{
		Type: extensions.RollingUpdateDeploymentStrategyType,
		RollingUpdate: &extensions.RollingUpdateDeployment{
			MaxSurge:       intstr.FromString("0%"),
			MaxUnavailable: intstr.FromInt(0),
		},
	}
	errorCases["may not be 0 when `maxSurge` is 0"] = invalidRollingUpdateDeployment

	// MaxUnavailable should not be more than 100%.
	invalidMaxUnavailableDeployment := validDeployment()
	invalidMaxUnavailableDeployment.Spec.Strategy = extensions.DeploymentStrategy{
		Type: extensions.RollingUpdateDeploymentStrategyType,
		RollingUpdate: &extensions.RollingUpdateDeployment{
			MaxUnavailable: intstr.FromString("110%"),
		},
	}
	errorCases["must not be greater than 100%"] = invalidMaxUnavailableDeployment

	for k, v := range errorCases {
		errs := ValidateDeployment(v)
		if len(errs) == 0 {
			t.Errorf("[%s] expected failure", k)
		} else if !strings.Contains(errs[0].Error(), k) {
			t.Errorf("unexpected error: %q, expected: %q", errs[0].Error(), k)
		}
	}
}

func TestValidateJob(t *testing.T) {
	validSelector := &extensions.LabelSelector{
		MatchLabels: map[string]string{"a": "b"},
	}
	validPodTemplateSpec := api.PodTemplateSpec{
		ObjectMeta: api.ObjectMeta{
			Labels: validSelector.MatchLabels,
		},
		Spec: api.PodSpec{
			RestartPolicy: api.RestartPolicyOnFailure,
			DNSPolicy:     api.DNSClusterFirst,
			Containers:    []api.Container{{Name: "abc", Image: "image", ImagePullPolicy: "IfNotPresent"}},
		},
	}
	successCases := []extensions.Job{
		{
			ObjectMeta: api.ObjectMeta{
				Name:      "myjob",
				Namespace: api.NamespaceDefault,
			},
			Spec: extensions.JobSpec{
				Selector: validSelector,
				Template: validPodTemplateSpec,
			},
		},
	}
	for _, successCase := range successCases {
		if errs := ValidateJob(&successCase); len(errs) != 0 {
			t.Errorf("expected success: %v", errs)
		}
	}
	negative := -1
	negative64 := int64(-1)
	errorCases := map[string]extensions.Job{
		"spec.parallelism:must be greater than or equal to 0": {
			ObjectMeta: api.ObjectMeta{
				Name:      "myjob",
				Namespace: api.NamespaceDefault,
			},
			Spec: extensions.JobSpec{
				Parallelism: &negative,
				Selector:    validSelector,
				Template:    validPodTemplateSpec,
			},
		},
		"spec.completions:must be greater than or equal to 0": {
			ObjectMeta: api.ObjectMeta{
				Name:      "myjob",
				Namespace: api.NamespaceDefault,
			},
			Spec: extensions.JobSpec{
				Completions: &negative,
				Selector:    validSelector,
				Template:    validPodTemplateSpec,
			},
		},
		"spec.activeDeadlineSeconds:must be greater than or equal to 0": {
			ObjectMeta: api.ObjectMeta{
				Name:      "myjob",
				Namespace: api.NamespaceDefault,
			},
			Spec: extensions.JobSpec{
				ActiveDeadlineSeconds: &negative64,
				Selector:              validSelector,
				Template:              validPodTemplateSpec,
			},
		},
		"spec.selector:Required value": {
			ObjectMeta: api.ObjectMeta{
				Name:      "myjob",
				Namespace: api.NamespaceDefault,
			},
			Spec: extensions.JobSpec{
				Template: validPodTemplateSpec,
			},
		},
		"spec.template.metadata.labels: Invalid value: {\"y\":\"z\"}: `selector` does not match template `labels`": {
			ObjectMeta: api.ObjectMeta{
				Name:      "myjob",
				Namespace: api.NamespaceDefault,
			},
			Spec: extensions.JobSpec{
				Selector: validSelector,
				Template: api.PodTemplateSpec{
					ObjectMeta: api.ObjectMeta{
						Labels: map[string]string{"y": "z"},
					},
					Spec: api.PodSpec{
						RestartPolicy: api.RestartPolicyOnFailure,
						DNSPolicy:     api.DNSClusterFirst,
						Containers:    []api.Container{{Name: "abc", Image: "image", ImagePullPolicy: "IfNotPresent"}},
					},
				},
			},
		},
		"spec.template.spec.restartPolicy: Unsupported value": {
			ObjectMeta: api.ObjectMeta{
				Name:      "myjob",
				Namespace: api.NamespaceDefault,
			},
			Spec: extensions.JobSpec{
				Selector: validSelector,
				Template: api.PodTemplateSpec{
					ObjectMeta: api.ObjectMeta{
						Labels: validSelector.MatchLabels,
					},
					Spec: api.PodSpec{
						RestartPolicy: api.RestartPolicyAlways,
						DNSPolicy:     api.DNSClusterFirst,
						Containers:    []api.Container{{Name: "abc", Image: "image", ImagePullPolicy: "IfNotPresent"}},
					},
				},
			},
		},
	}

	for k, v := range errorCases {
		errs := ValidateJob(&v)
		if len(errs) == 0 {
			t.Errorf("expected failure for %s", k)
		} else {
			s := strings.Split(k, ":")
			err := errs[0]
			if err.Field != s[0] || !strings.Contains(err.Error(), s[1]) {
				t.Errorf("unexpected error: %v, expected: %s", err, k)
			}
		}
	}
}

type ingressRules map[string]string

func TestValidateIngress(t *testing.T) {
	defaultBackend := extensions.IngressBackend{
		ServiceName: "default-backend",
		ServicePort: intstr.FromInt(80),
	}

	newValid := func() extensions.Ingress {
		return extensions.Ingress{
			ObjectMeta: api.ObjectMeta{
				Name:      "foo",
				Namespace: api.NamespaceDefault,
			},
			Spec: extensions.IngressSpec{
				Backend: &extensions.IngressBackend{
					ServiceName: "default-backend",
					ServicePort: intstr.FromInt(80),
				},
				Rules: []extensions.IngressRule{
					{
						Host: "foo.bar.com",
						IngressRuleValue: extensions.IngressRuleValue{
							HTTP: &extensions.HTTPIngressRuleValue{
								Paths: []extensions.HTTPIngressPath{
									{
										Path:    "/foo",
										Backend: defaultBackend,
									},
								},
							},
						},
					},
				},
			},
			Status: extensions.IngressStatus{
				LoadBalancer: api.LoadBalancerStatus{
					Ingress: []api.LoadBalancerIngress{
						{IP: "127.0.0.1"},
					},
				},
			},
		}
	}
	servicelessBackend := newValid()
	servicelessBackend.Spec.Backend.ServiceName = ""
	invalidNameBackend := newValid()
	invalidNameBackend.Spec.Backend.ServiceName = "defaultBackend"
	noPortBackend := newValid()
	noPortBackend.Spec.Backend = &extensions.IngressBackend{ServiceName: defaultBackend.ServiceName}
	noForwardSlashPath := newValid()
	noForwardSlashPath.Spec.Rules[0].IngressRuleValue.HTTP.Paths = []extensions.HTTPIngressPath{
		{
			Path:    "invalid",
			Backend: defaultBackend,
		},
	}
	noPaths := newValid()
	noPaths.Spec.Rules[0].IngressRuleValue.HTTP.Paths = []extensions.HTTPIngressPath{}
	badHost := newValid()
	badHost.Spec.Rules[0].Host = "foobar:80"
	badRegexPath := newValid()
	badPathExpr := "/invalid["
	badRegexPath.Spec.Rules[0].IngressRuleValue.HTTP.Paths = []extensions.HTTPIngressPath{
		{
			Path:    badPathExpr,
			Backend: defaultBackend,
		},
	}
	badPathErr := fmt.Sprintf("spec.rules[0].http.paths[0].path: Invalid value: '%v'", badPathExpr)
	hostIP := "127.0.0.1"
	badHostIP := newValid()
	badHostIP.Spec.Rules[0].Host = hostIP
	badHostIPErr := fmt.Sprintf("spec.rules[0].host: Invalid value: '%v'", hostIP)

	errorCases := map[string]extensions.Ingress{
		"spec.backend.serviceName: Required value":        servicelessBackend,
		"spec.backend.serviceName: Invalid value":         invalidNameBackend,
		"spec.backend.servicePort: Invalid value":         noPortBackend,
		"spec.rules[0].host: Invalid value":               badHost,
		"spec.rules[0].http.paths: Required value":        noPaths,
		"spec.rules[0].http.paths[0].path: Invalid value": noForwardSlashPath,
	}
	errorCases[badPathErr] = badRegexPath
	errorCases[badHostIPErr] = badHostIP

	for k, v := range errorCases {
		errs := ValidateIngress(&v)
		if len(errs) == 0 {
			t.Errorf("expected failure for %q", k)
		} else {
			s := strings.Split(k, ":")
			err := errs[0]
			if err.Field != s[0] || !strings.Contains(err.Error(), s[1]) {
				t.Errorf("unexpected error: %q, expected: %q", err, k)
			}
		}
	}
}

func TestValidateIngressStatusUpdate(t *testing.T) {
	defaultBackend := extensions.IngressBackend{
		ServiceName: "default-backend",
		ServicePort: intstr.FromInt(80),
	}

	newValid := func() extensions.Ingress {
		return extensions.Ingress{
			ObjectMeta: api.ObjectMeta{
				Name:            "foo",
				Namespace:       api.NamespaceDefault,
				ResourceVersion: "9",
			},
			Spec: extensions.IngressSpec{
				Backend: &extensions.IngressBackend{
					ServiceName: "default-backend",
					ServicePort: intstr.FromInt(80),
				},
				Rules: []extensions.IngressRule{
					{
						Host: "foo.bar.com",
						IngressRuleValue: extensions.IngressRuleValue{
							HTTP: &extensions.HTTPIngressRuleValue{
								Paths: []extensions.HTTPIngressPath{
									{
										Path:    "/foo",
										Backend: defaultBackend,
									},
								},
							},
						},
					},
				},
			},
			Status: extensions.IngressStatus{
				LoadBalancer: api.LoadBalancerStatus{
					Ingress: []api.LoadBalancerIngress{
						{IP: "127.0.0.1", Hostname: "foo.bar.com"},
					},
				},
			},
		}
	}
	oldValue := newValid()
	newValue := newValid()
	newValue.Status = extensions.IngressStatus{
		LoadBalancer: api.LoadBalancerStatus{
			Ingress: []api.LoadBalancerIngress{
				{IP: "127.0.0.2", Hostname: "foo.com"},
			},
		},
	}
	invalidIP := newValid()
	invalidIP.Status = extensions.IngressStatus{
		LoadBalancer: api.LoadBalancerStatus{
			Ingress: []api.LoadBalancerIngress{
				{IP: "abcd", Hostname: "foo.com"},
			},
		},
	}
	invalidHostname := newValid()
	invalidHostname.Status = extensions.IngressStatus{
		LoadBalancer: api.LoadBalancerStatus{
			Ingress: []api.LoadBalancerIngress{
				{IP: "127.0.0.1", Hostname: "127.0.0.1"},
			},
		},
	}

	errs := ValidateIngressStatusUpdate(&newValue, &oldValue)
	if len(errs) != 0 {
		t.Errorf("Unexpected error %v", errs)
	}

	errorCases := map[string]extensions.Ingress{
		"status.loadBalancer.ingress[0].ip: Invalid value":       invalidIP,
		"status.loadBalancer.ingress[0].hostname: Invalid value": invalidHostname,
	}
	for k, v := range errorCases {
		errs := ValidateIngressStatusUpdate(&v, &oldValue)
		if len(errs) == 0 {
			t.Errorf("expected failure for %s", k)
		} else {
			s := strings.Split(k, ":")
			err := errs[0]
			if err.Field != s[0] || !strings.Contains(err.Error(), s[1]) {
				t.Errorf("unexpected error: %q, expected: %q", err, k)
			}
		}
	}
}

func TestValidateClusterAutoscaler(t *testing.T) {
	successCases := []extensions.ClusterAutoscaler{
		{
			ObjectMeta: api.ObjectMeta{
				Name:      "ClusterAutoscaler",
				Namespace: api.NamespaceDefault,
			},
			Spec: extensions.ClusterAutoscalerSpec{
				MinNodes: 1,
				MaxNodes: 5,
				TargetUtilization: []extensions.NodeUtilization{
					{
						Resource: extensions.CpuRequest,
						Value:    0.7,
					},
				},
			},
		},
	}
	for _, successCase := range successCases {
		if errs := ValidateClusterAutoscaler(&successCase); len(errs) != 0 {
			t.Errorf("expected success: %v", errs)
		}
	}

	errorCases := map[string]extensions.ClusterAutoscaler{
		"must be 'ClusterAutoscaler'": {
			ObjectMeta: api.ObjectMeta{
				Name:      "TestClusterAutoscaler",
				Namespace: api.NamespaceDefault,
			},
			Spec: extensions.ClusterAutoscalerSpec{
				MinNodes: 1,
				MaxNodes: 5,
				TargetUtilization: []extensions.NodeUtilization{
					{
						Resource: extensions.CpuRequest,
						Value:    0.7,
					},
				},
			},
		},
		"must be 'default'": {
			ObjectMeta: api.ObjectMeta{
				Name:      "ClusterAutoscaler",
				Namespace: "test",
			},
			Spec: extensions.ClusterAutoscalerSpec{
				MinNodes: 1,
				MaxNodes: 5,
				TargetUtilization: []extensions.NodeUtilization{
					{
						Resource: extensions.CpuRequest,
						Value:    0.7,
					},
				},
			},
		},

		`must be greater than or equal to 0`: {
			ObjectMeta: api.ObjectMeta{
				Name:      "ClusterAutoscaler",
				Namespace: api.NamespaceDefault,
			},
			Spec: extensions.ClusterAutoscalerSpec{
				MinNodes: -1,
				MaxNodes: 5,
				TargetUtilization: []extensions.NodeUtilization{
					{
						Resource: extensions.CpuRequest,
						Value:    0.7,
					},
				},
			},
		},
		"must be greater than or equal to `minNodes`": {
			ObjectMeta: api.ObjectMeta{
				Name:      "ClusterAutoscaler",
				Namespace: api.NamespaceDefault,
			},
			Spec: extensions.ClusterAutoscalerSpec{
				MinNodes: 10,
				MaxNodes: 5,
				TargetUtilization: []extensions.NodeUtilization{
					{
						Resource: extensions.CpuRequest,
						Value:    0.7,
					},
				},
			},
		},
		"Required value": {
			ObjectMeta: api.ObjectMeta{
				Name:      "ClusterAutoscaler",
				Namespace: api.NamespaceDefault,
			},
			Spec: extensions.ClusterAutoscalerSpec{
				MinNodes:          1,
				MaxNodes:          5,
				TargetUtilization: []extensions.NodeUtilization{},
			},
		},
	}

	for k, v := range errorCases {
		errs := ValidateClusterAutoscaler(&v)
		if len(errs) == 0 {
			t.Errorf("[%s] expected failure", k)
		} else if !strings.Contains(errs[0].Error(), k) {
			t.Errorf("unexpected error: %v, expected: %q", errs[0], k)
		}
	}
}

func TestValidateScale(t *testing.T) {
	successCases := []extensions.Scale{
		{
			ObjectMeta: api.ObjectMeta{
				Name:      "frontend",
				Namespace: api.NamespaceDefault,
			},
			Spec: extensions.ScaleSpec{
				Replicas: 1,
			},
		},
		{
			ObjectMeta: api.ObjectMeta{
				Name:      "frontend",
				Namespace: api.NamespaceDefault,
			},
			Spec: extensions.ScaleSpec{
				Replicas: 10,
			},
		},
		{
			ObjectMeta: api.ObjectMeta{
				Name:      "frontend",
				Namespace: api.NamespaceDefault,
			},
			Spec: extensions.ScaleSpec{
				Replicas: 0,
			},
		},
	}

	for _, successCase := range successCases {
		if errs := ValidateScale(&successCase); len(errs) != 0 {
			t.Errorf("expected success: %v", errs)
		}
	}

	errorCases := []struct {
		scale extensions.Scale
		msg   string
	}{
		{
			scale: extensions.Scale{
				ObjectMeta: api.ObjectMeta{
					Name:      "frontend",
					Namespace: api.NamespaceDefault,
				},
				Spec: extensions.ScaleSpec{
					Replicas: -1,
				},
			},
			msg: "must be greater than or equal to 0",
		},
	}

	for _, c := range errorCases {
		if errs := ValidateScale(&c.scale); len(errs) == 0 {
			t.Errorf("expected failure for %s", c.msg)
		} else if !strings.Contains(errs[0].Error(), c.msg) {
			t.Errorf("unexpected error: %v, expected: %s", errs[0], c.msg)
		}
	}
}

func newInt(val int) *int {
	p := new(int)
	*p = val
	return p
}

func TestValidateConfigMap(t *testing.T) {
	newConfigMap := func(name, namespace string, data map[string]string) extensions.ConfigMap {
		return extensions.ConfigMap{
			ObjectMeta: api.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Data: data,
		}
	}

	var (
		validConfigMap = newConfigMap("validname", "validns", map[string]string{"key": "value"})
		maxKeyLength   = newConfigMap("validname", "validns", map[string]string{strings.Repeat("a", 253): "value"})

		emptyName        = newConfigMap("", "validns", nil)
		invalidName      = newConfigMap("NoUppercaseOrSpecialCharsLike=Equals", "validns", nil)
		emptyNs          = newConfigMap("validname", "", nil)
		invalidNs        = newConfigMap("validname", "NoUppercaseOrSpecialCharsLike=Equals", nil)
		invalidKey       = newConfigMap("validname", "validns", map[string]string{"a..b": "value"})
		leadingDotKey    = newConfigMap("validname", "validns", map[string]string{".ab": "value"})
		dotKey           = newConfigMap("validname", "validns", map[string]string{".": "value"})
		doubleDotKey     = newConfigMap("validname", "validns", map[string]string{"..": "value"})
		overMaxKeyLength = newConfigMap("validname", "validns", map[string]string{strings.Repeat("a", 254): "value"})
	)

	tests := map[string]struct {
		cfg     extensions.ConfigMap
		isValid bool
	}{
		"valid":               {validConfigMap, true},
		"max key length":      {maxKeyLength, true},
		"leading dot key":     {leadingDotKey, true},
		"empty name":          {emptyName, false},
		"invalid name":        {invalidName, false},
		"invalid key":         {invalidKey, false},
		"empty namespace":     {emptyNs, false},
		"invalid namespace":   {invalidNs, false},
		"dot key":             {dotKey, false},
		"double dot key":      {doubleDotKey, false},
		"over max key length": {overMaxKeyLength, false},
	}

	for name, tc := range tests {
		errs := ValidateConfigMap(&tc.cfg)
		if tc.isValid && len(errs) > 0 {
			t.Errorf("%v: unexpected error: %v", name, errs)
		}
		if !tc.isValid && len(errs) == 0 {
			t.Errorf("%v: unexpected non-error", name)
		}
	}
}

func TestValidateConfigMapUpdate(t *testing.T) {
	newConfigMap := func(version, name, namespace string, data map[string]string) extensions.ConfigMap {
		return extensions.ConfigMap{
			ObjectMeta: api.ObjectMeta{
				Name:            name,
				Namespace:       namespace,
				ResourceVersion: version,
			},
			Data: data,
		}
	}

	var (
		validConfigMap = newConfigMap("1", "validname", "validns", map[string]string{"key": "value"})
		noVersion      = newConfigMap("", "validname", "validns", map[string]string{"key": "value"})
	)

	cases := []struct {
		name    string
		newCfg  extensions.ConfigMap
		oldCfg  extensions.ConfigMap
		isValid bool
	}{
		{
			name:    "valid",
			newCfg:  validConfigMap,
			oldCfg:  validConfigMap,
			isValid: true,
		},
		{
			name:    "invalid",
			newCfg:  noVersion,
			oldCfg:  validConfigMap,
			isValid: false,
		},
	}

	for _, tc := range cases {
		errs := ValidateConfigMapUpdate(&tc.newCfg, &tc.oldCfg)
		if tc.isValid && len(errs) > 0 {
			t.Errorf("%v: unexpected error: %v", tc.name, errs)
		}
		if !tc.isValid && len(errs) == 0 {
			t.Errorf("%v: unexpected non-error", tc.name)
		}
	}
}
