package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"

	quic "github.com/lucas-clemente/quic-go"
)

const addr = "localhost:4242"

// We start a server echoing data on the first stream the client opens,
// then connect with a client, send the message, and wait for its receipt.
func main() {
	echoServer()

}

// Start a server that echos all data on the first stream opened by the client
func echoServer() error {

	listener, err := quic.ListenAddr(addr, generateTLSConfig(), nil)

	if err != nil {
		return err
	}
	for {
		sess, err := listener.Accept()
		if err != nil {
			return err
		}

		fmt.Println("accepted a connection")
		stream, err := sess.AcceptStream()

		for i := 1; i < 10; i++ {

			buf := make([]byte, 6)

			// fmt.Println("here")
			_, err = io.ReadFull(stream, buf)

			// fmt.Println("here1")
			// message := string(buf[:6])

			_, err = io.Writer.Write(stream, buf)
			if err != nil {
				return err
			}
			// fmt.Println("here2")
		}
	}
	return err
}

// A wrapper for io.Writer that also logs the message.
type loggingWriter struct{ io.Writer }

func (w loggingWriter) Write(b []byte) (int, error) {
	fmt.Printf("Server: Got '%s'\n", string(b))
	return w.Writer.Write(b)
}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}
}
