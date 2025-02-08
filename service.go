package main

import (
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dop251/goja"
	"github.com/fsnotify/fsnotify"
	"github.com/samber/lo"
)

type Service struct {
	VM       *goja.Runtime
	Program  *goja.Program
	FilePath string
}

var (
	services     = make(map[string]*Service)
	servicesLock sync.RWMutex
)

func loadServicesAndWatch(dir string) {
	dirsToWatch, serviceFiles := lo.Must2(walkDir(cfg.ServicesDir))

	for _, servicePath := range serviceFiles {
		serviceName := strings.TrimPrefix(servicePath, dir)
		serviceName = strings.TrimPrefix(serviceName, "/")

		servicesLock.Lock()
		err := loadService(serviceName, servicePath)
		servicesLock.Unlock()

		if err != nil {
			slog.Error("loadServicesAndWatch() loadService", "path", servicePath, "err", err)
		} else {
			slog.Info("loadServicesAndWatch() loaded service", "path", servicePath)
		}
	}

	go watchServicesDir(cfg.ServicesDir, dirsToWatch)
}

func walkDir(dir string) (subDirs []string, files []string, err error) {
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			subDirs = append(subDirs, path)
			return nil
		}

		files = append(files, path)
		return nil
	})
	return
}

func loadService(serviceName, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", false))

	program, err := goja.Compile(filePath, string(data), false)
	if err != nil {
		return err
	}

	services[serviceName] = &Service{
		VM:       vm,
		Program:  program,
		FilePath: filePath,
	}
	return nil
}

func watchServicesDir(dir string, dirsToWatch []string) {
	watcher := lo.Must(fsnotify.NewWatcher())
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				slog.Info("watchServicesDir()", "event", event)

				serviceName := strings.TrimPrefix(event.Name, dir)
				serviceName = strings.TrimPrefix(serviceName, "/")

				switch {
				case event.Has(fsnotify.Create) || event.Has(fsnotify.Write):
					if lo.Must(os.Stat(event.Name)).IsDir() {
						watcher.Add(event.Name)
						slog.Info("watcher.Add()", "watching", event.Name)
						break
					}

					if !strings.HasSuffix(event.Name, ".js") {
						slog.Info("watchServicesDir() not a js file", "name", event.Name)
						break
					}

					servicesLock.Lock()
					err := loadService(serviceName, event.Name)
					servicesLock.Unlock()

					if err != nil {
						slog.Error("watchServicesDir() loadService", "name", event.Name, "err", err)
					} else {
						slog.Info("watchServicesDir() loaded service", "name", event.Name)
					}
				case event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename):
					servicesLock.Lock()
					delete(services, serviceName)
					servicesLock.Unlock()

					slog.Info("watchServicesDir() deleted service", "name", serviceName)
				default:
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				slog.Error("watchServicesDir()", "error", err)
			}
		}
	}()

	for _, d := range dirsToWatch {
		watcher.Add(d)
		slog.Info("watcher.Add()", "watching", d)
	}
	select {}
}
