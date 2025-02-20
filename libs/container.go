package libs

import (
	"fmt"
	"log"
	"os"
	"sync"
	"text/template"

	"github.com/fungolang/screw"
	"github.com/spf13/viper"
)

var (
	onceEnv   sync.Once
	globalEnv *Container
)

type Container struct {
	ConfigPath string   `screw:"-c;--ConfigPath" usage:"config file path"`
	WorkDir    string   `screw:"-w;--WorkDir" usage:"work dir"`
	BackupDir  string   `screw:"-z;--BackupDir" usage:"backup dir"`
	TmplDir    string   `screw:"-l;--TmplDir" usage:"templates dir"`
	Operation  string   `screw:"-o;--Operation" usage:"Operation command"`
	ManiFests  []string `screw:"-f;--ManiFests" usage:"ManiFests name"`
	ImageAPI   string   `screw:"-a;--ImageAPI" usage:"oneapi image name"`
	ImagePG    string   `screw:"-g;--ImagePG" usage:"postgres image name"`
	ImageMG    string   `screw:"-m;--ImageMG" usage:"mongo image name"`
	ImageMySql string   `screw:"-q;--ImageMySql" usage:"mysql image name"`
	ImageGPT   string   `screw:"-t;--ImageGPT" usage:"fastgpt image name"`
	BaseURL    string   `screw:"-b;--BaseURL" usage:"opengpt api url"`
	GptPass    string   `screw:"-s;--GptPass" usage:"database password"`
	ApiKey     string   `screw:"-k;--ApiKey" usage:"opengpt api key"`
	RootKey    string   `screw:"-r;--RootKey" usage:"fastgpt root key"`
	DbUser     string   `screw:"-u;--DbUser" usage:"database user"`
	DbPass     string   `screw:"-p;--DbPass" usage:"database password"`
	DataDir    string   `screw:"-d;--DataDir" usage:"data directory"`
}

func GetEnv() *Container {
	return globalEnv
}

func NewContainer() *Container {
	onceEnv.Do(func() {
		env := &Container{
			WorkDir:    "./run",
			BackupDir:  "./backup",
			TmplDir:    "./tmpl",
			ManiFests:  []string{"gpt"},
			ConfigPath: "deploy.json",
			ImageAPI:   "justsong/one-api:latest",
			ImageGPT:   "ghcr.io/labring/fastgpt:latest",
			ImagePG:    "pgvector/pgvector:0.7.0-pg15",
			ImageMG:    "mongo:5.0.18",
			ImageMySql: "mysql:8.0.36",
			GptPass:    "admin",
		}

		screw.Bind(env)
		viper.SetConfigFile(env.ConfigPath)
		if err := viper.ReadInConfig(); err != nil {
			log.Fatal("cannot read configuration")
		}
		if err := viper.Unmarshal(env); err != nil {
			log.Fatal("environment can't be loaded: ", err)
		}
		globalEnv = env
	})
	return globalEnv
}

func ApplyTemplate(filePath string, tmpl *template.Template, data any) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	return nil
}
