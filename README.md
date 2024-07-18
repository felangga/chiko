# Chiko - TUI gRPC client

Chiko is a TUI (Terminal User Interface) gRPC client. It is a simple tool to interact with gRPC services using a beautiful terminal interface. This project using [grpcurl](https://github.com/fullstorydev/grpcurl) library to interact with gRPC services, and combine with the beautiful terminal interface using [rivo's tview](https://github.com/rivo/tview) library. I love using `grpcurl` to interact with the gRPC services, but I'm bad at remembering the flags and the syntax. So, I created this tool to help me interact with the gRPC services easily.

https://github.com/user-attachments/assets/acc4be04-6be6-4743-ad30-ddfe1bcc229d

## Install 
You can visit the [Release Page](https://github.com/felangga/chiko/releases), and select the version that you want to download, the architecture and the operating system that you are using.

### Homebrew
Currently `chiko` is not available on the `homebrew-core`, so you can install directly from our repository by using:
```
brew install felangga/chiko/chiko
```

### Go
```
go install github.com/felangga/chiko@latest
```

### Manual
```
git clone https://github.com/felangga/chiko
cd chiko
go run ./...
```

## Features
- List all gRPC methods using server reflections
- Generate sample request payload 
- Bookmark support

## To Do
- ~Add support for metadata (headers)~
- Add support to import `.proto` files for services that don't support server reflection
- Add support for any authorization types
- Add dump log to file feature
- Add an option to protect or lock the bookmark library
- Add support to import and export from and to `grpcurl` command
