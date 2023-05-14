
# Docker Image for PluginWebPoint

This is the introduction about the docker images for this project

I recommend you to use nginx as the reverse proxy to integrate the APIs and the web view

## Database

You should create a docker bridge network to connect the containers to the database.  
You can only use mysql database for now (you can use `mysql:latest` image),
and you need to run [init.sql](../init.sql) after the database is created.

## Examples

#### nginx site file:

```nginx
server {
	listen       443 ssl;
	server_name  mcdreforged.example.com;
	charset UTF-8;

	ssl_certificate      /home/example/.acme.sh/mcdreforged.example.com/fullchain.cer;
	ssl_certificate_key  /home/example/.acme.sh/mcdreforged.example.com/mcdreforged.example.com.key;

	location /v1/ {
		rewrite ^/v1/(.*)$ /$1 break;
		proxy_redirect   off;
		proxy_set_header X-Forwarded-Proto $scheme;
		proxy_set_header X-Forwarded-Host $host:$server_port;
		proxy_set_header X-Forwarded-Server $host;
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
		proxy_set_header X-Forwarded-Port $server_port;
		proxy_set_header Host $host;
		proxy_set_header X-Real-IP $remote_addr;
		proxy_pass       http://127.0.0.1:2082;
	}

	location / {
		include uwsgi_params;
		proxy_set_header X-Forwarded-Proto $scheme;
		proxy_set_header X-Forwarded-Host $host:$server_port;
		proxy_set_header X-Forwarded-Server $host;
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
		proxy_set_header X-Forwarded-Port $server_port;
		proxy_set_header Host $host;
		proxy_set_header X-Real-IP $remote_addr;
		proxy_pass       http://127.0.0.1:2080;
	}
}
```

#### Env file: _(settings.env)_
```env
GH_CLI_ID=***
GH_CLI_SEC=***

DB_USER=pwp_user
DB_PASSWD=the_password_or_leave_it_empty
DB_ADDR=tcp(ip.to.the.database.contianer:3306)
DB_NAME=pluginDatabase
```

#### Script to start github updater:
```bash
#!/usr/bin/bash

docker create --name pwp_ghupdater \
 --restart=on-failure \
 --env-file ./settings.env \
 --network=pwebpoint \
 craftmine/pwp:ghupdater

docker start pwp_ghupdater

```

#### Script to start web hoster:
```bash
#!/usr/bin/bash

docker create --name pwp_web \
 --restart=always \
 --env-file ./settings.env \
 --network=pwebpoint \
 -p 2080:80 \
 craftmine/pwp:web

docker start pwp_web

```

#### Script to start APIs:
```bash
#!/usr/bin/bash

docker create --name pwp_v1 \
 --restart=always \
 --env-file ./settings.env \
 --network=pwebpoint \
 -p 2082:80 \
 craftmine/pwp:v1

docker start pwp_v1

```
