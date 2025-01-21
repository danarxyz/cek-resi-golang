package utils

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/goodsign/monday"
	"google.golang.org/grpc/credentials"
)

func ParseTime(dateTimeStr string) time.Time {
	layout := "02 Jan 15:04 MST"
	t, err := monday.Parse(layout, dateTimeStr, monday.LocaleIdID) // Gunakan locale Indonesia
	if err != nil {
		fmt.Println(err)
		return time.Time{}
	}
	return t
}

func LoadTLSCredentials() (credentials.TransportCredentials, error) {
	serverCert, err := tls.LoadX509KeyPair("certs/server-cert.pem", "certs/server-key.pem")
	if err != nil {
		log.Fatalln("Failed to read server certificate:", err)
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}
