# Novel Q&A Assistant

**Novel Q&A Assistant** is a web application that allows users to upload novels in `.txt` format and ask questions about their content using local LLMs (Large Language Models) via [Ollama](https://ollama.com/). The app chunks uploaded novels, stores them in a simple vector database (ChromaDB-like), and uses context retrieval to provide accurate answers.

---

## Features

- 📤 Upload `.txt` novels via web interface
- 🔍 Ask questions about uploaded novels
- 🤖 Uses local LLMs (phi3, llama3, mistral, gemma) via Ollama
- 🧠 Simple context retrieval using keyword search (can be extended to real embeddings)

---

## Getting Started

### Prerequisites

- Go 1.24+
- Ollama running locally (or in Docker)

### Running Locally

1. **Start Ollama**  
   Make sure Ollama is running on [http://localhost:11434](http://localhost:11434) and has your desired models pulled (e.g., `ollama pull phi3`).

2. **Build and Run the App**
   ```sh
   go build -o novel-qa .
   ./novel-qa
   ```
   The app will be available at [http://localhost:8080](http://localhost:8080).


---

## Usage

1. **Upload a Novel**  
   Use the "Upload New Novel" section to upload a `.txt` file.

2. **Ask Questions**  
   Enter your question and select a model. The app retrieves relevant context and queries the LLM for an answer.

---

## Project Structure

- `main.go` — Entry point, sets up routes and services
- `handlers` — HTTP handlers for Q&A and uploads
- `models` — Request/response models
- `services` — Core logic: novel chunking, context retrieval, Ollama API
- `templates` — HTML templates
- `static` — CSS and static assets
- `novels/` — Uploaded novels (created at runtime)
- `chroma_db/` — Simple vector DB (created at runtime)

---

## License

This project is licensed under the GNU GPL v3.

---

## Acknowledgements

- [Ollama](https://ollama.com/) for local LLM serving
- [Gin](https://github.com/gin-gonic/gin) web framework

---

## Author

[Eusuf Kanchwala](https://github.com/kweusuf)

---

**Happy reading & questioning!**