package common

import (
	// "bytes"
	"crypto/tls"
	"crypto/x509"
	// "fmt"
	"io/ioutil"
	"log"
	"net/http"
	// "reflect"
	// "strconv"
)

func HttpClientChkSrv(url string) []byte {
	// get CA cert
	caCert, err := ioutil.ReadFile("./cert/ca.crt")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair("./cert/server.crt", "./cert/server.key")
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
	}
	res, err := client.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	// log.Printf("%s health check result: %v\n", url, string(body))
	// return string(body)
	return body
}
