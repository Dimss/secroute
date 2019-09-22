package routewebhook

import (
	"encoding/json"
	"fmt"
	routev1 "github.com/openshift/api/route/v1"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"net/http"
	routev1Configs "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
)

var (
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()
)

const (
	addEdgeTls string = `[
    { "op":"add","path":"/spec/tls","value":{"termination":"edge"} }
  ]`
)

func ValidateRouteWebHookHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Handling route webhook")
	var body []byte
	// Read request body
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	// K8S sends POST request with the admission webhook data,
	// the body can't be empty, but if it is,
	// further processing will be stopped and empty
	// admission response will be sent to K8S API
	if len(body) == 0 {
		errMessage := "The body is empty, can't proceed the request"
		sendAdmissionValidationResponse(w, false, errMessage)
		logrus.Errorf(errMessage)
		return
	}
	// This object gonna hold actual route
	var route = routev1.Route{}
	ar := v1beta1.AdmissionReview{}
	// Try to decode body into Admission Review object
	if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
		logrus.Errorf("Error during deserializing request body: %v", err)
		sendAdmissionValidationResponse(w, false, "error during deserializing request body")
		return
	}
	// Try to unmarshal Admission Review raw object to Router
	if err := json.Unmarshal(ar.Request.Object.Raw, &route); err != nil {
		errMessage := "Error during unmarshaling request body"
		logrus.Error(errMessage)
		sendAdmissionValidationResponse(w, false, errMessage)
		return
	}
	if route.Spec.TLS == nil {
		errMessage := fmt.Sprintf("Creation of insecure routes are forbiden, route: %v", route.Name)
		logrus.Warn(errMessage)
		sendAdmissionValidationResponse(w, false, errMessage)
	} else {
		sendAdmissionValidationResponse(w, true, "Router is secure, proceed request")
	}
}

func MutateRouteWebHookHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Handling route webhook")
	var body []byte
	// Read request body
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	// K8S sends POST request with the admission webhook data,
	// the body can't be empty, but if it is,
	// further processing will be stopped and empty
	// admission response will be sent to K8S API
	if len(body) == 0 {
		errMessage := "The body is empty, can't proceed the request"
		sendAdmissionValidationResponse(w, false, errMessage)
		logrus.Errorf(errMessage)
		return
	}
	// This object gonna hold actual route
	var route = routev1.Route{}
	ar := v1beta1.AdmissionReview{}
	// Try to decode body into Admission Review object
	if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
		logrus.Errorf("Error during deserializing request body: %v", err)
		sendAdmissionValidationResponse(w, false, "error during deserializing request body")
		return
	}
	// Try to unmarshal Admission Review raw object to Router
	if err := json.Unmarshal(ar.Request.Object.Raw, &route); err != nil {
		errMessage := "Error during unmarshaling request body"
		logrus.Error(errMessage)
		sendAdmissionValidationResponse(w, false, errMessage)
		return
	}

	if route.Spec.TLS != nil {
		sendAdmissionValidationResponse(w, true, "Router is secure, proceed request")
	} else {
		sendAdmissionMutationRouterResponse(ar.Request.UID, w)
	}
}

func CreateRouteOnServiceWebHookHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Handling service webhook")
	var body []byte
	// Read request body
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}
	// K8S sends POST request with the admission webhook data,
	// the body can't be empty, but if it is,
	// further processing will be stopped and empty
	// admission response will be sent to K8S API
	if len(body) == 0 {
		errMessage := "The body is empty, can't proceed the request"
		sendAdmissionValidationResponse(w, false, errMessage)
		logrus.Errorf(errMessage)
		return
	}

	ar := v1beta1.AdmissionReview{}
	// Try to decode body into Admission Review object
	if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
		logrus.Errorf("Error during deserializing request body: %v", err)
		sendAdmissionValidationResponse(w, false, "error during deserializing request body")
		return
	}

	if ar.Request.Operation == "CREATE" {
		// This object gonna hold actual service
		var service = corev1.Service{}
		// Try to unmarshal Admission Review raw object to Router
		if err := json.Unmarshal(ar.Request.Object.Raw, &service); err != nil {
			errMessage := "Error during unmarshaling request body"
			logrus.Error(errMessage)
			sendAdmissionValidationResponse(w, false, errMessage)
			return
		}

		if value, ok := service.Labels["addRoute"]; ok {
			if value == "true" {
				err := CreateRouteForService(service.Name, service.Namespace)
				if err != nil {
					sendAdmissionValidationResponse(w, false, err.Error())
				}
			}
		}
	}

	if ar.Request.Operation == "DELETE" {
		err := DeleteRouteForService(ar.Request.Name, ar.Request.Namespace)
		if err != nil {
			sendAdmissionValidationResponse(w, false, err.Error())
		}
	}
	// All good proceed with request
	sendAdmissionValidationResponse(w, true, "All good, proceed request")
}

