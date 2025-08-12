# 🐳 Docker/Podman Setup for Novel Q&A Assistant

This guide explains how to dockerize and run the Novel Q&A Assistant application using **Docker** or **Podman**.

## 📋 Prerequisites

- **Docker** OR **Podman** installed on your system
- **Docker Compose** OR **Podman Compose** installed
- At least 4GB of RAM available for Ollama models

## 🚀 Quick Start

### 1. Production Setup

```bash
# Build and run the production application
make build
make run

# Or use compose directly
podman-compose up -d
# or
docker-compose up -d
```

### 2. Development Setup (with hot reloading)

```bash
# Build and run the development application
make build-dev
make run-dev

# Or use compose directly
podman-compose -f docker-compose.dev.yml up -d
# or
docker-compose -f docker-compose.dev.yml up -d
```

## 🏗️ Docker/Podman Architecture

The application uses a multi-service architecture:

- **novel-qa**: Your Go application (port 8080)
- **ollama**: AI model service (port 11434)
- **Persistent volumes**: For novels, ChromaDB data, and Ollama models

## 📁 File Structure

```
.
├── Dockerfile              # Production Docker image
├── Dockerfile.dev         # Development Docker image
├── docker-compose.yml     # Production services
├── docker-compose.dev.yml # Development services
├── .air.toml             # Hot reloading config
├── .dockerignore         # Docker build exclusions
├── Makefile              # Common commands (Auto-detection)
└── run.sh                # Runtime selector script
```

## 🔧 Available Commands

### 🎯 **Auto-Detection (Recommended)**

The Makefile automatically detects and uses the best available runtime:

```bash
# Show which runtime is being used
make runtime-info

# Auto-detect and use best runtime
make build          # Build production image
make run            # Start production services
make run-dev        # Start development services
```

### 🐳 **Using the Runtime Selector Script**

The `run.sh` script provides an easy way to choose your runtime:

```bash
# Show help
./run.sh help

# Auto-detect runtime
./run.sh build
./run.sh run

# Force specific runtime
./run.sh build docker     # Force Docker
./run.sh run podman       # Force Podman
./run.sh build auto       # Auto-detect (default)
```

### 🔧 **Using Makefile Directly**

```bash
# Show all available commands
make help

# Production commands
make build          # Build production image
make run            # Start production services
make stop-prod      # Stop production services
make logs           # View production logs

# Development commands
make build-dev      # Build development image
make run-dev        # Start development services
make stop-dev       # Stop development services
make logs-dev       # View development logs

# Utility commands
make shell          # Open shell in production container
make shell-dev      # Open shell in development container
make status         # Show container status
make clean          # Clean up everything

# Force specific runtime
make docker build   # Force Docker
make podman build   # Force Podman
```

### 🐳 **Using Compose Directly**

```bash
# With Podman
podman-compose up -d
podman-compose down
podman-compose logs -f

# With Docker
docker-compose up -d
docker-compose down
docker-compose logs -f
```

## 🎯 Development vs Production

### Development Mode
- **Hot reloading** with Air
- **Source code mounted** for live editing
- **Debug mode** enabled
- **Faster iteration** cycle

### Production Mode
- **Optimized binary** (smaller image)
- **Security** with non-root user
- **Health checks** enabled
- **Production-ready** configuration

## 📊 Ollama Models

After starting the services, pull the AI models:

```bash
# Pull all supported models
make pull-ollama

# Or pull individually
podman-compose exec ollama ollama pull phi3
podman-compose exec ollama ollama pull llama3
podman-compose exec ollama ollama pull mistral
podman-compose exec ollama ollama pull gemma
```

## 🔍 Monitoring & Debugging

### View Logs
```bash
# All services
make logs

# Specific service
podman-compose logs -f novel-qa
podman-compose logs -f ollama
```

### Container Shell Access
```bash
# Production container
make shell

# Development container
make shell-dev
```

### Health Check
The application includes a health check endpoint:
```bash
curl http://localhost:8080/
```

## 🗄️ Data Persistence

The following data is persisted across container restarts:

- **Novels**: `/app/novels` → `novel_data` volume
- **ChromaDB**: `/app/chroma_db` → `chroma_data` volume
- **Ollama Models**: `/root/.ollama` → `ollama_data` volume

## 🚨 Troubleshooting

### Podman-Specific Issues

1. **Permission denied errors**
   ```bash
   # Podman runs rootless by default
   # If you get permission errors, try:
   podman unshare chown -R 1001:1001 /path/to/volume
   ```

2. **Port binding issues**
   ```bash
   # Podman might need explicit port binding
   # Check if ports are available:
   netstat -tulpn | grep :8080
   ```

3. **Volume mounting issues**
   ```bash
   # Podman volumes might need different paths
   # Use absolute paths in docker-compose.yml
   ```

### Common Issues

1. **Port already in use**
   ```bash
   # Check what's using port 8080
   lsof -i :8080
   # Stop conflicting service or change port in docker-compose.yml
   ```

2. **Out of memory**
   ```bash
   # Ollama needs at least 4GB RAM
   # Check system memory allocation
   ```

3. **Permission denied**
   ```bash
   # Clean up and rebuild
   make clean
   make build
   make run
   ```

### Reset Everything
```bash
# Stop and remove everything
make clean

# Rebuild from scratch
make build
make run
```

## 🔒 Security Considerations

- **Non-root user** in production container
- **Read-only** file system where possible
- **Network isolation** between services
- **No secrets** in Docker images
- **Rootless containers** with Podman (default)

## 📈 Scaling

For production deployment:

1. **Use external database** instead of file-based ChromaDB
2. **Add load balancer** for multiple instances
3. **Use Docker Swarm, Kubernetes, or Podman Pods** for orchestration
4. **Implement proper logging** (ELK stack)
5. **Add monitoring** (Prometheus + Grafana)

## 🌐 Environment Variables

You can customize the application behavior:

```bash
# In docker-compose.yml
environment:
  - OLLAMA_HOST=http://ollama:11434
  - GIN_MODE=debug  # Only in development
```

## 🐳 Podman vs Docker

### Podman Advantages
- **Rootless** by default (more secure)
- **No daemon** required
- **Compatible** with Docker commands
- **Better** for CI/CD environments

### Docker Advantages
- **More mature** ecosystem
- **Better** desktop integration
- **More** third-party tools

### Compatibility
- **Same commands** work with both
- **Same compose files** work with both
- **Easy switching** between them

## 🎯 **Runtime Selection Strategy**

The system automatically chooses the best runtime:

1. **Docker** (if available) - More mature ecosystem
2. **Podman** (if Docker unavailable) - Modern alternative
3. **Manual override** - Force specific runtime when needed

### **Why This Order?**
- **Docker** has better tooling and integration
- **Podman** is excellent for development and CI/CD
- **Auto-detection** ensures it "just works"

## 📝 Next Steps

1. **Start the application**: `./run.sh run` or `make run`
2. **Pull AI models**: `make pull-ollama`
3. **Access the web interface**: http://localhost:8080
4. **Upload novels** and start asking questions!

For more information, see the main [README.md](README.md) file.

