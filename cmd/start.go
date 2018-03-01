package cmd

import (
	"log"
	"os"
	"os/signal"

	"github.com/oneut/lrp/livereloadproxy"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start live reload proxy",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func run() {
	log.Println("Start live reload proxy")
	lrp := livereloadproxy.NewLivereloadProxy()
	lrp.Run()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	for {
		select {
		case <-sigChan:
			lrp.Stop()
			os.Exit(0)
		}
	}
}
