package pki

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"strings"
	"time"
)

const (
	rootKeySize           = 4096
	subKeySize            = 4096
	endEntityKeySize      = 2048
	defaultOrg            = "My Org"
	defaultCountry        = "US"
	defaultProvince       = "CA"
	defaultLocality       = "San Francisco"
	defaultStreet         = "Market Street"
	defaultPostalCode     = "94105"
	defaultCommonName     = "CryptKeeper"
	rootCAValidity        = 10 * 365 * 24 * time.Hour // 10 years
	subCAValidity         = 5 * 365 * 24 * time.Hour  // 5 years
	defaultEntityValidity = 1 * 365 * 24 * time.Hour  // 1 year
)

func GenerateRootCA(caCert *x509.Certificate, caKey *rsa.PrivateKey, caAuthority string) (*x509.Certificate, *rsa.PrivateKey, error) {
	rootKey, err := rsa.GenerateKey(rand.Reader, rootKeySize)
	if err != nil {
		return nil, nil, err
	}

	rootCertTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:       caCert.Issuer.Organization,
			OrganizationalUnit: caCert.Issuer.OrganizationalUnit,
			Country:            caCert.Issuer.Country,
			Province:           caCert.Issuer.Province,
			Locality:           caCert.Issuer.Locality,
			StreetAddress:      caCert.Issuer.StreetAddress,
			PostalCode:         caCert.Issuer.PostalCode,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(rootCAValidity),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	rootCertBytes, err := x509.CreateCertificate(rand.Reader, rootCertTemplate, caCert, &rootKey.PublicKey, caKey)
	if err != nil {
		return nil, nil, err
	}

	rootCert, err := x509.ParseCertificate(rootCertBytes)
	if err != nil {
		return nil, nil, err
	}

	return rootCert, rootKey, nil
}

func GenerateSubCA(rootCert *x509.Certificate, rootKey *rsa.PrivateKey, path string) (*x509.Certificate, *rsa.PrivateKey, error) {
	subKey, err := rsa.GenerateKey(rand.Reader, subKeySize)
	if err != nil {
		return nil, nil, err
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, nil, err
	}

	subCertTemplate := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization:       rootCert.Issuer.Organization,
			OrganizationalUnit: []string{path},
			Country:            rootCert.Issuer.Country,
			Province:           rootCert.Issuer.Province,
			Locality:           rootCert.Issuer.Locality,
			StreetAddress:      rootCert.Issuer.StreetAddress,
			PostalCode:         rootCert.Issuer.PostalCode,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(subCAValidity),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            0,
		MaxPathLenZero:        true,
	}

	subCertBytes, err := x509.CreateCertificate(rand.Reader, subCertTemplate, rootCert, &subKey.PublicKey, rootKey)
	if err != nil {
		return nil, nil, err
	}

	subCert, err := x509.ParseCertificate(subCertBytes)
	if err != nil {
		return nil, nil, err
	}

	return subCert, subKey, nil
}

func GenerateEndEntityCertificateNew(ca *x509.Certificate, caPrivateKey *rsa.PrivateKey, commonName string, organizationName string, dnsNames, ipAddresses, emailAddresses []string, expiresAt time.Time) (*x509.Certificate, *rsa.PrivateKey, error) {
	certPrivateKey, err := rsa.GenerateKey(rand.Reader, endEntityKeySize)
	if err != nil {
		return nil, nil, err
	}

	ipAddressesFmt := []net.IP{}
	for _, ip := range ipAddresses {
		parsedIP := net.ParseIP(strings.Trim(ip, " "))
		if parsedIP != nil {
			ipAddressesFmt = append(ipAddressesFmt, parsedIP)
		}
	}
	// Create SAN extension in DER format
	sanExtension, err := createSANExtension(dnsNames, ipAddressesFmt, emailAddresses)
	if err != nil {
		return nil, nil, err
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, nil, err
	}

	entityCertTemplate := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: commonName,
			// Organization: []string{organizationName},
			Organization: func() []string {
				if organizationName != "" {
					return []string{organizationName}
				}
				return nil
			}(),
		},
		NotBefore:   time.Now().Add(-5 * time.Minute), //Ensure the NotBefore Date Accounts for Clock Skew
		NotAfter:    expiresAt,
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		IsCA:        false,

		// Add SAN extension here
		ExtraExtensions: []pkix.Extension{sanExtension},
		DNSNames:        dnsNames,
		IPAddresses:     ipAddressesFmt,
		EmailAddresses:  emailAddresses,
	}

	entityCertBytes, err := x509.CreateCertificate(rand.Reader, entityCertTemplate, ca, &certPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		return nil, nil, err
	}

	entityCert, err := x509.ParseCertificate(entityCertBytes)
	if err != nil {
		return nil, nil, err
	}

	return entityCert, certPrivateKey, nil
}

