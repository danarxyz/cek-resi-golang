package utils

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/goodsign/monday"
	"github.com/goravel/framework/facades"
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
	env := facades.Config()

	serverCert, err := tls.LoadX509KeyPair(env.GetString("SERVER_CERT", "certs/server-cert.pem"), env.GetString("SERVER_KEY", "certs/server-key.pem"))
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

func CreateListener(host string) (net.Listener, error) {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalln("Failed to create listener:", err)
	}
	return listener, nil
}
