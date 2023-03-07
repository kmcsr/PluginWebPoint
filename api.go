
package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	BASE_DIR string = "/opt/pwebpoint"
	CACHE_DIR string = filepath.Join(BASE_DIR, "caches")
	PLUGIN_CACHE_DIR string = filepath.Join(CACHE_DIR, "plugin")
)

var httpClient = &http.Client{
	Timeout: time.Second * 5,
}

type PluginLabels struct {
	Information bool `json:"information,omitempty"`
	Tool        bool `json:"tool,omitempty"`
	Management  bool `json:"management,omitempty"`
	API         bool `json:"api,omitempty"`
}

type PluginInfo struct {
	Id         string       `json:"id"`
	Name       string       `json:"name"`
	Version    string       `json:"version"`
	Authors    []string     `json:"authors"`
	Desc       string       `json:"desc,omitempty"`
	Desc_zhCN  string       `json:"desc_zhcn,omitempty"`
	CreateAt   time.Time    `json:"createAt"`
	LastUpdate time.Time    `json:"lastUpdate"`
	Repo       string       `json:"repo,omitempty"`
	Link       string       `json:"link,omitempty"`
	Labels     PluginLabels `json:"labels"`
	Downloads  int64        `json:"downloads"`
}

type PluginRelease struct {
	Id        string    `json:"id"`
	Tag       string    `json:"tag"`
	Enabled   bool      `json:"enabled"`
	Stable    bool      `json:"stable"`
	Size      int64     `json:"size"`
	Uploaded  time.Time `json:"uploaded"`
	FileName  string    `json:"filename"`
	Downloads int       `json:"downloads"`
	GithubUrl string    `json:"github_url"`
}

type PluginListOpt struct{
	FilterBy string
	Tags     []string
	SortBy   string
	Reversed bool
}

