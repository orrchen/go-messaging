#!/bin/bash
# Transform long options to short ones
for arg in "$@"; do
  shift
  case "$arg" in
    "--tag")   set -- "$@" "-t" ;;
    "--repository")   set -- "$@" "-r" ;;
    "--name")   set -- "$@" "-n" ;;
    "--help")   set -- "$@" "-h" ;;
    *)        set -- "$@" "$arg"
  esac
done

OPTIND=1

while getopts "r:t:i:n:h" opt; do
  case $opt in
    h) echo "usage build.sh --tag <tag>  --name <name> --repository <repository>" 
    exit
    ;;
    t) tag="$OPTARG"
    ;;
    r) repository="$OPTARG"
    ;;
    n) name="$OPTARG"
    ;;
    \?) echo "Invalid option -$OPTARG" >&2
    ;;
  esac
done

shift $(expr $OPTIND - 1) # remove options from positional parameters


if [ -z ${tag+x} ]; then
    echo "tag is not set"
    exit;
fi

if [ -z ${name+x} ]; then
    echo "name is not set"
    exit;
fi

if [ -z ${repository+x} ]; then
    echo "repository is not set"
    exit;
fi

GOOS=linux GOARCH=amd64 go build -o app/tcp-server.linux main.go
docker build -t ${name} .
docker tag ${name} ${repository}:${tag}
docker push ${repository}:${tag}


