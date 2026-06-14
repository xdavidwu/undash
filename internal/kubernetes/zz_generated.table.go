package kubernetes

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	admissionregistrationv1 "k8s.io/kubernetes/pkg/apis/admissionregistration/v1"
	admissionregistrationv1alpha1 "k8s.io/kubernetes/pkg/apis/admissionregistration/v1alpha1"
	admissionregistrationv1beta1 "k8s.io/kubernetes/pkg/apis/admissionregistration/v1beta1"
	apidiscoveryv2 "k8s.io/kubernetes/pkg/apis/apidiscovery/v2"
	apidiscoveryv2beta1 "k8s.io/kubernetes/pkg/apis/apidiscovery/v2beta1"
	apiserverinternalv1alpha1 "k8s.io/kubernetes/pkg/apis/apiserverinternal/v1alpha1"
	appsv1 "k8s.io/kubernetes/pkg/apis/apps/v1"
	appsv1beta1 "k8s.io/kubernetes/pkg/apis/apps/v1beta1"
	appsv1beta2 "k8s.io/kubernetes/pkg/apis/apps/v1beta2"
	authenticationv1 "k8s.io/kubernetes/pkg/apis/authentication/v1"
	authenticationv1alpha1 "k8s.io/kubernetes/pkg/apis/authentication/v1alpha1"
	authenticationv1beta1 "k8s.io/kubernetes/pkg/apis/authentication/v1beta1"
	authorizationv1 "k8s.io/kubernetes/pkg/apis/authorization/v1"
	authorizationv1beta1 "k8s.io/kubernetes/pkg/apis/authorization/v1beta1"
	autoscalingv1 "k8s.io/kubernetes/pkg/apis/autoscaling/v1"
	autoscalingv2 "k8s.io/kubernetes/pkg/apis/autoscaling/v2"
	autoscalingv2beta1 "k8s.io/kubernetes/pkg/apis/autoscaling/v2beta1"
	autoscalingv2beta2 "k8s.io/kubernetes/pkg/apis/autoscaling/v2beta2"
	batchv1 "k8s.io/kubernetes/pkg/apis/batch/v1"
	batchv1beta1 "k8s.io/kubernetes/pkg/apis/batch/v1beta1"
	certificatesv1 "k8s.io/kubernetes/pkg/apis/certificates/v1"
	certificatesv1alpha1 "k8s.io/kubernetes/pkg/apis/certificates/v1alpha1"
	certificatesv1beta1 "k8s.io/kubernetes/pkg/apis/certificates/v1beta1"
	coordinationv1 "k8s.io/kubernetes/pkg/apis/coordination/v1"
	coordinationv1beta1 "k8s.io/kubernetes/pkg/apis/coordination/v1beta1"
	corev1 "k8s.io/kubernetes/pkg/apis/core/v1"
	discoveryv1 "k8s.io/kubernetes/pkg/apis/discovery/v1"
	discoveryv1beta1 "k8s.io/kubernetes/pkg/apis/discovery/v1beta1"
	eventsv1 "k8s.io/kubernetes/pkg/apis/events/v1"
	eventsv1beta1 "k8s.io/kubernetes/pkg/apis/events/v1beta1"
	extensionsv1beta1 "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	flowcontrolv1 "k8s.io/kubernetes/pkg/apis/flowcontrol/v1"
	flowcontrolv1beta1 "k8s.io/kubernetes/pkg/apis/flowcontrol/v1beta1"
	flowcontrolv1beta2 "k8s.io/kubernetes/pkg/apis/flowcontrol/v1beta2"
	flowcontrolv1beta3 "k8s.io/kubernetes/pkg/apis/flowcontrol/v1beta3"
	imagepolicyv1alpha1 "k8s.io/kubernetes/pkg/apis/imagepolicy/v1alpha1"
	networkingv1 "k8s.io/kubernetes/pkg/apis/networking/v1"
	networkingv1alpha1 "k8s.io/kubernetes/pkg/apis/networking/v1alpha1"
	networkingv1beta1 "k8s.io/kubernetes/pkg/apis/networking/v1beta1"
	nodev1 "k8s.io/kubernetes/pkg/apis/node/v1"
	nodev1alpha1 "k8s.io/kubernetes/pkg/apis/node/v1alpha1"
	nodev1beta1 "k8s.io/kubernetes/pkg/apis/node/v1beta1"
	policyv1 "k8s.io/kubernetes/pkg/apis/policy/v1"
	policyv1beta1 "k8s.io/kubernetes/pkg/apis/policy/v1beta1"
	rbacv1 "k8s.io/kubernetes/pkg/apis/rbac/v1"
	rbacv1alpha1 "k8s.io/kubernetes/pkg/apis/rbac/v1alpha1"
	rbacv1beta1 "k8s.io/kubernetes/pkg/apis/rbac/v1beta1"
	resourcev1alpha2 "k8s.io/kubernetes/pkg/apis/resource/v1alpha2"
	schedulingv1 "k8s.io/kubernetes/pkg/apis/scheduling/v1"
	schedulingv1alpha1 "k8s.io/kubernetes/pkg/apis/scheduling/v1alpha1"
	schedulingv1beta1 "k8s.io/kubernetes/pkg/apis/scheduling/v1beta1"
	storagev1 "k8s.io/kubernetes/pkg/apis/storage/v1"
	storagev1alpha1 "k8s.io/kubernetes/pkg/apis/storage/v1alpha1"
	storagev1beta1 "k8s.io/kubernetes/pkg/apis/storage/v1beta1"
	storagemigrationv1alpha1 "k8s.io/kubernetes/pkg/apis/storagemigration/v1alpha1"
)

// apiextensionsv1 objects to metav1.Table functions
var ApiextensionsV1CustomResourceDefinitionToTable = ToTableFunc(apiextensionsv1.Convert_v1_CustomResourceDefinition_To_apiextensions_CustomResourceDefinition)

// apiextensionsv1 lists to metav1.Table functions
var ApiextensionsV1CustomResourceDefinitionListToTable = ToTableFunc(apiextensionsv1.Convert_v1_CustomResourceDefinitionList_To_apiextensions_CustomResourceDefinitionList)

// apiextensionsv1beta1 objects to metav1.Table functions
var ApiextensionsV1beta1CustomResourceDefinitionToTable = ToTableFunc(apiextensionsv1beta1.Convert_v1beta1_CustomResourceDefinition_To_apiextensions_CustomResourceDefinition)

// apiextensionsv1beta1 lists to metav1.Table functions
var ApiextensionsV1beta1CustomResourceDefinitionListToTable = ToTableFunc(apiextensionsv1beta1.Convert_v1beta1_CustomResourceDefinitionList_To_apiextensions_CustomResourceDefinitionList)

// authenticationv1 objects to metav1.Table functions
var (
	AuthenticationV1SelfSubjectReviewToTable = ToTableFunc(authenticationv1.Convert_v1_SelfSubjectReview_To_authentication_SelfSubjectReview)
	AuthenticationV1TokenRequestToTable      = ToTableFunc(authenticationv1.Convert_v1_TokenRequest_To_authentication_TokenRequest)
	AuthenticationV1TokenReviewToTable       = ToTableFunc(authenticationv1.Convert_v1_TokenReview_To_authentication_TokenReview)
)

