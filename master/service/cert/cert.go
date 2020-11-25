package cert

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/providers/dns/cloudflare"
	"github.com/go-acme/lego/v4/registration"
	"gitlab.com/systemz/aimpanel2/master/config"
	"gitlab.com/systemz/aimpanel2/master/model"
	"io"
	"math/big"
	"net"
	"os"
	"time"
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

	c, err := aes.NewCipher([]byte(config.LE_SECRET))
	if err != nil {
		return err
	}

	aesGCM, err := cipher.NewGCM(c)
	if err != nil {
		return err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	encryptedPrivateKey := aesGCM.Seal(nonce, nonce, cert.PrivateKey, nil)

	dbCert := &model.Cert{
		Cert:       string(cert.Certificate),
		PrivateKey: fmt.Sprintf("%x", encryptedPrivateKey),
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

//https://gist.github.com/samuel/8b500ddd3f6118d052b5e6bc16bc4c09
func CreateSelfSignedCert(domain model.CertDomain) error {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Aimpanel"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 365),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	if ip := net.ParseIP(domain.Name); ip != nil {
		template.IPAddresses = append(template.IPAddresses, ip)
	} else {
		template.DNSNames = append(template.DNSNames, domain.Name)
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(privateKey), privateKey)
	if err != nil {
		return err
	}

	dbCert := &model.Cert{
		DomainId: domain.ID,
	}

	out := &bytes.Buffer{}

	pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	dbCert.Cert = out.String()

	out.Reset()

	pem.Encode(out, pemBlockForKey(privateKey))
	dbCert.PrivateKey = out.String()

	err = model.Put(dbCert)
	if err != nil {
		return err
	}

	return nil
}

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}
