
package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"time"
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

type DependMap map[string]VersionCondList
type RequireMap map[string]string

type PluginInfo struct {
	Id           string       `json:"id"`
	Name         string       `json:"name"`
	Version      Version      `json:"version"`
	Authors      []string     `json:"authors"`
	Desc         string       `json:"desc,omitempty"`
	Desc_zhCN    string       `json:"desc_zhCN,omitempty"`
	CreateAt     time.Time    `json:"createAt"`
	LastRelease  *time.Time   `json:"lastRelease,omitempty"`
	Repo         string       `json:"repo,omitempty"`
	RepoBranch   string       `json:"repoBranch,omitempty"`
	RepoSubdir   string       `json:"repoSubdir,omitempty"`
	Link         string       `json:"link,omitempty"`
	Labels       PluginLabels `json:"labels"`
	Downloads    int64        `json:"downloads"`
	Dependencies DependMap    `json:"dependencies,omitempty"`
	Requirements RequireMap   `json:"requirements,omitempty"`
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

type Content struct{
	Data func()([]byte, error)
	CloseFunc func()(error)
	ModTime string
	URLPrefix string
	DataURLPrefix string
}

func (c Content)Close()(error){
	if c.CloseFunc == nil {
		return nil
	}
	return c.CloseFunc()
}

type API interface {
	GetLastUpdateTime()(modTime time.Time, err error)
	GetPluginLastUpdateTime(id string)(modTime time.Time, err error)
	GetPluginCounts(opt PluginListOpt)(count PluginCounts, err error)
	GetPluginList(opt PluginListOpt)(infos []*PluginInfo, err error)
	GetPluginIdList(opt PluginListOpt)(ids []string, err error)
	GetPluginInfo(id string, version string)(info *PluginInfo, err error)
	GetPluginInfos(id string)(info []*PluginInfo, err error)
	GetPluginReadme(id string)(content Content, err error)
	GetPluginReleases(id string)(releases []*PluginRelease, err error)
	GetPluginRelease(id string, tag Version)(release *PluginRelease, err error)
	GetPluginReleaseAsset(id string, tag Version, filename string)(rc io.ReadSeekCloser, modTime time.Time, err error)
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