// admissionregistrationv1 objects to metav1.Table functions
var (
	AdmissionregistrationV1MutatingWebhookConfigurationToTable     = ToTableFunc(admissionregistrationv1.Convert_v1_MutatingWebhookConfiguration_To_admissionregistration_MutatingWebhookConfiguration)
	AdmissionregistrationV1ValidatingAdmissionPolicyBindingToTable = ToTableFunc(admissionregistrationv1.Convert_v1_ValidatingAdmissionPolicyBinding_To_admissionregistration_ValidatingAdmissionPolicyBinding)
	AdmissionregistrationV1ValidatingAdmissionPolicyToTable        = ToTableFunc(admissionregistrationv1.Convert_v1_ValidatingAdmissionPolicy_To_admissionregistration_ValidatingAdmissionPolicy)
	AdmissionregistrationV1ValidatingWebhookConfigurationToTable   = ToTableFunc(admissionregistrationv1.Convert_v1_ValidatingWebhookConfiguration_To_admissionregistration_ValidatingWebhookConfiguration)
)

// admissionregistrationv1 lists to metav1.Table functions
var (
	AdmissionregistrationV1MutatingWebhookConfigurationListToTable     = ToTableFunc(admissionregistrationv1.Convert_v1_MutatingWebhookConfigurationList_To_admissionregistration_MutatingWebhookConfigurationList)
	AdmissionregistrationV1ValidatingAdmissionPolicyBindingListToTable = ToTableFunc(admissionregistrationv1.Convert_v1_ValidatingAdmissionPolicyBindingList_To_admissionregistration_ValidatingAdmissionPolicyBindingList)
	AdmissionregistrationV1ValidatingAdmissionPolicyListToTable        = ToTableFunc(admissionregistrationv1.Convert_v1_ValidatingAdmissionPolicyList_To_admissionregistration_ValidatingAdmissionPolicyList)
	AdmissionregistrationV1ValidatingWebhookConfigurationListToTable   = ToTableFunc(admissionregistrationv1.Convert_v1_ValidatingWebhookConfigurationList_To_admissionregistration_ValidatingWebhookConfigurationList)
)

// admissionregistrationv1alpha1 objects to metav1.Table functions
var (
	AdmissionregistrationV1alpha1ValidatingAdmissionPolicyBindingToTable = ToTableFunc(admissionregistrationv1alpha1.Convert_v1alpha1_ValidatingAdmissionPolicyBinding_To_admissionregistration_ValidatingAdmissionPolicyBinding)
	AdmissionregistrationV1alpha1ValidatingAdmissionPolicyToTable        = ToTableFunc(admissionregistrationv1alpha1.Convert_v1alpha1_ValidatingAdmissionPolicy_To_admissionregistration_ValidatingAdmissionPolicy)
)

// admissionregistrationv1alpha1 lists to metav1.Table functions
var (
	AdmissionregistrationV1alpha1ValidatingAdmissionPolicyBindingListToTable = ToTableFunc(admissionregistrationv1alpha1.Convert_v1alpha1_ValidatingAdmissionPolicyBindingList_To_admissionregistration_ValidatingAdmissionPolicyBindingList)
	AdmissionregistrationV1alpha1ValidatingAdmissionPolicyListToTable        = ToTableFunc(admissionregistrationv1alpha1.Convert_v1alpha1_ValidatingAdmissionPolicyList_To_admissionregistration_ValidatingAdmissionPolicyList)
)

// admissionregistrationv1beta1 objects to metav1.Table functions
var (
	AdmissionregistrationV1beta1MutatingWebhookConfigurationToTable     = ToTableFunc(admissionregistrationv1beta1.Convert_v1beta1_MutatingWebhookConfiguration_To_admissionregistration_MutatingWebhookConfiguration)
	AdmissionregistrationV1beta1ValidatingAdmissionPolicyBindingToTable = ToTableFunc(admissionregistrationv1beta1.Convert_v1beta1_ValidatingAdmissionPolicyBinding_To_admissionregistration_ValidatingAdmissionPolicyBinding)
	AdmissionregistrationV1beta1ValidatingAdmissionPolicyToTable        = ToTableFunc(admissionregistrationv1beta1.Convert_v1beta1_ValidatingAdmissionPolicy_To_admissionregistration_ValidatingAdmissionPolicy)
	AdmissionregistrationV1beta1ValidatingWebhookConfigurationToTable   = ToTableFunc(admissionregistrationv1beta1.Convert_v1beta1_ValidatingWebhookConfiguration_To_admissionregistration_ValidatingWebhookConfiguration)
)

// admissionregistrationv1beta1 lists to metav1.Table functions
var (
	AdmissionregistrationV1beta1MutatingWebhookConfigurationListToTable     = ToTableFunc(admissionregistrationv1beta1.Convert_v1beta1_MutatingWebhookConfigurationList_To_admissionregistration_MutatingWebhookConfigurationList)
	AdmissionregistrationV1beta1ValidatingAdmissionPolicyBindingListToTable = ToTableFunc(admissionregistrationv1beta1.Convert_v1beta1_ValidatingAdmissionPolicyBindingList_To_admissionregistration_ValidatingAdmissionPolicyBindingList)
	AdmissionregistrationV1beta1ValidatingAdmissionPolicyListToTable        = ToTableFunc(admissionregistrationv1beta1.Convert_v1beta1_ValidatingAdmissionPolicyList_To_admissionregistration_ValidatingAdmissionPolicyList)
	AdmissionregistrationV1beta1ValidatingWebhookConfigurationListToTable   = ToTableFunc(admissionregistrationv1beta1.Convert_v1beta1_ValidatingWebhookConfigurationList_To_admissionregistration_ValidatingWebhookConfigurationList)
)

// apidiscoveryv2 objects to metav1.Table functions
var ApidiscoveryV2APIGroupDiscoveryToTable = ToTableFunc(apidiscoveryv2.Convert_v2_APIGroupDiscovery_To_apidiscovery_APIGroupDiscovery)

// apidiscoveryv2 lists to metav1.Table functions
var ApidiscoveryV2APIGroupDiscoveryListToTable = ToTableFunc(apidiscoveryv2.Convert_v2_APIGroupDiscoveryList_To_apidiscovery_APIGroupDiscoveryList)

// apidiscoveryv2beta1 objects to metav1.Table functions
var ApidiscoveryV2beta1APIGroupDiscoveryToTable = ToTableFunc(apidiscoveryv2beta1.Convert_v2beta1_APIGroupDiscovery_To_apidiscovery_APIGroupDiscovery)

// apidiscoveryv2beta1 lists to metav1.Table functions
var ApidiscoveryV2beta1APIGroupDiscoveryListToTable = ToTableFunc(apidiscoveryv2beta1.Convert_v2beta1_APIGroupDiscoveryList_To_apidiscovery_APIGroupDiscoveryList)

// apiserverinternalv1alpha1 objects to metav1.Table functions
var ApiserverinternalV1alpha1StorageVersionToTable = ToTableFunc(apiserverinternalv1alpha1.Convert_v1alpha1_StorageVersion_To_apiserverinternal_StorageVersion)

// apiserverinternalv1alpha1 lists to metav1.Table functions
var ApiserverinternalV1alpha1StorageVersionListToTable = ToTableFunc(apiserverinternalv1alpha1.Convert_v1alpha1_StorageVersionList_To_apiserverinternal_StorageVersionList)

