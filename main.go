// main.go
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	sq "github.com/Masterminds/squirrel"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	db *sql.DB
)

func init() {
	log.WithField("status", "starting").Debug("initialize")
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	viper.SetEnvPrefix("sample")
	viper.AutomaticEnv()
	getSecret()
	db = connectRDS()

	log.WithField("status", "success").Debug("initialize")

}
func connectRDS() (db *sql.DB) {
	log.WithField("status", "starting").Info("connectRDS")

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		viper.GetString("DB-USER"),
		viper.GetString("DB-PASSWORD"),
		viper.GetString("DB-HOST"),
		viper.GetString("DB-DEFAULT"),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.WithField("status", "success").Info("connectRDS")
	return
}

func say(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	lp := request.QueryStringParameters["lp"]
	value, err := readLpSecure(lp)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: http.StatusInternalServerError,
		}, err
	}
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Leadership Principles number %s is %s", lp, value),
		StatusCode: http.StatusOK,
	}, nil
}

func readLp(input string) (output string, err error) {
	log.WithField("status", "starting").Info("readLp")

	query := fmt.Sprintf("SELECT value from lp where id='%s'", input)
	err = db.QueryRow(query).Scan(&output)

	log.WithField("status", "success").Info("readLp")
	return
}

func readLpSecure(input string) (output string, err error) {
	log.WithField("status", "starting").Info("readLpSecure")
	sql, args, err := sq.Select("value").From("lp").Where(sq.Eq{"id": input}).ToSql()
	log.Debug(sql)
	log.Debug(args)
	err = db.QueryRow(sql, args).Scan(&output)

	log.WithField("status", "success").Info("readLpSecure")
	return
}

func getSecret() {
	log.WithField("status", "starting").Debug("getSecret")
	type dbCredential struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		Host      string `json:"host"`
		Port      int    `json:"port"`
		DBDefault string `json:"dbname"`
	}

	secretDBName := viper.GetString("secret_manager_db")
	region := viper.GetString("region")
	svc := secretsmanager.New(session.New(),
		aws.NewConfig().WithRegion(region))

	//get db
	inputDb := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretDBName),
	}

	resultDb, err := svc.GetSecretValue(inputDb)
	if err != nil {
		log.Fatalln(err)
	}

	db := dbCredential{}
	err = json.Unmarshal([]byte(*resultDb.SecretString), &db)
	if err != nil {
		log.Fatalln(err)
	}
	//

	viper.Set("DB-USER", db.Username)
	viper.Set("DB-PASSWORD", db.Password)
	viper.Set("DB-HOST", fmt.Sprintf("%s:%d", db.Host, db.Port))
	viper.Set("DB-DEFAULT", db.DBDefault)
	log.WithField("status", "success").Info("getSecret")
}

func main() {
	lambda.Start(say)
}
