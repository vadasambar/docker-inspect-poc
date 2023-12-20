### What
PoC to get docker image size before pulling the image. It connects with a docker registry (doesn't have to be dockerhub) directly, gets the image metadata (equivalent of `docker manifest inspect <image> -v`) and prints the total image size (compressed). It also prints the execution time of the program. 

### Why
Many times you want to know the size of the image before pulling it so that you can warn your users by emitting a pod event if the image size is too big. The main intention for this PoC is to evaluate adding such a feature to [warmmetal/csi-driver-image](https://github.com/warm-metal/csi-driver-image/tree/master)

### How to use
```
$ go run main.go ubuntu:latest
docker.io/library/ubuntu:latest is 28MB in size
======================
that took 12.435566417 seconds

$ go run main.go warmmetal/csi-image:v0.6.3 
docker.io/warmmetal/csi-image:v0.6.3 is 40MB in size
======================
that took 2.977110167 seconds
```

I am using `1.21.1` version of Golang.