// corev1 objects to metav1.Table functions
var (
	CoreV1BindingToTable               = ToTableFunc(corev1.Convert_v1_Binding_To_core_Binding)
	CoreV1ComponentStatusToTable       = ToTableFunc(corev1.Convert_v1_ComponentStatus_To_core_ComponentStatus)
	CoreV1ConfigMapToTable             = ToTableFunc(corev1.Convert_v1_ConfigMap_To_core_ConfigMap)
	CoreV1EndpointsToTable             = ToTableFunc(corev1.Convert_v1_Endpoints_To_core_Endpoints)
	CoreV1EventToTable                 = ToTableFunc(corev1.Convert_v1_Event_To_core_Event)
	CoreV1LimitRangeToTable            = ToTableFunc(corev1.Convert_v1_LimitRange_To_core_LimitRange)
	CoreV1NamespaceToTable             = ToTableFunc(corev1.Convert_v1_Namespace_To_core_Namespace)
	CoreV1NodeToTable                  = ToTableFunc(corev1.Convert_v1_Node_To_core_Node)
	CoreV1PersistentVolumeClaimToTable = ToTableFunc(corev1.Convert_v1_PersistentVolumeClaim_To_core_PersistentVolumeClaim)
	CoreV1PersistentVolumeToTable      = ToTableFunc(corev1.Convert_v1_PersistentVolume_To_core_PersistentVolume)
	CoreV1PodStatusResultToTable       = ToTableFunc(corev1.Convert_v1_PodStatusResult_To_core_PodStatusResult)
	CoreV1PodTemplateToTable           = ToTableFunc(corev1.Convert_v1_PodTemplate_To_core_PodTemplate)
	CoreV1PodToTable                   = ToTableFunc(corev1.Convert_v1_Pod_To_core_Pod)
	CoreV1RangeAllocationToTable       = ToTableFunc(corev1.Convert_v1_RangeAllocation_To_core_RangeAllocation)
	CoreV1ReplicationControllerToTable = ToTableFunc(corev1.Convert_v1_ReplicationController_To_apps_ReplicaSet)
	CoreV1ResourceQuotaToTable         = ToTableFunc(corev1.Convert_v1_ResourceQuota_To_core_ResourceQuota)
	CoreV1SecretToTable                = ToTableFunc(corev1.Convert_v1_Secret_To_core_Secret)
	CoreV1ServiceAccountToTable        = ToTableFunc(corev1.Convert_v1_ServiceAccount_To_core_ServiceAccount)
	CoreV1ServiceToTable               = ToTableFunc(corev1.Convert_v1_Service_To_core_Service)
)

// corev1 lists to metav1.Table functions
var (
	CoreV1ComponentStatusListToTable       = ToTableFunc(corev1.Convert_v1_ComponentStatusList_To_core_ComponentStatusList)
	CoreV1ConfigMapListToTable             = ToTableFunc(corev1.Convert_v1_ConfigMapList_To_core_ConfigMapList)
	CoreV1EndpointsListToTable             = ToTableFunc(corev1.Convert_v1_EndpointsList_To_core_EndpointsList)
	CoreV1EventListToTable                 = ToTableFunc(corev1.Convert_v1_EventList_To_core_EventList)
	CoreV1LimitRangeListToTable            = ToTableFunc(corev1.Convert_v1_LimitRangeList_To_core_LimitRangeList)
	CoreV1ListToTable                      = ToTableFunc(corev1.Convert_v1_List_To_core_List)
	CoreV1NamespaceListToTable             = ToTableFunc(corev1.Convert_v1_NamespaceList_To_core_NamespaceList)
	CoreV1NodeListToTable                  = ToTableFunc(corev1.Convert_v1_NodeList_To_core_NodeList)
	CoreV1PersistentVolumeClaimListToTable = ToTableFunc(corev1.Convert_v1_PersistentVolumeClaimList_To_core_PersistentVolumeClaimList)
	CoreV1PersistentVolumeListToTable      = ToTableFunc(corev1.Convert_v1_PersistentVolumeList_To_core_PersistentVolumeList)
	CoreV1PodListToTable                   = ToTableFunc(corev1.Convert_v1_PodList_To_core_PodList)
	CoreV1PodTemplateListToTable           = ToTableFunc(corev1.Convert_v1_PodTemplateList_To_core_PodTemplateList)
	CoreV1ReplicationControllerListToTable = ToTableFunc(corev1.Convert_v1_ReplicationControllerList_To_core_ReplicationControllerList)
	CoreV1ResourceQuotaListToTable         = ToTableFunc(corev1.Convert_v1_ResourceQuotaList_To_core_ResourceQuotaList)
	CoreV1SecretListToTable                = ToTableFunc(corev1.Convert_v1_SecretList_To_core_SecretList)
	CoreV1ServiceAccountListToTable        = ToTableFunc(corev1.Convert_v1_ServiceAccountList_To_core_ServiceAccountList)
	CoreV1ServiceListToTable               = ToTableFunc(corev1.Convert_v1_ServiceList_To_core_ServiceList)
)

// appsv1 objects to metav1.Table functions
var (
	AppsV1ControllerRevisionToTable = ToTableFunc(appsv1.Convert_v1_ControllerRevision_To_apps_ControllerRevision)
	AppsV1DaemonSetToTable          = ToTableFunc(appsv1.Convert_v1_DaemonSet_To_apps_DaemonSet)
	AppsV1DeploymentToTable         = ToTableFunc(appsv1.Convert_v1_Deployment_To_apps_Deployment)
	AppsV1ReplicaSetToTable         = ToTableFunc(appsv1.Convert_v1_ReplicaSet_To_apps_ReplicaSet)
	AppsV1StatefulSetToTable        = ToTableFunc(appsv1.Convert_v1_StatefulSet_To_apps_StatefulSet)
)

// appsv1 lists to metav1.Table functions
var (
	AppsV1ControllerRevisionListToTable = ToTableFunc(appsv1.Convert_v1_ControllerRevisionList_To_apps_ControllerRevisionList)
	AppsV1DaemonSetListToTable          = ToTableFunc(appsv1.Convert_v1_DaemonSetList_To_apps_DaemonSetList)
	AppsV1DeploymentListToTable         = ToTableFunc(appsv1.Convert_v1_DeploymentList_To_apps_DeploymentList)
	AppsV1ReplicaSetListToTable         = ToTableFunc(appsv1.Convert_v1_ReplicaSetList_To_apps_ReplicaSetList)
	AppsV1StatefulSetListToTable        = ToTableFunc(appsv1.Convert_v1_StatefulSetList_To_apps_StatefulSetList)
)

// appsv1beta1 objects to metav1.Table functions
var (
	AppsV1beta1ControllerRevisionToTable = ToTableFunc(appsv1beta1.Convert_v1beta1_ControllerRevision_To_apps_ControllerRevision)
	AppsV1beta1DeploymentToTable         = ToTableFunc(appsv1beta1.Convert_v1beta1_Deployment_To_apps_Deployment)
	AppsV1beta1ScaleToTable              = ToTableFunc(appsv1beta1.Convert_v1beta1_Scale_To_autoscaling_Scale)
	AppsV1beta1StatefulSetToTable        = ToTableFunc(appsv1beta1.Convert_v1beta1_StatefulSet_To_apps_StatefulSet)
)

// appsv1beta1 lists to metav1.Table functions
var (
	AppsV1beta1ControllerRevisionListToTable = ToTableFunc(appsv1beta1.Convert_v1beta1_ControllerRevisionList_To_apps_ControllerRevisionList)
	AppsV1beta1DeploymentListToTable         = ToTableFunc(appsv1beta1.Convert_v1beta1_DeploymentList_To_apps_DeploymentList)
	AppsV1beta1StatefulSetListToTable        = ToTableFunc(appsv1beta1.Convert_v1beta1_StatefulSetList_To_apps_StatefulSetList)
)

