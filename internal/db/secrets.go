package db

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cryptkeeperhq/cryptkeeper/internal/crypt"
	"github.com/cryptkeeperhq/cryptkeeper/internal/models"
	"github.com/cryptkeeperhq/cryptkeeper/internal/utils"
	"github.com/go-pg/pg/v10"
	"github.com/tink-crypto/tink-go/v2/aead"
	"github.com/tink-crypto/tink-go/v2/keyset"

	enginedb "github.com/cryptkeeperhq/cryptkeeper/internal/engine/database"
	enginepki "github.com/cryptkeeperhq/cryptkeeper/internal/engine/pki"
	enginetransit "github.com/cryptkeeperhq/cryptkeeper/internal/engine/transit"
)

func GetTransitEncryptionKey(keyId string, version int) (models.Secret, error) {
	var secret models.Secret

	query := DB.Model(&secret).Distinct().
		ColumnExpr("secret.*, p.path, secret.metadata->>'key_type' AS key_type").
		Join("JOIN paths p ON p.id = secret.path_id").
		Where("secret.id = ?", keyId).
		Order("version DESC")

	if version != 0 {
		query.Where("version = ?", version)
	}

	err := query.Limit(1).Select()

	return secret, err
}

func GetSecretByID(id string, version int64) (models.Secret, error) {

	var secret models.Secret
	query := DB.Model(&secret).
		Column("secret.*", "p.path").
		Join("JOIN paths p ON p.id = secret.path_id").
		Where("secret.id = ?", id).
		Order("version DESC")

	if version != 0 {
		query.Where("version = ?", version)
	}

	err := query.Limit(1).Select()

	return secret, err

}

// func GetSecretByPathID(pathID string, key string, version int64) (models.Secret, error) {

// 	var secret models.Secret
// 	query := DB.Model(&secret).
// 		Column("secret.*", "p.path").
// 		Join("JOIN paths p ON p.id = secret.path_id").
// 		Where("secret.key = ?", key).
// 		Where("p.id = ?", pathID).
// 		Order("version DESC")

// 	if version != 0 {
// 		query.Where("version = ?", version)
// 	}

// 	err := query.Limit(1).Select()

// 	return secret, err

// }

func GetSecret(path string, key string, version int64) (models.Secret, error) {
	var secret models.Secret
	query := DB.Model(&secret).
		Column("secret.*", "p.path").
		Join("JOIN paths p ON p.id = secret.path_id").
		Where("secret.key = ?", key).
		Where("p.path = ?", path).
		Order("version DESC")

	if version != 0 {
		query.Where("version = ?", version)
	}

	err := query.Limit(1).Select()

	return secret, err

}

func DeleteSecret(secret models.Secret) error {
	err := DB.RunInTransaction(context.Background(), func(tx *pg.Tx) error {

		// Track the deleted secret and its access control entries
		deletion := models.SecretDeletion{
			SecretID:       secret.ID,
			PathID:         secret.PathID,
			Key:            secret.Key,
			Version:        secret.Version,
			EncryptedDEK:   secret.EncryptedDEK,
			EncryptedValue: secret.EncryptedValue,
			Metadata:       secret.Metadata,
			DeletedAt:      time.Now(),
		}

		_, err := DB.Model(&deletion).Insert()
		if err != nil {
			return err
		}

		_, err = tx.Model(&secret).Where("id = ? and version = ?", secret.ID, secret.Version).Delete()
		return err

	})

	return err
}

