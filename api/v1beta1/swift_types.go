/*
Copyright 2022.

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

package v1beta1

import (
	"github.com/openstack-k8s-operators/lib-common/modules/common/util"
	condition "github.com/openstack-k8s-operators/lib-common/modules/common/condition"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Container image fall-back defaults
	ContainerImageAccount = "quay.io/podified-antelope-centos9/openstack-swift-account:current-podified"
	ContainerImageContainer = "quay.io/podified-antelope-centos9/openstack-swift-container:current-podified"
	ContainerImageObject = "quay.io/podified-antelope-centos9/openstack-swift-object:current-podified"
	ContainerImageProxy = "quay.io/podified-antelope-centos9/openstack-swift-proxy-server:current-podified"
	ContainerImageMemcached = "quay.io/podified-antelope-centos9/openstack-memcached:current-podified"
)

// SwiftSpec defines the desired state of Swift
type SwiftSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Required
	// +kubebuilder:default=swift-ring
	// Name of ConfigMap containing Swift rings
	RingConfigMap string `json:"ringConfigMap"`

	// +kubebuilder:validation:Required
        // SwiftRing - Spec definition for the Ring service of this Swift deployment
        SwiftRing SwiftRingSpec `json:"swiftRing"`

	// +kubebuilder:validation:Required
        // SwiftStorage - Spec definition for the Storage service of this Swift deployment
        SwiftStorage SwiftStorageSpec `json:"swiftStorage"`

	// +kubebuilder:validation:Required
        // SwiftProxy - Spec definition for the Proxy service of this Swift deployment
        SwiftProxy SwiftProxySpec `json:"swiftProxy"`

	// +kubebuilder:validation:Required
	// +kubebuilder:default=swift-conf
	// Name of Secret containing swift.conf
	SwiftConfSecret string `json:"swiftConfSecret,omitempty"`
}

// SwiftStatus defines the observed state of Swift
type SwiftStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Conditions
	Conditions condition.Conditions `json:"conditions,omitempty" optional:"true"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Swift is the Schema for the swifts API
type Swift struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SwiftSpec   `json:"spec,omitempty"`
	Status SwiftStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SwiftList contains a list of Swift
type SwiftList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Swift `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Swift{}, &SwiftList{})
}

// RbacConditionsSet - set the conditions for the rbac object
func (instance Swift) RbacConditionsSet(c *condition.Condition) {
	instance.Status.Conditions.Set(c)
}

// RbacNamespace - return the namespace
func (instance Swift) RbacNamespace() string {
	return instance.Namespace
}

// RbacResourceName - return the name to be used for rbac objects (serviceaccount, role, rolebinding)
func (instance Swift) RbacResourceName() string {
	return "swift-" + instance.Name
}

// IsReady - returns true if all subresources Ready condition is true
func (instance Swift) IsReady() bool {
	return instance.Status.Conditions.IsTrue(SwiftRingReadyCondition) &&
		instance.Status.Conditions.IsTrue(SwiftStorageReadyCondition) &&
		instance.Status.Conditions.IsTrue(SwiftProxyReadyCondition)
}

// SetupDefaults - initializes any CRD field defaults based on environment variables (the defaulting mechanism itself is implemented via webhooks)
func SetupDefaults() {
	// Acquire environmental defaults and initialize Swift defaults with them
	swiftDefaults := SwiftDefaults{
		AccountContainerImageURL:	util.GetEnvVar("SWIFT_ACCOUNT_IMAGE_URL_DEFAULT", ContainerImageAccount),
		ContainerContainerImageURL:	util.GetEnvVar("SWIFT_CONTAINER_IMAGE_URL_DEFAULT", ContainerImageContainer),
		ObjectContainerImageURL:	util.GetEnvVar("SWIFT_OBJECT_IMAGE_URL_DEFAULT", ContainerImageObject),
		ProxyContainerImageURL:		util.GetEnvVar("SWIFT_PROXY_IMAGE_URL_DEFAULT", ContainerImageProxy),
		MemcachedContainerImageURL:	util.GetEnvVar("SWIFT_MEMCACHED_IMAGE_URL_DEFAULT", ContainerImageMemcached),
	}

	SetupSwiftDefaults(swiftDefaults)
}
