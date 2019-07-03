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
	plugins []Plugin
)

func main() {
	log.Info("Starting VueHoster")

	ex, _ := os.Executable()
	currentDir = filepath.Dir(ex)
	pluginsDir = path.Join(currentDir, "plugins")
	if *flag.Bool("debug", false, "enables logging of debug") {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
		log.SetLevel(log.DebugLevel)
	}
	log.SetLevel(log.DebugLevel)

	log.Debug("currentDir:", currentDir)
	log.Debug("pluginsDir:", pluginsDir)

	handlePlugins()

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
				log.Debug("FS:", event)
				handlePlugins()

			case err := <-watcher.Errors:
				log.Fatal("Error received from file watcher", err)
			}
		}
	})

	// init plugin scan
	handlePlugins()

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}

func handlePlugins() {
	if len(plugins) != 0 {
		log.Info("Shutting down all plugins")

		for _, plug := range plugins {
			plug.Cancel()
		}
	}

	log.Info("Scanning for plugins")

	files, err := ioutil.ReadDir(pluginsDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".so" {
			log.Debug("Found plugin:", f.Name())

			// open the so file to load the symbols
			symbols, err := plugin.Open(path.Join(pluginsDir, f.Name()))
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

			plugins = append(plugins, plug)
		}
	}

	log.Info("Running all plugins")

	for _, plug := range plugins {
		g.Go(plug.Run)
	}
}