func sendAdmissionMutationRouterResponse(uuid types.UID, w http.ResponseWriter) {

	// Compose admission response
	admissionResponse := &v1beta1.AdmissionResponse{}
	admissionResponse.Allowed = true
	admissionResponse.Patch = []byte(addEdgeTls)
	pt := v1beta1.PatchTypeJSONPatch
	admissionResponse.PatchType = &pt
	// Compose admission review
	admissionReview := v1beta1.AdmissionReview{}
	admissionReview.Response = admissionResponse
	admissionReview.Response.UID = uuid

	resp, err := json.Marshal(admissionReview)
	if err != nil {
		logrus.Errorf("Error during marshaling admissionResponse object: %v", err)
		http.Error(w, fmt.Sprintf("Error during marshaling admissionResponse object: %w", err), http.StatusInternalServerError)
	}
	logrus.Info("Sending response to API server")
	if _, err := w.Write(resp); err != nil {
		logrus.Errorf("Can't write response: %v", err)
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}
}

func sendAdmissionValidationResponse(w http.ResponseWriter, isAllowed bool, message string) {
	var admissionResponse *v1beta1.AdmissionResponse
	admissionResponse = &v1beta1.AdmissionResponse{Allowed: isAllowed, Result: &metav1.Status{Message: message}}
	admissionReview := v1beta1.AdmissionReview{}
	admissionReview.Response = admissionResponse
	resp, err := json.Marshal(admissionReview)
	if err != nil {
		logrus.Errorf("Error during marshaling admissionResponse object: %v", err)
		http.Error(w, fmt.Sprintf("Error during marshaling admissionResponse object: %w", err), http.StatusInternalServerError)
	}
	logrus.Info("Sending response to API server")
	if _, err := w.Write(resp); err != nil {
		logrus.Errorf("Can't write response: %v", err)
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}
}

func CreateRouteForService(serviceName string, namespace string) error {
	var route = routev1.Route{}
	route.Name = fmt.Sprintf("%s-route", serviceName)
	route.Spec.To = routev1.RouteTargetReference{Kind: "Service", Name: serviceName}
	routerv1Client, err := routev1Configs.NewForConfig(getClientcmdConfigs())
	if err != nil {

		logrus.Errorf(err.Error())
		return err

	}
	_, err = routerv1Client.Routes(namespace).Create(&route)
	if err != nil {
		logrus.Errorf(err.Error())
		return err
	}
	return nil
}

func DeleteRouteForService(serviceName string, namespace string) error {

	routerv1Client, err := routev1Configs.NewForConfig(getClientcmdConfigs())
	if err != nil {
		logrus.Errorf(err.Error())
		return err
	}
	routeNameForDelete := fmt.Sprintf("%s-route", serviceName)
	err = routerv1Client.Routes(namespace).Delete(routeNameForDelete, nil)
	if err != nil {
		logrus.Errorf(err.Error())
		return err
	}
	return nil
}

func getClientcmdConfigs() *rest.Config {
	conf := viper.GetString("kubeconfig")
	if conf == "useInClusterConfig" {
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		return config
	} else {
		config, err := clientcmd.BuildConfigFromFlags("", conf)
		if err != nil {
			panic(err.Error())
		}
		return config
	}
}
