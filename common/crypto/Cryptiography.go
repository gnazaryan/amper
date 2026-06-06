package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"time"
)

// GenerateKeyPair generates a new key pair
func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil
	}
	return privkey, &privkey.PublicKey
}

// PrivateKeyToBytes private key to bytes
func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	if priv == nil {
		return nil
	}

	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	return privBytes
}

// PublicKeyToBytes public key to bytes
func PublicKeyToBytes(pub *rsa.PublicKey) []byte {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	if pub == nil {
		return nil
	}

	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) []byte {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, cryptorand.Reader, priv, ciphertext, nil)
	if err != nil {
		return nil
	}
	return plaintext
}

// BytesToPrivateKey bytes to private key
func BytesToPrivateKey(priv []byte) *rsa.PrivateKey {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		return nil
	}
	return key
}

//CreateHash -
func CreateHash(key string) string {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

//CreateHashByte -
func CreateHashByte(key []byte) string {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	hasher := md5.New()
	hasher.Write(key)
	return hex.EncodeToString(hasher.Sum(nil))
}

//UUID -
func UUID() string {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	uuid := make([]byte, 16)
	n, err := io.ReadFull(cryptorand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return ""
	}
	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40
	result := fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
	return result
}

//Decrypt -
func Decrypt(data []byte, passphrase string) []byte {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	if data == nil {
		return nil
	}

	key := []byte(CreateHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil
	}
	return plaintext
}

//EncryptWithPublicKey -
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) []byte {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	if pub == nil {
		return nil
	}

	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {

		return nil
	}
	return ciphertext
}

//Passphrase to encrypt data
var Passphrase = "amper_passpphrase"

//Encrypt -
func Encrypt(data []byte, passphrase string) []byte {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	block, _ := aes.NewCipher([]byte(CreateHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

// BytesToPublicKey bytes to public key
func BytesToPublicKey(pub []byte) *rsa.PublicKey {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)

	b := block.Bytes
	var err error
	if enc {
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			os.Exit(2)
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		os.Exit(2)
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		os.Exit(2)
	}
	return key
}

//ByteCountSI -
func ByteCountSI(b int64) string {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

//ByteCountIEC -
func ByteCountIEC(b int64) string {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

//Certificate -
type Certificate struct {
	CACert     []byte
	CAKey      []byte
	ServerCert []byte
	ServerKey  []byte
}

//GenerateCA -
func GenerateCA(keySize int) (cert Certificate, err error) {
	// step: generate a serial number
	serial, err := rand.Int(rand.Reader, (&big.Int{}).Exp(big.NewInt(2), big.NewInt(159), nil))
	if err != nil {
		return cert, err
	}

	ca := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			Organization:       []string{"Afis.me"},
			OrganizationalUnit: []string{"Certificate Authority"},
			Country:            []string{"ID"},
			Province:           []string{"East Java"},
			Locality:           []string{"Surabaya"},
			CommonName:         "afis.me",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// create our private and public key
	caPrivKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return cert, err
	}

	// create the CA
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return cert, err
	}

	// pem encode
	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})

	return Certificate{
		CACert: caPEM.Bytes(),
		CAKey:  caPrivKeyPEM.Bytes(),
	}, nil

}

/*//GenerateCert -
func GenerateCert(keySize int, domain []string, alips []string) (cert Certificate, err error) {

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(caCert))
	if !ok {
		panic("failed to parse root certificate")
	}

	block, _ := pem.Decode([]byte(caCert))
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	cacert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}

	ca := &x509.Certificate{
		SerialNumber:          cacert.SerialNumber,
		Subject:               cacert.Subject,
		NotBefore:             cacert.NotBefore,
		NotAfter:              cacert.NotAfter,
		IsCA:                  true,
		ExtKeyUsage:           cacert.ExtKeyUsage,
		KeyUsage:              cacert.KeyUsage,
		BasicConstraintsValid: true,
	}

	serial, err := rand.Int(rand.Reader, (&big.Int{}).Exp(big.NewInt(2), big.NewInt(159), nil))
	if err != nil {
		return cert, err
	}

	ipstemp := []net.IP{net.IPv4(127, 0, 0, 1)}
	// get all ip address in this server
	if val, err := psnet.Interfaces(); err == nil {
		for _, s := range val {
			if s.Name != "lo" {
				for _, sv := range s.Addrs {
					if ii, _, err := net.ParseCIDR(sv.Addr); err == nil {
						if ip := ii.To4(); ip != nil {
							ipstemp = append(ipstemp, ip)
						}
					}
				}
			}
		}
	}

	for _, sv := range alips {
		if ii, _, err := net.ParseCIDR(sv); err == nil {
			if ip := ii.To4(); ip != nil {
				ipstemp = append(ipstemp, ip)
			}
		} else {
			if ipv := net.ParseIP(sv); ipv != nil {
				ipstemp = append(ipstemp, ipv)
			}
		}
	}

	//check duplicating ips
	var listips []string
	for _, s := range ipstemp {
		listips = append(listips, s.String())
	}

	ips := []net.IP{}
	for _, s := range utility.UniqueString(listips) {
		if v := net.ParseIP(s); v != nil {
			ips = append(ips, v)
		}
	}

	// set up our server certificate
	clcert := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			Organization:       []string{"Afis.me"},
			OrganizationalUnit: []string{"Trading apps"},
			Country:            []string{"ID"},
			Province:           []string{"East Java"},
			Locality:           []string{"Surabaya"},
			CommonName:         "self-signed.afis.me",
		},
		DNSNames:     utility.UniqueString(domain),
		IPAddresses:  ips,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return cert, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, clcert, ca, &certPrivKey.PublicKey, BytesToPrivateKey([]byte(caKeys)))
	if err != nil {
		return cert, err
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	return Certificate{
		CACert:     []byte(caCert),
		ServerCert: certPEM.Bytes(),
		ServerKey:  certPrivKeyPEM.Bytes(),
	}, nil

}*/
