## Running with Supervisor

I enjoy supervisor for its ease of use and the control it offers on the process.
You can run the goploader server using a configuration file that could look like
that :

```ini
[program:goploader]
directory=/home/youruser/goploader/
command=/home/youruser/goploader/server
user=youruser
autostart=true
autorestart=true
stdout_logfile=/var/log/supervisor/goploader_stdout.log
stderr_logfile=/var/log/supervisor/goploader_stderr.log
```

Of course you need to replace `youruser` with the user you installed goploader 
with.

## Running with Systemd

You can also use a Systemd Unit to launch Goploader on boot.

```ini
[Unit]
Description=goploader

[Service]
Type=simple
User=youruser
WorkingDirectory=/home/youruser/goploader
ExecStart=/home/youruser/goploader/server


[Install]
WantedBy=multi-user.target
```

As for supervisor, you'll need to replace `youruser` with the user you installed
goploader with.