// appsv1beta2 objects to metav1.Table functions
var (
	AppsV1beta2ControllerRevisionToTable = ToTableFunc(appsv1beta2.Convert_v1beta2_ControllerRevision_To_apps_ControllerRevision)
	AppsV1beta2DaemonSetToTable          = ToTableFunc(appsv1beta2.Convert_v1beta2_DaemonSet_To_apps_DaemonSet)
	AppsV1beta2DeploymentToTable         = ToTableFunc(appsv1beta2.Convert_v1beta2_Deployment_To_apps_Deployment)
	AppsV1beta2ReplicaSetToTable         = ToTableFunc(appsv1beta2.Convert_v1beta2_ReplicaSet_To_apps_ReplicaSet)
	AppsV1beta2ScaleToTable              = ToTableFunc(appsv1beta2.Convert_v1beta2_Scale_To_autoscaling_Scale)
	AppsV1beta2StatefulSetToTable        = ToTableFunc(appsv1beta2.Convert_v1beta2_StatefulSet_To_apps_StatefulSet)
)

// appsv1beta2 lists to metav1.Table functions
var (
	AppsV1beta2ControllerRevisionListToTable = ToTableFunc(appsv1beta2.Convert_v1beta2_ControllerRevisionList_To_apps_ControllerRevisionList)
	AppsV1beta2DaemonSetListToTable          = ToTableFunc(appsv1beta2.Convert_v1beta2_DaemonSetList_To_apps_DaemonSetList)
	AppsV1beta2DeploymentListToTable         = ToTableFunc(appsv1beta2.Convert_v1beta2_DeploymentList_To_apps_DeploymentList)
	AppsV1beta2ReplicaSetListToTable         = ToTableFunc(appsv1beta2.Convert_v1beta2_ReplicaSetList_To_apps_ReplicaSetList)
	AppsV1beta2StatefulSetListToTable        = ToTableFunc(appsv1beta2.Convert_v1beta2_StatefulSetList_To_apps_StatefulSetList)
)

// authenticationv1alpha1 objects to metav1.Table functions
var AuthenticationV1alpha1SelfSubjectReviewToTable = ToTableFunc(authenticationv1alpha1.Convert_v1alpha1_SelfSubjectReview_To_authentication_SelfSubjectReview)

// authenticationv1beta1 objects to metav1.Table functions
var (
	AuthenticationV1beta1SelfSubjectReviewToTable = ToTableFunc(authenticationv1beta1.Convert_v1beta1_SelfSubjectReview_To_authentication_SelfSubjectReview)
	AuthenticationV1beta1TokenReviewToTable       = ToTableFunc(authenticationv1beta1.Convert_v1beta1_TokenReview_To_authentication_TokenReview)
)

// authorizationv1 objects to metav1.Table functions
var (
	AuthorizationV1LocalSubjectAccessReviewToTable = ToTableFunc(authorizationv1.Convert_v1_LocalSubjectAccessReview_To_authorization_LocalSubjectAccessReview)
	AuthorizationV1SelfSubjectAccessReviewToTable  = ToTableFunc(authorizationv1.Convert_v1_SelfSubjectAccessReview_To_authorization_SelfSubjectAccessReview)
	AuthorizationV1SelfSubjectRulesReviewToTable   = ToTableFunc(authorizationv1.Convert_v1_SelfSubjectRulesReview_To_authorization_SelfSubjectRulesReview)
	AuthorizationV1SubjectAccessReviewToTable      = ToTableFunc(authorizationv1.Convert_v1_SubjectAccessReview_To_authorization_SubjectAccessReview)
)

// authorizationv1beta1 objects to metav1.Table functions
var (
	AuthorizationV1beta1LocalSubjectAccessReviewToTable = ToTableFunc(authorizationv1beta1.Convert_v1beta1_LocalSubjectAccessReview_To_authorization_LocalSubjectAccessReview)
	AuthorizationV1beta1SelfSubjectAccessReviewToTable  = ToTableFunc(authorizationv1beta1.Convert_v1beta1_SelfSubjectAccessReview_To_authorization_SelfSubjectAccessReview)
	AuthorizationV1beta1SelfSubjectRulesReviewToTable   = ToTableFunc(authorizationv1beta1.Convert_v1beta1_SelfSubjectRulesReview_To_authorization_SelfSubjectRulesReview)
	AuthorizationV1beta1SubjectAccessReviewToTable      = ToTableFunc(authorizationv1beta1.Convert_v1beta1_SubjectAccessReview_To_authorization_SubjectAccessReview)
)

// autoscalingv1 objects to metav1.Table functions
var (
	AutoscalingV1HorizontalPodAutoscalerToTable = ToTableFunc(autoscalingv1.Convert_v1_HorizontalPodAutoscaler_To_autoscaling_HorizontalPodAutoscaler)
	AutoscalingV1ScaleToTable                   = ToTableFunc(autoscalingv1.Convert_v1_Scale_To_autoscaling_Scale)
)

// autoscalingv1 lists to metav1.Table functions
var AutoscalingV1HorizontalPodAutoscalerListToTable = ToTableFunc(autoscalingv1.Convert_v1_HorizontalPodAutoscalerList_To_autoscaling_HorizontalPodAutoscalerList)

// autoscalingv2 objects to metav1.Table functions
var AutoscalingV2HorizontalPodAutoscalerToTable = ToTableFunc(autoscalingv2.Convert_v2_HorizontalPodAutoscaler_To_autoscaling_HorizontalPodAutoscaler)

// autoscalingv2 lists to metav1.Table functions
var AutoscalingV2HorizontalPodAutoscalerListToTable = ToTableFunc(autoscalingv2.Convert_v2_HorizontalPodAutoscalerList_To_autoscaling_HorizontalPodAutoscalerList)

// autoscalingv2beta1 objects to metav1.Table functions
var AutoscalingV2beta1HorizontalPodAutoscalerToTable = ToTableFunc(autoscalingv2beta1.Convert_v2beta1_HorizontalPodAutoscaler_To_autoscaling_HorizontalPodAutoscaler)

// autoscalingv2beta1 lists to metav1.Table functions
var AutoscalingV2beta1HorizontalPodAutoscalerListToTable = ToTableFunc(autoscalingv2beta1.Convert_v2beta1_HorizontalPodAutoscalerList_To_autoscaling_HorizontalPodAutoscalerList)

// autoscalingv2beta2 objects to metav1.Table functions
var AutoscalingV2beta2HorizontalPodAutoscalerToTable = ToTableFunc(autoscalingv2beta2.Convert_v2beta2_HorizontalPodAutoscaler_To_autoscaling_HorizontalPodAutoscaler)

// autoscalingv2beta2 lists to metav1.Table functions
var AutoscalingV2beta2HorizontalPodAutoscalerListToTable = ToTableFunc(autoscalingv2beta2.Convert_v2beta2_HorizontalPodAutoscalerList_To_autoscaling_HorizontalPodAutoscalerList)

// batchv1 objects to metav1.Table functions
var (
	BatchV1CronJobToTable = ToTableFunc(batchv1.Convert_v1_CronJob_To_batch_CronJob)
	BatchV1JobToTable     = ToTableFunc(batchv1.Convert_v1_Job_To_batch_Job)
)

// batchv1 lists to metav1.Table functions
var (
	BatchV1CronJobListToTable = ToTableFunc(batchv1.Convert_v1_CronJobList_To_batch_CronJobList)
	BatchV1JobListToTable     = ToTableFunc(batchv1.Convert_v1_JobList_To_batch_JobList)
)

