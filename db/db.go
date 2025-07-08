package db

import (
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/duanchi/min/v2/context"
	"github.com/duanchi/min/v2/db/xorm"
	config2 "github.com/duanchi/min/v2/types/config"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"xorm.io/core"

	// _ "github.com/mattn/go-sqlite3"
	"log"
	"net/url"
	"strings"
)

var Connection *xorm.Engine
var Connections map[string]*xorm.Engine
var dbConfig config2.Db

func Init() {
	var err error

	dbConfig = context.GetApplicationContext().GetConfig("Db").(config2.Db)

	if len(dbConfig.Sources) > 0 {
		Connections = map[string]*xorm.Engine{}
		for name, sourceConfig := range dbConfig.Sources {
			parsedDsn, _ := url.Parse(sourceConfig.Dsn)
			Connections[name], err = connect(parsedDsn, sourceConfig)
			fmt.Println("Data Source [" + name + "] Inited!")
			if err != nil {
				fmt.Println(err.Error())
				break
			}
			if name == "default" {
				Connection = Connections[name]
			}
		}
	} else {
		parsedDsn, _ := url.Parse(context.GetApplicationContext().GetConfig("Db.Dsn").(string))
		Connection, err = connect(parsedDsn, config2.DbConfig{
			Dsn:        context.GetApplicationContext().GetConfig("Db.Dsn").(string),
			MigrateSQL: context.GetApplicationContext().GetConfig("Db.MigrateSQL").(string),
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func Engine(name string) *xorm.Engine {
	if _, has := Connections[name]; has {
		return Connections[name]
	} else {
		if dbConfig.SelectEngine == nil {
			return nil
		} else {
			dsn, err := dbConfig.SelectEngine(name)
			if err != nil {
				panic(err)
			}

			err = NewEngine(name, config2.DbConfig{
				Dsn: dsn,
			})
			if err != nil {
				panic(err)
			}
			return Connections[name]
		}
	}
}

func Remove(name string) {
	if _, has := Connections[name]; has {
		Engine(name).Close()
		delete(Connections, name)
	}
}

func NewEngine(name string, sourceConfig config2.DbConfig) (err error) {
	parsedDsn, _ := url.Parse(sourceConfig.Dsn)
	Connections[name], err = connect(parsedDsn, sourceConfig)

	return err
}

func connect(dsnUrl *url.URL, dbConfig config2.DbConfig) (connection *xorm.Engine, err error) {

	defer func() {
		e := recover()
		if e != nil {
			fmt.Printf("%s", e)
			log.Fatal(err)
		}
		return
	}()

	switch dsnUrl.Scheme {
	case "postgres":
		password, _ := dsnUrl.User.Password()
		dbStack := strings.Split(strings.Trim(dsnUrl.Path, "/"), "/")
		prefix := dsnUrl.Query().Get("prefix")
		dbname := ""
		schema := ""
		if len(dbStack) > 1 {
			dbname = dbStack[0]
			schema = dbStack[1]
		} else {
			dbname = dbStack[0]
		}

		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			dsnUrl.Hostname(),
			dsnUrl.Port(),
			dsnUrl.User.Username(),
			password,
			dbname,
			dsnUrl.Query().Get("sslmode"),
		)

		connection, err = xorm.NewEngine("postgres", dsn)
		if err != nil {
			panic(fmt.Sprintf("Database Init Error %s", dsn))
		}

		if schema != "" {
			connection.SetSchema(schema)
		}

		if prefix != "" {
			connection.SetTableMapper(core.NewPrefixMapper(core.SnakeMapper{}, prefix))
		}

		err = connection.Ping()
		if err != nil {
			log.Fatal(err)
			return
		}

	case "mysql":

		host := dsnUrl.Host
		prefix := dsnUrl.Query().Get("prefix")
		dsnUrl.Query().Del("prefix")
		query := dsnUrl.Query().Encode()

		if query == "" {
			query = ""
		} else {
			query = "?" + query
		}

		if host[0:1] == "/" {
			host = "unix(" + host + ")"
		} else {
			host = "tcp(" + host + ")"
		}

		dsn := dsnUrl.User.String() + "@" + host + dsnUrl.Path + query

		connection, err = xorm.NewEngine("mysql", dsn)
		if err != nil {
			panic(fmt.Sprintf("Database Init Error %s", dsn))
		}

		if prefix != "" {
			connection.SetTableMapper(core.NewPrefixMapper(core.SnakeMapper{}, prefix))
		}

		err = connection.Ping()
		if err != nil {
			log.Fatal(err)
			return
		}

	case "sqlserver":
		fallthrough
	case "mssql":
		connection, err = xorm.NewEngine("sqlserver", dsnUrl.String())

		err = connection.Ping()
		if err != nil {
			log.Fatal(err)
			return
		}
		/*case "sqlite":
		dbFile := dbConfig.Dsn[9:]
		isNewFile := false
		connection, err = xorm.NewEngine("sqlite3", dbFile)
		if _, fileError := os.Stat(dbFile); fileError != nil {
			f, createErr := os.Create(dbFile)
			isNewFile = true
			defer f.Close()
			if createErr != nil {
				fmt.Println("Create DB file Error in " + dbFile)
				return
			}
		}
		err = connection.Ping()
		if err != nil {
			log.Fatal(err)
			return
		}
		if isNewFile && dbConfig.MigrateSQL != "" {
			sql, readErr := ioutil.ReadFile(dbConfig.MigrateSQL)
			if readErr != nil {
				log.Fatal(readErr)
				return
			}
			_, importErr := connection.Import(bytes.NewReader(sql))

			if importErr == nil {
				fmt.Println("Import DB successful!")
			} else {
				fmt.Println("Import DB error, " + importErr.Error())
			}

			return
		}*/
	}

	if err == nil {
		fmt.Println("connect database success!")
	}
	if context.GetApplicationContext().GetConfig("Env").(string) == "development" {
		connection.ShowSQL()
	}

	return
}
