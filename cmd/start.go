package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	"github.com/oneut/lrp/config"
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

var yamlFile string

const defaultYamlFile string = "./lrp.yml"

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().StringVarP(&yamlFile, "file", "f", "", "yaml file. default ./lrp.yml")
}

func run() {
	log.Println("Start live reload proxy")

	cfg := config.CreateConfig(getYamlBuf())
	lrp := livereloadproxy.NewLivereloadProxy(cfg)
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

func getYamlBuf() []byte {
	buf, err := ioutil.ReadFile(getYamlFile())
	if err != nil {
		panic(err)
	}

	return buf
}

func getYamlFile() string {
	if yamlFile == "" {
		return defaultYamlFile
	}

	return yamlFile
}
