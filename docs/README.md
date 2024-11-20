# Initialize CryptKeeper


```sh
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout keyfile.pem -out certfile.pem
```

Or use Certbot

```sh
brew install certbot
sudo certbot certonly --standalone -d yourdomain.com
```

Certbot will generate the certificate and key files in the `/etc/letsencrypt/live/yourdomain.com/` directory.

Update TLS config
```yaml
tls:
  cert_file: "/etc/letsencrypt/live/yourdomain.com/fullchain.pem"
  key_file: "/etc/letsencrypt/live/yourdomain.com/privkey.pem"
```


```ssh
openssl pkcs12 -in certificate.p12 -clcerts -nodes -passin pass:"password"
openssl x509 -in ca.pem -text -noout


openssl pkcs12 -in certificate.p12 -nokeys -out extracted-cert.pem  -passin pass:"password"
openssl x509 -in extracted-cert.pem -noout -text
openssl x509 -in extracted-cert.pem -noout -dates -subject -issuer


openssl pkcs12 -in certificate.p12 -out OUTFILE.crt -nodes
openssl pkcs12 -in certificate.p12 -out OUTFILE.key -nodes -nocerts
openssl pkcs12 -in INFILE.p12 -out OUTFILE.crt -nokeys



openssl pkcs12 -info -in certificate.p12 -nodes -nocerts
```