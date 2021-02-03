/*
    Copyright (C) 2020 Accurics, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package initialize

import (
	"fmt"
	"os"
	"strings"

	"github.com/accurics/terrascan/pkg/config"
	"go.uber.org/zap"
	"gopkg.in/src-d/go-git.v4"
	gitConfig "gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

var (
	basePath       = config.GetPolicyRepoPath()
	basePolicyPath = config.GetPolicyBasePath()
	repoURL        = config.GetPolicyRepoURL()
	branch         = config.GetPolicyBranch()
)

// Run initializes terrascan if not done already
func Run() error {

	zap.S().Debug("initializing terrascan")

	// check if policy paths exist
	if path, err := os.Stat(basePolicyPath); err == nil && path.IsDir() {
		return refreshPolicies()
	}

	// download policies
	os.RemoveAll(basePath)
	if err := DownloadPolicies(); err != nil {
		return err
	}

	zap.S().Debug("intialized successfully")
	return nil
}

// DownloadPolicies clones the policies to a local folder
func DownloadPolicies() error {

	// clone the repo
	r, err := git.PlainClone(basePath, false, &git.CloneOptions{
		URL: repoURL,
	})
	if err != nil {
		zap.S().Errorf("failed to clone repository %s. error: '%v'", repoURL, err)
		return err
	}

	// fetch references
	err = fetch(r)
	if err != nil {
		zap.S().Errorf("failed to fetch, fetchURL: %s. error: '%v'", repoURL, err)
		return err
	}

	// checkout policies branch
	err = checkout(r)
	if err != nil {
		zap.S().Errorf("failed to checkout to branch: %s. error: '%v'", branch, err)
		return err
	}

	return nil
}

// this function will either pull or cal download policies
func refreshPolicies() error {
	r, err := git.PlainOpen(basePath)
	if err != nil {
		return err
	}

	remote, err := r.Remote("origin")
	if err != nil {
		return err
	}
	remoteConfig := remote.Config()
	if err := remoteConfig.Validate(); err != nil {
		return err
	}

	// size of the urls cannot be empty as it is validted above
	if strings.EqualFold(remoteConfig.URLs[0], repoURL) {
		_, err := r.Branch(branch)
		if err != nil {
			// branch does not exist
			if err := fetch(r); err != nil {
				return err
			}
			if err := checkout(r); err != nil {
				return err
			}
		}
		// branch exists, pull from remote repository
		if err := pull(r); err != nil {
			return err
		}

	} else {
		// repoURL is not same as remote origin fetch url
		// delete the basePath and download policies
		os.RemoveAll(basePath)
		return DownloadPolicies()
	}

	return nil
}

func fetch(r *git.Repository) error {
	// fetch references
	err := r.Fetch(&git.FetchOptions{
		RefSpecs: []gitConfig.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
	})
	if err != nil {
		zap.S().Errorf("failed to fetch references from repo. error: '%v'", err)
		return err
	}
	return nil
}

func checkout(r *git.Repository) error {
	// get the work tree
	w, err := r.Worktree()
	if err != nil {
		zap.S().Errorf("failed to create working tree. error: '%v'", err)
		return err
	}

	// checkout policies branch
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
		Force:  true,
	})
	if err != nil {
		zap.S().Errorf("failed to checkout branch '%v'. error: '%v'", branch, err)
		return err
	}
	return nil
}

func pull(r *git.Repository) error {
	// create working tree
	w, err := r.Worktree()
	if err != nil {
		zap.S().Errorf("failed to create working tree. error: '%v'", err)
		return err
	}
	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
	})
	if err != nil {
		if strings.EqualFold(err.Error(), git.NoErrAlreadyUpToDate.Error()) {
			// repo is up to date
			zap.S().Info("repository is already up to date")
			return nil
		}
		return err
	}
	return nil
}
