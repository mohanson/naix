# Dahlia

Dahlia is an encrypted port forwarding protocol. Unlike common port forwarding tools, it needs to configure a server and a client, and the communication between the server and the client is encrypted to bypass firewall detection.

# Tutorials

Forward data from port 20002 to 20000:

```sh
$ dahlia server -s 127.0.0.1:20000 -l :20001 -k password
$ dahlia client -s 127.0.0.1:20001 -l :20002 -k password
```
