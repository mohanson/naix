# Dahlia

Dahlia is an encrypted port forwarding protocol.

# Tutorials

Forward data from port 20002 to 20000:

```sh
$ dahlia -s 127.0.0.1:20000 -l :20001 -k password -c xo
$ dahlia -s 127.0.0.1:20001 -l :20002 -k password -c ox
```
