<div align="center">

# 🐥 Chiko - Beautiful gRPC TUI Client

**The developer-friendly terminal interface for gRPC that makes API testing a joy**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)](LICENSE)
[![Release](https://img.shields.io/github/v/release/felangga/chiko?style=for-the-badge&logo=github)](https://github.com/felangga/chiko/releases)
[![Stars](https://img.shields.io/github/stars/felangga/chiko?style=for-the-badge&logo=github)](https://github.com/felangga/chiko/stargazers)

</div>

---

## 🎯 Why Chiko?

Tired of memorizing complex `grpcurl` commands? Fed up with switching between terminal and documentation? **Chiko** transforms your gRPC testing experience into something beautiful and intuitive.

Built with the power of [grpcurl](https://github.com/fullstorydev/grpcurl) and the elegance of [tview](https://github.com/rivo/tview), Chiko brings you a stunning terminal interface that makes gRPC testing feel like magic ✨

![image](https://github.com/user-attachments/assets/72c74248-8ab3-4c68-a846-8925bfb2fc80)

## 🚀 What Makes Chiko Special

### 🔍 **Smart Server Reflection**
![image](https://github.com/user-attachments/assets/fe63a771-87e5-48d3-9ea8-e85abfe9ed8c)

Discover and browse gRPC endpoints automatically! No more digging through documentation - if your server supports reflection, Chiko shows you everything at a glance.

### 🔐 **Seamless Authorization**
![image](https://github.com/user-attachments/assets/0872e00d-493b-4ca9-ad13-4b46299bf003)

Secure your requests with built-in Bearer token support. Authentication made simple and secure.

### 📋 **Rich Metadata Support**
![image](https://github.com/user-attachments/assets/91987536-52ff-46d0-a3b9-a901a5e17256)

Add custom headers and metadata to your requests with an intuitive interface. No more command-line gymnastics!

### ⚡ **Instant Payload Generation**
![image](https://github.com/user-attachments/assets/b560a034-2419-4a80-920a-4e237b70e61b)

Get perfectly formatted request templates with a single click. Say goodbye to manual JSON crafting!

### 📚 **Smart Bookmarks**
![image](https://github.com/user-attachments/assets/fef777ae-1500-48c6-991f-0cc3b125390a)

Save your favorite requests as bookmarks and replay them instantly. Perfect for API regression testing and development workflows.

## 📦 Installation

Choose your preferred installation method:

### 🍺 Homebrew (Recommended)
```bash
brew install felangga/chiko/chiko
```

### 🐹 Go Install
```bash
go install github.com/felangga/chiko/cmd/chiko@latest
```

### 🔧 From Source
```bash
git clone https://github.com/felangga/chiko
cd chiko
go run ./...
```

### 📥 Pre-built Binaries
Download the latest release from our [Release Page](https://github.com/felangga/chiko/releases) for your platform and architecture.

## 🎮 Quick Start

1. **Launch Chiko**
   ```bash
   chiko
   ```

2. **Connect to your gRPC server**
   - Enter your server URL
   - Set up authentication if needed
   - Browse available services via reflection

3. **Make your first request**
   - Select a method from the sidebar
   - Generate a sample payload
   - Hit send and see the magic! ✨

## 🗺️ Roadmap
### ✅ Completed
- ~~Metadata and headers support~~
- ~~Log dumping to file~~
- ~~Import/export grpcurl commands~~

### 🚧 In Progress
- 📄 **Proto file import** - Support for services without reflection
- 🔒 **Enhanced authentication** - OAuth, API keys, and more
- 🛡️ **Bookmark security** - Password protection for sensitive requests

### 💡 Future Ideas
- 🎨 **Themes and customization**
- 📊 **Request analytics and performance metrics**
- 🔄 **Request history and replay**
- 🌐 **Multi-server workspace management**

---

## 🤝 Contributing

We love contributions! Whether it's:
- 🐛 Bug reports
- 💡 Feature requests  
- 📖 Documentation improvements
- 🔧 Code contributions

Check out our [contributing guidelines](CONTRIBUTING.md) to get started.

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [grpcurl](https://github.com/fullstorydev/grpcurl) - The powerful gRPC command-line tool
- [tview](https://github.com/rivo/tview) - The amazing TUI library for Go
- All our [contributors](https://github.com/felangga/chiko/contributors) who make this project better

---

<div align="center">

**Made with ❤️ by developers, for developers**

[⭐ Star us on GitHub](https://github.com/felangga/chiko) | [🐛 Report Issues](https://github.com/felangga/chiko/issues) | [💬 Discussions](https://github.com/felangga/chiko/discussions)

</div>

