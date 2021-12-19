# STOMP Chat

Chat client to operate through a STOMP pub/sub server. Designed to demonstrate my project [stomper](https://github.com/tydar/stomper).

Todo to finish:

1) Allow runtime configuration of hostname and port
2) Allow runtime configuration of username

Might do for fun:

1) Support multiple channels
2) Support DMs
	* Have a pub/sub channel per DM that is created with a dynamic UUID to avoid snooping ?

Not really planning to write a dedicated chat server, so options are limited feature-wise.

Borrows `frames.go` from stomper. Also looked at [termoose/irccloud](https://github.com/termoose/irccloud) to help me grok writing a chat client & using [tview](https://github.comrivo/tview) while running multiple goroutines.
