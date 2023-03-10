
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/kmcsr/go-logger"
	"github.com/kmcsr/go-logger/logrus"
	"github.com/kmcsr/PluginWebPoint/api"
)

var loger logger.Logger = initLogger()

func initLogger()(loger logger.Logger){
	loger = logrus.Logger
	if os.Getenv("DEBUG") == "true" {
		loger.SetLevel(logger.TraceLevel)
	}
	return
}

var DB *sql.DB = initDB()

const (
	maxConn = 16
)

func initDB()(DB *sql.DB){
	username := os.Getenv("DB_USER")
	passwd := os.Getenv("DB_PASSWD")
	address := os.Getenv("DB_ADDR")
	database := os.Getenv("DB_NAME")

	loger.Debug("Connecting to db %s@%s/%s", username, address, database)

	var err error
	if DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s", username, passwd, address, database)); err != nil {
		loger.Fatalf("Cannot connect to database: %v", err)
	}
	DB.SetConnMaxLifetime(time.Minute * 3)
	DB.SetMaxOpenConns(maxConn * 2 + 3)
	DB.SetMaxIdleConns(maxConn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 3)
	defer cancel()
	if err = DB.PingContext(ctx); err != nil {
		loger.Fatalf("Cannot ping to database: %v", err)
	}
	return
}

func getDBConn(ctx context.Context)(conn *sql.Conn, err error){
	conn, err = DB.Conn(ctx)
	if err != nil {
		loger.Errorf("Cannot get new conn: %T, %#v", err, err)
		return
	}
	return
}

var (
	target string = "https://github.com/MCDReforged/PluginCatalogue"
	targetRaw string = "https://raw.githubusercontent.com/MCDReforged/PluginCatalogue/" // meta/{{plugin_id}}/meta.json
)

type Author struct{
	Name string `json:"name"`
	Link string `json:"link"`
}

type Labels []string

func (l Labels)HasInformation()(bool){
	for _, a := range l {
		if a == "information" {
			return true
		}
	}
	return false
}

func (l Labels)HasTool()(bool){
	for _, a := range l {
		if a == "tool" {
			return true
		}
	}
	return false
}

func (l Labels)HasManagement()(bool){
	for _, a := range l {
		if a == "management" {
			return true
		}
	}
	return false
}

func (l Labels)HasAPI()(bool){
	for _, a := range l {
		if a == "api" {
			return true
		}
	}
	return false
}

type PluginInfo struct{
	Disable bool `json:"disable"`
	Id string `json:"id"`
	Authors []Author `json:"authors"`
	Repo string `json:"repository"`
	Branch string `json:"branch"`
	RelatedPath string `json:"related_path"`
	Labels Labels `json:"labels"`
}

type DependMap map[string]api.VersionCond
type Requirements []string

type PluginMeta struct{
	Id string `json:"id"`
	Name string `json:"name"`
	Version string `json:"version"`
	Repo string `json:"repository"`
	Authors []string `json:"authors"`
	Deps DependMap `json:"dependencies"`
	Reqs Requirements `json:"requirements"`
	Desc any `json:"description"`
}

type Asset struct{
	Name string `json:"name"`
	Size int64 `json:"size"`
	DownloadCount int64 `json:"download_count"`
	CreateAt time.Time `json:"create_at"`
	BrowserDownloadUrl string `json:"browser_download_url"`
}

type Release struct{
	Url string `json:"url"`
	Name string `json:"name"`
	TagName string `json:"tag_name"`
	CreateAt time.Time `json:"create_at"`
	Assets []Asset `json:"assets"`
	Description string `json:"description"`
	Prerelease bool `json:"prerelease"`
	ParsedVersion string `json:"parsed_version"`
}

type PluginRelease struct{
	SchemaVersion int `json:"schema_version"`
	Id string `json:"id"`
	LatestVersion string `json:"latest_version"`
	Releases []Release `json:"releases"`
	ReleaseMeta map[string]PluginMeta `json:"release_meta"`
}

func GetPluginMetaJson(id string)(meta PluginMeta, err error){
	p, err := url.JoinPath(targetRaw, "meta", id, "meta.json")
	if err != nil {
		return
	}
	loger.Infof("Getting %q", p)
	resp, err := http.DefaultClient.Get(p)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if err = json.Unmarshal(body, &meta); err != nil {
		return
	}
	return
}

func GetPluginReleaseJson(id string)(meta PluginRelease, err error){
	p, err := url.JoinPath(targetRaw, "meta", id, "release.json")
	if err != nil {
		return
	}
	loger.Infof("Getting %q", p)
	resp, err := http.DefaultClient.Get(p)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if err = json.Unmarshal(body, &meta); err != nil {
		return
	}
	if meta.SchemaVersion != 4 {
		err = fmt.Errorf("Unexpect schema version %d, expect 4", meta.SchemaVersion)
		return
	}
	return
}

func ExecTx(tx *sql.Tx, cmd string, args ...any)(res sql.Result, err error){
	loger.Debugf("Exec sql cmd: %s\n  args: %v", cmd, args)
	for {
		if res, err = tx.Exec(cmd, args...); err != nil {
			if e, ok := err.(*mysql.MySQLError); ok {
				switch e.Number {
				case 1213:
					continue
				}
			}
		}
		return
	}
}

