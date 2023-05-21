
# Run with docker compose


## Steps to setup

Before continue following the steps, you have to ensure your are in unix platforms,
and [docker](https://docs.docker.com/get-docker/) is already installed.

1.
	make a new directory
	```shell
	mkdir pwp
	cd pwp
	```

2.
	download the `docker-compose.yml` and `init.sql`
	```shell
	wget https://raw.githubusercontent.com/kmcsr/PluginWebPoint/master/docker/docker-compose.yml
	wget https://raw.githubusercontent.com/kmcsr/PluginWebPoint/master/init.sql
	```
	or
	```shell
	curl --output docker-compose.yml https://raw.githubusercontent.com/kmcsr/PluginWebPoint/master/docker/docker-compose.yml
	curl --output init.sql https://raw.githubusercontent.com/kmcsr/PluginWebPoint/master/init.sql
	```

3.
	If you don't want to let your IP to reach the github API limit, you have to register a github application and save them into `gh_secrets.json`
	```json
	{
		"id": "<Your github application id>",
		"secret": "<Your github application secret or leave it empty for manual auth>",
	}
	```

4.
	Login your docker user
	```shell
	docker login
	```

5.
	Start the compose
	```shell
	$ docker compose up -d
	[+] Running 5/5
	 ⠿ Container docker-db-1             Created
	 ⠿ Container docker-api_v1-1         Created
	 ⠿ Container docker-ghupdater-1      Created
	 ⠿ Container docker-web-1            Created
	 ⠿ Container docker-reverse_proxy-1  Created  
	```

	You will get the db-container's name which in this example is `docker-db-1`,
	then you have to init the database by executing the `init.sql`
	```shell
	docker exec -i docker-db-1 mysql -u root <./init.sql
	```
	and restart ghupdater (or actually you can wait it auto refresh for an hour)
	```shell
	docker restart docker-ghupdater-1
	```

6.
	Finally, you can access the website at `http://localhost:2080`
