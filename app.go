package main

import (
	"flag"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"plugin"
)

type Plugin interface {
	Run() error
	Cancel()
}

var (
	g errgroup.Group
	currentDir string
	pluginsDir string
	plugins map[string]Plugin
)

func main() {
	startup()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Failed to create new watcher", err)
	}
	defer watcher.Close()

	err = watcher.Add(pluginsDir)
	if err != nil {
		log.Fatal("Failed to watch directory", err)
	}

	g.Go(func() error {
		for {
			select {
			case event:= <-watcher.Events:
				log.Debug("FS:", event.String())
				updatePlugin(event.Name)

			case err := <-watcher.Errors:
				log.Fatal("Error received from file watcher", err)
			}
		}
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}

func startup() {
	log.Info("Starting VueHoster")

	plugins = make(map[string]Plugin)
	ex, _ := os.Executable()
	currentDir = filepath.Dir(ex)
	pluginsDir = path.Join(currentDir, "plugins")

	debugFlag := flag.Bool("debug", false, "enables logging of debug")
	flag.Parse()
	if *debugFlag {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
		log.SetLevel(log.DebugLevel)
	}

	log.Debug("currentDir:", currentDir)
	log.Debug("pluginsDir:", pluginsDir)

	log.Info("Scanning for plugins")

	files, err := ioutil.ReadDir(pluginsDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		loadPlugin(f.Name())
	}
}

func updatePlugin(name string) {
	plug, exists := plugins[name]

	if exists {
		plug.Cancel()
	}

	loadPlugin(name)
}

func loadPlugin(name string) {
	if filepath.Ext(name) == ".so" {
		log.Debug("Plugin plugin:", name)

		// open the so file to load the symbols
		symbols, err := plugin.Open(path.Join(pluginsDir, name))
		if err != nil {
			log.Fatal("Failed to open plugin", err)
		}

		// look up a symbol
		symb, err := symbols.Lookup("Plugin")
		if err != nil {
			log.Fatal("Plugin does not contain required functions")
		}

		// assert symbol has correct type
		var plug Plugin
		plug, ok := symb.(Plugin)
		if !ok {
			log.Fatal("unexpected type from plugin symbol")
		}

		plugins[name] = plug
		g.Go(plug.Run)
	}
}