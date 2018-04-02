package main

import (
	"log"
	"os"
	"path"
	"strings"
	"gitlab.com/blk-io/crux/api"
	"gitlab.com/blk-io/crux/config"
	"gitlab.com/blk-io/crux/enclave"
	"gitlab.com/blk-io/crux/server"
	"gitlab.com/blk-io/crux/storage"
)

func main() {

	config.InitFlags()

	args := os.Args
	if len(args) == 1 {
		exit()
	}

	for _, arg := range args[1:] {
		if strings.Contains(arg, ".conf") {
			err := config.LoadConfig(os.Args[0])
			if err != nil {
				log.Fatalln(err)
			}
			break
		}
	}
	config.ParseCommandLine()

	keyFile := config.GetString(config.GenerateKeys)
	if keyFile != "" {
		err := enclave.DoKeyGeneration(keyFile)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Key pair successfully written to %s", keyFile)
		os.Exit(0)
	}

	workDir := config.GetString(config.WorkDir)
	dbStorage := config.GetString(config.Storage)
	storagePath := path.Join(workDir, dbStorage)

	db, err := storage.Init(storagePath)
	if err != nil {
		log.Fatalf("Unable to initialise storage, error: %v\n", err)
	}

	otherNodes := config.GetStringSlice(config.OtherNodes)
	url := config.GetString(config.Url)
	if url == "" {
		log.Fatalln("URL must be specified")
	}

	pi := api.LoadPartyInfo(url, otherNodes)

	privKeyFiles := config.GetStringSlice(config.PrivateKeys)
	pubKeyFiles := config.GetStringSlice(config.PublicKeys)

	if len(privKeyFiles) != len(pubKeyFiles) {
		log.Fatalln("Private keys provided must have corresponding public keys")
	}

	if len(privKeyFiles) == 0 {
		log.Fatalln("Node key files must be provided")
	}

	enc := enclave.Init(db, pubKeyFiles, privKeyFiles, pi)

	port := config.GetInt(config.Port)
	if port < 0 {
		log.Fatalln("Port must be specified")
	}

	_, err = server.Init(enc, port)
	if err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}

	pi.PollPartyInfo()
}

func exit() {
	config.Usage()
	os.Exit(1)
}