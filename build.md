# Build Docs

## Linux

### Install deps

```bash
apt update && \
apt install -y \
pkg-config \
build-essential \
gcc \
libgl1-mesa-dev \
libx11-dev \
libxcursor-dev \
libxrandr-dev \
libxinerama-dev \
libxinerama-dev \
libxi-dev \
libxrandr-dev \
libxinerama-dev \
libx11-dev \
libxxf86vm-dev \
wget
```

### Install Go

```bash
wget https://go.dev/dl/go1.21.4.linux-amd64.tar.gz && \
tar -C /usr/local -xzf go1.21.4.linux-amd64.tar.gz && \
export PATH=$PATH:/usr/local/go/bin
```

### arm64 ubuntu-latest

```bash
CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -o linux_arm64 -buildvcs=false
```

### amd64 ubuntu-latest

```bash
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o linux_amd64 -buildvcs=false
```

## MacOS

### arm64

```bash
CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o darwin_arm64
```

### intel

```bash
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o darwin_amd64
```

## Windows

```bash
# Don't care
```

## Release

### Directories

```bash
mkdir -p build/{darwin,linux} && \
mkdir -p build/darwin/{arm64,x86_64} && \
mkdir -p build/linux/{arm64,x86_64} && \
mkdir build/tar
```

`Build`

### Copy README

```bash
for dir in build/{linux,darwin}/{arm64,x86_64}; do
    cp README.md "$dir"/
done
```

### Tar GZip

```bash
for dir in build/{linux,darwin}/{arm64,x86_64}
    do tar -czvf "build/tar/khaossweeper_$(basename $(dirname $dir))_$(basename $dir).tar.gz" -C "$(dirname $dir)" "$(basename $dir)" 
done
```

### ShaSum

```bash
shasum -s 256 build/tar/* > build/tar/khaossweeper_${VERSION}_checksums.txt
```
