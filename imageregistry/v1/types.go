package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	operatorv1 "github.com/openshift/api/operator/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ConfigList is a slice of Config objects.
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type ConfigList struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard list's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata"`
	Items           []Config `json:"items"`
}

const (
	// StorageManagementStateManaged indicates the operator is managing the underlying storage.
	StorageManagementStateManaged = "Managed"
	// StorageManagementStateUnmanaged indicates the operator is not managing the underlying
	// storage.
	StorageManagementStateUnmanaged = "Unmanaged"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Config is the configuration object for a registry instance managed by
// the registry operator
//
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=configs,scope=Cluster
// +openshift:api-approved.openshift.io=https://github.com/openshift/api/pull/519
// +openshift:file-pattern=operatorOrdering=00
type Config struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata"`

	Spec ImageRegistrySpec `json:"spec"`
	// +optional
	Status ImageRegistryStatus `json:"status,omitempty"`
}

// ImageRegistrySpec defines the specs for the running registry.
type ImageRegistrySpec struct {
	// operatorSpec allows operator specific configuration to be made.
	operatorv1.OperatorSpec `json:",inline"`
	// httpSecret is the value needed by the registry to secure uploads, generated by default.
	// +optional
	HTTPSecret string `json:"httpSecret,omitempty"`
	// proxy defines the proxy to be used when calling master api, upstream
	// registries, etc.
	// +optional
	Proxy ImageRegistryConfigProxy `json:"proxy,omitempty"`
	// storage details for configuring registry storage, e.g. S3 bucket
	// coordinates.
	// +optional
	Storage ImageRegistryConfigStorage `json:"storage,omitempty"`
	// readOnly indicates whether the registry instance should reject attempts
	// to push new images or delete existing ones.
	// +optional
	ReadOnly bool `json:"readOnly,omitempty"`
	// disableRedirect controls whether to route all data through the Registry,
	// rather than redirecting to the backend.
	// +optional
	DisableRedirect bool `json:"disableRedirect,omitempty"`
	// requests controls how many parallel requests a given registry instance
	// will handle before queuing additional requests.
	// +optional
	// +structType=atomic
	Requests ImageRegistryConfigRequests `json:"requests,omitempty"`
	// defaultRoute indicates whether an external facing route for the registry
	// should be created using the default generated hostname.
	// +optional
	DefaultRoute bool `json:"defaultRoute,omitempty"`
	// routes defines additional external facing routes which should be
	// created for the registry.
	// +optional
	// +listType=atomic
	Routes []ImageRegistryConfigRoute `json:"routes,omitempty"`
	// replicas determines the number of registry instances to run.
	Replicas int32 `json:"replicas"`
	// logging is deprecated, use logLevel instead.
	// +optional
	Logging int64 `json:"logging,omitempty"`
	// resources defines the resource requests+limits for the registry pod.
	// +optional
	// +structType=atomic
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
	// nodeSelector defines the node selection constraints for the registry
	// pod.
	// +optional
	// +mapType=atomic
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// tolerations defines the tolerations for the registry pod.
	// +optional
	// +listType=atomic
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
	// rolloutStrategy defines rollout strategy for the image registry
	// deployment.
	// +optional
	// +kubebuilder:validation:Pattern=`^(RollingUpdate|Recreate)$`
	RolloutStrategy string `json:"rolloutStrategy,omitempty"`
	// affinity is a group of node affinity scheduling rules for the image registry pod(s).
	// +optional
	// +structType=atomic
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
	// topologySpreadConstraints specify how to spread matching pods among the given topology.
	// +optional
	// +listType=atomic
	TopologySpreadConstraints []corev1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
}

// ImageRegistryStatus reports image registry operational status.
type ImageRegistryStatus struct {
	operatorv1.OperatorStatus `json:",inline"`

	// storageManaged is deprecated, please refer to Storage.managementState
	StorageManaged bool `json:"storageManaged"`
	// storage indicates the current applied storage configuration of the
	// registry.
	Storage ImageRegistryConfigStorage `json:"storage"`
}

