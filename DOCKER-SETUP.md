# 🐳 Docker/Podman Quick Reference

## 🚀 **One-Command Setup**

```bash
# Auto-detect best runtime and start
./run.sh run

# Or use Makefile
make run
```

## 🔧 **Runtime Selection**

### **Auto-Detection (Recommended)**
```bash
./run.sh build          # Auto-detect runtime
./run.sh run            # Auto-detect runtime
make build              # Auto-detect runtime
make run                # Auto-detect runtime
```

### **Force Specific Runtime**
```bash
./run.sh build docker   # Force Docker
./run.sh run podman     # Force Podman
make docker build       # Force Docker
make podman build       # Force Podman
```

## 📋 **Essential Commands**

```bash
# Build & Run
./run.sh build          # Build image
./run.sh run            # Start services
./run.sh stop           # Stop services

# Development
./run.sh run-dev        # Start with hot reloading
./run.sh logs-dev       # View dev logs

# Management
./run.sh status         # Show container status
./run.sh logs           # View logs
./run.sh clean          # Clean everything

# AI Models
make pull-ollama        # Download AI models
```

## 🎯 **Runtime Priority**

1. **Docker** (if available) - Better tooling
2. **Podman** (if Docker unavailable) - Modern alternative
3. **Manual override** - Force specific choice

## 🔍 **Check What's Available**

```bash
./run.sh help           # Show all options
make runtime-info       # Show detected runtime
make help               # Show Makefile commands
```

## 🚨 **Troubleshooting**

```bash
# Reset everything
./run.sh clean

# Check runtime
./run.sh help

# Force specific runtime
./run.sh build podman   # If Docker has issues
```

## 💡 **Pro Tips**

- **First time**: Use `./run.sh run` (auto-detection)
- **Development**: Use `./run.sh run-dev` (hot reloading)
- **Issues**: Force Podman with `./run.sh run podman`
- **Clean slate**: Use `./run.sh clean` then rebuild

## 🌐 **Access Points**

- **Web App**: http://localhost:8080
- **Ollama API**: http://localhost:11434

---

**🎯 The system automatically chooses the best runtime for you!**
