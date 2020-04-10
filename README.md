# Deploy Composite Vuls using Docker-Compose
You don't have to individually deploy the vuls dictionary, db and executables anymore. This docker-compose file helps you to deploy in one-go.

## How is it organized?
The following are the folders and files
### Dockerfiles/Dockerfile:
This contains the composite Dockerfile. It is a multi-stage build file, that finally is built on Alpine 3.11. Seperate them into individual units or add more, if you have to bring in other dictionaries or dbs. 

### etc:
This folder can contain all the env variables that will need to be passed through the docker-compose file into the running instance. For now, there is none

### src/vuls:
This contains all the source code required for the build. These are repo submodules. Refer to the "How to run" section. This folder also contains the start.sh script. This script downloads the db files, if not already present in the shared volume /vuls/data. If you want to fetch these dbs again, you can simply delete the files in /vuls/data and we will automatically fetch these dbs at runtime. When you first start, it can take 15 minutes to download all the db files.

### docker-compose.yml
Runs the composite Dockerfile. The docker-compose file runs the vuls server. The default listen port has been modified to port 80. 

You will need to run the following commands:

```
sudo docker network create vuls_default --driver bridge --scope local
sudo docker network create vuls_internal
```
