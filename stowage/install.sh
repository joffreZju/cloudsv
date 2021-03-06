#!/bin/bash

path=$PWD
oldGOPATH=$GOPATH
export GOPATH=$PWD

GoInstall() {
	#go fmt src/$1
    	#go tool vet src/$1
    	#golint src/$1
	timestamp=`date +%Y-%m-%dT%H:%M:%S` 
	appname=`echo $1|awk -F '/' '{print $NF}'`
	appname=${appname}".exe"
	gitversion=`git log -n 1|grep commit | awk -F ' ' '{print $2}'`
	go build -ldflags "-X main.Version=${timestamp}_${gitversion} -X common/lib/deamon.Version=${timestamp}_${gitversion}" -o "$PWD/bin/$appname" $1
	#go install -ldflags "-X main.Version=${timestamp}_${gitversion} -X common/lib/deamon.Version=${timestamp}_${gitversion}" $1
	#go build -o "$PWD/bin/$appname" $1

	if [ $? -eq 0 ]; then 
		echo -e "\033[32m ====== go install $1 success ==== \033[0m"
	else
		echo $ret
		echo -e "\033[31m ====== go install $1 error ====== \033[0m"
	fi
	
}

((pos+=1)) ; bins[$pos]="stowage_web"

for ((i=1; i<=${#bins[*]}; i++)) 
do
	bin=${bins[$i]}
	if [ "$1" = "" ]; then
		GoInstall $bin
	elif [[ $bin == *$1* ]]; then 
		GoInstall $bin
	fi
done

export GOPATH=$oldGOPATH

bin/stowage_web.exe -c conf/stowage_web.conf