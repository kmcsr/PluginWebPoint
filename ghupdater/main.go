
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kmcsr/go-logger"
	"github.com/kmcsr/go-logger/logrus"
)

var loger logger.Logger = logrus.Logger

var DB *sql.DB = initDB()

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
	DB.SetMaxOpenConns(16)
	DB.SetMaxIdleConns(16)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 3)
	defer cancel()
	if err = DB.PingContext(ctx); err != nil {
		loger.Fatalf("Cannot ping to database: %v", err)
	}
	return
}

var (
	target string = "https://github.com/MCDReforged/PluginCatalogue"
	targetRaw string = "raw.githubusercontent.com/MCDReforged/PluginCatalogue/" // meta/{{plugin_id}}/meta.json
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

type Dependencies map[string]string
type Requirements []string

type PluginMeta struct{
	Id string `json:"id"`
	Name string `json:"name"`
	Version string `json:"version"`
	Repo string `json:"repository"`
	Authors []string `json:"authors"`
	Deps Dependencies `json:"Dependencies"`
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
	p := "https://" + path.Join(targetRaw, "meta", id, "meta.json")
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
	p := "https://" + path.Join(targetRaw, "meta", id, "release.json")
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

func updateSql(info PluginInfo, meta PluginMeta, releases PluginRelease)(err error){
	const updateCmd = "UPDATE plugins SET " +
			"`name`=?," +
			"`enabled`=?," +
			"`version`=?," +
			"`authors`=?," +
			"`desc`=?," +
			"`desc_zhCN`=?," +
			"`repo`=?," +
			"`link`=?," +
			"`label_information`=?," +
			"`label_tool`=?," +
			"`label_management`=?," +
			"`label_api`=?," +
			"`last_sync`=? " +
		"WHERE `id`=? AND `github_sync`=TRUE"
	const insertCmd = "INSERT INTO plugins (`id`,`name`,`enabled`,`version`,`authors`,`desc`,`desc_zhCN`,`repo`,`link`," +
			"`label_information`,`label_tool`,`label_management`,`label_api`,`github_sync`,`last_sync`)" +
			"VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,TRUE,?)"
	const insertReleaseCmd = "INSERT INTO plugin_releases (`id`,`tag`,`enabled`,`stable`,`size`,`filename`,`downloads`," +
			"`github_url`)" +
			"VALUES (?,?,TRUE,?,?,?,?,?)"

	now := time.Now()

	tx, err := DB.Begin()
	if err != nil {
		loger.Errorf("Error when new Tx")
		return
	}
	defer tx.Rollback()
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
	link := path.Join(info.Repo, "tree", info.Branch, info.RelatedPath)
	var res sql.Result
	loger.Infof("Updating database for %s", info.Id)
	if res, err = tx.Exec(updateCmd, meta.Name, !info.Disable, meta.Version,
		strings.Join(meta.Authors, ","), desc, desc_zhCN, info.Repo, link,
		info.Labels.HasInformation(), info.Labels.HasTool(), info.Labels.HasManagement(), info.Labels.HasAPI(),
		now, info.Id); err != nil {
		loger.Errorf("Error when update meta")
		return
	}
	var n int64
	if n, err = res.RowsAffected(); err != nil {
		return
	}
	if n == 0 {
		loger.Infof("Insert into database for %s", info.Id)
		if res, err = tx.Exec(insertCmd, info.Id, meta.Name, !info.Disable, meta.Version,
			strings.Join(meta.Authors, ","), desc, desc_zhCN, info.Repo, link,
			info.Labels.HasInformation(), info.Labels.HasTool(), info.Labels.HasManagement(), info.Labels.HasAPI(),
			now); err != nil {
			loger.Errorf("Error when insert meta into sql")
			return
		}
	}
	for _, release := range releases.Releases {
		assets := release.Assets[0]
		if _, err = tx.Exec(insertReleaseCmd, info.Id, release.ParsedVersion, release.Prerelease,
			assets.Size, assets.Name, assets.DownloadCount, assets.BrowserDownloadUrl); err != nil {
			loger.Errorf("Error when insert release into sql")
			return
		}
	}
	if err = tx.Commit(); err != nil {
		loger.Errorf("Error when commit Tx")
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
					loger.Errorf("Cannot get meta json: %v", err)
					return
				}
				releases, err := GetPluginReleaseJson(info.Id)
				if err = updateSql(info, meta, releases); err != nil {
					loger.Errorf("Cannot update to database: %v", err)
				}
			}(info)
		}
	}
	wg.Wait()
}
