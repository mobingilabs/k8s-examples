package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"

	"k8s.io/client-go/tools/clientcmd"

	restclient "k8s.io/client-go/rest"

	"n1ce/k8s/pkg/kubernetes"
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
	fmt.Println(client.Get().Suffix("nodes").URL().String())
}

func setConfigDefaults(config *restclient.Config) {
	gv := v1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}
	config.UserAgent = restclient.DefaultKubernetesUserAgent()
}
