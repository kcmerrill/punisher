[![Build Status](https://travis-ci.org/kcmerrill/punisher.svg?branch=master)](https://travis-ci.org/kcmerrill/punisher) [![Go Report Card](https://goreportcard.com/badge/github.com/kcmerrill/punisher)](https://goreportcard.com/report/github.com/kcmerrill/punisher) ![Punisher](assets/punisher.jpg "Punisher") 
Punisher is a simple app that will run a command repeatedly. The need for this app came about because I'm always finding myself writing one off scripts to run commands to inject a lot of data. Be it hammering a mysql database to test [gh-ost](https://github.com/github/gh-ost) migrations, testing queues like [genie](https://github.com/kcmerrill/genie) or others like [fetch-proxy](https://github.com/kcmerrill/fetch-proxy).

## Usage

```bash
$> punisher --workers <int:100> --nice <duration:0s> --duration <duration:forever> <command to run>

#example
$> punisher curl https://kcmerrill.com
```

### Parameters
* **workers** are how many `threads` to spin up
* **duration** how long should this run for? Default is forever.
* **nice** golang `duration` after a command is run(A simple way to throttle or batch)

