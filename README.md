# Robolucha publisher

Serves match state thru websocket to the users

# channels used 

- luchador.ID.message, console messages for specific luchador
- match.ID.state, match state updates
- match.ID.event, match events, e.g. luchador kill, match start, match end.


## Local environment setup

Create symbolic link from workspace to gopath
```
	export WIN=/mnt/c/Users/hamil/code
	mkdir -p $WIN/go/src/gitlab.com/robolucha
	ln -s $WIN/robolucha-api $WIN/go/src/gitlab.com/robolucha
	export PATH=$PATH:$WIN/go/bin
	cd $WIN/go/src/gitlab.com/robolucha/robolucha-publisher
	go get -v	
```

when using Ubuntu on Windows, create the symlink on Windows to be regonized on both systems
``
mklink /d C:\Users\hamil\code\go\src\gitlab.com\robolucha\robolucha-publisher C:\Users\hamil\code\robolucha-publisher 
```
