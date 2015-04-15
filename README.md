# docker-rsync

docker-rsync recursively watches directories for changes and copies
changes to a docker-machine. Internally [FSEvents API](https://developer.apple.com/library/mac/documentation/Darwin/Reference/FSEvents_Ref/) 
is used for now, so this tool will only work under Mac OSX.

## Installation

```bash
brew install docker-machine
curl 
```

## Usage

```bash
docker-machine create my-machine123 -d virtualbox

cd directory-to-sync-with-rsync
docker-rsync my-machine123
```