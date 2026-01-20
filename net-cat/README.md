# TCPChat

TCP chat server written in Go with username, multiple rooms, and logging.

### Usage
go build -o TCPChat main.go <br>
./TCPChat [port]

The default port is 8989

### Client-side commands :<br>

/help — displays available commands<br>
/quit — leave the server<br>
/newname newusername — change username<br>
/join roomname — switch room<br>
/list — list available rooms<br>

### Logs

Logs are stored in the logfiles directory, under this format:
log_YYYY-MM-DD_HH-MM-SS.log

<br>
<img src="./img/team_picture.png" alt="team picture">