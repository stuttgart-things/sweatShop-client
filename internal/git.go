/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	http "github.com/go-git/go-git/v5/plumbing/transport/http"

	memfs "github.com/go-git/go-billy/v5/memfs"
	memory "github.com/go-git/go-git/v5/storage/memory"

	"github.com/go-git/go-git/v5"
)

func GetGitRevision(gitUrl string) (revisionDetails map[string]interface{}) {

	// Clones the given repository, creating the remote, the local branches
	// and fetching the objects, everything in memory:
	Info("git clone " + gitUrl)
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: gitUrl,
	})
	CheckIfError(err)

	ref, err := r.Head()
	CheckIfError(err)

	// ... retrieves the commit history
	since := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	until := time.Date(2035, 7, 30, 0, 0, 0, 0, time.UTC)

	cIter, err := r.Log(&git.LogOptions{From: ref.Hash(), Since: &since, Until: &until})
	CheckIfError(err)

	commit, _ := cIter.Next()
	cIter.Close()

	CheckIfError(err)

	revisionDetails = make(map[string]interface{})

	str1 := fmt.Sprint(commit.Author)
	re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)

	submatchall := re.FindAllString(str1, -1)
	for _, element := range submatchall {
		fmt.Println(element)
	}

	repoName, _ := GetRegexSubMatch(gitUrl, `^(.*/)?(?:$|(.+?)(?:(\.[^.]*$)|$))`, 2)
	author, _ := GetRegexSubMatch(fmt.Sprint(commit.Author), `\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`, 0)
	date, _ := GetRegexSubMatch(fmt.Sprint(commit.Author), `([a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+\.[a-zA-Z0-9_-]+)`, 0)

	revisionDetails["name"] = repoName
	revisionDetails["url"] = gitUrl
	revisionDetails["date"] = author
	revisionDetails["message"] = fmt.Sprintf(commit.Message)
	revisionDetails["id"] = fmt.Sprint(commit.Hash)
	revisionDetails["author"] = date

	return
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))

}

func GetYachtConfig(repository, filename string, auth *http.BasicAuth) string {

	// Init memory storage and fs
	storer := memory.NewStorage()
	fs := memfs.New()

	// Clone repo into memfs
	r, err := git.Clone(storer, fs, &git.CloneOptions{
		URL:  repository,
		Auth: auth,
	})

	if err != nil {
		fmt.Println("Could not git clone repository %s: %w", repository, err)
	}
	fmt.Println("Repository cloned")

	// Get git default worktree
	w, err := r.Worktree()
	if err != nil {
		fmt.Println("Could not get git worktree: %w", err)
	}

	fmt.Println(w)
	files, _ := fs.ReadDir("/")

	for _, file := range files {
		fmt.Println(file.Name())
	}

	file, _ := fs.Open(filename)
	fileContent, _ := ioutil.ReadAll(file)

	return string(fileContent)

}

func GetRegexSubMatch(scanText, regexPattern string, group int) (string, bool) {

	rgx := regexp.MustCompile(regexPattern)
	rs := rgx.FindStringSubmatch(scanText)

	if len(rs) == 0 {
		return "", false
	}

	return strings.Trim(rs[group], " "), true
}
