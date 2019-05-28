# Little Load Balancer

The little load balancer is a simple load-balancer (originally designed for `kube-apiserver`) and is designed to quickly instantiate a lb from the cli.

## Example Usage

In this example environment i've deployed a small server (raspberry Pi is an option) that is given the IP address 192.168.1.1 and this is where `llb` will run. 

I've two Kubernetes masters that are ready to be deployed with the addresses (`192.168.1.110` / `.111`), to create the load balancer we run the following command:

```
llb server -e 192.168.0.110:6443,192.168.0.111:6443 -t tcp -l 5 -p 6443
```

This will expose a front end on port 6443 (`-p 6443`) to the endpoints defined (`-e <address:port>`) using commas to define multiple endpoints. Logging is turned up to debug with `-l 5`

```
               -----------------
               |192.168.1.1:6443|
               -----------------
                   |        |
   -------------------    -------------------
   |192.168.1.110:6443|   |192.168.1.111:6443|
   -------------------    -------------------
             
```

## Getting LLB

`go get github.com/thebsdbox/llb`
