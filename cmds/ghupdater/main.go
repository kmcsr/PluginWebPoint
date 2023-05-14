
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
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/kmcsr/go-logger"
	"github.com/kmcsr/go-logger/logrus"
	"github.com/kmcsr/PluginWebPoint/api"
)

var packageNameRe = regexp.MustCompile(`^[a-zA-Z_.-]+[a-zA-Z0-9_.-]*`)

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
	if DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s?parseTime=true", username, passwd, address, database)); err != nil {
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
	Branch string `json:"branch"`
	RelatedPath string `json:"related_path"`
	Authors []string `json:"authors"`
	Deps DependMap `json:"dependencies"`
	Reqs Requirements `json:"requirements"`
	Desc any `json:"description"`
}

type Asset struct{
	Name string `json:"name"`
	Size int64 `json:"size"`
	DownloadCount int64 `json:"download_count"`
	Url string `json:"url"`
	CreateAt time.Time `json:"created_at"`
	BrowserDownloadUrl string `json:"browser_download_url"`
}

type Release struct{
	Name string `json:"name"`
	TagName string `json:"tag_name"`
	CreateAt time.Time `json:"created_at"`
	Assets []Asset `json:"assets"`
	Description string `json:"description"`
	Prerelease bool `json:"prerelease"`
	ParsedVersion string `json:"parsed_version"`
	Meta any `json:"meta"`
}

