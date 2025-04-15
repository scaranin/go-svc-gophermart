package config

import (
	"flag"

	"github.com/caarlos0/env"
)

// Конфигурация сервиса GopherMart
type ConfigGM struct {
	// Адрес и порт запуска сервиса
	ServerURL string `env:"RUN_ADDRESS"`
	// Адрес подключения к базе данных
	DSN string `env:"DATABASE_URI"`
	// Адрес системы расчёта начислений
	ASAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

// Возвращает предзаполненную структуру по умолчанию
func GetDefaultConfigGM() ConfigGM {
	return ConfigGM{
		ServerURL: "localhost:8080",
		DSN:       "postgres://postgres:admin@localhost:5432/postgres",
		ASAddress: "http://localhost:8080",
	}

}

// Считываем параметры конфигурации в приоритете:
//
// - Из флагов запуска
//
// - Из переменных окружения
func NewConfig() (ConfigGM, error) {
	var Cfg ConfigGM

	err := env.Parse(&Cfg)

	if err != nil {
		return Cfg, err
	}

	defCfgGM := GetDefaultConfigGM()

	if flag.Lookup("a") == nil {
		flag.StringVar(&defCfgGM.ServerURL, "a", "localhost:8080", "RUN_ADDRESS")
	}
	if flag.Lookup("d") == nil {
		flag.StringVar(&defCfgGM.DSN, "d", "postgres://postgres:admin@localhost:5432/postgres", "DATABASE_URI")
	}
	if flag.Lookup("r") == nil {
		flag.StringVar(&defCfgGM.ASAddress, "r", "http://localhost:8080", "ACCRUAL_SYSTEM_ADDRESS")
	}
	flag.Parse()

	if len(Cfg.ServerURL) == 0 {
		Cfg.ServerURL = defCfgGM.ServerURL
	}

	if len(Cfg.DSN) == 0 {
		Cfg.DSN = defCfgGM.DSN
	}

	if len(Cfg.ASAddress) == 0 {
		Cfg.ASAddress = defCfgGM.ASAddress
	}

	return Cfg, err
}