func WriteSecret(userID string, request models.Secret, c crypt.CryptographicOperations) (models.Secret, error) {

	// Fetch the path to determine the engine type
	var path models.Path
	err := DB.Model(&path).Where("id = ?", request.PathID).Select()
	if err != nil {
		return models.Secret{}, err
	}

	// Decrypt the path key
	decryptedPathKeyHandle, err := c.DecryptPathKey(path.KeyData)
	if err != nil {
		log.Printf("failed to decrypt path key: %v\n", err)
		return models.Secret{}, err
	}

	secretEngine := path.EngineType
	// Handle secret creation based on engine type
	switch secretEngine {
	case "kv":
		request.EncryptedDEK, request.EncryptedValue, err = c.EncryptSecretValue(request.Value, decryptedPathKeyHandle)
		if err != nil {
			log.Fatalf("failed to encrypt secret value: %v", err)
		}
	case "pki":
		// Create PKI secret
		var subCA models.SubCA
		err := DB.Model(&subCA).Where("path_id = ?", path.ID).First()
		if err != nil {
			return models.Secret{}, errors.New("subCA not found")
		}

		// if path.Metadata["root_ca"].(string) == "cryptkeeper_ca" {
		// 	err := DB.Model(&subCA).Where("path_id = ?", path.ID).First()
		// 	if err != nil {
		// 		return models.Secret{}, errors.New("subCA not found")
		// 	}
		// } else if path.Metadata["root_ca"].(string) == "lets_encrypt_staging" {
		// 	// Below values should come from metadata
		// 	certmagic.DefaultACME.Agreed = true
		// 	certmagic.DefaultACME.Email = "vdparikh@gmail.com"
		// 	certmagic.DefaultACME.CA = certmagic.LetsEncryptStagingCA // Use staging for testing

		// 	certCfg := certmagic.NewDefault()
		// 	err := certCfg.ManageSync(context.Background(), []string{request.Key})
		// 	if err != nil {
		// 		return models.Secret{}, err
		// 	}

		// 	certResource, err := certCfg.CacheManagedCertificate(context.Background(), request.Key)
		// 	if err != nil {
		// 		return models.Secret{}, err
		// 	}

		// 	subCA.SubCACert = []byte(utils.ConvertCertificateToPEM(certResource.Certificate.Leaf.Raw))
		// 	privateKey, _ := utils.ConvertPrivateKeyToPEM(certResource.PrivateKey)
		// 	subCA.SubCAKey = []byte(privateKey)
		// 	// subCA.SubCACert = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certResource.Certificate.Leaf.Raw})
		// 	// subCA.SubCAKey = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(certResource.PrivateKey.(*rsa.PrivateKey))})
		// } else {
		// 	fmt.Println("Genertaing CA for ", path.Metadata["root_ca"].(string))
		// 	var ca models.CertificateAuthority
		// 	if err := DB.Model(&ca).Where("id = ?", utils.ToInt(path.Metadata["root_ca"].(string))).Select(); err != nil {
		// 		return models.Secret{}, errors.New("invalid root ca id")
		// 	}

		// 	// Decode PEM CA certificate
		// 	block, _ := pem.Decode([]byte(ca.CACert))
		// 	if block == nil || block.Type != "CERTIFICATE" {
		// 		return models.Secret{}, errors.New("failed to decode PEM block containing the certificate")
		// 	}

		// 	subCA.SubCACert = block.Bytes

		// 	// Decode PEM CA private key
		// 	block, _ = pem.Decode([]byte(ca.PrivateKey))
		// 	if block == nil {
		// 		return models.Secret{}, errors.New("failed to decode PEM block containing the private key")

		// 	}
		// 	subCA.SubCAKey = block.Bytes
		// }

		subCert, err := x509.ParseCertificate(subCA.SubCACert)
		if err != nil {
			return models.Secret{}, errors.New("invalid SubCA Certificate")
		}

		// subKey, err := x509.ParsePKCS1PrivateKey(subCA.SubCAKey)
		// if err != nil {
		// 	return models.Secret{}, errors.New("invalid SubCA Private Key")
		// }

		var subKey interface{}
		subKey, err = x509.ParsePKCS1PrivateKey(subCA.SubCAKey)
		if err != nil {
			subKey, err = x509.ParsePKCS8PrivateKey(subCA.SubCAKey)
			if err != nil {
				log.Println("Invalid CA Private ", err.Error())
				return models.Secret{}, errors.New("invalid CA Private ")
			}
		}

		//TODO: Update to match the path or secret data
		// leaseDuration := request.Metadata["max_lease_time"].(string)
		expiresAt := time.Now().AddDate(0, 0, 365)
		if val, ok := request.Metadata["validity_period"]; ok {
			// TODO: More logic to check against path max_lease_time
			expiresAt = time.Now().AddDate(0, 0, utils.ToInt(val.(string)))
		}

		request.ExpiresAt = &expiresAt

		ou := path.Path
		if val, ok := request.Metadata["organization"]; ok {
			// TODO: More logic to check against path max_lease_time
			ou = val.(string)
		}

		// Define SANs: DNS Names, IP Addresses, Email Addresses
		// dnsNames := []string{"www.example.com", "example.org"}
		// ipAddresses := []net.IP{net.ParseIP("192.168.1.1"), net.ParseIP("10.0.0.1")}
		// emailAddresses := []string{"admin@example.com"}

		dnsNames := []string{}
		if val, ok := request.Metadata["dns_names"]; ok {
			dnsNames = strings.Split(val.(string), ",")
		}

		ipAddresses := []string{}
		if val, ok := request.Metadata["ip_addresses"]; ok {
			ipAddresses = strings.Split(val.(string), ",")
		}

		emailAddresses := []string{}
		if val, ok := request.Metadata["email_addresses"]; ok {
			emailAddresses = strings.Split(val.(string), ",")
		}

		fmt.Println(dnsNames, ipAddresses, emailAddresses)
		entityCert, entityKey, err := enginepki.GenerateEndEntityCertificateNew(subCert, subKey.(*rsa.PrivateKey), request.Key, ou, dnsNames, ipAddresses, emailAddresses, expiresAt)
		if err != nil {
			fmt.Println(err)
			return models.Secret{}, errors.New("failed to generate end-entity certificate")
		}

		// Private Key: This is the private part of the certificate that should be kept secure. It is used to sign data and decrypt information encrypted with the public certificate.
		// Private Key: Users use the private key to sign data or decrypt information. For example, in a web server scenario, the server uses the private key to decrypt incoming HTTPS traffic.
		// certPEM := utils.ConvertCertificateToPEM(entityCert.Raw)
		certPEM := entityCert.Raw
		// privateKeyPEM, err := utils.ConvertPrivateKeyToPEM(entityKey)
		// if err != nil {
		// 	return models.Secret{}, err
		// }

		privateKeyPEMBytes := x509.MarshalPKCS1PrivateKey(entityKey)
		request.Value = base64.StdEncoding.EncodeToString(privateKeyPEMBytes)

		request.EncryptedDEK, request.EncryptedValue, err = c.EncryptSecretValue(request.Value, decryptedPathKeyHandle)
		if err != nil {
			return models.Secret{}, err
		}

		// Public Certificate (entityCert.Raw): This is the public part of the certificate that can be shared with others. It is used to verify the identity of the certificate holder.
		// Public Certificate: Users can install the certificate on their systems or applications that need to verify their identity. For example, a web server can use the certificate to establish HTTPS connections.
		// metadata := make(map[string]interface{})
		request.Metadata["public_key"] = base64.StdEncoding.EncodeToString([]byte(certPEM))
		request.Metadata["expiration"] = entityCert.NotAfter
		request.Metadata["subject"] = entityCert.Subject
		request.Metadata["issuer"] = entityCert.Issuer
		request.Metadata["serial_number"] = entityCert.SerialNumber.String()

		// request.Metadata = metadata

	case "database":
		// Create role in the database using the connection string and role template
		user := path.Metadata["username"].(string)
		password := path.Metadata["password"].(string)

		connectionString := path.Metadata["connection_string"].(string)
		connectionString = strings.ReplaceAll(connectionString, "{{username}}", user)
		connectionString = strings.ReplaceAll(connectionString, "{{password}}", password)

		roleTemplate := path.Metadata["role_template"].(string)

		fmt.Println(roleTemplate)

		// Generate database credentials
		var dbUsername, dbPassword string

		dbUsername = utils.GenerateUsername()

		// TODO: FIX THIS FOR PG USER ROTATION
		// if request.ID != "" {
		// 	// This is triggered for first time creation.
		// 	// plaintextValue, err := utils.DecryptSecretValue(request, encryptionKey)
		// 	plaintextValue, err := c.DecryptSecretValue(request.EncryptedDEK, request.EncryptedValue, decryptedPathKeyHandle)

		// 	if err != nil {
		// 		return models.Secret{}, err
		// 	}

		// 	credSegements := strings.Split(string(plaintextValue), ":")
		// 	dbUsername = credSegements[0]
		// 	roleTemplate = "ALTER ROLE {{name}} WITH PASSWORD '{{password}}';"
		// } else {
		// 	dbUsername = utils.GenerateUsername()
		// }

		dbPassword = utils.GeneratePassword()
		// TODO: Move this line to transaction as we should only create role in DB if the secret creation is successful.
		err = enginedb.CreateRoleInDatabase(connectionString, roleTemplate, dbUsername, dbPassword)
		if err != nil {
			return models.Secret{}, err
		}

		request.Value = fmt.Sprintf("%s:%s", dbUsername, dbPassword)

		// Encrypt and store the credentials
		request.EncryptedDEK, request.EncryptedValue, err = c.EncryptSecretValue(request.Value, decryptedPathKeyHandle)

		// encryptedValue, err := keymanagement.Encrypt(masterKey, []byte(fmt.Sprintf("%s:%s", dbUsername, dbPassword)))
		if err != nil {
			return models.Secret{}, err
		}
	case "transit":
		// Create Transit secret
		keyType, ok := request.Metadata["key_type"].(string)
		if !ok {
			return models.Secret{}, errors.New("failed getting key_type")
		}

		// Generate key based on key type
		transitHandler := enginetransit.Handler{
			// PathKeyHandle: decryptedPathKeyHandle,
			KeyType: keyType,
		}

		keyTemplate, err := transitHandler.GenerateKeyTemplate()
		if err != nil {
			return models.Secret{}, fmt.Errorf("error generating template: %v", err)
		}

		// templateString, err := proto.Marshal(keyTemplate)

		// if err != nil {
		// 	fmt.Println("ERROR!!!", "failed to Marshal KeyTemplate")
		// 	return models.Secret{}, err
		// }

		// request.EncryptedDEK, request.EncryptedValue, err = c.EncryptSecretValue(string(templateString), decryptedPathKeyHandle)

		// Encrypt the template using Path Key
		pathAead, _ := aead.New(decryptedPathKeyHandle)
		kh, err := keyset.NewHandle(keyTemplate)
		if err != nil {
			return models.Secret{}, fmt.Errorf("error generating key handle: %v", err)
		}

		var buffer bytes.Buffer
		writer := keyset.NewBinaryWriter(&buffer)
		if err := kh.Write(writer, pathAead); err != nil {
			return models.Secret{}, err
		}

		serializedKey := buffer.Bytes()
		request.EncryptedDEK, err = pathAead.Encrypt(serializedKey, nil)

		if err != nil {
			return models.Secret{}, err
		}

	default:
		return models.Secret{}, errors.New("unsupported engine")
	}

	// Determine the next version
	var existingSecret models.Secret
	err = DB.Model(&existingSecret).
		Column("version").
		Where("path_id = ? AND key = ?", request.PathID, request.Key).
		Order("version DESC").
		Limit(1).
		Select()

	if err != nil && err != pg.ErrNoRows {
		return request, err
	}

	// secretUUID := ""

	// if existingSecret.ID == "" {
	// 	secretUUID = uuid.New().String()
	// }

	// Create the secret
	now := time.Now()
	secret := models.Secret{
		ID:               request.ID,
		PathID:           request.PathID,
		Key:              request.Key,
		Version:          existingSecret.Version + 1,
		EncryptedDEK:     request.EncryptedDEK,
		EncryptedValue:   request.EncryptedValue,
		Metadata:         request.Metadata,
		IsOneTime:        request.IsOneTime,
		Checksum:         request.Checksum,
		ExpiresAt:        request.ExpiresAt,
		Tags:             request.Tags,
		RotationInterval: request.RotationInterval,
		IsMultiValue:     request.IsMultiValue,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		CreatedBy:        request.CreatedBy,
		LastRotatedAt:    &now,
	}

	err = DB.RunInTransaction(context.Background(), func(tx *pg.Tx) error {
		// Insert the new secret
		_, err := tx.Model(&secret).Insert()
		if err != nil {
			return err
		}

		// // Assign access control to the user
		// userAccess := models.SecretAccess{
		// 	SecretID:    secret.ID,
		// 	UserID:      &userID,
		// 	AccessLevel: "owner",
		// }
		// _, err = tx.Model(&userAccess).Insert()
		// if err != nil {
		// 	return err
		// }

		// // Create keys for transit engine
		// if secretEngine == "transit" {
		// 	keyType, ok := request.Metadata["key_type"].(string)
		// 	if !ok {
		// 		return errors.New("invalid key type in metadata")
		// 	}

		// 	// Generate key based on key type
		// 	// TODO: Should we encrypt the transit keys as well using the Path?
		// 	generatedKey, err := tink.GenerateKeyData(keyType, decryptedPathKeyHandle)
		// 	fmt.Println(string(generatedKey))

		// 	// encryptedDEK, err := crypto.EncryptDEK(pathKey, dek)
		// 	// if err != nil {
		// 	// 	return []byte{}, []byte{}, err
		// 	// }

		// 	// // Base64 encode the encrypted DEK
		// 	// encryptedDEKBase64 := base64.StdEncoding.EncodeToString(encryptedDEK)

		// 	if err != nil {
		// 		return err
		// 	}

		// 	initialKeyVersion := models.TransitKeyVersion{
		// 		KeyID:   int(secret.ID),
		// 		Version: secret.Version,
		// 		KeyData: generatedKey,
		// 		// KeyDataString: base64.StdEncoding.EncodeToString(generatedKey),
		// 		KeyType:   keyType,
		// 		CreatedAt: time.Now(),
		// 		UpdatedAt: time.Now(),
		// 	}

		// 	_, err = DB.Model(&initialKeyVersion).Insert()
		// 	if err != nil {
		// 		return err
		// 	}

		// 	// pubKey := utils.GetPublicKey(generatedKey)
		// 	// secret.Metadata = metadata
		// 	// secret.UpdatedAt = time.Now()

		// 	// _, err = h.DB.Model(&secret).WherePK().Update()
		// 	// if err != nil {
		// 	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	// 	return
		// 	// }

		// }

		// Assign access control to the user's groups
		// groupIDs, err := h.getUserGroups(userID)
		// var groupIDs []int64
		// err = DB.Model((*models.UserGroup)(nil)).
		// 	Column("group_id").
		// 	Where("user_id = ?", userID).
		// 	Select(&groupIDs)

		// if err != nil {
		// 	return err
		// }

		// for _, groupID := range groupIDs {
		// 	groupAccess := models.SecretAccess{
		// 		SecretID:    secret.ID,
		// 		GroupID:     &groupID,
		// 		AccessLevel: "member",
		// 	}
		// 	_, err = tx.Model(&groupAccess).Insert()
		// 	if err != nil {
		// 		return err
		// 	}
		// }

		return nil
	})

	if err != nil {
		return request, err
	}

	return secret, err
}

