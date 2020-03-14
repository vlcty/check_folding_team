# check_folding_team

A check script developed according to the Nagios developer guidelines. Can be used for example in Icinga 2. This script doesn't actually check anything but extracts the following Folding@Home team statistics:

- Team's points aka credits
- Team's completed work units
- Team's active CPU/GPU count in the last 50 days

In addition for each donor the following informations are extracted:

- Donor's points aka credits for that team
- Donor's completed work units for that team

The whole script is written in native go. No external libs needed.

## Install and usage

Clone the repo. Then you can run it from source (if you have go installed):

```
:-$ go run check_folding_team.go -team 12345
```

Otherwise you can compile it for Linux on amd64:

```
:-$ GOOS=linux GOARCH=amd64 go build check_folding_team.go
```

## Please rate limit your requests

The folding at home stats pages are not refreshed on the fly. I think I read somewhere that they are generated every hour. I've defined 30 minutes between each crawl. Please do the same!

## Example implementation into Icinga 2

Example check command:

```
object CheckCommand "folding_team" {
	import "plugin-check-command"

	command = [ PluginDir + "/check_folding_team", "-team", "yourteamid" ]
}
```

Please not that I've hardcoded my team id. If you have multiple teams you must rewrite it to take actual arguments. And the example service:

```
apply Service "folding_team" {
  import "generic-service"

  check_command = "folding_team"
  check_timeout = 120
  check_interval = 30m

  assign where host.name == "somehost"
}
```
