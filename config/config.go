package config

import (
	"fmt"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"regexp"

	"foodtastechess/directory"
	"foodtastechess/logger"
)

type ConfigProvider interface {
	PreProvide(provide directory.Provider) error
}

type viperProvider struct {
	log *logging.Logger

	name  string
	paths []string

	configTypes  map[string]reflect.Type
	configValues map[string]interface{}
}

var (
	defaultPaths []string = []string{"./"}
)

func NewConfigProvider(nameAndPaths ...string) ConfigProvider {
	var cfg *viperProvider = new(viperProvider)
	cfg.initNameAndPaths(nameAndPaths...)
	cfg.readConfig()
	cfg.initLogging()
	cfg.initSections()
	cfg.initValues()
	return cfg
}

func (cfg *viperProvider) PreProvide(provide directory.Provider) error {
	cfg.log.Debug("Starting Pre-Provide")
	for section, value := range cfg.configValues {
		serviceName := fmt.Sprintf("%sConfig", section)

		config := reflect.ValueOf(value).Elem().Interface()
		cfg.log.Debug("Providing config \"%s\"", serviceName)
		err := provide(serviceName, config)
		if err != nil {
			return err
		}
	}
	cfg.log.Debug("")

	return nil
}

func (cfg *viperProvider) initNameAndPaths(nameAndPaths ...string) {
	switch len(nameAndPaths) {
	case 0:
		cfg.name = "config"
		cfg.paths = defaultPaths
	case 1:
		cfg.name = nameAndPaths[0]
		cfg.paths = defaultPaths
	default:
		cfg.name = nameAndPaths[0]
		cfg.paths = nameAndPaths[1:]
	}
}

func (cfg *viperProvider) initLogging() {
	var C logger.LoggerConfig
	err := viper.MarshalKey("logger", &C)
	if err != nil {
		panic(fmt.Errorf("Can't parse: %s \n", err))
	}

	logger.InitLog(C)

	cfg.log = logger.Log("config")
}

func (cfg *viperProvider) initSections() {
	cfg.configTypes = make(map[string]reflect.Type)

	cfg.addSection("server", ServerConfig{})
	cfg.addSection("session", SessionConfig{})
	cfg.addSection("database", DatabaseConfig{})
	cfg.addSection("auth", AuthConfig{})
	cfg.addSection("cache", QueriesCacheConfig{})
	cfg.addSection("fixtures", FixturesConfig{})
}

func (cfg *viperProvider) addSection(section string, configStruct interface{}) {
	var configType reflect.Type
	configType = reflect.TypeOf(configStruct)
	if configType.Kind() != reflect.Struct {
		cfg.log.Error(
			fmt.Sprintf(
				"Could not add section with name \"%s\", value "+
					"passed for `configStruct` not a struct",
				section,
			),
		)
	}
	cfg.configTypes[section] = reflect.TypeOf(configStruct)
}

func (cfg *viperProvider) readConfig() {
	viper.SetConfigName(cfg.name)
	for _, path := range cfg.paths {
		viper.AddConfigPath(path)
	}

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func (cfg *viperProvider) initValues() {
	cfg.configValues = make(map[string]interface{})
	for _, key := range viper.AllKeys() {
		configType, present := cfg.configTypes[key]
		if !present {
			continue
		}

		config := reflect.New(configType).Interface()
		err := viper.MarshalKey(key, config)
		if err != nil {
			cfg.log.Fatalf("Can't parse section %s, got: %v", key, err)
		}

		cfg.substituteEnvVars(config)

		cfg.configValues[key] = config
	}
}

func (cfg *viperProvider) substituteEnvVars(config interface{}) {
	// config is a pointer to a struct (assertion made in
	// initialization)

	re := regexp.MustCompile("\\$(.*)")

	value := reflect.ValueOf(config).Elem()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)

		if field.Type().Kind() != reflect.String {
			continue
		}

		str := field.String()
		submatch := re.FindStringSubmatch(str)
		if submatch == nil {
			continue
		}

		envVar := submatch[1]
		newValue := os.Getenv(envVar)
		field.SetString(newValue)

	}
}

type ConfigTestProvider struct {
	ConfigProvider ConfigProvider
}

func (c *ConfigTestProvider) InitTestConfig() {
	c.ConfigProvider = NewConfigProvider("testconfig", "../")
}
