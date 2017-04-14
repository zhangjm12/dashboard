package chart

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"k8s.io/helm/pkg/helm/helmpath"
	"k8s.io/helm/pkg/repo"
)

// RepositorySpec is a specification for a repository.
type RepositorySpec struct {
	// Name of the repository.
	RepoName string `json:"name"`
	// URL of the repository.
	RepoURL string `json:"url"`
	// Phase of the repository.
	Phase string `json:"phase"`
}

// RepositoryListSpec is a specification for a repository.
type RepositoryListSpec struct {
	// List of repository names.
	Repositories []RepositorySpec `json:"repositories"`
	// Number of repositories.
	TotalItems int `json:"totalItems"`
}

// RepositoryListSpec is a specification for a repository.
type RepositoryChartListSpec struct {
	// List of charts.
	Charts []ChartSpec `json:"charts"`
}

// Chartspec representation view of a chart.
type ChartSpec struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	FullURL     string `json:"fullURL"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// AddRepository adds a repository.
func AddRepository(spec *RepositorySpec) error {
	helmHome := helmpath.Home(homePath())
	log.Printf("Adding repo named: %q with url: %q", spec.RepoName, spec.RepoURL)
	return addRepository(spec.RepoName, spec.RepoURL, helmHome)
}

// DeleteRepository removes a repository.
func DeleteRepository(repoName string) error {
	helmHome := helmpath.Home(homePath())
	return removeRepoLine(repoName, helmHome)
}

// GetRepository gets a specfied repository.
func GetRepository(repoName string) (*RepositorySpec, error) {
	helmHome := helmpath.Home(homePath())
	ensureHome(helmHome)
	repoSpec := &RepositorySpec{
		RepoName: repoName,
		RepoURL:  "",
		Phase:    "Available",
	}
	f, err := repo.LoadRepositoriesFile(helmHome.RepositoryFile())
	if err != nil {
		return repoSpec, err
	}
	for _, r := range f.Repositories {
		if r.Name == repoName {
			repoSpec.RepoURL = r.URL
		}
	}
	return repoSpec, nil
}

// GetRepositoryList get a list of repository.
func GetRepositoryList() (*RepositoryListSpec, error) {
	helmHome := helmpath.Home(homePath())
	ensureHome(helmHome)
	repoList := &RepositoryListSpec{
		Repositories: make([]RepositorySpec, 0),
	}
	f, err := repo.LoadRepositoriesFile(helmHome.RepositoryFile())
	if err != nil {
		return repoList, err
	}
	for _, r := range f.Repositories {
		repoSpec := &RepositorySpec{
			RepoName: r.Name,
			RepoURL:  r.URL,
			Phase:    "Available",
		}
		repoList.Repositories = append(repoList.Repositories, *repoSpec)
		repoList.TotalItems = repoList.TotalItems + 1
	}
	return repoList, nil
}

// GetRepositoryCharts get charts in a repository.
func GetRepositoryCharts(repoName string) (*RepositoryChartListSpec, error) {
	helmHome := helmpath.Home(homePath())
	chartList := &RepositoryChartListSpec{
		Charts: make([]ChartSpec, 0),
	}
	i, err := repo.LoadIndexFile(helmHome.CacheIndex(repoName))
	if err != nil {
		return chartList, err
	}
	for _, e := range i.Entries {
		for _, c := range e {
			if c == nil {
				continue
			}
			icon := c.Icon
			if icon == "" {
				icon = "https://deis.com/assets/images/svg/helm-logo.svg"
			}
			desc := c.Description
			if len(desc) > 45 {
				desc = desc[0:41] + "..."
			}
			chart := &ChartSpec{
				Name:        c.Name,
				Version:     c.Version,
				FullURL:     c.URLs[0],
				Description: desc,
				Icon:        icon,
			}
			chartList.Charts = append(chartList.Charts, *chart)
		}
	}
	return chartList, nil
}

func index(dir, url string) error {
	chartRepo, err := repo.LoadChartRepository(dir, url)
	if err != nil {
		return err
	}
	return chartRepo.Index()
}

func addRepository(name, url string, home helmpath.Home) error {
	entry := repo.Entry{
		Name:  name,
		URL:   url,
		Cache: fmt.Sprintf("%s-index.yaml", name),
	}
	cif := home.CacheIndex(name)
	rep, _ := repo.NewChartRepository(&entry)
	if err := rep.DownloadIndexFile(cif); err != nil {
		return fmt.Errorf("Looks like %q is not a valid chart repository or cannot be reached: %s", url, err.Error())
	}

	return insertRepoLine(name, url, home)
}

func removeRepoLine(name string, home helmpath.Home) error {
	repoFile := home.RepositoryFile()
	r, err := repo.LoadRepositoriesFile(repoFile)
	if err != nil {
		return err
	}

	if !r.Remove(name) {
		return fmt.Errorf("no repo named %q found", name)
	}
	if err := r.WriteFile(repoFile, 0644); err != nil {
		return err
	}

	if err := removeRepoCache(name, home); err != nil {
		return err
	}

	fmt.Printf("%q has been removed from your repositories", name)

	return nil
}

func removeRepoCache(name string, home helmpath.Home) error {
	if _, err := os.Stat(home.CacheIndex(name)); err == nil {
		err = os.Remove(home.CacheIndex(name))
		if err != nil {
			return err
		}
	}
	return nil
}

func insertRepoLine(name, url string, home helmpath.Home) error {
	cif := home.CacheIndex(name)
	f, err := repo.LoadRepositoriesFile(home.RepositoryFile())
	if err != nil {
		return err
	}

	if f.Has(name) {
		return fmt.Errorf("The repository name you provided (%s) already exists. Please specify a different name.", name)
	}
	f.Add(&repo.Entry{
		Name:  name,
		URL:   url,
		Cache: filepath.Base(cif),
	})
	return f.WriteFile(home.RepositoryFile(), 0644)
}
