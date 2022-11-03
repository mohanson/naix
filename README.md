# Naix

Naix is an encrypted port forwarding protocol. Unlike common port forwarding tools, it needs to configure a server and a client, and the communication between the server and the client is encrypted to bypass firewall detection.

```sh
$ go build github.com/mohanson/naix/cmd/naix

# Port forwarding from 20002 to 20000:
$ naix server -s 127.0.0.1:20000 -l :20001 -k password
$ naix client -s 127.0.0.1:20001 -l :20002 -k password
```

Naix is a sub-project of [Daze](https://github.com/mohanson/daze).
