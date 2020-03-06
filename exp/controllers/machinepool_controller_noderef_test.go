/*
Copyright 2019 The Kubernetes Authors.

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
	"context"
	"testing"

	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/log"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
)

func TestMachinePoolGetNodeReference(t *testing.T) {
	g := NewWithT(t)

	g.Expect(clusterv1.AddToScheme(scheme.Scheme)).To(Succeed())

	r := &MachinePoolReconciler{
		Client:   fake.NewFakeClientWithScheme(scheme.Scheme),
		Log:      log.Log,
		recorder: record.NewFakeRecorder(32),
	}

	nodeList := []runtime.Object{
		&corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: "node-1",
			},
			Spec: corev1.NodeSpec{
				ProviderID: "aws://us-east-1/id-node-1",
			},
		},
		&corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: "node-2",
			},
			Spec: corev1.NodeSpec{
				ProviderID: "aws://us-west-2/id-node-2",
			},
		},
		&corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: "gce-node-2",
			},
			Spec: corev1.NodeSpec{
				ProviderID: "gce://us-central1/gce-id-node-2",
			},
		},
		&corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: "azure-node-4",
			},
			Spec: corev1.NodeSpec{
				ProviderID: "azure://westus2/id-node-4",
			},
		},
	}

	client := fake.NewFakeClientWithScheme(scheme.Scheme, nodeList...)

	testCases := []struct {
		name           string
		providerIDList []string
		expected       *getNodeReferencesResult
		err            error
	}{
		{
			name:           "valid provider id, valid aws node",
			providerIDList: []string{"aws://us-east-1/id-node-1"},
			expected: &getNodeReferencesResult{
				references: []corev1.ObjectReference{
					{Name: "node-1"},
				},
			},
		},
		{
			name:           "valid provider id, valid aws node",
			providerIDList: []string{"aws://us-west-2/id-node-2"},
			expected: &getNodeReferencesResult{
				references: []corev1.ObjectReference{
					{Name: "node-2"},
				},
			},
		},
		{
			name:           "valid provider id, valid gce node",
			providerIDList: []string{"gce://us-central1/gce-id-node-2"},
			expected: &getNodeReferencesResult{
				references: []corev1.ObjectReference{
					{Name: "gce-node-2"},
				},
			},
		},
		{
			name:           "valid provider id, valid azure node",
			providerIDList: []string{"azure://westus2/id-node-4"},
			expected: &getNodeReferencesResult{
				references: []corev1.ObjectReference{
					{Name: "azure-node-4"},
				},
			},
		},
		{
			name:           "valid provider ids, valid azure and aws nodes",
			providerIDList: []string{"aws://us-east-1/id-node-1", "azure://westus2/id-node-4"},
			expected: &getNodeReferencesResult{
				references: []corev1.ObjectReference{
					{Name: "node-1"},
					{Name: "azure-node-4"},
				},
			},
		},
		{
			name:           "valid provider id, no node found",
			providerIDList: []string{"aws:///id-node-100"},
			expected:       nil,
			err:            ErrNoAvailableNodes,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			gt := NewWithT(t)

			result, err := r.getNodeReferences(context.TODO(), client, test.providerIDList)
			if test.err == nil {
				g.Expect(err).To(BeNil())
			} else {
				gt.Expect(err).NotTo(BeNil())
				gt.Expect(err).To(Equal(test.err), "Expected error %v, got %v", test.err, err)
			}

			if test.expected == nil && len(result.references) == 0 {
				return
			}

			gt.Expect(len(result.references)).To(Equal(len(test.expected.references)), "Expected NodeRef count to be %v, got %v", len(result.references), len(test.expected.references))

			for n := range test.expected.references {
				gt.Expect(result.references[n].Name).To(Equal(test.expected.references[n].Name), "Expected NodeRef's name to be %v, got %v", result.references[n].Name, test.expected.references[n].Name)
				gt.Expect(result.references[n].Namespace).To(Equal(test.expected.references[n].Namespace), "Expected NodeRef's namespace to be %v, got %v", result.references[n].Namespace, test.expected.references[n].Namespace)
			}
		})

	}
}
