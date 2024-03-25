package main

import (
	"context"
	"fmt"
	"flag"
	"log"
	"os"
	fm "task/flood-monitoring"
	"task/repository/postgres"
	"time"

	"github.com/spf13/viper"
	_ "github.com/lib/pq"
)

func main() {
	userID, checksExpected := flagsParse()
	if err := initConfig(); err != nil {
		log.Fatalf("failed to intialize config: %s", err.Error())
	}
	repos, err := postgres.New(postgres.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		log.Fatal(err)
	}

	k := viper.GetInt("checks_limit")
	n := viper.GetDuration("control_duration") * time.Second
	controller := fm.NewFloodMonitoring(repos, k)
	ctx, cancel := context.WithTimeout(context.Background(), n)
	defer cancel()
	
	for i := 1; i < checksExpected; i++ {
		res, err := controller.Check(ctx, userID)
		if err != nil {
			log.Fatal(err)
		}

		if res {
			fmt.Printf("User %d is NOT flooding\n", userID)
		} else {
			fmt.Printf("User %d is flooding\n", userID)
		}
	}
}

// FloodControl интерфейс, который нужно реализовать.
// Рекомендуем создать директорию-пакет, в которой будет находиться реализация.
type FloodControl interface {
	// Check возвращает false если достигнут лимит максимально разрешенного
	// кол-ва запросов согласно заданным правилам флуд контроля.
	Check(ctx context.Context, userID int64) (bool, error)
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func flagsParse() (int64, int) {
	var userID int64
	var limit int
	flag.Int64Var(&userID, "id", 0, "userID")
	flag.IntVar(&limit, "l", 0, "Checks Expected")
	flag.Parse()

	return userID, limit
}