package main

import (
	"flag"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", `configs\config.yaml`, "Path to configuration file")
}

func main() {

	//config, err := config.New(configFile)
	//if err != nil {
	//	log.Fatalln(err)
	//}

	//logg := logger.New(config.Logger.Level)


	//fmt.Println(config)

	//storage := memorystorage.New()
	//calendar := app.New(logg, storage)
	//
	//server := internalhttp.NewServer(calendar)
	//
	//go func() {
	//	signals := make(chan os.Signal, 1)
	//	signal.Notify(signals)
	//
	//	<-signals
	//	signal.Stop(signals)
	//
	//	if err := server.Stop(); err != nil {
	//		logger.Error("failed to stop http server: " + err.String())
	//	}
	//}()
	//
	//if err := server.Start(); err != nil {
	//	logger.Error("failed to start http server: " + err.String())
	//	os.Exit(1)
	//}
}
