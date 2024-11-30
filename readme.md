relay is a simple command line file transfer over local or global network using tcp protocol.
It don't require any dependencies and is not platform specific. It should work on linux, windows and macOS.

# INSTALLATION
You can install this tool using go command `go install https://github.com/AmirMirzayi/relay` or download binary from [Releases](https://github.com/AmirMirzayi/relay/releases).

# USAGE
> [!IMPORTANT]
> Make sure application has privilege to connect to network and change filesystem.

> [!TIP]
> If you increase the size of read/write buffer. Operating system can touch the I/O device less. And it can read/write larger blocks in each operation. Set buffer size on powers of two with non-negative exponents to achieve more performance and speed.

> [!NOTE]
> Timeout flag receives s, m and h in given timeout format which are Second, Minute And Hour. For example: 72h3m0.5s is 72Hours and 3 Minutes and 0.5 Second.

## Flags
```
-h, --help                   help for relay
-t, --timeout   duration     connection timeout (default 30s)
-i, --ip        ip           sender machine binding ip address (default 0.0.0.0)
-p, --port      int          application running port (default 55555)
-b, --buffer    int          buffer size in byte (default 1048576)
-w, --width     int          progress bar width (default 25)
-l, --silent    bool         silent transfer (default false)
```

## Host serves to send files
send command must have at least 1 argument which has file or directory path
```
relay send [-p 12345 | -i 127.0.0.1 | -b 1024 | -t 120s | -w 25 | -l false] some_file.ext other_file2.ext some_directory_within_subdirectories
```

### Flags
```
-s, --save string   files save path (default "/home/$(USER)/relay")
```

## Connect to the host to receive files
```
relay receive -i 127.0.0.1 [-p 12345 | -b 1024 | -t 120s | -w 25 | -s /home | -l false]
```