func GetSecrets(pathId int) ([]models.Secret, error) {
	var secrets []models.Secret
	err := DB.Model(&secrets).Distinct().
		// Join("JOIN secret_accesses sa ON sa.secret_id = secret.id").
		Where("secret.path_id = ?", pathId).
		// Where("sa.user_id = ? OR sa.group_id IN (SELECT group_id FROM user_groups WHERE user_id = ?)", userID, userID).
		Select()

	return secrets, err
}

func GetSecretsByPathName(pathName string) ([]models.Secret, error) {
	var secrets []models.Secret

	// Subquery to get the latest version for each secret
	subquery := DB.Model((*models.Secret)(nil)).
		Column("key").
		ColumnExpr("MAX(version) AS max_version").
		Group("key")

	err := DB.Model(&secrets).Distinct().
		Column("secret.*", "p.path").
		Join("JOIN paths p ON p.id = secret.path_id").
		Join("JOIN (?) AS sq ON sq.key = secret.key AND sq.max_version = secret.version", subquery).
		Where("p.path = ? and (secret.expires_at is null OR expires_at > now()::date)", pathName).
		Select()

	return secrets, err
}

// func GetSecretsByPathID(id int64) ([]models.Secret, error) {
// 	var secrets []models.Secret
// 	err := DB.Model(&secrets).Distinct().
// 		Column("secret.*", "p.path").
// 		Join("JOIN paths p ON p.id = secret.path_id").
// 		Where("p.id = ?", id).
// 		// Where("sa.user_id = ? OR sa.group_id IN (SELECT group_id FROM user_groups WHERE user_id = ?)", userID, userID).
// 		Select()

// 	return secrets, err
// }