// batchv1beta1 objects to metav1.Table functions
var BatchV1beta1CronJobToTable = ToTableFunc(batchv1beta1.Convert_v1beta1_CronJob_To_batch_CronJob)

// batchv1beta1 lists to metav1.Table functions
var BatchV1beta1CronJobListToTable = ToTableFunc(batchv1beta1.Convert_v1beta1_CronJobList_To_batch_CronJobList)

// certificatesv1 objects to metav1.Table functions
var CertificatesV1CertificateSigningRequestToTable = ToTableFunc(certificatesv1.Convert_v1_CertificateSigningRequest_To_certificates_CertificateSigningRequest)

// certificatesv1 lists to metav1.Table functions
var CertificatesV1CertificateSigningRequestListToTable = ToTableFunc(certificatesv1.Convert_v1_CertificateSigningRequestList_To_certificates_CertificateSigningRequestList)

// certificatesv1alpha1 objects to metav1.Table functions
var CertificatesV1alpha1ClusterTrustBundleToTable = ToTableFunc(certificatesv1alpha1.Convert_v1alpha1_ClusterTrustBundle_To_certificates_ClusterTrustBundle)

// certificatesv1alpha1 lists to metav1.Table functions
var CertificatesV1alpha1ClusterTrustBundleListToTable = ToTableFunc(certificatesv1alpha1.Convert_v1alpha1_ClusterTrustBundleList_To_certificates_ClusterTrustBundleList)

// certificatesv1beta1 objects to metav1.Table functions
var CertificatesV1beta1CertificateSigningRequestToTable = ToTableFunc(certificatesv1beta1.Convert_v1beta1_CertificateSigningRequest_To_certificates_CertificateSigningRequest)

// certificatesv1beta1 lists to metav1.Table functions
var CertificatesV1beta1CertificateSigningRequestListToTable = ToTableFunc(certificatesv1beta1.Convert_v1beta1_CertificateSigningRequestList_To_certificates_CertificateSigningRequestList)

// coordinationv1 objects to metav1.Table functions
var CoordinationV1LeaseToTable = ToTableFunc(coordinationv1.Convert_v1_Lease_To_coordination_Lease)

// coordinationv1 lists to metav1.Table functions
var CoordinationV1LeaseListToTable = ToTableFunc(coordinationv1.Convert_v1_LeaseList_To_coordination_LeaseList)

// coordinationv1beta1 objects to metav1.Table functions
var CoordinationV1beta1LeaseToTable = ToTableFunc(coordinationv1beta1.Convert_v1beta1_Lease_To_coordination_Lease)

// coordinationv1beta1 lists to metav1.Table functions
var CoordinationV1beta1LeaseListToTable = ToTableFunc(coordinationv1beta1.Convert_v1beta1_LeaseList_To_coordination_LeaseList)

// discoveryv1 objects to metav1.Table functions
var DiscoveryV1EndpointSliceToTable = ToTableFunc(discoveryv1.Convert_v1_EndpointSlice_To_discovery_EndpointSlice)

// discoveryv1 lists to metav1.Table functions
var DiscoveryV1EndpointSliceListToTable = ToTableFunc(discoveryv1.Convert_v1_EndpointSliceList_To_discovery_EndpointSliceList)

// discoveryv1beta1 objects to metav1.Table functions
var DiscoveryV1beta1EndpointSliceToTable = ToTableFunc(discoveryv1beta1.Convert_v1beta1_EndpointSlice_To_discovery_EndpointSlice)

// discoveryv1beta1 lists to metav1.Table functions
var DiscoveryV1beta1EndpointSliceListToTable = ToTableFunc(discoveryv1beta1.Convert_v1beta1_EndpointSliceList_To_discovery_EndpointSliceList)

// eventsv1 objects to metav1.Table functions
var EventsV1EventToTable = ToTableFunc(eventsv1.Convert_v1_Event_To_core_Event)

// eventsv1 lists to metav1.Table functions
var EventsV1EventListToTable = ToTableFunc(eventsv1.Convert_v1_EventList_To_core_EventList)

// eventsv1beta1 objects to metav1.Table functions
var EventsV1beta1EventToTable = ToTableFunc(eventsv1beta1.Convert_v1beta1_Event_To_core_Event)

// eventsv1beta1 lists to metav1.Table functions
var EventsV1beta1EventListToTable = ToTableFunc(eventsv1beta1.Convert_v1beta1_EventList_To_core_EventList)

// extensionsv1beta1 objects to metav1.Table functions
var (
	ExtensionsV1beta1DaemonSetToTable     = ToTableFunc(extensionsv1beta1.Convert_v1beta1_DaemonSet_To_apps_DaemonSet)
	ExtensionsV1beta1DeploymentToTable    = ToTableFunc(extensionsv1beta1.Convert_v1beta1_Deployment_To_apps_Deployment)
	ExtensionsV1beta1IngressToTable       = ToTableFunc(extensionsv1beta1.Convert_v1beta1_Ingress_To_networking_Ingress)
	ExtensionsV1beta1NetworkPolicyToTable = ToTableFunc(extensionsv1beta1.Convert_v1beta1_NetworkPolicy_To_networking_NetworkPolicy)
	ExtensionsV1beta1ReplicaSetToTable    = ToTableFunc(extensionsv1beta1.Convert_v1beta1_ReplicaSet_To_apps_ReplicaSet)
	ExtensionsV1beta1ScaleToTable         = ToTableFunc(extensionsv1beta1.Convert_v1beta1_Scale_To_autoscaling_Scale)
)

// extensionsv1beta1 lists to metav1.Table functions
var (
	ExtensionsV1beta1DaemonSetListToTable     = ToTableFunc(extensionsv1beta1.Convert_v1beta1_DaemonSetList_To_apps_DaemonSetList)
	ExtensionsV1beta1DeploymentListToTable    = ToTableFunc(extensionsv1beta1.Convert_v1beta1_DeploymentList_To_apps_DeploymentList)
	ExtensionsV1beta1IngressListToTable       = ToTableFunc(extensionsv1beta1.Convert_v1beta1_IngressList_To_networking_IngressList)
	ExtensionsV1beta1NetworkPolicyListToTable = ToTableFunc(extensionsv1beta1.Convert_v1beta1_NetworkPolicyList_To_networking_NetworkPolicyList)
	ExtensionsV1beta1ReplicaSetListToTable    = ToTableFunc(extensionsv1beta1.Convert_v1beta1_ReplicaSetList_To_apps_ReplicaSetList)
)

// flowcontrolv1 objects to metav1.Table functions
var (
	FlowcontrolV1FlowSchemaToTable                 = ToTableFunc(flowcontrolv1.Convert_v1_FlowSchema_To_flowcontrol_FlowSchema)
	FlowcontrolV1PriorityLevelConfigurationToTable = ToTableFunc(flowcontrolv1.Convert_v1_PriorityLevelConfiguration_To_flowcontrol_PriorityLevelConfiguration)
)

// flowcontrolv1 lists to metav1.Table functions
var (
	FlowcontrolV1FlowSchemaListToTable                 = ToTableFunc(flowcontrolv1.Convert_v1_FlowSchemaList_To_flowcontrol_FlowSchemaList)
	FlowcontrolV1PriorityLevelConfigurationListToTable = ToTableFunc(flowcontrolv1.Convert_v1_PriorityLevelConfigurationList_To_flowcontrol_PriorityLevelConfigurationList)
)

