
package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type PluginLabels struct {
	Information bool `json:"information"`
	Tool        bool `json:"tool"`
	Management  bool `json:"management"`
	API         bool `json:"api"`
}

type PluginMeta struct {
	Id      string       `json:"id"`
	Name    string       `json:"name"`
	Version string       `json:"version"`
	Authors []string     `json:"authors"`
	Desc    string       `json:"desc"`
	LastUpdate time.Time `json:"lastUpdate"`
	Labels  PluginLabels `json:"labels"`
}

type API interface {
	GetPluginList(filterBy string, tags []string, sortBy string, reversed bool)(metas []*PluginMeta, err error)
}

var APIIns API = NewMySqlAPI()

type MySqlAPI struct {
	DB *sql.DB
}

func NewMySqlAPI()(api *MySqlAPI){
	var err error
	api = new(MySqlAPI)

	username := os.Getenv("DB_USER")
	passwd := os.Getenv("DB_PASSWD")
	address := os.Getenv("DB_ADDR")
	database := os.Getenv("DB_NAME")

	loger.Debug("Connecting to db %s@%s/%s", username, address, database)

	if api.DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", username, passwd, address, database)); err != nil {
		loger.Fatalf("Cannot connect to database: %v", err)
	}
	api.DB.SetConnMaxLifetime(time.Minute * 3)
	api.DB.SetMaxOpenConns(10)
	api.DB.SetMaxIdleConns(10)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 3)
	defer cancel()
	if err = api.DB.PingContext(ctx); err != nil {
		loger.Fatalf("Cannot ping to database: %v", err)
	}
	return
}

func (api *MySqlAPI)GetPluginList(filterBy string, tags []string, sortBy string, reversed bool)(metas []*PluginMeta, err error){
	cmd := "SELECT `id`,`name`,`version`,`authors`,`desc`,`lastUpdate`," + 
		"`label_information`,`label_tool`,`label_management`,`label_api` FROM plugins WHERE `enabled`=TRUE"
	args := []any{}
	if len(filterBy) > 0 {
		cmd0, args0 := parseFilterBy(filterBy)
		if len(cmd0) > 0 {
			cmd += " AND (" + cmd0 + ")"
			args = append(args, args0...)
		}
	}
	if len(tags) > 0 {
		cmds := []string{}
		for _, t := range tags {
			switch t {
			case "management", "tool", "information", "api":
				cmds = append(cmds, "`label_" + t + "`=TRUE")
			default:
				return nil, fmt.Errorf("Unexpect param tags=%q", t)
			}
		}
		if len(cmds) > 0 {
			cmd += " AND (" + strings.Join(cmds, " OR ") + ")"
		}
	}
	switch sortBy {
	case "":
	case "id", "name", "authors", "lastUpdate":
		cmd += " ORDER BY `" + sortBy + "`"
		if reversed {
			cmd += " DESC"
		}
	default:
		return nil, fmt.Errorf("Unexpect param sortBy=%q", sortBy)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	var rows *sql.Rows
	loger.Debugf("exec sql: %q", cmd)
	loger.Debugf("  args: %v", args)
	if rows, err = api.DB.QueryContext(ctx, cmd, args...); err != nil {
		loger.Debugf("sql error:", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			meta PluginMeta
			authors string
			lastUpdate sql.NullTime
		)
		if err = rows.Scan(&meta.Id, &meta.Name, &meta.Version, &authors, &meta.Desc, &lastUpdate,
			&meta.Labels.Information, &meta.Labels.Tool, &meta.Labels.Management, &meta.Labels.API); err != nil {
			return
		}
		if lastUpdate.Valid {
			meta.LastUpdate = lastUpdate.Time
		}
		meta.Authors = strings.Split(authors, ",")
		metas = append(metas, &meta)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

func parseFilterBy(filter string)(cmd string, args []any){
	cmds := []string{}
	for _, a := range strings.Split(filter, " ") {
		ok := true
		switch {
		case a == "":
		case strings.HasPrefix(a, "id:"):
			a = a[len("id:"):]
			cmds = append(cmds, "`id` LIKE ?")
			f := "%" + a + "%"
			args = append(args, f)
		case strings.HasPrefix(a, "name:"):
			a = a[len("name:"):]
			cmds = append(cmds, "`name` LIKE ?")
			f := "%" + a + "%"
			args = append(args, f)
		case strings.HasPrefix(a, "@"):
			a = a[len("@"):]
			ok = false
			fallthrough
		case strings.HasPrefix(a, "author:"):
			if ok {
				a = a[len("author:"):]
				ok = false
			}
			fallthrough
		case strings.HasPrefix(a, "authors:"):
			if ok {
				a = a[len("authors:"):]
			}
			for _, a := range strings.Split(a, ",") {
				cmds = append(cmds, "`authors` LIKE ?")
				f := "%" + a + "%"
				args = append(args, f)
			}
		case strings.HasPrefix(a, "description:"):
			a = a[len("description:"):]
			ok = false
			fallthrough
		case strings.HasPrefix(a, "desc:"):
			if ok {
				a = a[len("desc:"):]
			}
			cmds = append(cmds, "`desc` LIKE ?")
			f := "%" + a + "%"
			args = append(args, f)
		default:
			cmds = append(cmds, "`id` LIKE ? OR `name` LIKE ? OR `authors` LIKE ? OR `desc` LIKE ?")
			f := "%" + a + "%"
			args = append(args, f, f, f, f)
		}
	}
	cmd = strings.Join(cmds, " OR ")
	return
}
