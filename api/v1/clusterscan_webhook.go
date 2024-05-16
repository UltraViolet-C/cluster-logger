/*
Copyright 2024.

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

package v1

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var clusterscanlog = logf.Log.WithName("clusterscan-resource")

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *ClusterScan) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-log-my-domain-v1-clusterscan,mutating=true,failurePolicy=fail,sideEffects=None,groups=log.my.domain,resources=clusterscans,verbs=create;update,versions=v1,name=mclusterscan.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &ClusterScan{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *ClusterScan) Default() {
	clusterscanlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
	if r.Spec.Name == "" {
		r.Spec.Name = "BasicClusterScan"
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
//+kubebuilder:webhook:path=/validate-log-my-domain-v1-clusterscan,mutating=false,failurePolicy=fail,sideEffects=None,groups=log.my.domain,resources=clusterscans,verbs=create;update,versions=v1,name=vclusterscan.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &ClusterScan{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *ClusterScan) ValidateCreate() (admission.Warnings, error) {
	clusterscanlog.Info("validate create", "name", r.Name)

	return nil, r.validateClusterScan()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *ClusterScan) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	clusterscanlog.Info("validate update", "name", r.Name)

	return nil, r.validateClusterScan()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *ClusterScan) ValidateDelete() (admission.Warnings, error) {
	clusterscanlog.Info("validate delete", "name", r.Name)

	// no need to validate the scan on delete
	return nil, nil
}

func (r *ClusterScan) validateClusterScan() error {
	var allErrs field.ErrorList
	if err := r.validateClusterScanSpec(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "log.my.domain", Kind: "ClusterScan"},
		r.Name, allErrs)
}

func (r *ClusterScan) validateClusterScanSpec() *field.Error {
	// check if version is correctly formatted. For the purpose of this tool, it will just check that it is "v1"
	if r.Spec.Version != "v1" {
		return field.Invalid(field.NewPath("spec").Child("name"), r.Spec.Version, "Incorrect scan version")
	}
	// in a more fleshed out tool, this would validate the full clusterScan but since the scan does not actually represent
	// anything, I just validate the version string.
	return nil
}
