package cert

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/master/config"
)

var (
	User   LeUser
	Client *lego.Client
)

type LeUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *LeUser) GetEmail() string {
	return u.Email
}
func (u LeUser) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *LeUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func InitLego() {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		logrus.Fatal(err)
	}

	User = LeUser{
		Email: "aimpanel@aimpanel.pro",
		key:   privateKey,
	}

	cfg := lego.NewConfig(&User)
	cfg.Certificate.KeyType = certcrypto.RSA2048

	if config.DEV_MODE {
		cfg.CADirURL = "https://acme-staging-v02.api.letsencrypt.org/directory"
	} else {
		cfg.CADirURL = "https://acme-v02.api.letsencrypt.org/directory"
	}

	//Create lego client
	Client, err = lego.NewClient(cfg)
	if err != nil {
		logrus.Fatal(err)
	}

	//Create http solver
	err = Client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", "5002"))
	if err != nil {
		logrus.Fatal(err)
	}

	//Register client
	reg, err := Client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		logrus.Fatal(err)
	}

	User.Registration = reg
}

func CreateCertForDomain(domain string) {
	request := certificate.ObtainRequest{
		Domains: []string{domain},
		Bundle:  true,
	}

	certificates, err := Client.Certificate.Obtain(request)
	if err != nil {
		logrus.Fatal(err)
	}

	// Each certificate comes back with the cert bytes, the bytes of the client's
	// private key, and a certificate URL. SAVE THESE TO DISK.
	fmt.Printf("%#v\n", certificates)
}
