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