// ImageRegistryConfigProxy defines proxy configuration to be used by registry.
type ImageRegistryConfigProxy struct {
	// http defines the proxy to be used by the image registry when
	// accessing HTTP endpoints.
	// +optional
	HTTP string `json:"http,omitempty"`
	// https defines the proxy to be used by the image registry when
	// accessing HTTPS endpoints.
	// +optional
	HTTPS string `json:"https,omitempty"`
	// noProxy defines a comma-separated list of host names that shouldn't
	// go through any proxy.
	// +optional
	NoProxy string `json:"noProxy,omitempty"`
}

// ImageRegistryConfigStorageS3CloudFront holds the configuration
// to use Amazon Cloudfront as the storage middleware in a registry.
// https://docs.docker.com/registry/configuration/#cloudfront
type ImageRegistryConfigStorageS3CloudFront struct {
	// baseURL contains the SCHEME://HOST[/PATH] at which Cloudfront is served.
	BaseURL string `json:"baseURL"`
	// privateKey points to secret containing the private key, provided by AWS.
	PrivateKey corev1.SecretKeySelector `json:"privateKey"`
	// keypairID is key pair ID provided by AWS.
	KeypairID string `json:"keypairID"`
	// duration is the duration of the Cloudfront session.
	// +optional
	// +kubebuilder:validation:Format=duration
	Duration metav1.Duration `json:"duration,omitempty"`
}

// ImageRegistryConfigStorageEmptyDir is an place holder to be used when
// when registry is leveraging ephemeral storage.
type ImageRegistryConfigStorageEmptyDir struct{}

// S3TrustedCASource references a config map with a CA certificate bundle in
// the "openshift-config" namespace. The key for the bundle in the
// config map is "ca-bundle.crt".
type S3TrustedCASource struct {
	// name is the metadata.name of the referenced config map.
	// This field must adhere to standard config map naming restrictions.
	// The name must consist solely of alphanumeric characters, hyphens (-)
	// and periods (.). It has a maximum length of 253 characters.
	// If this field is not specified or is empty string, the default trust
	// bundle will be used.
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=`^$|^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`
	// +optional
	Name string `json:"name"`
}

// ImageRegistryConfigStorageS3 holds the information to configure
// the registry to use the AWS S3 service for backend storage
// https://docs.docker.com/registry/storage-drivers/s3/
type ImageRegistryConfigStorageS3 struct {
	// bucket is the bucket name in which you want to store the registry's
	// data.
	// Optional, will be generated if not provided.
	// +optional
	Bucket string `json:"bucket,omitempty"`
	// region is the AWS region in which your bucket exists.
	// Optional, will be set based on the installed AWS Region.
	// +optional
	Region string `json:"region,omitempty"`
	// regionEndpoint is the endpoint for S3 compatible storage services.
	// It should be a valid URL with scheme, e.g. https://s3.example.com.
	// Optional, defaults based on the Region that is provided.
	// +optional
	RegionEndpoint string `json:"regionEndpoint,omitempty"`
	// encrypt specifies whether the registry stores the image in encrypted
	// format or not.
	// Optional, defaults to false.
	// +optional
	Encrypt bool `json:"encrypt,omitempty"`
	// keyID is the KMS key ID to use for encryption.
	// Optional, Encrypt must be true, or this parameter is ignored.
	// +optional
	KeyID string `json:"keyID,omitempty"`
	// cloudFront configures Amazon Cloudfront as the storage middleware in a
	// registry.
	// +optional
	// +structType=atomic
	CloudFront *ImageRegistryConfigStorageS3CloudFront `json:"cloudFront,omitempty"`
	// virtualHostedStyle enables using S3 virtual hosted style bucket paths with
	// a custom RegionEndpoint
	// Optional, defaults to false.
	// +optional
	VirtualHostedStyle bool `json:"virtualHostedStyle"`
	// trustedCA is a reference to a config map containing a CA bundle. The
	// image registry and its operator use certificates from this bundle to
	// verify S3 server certificates.
	//
	// The namespace for the config map referenced by trustedCA is
	// "openshift-config". The key for the bundle in the config map is
	// "ca-bundle.crt".
	// +optional
	TrustedCA S3TrustedCASource `json:"trustedCA"`
}

