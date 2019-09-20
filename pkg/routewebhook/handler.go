package routewebhook

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"net/http"
)

var (
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()
)

func WebHookHandler(w http.ResponseWriter, r *http.Request, adUsersChan chan string) {
	//logrus.Info("Handling oauthtokenwebhook webhook")
	//var body []byte
	//// Read request body
	//if r.Body != nil {
	//	if data, err := ioutil.ReadAll(r.Body); err == nil {
	//		body = data
	//	}
	//}
	//// K8S sends POST request with the admission webhook data,
	//// the body can't be empty, but if it is,
	//// further processing will be stopped and empty
	//// admission response will be sent to K8S API
	//if len(body) == 0 {
	//	sendAdmissionResponse(w)
	//	logrus.Errorf("SKIPPING USER PROCESSING DURING PREVIOUS ERRORS!")
	//	return
	//}
	//// This object gonna hold actual username
	//var oauthToken oauthv1.OAuthAccessToken
	//ar := v1beta1.AdmissionReview{}
	//// Try to decode body into Admission Review object
	//if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
	//	logrus.Errorf("Error during deserializing request body: %v", err)
	//	logrus.Errorf("SKIPPING USER PROCESSING DURING PREVIOUS ERRORS!")
	//	sendAdmissionResponse(w)
	//	return
	//}
	//// Try to unmarshal Admission Review raw object to OAuthAccessToken
	//if err := json.Unmarshal(ar.Request.Object.Raw, &oauthToken); err != nil {
	//	logrus.Error("Error during unmarshaling request body")
	//	logrus.Errorf("SKIPPING USER PROCESSING DURING PREVIOUS ERRORS!")
	//	sendAdmissionResponse(w)
	//	return
	//}
	//sendAdmissionResponse(w)
	//logrus.Infof("Passing user: %s to channel for further processing", oauthToken.UserName)
	//// Write AD user to channel for further processing
	//adUsersChan <- oauthToken.UserName
	//logrus.Info("request is done. . . .")
}

func sendAdmissionResponse(w http.ResponseWriter) {

}

func LivenessHandler(w http.ResponseWriter, r *http.Request) {


}