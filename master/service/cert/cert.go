package cert

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns/cloudflare"
	"github.com/go-acme/lego/v4/registration"
	"gitlab.com/systemz/aimpanel2/master/config"
	"gitlab.com/systemz/aimpanel2/master/model"
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

func InitLego() error {
	//Create cloudflare provider
	cfProvider, err := createCloudflareProvider()
	if err != nil {
		return err
	}

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	User = LeUser{
		Email: config.LE_EMAIL,
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
		return err
	}

	err = Client.Challenge.SetDNS01Provider(cfProvider)

	//Register client
	reg, err := Client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return err
	}

	User.Registration = reg

	return nil
}

func createCloudflareProvider() (*cloudflare.DNSProvider, error) {
	if config.CLOUDFLARE_AUTH_TOKEN == "" || config.CLOUDFLARE_ZONE_TOKEN == "" {
		return nil, errors.New("CLOUDFLARE_AUTH_TOKEN or CLOUDFLARE_ZONE_TOKEN not set")
	}

	cloudflareCfg := cloudflare.NewDefaultConfig()
	cloudflareCfg.AuthToken = config.CLOUDFLARE_AUTH_TOKEN
	cloudflareCfg.ZoneToken = config.CLOUDFLARE_ZONE_TOKEN

	provider, err := cloudflare.NewDNSProviderConfig(cloudflareCfg)
	if err != nil {
		return nil, err
	}

	return provider, nil
}

func CreateCertForDomain(domain model.CertDomain) error {
	request := certificate.ObtainRequest{
		Domains: []string{domain.Name},
		Bundle:  false,
	}

	cert, err := Client.Certificate.Obtain(request)
	if err != nil {
		return err
	}

	dbCert := &model.Cert{
		Cert:       string(cert.Certificate),
		PrivateKey: string(cert.PrivateKey),
		DomainId:   domain.ID,
	}

	err = model.Put(dbCert)
	if err != nil {
		return err
	}

	return nil
}

func CreateCerts() error {
	domains, err := model.GetCertDomains()
	if err != nil {
		return err
	}

	for _, domain := range domains {
		if err := CreateCertForDomain(domain); err != nil {
			return err
		}
	}

	return nil
}
