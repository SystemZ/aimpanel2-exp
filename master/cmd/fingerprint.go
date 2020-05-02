package cmd

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"log"
	"net"
	"net/http"
)

func init() {
	rootCmd.AddCommand(fingerPrintCmd)
}

var fingerPrintCmd = &cobra.Command{
	Use:     "fingerprint [addr]",
	Short:   "Get fingerprint for specific address",
	Long:    "",
	Example: "fingerprint https://aimpanel.local",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := &http.Client{}
		client.Transport = &http.Transport{
			DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				conn, err := tls.Dial(network, addr, &tls.Config{InsecureSkipVerify: true})
				if err != nil {
					return conn, err
				}
				connState := conn.ConnectionState()
				for _, peerCert := range connState.PeerCertificates {
					der, err := x509.MarshalPKIXPublicKey(peerCert.PublicKey)
					if err != nil {
						log.Fatal(err)
					}
					hash := sha256.Sum256(der)
					logrus.Info(peerCert.Issuer)
					logrus.Info("Fingerprint: " + hex.EncodeToString(hash[:]))
					logrus.Info("")
				}
				return conn, nil
			},
		}

		_, err := client.Get(args[0])
		if err != nil {
			logrus.Error(err)
		}
	},
}
