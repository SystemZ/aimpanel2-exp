package cmd

import (
	"context"
	"crypto/tls"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/lib/ahttp"
	"net"
	"net/http"
)

func init() {
	rootCmd.AddCommand(fingerPrintCmd)
}

var fingerPrintCmd = &cobra.Command{
	Use:     "fingerprint [addr]",
	Short:   "Get fingerprint for specific address",
	Example: "fingerprint https://127.0.0.1:3000",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("You can use bash for this too:")
		logrus.Info("openssl x509 -noout -in crt.pem -fingerprint -sha256")
		logrus.Info("echo | openssl s_client -connect example.com:443 |& openssl x509 -fingerprint -sha256 -noout")
		client := &http.Client{}
		client.Transport = &http.Transport{
			DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				conn, err := tls.Dial(network, addr, &tls.Config{InsecureSkipVerify: true})
				if err != nil {
					return conn, err
				}
				connState := conn.ConnectionState()
				for _, peerCert := range connState.PeerCertificates {
					logrus.Info(peerCert.Issuer)
					hash := ahttp.GenerateCertFingerprint(peerCert.Raw)
					logrus.Info("Fingerprint: " + hash)
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