// flowcontrolv1beta1 objects to metav1.Table functions
var (
	FlowcontrolV1beta1FlowSchemaToTable                 = ToTableFunc(flowcontrolv1beta1.Convert_v1beta1_FlowSchema_To_flowcontrol_FlowSchema)
	FlowcontrolV1beta1PriorityLevelConfigurationToTable = ToTableFunc(flowcontrolv1beta1.Convert_v1beta1_PriorityLevelConfiguration_To_flowcontrol_PriorityLevelConfiguration)
)

// flowcontrolv1beta1 lists to metav1.Table functions
var (
	FlowcontrolV1beta1FlowSchemaListToTable                 = ToTableFunc(flowcontrolv1beta1.Convert_v1beta1_FlowSchemaList_To_flowcontrol_FlowSchemaList)
	FlowcontrolV1beta1PriorityLevelConfigurationListToTable = ToTableFunc(flowcontrolv1beta1.Convert_v1beta1_PriorityLevelConfigurationList_To_flowcontrol_PriorityLevelConfigurationList)
)

// flowcontrolv1beta2 objects to metav1.Table functions
var (
	FlowcontrolV1beta2FlowSchemaToTable                 = ToTableFunc(flowcontrolv1beta2.Convert_v1beta2_FlowSchema_To_flowcontrol_FlowSchema)
	FlowcontrolV1beta2PriorityLevelConfigurationToTable = ToTableFunc(flowcontrolv1beta2.Convert_v1beta2_PriorityLevelConfiguration_To_flowcontrol_PriorityLevelConfiguration)
)

// flowcontrolv1beta2 lists to metav1.Table functions
var (
	FlowcontrolV1beta2FlowSchemaListToTable                 = ToTableFunc(flowcontrolv1beta2.Convert_v1beta2_FlowSchemaList_To_flowcontrol_FlowSchemaList)
	FlowcontrolV1beta2PriorityLevelConfigurationListToTable = ToTableFunc(flowcontrolv1beta2.Convert_v1beta2_PriorityLevelConfigurationList_To_flowcontrol_PriorityLevelConfigurationList)
)

// flowcontrolv1beta3 objects to metav1.Table functions
var (
	FlowcontrolV1beta3FlowSchemaToTable                 = ToTableFunc(flowcontrolv1beta3.Convert_v1beta3_FlowSchema_To_flowcontrol_FlowSchema)
	FlowcontrolV1beta3PriorityLevelConfigurationToTable = ToTableFunc(flowcontrolv1beta3.Convert_v1beta3_PriorityLevelConfiguration_To_flowcontrol_PriorityLevelConfiguration)
)

// flowcontrolv1beta3 lists to metav1.Table functions
var (
	FlowcontrolV1beta3FlowSchemaListToTable                 = ToTableFunc(flowcontrolv1beta3.Convert_v1beta3_FlowSchemaList_To_flowcontrol_FlowSchemaList)
	FlowcontrolV1beta3PriorityLevelConfigurationListToTable = ToTableFunc(flowcontrolv1beta3.Convert_v1beta3_PriorityLevelConfigurationList_To_flowcontrol_PriorityLevelConfigurationList)
)

// imagepolicyv1alpha1 objects to metav1.Table functions
var ImagepolicyV1alpha1ImageReviewToTable = ToTableFunc(imagepolicyv1alpha1.Convert_v1alpha1_ImageReview_To_imagepolicy_ImageReview)

// networkingv1 objects to metav1.Table functions
var (
	NetworkingV1IngressClassToTable  = ToTableFunc(networkingv1.Convert_v1_IngressClass_To_networking_IngressClass)
	NetworkingV1IngressToTable       = ToTableFunc(networkingv1.Convert_v1_Ingress_To_networking_Ingress)
	NetworkingV1NetworkPolicyToTable = ToTableFunc(networkingv1.Convert_v1_NetworkPolicy_To_networking_NetworkPolicy)
)

// networkingv1 lists to metav1.Table functions
var (
	NetworkingV1IngressClassListToTable  = ToTableFunc(networkingv1.Convert_v1_IngressClassList_To_networking_IngressClassList)
	NetworkingV1IngressListToTable       = ToTableFunc(networkingv1.Convert_v1_IngressList_To_networking_IngressList)
	NetworkingV1NetworkPolicyListToTable = ToTableFunc(networkingv1.Convert_v1_NetworkPolicyList_To_networking_NetworkPolicyList)
)

// networkingv1alpha1 objects to metav1.Table functions
var (
	NetworkingV1alpha1IPAddressToTable   = ToTableFunc(networkingv1alpha1.Convert_v1alpha1_IPAddress_To_networking_IPAddress)
	NetworkingV1alpha1ServiceCIDRToTable = ToTableFunc(networkingv1alpha1.Convert_v1alpha1_ServiceCIDR_To_networking_ServiceCIDR)
)

// networkingv1alpha1 lists to metav1.Table functions
var (
	NetworkingV1alpha1IPAddressListToTable   = ToTableFunc(networkingv1alpha1.Convert_v1alpha1_IPAddressList_To_networking_IPAddressList)
	NetworkingV1alpha1ServiceCIDRListToTable = ToTableFunc(networkingv1alpha1.Convert_v1alpha1_ServiceCIDRList_To_networking_ServiceCIDRList)
)

// networkingv1beta1 objects to metav1.Table functions
var (
	NetworkingV1beta1IngressClassToTable = ToTableFunc(networkingv1beta1.Convert_v1beta1_IngressClass_To_networking_IngressClass)
	NetworkingV1beta1IngressToTable      = ToTableFunc(networkingv1beta1.Convert_v1beta1_Ingress_To_networking_Ingress)
)

// networkingv1beta1 lists to metav1.Table functions
var (
	NetworkingV1beta1IngressClassListToTable = ToTableFunc(networkingv1beta1.Convert_v1beta1_IngressClassList_To_networking_IngressClassList)
	NetworkingV1beta1IngressListToTable      = ToTableFunc(networkingv1beta1.Convert_v1beta1_IngressList_To_networking_IngressList)
)

// nodev1 objects to metav1.Table functions
var NodeV1RuntimeClassToTable = ToTableFunc(nodev1.Convert_v1_RuntimeClass_To_node_RuntimeClass)

// nodev1 lists to metav1.Table functions
var NodeV1RuntimeClassListToTable = ToTableFunc(nodev1.Convert_v1_RuntimeClassList_To_node_RuntimeClassList)

// nodev1alpha1 objects to metav1.Table functions
var NodeV1alpha1RuntimeClassToTable = ToTableFunc(nodev1alpha1.Convert_v1alpha1_RuntimeClass_To_node_RuntimeClass)

// nodev1alpha1 lists to metav1.Table functions
var NodeV1alpha1RuntimeClassListToTable = ToTableFunc(nodev1alpha1.Convert_v1alpha1_RuntimeClassList_To_node_RuntimeClassList)

// nodev1beta1 objects to metav1.Table functions
var NodeV1beta1RuntimeClassToTable = ToTableFunc(nodev1beta1.Convert_v1beta1_RuntimeClass_To_node_RuntimeClass)

// nodev1beta1 lists to metav1.Table functions
var NodeV1beta1RuntimeClassListToTable = ToTableFunc(nodev1beta1.Convert_v1beta1_RuntimeClassList_To_node_RuntimeClassList)

// policyv1 objects to metav1.Table functions
var (
	PolicyV1EvictionToTable            = ToTableFunc(policyv1.Convert_v1_Eviction_To_policy_Eviction)
	PolicyV1PodDisruptionBudgetToTable = ToTableFunc(policyv1.Convert_v1_PodDisruptionBudget_To_policy_PodDisruptionBudget)
)

