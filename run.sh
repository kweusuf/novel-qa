#!/bin/bash

# Novel Q&A Assistant Runtime Selector
# This script helps you choose between Docker and Podman

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to show available runtimes
show_runtimes() {
    echo -e "${BLUE}🔍 Available container runtimes:${NC}"
    if command_exists docker; then
        echo -e "  ✅ Docker: $(which docker)"
    else
        echo -e "  ❌ Docker: Not found"
    fi
    
    if command_exists podman; then
        echo -e "  ✅ Podman: $(which podman)"
    else
        echo -e "  ❌ Podman: Not found"
    fi
    
    if command_exists docker-compose; then
        echo -e "  ✅ Docker Compose: $(which docker-compose)"
    else
        echo -e "  ❌ Docker Compose: Not found"
    fi
    
    if command_exists podman-compose; then
        echo -e "  ✅ Podman Compose: $(which podman-compose)"
    else
        echo -e "  ❌ Podman Compose: Not found"
    fi
    echo ""
}

# Function to show help
show_help() {
    echo -e "${BLUE}📚 Novel Q&A Assistant - Runtime Selector${NC}"
    echo ""
    echo "Usage: $0 [COMMAND] [RUNTIME]"
    echo ""
    echo "Commands:"
    echo "  build     Build the container image"
    echo "  run       Start the application"
    echo "  stop      Stop the application"
    echo "  logs      View logs"
    echo "  status    Show container status"
    echo "  clean     Clean up everything"
    echo "  help      Show this help message"
    echo ""
    echo "Runtimes:"
    echo "  docker    Force use of Docker"
    echo "  podman    Force use of Podman"
    echo "  auto      Auto-detect (default)"
    echo ""
    echo "Examples:"
    echo "  $0 build          # Build with auto-detected runtime"
    echo "  $0 run docker     # Run with Docker"
    echo "  $0 run podman     # Run with Podman"
    echo "  $0 build auto     # Build with auto-detected runtime"
    echo ""
    show_runtimes
}

# Function to run make command with specific runtime
run_make() {
    local command=$1
    local runtime=$2
    
    case $runtime in
        "docker")
            if ! command_exists docker; then
                echo -e "${RED}❌ Docker not found!${NC}"
                exit 1
            fi
            if ! command_exists docker-compose; then
                echo -e "${RED}❌ Docker Compose not found!${NC}"
                exit 1
            fi
            echo -e "${GREEN}🐳 Using Docker runtime${NC}"
            make DOCKER=docker COMPOSE=docker-compose $command
            ;;
        "podman")
            if ! command_exists podman; then
                echo -e "${RED}❌ Podman not found!${NC}"
                exit 1
            fi
            if ! command_exists podman-compose; then
                echo -e "${RED}❌ Podman Compose not found!${NC}"
                exit 1
            fi
            echo -e "${GREEN}🦭 Using Podman runtime${NC}"
            make DOCKER=podman COMPOSE=podman-compose $command
            ;;
        "auto"|"")
            echo -e "${GREEN}🔍 Auto-detecting runtime...${NC}"
            make $command
            ;;
        *)
            echo -e "${RED}❌ Unknown runtime: $runtime${NC}"
            echo -e "${YELLOW}💡 Use 'docker', 'podman', or 'auto'${NC}"
            exit 1
            ;;
    esac
}

# Main script logic
main() {
    # Show help if no arguments
    if [ $# -eq 0 ]; then
        show_help
        exit 0
    fi
    
    # Parse command and runtime
    local command=""
    local runtime="auto"
    
    case $1 in
        "help"|"-h"|"--help")
            show_help
            exit 0
            ;;
        "build"|"run"|"stop"|"logs"|"status"|"clean"|"build-dev"|"run-dev"|"stop-dev"|"logs-dev"|"shell"|"shell-dev"|"pull-ollama"|"restart"|"restart-dev")
            command=$1
            runtime=${2:-"auto"}
            ;;
        *)
            echo -e "${RED}❌ Unknown command: $1${NC}"
            echo -e "${YELLOW}💡 Run '$0 help' for usage information${NC}"
            exit 1
            ;;
    esac
    
    # Show runtime info
    show_runtimes
    
    # Run the command
    run_make $command $runtime
}

# Run main function with all arguments
main "$@"
