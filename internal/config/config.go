package config

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App    AppConfig
	DB     DBConfig
	Redis  RedisConfig
	JWT    JWTConfig
	Upload UploadConfig
	Rate   RateConfig
	CORS   CORSConfig
	Audit  AuditConfig
}

type AppConfig struct {
	Env  string
	Port int
	Name string
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

func (d DBConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode)
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret         string
	AccessTTLMin   int
	RefreshTTLHour int
}

type UploadConfig struct {
	Dir           string
	MaxMB         int
	PublicBaseURL string
}

func (u UploadConfig) MaxBytes() int64 {
	return int64(u.MaxMB) * 1024 * 1024
}

type RateConfig struct {
	PerMin int
}

type CORSConfig struct {
	Origins []string
}

type AuditConfig struct {
	Enabled bool
}

func Load(path string) (*Config, error) {
	_ = godotenv.Load(path) // ignore if .env missing

	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	setDefaults(v)

	if err := v.BindEnv("App.Env", "APP_ENV"); err != nil {
		return nil, err
	}
	if err := v.BindEnv("App.Port", "APP_PORT"); err != nil {
		return nil, err
	}
	if err := v.BindEnv("App.Name", "APP_NAME"); err != nil {
		return nil, err
	}
	if err := v.BindEnv("DB.Host", "DB_HOST"); err != nil {
		return nil, err
	}
	if err := v.BindEnv("DB.Port", "DB_PORT"); err != nil {
		return nil, err
	}
	if err := v.BindEnv("DB.User", "DB_USER"); err != nil {
		return nil, err
	}
	if err := v.BindEnv("DB.Password", "DB_PASSWORD"); err != nil {
		return nil, err
	}
	if err := v.BindEnv("DB.Name", "DB_NAME"); err != nil {
		return nil, err
	}
	if err := v.BindEnv("DB.SSLMode", "DB_SSLMODE"); err != nil {
		return nil, err
	}
	if err := v.BindEnv("Redis.Addr", "REDIS_ADDR"); err != nil {
		return nil, err
	}
	if err := v.BindEnv("Redis.Password", "REDIS_PASSWORD"); err != nil {
		return nil, err
	}
	if err := v.BindEnv("Redis.DB", "REDIS_DB"); err != nil {
		return nil, err
	}
	if err := v.BindEnv("JWT.Secret", "JWT_SECRET"); err != nil {
		return nil, err
	}
	if err := v.BindEnv("JWT.AccessTTLMin", "JWT_ACCESS_TTL_MIN"); err != nil {
		return nil, err
	}
	if err := v.BindEnv("JWT.RefreshTTLHour", "JWT_REFRESH_TTL_HOUR"); err != nil {
		return nil, err
	}
	if err := v.BindEnv("Upload.Dir", "UPLOAD_DIR"); err != nil {
		return nil, err
	}
	if err := v.BindEnv("Upload.MaxMB", "MAX_UPLOAD_MB"); err != nil {
		return nil, err
	}
	if err := v.BindEnv("Upload.PublicBaseURL", "PUBLIC_BASE_URL"); err != nil {
		return nil, err
	}
	if err := v.BindEnv("Rate.PerMin", "RATE_LIMIT_PER_MIN"); err != nil {
		return nil, err
	}
	if err := v.BindEnv("CORS.Origins", "CORS_ORIGINS"); err != nil {
		return nil, err
	}
	if err := v.BindEnv("Audit.Enabled", "AUDIT_ENABLED"); err != nil {
		return nil, err
	}

	c := &Config{}
	if err := v.Unmarshal(c); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	c.CORS.Origins = splitComma(v.GetString("CORS.Origins"))
	return c, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("App.Env", "development")
	v.SetDefault("App.Port", 8080)
	v.SetDefault("App.Name", "go-tutorials")
	v.SetDefault("DB.Host", "localhost")
	v.SetDefault("DB.Port", 5432)
	v.SetDefault("DB.User", "postgres")
	v.SetDefault("DB.Password", "postgres_password")
	v.SetDefault("DB.Name", "claude_code_flutter")
	v.SetDefault("DB.SSLMode", "disable")
	v.SetDefault("Redis.Addr", "localhost:6379")
	v.SetDefault("Redis.DB", 0)
	v.SetDefault("JWT.AccessTTLMin", 15)
	v.SetDefault("JWT.RefreshTTLHour", 168)
	v.SetDefault("Upload.Dir", "./uploads")
	v.SetDefault("Upload.MaxMB", 10)
	v.SetDefault("Upload.PublicBaseURL", "http://localhost:8080")
	v.SetDefault("Rate.PerMin", 60)
	v.SetDefault("CORS.Origins", "http://localhost:3000,http://localhost:5173")
	v.SetDefault("Audit.Enabled", true)
}

func splitComma(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
