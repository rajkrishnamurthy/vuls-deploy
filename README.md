# Deploy Composite Vuls using Docker-Compose
You don't have to individually deploy the vuls dictionary, db and executables anymore. This docker-compose file helps you to deploy in one-go.

## How is it organized?
The following are the folders and files
### Dockerfiles/Dockerfile:
This contains the composite Dockerfile. It is a multi-stage build file, that finally is built on Alpine 3.11. Seperate them into individual units or add more, if you have to bring in other dictionaries or dbs. 

### etc:
This folder can contain all the env variables that will need to be passed through the docker-compose file into the running instance. For now, there is none

### src/vuls:
This contains all the source code required for the build. These are repo submodules. Refer to the "How to run" section. This folder also contains the start.sh script. This script downloads the db files, if not already present in the shared volume /vuls/data. If you want to fetch these dbs again, you can simply delete the files in /vuls/data and we will automatically fetch these dbs at runtime. When you first start, it can take 30 minutes to download all the db files.

### docker-compose.yml
Runs the composite Dockerfile. The docker-compose file runs the vuls server. The default listen port has been modified to port 80. You will also find that this maps the /vuls/data volume to the localhost on which the docker-compose is run. Because the volume is mapped, the subsequent runs will be instantaneous.


## Pre-Requisites:

You will need to run the following commands:
```
sudo docker network create vuls_default --driver bridge --scope local
sudo docker network create vuls_internal
```
You will see this network used in the docker-compose file. Inorder to establish any connectivity to other containers, say "customapp", you will need to do the following:
```
sudo docker network connect vuls_default customapp
```

## How do we start:

1. Check to make sure that you mount the volume where you want
2. Run the docker network command as stated above or remove the network sections from docker-compose file
3. Clone the repo: 
```
    git clone <repo>
```
4. We will now need to get all the submodules: go-cve-dictionary, go-exploitdb, gost etc. To do this, 
```
    cd <repo>
    git submodule update --remote
```
5. Build the docker-compose file. You may need to use sudo depending on how you installed docker and docker-compose. Docker-Compose, by default, looks for docker-compose.yml. So, you can omit the -f <docker-compose-file.yml> argument, if you not changing the file
```
    sudo docker-compose build -f docker-compose.yml
```
6. Run docker-compose. This will check if the db files are present and if not download and run the injection commands. This may take a LONG time. However, once downloaded, and the volume is set, the subsequent runs will simply skip the datafile injection project. If you do not modify the run command, it will run vuls server on port 80. 
```
    sudo docker-compose up -d
```