func GenerateCABundle(subCACert *x509.Certificate, rootCACert *x509.Certificate) ([]byte, error) {
	var bundle bytes.Buffer

	// Encode Sub CA Certificate in PEM format
	err := pem.Encode(&bundle, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: subCACert.Raw,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to encode Sub CA certificate: %v", err)
	}

	// Encode Root CA Certificate in PEM format
	err = pem.Encode(&bundle, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: rootCACert.Raw,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to encode Root CA certificate: %v", err)
	}

	return bundle.Bytes(), nil
}

// Helper function to create a SAN extension
func createSANExtension(dnsNames []string, ipAddresses []net.IP, emailAddresses []string) (pkix.Extension, error) {
	var rawValues []asn1.RawValue

	// Add DNS Names to SAN extension
	for _, dnsName := range dnsNames {
		rawValues = append(rawValues, asn1.RawValue{
			Class: asn1.ClassContextSpecific,
			Tag:   2, // Tag 2 for DNS names
			Bytes: []byte(dnsName),
		})
	}

	// Add IP Addresses to SAN extension
	for _, ip := range ipAddresses {
		rawValues = append(rawValues, asn1.RawValue{
			Class: asn1.ClassContextSpecific,
			Tag:   7, // Tag 7 for IP addresses
			Bytes: ip,
		})
	}

	// Add Email Addresses to SAN extension
	for _, email := range emailAddresses {
		rawValues = append(rawValues, asn1.RawValue{
			Class: asn1.ClassContextSpecific,
			Tag:   1, // Tag 1 for email addresses
			Bytes: []byte(email),
		})
	}

	// Marshal SAN values to DER format
	sanBytes, err := asn1.Marshal(rawValues)
	if err != nil {
		return pkix.Extension{}, err
	}

	// Return SAN extension
	return pkix.Extension{
		Id:       asn1.ObjectIdentifier{2, 5, 29, 17}, // OID for Subject Alternative Name
		Critical: false,
		Value:    sanBytes,
	}, nil
}

// func GenerateEndEntityCertificate(subCert *x509.Certificate, subKey *rsa.PrivateKey, commonName string, expiresAt *time.Time) (*x509.Certificate, *rsa.PrivateKey, error) {
// 	entityKey, err := rsa.GenerateKey(rand.Reader, endEntityKeySize)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	notAfter := time.Now().Add(defaultEntityValidity)
// 	if expiresAt != nil {
// 		notAfter = *expiresAt
// 	}

// 	entityCertTemplate := &x509.Certificate{
// 		SerialNumber: big.NewInt(time.Now().UnixNano()),
// 		Subject: pkix.Name{
// 			CommonName:   commonName,
// 			Organization: []string{defaultCommonName},
// 		},
// 		NotBefore:   time.Now(),
// 		NotAfter:    notAfter,
// 		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
// 		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
// 		DNSNames:    []string{"localhost", defaultCommonName}, // Add localhost here
// 		IsCA:        false,
// 	}

// 	entityCertBytes, err := x509.CreateCertificate(rand.Reader, entityCertTemplate, subCert, &entityKey.PublicKey, subKey)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	entityCert, err := x509.ParseCertificate(entityCertBytes)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return entityCert, entityKey, nil
// }
