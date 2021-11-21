package cmd

var HelpStr string = UsageStr + FlagStr + CommandStr

var CommandStr string = `Commands:
  help    list this help message
  listen  run the bot

`

var FlagStr string = `Flags:
  -h --help | show this help message
`

var UsageStr string = `plexer - discord bot to download things to your plex server 
Usage:
  plexer [--help] <command> [<args>]
`
