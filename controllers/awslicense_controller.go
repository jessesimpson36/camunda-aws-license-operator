/*
Copyright 2023 Distribution Team.

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
	appsv1 "k8s.io/api/apps/v1"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	//"github.com/aws/aws-sdk-go-v2/aws"
	//"github.com/aws/aws-sdk-go/aws/session"
	//"github.com/aws/aws-sdk-go-v2/config"
	//"github.com/aws/aws-sdk-go/service/licensemanager"
	"k8s.io/apimachinery/pkg/types"

	camundaiov1alpha1 "camunda.io/camunda-aws-license-operator/api/v1alpha1"
)

// AWSLicenseReconciler reconciles a AWSLicense object
type AWSLicenseReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *AWSLicenseReconciler) ScaleAllDeploymentsUp(ctx context.Context, req ctrl.Request, log logr.Logger) (ctrl.Result, error) {
	existingDeployment := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{
		Namespace: req.Namespace,
		Name:      "zeebe-0",
	}, existingDeployment)
	if err == nil {
		log.Info("Setting " + "zeebe-0 " + " to 0 replicas")
		var newReplicaCount int32 = 1

		existingDeployment.Spec.Replicas = &newReplicaCount
		err = r.Update(ctx, existingDeployment)
		if err != nil {
			log.Error(err, "Failed to update "+"zeebe-0 "+" deployment")
			return ctrl.Result{Requeue: true}, err
		}
	}
	return ctrl.Result{}, nil
}

func (r *AWSLicenseReconciler) ScaleAllDeploymentsDown(ctx context.Context, req ctrl.Request, log logr.Logger) (ctrl.Result, error) {
	existingDeployment := &appsv1.Deployment{}
	err := r.Get(ctx, types.NamespacedName{
		Namespace: req.Namespace,
		Name:      "cpt-operate",
	}, existingDeployment)
	if err == nil {
		log.Info("Setting " + "cpt-operate" + " to 0 replicas")
		var newReplicaCount int32 = 0

		existingDeployment.Spec.Replicas = &newReplicaCount
		err = r.Update(ctx, existingDeployment)
		if err != nil {
			log.Error(err, "Failed to update "+"cpt-operate "+" deployment")
			return ctrl.Result{Requeue: true}, err
		}
	}
	return ctrl.Result{}, nil
}

func (r *AWSLicenseReconciler) DetermineAWSLicenseValidity(ctx context.Context, req ctrl.Request, log logr.Logger) (ctrl.Result, error) {
	// Load the Shared AWS Configuration (~/.aws/config)
	//	cfg, err := config.LoadDefaultConfig(context.TODO())
	//	if err != nil {
	//		log.Error(err, "Failed to load AWS Configuration from ~/.aws/config")
	//	}

	//	session := session.Must(session.NewSession())
	//	licenseManager := licensemanager.New(session)
	//	arnString := "blah"
	//
	//	checkoutLicenseInput := licensemanager.CheckoutLicenseInput{
	//		Beneficiary: "",
	//		CheckoutType: "",
	//		ClientToken: "",
	//		Entitlements: [],
	//		KeyFingerprint: "",
	//		NodeId: "",
	//		ProductSKU: "",
	//	}
	//
	//
	//	licenseArn := licensemanager.GetLicenseInput{
	//		LicenseArn: &arnString,
	//	}
	//	licenseOutput, err := licenseManager.GetLicense(&licenseArn)
	//	if err != nil {
	//		log.Error(err, "Failed to get license")
	//	}

	//license := *licenseOutput.License
	//status := *license.Status
	status := "NOT_AVAILABLE"
	if status == "AVAILABLE" || status == "PENDING_AVAILABLE" {
		r.ScaleAllDeploymentsUp(ctx, req, log)
	} else {
		r.ScaleAllDeploymentsDown(ctx, req, log)
	}
	return ctrl.Result{}, nil
}

//+kubebuilder:rbac:groups=camunda.io,resources=awslicenses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=camunda.io,resources=awslicenses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=camunda.io,resources=awslicenses/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AWSLicense object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *AWSLicenseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log, _ := logr.FromContext(ctx)
	r.DetermineAWSLicenseValidity(ctx, req, log)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AWSLicenseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&camundaiov1alpha1.AWSLicense{}).
		Complete(r)
}
