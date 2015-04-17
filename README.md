# docker-rsync

docker-rsync recursively watches directories for changes and copies
changes to a docker-machine. It is a drop in replacement for the 
existing boot2docker vboxsf feature.

Please note though that syncing happens only in one direction, 
from your local machine to the boot2docker machine. If you want to sync
from a Docker container back to your local machine, docker-rsync is not
the tool you're looking for. 

__Is it fast?__ Yes! While the initial sync might take some seconds
(depending on the number of files you want to sync), following syncs are 
super fast (compared to vboxsf & NFS). A one file sync usually takes less than 100ms.

docker-sync relies on [FSEvents API](https://developer.apple.com/library/mac/documentation/Darwin/Reference/FSEvents_Ref/), 
so this tool will only work under Mac OSX. You also need [docker-machine](https://github.com/docker/machine) installed.


## Installation

```bash
brew install docker-machine

brew tap synack/docker
brew install docker-rsync
```

## Usage

```bash
docker-machine create my-machine123 -d virtualbox

cd sync-this-directory
docker-rsync my-machine123

echo "git" >> .rsyncignore
```
