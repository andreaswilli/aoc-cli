package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
)

type Config struct {
	Engines []*Engine `json:"engines"`
}

type EngineManager struct {
	engines map[string]*Engine
}

func NewEngineManager(fsys fs.FS) (*EngineManager, error) {
	config := Config{}
	jsonBytes, err := fs.ReadFile(fsys, "aoc-cli.json")

	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, errors.New(
				"'aoc-cli.json' not found, it must exist in the root of your project",
			)
		}
		return nil, err
	}

	json.Unmarshal(jsonBytes, &config)

	if len(config.Engines) == 0 {
		return nil, errors.New("no engines are defined in 'aoc-cli.json'")
	}

	engines := make(map[string]*Engine)

	for _, engine := range config.Engines {
		engines[engine.Name] = engine
	}

	return &EngineManager{engines: engines}, nil
}

func (em *EngineManager) Get(name string) (*Engine, error) {
	engine, ok := em.engines[name]

	if !ok {
		return nil, fmt.Errorf("engine %q not found", name)
	}

	return engine, nil
}

func (em *EngineManager) FindAppropriateEngine(fsys fs.FS) *Engine {
	for _, engine := range em.engines {
		if _, err := fs.Stat(fsys, engine.EntryFile); err == nil && !os.IsNotExist(err) {
			return engine
		}
	}
	return nil
}
