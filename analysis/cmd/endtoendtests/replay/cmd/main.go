package main

import (
	"coralreefci/analysis/cmd/endtoendtests/replay"
	"coralreefci/engine/backend"
	"coralreefci/engine/ingestor"
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"net/url"
	"time"
)

func main() {
	runBacktestFlag := flag.Bool("runbacktest", false, "runs the end to end backtest")
	loadArchiveFlag := flag.Bool("loadarchive", false, "load archive into the database")
	archivePathFlag := flag.String("archivepath", "", "location of github archive")
	flag.Parse()

	if !*runBacktestFlag && !*loadArchiveFlag {
		fmt.Println("Usage: ./cmd --loadarchive=true --archivepath=/home/michael/Data/GithubArchive/")
		fmt.Println("Usage: ./cmd --runbacktest=true")
		return
	}

	bufferPool := ingestor.NewPool()
	db := ingestor.Database{BufferPool: bufferPool}
	db.Open()

	bs := replay.BacktestServer{DB: &db}
	go bs.Start()

	dispatcher := ingestor.Dispatcher{}
	dispatcher.Start(5)
	ingestorServer := ingestor.IngestorServer{}
	go ingestorServer.Start()

	if *loadArchiveFlag && *archivePathFlag != "" {
		bs.LoadArchive(*archivePathFlag)
	}

	if *runBacktestFlag {
		backendServer := backend.BackendServer{}
		client := github.NewClient(nil)
		url, _ := url.Parse("http://localhost:8000/")
		client.BaseURL = url
		client.UploadURL = url

		backendServer.Repos = new(backend.ActiveRepos)
		backendServer.Repos.Actives = make(map[int]*backend.ArchRepo)
		backendServer.Repos.Actives[26295345] = new(backend.ArchRepo)
		backendServer.Repos.Actives[26295345].Hive = new(backend.ArchHive)
		backendServer.Repos.Actives[26295345].Hive.Blender = new(backend.Blender)
		backendServer.NewModel(26295345)
		backendServer.Repos.Actives[26295345].Client = client
		go backendServer.Start()

		bs.AddRepo(26295345, "dotnet", "corefx")
		bs.AddRepo(724712, "rust-lang", "rust")
		bs.StreamWebhookEvents()
		time.Sleep(15 * time.Second)

		backendServer.OpenSQL()
		backendServer.Timer()

		time.Sleep(15 * time.Second)
	}
}
