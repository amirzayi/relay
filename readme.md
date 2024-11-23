relay is a simple command line file transfer over local or global network using tcp protocol.
It don't require any dependencies and is not platform specific. It should work on linux, windows and macOS.

# INSTALLATION
You can install this tool using go command `go install https://github.com/AmirMirzayi/relay"` or download binary from [Releases][https://github.com/AmirMirzayi/relay/releases].

# USAGE
> [!IMPORTANT]
> Make sure application has privilege to connect to network and change filesystem.

## Flags
```
-h, --help              help for relay
-i, --ip        ip      sender machine binding ip address (default 0.0.0.0)
-p, --port      int     application running port(default 55555)
-t, --timeout   int     connection timeout in second (default 30)
-w, --width     int     progress bar width (default 25)
```

## Host serves to send files
send command must have at least 1 argument which has file or directory path
```
relay send [-p 12345 | -i 127.0.0.1 | -t 120 | -w 25] some_file.ext other_file2.ext some_directory_within_subdirectories
```

### Flags
```
-s, --save string   files save path (default "/home/$(USER)/relay")
```

## Connect to host to receive files
```
relay receive -i 127.0.0.1 [-p 12345 | -t 120 | -w 25 | -s /home]
```
