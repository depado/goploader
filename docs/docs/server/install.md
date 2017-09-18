## Web Setup

There are currently two ways of configuring a goploader server. You can either 
run the binary a first time and follow the setup instructions in your browser 
(it will first run on the 8008 port, so make sure this port is available).

```shell
$ ./server --initial
Please go to http://127.0.0.1:8008 to setup goploader.
```

You'll then need to fill the form that is shown to you. 
This step will generate a `conf.yml` file next to your binary file. This is the 
configuration file of your service. You can freely modify it, keep in mind that 
you'll need to restart the service each time you modify it so the changes are 
applied. 

## The conf.yml file

The other solution is to directly create a `conf.yml` file next to the binary 
(or wherever you want, you can actually specify a configuration file with the 
`-c/--conf flag`) 

The next sections describe how to write and custmomize this file.

### Server


| Key         | Description                                                                             | Default | Example   |
|-------------|-----------------------------------------------------------------------------------------|---------|-----------|
| name_server | The server name to use. This is used to return correct URLs to the clients.             |         | gpldr.in  |
| port        | Port on which the server should listen.                                                 | 8080    | 8003      |
| host        | Host on which the server should listen.                                                 |         | 127.0.0.1 |
| append_port | Set to `true` if you need to append the port to the `name_server` in the returned URLs. | false   | false     |

### HTTPS

| Key             | Description                                                              | Default | Example     |
|-----------------|--------------------------------------------------------------------------|---------|-------------|
| serve_https     | Set to `true` if you want the server to serve its own HTTPS certificate. | false   | false       |
| ssl_cert        | Path to the SSL Certificate. Required if you set `serve_https` to true.  | -       | ssl/my.cert |
| ssl_private_key | Path to your SSL private key. Required if you set `serve_https` to true. | -        | ssl/my.pem  |

### Files and Database

| Key | Description | Default | Example |
|---------------|----------------------------------------------------------------------------------------------------------------------------------------------|--------------|--------------|
| upload_dir | Path to the (local or absolute) directory where the files should be stored. | up/ | up/ |
| db | Path to the database. (The database being a single file) | goploader.db | goploader.db |
| uniuri_length | Length of the unique ID of the resources, which is part of the returned URL for each uploaded file. | 10 | 10 |
| key_length | Size of the AES cipher key.  | 16 | 16 |
| size_limit | Maximum size of sent files in MB. | 20 | 20 |
| disk_quota | Maximum total size the uploaded files occupy on disk in GB. 0 to disable. | 0 | 3 |
| loglevel | Defines a level for the logs. Choices are "debug", "info", "error","disabled". The names are pretty explicit. "info" is the most useful one. | "info" | "info" |

### Privacy

| Key | Description | Default | Example |
|----------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------|---------|
| sensitive_mode | If web interface is activated, hide all information and only display the file upload form. | false | true |
| stats | Display or not the statistics on the main page (number of files uploaded and total size). Note that, even when setting this to false, statistics will still be gathered. | false | true |
| no_web | When set to true, this option completely disable the HTML pages. Upload will only be possible using curl and/or client. | false | false |
| fulldoc | Display or not the link to the full documentation of the project (which you are reading right now) | true | true |

### Others

| Key | Description | Default | Example |
|--------------------|---------------------------------------------------------------------------------------------------------------|---------|---------|
| always_download | Disable viewing the file in the browser and always offer to download the files | false | true |
| disable_encryption | Disable encryption. Files will be written to disk much faster, but there won't be any server-side encryption. | false | false |
| prometheus_enabled | Enabled Prometheus metrics. This exposes the /metrics endpoint which will expose data for prometheus | false | true |

## Sensitive Mode

This is actually a really important setting. If you know the server will be used
to host sensitive information, please set this to `true`. This will disable any
information if the web interface is activated. For example, there won't be any
footer, the project's name won't be displayed, only the file upload form will. 
This is a simple way of not affiliating me to your project/company/organization.
This setting will obviously override the fulldoc and stats settings as there 
won't be anything displayed on the page.