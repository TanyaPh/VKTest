package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	fm "task/flood-monitoring"
	"task/repository/postgres"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func main() {
	userID, requestsExpected := flagsParse()
	if err := initConfig(); err != nil {
		log.Fatalf("failed to intialize config: %s", err.Error())
	}
	repos, err := postgres.New(postgres.Config{
		Host:    viper.GetString("db.host"),
		Port:    viper.GetString("db.port"),
		DBName:  viper.GetString("db.dbname"),
		SSLMode: viper.GetString("db.sslmode"),
	})
	if err != nil {
		log.Fatal(err)
	}

	k := viper.GetInt("checks_limit")
	n := viper.GetInt("control_duration")
	controller := fm.NewFloodMonitoring(repos, k, n)

	for i := 0; i < requestsExpected; i++ {
		res, err := controller.Check(context.Background(), userID)
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
	flag.Int64Var(&userID, "id", 0, "user ID")
	flag.IntVar(&limit, "r", 0, "requests expected")
	flag.Parse()

	return userID, limit
}
