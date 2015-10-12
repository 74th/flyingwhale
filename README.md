# Flying Whale

This makes a container as a command line tool from many package management systems.

## what's this?

If you want to use a command (ex: "marked" markdown compiler) without installing nodejs, you type this.

```
whale npm install marked
```

You can use it!

```
marked README.md
```

## how it work

It creates a docker container installed the command by using package managers, and adds a script calling the container into /usr/local/bin/.

```
#!/bin/sh
# This script was created by flying docker 0.1
docker run -it --rm -v `pwd`:/src --workdir=/src --entrypoint=marked whale-npm-marked $*
```

## install

Available for MacOS and Linux. https://github.com/74th/flyingwhale/releases

## supporting package manages

* whale npm install <package-name>
* whale yum install <package-name>
* whale apt-get install <package-name>

## features

* windows
* behind proxy
