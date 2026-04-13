package cli

import (
	"encoding/json"
	"fmt"

	"filechat/internal/cache"
	"filechat/internal/chat"
	"filechat/internal/config"
	"filechat/internal/extract"
	"filechat/internal/ingest"
	"filechat/internal/ollama"
	"filechat/internal/search"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var chatConf = config.Init()
var ChatCmd = &cobra.Command{
	Use:   "filechat [files...]",
	Short: "chat with pdf files",
	Args:  cobra.MinimumNArgs(1),
	Run:   runChat,
}

func InitChat() {
	conf := chatConf
	ChatCmd.Flags().StringVar(
		&conf.CacheDirPrefix,
		"cache-dir-prefix",
		conf.CacheDirPrefix,
		"Cache dir prefix",
	)
	ChatCmd.Flags().StringVar(
		&conf.OllamaEndpoint,
		"ollama-endpoint",
		conf.OllamaEndpoint,
		"Ollama endpoint url",
	)
	ChatCmd.Flags().IntVar(
		&conf.EmbeddingLength,
		"embedding-length",
		conf.EmbeddingLength,
		"Embedding length",
	)
	ChatCmd.Flags().IntVarP(
		&conf.ChunkTokenLength,
		"chunk-token-length", "l",
		conf.ChunkTokenLength,
		"Chunk token length",
	)
	ChatCmd.Flags().IntVarP(
		&conf.ChunkTokenOverlap,
		"chunk-token-overlap", "o",
		conf.ChunkTokenOverlap,
		"Chunk token overlap",
	)
	ChatCmd.Flags().IntVarP(
		&conf.SearchKeepNChunks,
		"search-keep-n-chunks", "n",
		conf.SearchKeepNChunks,
		"Number of search results to keep",
	)
	ChatCmd.Flags().StringVar(
		&conf.RagPromptTemplate,
		"rag-prompt-template",
		conf.RagPromptTemplate,
		"RAG prompt template",
	)
	ChatCmd.Flags().StringVarP(
		&conf.ReasoningModel,
		"reasoning-model", "r",
		conf.ReasoningModel,
		"Reasoning model",
	)
	ChatCmd.Flags().StringVarP(
		&conf.EmbeddingModel,
		"embedding-model", "e",
		conf.EmbeddingModel,
		"Embedding model",
	)
	ChatCmd.MarkFlagRequired("reasoning-model")
	ChatCmd.MarkFlagRequired("embedding-model")
}

func runChat(cmd *cobra.Command, args []string) {
	conf := chatConf
	ctx := cmd.Context()
	ai, err := ollama.New(ctx, conf)
	if err != nil {
		log.Fatal(err)
	}
	search := search.New()
	ingester := ingest.New(conf, ai)
	cache := cache.New(conf.CacheDirPrefix)

	db := make(ingest.Registry)
	log.Info("Indexing documents...")
	for _, path := range args {
		text, err := extract.PDF(path)
		if err != nil {
			log.Fatal(err)
		}
		// If possible use cached data
		// instead of creating snapshot
		// for the document from scratch
		var snapshot ingest.Snapshot
		cacheKey := ingester.Hash(text)
		data, err := cache.Get(cacheKey)
		if err != nil {
			log.Fatal(err)
		} else if len(data) > 0 {
			log.Infof("Found cached data %q", path)
			err := json.Unmarshal(data, &snapshot)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			// If there is no cached data
			// build snapshot from scratch
			log.Infof("Creating embeddings %q...", path)
			res, err := ingester.Do(ctx, text)
			if err != nil {
				log.Fatal(err)
			}
			snapshot = *res
			// Cache a snapshot for future
			bytes, err := json.Marshal(res)
			if err != nil {
				log.Fatal(err)
			}
			err = cache.Set(cacheKey, bytes)
			if err != nil {
				log.Fatal(err)
			}
		}
		// Adding document to search index
		// so that we can perform a vector
		// similarity search against it
		doc := snapshot.Document
		db[doc.Hash] = snapshot
		for _, chunk := range doc.Chunks {
			search.Add(
				doc.Hash,
				chunk.ID,
				chunk.Vec,
			)
		}
	}

	session := chat.New(conf, ai, db, search)
	if err := session.Start(ctx); err != nil {
		log.Fatal(err)
	}
	fmt.Println("👋")
}
