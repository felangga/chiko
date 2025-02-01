# Chiko - TUI gRPC client

Chiko is a TUI (Terminal User Interface) gRPC client. It is a simple tool that interacts with gRPC services using a beautiful terminal interface. This project uses [grpcurl](https://github.com/fullstorydev/grpcurl) library to interact with gRPC services, and combine with the beautiful terminal interface using [rivo's tview](https://github.com/rivo/tview) library. I love using `grpcurl` to interact with the gRPC services, but I'm bad at remembering the flags and the syntax. So, I created this tool to help me interact with the gRPC services easily.

![image](https://github.com/user-attachments/assets/72c74248-8ab3-4c68-a846-8925bfb2fc80)

## Features
### Server Reflection Support 
![image](https://github.com/user-attachments/assets/fe63a771-87e5-48d3-9ea8-e85abfe9ed8c)
*You can browse the server's endpoints if the server supports server reflection, currently manual proto import is not yet supported*
### Authorization
![image](https://github.com/user-attachments/assets/0872e00d-493b-4ca9-ad13-4b46299bf003)
*Currently only Bearer authorization is supported*
### Metadata
![image](https://github.com/user-attachments/assets/91987536-52ff-46d0-a3b9-a901a5e17256)
*You can append metadata into the request header*
### Generate Request Payload
![image](https://github.com/user-attachments/assets/b560a034-2419-4a80-920a-4e237b70e61b)
*You can easily get the request template format by clicking the `Generate Sample` button*
### Bookmark Support
![image](https://github.com/user-attachments/assets/fef777ae-1500-48c6-991f-0cc3b125390a)
*You can save your payload request into a bookmark, so you can easily invoke them in the future*

## Install 
You can visit the [Release Page](https://github.com/felangga/chiko/releases), and select the version that you want to download, the architecture, and the operating system that you are using.

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

## To Do
- ~Add support for metadata (headers)~
- ~Add dump log to file feature~
- Add support to import `.proto` files for services that don't support server reflection
- Add support for any authorization types
- Add an option to protect or lock the bookmark library
- Add support to import and export from and to `grpcurl` command
