package connect

import (
	"fmt"
	"log"

	"../config"
	"../structures"
	uuid "github.com/satori/go.uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var connection *gorm.DB

func InitializeDatabase() {
	connection = ConnectORM(CreateString())
}

func CreateString() string {

	var configuration config.Config
	configuration.LoadEnv()

	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=verify-full&sslrootcert=%s&options=--cluster=%s",
		configuration.Engine_sql,
		configuration.Username,
		configuration.Password,
		configuration.Host,
		configuration.Port,
		configuration.Database,
		configuration.SSL_root_cert,
		configuration.Cluster,
	)
}

func ConnectORM(stringConnection string) *gorm.DB {
	connection, err := gorm.Open(postgres.Open(stringConnection))
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return connection
}

func CloseConnection() {
	sqlDB, err := connection.DB()
	if err != nil {
		log.Fatal("No se cerro")
	}
	sqlDB.Close()
}

func GetUser(id string) structures.User {
	user := structures.User{}
	connection.Where("id = ?", id).First(&user)
	user.Scores = GetUserScores(user.Id)
	return user
}

func GetUserByUsername(username string) structures.User {
	user := structures.User{}
	connection.Where("username = ?", username).First(&user)
	user.Scores = GetUserScores(user.Id)
	return user
}

func GetUserScores(user_id uuid.UUID) []structures.Score {
	scores := []structures.Score{}
	connection.Order("created_at desc").Where("user_id = ?", user_id).Find(&scores)
	return scores
}

func CreateUser(user structures.User) structures.User {
	connection.Create(&user)
	return user
}

func NewScore(score structures.Score) structures.User {
	connection.Create(&score)
	return GetUser(score.User_id.String())
}

func GetBestScores() []structures.BestScores {
	scores := []structures.BestScores{}
	connection.Model(&structures.Score{}).Select("users.username, scores.score, scores.created_at").Joins("left join users on scores.user_id = users.id").Order("scores.score desc"). /*.Limit(10)*/ Find(&scores)
	return scores
}

// R  63461839
// of 63462067
/**
CREATE DATABASE game;

SET DATABASE = game;

CREATE TABLE IF NOT EXISTS users (id UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,username string,created_at TIMESTAMPTZ DEFAULT now());
CREATE TABLE IF NOT EXISTS scores (id UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,user_id UUID NOT NULL REFERENCES users(id),score int NOT NULL,created_at TIMESTAMPTZ DEFAULT now());

INSERT INTO "scores" ("user_id", "score") VALUES ('91b287a0-c8d7-11eb-af62-5c3a459d9f2f' , 890);

*/