// ImageRegistryConfigStorageGCS holds GCS configuration.
type ImageRegistryConfigStorageGCS struct {
	// bucket is the bucket name in which you want to store the registry's
	// data.
	// Optional, will be generated if not provided.
	// +optional
	Bucket string `json:"bucket,omitempty"`
	// region is the GCS location in which your bucket exists.
	// Optional, will be set based on the installed GCS Region.
	// +optional
	Region string `json:"region,omitempty"`
	// projectID is the Project ID of the GCP project that this bucket should
	// be associated with.
	// +optional
	ProjectID string `json:"projectID,omitempty"`
	// keyID is the KMS key ID to use for encryption.
	// Optional, buckets are encrypted by default on GCP.
	// This allows for the use of a custom encryption key.
	// +optional
	KeyID string `json:"keyID,omitempty"`
}

// ImageRegistryConfigStorageSwift holds the information to configure
// the registry to use the OpenStack Swift service for backend storage
// https://docs.docker.com/registry/storage-drivers/swift/
type ImageRegistryConfigStorageSwift struct {
	// authURL defines the URL for obtaining an authentication token.
	// +optional
	AuthURL string `json:"authURL,omitempty"`
	// authVersion specifies the OpenStack Auth's version.
	// +optional
	AuthVersion string `json:"authVersion,omitempty"`
	// container defines the name of Swift container where to store the
	// registry's data.
	// +optional
	Container string `json:"container,omitempty"`
	// domain specifies Openstack's domain name for Identity v3 API.
	// +optional
	Domain string `json:"domain,omitempty"`
	// domainID specifies Openstack's domain id for Identity v3 API.
	// +optional
	DomainID string `json:"domainID,omitempty"`
	// tenant defines Openstack tenant name to be used by registry.
	// +optional
	Tenant string `json:"tenant,omitempty"`
	// tenant defines Openstack tenant id to be used by registry.
	// +optional
	TenantID string `json:"tenantID,omitempty"`
	// regionName defines Openstack's region in which container exists.
	// +optional
	RegionName string `json:"regionName,omitempty"`
}

// ImageRegistryConfigStoragePVC holds Persistent Volume Claims data to
// be used by the registry.
type ImageRegistryConfigStoragePVC struct {
	// claim defines the Persisent Volume Claim's name to be used.
	// +optional
	Claim string `json:"claim,omitempty"`
}

// ImageRegistryConfigStorageAzure holds the information to configure
// the registry to use Azure Blob Storage for backend storage.
type ImageRegistryConfigStorageAzure struct {
	// accountName defines the account to be used by the registry.
	// +optional
	AccountName string `json:"accountName,omitempty"`
	// container defines Azure's container to be used by registry.
	// +optional
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:MinLength=3
	// +kubebuilder:validation:Pattern=`^[0-9a-z]+(-[0-9a-z]+)*$`
	Container string `json:"container,omitempty"`
	// cloudName is the name of the Azure cloud environment to be used by the
	// registry. If empty, the operator will set it based on the infrastructure
	// object.
	// +optional
	CloudName string `json:"cloudName,omitempty"`
	// networkAccess defines the network access properties for the storage account.
	// Defaults to type: External.
	// +kubebuilder:default={"type": "External"}
	// +optional
	NetworkAccess *AzureNetworkAccess `json:"networkAccess,omitempty"`
}

// AzureNetworkAccess defines the network access properties for the storage account.
// +kubebuilder:validation:XValidation:rule="has(self.type) && self.type == 'Internal' ?  true : !has(self.internal)",message="internal is forbidden when type is not Internal"
// +union
type AzureNetworkAccess struct {
	// type is the network access level to be used for the storage account.
	// type: Internal means the storage account will be private, type: External
	// means the storage account will be publicly accessible.
	// Internal storage accounts are only exposed within the cluster's vnet.
	// External storage accounts are publicly exposed on the internet.
	// When type: Internal is used, a vnetName, subNetName and privateEndpointName
	// may optionally be specified. If unspecificed, the image registry operator
	// will discover vnet and subnet names, and generate a privateEndpointName.
	// Defaults to "External".
	// +kubebuilder:default:="External"
	// +unionDiscriminator
	// +optional
	Type AzureNetworkAccessType `json:"type,omitempty"`
	// internal defines the vnet and subnet names to configure a private
	// endpoint and connect it to the storage account in order to make it
	// private.
	// when type: Internal and internal is unset, the image registry operator
	// will discover vnet and subnet names, and generate a private endpoint
	// name.
	// +optional
	Internal *AzureNetworkAccessInternal `json:"internal,omitempty"`
}

