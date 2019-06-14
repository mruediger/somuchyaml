# Somuchyaml

Somuchyaml is a simple mattermost bot that post the below picture whenever someone mentions yaml.

![the some much yaml commic](yaml.jpg)

## Usage

```
usage: somuchyaml --server=SERVER --websocket=WEBSOCKET --username=USERNAME --password=PASSWORD [<flags>]

A mattermost bot.

Flags:
  --help                   Show context-sensitive help (also try --help-long and --help-man).
  --server=SERVER          server url
  --websocket=WEBSOCKET    the websocket url used for listening
  --username=USERNAME      username to connect to the server
  --password=PASSWORD      password to connect to the server
  --team="darksystem"      the team on the server
  --channel="town-square"  the channel the bot will listen
  --goldenfile             Enable golden file output.
```


Instead of the command line parameters, environment variables can be used:

```
SOMUCHYAML_SERVER
SOMUCHYAML_WEBSOCKET
SOMUCHYAML_USERNAME
SOMUCHYAML_PASSWORD
SOMUCHYAML_TEAM
SOMUCHYAML_CHANNEL
```