// policyv1 lists to metav1.Table functions
var PolicyV1PodDisruptionBudgetListToTable = ToTableFunc(policyv1.Convert_v1_PodDisruptionBudgetList_To_policy_PodDisruptionBudgetList)

// policyv1beta1 objects to metav1.Table functions
var (
	PolicyV1beta1EvictionToTable            = ToTableFunc(policyv1beta1.Convert_v1beta1_Eviction_To_policy_Eviction)
	PolicyV1beta1PodDisruptionBudgetToTable = ToTableFunc(policyv1beta1.Convert_v1beta1_PodDisruptionBudget_To_policy_PodDisruptionBudget)
)

// policyv1beta1 lists to metav1.Table functions
var PolicyV1beta1PodDisruptionBudgetListToTable = ToTableFunc(policyv1beta1.Convert_v1beta1_PodDisruptionBudgetList_To_policy_PodDisruptionBudgetList)

// rbacv1 objects to metav1.Table functions
var (
	RbacV1ClusterRoleBindingToTable = ToTableFunc(rbacv1.Convert_v1_ClusterRoleBinding_To_rbac_ClusterRoleBinding)
	RbacV1ClusterRoleToTable        = ToTableFunc(rbacv1.Convert_v1_ClusterRole_To_rbac_ClusterRole)
	RbacV1RoleBindingToTable        = ToTableFunc(rbacv1.Convert_v1_RoleBinding_To_rbac_RoleBinding)
	RbacV1RoleToTable               = ToTableFunc(rbacv1.Convert_v1_Role_To_rbac_Role)
)

// rbacv1 lists to metav1.Table functions
var (
	RbacV1ClusterRoleBindingListToTable = ToTableFunc(rbacv1.Convert_v1_ClusterRoleBindingList_To_rbac_ClusterRoleBindingList)
	RbacV1ClusterRoleListToTable        = ToTableFunc(rbacv1.Convert_v1_ClusterRoleList_To_rbac_ClusterRoleList)
	RbacV1RoleBindingListToTable        = ToTableFunc(rbacv1.Convert_v1_RoleBindingList_To_rbac_RoleBindingList)
	RbacV1RoleListToTable               = ToTableFunc(rbacv1.Convert_v1_RoleList_To_rbac_RoleList)
)

// rbacv1alpha1 objects to metav1.Table functions
var (
	RbacV1alpha1ClusterRoleBindingToTable = ToTableFunc(rbacv1alpha1.Convert_v1alpha1_ClusterRoleBinding_To_rbac_ClusterRoleBinding)
	RbacV1alpha1ClusterRoleToTable        = ToTableFunc(rbacv1alpha1.Convert_v1alpha1_ClusterRole_To_rbac_ClusterRole)
	RbacV1alpha1RoleBindingToTable        = ToTableFunc(rbacv1alpha1.Convert_v1alpha1_RoleBinding_To_rbac_RoleBinding)
	RbacV1alpha1RoleToTable               = ToTableFunc(rbacv1alpha1.Convert_v1alpha1_Role_To_rbac_Role)
)

// rbacv1alpha1 lists to metav1.Table functions
var (
	RbacV1alpha1ClusterRoleBindingListToTable = ToTableFunc(rbacv1alpha1.Convert_v1alpha1_ClusterRoleBindingList_To_rbac_ClusterRoleBindingList)
	RbacV1alpha1ClusterRoleListToTable        = ToTableFunc(rbacv1alpha1.Convert_v1alpha1_ClusterRoleList_To_rbac_ClusterRoleList)
	RbacV1alpha1RoleBindingListToTable        = ToTableFunc(rbacv1alpha1.Convert_v1alpha1_RoleBindingList_To_rbac_RoleBindingList)
	RbacV1alpha1RoleListToTable               = ToTableFunc(rbacv1alpha1.Convert_v1alpha1_RoleList_To_rbac_RoleList)
)

// rbacv1beta1 objects to metav1.Table functions
var (
	RbacV1beta1ClusterRoleBindingToTable = ToTableFunc(rbacv1beta1.Convert_v1beta1_ClusterRoleBinding_To_rbac_ClusterRoleBinding)
	RbacV1beta1ClusterRoleToTable        = ToTableFunc(rbacv1beta1.Convert_v1beta1_ClusterRole_To_rbac_ClusterRole)
	RbacV1beta1RoleBindingToTable        = ToTableFunc(rbacv1beta1.Convert_v1beta1_RoleBinding_To_rbac_RoleBinding)
	RbacV1beta1RoleToTable               = ToTableFunc(rbacv1beta1.Convert_v1beta1_Role_To_rbac_Role)
)

// rbacv1beta1 lists to metav1.Table functions
var (
	RbacV1beta1ClusterRoleBindingListToTable = ToTableFunc(rbacv1beta1.Convert_v1beta1_ClusterRoleBindingList_To_rbac_ClusterRoleBindingList)
	RbacV1beta1ClusterRoleListToTable        = ToTableFunc(rbacv1beta1.Convert_v1beta1_ClusterRoleList_To_rbac_ClusterRoleList)
	RbacV1beta1RoleBindingListToTable        = ToTableFunc(rbacv1beta1.Convert_v1beta1_RoleBindingList_To_rbac_RoleBindingList)
	RbacV1beta1RoleListToTable               = ToTableFunc(rbacv1beta1.Convert_v1beta1_RoleList_To_rbac_RoleList)
)

// resourcev1alpha2 objects to metav1.Table functions
var (
	ResourceV1alpha2PodSchedulingContextToTable    = ToTableFunc(resourcev1alpha2.Convert_v1alpha2_PodSchedulingContext_To_resource_PodSchedulingContext)
	ResourceV1alpha2ResourceClaimParametersToTable = ToTableFunc(resourcev1alpha2.Convert_v1alpha2_ResourceClaimParameters_To_resource_ResourceClaimParameters)
	ResourceV1alpha2ResourceClaimTemplateToTable   = ToTableFunc(resourcev1alpha2.Convert_v1alpha2_ResourceClaimTemplate_To_resource_ResourceClaimTemplate)
	ResourceV1alpha2ResourceClaimToTable           = ToTableFunc(resourcev1alpha2.Convert_v1alpha2_ResourceClaim_To_resource_ResourceClaim)
	ResourceV1alpha2ResourceClassParametersToTable = ToTableFunc(resourcev1alpha2.Convert_v1alpha2_ResourceClassParameters_To_resource_ResourceClassParameters)
	ResourceV1alpha2ResourceClassToTable           = ToTableFunc(resourcev1alpha2.Convert_v1alpha2_ResourceClass_To_resource_ResourceClass)
	ResourceV1alpha2ResourceSliceToTable           = ToTableFunc(resourcev1alpha2.Convert_v1alpha2_ResourceSlice_To_resource_ResourceSlice)
)

// resourcev1alpha2 lists to metav1.Table functions
var (
	ResourceV1alpha2PodSchedulingContextListToTable    = ToTableFunc(resourcev1alpha2.Convert_v1alpha2_PodSchedulingContextList_To_resource_PodSchedulingContextList)
	ResourceV1alpha2ResourceClaimListToTable           = ToTableFunc(resourcev1alpha2.Convert_v1alpha2_ResourceClaimList_To_resource_ResourceClaimList)
	ResourceV1alpha2ResourceClaimParametersListToTable = ToTableFunc(resourcev1alpha2.Convert_v1alpha2_ResourceClaimParametersList_To_resource_ResourceClaimParametersList)
	ResourceV1alpha2ResourceClaimTemplateListToTable   = ToTableFunc(resourcev1alpha2.Convert_v1alpha2_ResourceClaimTemplateList_To_resource_ResourceClaimTemplateList)
	ResourceV1alpha2ResourceClassListToTable           = ToTableFunc(resourcev1alpha2.Convert_v1alpha2_ResourceClassList_To_resource_ResourceClassList)
	ResourceV1alpha2ResourceClassParametersListToTable = ToTableFunc(resourcev1alpha2.Convert_v1alpha2_ResourceClassParametersList_To_resource_ResourceClassParametersList)
	ResourceV1alpha2ResourceSliceListToTable           = ToTableFunc(resourcev1alpha2.Convert_v1alpha2_ResourceSliceList_To_resource_ResourceSliceList)
)

