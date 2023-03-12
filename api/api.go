
package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

var (
	BASE_DIR string = "/opt/pwebpoint"
	PLUGIN_DIR string = filepath.Join(BASE_DIR, "plugin")
	CACHE_DIR string = filepath.Join(BASE_DIR, "caches")
	PLUGIN_CACHE_DIR string = filepath.Join(CACHE_DIR, "plugin")
)

var (
	ErrNotFound = errors.New("ErrNotFound")
)

type PluginCounts struct {
	Total       int `json:"total"`
	Information int `json:"information"`
	Tool        int `json:"tool"`
	Management  int `json:"management"`
	Api         int `json:"api"`
}

type PluginLabels struct {
	Information bool `json:"information,omitempty"`
	Tool        bool `json:"tool,omitempty"`
	Management  bool `json:"management,omitempty"`
	Api         bool `json:"api,omitempty"`
}

type DependMap map[string]VersionCond

type PluginInfo struct {
	Id           string       `json:"id"`
	Name         string       `json:"name"`
	Version      Version      `json:"version"`
	Authors      []string     `json:"authors"`
	Desc         string       `json:"desc,omitempty"`
	Desc_zhCN    string       `json:"desc_zhCN,omitempty"`
	CreateAt     time.Time    `json:"createAt"`
	LastUpdate   time.Time    `json:"lastUpdate"`
	Repo         string       `json:"repo,omitempty"`
	RepoBranch   string       `json:"repoBranch,omitempty"`
	RepoSubdir   string       `json:"repoSubdir,omitempty"`
	Link         string       `json:"link,omitempty"`
	Labels       PluginLabels `json:"labels"`
	Downloads    int64        `json:"downloads"`
	Dependencies DependMap    `json:"dependencies,omitempty"`
	GithubSync   bool         `json:"github_sync"`
	GhRepoOwner  string       `json:"ghRepoOwner,omitempty"`
	GhRepoName   string       `json:"ghRepoName,omitempty"`
	LastSync     *time.Time   `json:"last_sync,omitempty"`
}

type PluginRelease struct {
	Id        string    `json:"id"`
	Tag       Version   `json:"tag"`
	Enabled   bool      `json:"enabled"`
	Stable    bool      `json:"stable"`
	Size      int64     `json:"size"`
	Uploaded  time.Time `json:"uploaded"`
	FileName  string    `json:"filename"`
	Downloads int       `json:"downloads"`
	GithubUrl string    `json:"github_url"`
}

type PluginListOpt struct{
	FilterBy string   `json:"filterBy,omitempty"`
	Tags     []string `json:"tags,omitempty"`
	SortBy   string   `json:"sortBy,omitempty"`
	Reversed bool     `json:"reversed,omitempty"`
	Limit    int      `json:"limit,omitempty"`
	Offset   int      `json:"offset,omitempty"`
}

type API interface {
	GetPluginCounts(opt PluginListOpt)(count PluginCounts, err error)
	GetPluginList(opt PluginListOpt)(infos []*PluginInfo, err error)
	GetPluginInfo(id string)(info *PluginInfo, err error)
	GetPluginReadme(id string)(data []byte, prefix string, err error)
	GetPluginReleases(id string)(releases []*PluginRelease, err error)
	GetPluginRelease(id string, tag Version)(release *PluginRelease, err error)
	GetPluginReleaseAsset(id string, tag Version, filename string)(rc io.ReadSeekCloser, modTime time.Time, err error)
}

var Ins API = nil

type MySqlAPI struct {
	DB *sql.DB
}

var _ API = (*MySqlAPI)(nil)

