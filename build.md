# Linux

## Install deps
```
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

## Install Go
```
wget https://go.dev/dl/go1.21.4.linux-amd64.tar.gz && \
tar -C /usr/local -xzf go1.21.4.linux-amd64.tar.gz && \
export PATH=$PATH:/usr/local/go/bin
```

## arm64 ubuntu-latest
```
CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -o linux_arm64 -buildvcs=false
```

## amd64 ubuntu-latest
```
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o linux_amd64 -buildvcs=false
```

# MacOS

## arm64
```
CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o darwin_arm64
```

## intel
```
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o darwin_amd64
```

# Windows
```
# Don't care
```

# Release

## Directories
```
mkdir -p build/{darwin,linux} && \
mkdir -p build/darwin/{arm64,x86_64} && \
mkdir -p build/linux/{arm64,x86_64} && \
mkdir build/tar
```

`Build`

## Copy README
```
for dir in build/{linux,darwin}/{arm64,x86_64}; do
    cp README.md "$dir"/
done
```
## Tar GZip
```
for dir in build/{linux,darwin}/{arm64,x86_64}
do tar -czvf "build/tar/khaossweeper_$(basename $(dirname $dir))_$(basename $dir).tar.gz" -C "$(dirname $dir)" "$(basename $dir)" 
done
```
## ShaSum
```
shasum -s 256 build/tar/* > build/tar/khaossweeper_${VERSION}_checksums.txt
```