package main

import (
	"encoding/json"
	"io/ioutil"

	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"

	"k8s.io/client-go/tools/clientcmd"

	restclient "k8s.io/client-go/rest"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/metrics/pkg/apis/metrics"

	"github.com/mobingilabs/ouchan/services/ocean-alibaba/kubernetes"
)

func main() {
	c := newConfig()
	confByte, err := kubernetes.CreateKubeconf(c)
	if err != nil {
		panic(err)
	}
	client := newRESTClient(confByte)
	getMetric(client)
}

func newConfig() *kubernetes.Config {
	caCert, clientCert, clientKey := readCert()

	return &kubernetes.Config{
		APIServer:   "https://161.117.98.207:6443/",
		ClusterName: "kubernetes",
		CaCert:      caCert,
		ClientCert:  clientCert,
		ClientKey:   clientKey,
	}
}

func readCert() ([]byte, []byte, []byte) {
	data, err := ioutil.ReadFile("ca.json")
	if err != nil {
		panic(err)
	}

	var ca struct {
		Ca   string `json:"ca"`
		Cert string `json:"cert"`
		Key  string `json:"key"`
	}

	err = json.Unmarshal(data, &ca)
	if err != nil {
		panic(err)
	}

	return []byte(ca.Ca), []byte(ca.Cert), []byte(ca.Key)
}

func newRESTClient(kubeconf []byte) restclient.Interface {
	config, err := clientcmd.RESTConfigFromKubeConfig(kubeconf)
	if err != nil {
		panic(err)
	}
	setConfigDefaults(config)

	client, err := restclient.RESTClientFor(config)
	if err != nil {
		panic(err)
	}

	return client
}

func getMetric(client restclient.Interface) {
	results := metrics.NodeMetricsList{}
	err := client.Get().Suffix("nodes").Do().Into(&results)
	if err != nil {
		panic(err)
	}
}

func setConfigDefaults(config *restclient.Config) {
	gv := schema.GroupVersion{Group: "", Version: "v1alpha1"}
	config.GroupVersion = &gv
	config.APIPath = "/apis/metrics"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	config.UserAgent = restclient.DefaultKubernetesUserAgent()
}
