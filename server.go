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
	cfg := &tls.Config{}
	cfg.GetCertificate = GetCertificate

	// Request client certificate from the client
	// cfg.ClientAuth = tls.RequestClientCert

	//this func is called post TLS cert. verification with raw certs.
	cfg.VerifyPeerCertificate = func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		for _, v := range rawCerts {
			//fmt.Printf("Certificate verified:%v\n", v)
			_ = v
		}
		return nil
	}

	l, err := tls.Listen("tcp4", ":443", cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	var connections map[string]net.Conn = make(map[string]net.Conn, 1000)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func(c net.Conn) {
			buf := make([]byte, 512)
			n, _ := c.Read(buf)
			fmt.Printf("accepting new conection:%q\n", buf[:n])

			io.WriteString(c, `<!DOCTYPE html>
								<html>
								    <head>
								        <title>Standard response</title>
								    </head>
								    <body>
								Hello World!
								    </body>
								</html>`)

			connections[c.LocalAddr().String()] = c
			for key, _ := range connections {
				fmt.Printf("\nconnections:%+v\n", key)
			}
			// Shut down the connection.
			c.Close()
		}(conn)
	}
}

//Allows working with the client hello message post TLS cert verify
func GetCertificate(clientHelloInfo *tls.ClientHelloInfo) (*tls.Certificate, error) {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		fmt.Printf("LoadX509KeyPair() failed:%v\n", err)
		return nil, nil
	}

	fmt.Printf("cert:%+v\n", cert.Leaf)

	fmt.Printf("GetCertificate: client Hello:%v \n", clientHelloInfo.ServerName)
	return &cert, nil
}
