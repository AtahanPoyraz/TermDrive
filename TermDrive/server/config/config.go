package config

import (
	"log"
	"os"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

// Configuration holds the application's configuration settings, which are loaded
// from a YAML file. It includes server configuration, TermDrive settings,
// CORS settings, database settings, and the JWT secret key.
type Configuration struct {
	Server struct {
		Host     string `yaml:"host"`
		Port     int64  `yaml:"port"`
		Timeouts struct {
			ReadTimeout       int64 `yaml:"read_timeout"`
			WriteTimeout      int64 `yaml:"write_timeout"`
			IdleTimeout       int64 `yaml:"idle_timeout"`
			ReadHeaderTimeout int64 `yaml:"read_header_timeout"`
		} `yaml:"timeouts"`
	} `yaml:"server"`

	TermDrive struct {
		StoragePath string `yaml:"storage_path"`
		UploadSize  int    `yaml:"upload_size"`
	} `yaml:"term_drive"`

	Cors struct {
		AllowedOrigins []string `yaml:"allowed_origins"`
		AllowedMethods []string `yaml:"allowed_methods"`
		AllowedHeaders []string `yaml:"allowed_headers"`
	} `yaml:"cors"`

	Database struct {
		Driver   string `yaml:"driver"`
		Host     string `yaml:"host"`
		Port     int64  `yaml:"port"`
		Name     string `yaml:"name"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		SSLMode  string `yaml:"ssl_mode"`
		TimeZone string `yaml:"time_zone"`
	} `yaml:"database"`

	Security struct {
		JwtSecretKey string `yaml:"jwt_secret_key"`
	} `yaml:"security"`
}

// LoadConfig loads the configuration from a YAML file located at the given path
// and unmarshals it into a Configuration struct. Returns the configuration struct
// and any error encountered while reading or unmarshaling the file.
func LoadConfig(path string) (*Configuration, error) {
	config := &Configuration{}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(content, config); err != nil {
		return nil, err
	}

	return config, nil
}

// Dependencies holds the application dependencies, including the application
// context and system context. This is used to provide required services
// across different components of the application.
type Dependencies struct {
	AppContext    *AppContext
	SystemContext *SystemContext
}

// AppContext holds the application's configuration and logger, providing access
// to application settings and logging functionality.
type AppContext struct {
	Configuration *Configuration
	Logger        *log.Logger
}

// SystemContext holds the router and database connection for routing and database
// operations within the system. It is used to pass these dependencies throughout
// the system.
type SystemContext struct {
	Router   *mux.Router
	Database *gorm.DB
}

// NewAppContext initializes and returns a new AppContext that holds the application's
// configuration and logger for logging purposes.
func NewAppContext(configuration *Configuration, logger *log.Logger) *AppContext {
	return &AppContext{Configuration: configuration, Logger: logger}
}

// NewSystemContext initializes and returns a new SystemContext that holds the router
// and database connection used by the system for routing and database operations.
func NewSystemContext(router *mux.Router, database *gorm.DB) *SystemContext {
	return &SystemContext{Router: router, Database: database}
}

// contextKey is a custom type used to define keys for storing and retrieving
// values from the context.
type contextKey string

// UserContextKey is a constant representing the key used to store the user
// in the context.
const UserContextKey = contextKey("user")
