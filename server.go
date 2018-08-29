package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"html"
	"log"
	"net/http"
	"time"
)

func main() {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		fmt.Printf("LoadX509KeyPair() failed:%v", err)
		return
	}

	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}
	// Request client certificate from the client
	cfg.ClientAuth = tls.RequestClientCert

	//this func is called post TLS cert. verification with raw certs.
	cfg.VerifyPeerCertificate = func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		for _, v := range rawCerts {
			fmt.Printf("Certificate verified:%+v\n", v)
		}
		return nil
	}

	srv := &http.Server{
		TLSConfig:    cfg,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}

	http.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from tls bit, %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(srv.ListenAndServeTLS("", ""))
}
