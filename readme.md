# Go watchdog example with net/rpc, called async

This is a example on how to implement a watchdog timer using golang with RPC as a heartbeat. The server will run a watchdog timer, set my thisapi.Init . You can change the watchdog timeout by setting the value in millisecond but the default is 5000 milliseconds (5 seconds). The server will then run a gorutine that will count down the timer until it reaches 0 or below and exit. You then use the client to send heartbeat via a RPC every 2 seconds or whatever value you set the sleep in the main loop.

## Usage

Run the server first
` $ go run ./server/. `

then run the client
` $ go run ./client/. `

<video width="640" height="480" controls>
  <source src="./useage_example.mp4" type="video/mp4">
</video>

## edit config
the config data is a yaml file located in the data directory. It is set to local address "127.0.0.1" with port 1234 and 1235 for the client and server respectively. You can change the server or client config in the yaml file.

```
client:
  ipv4_address: 127.0.0.1
  port: 1234
server:
  ipv4_address: 127.0.0.1
  port: 1235
```