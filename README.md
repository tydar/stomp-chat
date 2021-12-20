# STOMP Chat

Chat client to operate through a STOMP pub/sub server. Designed to demonstrate my project [stomper](https://github.com/tydar/stomper).

Might do for fun:

1) Support multiple channels
2) Support DMs
	* Have a pub/sub channel per DM that is created with a dynamic UUID to avoid snooping ?

Not really planning to write a dedicated chat server, so options are limited feature-wise.

Minor things I may clean up:

1) inconsistency in message spacing
2) additional UI controls

Borrows `frames.go` from stomper. Also looked at [termoose/irccloud](https://github.com/termoose/irccloud) to help me grok writing a chat client & using [tview](https://github.comrivo/tview) while running multiple goroutines.

## To run

1. Pull the [stomper](https://github.com/tydar/stomper) image from GHCR and clone this repo:

```bash
$ docker pull ghcr.io/tydar/stomper:main
$ git clone --depth=1 https://github.com/tydar/stomper.git
```

2. Build your chat client:

```bash
$ go build -o stomp-chat .
```

3. Run a stomper server with the appropriate environment variables for testing.

```bash
docker run -p 32801:32801 \
-e STOMPER_HOSTNAME=0.0.0.0 \
-e STOMPER_TOPICS=/channel/main \
-e STOMPER_TCPDEADLINE=0 \
ghcr.io/tydar/stomper:main
```

4. Run stomp-chat:

```bash
$ ./stomp-chat -uname=a_chat_user
```

A full listing of command line flags is available by running this command:

```bash
$ ./stomp-chat -h
-host string
	hostname of stomp server (default "localhost")
-log string
	filename for logging; destructively created each run (default "chat.log")
-port int
	port number for stomp server (default 32801)
-uname string
	username (default "default_guy_123")
```




