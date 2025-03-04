package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dop251/goja"
	"github.com/samber/lo"
)

type Service struct {
	Program  *goja.Program
	FilePath string
}

var (
	services     = make(map[string]*Service)
	servicesLock sync.RWMutex
)

func loadServicesAndWatch(dir string) {
	_, serviceFiles := lo.Must2(walkDir(cfg.ServicesDir))

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

	go watchServicesDir(cfg.ServicesDir, "js")
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

	program, err := goja.Compile(filePath, string(data), false)
	if err != nil {
		return err
	}

	services[serviceName] = &Service{
		Program:  program,
		FilePath: filePath,
	}
	return nil
}

type watchexecTag struct {
	Kind     string `json:"kind"`
	Source   string `json:"source"`
	Simple   string `json:"simple"`
	Full     string `json:"full"`
	Absolute string `json:"absolute"`
	Filetype string `json:"filetype"`
}

type watchexecEvent struct {
	Tags []watchexecTag `json:"tags"`
}

type event struct {
	Source   string
	FsSimple string
	FsFull   string
	Path     string
}

func parseEvent(data []byte) (*event, error) {
	we := watchexecEvent{}
	if err := json.Unmarshal(data, &we); err != nil {
		return nil, err
	}

	if len(we.Tags) < 3 {
		return nil, fmt.Errorf("invalid event: %v", we)
	}

	return &event{
		Source:   we.Tags[0].Source,
		FsSimple: we.Tags[1].Simple,
		FsFull:   we.Tags[1].Full,
		Path:     we.Tags[2].Absolute,
	}, nil
}

func watchServicesDir(dir string, ext string) {
	dirPrefix := lo.Must(filepath.Abs(dir))

	args := []string{
		"-w", dir,
		"-e", ext,
		"--emit-events-to=json-stdio",
		"--only-emit-events",
	}

	slog.Info("watchexec", "args", args)
	watchCmd := exec.Command(cfg.Watchexec, args...)

	out := lo.Must(watchCmd.StdoutPipe())
	defer out.Close()

	lo.Must0(watchCmd.Start())

	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		e := lo.Must(parseEvent(scanner.Bytes()))
		slog.Info("watchServicesDir()", "event", e)

		switch e.FsSimple {
		case "create", "modify", "remove":
			serviceName := strings.TrimPrefix(e.Path, dirPrefix)
			serviceName = strings.TrimPrefix(serviceName, "/")

			servicesLock.Lock()
			_, err := os.Stat(e.Path)
			if err != nil {
				delete(services, serviceName)
			} else {
				err = loadService(serviceName, e.Path)
			}
			servicesLock.Unlock()

			if err != nil {
				slog.Error("watchServicesDir() loadService", "path", e.Path, "err", err)
			}
		default:
			slog.Warn("watchServicesDir() unknown event", "event", e)
		}
	}
}
