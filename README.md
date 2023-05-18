# corporate-headquarters

GO GO GOLEMS CORPORATE HEADQUARTERS - THIS IS WHERE WE WORK

## Building docker images with goreleaser 

The goreleaser target rule (`make goreleaser`) will build cross-platform
docker images for amd64 and arm64. While docker desktop comes with qemu binary 
emulation support, when installing docker through other sources (for example on Linux),
you need to first install the emulations:

```sh 
docker run --privileged --rm tonistiigi/binfmt --install arm64,riscv64,arm
```