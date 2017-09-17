## Serving with Caddy

```
up.example.com {
	proxy / 127.0.0.1:8002 {
		transparent
	}
	gzip
}
```

## Serving with Apache

```xml
<VirtualHost *:80>
    ServerName up.example.com
    DocumentRoot /data/goploader-server

    <Directory /data/goploader-server>
      Require all granted
    </Directory>

    ProxyPass / "http://localhost:8080/"
    ProxyPassReverse / "http://localhost:8080/"

    LogLevel warn
    ErrorLog /data/logs/apache2/goploader.log
    CustomLog /var/log/apache2/goploader.log combined
</VirtualHost>
```