type AzureNetworkAccessInternal struct {
	// networkResourceGroupName is the resource group name where the cluster's vnet
	// and subnet are. When omitted, the registry operator will use the cluster
	// resource group (from in the infrastructure status).
	// If you set a networkResourceGroupName on your install-config.yaml, that
	// value will be used automatically (for clusters configured with publish:Internal).
	// Note that both vnet and subnet must be in the same resource group.
	// It must be between 1 and 90 characters in length and must consist only of
	// alphanumeric characters, hyphens (-), periods (.) and underscores (_), and
	// not end with a period.
	// +kubebuilder:validation:MaxLength=90
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Pattern=`^[0-9A-Za-z_.-](?:[0-9A-Za-z_.-]*[0-9A-Za-z_-])?$`
	// +optional
	NetworkResourceGroupName string `json:"networkResourceGroupName,omitempty"`
	// vnetName is the name of the vnet the registry operates in. When omitted,
	// the registry operator will discover and set this by using the `kubernetes.io_cluster.<cluster-id>`
	// tag in the vnet resource. This tag is set automatically by the installer.
	// Commonly, this will be the same vnet as the cluster.
	// Advanced cluster network configurations should ensure the provided vnetName
	// is the vnet of the nodes where the image registry pods are running from.
	// It must be between 2 and 64 characters in length and must consist only of
	// alphanumeric characters, hyphens (-), periods (.) and underscores (_).
	// It must start with an alphanumeric character and end with an alphanumeric character or an underscore.
	// +kubebuilder:validation:MaxLength=64
	// +kubebuilder:validation:MinLength=2
	// +kubebuilder:validation:Pattern=`^[0-9A-Za-z][0-9A-Za-z_.-]*[0-9A-Za-z_]$`
	// +optional
	VNetName string `json:"vnetName,omitempty"`
	// subnetName is the name of the subnet the registry operates in. When omitted,
	// the registry operator will discover and set this by using the `kubernetes.io_cluster.<cluster-id>`
	// tag in the vnet resource, then using one of listed subnets.
	// Advanced cluster network configurations that use network security groups
	// to protect subnets should ensure the provided subnetName has access to
	// Azure Storage service.
	// It must be between 1 and 80 characters in length and must consist only of
	// alphanumeric characters, hyphens (-), periods (.) and underscores (_).
	// +kubebuilder:validation:MaxLength=80
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Pattern=`^[0-9A-Za-z](?:[0-9A-Za-z_.-]*[0-9A-Za-z_])?$`
	// +optional
	SubnetName string `json:"subnetName,omitempty"`
	// privateEndpointName is the name of the private endpoint for the registry.
	// When provided, the registry will use it as the name of the private endpoint
	// it will create for the storage account. When omitted, the registry will
	// generate one.
	// It must be between 2 and 64 characters in length and must consist only of
	// alphanumeric characters, hyphens (-), periods (.) and underscores (_).
	// It must start with an alphanumeric character and end with an alphanumeric character or an underscore.
	// +kubebuilder:validation:MaxLength=64
	// +kubebuilder:validation:MinLength=2
	// +kubebuilder:validation:Pattern=`^[0-9A-Za-z][0-9A-Za-z_.-]*[0-9A-Za-z_]$`
	// +optional
	PrivateEndpointName string `json:"privateEndpointName,omitempty"`
}

// AzureNetworkAccessType is the network access level to be used for the storage account.
// +kubebuilder:validation:Enum:="Internal";"External"
type AzureNetworkAccessType string

const (
	// AzureNetworkAccessTypeInternal means the storage account will be private
	AzureNetworkAccessTypeInternal AzureNetworkAccessType = "Internal"
	// AzureNetworkAccessTypeExternal means the storage account will be publicly accessible
	AzureNetworkAccessTypeExternal AzureNetworkAccessType = "External"
)

