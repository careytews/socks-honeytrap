package main

import (
	//"bytes"
	"crypto/tls"
	"crypto/x509"
	log "github.com/sirupsen/logrus"
	"github.com/trustnetworks/go-socks5"
	"io/ioutil"
	//"net/http"
	"os"
	"path"
	//"strings"
	"flag"
	"fmt"
)

var (
	//hostname, _ = os.Hostname()
	hostname = "Anonymous SHA2 Secure Server"

	dir      = path.Join(os.Getenv("HOME"), ".mitm")
	keyFile  = path.Join(dir, "ca-key.pem")
	certFile = path.Join(dir, "ca-cert.pem")
)

func main() {

	if logLevel, err := log.ParseLevel(os.Getenv("LOG_LEVEL")); err == nil {
		log.SetLevel(logLevel)
	}

	mitmOn := flag.Bool("m", false, "MitM HTTPS connections")
	// port 1080 is conventional port as defined in RFC1928
	port := flag.Int("p", 1080, "SOCKS5 serving port")
	flag.Parse()

	addr := fmt.Sprintf(":%d", *port)

	// Set up SOCKS5 server config
	conf := &socks5.Config{}

	if *mitmOn {
		log.Println("will mitm HTTPS")
		ca, err := loadCA()
		if err != nil {
			log.Fatal(err)
		}

		// Supplying a CA switches MitM on
		conf.CA = &ca
	}

	// Create SOCKS5 proxy
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	if err := server.ListenAndServe("tcp", addr); err != nil {
		panic(err)
	}
}

func loadCA() (cert tls.Certificate, err error) {
	// TODO(kr): check file permissions
	cert, err = tls.LoadX509KeyPair(certFile, keyFile)
	if os.IsNotExist(err) {
		cert, err = genCA()
	}
	if err == nil {
		cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	}
	return
}

func genCA() (cert tls.Certificate, err error) {
	err = os.MkdirAll(dir, 0700)
	if err != nil {
		return
	}
	certPEM, keyPEM, err := socks5.GenCA(hostname)
	if err != nil {
		return
	}
	cert, _ = tls.X509KeyPair(certPEM, keyPEM)
	err = ioutil.WriteFile(certFile, certPEM, 0400)
	if err == nil {
		err = ioutil.WriteFile(keyFile, keyPEM, 0400)
	}
	return cert, err
}
