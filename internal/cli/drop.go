package cli

import (
	"filechat/internal/cache"
	"filechat/internal/config"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var dropConf = config.Init()
var DropCmd = &cobra.Command{
	Use:   "drop [command]",
	Short: "Drop indexed documents cache",
	Run:   runDrop,
}

func InitDrop() {
	conf := dropConf
	DropCmd.Flags().StringVar(
		&conf.CacheDirPrefix,
		"cache-dir-prefix",
		conf.CacheDirPrefix,
		"Cache dir prefix",
	)
}

func runDrop(cmd *cobra.Command, args []string) {
	conf := dropConf
	cache := cache.New(conf.CacheDirPrefix)
	if err := cache.Drop(); err != nil {
		log.Fatal(err)
	}
	log.Info("Cache dropped 👌")
}
