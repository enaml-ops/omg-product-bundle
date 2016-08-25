package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/enaml-ops/omg-cli/utils"
)

var (
	host   = flag.String("host", "", "Comma-separated hostnames and IPs to generate a certificate for")
	prefix = flag.String("prefix", "", "filename prefix")
)

func main() {
	flag.Parse()

	if len(*host) == 0 {
		log.Fatalf("Missing required --host parameter")
	}

	if len(*prefix) == 0 {
		log.Fatalf("Missing required --prefix parameter")
	}

	hosts := strings.Split(*host, ",")
	caCert, cert, key, err := utils.GenerateCert(hosts)
	if err != nil {
		log.Fatal("error generating cert:", err)
	}

	certfileName := fmt.Sprintf("%s-cert.pem", *prefix)
	if err := ioutil.WriteFile(certfileName, []byte(cert), 0644); err != nil {
		log.Fatal("error writing cert:", err.Error())
	}
	log.Println("wrote", certfileName)

	keyfileName := fmt.Sprintf("%s-key.pem", *prefix)
	if err := ioutil.WriteFile(keyfileName, []byte(key), 0600); err != nil {
		log.Fatal("error writing key:", err.Error())
	}
	log.Println("wrote", keyfileName)

	cacertfileName := fmt.Sprintf("%s-ca-cert.pem", *prefix)
	if err := ioutil.WriteFile(cacertfileName, []byte(caCert), 0644); err != nil {
		log.Fatal("error writing CA cert:", err.Error())
	}
	log.Println("wrote", cacertfileName)
}