// ImageRegistryConfigStorageIBMCOS holds the information to configure
// the registry to use IBM Cloud Object Storage for backend storage.
type ImageRegistryConfigStorageIBMCOS struct {
	// bucket is the bucket name in which you want to store the registry's
	// data.
	// Optional, will be generated if not provided.
	// +optional
	Bucket string `json:"bucket,omitempty"`
	// location is the IBM Cloud location in which your bucket exists.
	// Optional, will be set based on the installed IBM Cloud location.
	// +optional
	Location string `json:"location,omitempty"`
	// resourceGroupName is the name of the IBM Cloud resource group that this
	// bucket and its service instance is associated with.
	// Optional, will be set based on the installed IBM Cloud resource group.
	// +optional
	ResourceGroupName string `json:"resourceGroupName,omitempty"`
	// resourceKeyCRN is the CRN of the IBM Cloud resource key that is created
	// for the service instance. Commonly referred as a service credential and
	// must contain HMAC type credentials.
	// Optional, will be computed if not provided.
	// +optional
	// +kubebuilder:validation:Pattern=`^crn:.+:.+:.+:cloud-object-storage:.+:.+:.+:resource-key:.+$`
	ResourceKeyCRN string `json:"resourceKeyCRN,omitempty"`
	// serviceInstanceCRN is the CRN of the IBM Cloud Object Storage service
	// instance that this bucket is associated with.
	// Optional, will be computed if not provided.
	// +optional
	// +kubebuilder:validation:Pattern=`^crn:.+:.+:.+:cloud-object-storage:.+:.+:.+::$`
	ServiceInstanceCRN string `json:"serviceInstanceCRN,omitempty"`
}

// EndpointAccessibility defines the Alibaba VPC endpoint for storage
type EndpointAccessibility string

// AlibabaEncryptionMethod defines an enumerable type for the encryption mode
type AlibabaEncryptionMethod string

const (
	// InternalEndpoint sets the VPC endpoint to internal
	InternalEndpoint EndpointAccessibility = "Internal"
	// PublicEndpoint sets the VPC endpoint to public
	PublicEndpoint EndpointAccessibility = "Public"

	// AES256 is an AlibabaEncryptionMethod. This means AES256 encryption
	AES256 AlibabaEncryptionMethod = "AES256"
	// KMS is an AlibabaEncryptionMethod. This means KMS encryption
	KMS AlibabaEncryptionMethod = "KMS"
)

// EncryptionAlibaba this a union type in kube parlance.  Depending on the value for the AlibabaEncryptionMethod,
// different pointers may be used
type EncryptionAlibaba struct {
	// Method defines the different encrytion modes available
	// Empty value means no opinion and the platform chooses the a default, which is subject to change over time.
	// Currently the default is `AES256`.
	// +kubebuilder:validation:Enum="KMS";"AES256"
	// +kubebuilder:default="AES256"
	// +optional
	Method AlibabaEncryptionMethod `json:"method"`

	// KMS (key management service) is an encryption type that holds the struct for KMS KeyID
	// +optional
	KMS *KMSEncryptionAlibaba `json:"kms,omitempty"`
}

type KMSEncryptionAlibaba struct {
	// KeyID holds the KMS encryption key ID
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	KeyID string `json:"keyID"`
}

// ImageRegistryConfigStorageAlibabaOSS holds Alibaba Cloud OSS configuration.
// Configures the registry to use Alibaba Cloud Object Storage Service for backend storage.
// More about oss, you can look at the [official documentation](https://www.alibabacloud.com/help/product/31815.htm)
type ImageRegistryConfigStorageAlibabaOSS struct {
	// Bucket is the bucket name in which you want to store the registry's data.
	// About Bucket naming, more details you can look at the [official documentation](https://www.alibabacloud.com/help/doc-detail/257087.htm)
	// Empty value means no opinion and the platform chooses the a default, which is subject to change over time.
	// Currently the default will be autogenerated in the form of <clusterid>-image-registry-<region>-<random string 27 chars>
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:MinLength=3
	// +kubebuilder:validation:Pattern=`^[0-9a-z]+(-[0-9a-z]+)*$`
	// +optional
	Bucket string `json:"bucket,omitempty"`
	// Region is the Alibaba Cloud Region in which your bucket exists.
	// For a list of regions, you can look at the [official documentation](https://www.alibabacloud.com/help/doc-detail/31837.html).
	// Empty value means no opinion and the platform chooses the a default, which is subject to change over time.
	// Currently the default will be based on the installed Alibaba Cloud Region.
	// +optional
	Region string `json:"region,omitempty"`
	// EndpointAccessibility specifies whether the registry use the OSS VPC internal endpoint
	// Empty value means no opinion and the platform chooses the a default, which is subject to change over time.
	// Currently the default is `Internal`.
	// +kubebuilder:validation:Enum="Internal";"Public";""
	// +kubebuilder:default="Internal"
	// +optional
	EndpointAccessibility EndpointAccessibility `json:"endpointAccessibility,omitempty"`
	// Encryption specifies whether you would like your data encrypted on the server side.
	// More details, you can look cat the [official documentation](https://www.alibabacloud.com/help/doc-detail/117914.htm)
	// +optional
	Encryption *EncryptionAlibaba `json:"encryption,omitempty"`
}

