# filechat

Chat with PDF files using local RAG powered by Ollama.

## How it works

When you give `filechat` one or more PDF files:
1. It loads each file and extracts text content
2. Splits text into small overlapping chunks
3. Creates embeddings for each of the chunk

When you ask a question:
1. Your question is turned into an embedding
2. Relevant chunks collected using similarity search
3. Thos chunks along with your question sent to LLM
4. LLM generates an answer based on that context

## Usage

```sh
filechat [command] [files...] [--flags]

COMMANDS

completion [command]       Generate the autocompletion script for the specified shell
drop [command] [--flags]   Drop indexed documents cache
help [command]             Help about any command

FLAGS

--cache-dir-prefix         Cache dir prefix (filechat)
-l --chunk-token-length    Chunk token length (256)
-o --chunk-token-overlap   Chunk token overlap (64)
--embedding-length         Embedding length (768)
-e --embedding-model       Embedding model
-h --help                  Help for filechat
--ollama-endpoint          Ollama endpoint url (http://localhost:11434)
--rag-prompt-template      Rag prompt template (
                            Answer the question using the context below.

                            Context:
                            {{ context }}

                            Question:
                            {{ question }}
                            )
-r --reasoning-model       Reasoning model
-n --search-keep-n-chunks  Number of search results to keep (3)
-v --version               Version for filechat
```

## Examples

Ask something about gardening:

```sh
go run ./cmd/cli/main.go \
	--reasoning-model=gemma4:e2b \
	--embedding-model=embeddinggemma:300m \
	./organic_gardening.pdf ./pest_control.pdf
```
