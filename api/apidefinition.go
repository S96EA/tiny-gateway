package api

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"tiny-gateway/proxy"
	"path"
	"sync"
)

type Definition struct {
	Name    string            `bson:"name" json:"name" valid:"required"`
	Active  bool              `bson:"active" json:"active"`
	Proxy   *proxy.Definition `bson:"proxy" json:"proxy"`
}

type Definitions struct {
	definitions map[string]*Definition
	sync.Mutex
}

func NewDefinitions() *Definitions{
	return &Definitions{
		definitions: make(map[string]*Definition),
	}
}

func (ds *Definitions) Set(name string, definition *Definition) {
	ds.Lock()
	defer ds.Unlock()

	ds.definitions[name] = definition
}

func (ds *Definitions) Get(name string) *Definition {
	ds.Lock()
	defer ds.Unlock()

	return ds.definitions[name]
}

func (ds *Definitions) GetDefinitions() []*Definition{
	ds.Lock()
	defer ds.Unlock()

	var defs []*Definition
	for _, def := range ds.definitions {
		defs = append(defs, def)
	}
	return defs
}

func LoadDefinitions(dir string) (*Definitions, error){
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	definitions := NewDefinitions()

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filePath := path.Join(dir, file.Name())

		v := viper.New()
		v.SetConfigFile(filePath)
		definition := &Definition{}

		if err := v.ReadInConfig(); err != nil {
			logrus.WithError(err).Errorf("read in config failed")
			return nil, err
		}

		if err := v.Unmarshal(definition); err != nil {
			logrus.WithError(err).Errorf("unmarshal failed")
			return nil, err
		}

		if definitions.Get(definition.Name) != nil {
			err := fmt.Errorf("duplicate config: %v", definition.Name)
			logrus.WithError(err).Errorf("duplicate config")
			return nil, err
		}

		definitions.Set(definition.Name, definition)
	}

	return definitions, nil
}