func updateSql(info PluginInfo, meta PluginMeta, releases PluginRelease)(err error){
	const insertCmd = "REPLACE INTO plugins (`id`,`name`,`enabled`,`version`,`authors`,`desc`,`desc_zhCN`,`repo`,`link`," +
		"`label_information`,`label_tool`,`label_management`,`label_api`,`github_sync`,`last_sync`)" +
		"VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,TRUE,?)"
	const removeDepenceCmd = "DELETE FROM plugin_dependencies WHERE `id`=?"
	const insertDepenceCmd = "INSERT INTO plugin_dependencies (`id`,`target`,`tag`)" +
		"VALUES (?,?,?)"
	const insertReleaseCmd = "INSERT IGNORE INTO plugin_releases (`id`,`tag`,`enabled`,`stable`,`size`,`filename`,`downloads`," +
		"`github_url`)" +
		"VALUES (?,?,TRUE,?,?,?,?,?)"

	now := time.Now().Format("2006-01-02 15:04:05")

	sort.Strings(meta.Authors)
	var (
		desc string
		desc_zhCN string
	)
	if d, ok := meta.Desc.(string); ok {
		desc = d
	}else if m, ok := meta.Desc.(map[string]any); ok {
		if d, ok := m["en_us"].(string); ok {
			desc = d
		}
		if d, ok := m["zh_cn"].(string); ok {
			desc_zhCN = d
		}
	}
	link, err := url.JoinPath(info.Repo, "tree", info.Branch, info.RelatedPath)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	conn, err := getDBConn(ctx)
	if err != nil {
		return
	}
	defer conn.Close()

	var flag sql.NullBool
	if err = conn.QueryRowContext(ctx,
		"SELECT 1 FROM plugins WHERE `id`=? AND `github_sync`=TRUE LIMIT 1", info.Id).Scan(&flag); err != nil {
		return
	}
	if flag.Valid && !flag.Bool {
		loger.Debugf("Plugin %s is not synced from github", info.Id)
		return
	}

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	defer tx.Rollback()

	loger.Infof("Insert into database for %s", info.Id)
	if _, err = ExecTx(tx, insertCmd, info.Id, meta.Name, !info.Disable, meta.Version,
		strings.Join(meta.Authors, ","), desc, desc_zhCN, info.Repo, link,
		info.Labels.HasInformation(), info.Labels.HasTool(), info.Labels.HasManagement(), info.Labels.HasAPI(),
		now); err != nil {
		return
	}
	if _, err = ExecTx(tx, removeDepenceCmd, info.Id); err != nil {
		return
	}
	for id, cond := range meta.Deps {
		if _, err = ExecTx(tx, insertDepenceCmd, info.Id, id, cond); err != nil {
			return
		}
	}
	for _, release := range releases.Releases {
		assets := release.Assets[0]
		if _, err = ExecTx(tx, insertReleaseCmd, info.Id, release.ParsedVersion, release.Prerelease,
			assets.Size, assets.Name, assets.DownloadCount, assets.BrowserDownloadUrl); err != nil {
			if e, ok := err.(*mysql.MySQLError); ok {
				if e.Number == 1062 {
					continue
				}
			}
			// loger.Errorf("Error when insert release into sql")
			return
		}
	}
	if err = tx.Commit(); err != nil {
		return
	}
	return
}

func main(){
	dir, err := os.MkdirTemp("", "gh_sync")
	if err != nil {
		loger.Fatalf("Cannot make temp dir: %v", err)
	}
	defer os.RemoveAll(dir)
	loger.Infof("Temp dir: %s", dir)

	cmd := exec.Command("git", "-C", dir, "clone", target)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		loger.Fatalf("Cannot execute git clone: %v", err)
	}
	pluginsPath := filepath.Join(dir, path.Base(target), "plugins")
	plugins, err := os.ReadDir(pluginsPath)
	if err != nil {
		loger.Panic(err)
	}

	var wg sync.WaitGroup
	for _, p := range plugins {
		if p.IsDir() {
			infop := filepath.Join(pluginsPath, p.Name(), "plugin_info.json")
			infob, err := os.ReadFile(infop)
			if err != nil {
				loger.Errorf("Cannot read %q: %v", infop, err)
				continue
			}
			var info PluginInfo
			if err = json.Unmarshal(infob, &info); err != nil {
				loger.Errorf("Cannot parse %q: %v", infop, err)
				continue
			}
			if info.Disable {
				loger.Warnf("disabled plugin %s", info.Id)
				continue
			}
			wg.Add(1)
			go func(info PluginInfo){
				defer wg.Done()
				meta, err := GetPluginMetaJson(info.Id)
				if err != nil {
					loger.Errorf("[%s] Cannot get meta json: %v", info.Id, err)
					return
				}
				releases, err := GetPluginReleaseJson(info.Id)
				if err = updateSql(info, meta, releases); err != nil {
					loger.Errorf("[%s] Cannot update to database: %v", info.Id, err)
				}
			}(info)
		}
	}
	wg.Wait()
}
