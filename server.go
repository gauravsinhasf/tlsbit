package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net"
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
			fmt.Printf("Certificate verified\n")
		}
		return nil
	}

	l, err := tls.Listen("tcp4", ":443", cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		// Wait for a connection.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn) {
			fmt.Printf("accepting new conection\n")
			// Echo all incoming data.
			io.Copy(c, c)
			// Shut down the connection.
			c.Close()
		}(conn)
	}
}
