# Devzat Courier

A plugin to transmit messages between multiple Devzat instances.

![Local test](https://github.com/Arkaeriit/devzat-courier/blob/master/demo.png?raw=true)

## User manual

For users, Devzat Courier is very seamless. Talk in an instance and the message will be transferred to every other instance. There is also a `courier` command to get basic information about the plugin's state.

## Admin manual

To run your own Devzat Courier, you must have a token for every Devzat instance you want to connect. Then, create a JSON configuration file in the following format:

```json
[
  {
    "Host": "devzat.bobignou.red:2222",
    "Token": "dvz@xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=",
    "Prefix": "ABD",
    "PrefixColor": "cyan",
    "NameColor": "yellow"
  },
  {
    "Host": "devzat.hackclub.com:5556",
    "Token": "dvz@xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx=",
    "Prefix": "D@",
    "PrefixColor": "yellow",
    "NameColor": "green"
  }
]
```

There is one object for each Devzat instance which must contain:
* `Host`: URL and gRPC port of the instance.
* `Token`: Devzat token granted in that instance. 
* `Prefix`: Prefix put before every message coming from that instance.
* `PrefixColor`: Color used to display the prefix.
* `nameColor`: Color used to display names of users in that instance.

The available colors are:
* black
* red
* green
* yellow
* blue
* purple
* cyan
* white

Invalid or missing colors will default to uncolored text.

Compile the code in this repository with `go build` and run the resulting executable with the configuration file as argument.

