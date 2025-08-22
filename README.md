# Novel Q&A Assistant

**Novel Q&A Assistant** is a web application that allows users to upload novels in `.txt` and `.epub` formats and ask questions about their content using local LLMs (Large Language Models) via [Ollama](https://ollama.com/). The app chunks uploaded novels, stores them in a simple vector database (ChromaDB-like), and uses context retrieval to provide accurate answers.

**🎉 New Feature: EPUB Support!** You can now upload EPUB books directly - the app will automatically extract text content from EPUB files for Q&A processing.

---

## Features

- 📤 Upload `.txt` and `.epub` novels via web interface
- 📖 **EPUB Processing**: Automatic text extraction from EPUB files using Go's standard library
- 🔍 Ask questions about uploaded novels (both TXT and EPUB)
- 🤖 Uses local LLMs (phi3, llama3, mistral, gemma) via Ollama
- 🧠 Simple context retrieval using keyword search (can be extended to real embeddings)
- ✨ **HTML Tag Cleaning**: Removes formatting tags from EPUB content for clean text processing

---

## Getting Started

### Prerequisites

- Go 1.24+
- Ollama running locally

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
   - Use the "Upload New Novel" section to upload `.txt` or `.epub` files
   - The app supports both plain text files and EPUB eBooks
   - EPUB files are automatically processed to extract readable text content

2. **Ask Questions**
   - Enter your question about the uploaded novels
   - Select your preferred AI model (phi3, llama3, mistral, or gemma)
   - The app retrieves relevant context and queries the LLM for an answer

### EPUB Processing Details

When you upload an EPUB file, the app:
- Opens the EPUB as a ZIP archive (EPUBs are ZIP files with a specific structure)
- Extracts HTML/XHTML content files from the EPUB
- Removes HTML tags and formatting to get clean, readable text
- Processes the text into chunks for efficient Q&A
- Stores the processed content in the database for future questions

**Supported EPUB Features:**
- ✅ Standard EPUB 2.0 and 3.0 formats
- ✅ HTML and XHTML content extraction
- ✅ Automatic HTML tag removal
- ✅ Chapter and section preservation
- ✅ Metadata extraction (title, author, etc.)

**Note:** The app uses Go's standard `archive/zip` library, so no external dependencies are required for EPUB processing.

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
