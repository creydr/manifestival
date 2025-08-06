package manifestival_test

import (
	"testing"

	. "github.com/manifestival/manifestival"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestByKind(t *testing.T) {
	// Create test resources in random order
	resources := []unstructured.Unstructured{
		{
			Object: map[string]interface{}{
				"apiVersion": "batch/v1",
				"kind":       "Job",
				"metadata": map[string]interface{}{
					"name": "test-job",
				},
			},
		},
		{
			Object: map[string]interface{}{
				"apiVersion": "apps/v1",
				"kind":       "Deployment",
				"metadata": map[string]interface{}{
					"name": "test-deployment",
				},
			},
		},
		{
			Object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Namespace",
				"metadata": map[string]interface{}{
					"name": "test-namespace",
				},
			},
		},
		{
			Object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Secret",
				"metadata": map[string]interface{}{
					"name": "test-secret",
				},
			},
		},
		{
			Object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Service",
				"metadata": map[string]interface{}{
					"name": "test-service",
				},
			},
		},
		{
			Object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "ConfigMap",
				"metadata": map[string]interface{}{
					"name": "test-configmap",
				},
			},
		},
		{
			Object: map[string]interface{}{
				"apiVersion": "apiextensions.k8s.io/v1",
				"kind":       "CustomResourceDefinition",
				"metadata": map[string]interface{}{
					"name": "test-crd",
				},
			},
		},
		{
			Object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "ServiceAccount",
				"metadata": map[string]interface{}{
					"name": "test-serviceaccount",
				},
			},
		},
	}

	manifest := Slice(resources)
	m, err := ManifestFrom(manifest)
	if err != nil {
		t.Fatalf("Failed to create manifest: %v", err)
	}

	// Sort by Kubernetes dependency order
	sorted := m.Sort(ByKindPriority())
	sortedResources := sorted.Resources()

	if len(sortedResources) != 8 {
		t.Errorf("Expected 8 resources, got %d", len(sortedResources))
	}

	// Check that dependencies come before dependents
	expectedOrder := []string{
		"Namespace",
		"ServiceAccount",
		"Secret",
		"ConfigMap",
		"CustomResourceDefinition",
		"Service",
		"Deployment",
		"Job",
	}

	for i, resource := range sortedResources {
		if resource.GetKind() != expectedOrder[i] {
			t.Errorf("Expected kind %s at position %d, got %s", expectedOrder[i], i, resource.GetKind())
		}
	}
}

func TestByKindWithUnknownKinds(t *testing.T) {
	// Create test resources including unknown kinds
	resources := []unstructured.Unstructured{
		{
			Object: map[string]interface{}{
				"apiVersion": "custom.io/v1",
				"kind":       "UnknownKind",
				"metadata": map[string]interface{}{
					"name": "test-unknown",
				},
			},
		},
		{
			Object: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Namespace",
				"metadata": map[string]interface{}{
					"name": "test-namespace",
				},
			},
		},
		{
			Object: map[string]interface{}{
				"apiVersion": "custom.io/v1",
				"kind":       "AnotherUnknown",
				"metadata": map[string]interface{}{
					"name": "test-another",
				},
			},
		},
		{
			Object: map[string]interface{}{
				"apiVersion": "apps/v1",
				"kind":       "Deployment",
				"metadata": map[string]interface{}{
					"name": "test-deployment",
				},
			},
		},
	}

	manifest := Slice(resources)
	m, err := ManifestFrom(manifest)
	if err != nil {
		t.Fatalf("Failed to create manifest: %v", err)
	}

	// Sort by Kubernetes dependency order
	sorted := m.Sort(ByKindPriority())
	sortedResources := sorted.Resources()

	if len(sortedResources) != 4 {
		t.Errorf("Expected 4 resources, got %d", len(sortedResources))
	}

	// Check that known kinds come first, unknown kinds come last
	expectedOrder := []string{
		"Namespace",
		"Deployment",
		"AnotherUnknown",
		"UnknownKind",
	}

	for i, resource := range sortedResources {
		if resource.GetKind() != expectedOrder[i] {
			t.Errorf("Expected kind %s at position %d, got %s", expectedOrder[i], i, resource.GetKind())
		}
	}
}
