module github.com/smantic/plexer

go 1.17

require (
	github.com/bwmarrin/discordgo v0.23.3-0.20211228023845-29269347e820
	golift.io/starr v0.11.12
)

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/rs/zerolog v1.26.1 // indirect
	github.com/smantic/libs/discord/imux v0.0.0-20220102185838-af74dc627f98 // indirect
	golang.org/x/crypto v0.0.0-20211215165025-cf75a172585e // indirect
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
)

replace golift.io/starr => github.com/smantic/starr v0.11.13-0.20211214061748-df7349f089fc

replace github.com/smantic/libs/discord/imux => ../libs/discord/imux