type API interface {
	GetPluginList(opt PluginListOpt)(infos []*PluginInfo, err error)
	GetPluginInfo(id string)(info *PluginInfo, err error)
	GetPluginReleases(id string)(releases []*PluginRelease, err error)
	GetPluginRelease(id string, tag string)(release *PluginRelease, err error)
	GetPluginReleaseAsset(id string, tag string, filename string)(rc io.ReadSeekCloser, modTime time.Time, err error)
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

	loger.Info("Connecting to db %s:*@%s/%s", username, address, database)

	if api.DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s?parseTime=true", username, passwd, address, database)); err != nil {
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

// TODO split pages
func (api *MySqlAPI)GetPluginList(opt PluginListOpt)(infos []*PluginInfo, err error){
	cmd := "SELECT `id`,`name`,`version`,`authors`,`desc`," +
		"CONVERT_TZ(`createAt`,@@session.time_zone,'+00:00') AS `utc_createAt`," +
		"CONVERT_TZ(`lastUpdate`,@@session.time_zone,'+00:00') AS `utc_lastUpdate`," +
		"`label_information`,`label_tool`,`label_management`,`label_api`" +
		" FROM plugins WHERE `enabled`=TRUE"
	const queryDownloadCmd = "SELECT SUM(`downloads`),`id` FROM plugin_releases GROUP BY `id`"
	args := []any{}
	if len(opt.FilterBy) > 0 {
		cmd0, args0 := parseFilterBy(opt.FilterBy)
		if len(cmd0) > 0 {
			cmd += " AND (" + cmd0 + ")"
			args = append(args, args0...)
		}
	}
	if len(opt.Tags) > 0 {
		cmds := []string{}
		for _, t := range opt.Tags {
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
	switch opt.SortBy {
	case "":
	case "id", "name", "authors", "lastUpdate":
		cmd += " ORDER BY `" + opt.SortBy + "`"
		if opt.Reversed {
			cmd += " DESC"
		}
	default:
		return nil, fmt.Errorf("Unexpect param sortBy=%q", opt.SortBy)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	infomap := make(map[string]*PluginInfo)

	var rows *sql.Rows
	loger.Debugf("exec sql: %q", cmd)
	loger.Debugf("  args: %v", args)
	if rows, err = api.DB.QueryContext(ctx, cmd, args...); err != nil {
		loger.Debugf("sql error: %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			info PluginInfo
			authors string
		)
		if err = rows.Scan(&info.Id, &info.Name, &info.Version, &authors, &info.Desc, &info.CreateAt, &info.LastUpdate,
			&info.Labels.Information, &info.Labels.Tool, &info.Labels.Management, &info.Labels.API); err != nil {
			return
		}
		info.Authors = strings.Split(authors, ",")
		infomap[info.Id] = &info
		infos = append(infos, &info)
	}
	if err = rows.Err(); err != nil {
		return
	}
	rows.Close()
	if rows, err = api.DB.QueryContext(ctx, queryDownloadCmd); err != nil {
		return
	}
	defer rows.Close()
	var (
		downloads sql.NullInt64
		pid string
	)
	for rows.Next() {
		if err = rows.Scan(&downloads, &pid); err != nil {
			return
		}
		if downloads.Valid {
			if info, ok := infomap[pid]; ok {
				info.Downloads = downloads.Int64
			}
		}
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

func (api *MySqlAPI)GetPluginInfo(id string)(info *PluginInfo, err error){
	const queryCmd = "SELECT `name`,`version`,`authors`,`desc`,`desc_zhCN`," +
		"CONVERT_TZ(`createAt`,@@session.time_zone,'+00:00') AS `utc_createAt`," +
		"CONVERT_TZ(`lastUpdate`,@@session.time_zone,'+00:00') AS `utc_lastUpdate`," +
		"`repo`,`link`,`label_information`,`label_tool`,`label_management`,`label_api`" +
		" FROM plugins WHERE `id`=? AND `enabled`=TRUE"
	const queryDownloadCmd = "SELECT SUM(`downloads`) FROM plugin_releases WHERE `id`=?"
	var (
		authors string
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	info = new(PluginInfo)
	if err = api.DB.QueryRowContext(ctx, queryCmd, id).
		Scan(&info.Name, &info.Version, &authors, &info.Desc, &info.Desc_zhCN, &info.CreateAt, &info.LastUpdate,
		&info.Repo, &info.Link,
		&info.Labels.Information, &info.Labels.Tool, &info.Labels.Management, &info.Labels.API); err != nil {
		return
	}
	var downloads sql.NullInt64
	if err = api.DB.QueryRowContext(ctx, queryDownloadCmd, id).Scan(&downloads); err != nil && err != sql.ErrNoRows {
		return
	}
	info.Id = id
	if downloads.Valid {
		info.Downloads = downloads.Int64
	}
	info.Authors = strings.Split(authors, ",")
	err = nil
	return
}

func (api *MySqlAPI)GetPluginReleases(id string)(releases []*PluginRelease, err error){
	const queryCmd = "SELECT `tag`,`enabled`,`stable`,`size`," +
		"CONVERT_TZ(`uploaded`,@@session.time_zone,'+00:00') AS `utc_uploaded`," +
		"`filename`,`downloads`,`github_url`" +
		" FROM plugin_releases WHERE `id`=?"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	var rows *sql.Rows
	if rows, err = api.DB.QueryContext(ctx, queryCmd, id); err != nil {
		loger.Debugf("sql error: %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			release PluginRelease
			ghUrl sql.NullString
		)
		if err = rows.Scan(&release.Tag, &release.Enabled, &release.Stable, &release.Size,
			&release.Uploaded, &release.FileName, &release.Downloads, &ghUrl); err != nil {
			return
		}
		release.Id = id
		if ghUrl.Valid {
			release.GithubUrl = ghUrl.String
		}
		releases = append(releases, &release)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

func (api *MySqlAPI)GetPluginRelease(id string, tag string)(release *PluginRelease, err error){
	const queryCmd = "SELECT `enabled`,`stable`,`size`," +
		"CONVERT_TZ(`uploaded`,@@session.time_zone,'+00:00') AS `utc_uploaded`," +
		"`filename`,`downloads`,`github_url`" +
		" FROM plugin_releases WHERE `id`=? AND `tag`=?"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	release = new(PluginRelease)
	var (
		downloads sql.NullInt64
		ghUrl sql.NullString
	)
	if err = api.DB.QueryRowContext(ctx, queryCmd, id, tag).
		Scan(&release.Enabled, &release.Stable, &release.Size, &release.Uploaded,
			&release.FileName, &downloads, &ghUrl); err != nil {
		return
	}
	release.Id = id
	release.Tag = tag
	if downloads.Valid {
		release.Downloads = (int)(downloads.Int64)
	}
	if ghUrl.Valid {
		release.GithubUrl = ghUrl.String
	}
	return
}

type StatusCodeErr struct{
	Code int
}

func (e *StatusCodeErr)Error()(string){
	return fmt.Sprintf("Unexpect http status code %d (%s)", e.Code, http.StatusText(e.Code))
}

type NopReadSeeker struct{
	io.ReadSeeker
}

func (NopReadSeeker)Close()(error){ return nil }

func (api *MySqlAPI)GetPluginReleaseAsset(id string, tag string, filename string)(rc io.ReadSeekCloser, modTime time.Time, err error){
	cache := filepath.Join(PLUGIN_CACHE_DIR, filepath.Clean(filepath.Join(id, tag, filename)))
	var fd *os.File
	if fd, err = os.Open(cache); err == nil {
		if stat, err := fd.Stat(); err == nil {
			modTime = stat.ModTime()
		}
		rc = fd
		return
	}
	var release *PluginRelease
	if release, err = api.GetPluginRelease(id, tag); err != nil {
		return
	}
	var resp *http.Response
	loger.Debugf("Downloading %q", release.GithubUrl)
	if resp, err = httpClient.Get(release.GithubUrl); err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = &StatusCodeErr{resp.StatusCode}
		return
	}
	var data []byte
	if data, err = io.ReadAll(resp.Body); err != nil {
		return
	}
	if err := os.MkdirAll(filepath.Dir(cache), 0755); err == nil {
		if err = os.WriteFile(cache, data, 0444); err != nil {
			loger.Warnf("Cannot write cache file %q: %v", cache, err)
		}else{
			loger.Infof("Cache %s(v%s):%s at %q: %v", id, tag, filename, cache, err)
		}
	}else{
		loger.Warnf("Cannot make cache dir %q: %v", filepath.Dir(cache), err)
	}
	rc = NopReadSeeker{bytes.NewReader(data)}
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
