# Robolucha publisher

Serves match state thru websocket to the users

# channels used 

- luchador.ID.message, console messages for specific luchador
- match.ID.state, match state updates
- match.ID.event, match events, e.g. luchador kill, match start, match end.


## Local environment setup

Create symbolic link from workspace to gopath
```
	ln -s /home/hamilton/Code/robolucha/robolucha-publisher /home/hamilton/go/src/gitlab.com/robolucha
```