func NewMySqlAPI(username string, passwd string, address string, database string)(api *MySqlAPI){
	var err error
	api = new(MySqlAPI)

	loger.Infof("Connecting to db %s:*@%s/%s", username, address, database)

	if api.DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s?parseTime=true", username, passwd, address, database)); err != nil {
		loger.Fatalf("Cannot connect to database: %v", err)
	}
	api.DB.SetConnMaxLifetime(time.Minute * 3)
	api.DB.SetMaxOpenConns(256)
	api.DB.SetMaxIdleConns(25)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 3)
	defer cancel()
	if err = api.DB.PingContext(ctx); err != nil {
		loger.Fatalf("Cannot ping to database: %v", err)
	}
	return
}

func (api *MySqlAPI)QueryContext(ctx context.Context, cmd string, args ...any)(rows *sql.Rows, err error){
	// loger.Debugf("Query sql cmd: %s\n  args: %v", cmd, args)
	for {
		if rows, err = api.DB.QueryContext(ctx, cmd, args...); err != nil {
			if e, ok := err.(*mysql.MySQLError); ok {
				switch e.Number {
				case 1213: // Ignore deadlock
					continue
				}
			}
		}
		return
	}
}

func (api *MySqlAPI)GetPluginCounts(opt PluginListOpt)(count PluginCounts, err error){
	const queryCmd = "SELECT COUNT(`id`) AS `count`," +
		"SUM(`label_information`) AS `count_information`," +
		"SUM(`label_tool`) AS `count_tool`," +
		"SUM(`label_management`) AS `count_management`," +
		"SUM(`label_api`) AS `count_api`" +
		"FROM plugins AS a WHERE `enabled`=TRUE"

	cmd := queryCmd
	args := []any{}
	cmd, args = opt.appendTextFilter(cmd, args)
	cmd, args = opt.appendTagFilter(cmd, args)
	cmd, args = opt.appendOrderBy(cmd, args)
	cmd, args = opt.appendLimit(cmd, args)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	var (
		total, ctInfo, ctTool, ctMng, ctApi sql.NullInt32
	)
	if err = api.DB.QueryRowContext(ctx, cmd, args...).Scan(&total,
		&ctInfo, &ctTool, &ctMng, &ctApi); err != nil && err != sql.ErrNoRows {
		return
	}
	if total.Valid {
		count.Total = (int)(total.Int32)
		count.Information = (int)(ctInfo.Int32)
		count.Tool = (int)(ctTool.Int32)
		count.Management = (int)(ctMng.Int32)
		count.Api = (int)(ctApi.Int32)
	}
	err = nil
	return
}

