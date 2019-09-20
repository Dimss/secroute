package routewebhook

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"k8s.io/api/admission/v1beta1"
	routev1 "github.com/openshift/api/route/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func ValidateWebHookHandler(w http.ResponseWriter, r *http.Request) {
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
		sendAdmissionValidationRouterResponse(w, false, errMessage)
		logrus.Errorf(errMessage)
		return
	}
	// This object gonna hold actual username
	var route = routev1.Route{}
	ar := v1beta1.AdmissionReview{}
	// Try to decode body into Admission Review object
	if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
		logrus.Errorf("Error during deserializing request body: %v", err)
		sendAdmissionValidationRouterResponse(w, false, "error during deserializing request body")
		return
	}
	// Try to unmarshal Admission Review raw object to Router
	if err := json.Unmarshal(ar.Request.Object.Raw, &route); err != nil {
		errMessage := "Error during unmarshaling request body"
		logrus.Error(errMessage)
		sendAdmissionValidationRouterResponse(w, false, errMessage)
		return
	}
	if route.Spec.TLS == nil {
		errMessage := fmt.Sprintf("Creation of insecure routes are forbiden, route: %v", route.Name)
		logrus.Warn(errMessage)
		sendAdmissionValidationRouterResponse(w, false, errMessage)
	} else {
		sendAdmissionValidationRouterResponse(w, true, "Router is secure, proceed request")
	}
}

func MutateWebHookHandler(w http.ResponseWriter, r *http.Request) {
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
		sendAdmissionValidationRouterResponse(w, false, errMessage)
		logrus.Errorf(errMessage)
		return
	}
	// This object gonna hold actual username
	var route = routev1.Route{}
	ar := v1beta1.AdmissionReview{}
	// Try to decode body into Admission Review object
	if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
		logrus.Errorf("Error during deserializing request body: %v", err)
		sendAdmissionValidationRouterResponse(w, false, "error during deserializing request body")
		return
	}
	// Try to unmarshal Admission Review raw object to Router
	if err := json.Unmarshal(ar.Request.Object.Raw, &route); err != nil {
		errMessage := "Error during unmarshaling request body"
		logrus.Error(errMessage)
		sendAdmissionValidationRouterResponse(w, false, errMessage)
		return
	}

	if route.Spec.TLS != nil {
		sendAdmissionValidationRouterResponse(w, true, "Router is secure, proceed request")
	} else {
		sendAdmissionMutationRouterResponse(ar.Request.UID, w)
	}
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

func sendAdmissionValidationRouterResponse(w http.ResponseWriter, isAllowed bool, message string) {
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

func LivenessHandler(w http.ResponseWriter, r *http.Request) {

}