const CurrenReleaseSchemaVersion = 7
type PluginRelease struct{
	SchemaVersion int `json:"schema_version"`
	Id string `json:"id"`
	LatestVersion string `json:"latest_version"`
	Releases []Release `json:"releases"`
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
	if meta.SchemaVersion != CurrenReleaseSchemaVersion {
		err = fmt.Errorf("Unexpect schema version %d, expect %d", meta.SchemaVersion, CurrenReleaseSchemaVersion)
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
	const insertCmd = "INSERT INTO plugins (`id`,`name`,`enabled`,`version`,`authors`,`desc`,`desc_zhCN`," +
		"`repo`,`repo_branch`,`repo_subdir`,`link`," +
		"`label_information`,`label_tool`,`label_management`,`label_api`," +
		"`createAt`,`lastRelease`,`github_sync`,`ghRepoOwner`,`ghRepoName`,`last_sync`)" +
		" VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,TRUE,?,?,?)"
	const updateCmd = "UPDATE plugins SET " +
			"`name`=?," +
			"`enabled`=?," +
			"`version`=?," +
			"`authors`=?," +
			"`desc`=?," +
			"`desc_zhCN`=?," +
			"`repo`=?," +
			"`repo_branch`=?," +
			"`repo_subdir`=?," +
			"`link`=?," +
			"`label_information`=?," +
			"`label_tool`=?," +
			"`label_management`=?," +
			"`label_api`=?," +
			"`lastRelease`=?," +
			"`ghRepoOwner`=?," +
			"`ghRepoName`=?," +
			"`last_sync`=?" +
			" WHERE `id`=?"
	const removeDepenceCmd = "DELETE FROM plugin_dependencies WHERE `id`=?"
	const insertDepenceCmd = "INSERT INTO plugin_dependencies (`id`,`target`,`tag`)" +
		" VALUES (?,?,?)"
	const removeRequireCmd = "DELETE FROM plugin_requirements WHERE `id`=?"
	const insertRequireCmd = "INSERT INTO plugin_requirements (`id`,`target`,`tag`)" +
		" VALUES (?,?,?)"
	const insertReleaseCmd = "REPLACE INTO plugin_releases (`id`,`tag`,`enabled`,`stable`,`size`,`uploaded`,`filename`,`downloads`," +
		"`github_url`)" +
		" VALUES (?,?,TRUE,?,?,?,?,?,?)"

	nowt := time.Now()
	now := nowt.Format("2006-01-02 15:04:05")

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
	if !strings.HasPrefix(info.Repo, "https://github.com/") {
		err = fmt.Errorf("Unexpect repo link (missing gh prefix): %q", info.Repo)
		return
	}
	var lastRelease sql.NullTime
	{
		for _, r := range releases.Releases {
			t := r.CreateAt
			if !lastRelease.Valid || t.After(lastRelease.Time) {
				lastRelease.Valid = true
				lastRelease.Time = t
			}
		}
	}
	var ghRepoOwner, ghRepoName string
	{
		b := info.Repo[len("https://github.com/"):]
		if b[len(b) - 1] == '/' {
			b = b[:len(b) - 1]
		}
		paths := strings.Split(b, "/")
		if len(paths) <= 1 {
			err = fmt.Errorf("Unexpect repo link (missing repo name): %q, expect 'https://github.com/{owner}/{name}'", info.Repo)
			return
		}
		if len(paths) > 2 {
			err = fmt.Errorf("Unexpect repo link (extra path): %q, expect 'https://github.com/{owner}/{name}'", info.Repo)
			return
		}
		ghRepoOwner, ghRepoName = paths[0], paths[1]
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
	if err = conn.QueryRowContext(ctx, "SELECT `github_sync`" +
		" FROM plugins WHERE `id`=?", info.Id).Scan(&flag); err != nil && err != sql.ErrNoRows {
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

	if flag.Valid {
		loger.Infof("[%s] Updating metadata", info.Id)
		if _, err = ExecTx(tx, updateCmd, meta.Name, !info.Disable, meta.Version,
			strings.Join(meta.Authors, ","), desc, desc_zhCN,
			info.Repo, info.Branch, info.RelatedPath, link,
			info.Labels.HasInformation(), info.Labels.HasTool(), info.Labels.HasManagement(), info.Labels.HasAPI(),
			lastRelease, ghRepoOwner, ghRepoName, now, info.Id); err != nil {
			return
		}
		if _, err = ExecTx(tx, removeDepenceCmd, info.Id); err != nil {
			return
		}
		if _, err = ExecTx(tx, removeRequireCmd, info.Id); err != nil {
			return
		}
	}else{
		loger.Infof("[%s] Insert into database", info.Id)
		if _, err = ExecTx(tx, insertCmd, info.Id, meta.Name, !info.Disable, meta.Version,
			strings.Join(meta.Authors, ","), desc, desc_zhCN,
			info.Repo, info.Branch, info.RelatedPath, link,
			info.Labels.HasInformation(), info.Labels.HasTool(), info.Labels.HasManagement(), info.Labels.HasAPI(),
			now, lastRelease, ghRepoOwner, ghRepoName, now); err != nil {
			return
		}
	}
	for id, cond := range meta.Deps {
		if _, err = ExecTx(tx, insertDepenceCmd, info.Id, id, cond); err != nil {
			return
		}
	}
	for _, req := range meta.Reqs {
		loger.Debugf("Parsing requirement %q", req)
		id := packageNameRe.FindString(req)
		if len(id) == 0 {
			loger.Warnf("Cannot parse python package requirement: %q", req)
			continue
		}
		cond := strings.ReplaceAll(req[len(id):], " ", "")
		if _, err = ExecTx(tx, insertRequireCmd, info.Id, id, cond); err != nil {
			return
		}
	}
	for _, release := range releases.Releases {
		for _, asset := range release.Assets {
			if strings.HasSuffix(asset.Name, ".mcdr") {
				loger.Debugf("inserting asset: %v", asset)
				if _, err = ExecTx(tx, insertReleaseCmd, info.Id, release.ParsedVersion, release.Prerelease,
					asset.Size, asset.CreateAt, asset.Name, asset.DownloadCount, asset.BrowserDownloadUrl); err != nil {
					// loger.Errorf("Error when insert release into sql")
					return
				}
				break
			}
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
				if err != nil {
					loger.Errorf("[%s] Release json error: %v", info.Id, err)
					return
				}
				if err = updateSql(info, meta, releases); err != nil {
					loger.Errorf("[%s] Cannot sync to database: %v", info.Id, err)
				}
			}(info)
		}
	}
	wg.Wait()
}
