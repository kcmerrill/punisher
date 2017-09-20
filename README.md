[![Build Status](https://travis-ci.org/kcmerrill/punisher.svg?branch=master)](https://travis-ci.org/kcmerrill/punisher) [![Go Report Card](https://goreportcard.com/badge/github.com/kcmerrill/punisher)](https://goreportcard.com/report/github.com/kcmerrill/punisher)

![Punisher](assets/punisher.jpg "Punisher")

Punisher is a simple app that will run a command repeatedly, for as long as necessary. The need for this app came about because I'm always finding myself writing one off scripts to run commands to inject a lot of data. Be it hammering a mysql database to test [gh-ost](https://github.com/github/gh-ost) migrations, testing queues like [crush](https://github.com/kcmerrill/crush) or others like [fetch-proxy](https://github.com/kcmerrill/fetch-proxy), or just to transfer data from one mysql table to another(big int migrations for example)

## Usage

The easiest way to get running with punisher is to run the binary with the command you want to execute.

```bash
#easiest usage
$> punisher <command_to_run>

#usage
$> punisher --workers <int:100> --nice <duration:0s> --duration <duration:forever> <command_to_run>

#example
$> punisher curl https://kcmerrill.com
```

[![asciicast](https://asciinema.org/a/135154.png)](https://asciinema.org/a/138531)

### Parameters
* `workers` are how many workers to run
* `duration` how long should this run for? Default is forever.
* **nice** golang `duration` after a command is run(A simple way to throttle or batch)
* **loop** If set to a non empty string, will create a loop and a counter(more info below)
* **loop-starts-at** Int to start a loop counter at (defaults to zero)
* **loop-ends-at** Int to end the loop counter at (defaults to max int)
* **loop-increment-by** Int to increment the loop counter by
* **retry** If set, will try to retry the command before moving on
* **verbose** Displays the command + success/failure + duration of command

### Templated Commands

When you run a command, you can pass in a template. This means you can alternate data, you can if/else. On top of that, `.Date` and `.UniqID` are both special template params that can be used.

```bash
#example using 'UniqID'
$> punisher curl http://crush.kcmerrill.com/demo/{{ .UniqID }}
```

### The loop

By default, the loop is disabled. To enable a loop counter, simply give the loop a name by providing a `--loop <name>` flag. This can be useful for a number of different things, and you can also nest punisher commands to provide nested loop commands. Give the loop name something uniq, as the string will be replaced in your command with the `LoopIndex`. This is so you can nest your loops/punisher commands if need be. 

You can access the loop counter by using a template `{{ .LoopIndex }}`. Or, if you wish to nest your loops

A basic example
```bash
$> punisher --loop :myspecialcounter  echo :myspecialcounter
# same as:
$> punisher --loop :myspecialcounter echo {{ .LoopIndex }}
```
A nested loop
```bash
$> punisher --loop :outer-loop punisher --loop :inner-loop --loop-ends-at 3 echo :outer-loop::inner-loop
# would echo:
0:0
0:1
0:2
0:3
1:0
1:1
1:2
1:3
# etc ...
```


## Binaries || Installation

[![MacOSX](https://raw.githubusercontent.com/kcmerrill/go-dist/master/assets/apple_logo.png "Mac OSX")](http://go-dist.kcmerrill.com/kcmerrill/punisher/mac/amd64) [![Linux](https://raw.githubusercontent.com/kcmerrill/go-dist/master/assets/linux_logo.png "Linux")](http://go-dist.kcmerrill.com/kcmerrill/punisher/linux/amd64)

via go:

`$ go get -u github.com/kcmerrill/punisher`