// ImageRegistryConfigStorage describes how the storage should be configured
// for the image registry.
type ImageRegistryConfigStorage struct {
	// emptyDir represents ephemeral storage on the pod's host node.
	// WARNING: this storage cannot be used with more than 1 replica and
	// is not suitable for production use. When the pod is removed from a
	// node for any reason, the data in the emptyDir is deleted forever.
	// +optional
	EmptyDir *ImageRegistryConfigStorageEmptyDir `json:"emptyDir,omitempty"`
	// s3 represents configuration that uses Amazon Simple Storage Service.
	// +optional
	S3 *ImageRegistryConfigStorageS3 `json:"s3,omitempty"`
	// gcs represents configuration that uses Google Cloud Storage.
	// +optional
	GCS *ImageRegistryConfigStorageGCS `json:"gcs,omitempty"`
	// swift represents configuration that uses OpenStack Object Storage.
	// +optional
	Swift *ImageRegistryConfigStorageSwift `json:"swift,omitempty"`
	// pvc represents configuration that uses a PersistentVolumeClaim.
	// +optional
	PVC *ImageRegistryConfigStoragePVC `json:"pvc,omitempty"`
	// azure represents configuration that uses Azure Blob Storage.
	// +optional
	Azure *ImageRegistryConfigStorageAzure `json:"azure,omitempty"`
	// ibmcos represents configuration that uses IBM Cloud Object Storage.
	// +optional
	IBMCOS *ImageRegistryConfigStorageIBMCOS `json:"ibmcos,omitempty"`
	// Oss represents configuration that uses Alibaba Cloud Object Storage Service.
	// +optional
	OSS *ImageRegistryConfigStorageAlibabaOSS `json:"oss,omitempty"`
	// managementState indicates if the operator manages the underlying
	// storage unit. If Managed the operator will remove the storage when
	// this operator gets Removed.
	// +optional
	// +kubebuilder:validation:Pattern=`^(Managed|Unmanaged)$`
	ManagementState string `json:"managementState,omitempty"`
}

// ImageRegistryConfigRequests defines registry limits on requests read and write.
// +structType=atomic
type ImageRegistryConfigRequests struct {
	// read defines limits for image registry's reads.
	// +optional
	Read ImageRegistryConfigRequestsLimits `json:"read,omitempty"`
	// write defines limits for image registry's writes.
	// +optional
	Write ImageRegistryConfigRequestsLimits `json:"write,omitempty"`
}

// ImageRegistryConfigRequestsLimits holds configuration on the max, enqueued
// and waiting registry's API requests.
type ImageRegistryConfigRequestsLimits struct {
	// maxRunning sets the maximum in flight api requests to the registry.
	// +optional
	MaxRunning int `json:"maxRunning,omitempty"`
	// maxInQueue sets the maximum queued api requests to the registry.
	// +optional
	MaxInQueue int `json:"maxInQueue,omitempty"`
	// maxWaitInQueue sets the maximum time a request can wait in the queue
	// before being rejected.
	// +optional
	// +kubebuilder:validation:Format=duration
	MaxWaitInQueue metav1.Duration `json:"maxWaitInQueue,omitempty"`
}

// ImageRegistryConfigRoute holds information on external route access to image
// registry.
type ImageRegistryConfigRoute struct {
	// name of the route to be created.
	Name string `json:"name"`
	// hostname for the route.
	// +optional
	Hostname string `json:"hostname,omitempty"`
	// secretName points to secret containing the certificates to be used
	// by the route.
	// +optional
	SecretName string `json:"secretName,omitempty"`
}
