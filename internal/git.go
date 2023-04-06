/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package internal

import (
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
)

func GetGitRevision(gitUrl string) (revisionDetails map[string]string) {

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

	revisionDetails = make(map[string]string)
	revisionDetails["message"] = fmt.Sprintf(commit.Message)
	revisionDetails["id"] = fmt.Sprint(commit.Hash)
	revisionDetails["author"] = fmt.Sprint(commit.Author)

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