func (api *MySqlAPI)GetPluginList(opt PluginListOpt)(infos []*PluginInfo, err error){
	const queryCmd = "SELECT a.`id`,a.`name`,a.`version`,a.`authors`,a.`desc`,a.`desc_zhCN`," +
		"CONVERT_TZ(a.`createAt`,@@session.time_zone,'+00:00') AS `utc_createAt`," +
		"CONVERT_TZ(a.`lastUpdate`,@@session.time_zone,'+00:00') AS `utc_lastUpdate`," +
		"`label_information`,`label_tool`,`label_management`,`label_api`," +
		"`github_sync`," +
		"CONVERT_TZ(`last_sync`,@@session.time_zone,'+00:00') AS `utc_last_sync`," +
		"SUM(b.`downloads`) AS `downloads`" +
		" FROM plugins as a LEFT JOIN plugin_releases as b" +
		" ON a.`id`=b.`id` WHERE a.`enabled`=TRUE"

	loger.Debugf("Getting plugin list with option %#v", opt)
	cmd := queryCmd
	args := []any{}
	cmd, args = opt.appendTextFilter(cmd, args)
	cmd, args = opt.appendTagFilter(cmd, args)
	cmd += " GROUP BY a.`id`"
	cmd, args = opt.appendOrderBy(cmd, args)
	cmd, args = opt.appendLimit(cmd, args)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	var rows *sql.Rows
	loger.Debugf("exec sql: %q", cmd)
	loger.Debugf("  args: %v", args)
	if rows, err = api.QueryContext(ctx, cmd, args...); err != nil {
		loger.Debugf("sql error: %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var (
			info PluginInfo
			authors string
			ghLastSync sql.NullTime
			downloads sql.NullInt64
		)
		if err = rows.Scan(&info.Id, &info.Name, &info.Version, &authors, &info.Desc, &info.Desc_zhCN, &info.CreateAt, &info.LastUpdate,
			&info.Labels.Information, &info.Labels.Tool, &info.Labels.Management, &info.Labels.Api,
			&info.GithubSync, &ghLastSync, &downloads); err != nil {
			return
		}
		info.Authors = strings.Split(authors, ",")
		if ghLastSync.Valid {
			info.LastSync = &ghLastSync.Time
		}
		if downloads.Valid {
			info.Downloads = downloads.Int64
		}
		infos = append(infos, &info)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

func (api *MySqlAPI)GetPluginInfo(id string)(info *PluginInfo, err error){
	const queryCmd = "SELECT a.`name`,a.`version`,a.`authors`,a.`desc`,a.`desc_zhCN`," +
		"CONVERT_TZ(a.`createAt`,@@session.time_zone,'+00:00') AS `utc_createAt`," +
		"CONVERT_TZ(a.`lastUpdate`,@@session.time_zone,'+00:00') AS `utc_lastUpdate`," +
		"a.`repo`,a.`repo_branch`,a.`repo_subdir`,a.`link`," +
		"`label_information`,`label_tool`,`label_management`,`label_api`," +
		"`github_sync`,`ghRepoOwner`,`ghRepoName`," +
		"CONVERT_TZ(`last_sync`,@@session.time_zone,'+00:00') AS `utc_last_sync`," +
		"SUM(b.`downloads`) AS `downloads`" +
		" FROM plugins as a LEFT JOIN plugin_releases as b" +
		" ON a.`id`=b.`id` WHERE a.`id`=? AND a.`enabled`=TRUE" +
		" GROUP BY a.`id`"
	const queryDependenciesCmd = "SELECT `target`,`tag` FROM plugin_dependencies WHERE `id`=?"
	var (
		authors string
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 7)
	defer cancel()

	var (
		ghLastSync sql.NullTime
		downloads sql.NullInt64
	)
	info = new(PluginInfo)
	if err = api.DB.QueryRowContext(ctx, queryCmd, id).
		Scan(&info.Name, &info.Version, &authors, &info.Desc, &info.Desc_zhCN, &info.CreateAt, &info.LastUpdate,
		&info.Repo, &info.RepoBranch, &info.RepoSubdir, &info.Link,
		&info.Labels.Information, &info.Labels.Tool, &info.Labels.Management, &info.Labels.Api,
		&info.GithubSync, &info.GhRepoOwner, &info.GhRepoName, &ghLastSync, &downloads); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return
	}
	info.Id = id
	if ghLastSync.Valid {
		info.LastSync = &ghLastSync.Time
	}
	if downloads.Valid {
		info.Downloads = downloads.Int64
	}
	info.Authors = strings.Split(authors, ",")
	info.Dependencies = make(DependMap, 3)
	var rows *sql.Rows
	if rows, err = api.QueryContext(ctx, queryDependenciesCmd, id); err != nil {
		return
	}
	defer rows.Close()
	var (
		pid string
		cond VersionCond
	)
	for rows.Next() {
		if err = rows.Scan(&pid, &cond); err != nil  {
			return
		}
		info.Dependencies[pid] = cond
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

func (api *MySqlAPI)GetPluginReadme(id string)(data []byte, prefix string, err error){
	var info *PluginInfo
	if info, err = api.GetPluginInfo(id); err != nil {
		return
	}
	if !info.GithubSync {
		filename := filepath.Join(PLUGIN_DIR, id, "README.MD")
		if data, err = os.ReadFile(filename); err != nil {
			return
		}
		return
	}
	// prefix only exists when fetching readme from github
	prefix = info.Link
	var res *http.Response
	baseurl, err := url.JoinPath("https://api.github.com", "repos",
		info.GhRepoOwner, info.GhRepoName, "readme")
	if err != nil {
		return
	}
	url0, err := url.JoinPath(baseurl, info.RepoSubdir)
	if err != nil {
		return
	}
	baseurl0 := baseurl + "?ref=" + info.RepoBranch
	url0 += "?ref=" + info.RepoBranch
	loger.Debugf("Getting readme for %s at %q", id, url0)
	res, err = githubCli.Get(url0)
	if e, ok := err.(*StatusCodeErr); ok && e.Code == http.StatusNotFound {
		loger.Debugf("Getting root readme for %s at %q", id, baseurl0)
		res, err = githubCli.Get(baseurl0)
	}
	if err != nil {
		if e, ok := err.(*StatusCodeErr); ok && e.Code == http.StatusNotFound {
			err = ErrNotFound
		}
		return
	}
	defer res.Body.Close()
	if data, err = io.ReadAll(res.Body); err != nil {
		return
	}
	var payload struct{
		Sha      string `json:"sha"`
		Size     int64  `json:"size"`
		Url      string `json:"url"`
		HtmlUrl  string `json:"html_url"`
		GitUrl   string `json:"git_url"`
		Download string `json:"download_url"`
		Type     string `json:"type"`
		Content  string `json:"content"`
		Encoding string `json:"encoding"`
	}
	if err = json.Unmarshal(data, &payload); err != nil {
		err = fmt.Errorf("JsonDecodeErr: %v", err)
		return
	}
	if payload.Encoding != "base64" {
		err = fmt.Errorf("Unexpect content enocding %q, expect base64", payload.Encoding)
		return
	}
	if data, err = base64.StdEncoding.DecodeString(payload.Content); err != nil {
		return
	}
	return
}

func (api *MySqlAPI)GetPluginReleases(id string)(releases []*PluginRelease, err error){
	const queryCmd = "SELECT `tag`,`enabled`,`stable`,`size`," +
		"CONVERT_TZ(`uploaded`,@@session.time_zone,'+00:00') AS `utc_uploaded`," +
		"`filename`,`downloads`,`github_url`" +
		" FROM plugin_releases WHERE `id`=? ORDER BY `uploaded` DESC"

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

func (api *MySqlAPI)GetPluginRelease(id string, tag Version)(release *PluginRelease, err error){
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
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
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

func (api *MySqlAPI)GetPluginReleaseAsset(id string, tag Version, filename string)(rc io.ReadSeekCloser, modTime time.Time, err error){
	filenam := filepath.Join(PLUGIN_DIR, id, "release", tag.String(), filepath.Clean(filename))
	var fd *os.File
	if fd, err = os.Open(filenam); err == nil {
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
	if len(release.GithubUrl) == 0 {
		err = ErrNotFound
		return
	}
	var resp *http.Response
	loger.Debugf("Downloading %q", release.GithubUrl)
	if resp, err = githubCli.Get(release.GithubUrl); err != nil {
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
	if err := os.MkdirAll(filepath.Dir(filenam), 0755); err == nil {
		if err = os.WriteFile(filenam, data, 0444); err != nil {
			loger.Warnf("Cannot write to asset file %q: %v", filenam, err)
		}else{
			loger.Infof("Cached %s(v%s):%s at %q: %v", id, tag, filename, filenam, err)
		}
	}else{
		loger.Warnf("Cannot make asset dir %q: %v", filepath.Dir(filenam), err)
	}
	rc = NopReadSeeker{bytes.NewReader(data)}
	return
}

func (opt PluginListOpt)appendTextFilter(cmd string, args []any)(string, []any){
	if len(opt.FilterBy) > 0 {
		cmd0, args0 := parseFilterBy(opt.FilterBy)
		if len(cmd0) > 0 {
			cmd += " AND (" + cmd0 + ")"
			args = append(args, args0...)
		}
	}
	return cmd, args
}

func parseFilterBy(filter string)(cmd string, args []any){
	cmds := []string{}
	for _, a := range strings.Split(filter, " ") {
		ok := true
		la := strings.ToLower(a)
		switch {
		case a == "":
		case strings.HasPrefix(la, "id:"):
			a = a[len("id:"):]
			cmds = append(cmds, "a.`id` LIKE ?")
			f := "%" + a + "%"
			args = append(args, f)
		case strings.HasPrefix(la, "n:"):
			a = a[len("n:"):]
			ok = false
			fallthrough
		case strings.HasPrefix(la, "name:"):
			if ok {
				a = a[len("name:"):]
			}
			cmds = append(cmds, "a.`name` LIKE ?")
			f := "%" + a + "%"
			args = append(args, f)
		case strings.HasPrefix(la, "@"):
			a = a[len("@"):]
			ok = false
			fallthrough
		case strings.HasPrefix(la, "a:"):
			if ok {
				a = a[len("a:"):]
				ok = false
			}
			fallthrough
		case strings.HasPrefix(la, "author:"):
			if ok {
				a = a[len("author:"):]
				ok = false
			}
			fallthrough
		case strings.HasPrefix(la, "authors:"):
			if ok {
				a = a[len("authors:"):]
			}
			for _, a := range strings.Split(a, ",") {
				cmds = append(cmds, "a.`authors` LIKE ?")
				f := "%" + a + "%"
				args = append(args, f)
			}
		case strings.HasPrefix(la, "d:"):
			a = a[len("d:"):]
			ok = false
			fallthrough
		case strings.HasPrefix(la, "desc:"):
			if ok {
				a = a[len("desc:"):]
				ok = false
			}
			fallthrough
		case strings.HasPrefix(la, "description:"):
			if ok {
				a = a[len("description:"):]
			}
			cmds = append(cmds, "a.`desc` LIKE ?")
			f := "%" + a + "%"
			args = append(args, f)
		default:
			cmds = append(cmds, "a.`id` LIKE ? OR a.`name` LIKE ? OR a.`authors` LIKE ? OR a.`desc` LIKE ?")
			f := "%" + a + "%"
			args = append(args, f, f, f, f)
		}
	}
	cmd = strings.Join(cmds, " OR ")
	return
}

func (opt PluginListOpt)appendTagFilter(cmd string, args []any)(string, []any){
	if len(opt.Tags) > 0 {
		cmds := []string{}
		for _, t := range opt.Tags {
			t = strings.ToLower(t)
			switch t {
			case "management", "tool", "information", "api":
				cmds = append(cmds, "`label_" + t + "`=TRUE")
			}
		}
		if len(cmds) > 0 {
			cmd += " AND (" + strings.Join(cmds, " OR ") + ")"
		}
	}
	return cmd, args
}

func (opt PluginListOpt)appendOrderBy(cmd string, args []any)(string, []any){
	switch strings.ToLower(opt.SortBy) {
	case "id", "name", "authors", "createAt", "lastUpdate":
		cmd += " ORDER BY a.`" + opt.SortBy + "`"
		rev := opt.Reversed
		switch opt.SortBy {
		case "lastUpdate":
			rev = !rev
		}
		if rev {
			cmd += " DESC"
		}
	case "downloads":
		cmd += " ORDER BY `" + opt.SortBy + "`"
		rev := opt.Reversed
		if !rev {
			cmd += " DESC"
		}
	}
	return cmd, args
}

func (opt PluginListOpt)appendLimit(cmd string, args []any)(string, []any){
	if opt.Limit > 0 || opt.Offset > 0 {
		if opt.Offset < 0 {
			opt.Offset = 0
		}
		cmd += " LIMIT ? OFFSET ?"
		args = append(args, opt.Limit, opt.Offset)
	}
	return cmd, args
}