// schedulingv1 objects to metav1.Table functions
var SchedulingV1PriorityClassToTable = ToTableFunc(schedulingv1.Convert_v1_PriorityClass_To_scheduling_PriorityClass)

// schedulingv1 lists to metav1.Table functions
var SchedulingV1PriorityClassListToTable = ToTableFunc(schedulingv1.Convert_v1_PriorityClassList_To_scheduling_PriorityClassList)

// schedulingv1alpha1 objects to metav1.Table functions
var SchedulingV1alpha1PriorityClassToTable = ToTableFunc(schedulingv1alpha1.Convert_v1alpha1_PriorityClass_To_scheduling_PriorityClass)

// schedulingv1alpha1 lists to metav1.Table functions
var SchedulingV1alpha1PriorityClassListToTable = ToTableFunc(schedulingv1alpha1.Convert_v1alpha1_PriorityClassList_To_scheduling_PriorityClassList)

// schedulingv1beta1 objects to metav1.Table functions
var SchedulingV1beta1PriorityClassToTable = ToTableFunc(schedulingv1beta1.Convert_v1beta1_PriorityClass_To_scheduling_PriorityClass)

// schedulingv1beta1 lists to metav1.Table functions
var SchedulingV1beta1PriorityClassListToTable = ToTableFunc(schedulingv1beta1.Convert_v1beta1_PriorityClassList_To_scheduling_PriorityClassList)

// storagev1 objects to metav1.Table functions
var (
	StorageV1CSIDriverToTable          = ToTableFunc(storagev1.Convert_v1_CSIDriver_To_storage_CSIDriver)
	StorageV1CSINodeToTable            = ToTableFunc(storagev1.Convert_v1_CSINode_To_storage_CSINode)
	StorageV1CSIStorageCapacityToTable = ToTableFunc(storagev1.Convert_v1_CSIStorageCapacity_To_storage_CSIStorageCapacity)
	StorageV1StorageClassToTable       = ToTableFunc(storagev1.Convert_v1_StorageClass_To_storage_StorageClass)
	StorageV1VolumeAttachmentToTable   = ToTableFunc(storagev1.Convert_v1_VolumeAttachment_To_storage_VolumeAttachment)
)

// storagev1 lists to metav1.Table functions
var (
	StorageV1CSIDriverListToTable          = ToTableFunc(storagev1.Convert_v1_CSIDriverList_To_storage_CSIDriverList)
	StorageV1CSINodeListToTable            = ToTableFunc(storagev1.Convert_v1_CSINodeList_To_storage_CSINodeList)
	StorageV1CSIStorageCapacityListToTable = ToTableFunc(storagev1.Convert_v1_CSIStorageCapacityList_To_storage_CSIStorageCapacityList)
	StorageV1StorageClassListToTable       = ToTableFunc(storagev1.Convert_v1_StorageClassList_To_storage_StorageClassList)
	StorageV1VolumeAttachmentListToTable   = ToTableFunc(storagev1.Convert_v1_VolumeAttachmentList_To_storage_VolumeAttachmentList)
)

// storagev1alpha1 objects to metav1.Table functions
var (
	StorageV1alpha1CSIStorageCapacityToTable    = ToTableFunc(storagev1alpha1.Convert_v1alpha1_CSIStorageCapacity_To_storage_CSIStorageCapacity)
	StorageV1alpha1VolumeAttachmentToTable      = ToTableFunc(storagev1alpha1.Convert_v1alpha1_VolumeAttachment_To_storage_VolumeAttachment)
	StorageV1alpha1VolumeAttributesClassToTable = ToTableFunc(storagev1alpha1.Convert_v1alpha1_VolumeAttributesClass_To_storage_VolumeAttributesClass)
)

// storagev1alpha1 lists to metav1.Table functions
var (
	StorageV1alpha1CSIStorageCapacityListToTable    = ToTableFunc(storagev1alpha1.Convert_v1alpha1_CSIStorageCapacityList_To_storage_CSIStorageCapacityList)
	StorageV1alpha1VolumeAttachmentListToTable      = ToTableFunc(storagev1alpha1.Convert_v1alpha1_VolumeAttachmentList_To_storage_VolumeAttachmentList)
	StorageV1alpha1VolumeAttributesClassListToTable = ToTableFunc(storagev1alpha1.Convert_v1alpha1_VolumeAttributesClassList_To_storage_VolumeAttributesClassList)
)

// storagev1beta1 objects to metav1.Table functions
var (
	StorageV1beta1CSIDriverToTable          = ToTableFunc(storagev1beta1.Convert_v1beta1_CSIDriver_To_storage_CSIDriver)
	StorageV1beta1CSINodeToTable            = ToTableFunc(storagev1beta1.Convert_v1beta1_CSINode_To_storage_CSINode)
	StorageV1beta1CSIStorageCapacityToTable = ToTableFunc(storagev1beta1.Convert_v1beta1_CSIStorageCapacity_To_storage_CSIStorageCapacity)
	StorageV1beta1StorageClassToTable       = ToTableFunc(storagev1beta1.Convert_v1beta1_StorageClass_To_storage_StorageClass)
	StorageV1beta1VolumeAttachmentToTable   = ToTableFunc(storagev1beta1.Convert_v1beta1_VolumeAttachment_To_storage_VolumeAttachment)
)

// storagev1beta1 lists to metav1.Table functions
var (
	StorageV1beta1CSIDriverListToTable          = ToTableFunc(storagev1beta1.Convert_v1beta1_CSIDriverList_To_storage_CSIDriverList)
	StorageV1beta1CSINodeListToTable            = ToTableFunc(storagev1beta1.Convert_v1beta1_CSINodeList_To_storage_CSINodeList)
	StorageV1beta1CSIStorageCapacityListToTable = ToTableFunc(storagev1beta1.Convert_v1beta1_CSIStorageCapacityList_To_storage_CSIStorageCapacityList)
	StorageV1beta1StorageClassListToTable       = ToTableFunc(storagev1beta1.Convert_v1beta1_StorageClassList_To_storage_StorageClassList)
	StorageV1beta1VolumeAttachmentListToTable   = ToTableFunc(storagev1beta1.Convert_v1beta1_VolumeAttachmentList_To_storage_VolumeAttachmentList)
)

// storagemigrationv1alpha1 objects to metav1.Table functions
var StoragemigrationV1alpha1StorageVersionMigrationToTable = ToTableFunc(storagemigrationv1alpha1.Convert_v1alpha1_StorageVersionMigration_To_storagemigration_StorageVersionMigration)

// storagemigrationv1alpha1 lists to metav1.Table functions
var StoragemigrationV1alpha1StorageVersionMigrationListToTable = ToTableFunc(storagemigrationv1alpha1.Convert_v1alpha1_StorageVersionMigrationList_To_storagemigration_StorageVersionMigrationList)
