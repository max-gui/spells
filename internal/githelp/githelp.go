package githelp

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"

	"github.com/max-gui/fileconvagt/pkg/fileops"
	"github.com/max-gui/logagent/pkg/logagent"
	"github.com/max-gui/spells/internal/pkg/constset"
)

// type Deployinfo struct {
// 	Buildenv         string `json:"BuildEnv"`
// 	Jenkinsnode      string `json:"jenkinsNode"`
// 	Gitrepositoryurl string `json:"GitRepositoryURL"`
// 	Releasename      string `json:"realseName"`
// 	Afversion        string `json:"AfVersion"`
// 	Gitbranch        string `json:"GitBranch"`
// 	Prfileactive     string `json:"prfileActive"`
// 	Dc               string `json:"dc"`
// 	Valstring        string
// 	Description      string
// 	Appname          string
// 	Appid            string
// }

type Writeinfo struct {
	Filepath string
	Content  string
	Del      bool
}

type RepoCloneResult struct {
	Name     string
	Repo     *git.Repository
	Msg      string
	Isupdate bool
}

func UpdateAll(c context.Context) map[string]RepoCloneResult {
	urlPathMap := []map[string]string{}
	log := logagent.Inst(c)
	urlPathMap = append(urlPathMap, map[string]string{"name": constset.Archname, "path": constset.Archpath, "url": constset.Archurl, "branch": "master"})
	urlPathMap = append(urlPathMap, map[string]string{"name": constset.Iacname, "path": constset.Iacpath, "url": constset.IacUrl, "branch": *constset.IacBranch})
	urlPathMap = append(urlPathMap, map[string]string{"name": constset.Templname, "path": constset.Templepath, "url": constset.Templurl, "branch": "master"})
	urlPathMap = append(urlPathMap, map[string]string{"name": constset.Dbname, "path": constset.DbPath, "url": constset.Dburl, "branch": "master"})

	results := make(map[string]RepoCloneResult)

	chain := make(chan RepoCloneResult, len(urlPathMap))
	var wg sync.WaitGroup
	for _, v := range urlPathMap {
		wg.Add(1)
		go func(info map[string]string) {
			defer func() {
				if e := recover(); e != nil {

					chain <- RepoCloneResult{Name: info["name"], Msg: fmt.Sprint(e)}
					wg.Done()
				}
			}()

			log.Println(info["name"])
			repo, err := CloneGetrepo(info["url"], info["branch"], info["path"], c)
			if err != nil && err != git.NoErrAlreadyUpToDate {
				log.Panic(err)
			}
			repores := RepoCloneResult{Name: info["name"], Repo: repo, Isupdate: true}
			if err == git.NoErrAlreadyUpToDate {
				repores.Isupdate = false
			}

			log.Println("done" + info["name"])
			chain <- repores
			wg.Done()
		}(v)
	}
	wg.Wait()
	close(chain)
	for v := range chain {
		if v.Repo == nil {
			log.Panic(v.Msg)
		}
		results[v.Name] = v
	}
	return results
}

// func Init(clsAppfig func()) {
// 	log.Println("initstart")
// 	doneArchrepo := make(chan int)
// 	doneiacrepo := make(chan int)
// 	doneTemplrepo := make(chan int)
// 	go func() {
// 		log.Println("inArchrepo")
// 		_, err := CloneGetrepo(*constset.Archurl, constset.Archpath)
// 		if err != nil {
// 			if err == git.NoErrAlreadyUpToDate {
// 				clsAppfig()

// 			} else {
// 				log.Panic(err)
// 			}
// 		}

// 		// GetProjWithTeam(constset.Archpath)
// 		log.Println("doneArchrepo")
// 		doneArchrepo <- 1
// 	}()
// 	go func() {
// 		log.Println("iniacrepo")
// 		_, err := CloneGetrepo(*constset.IacUrl, constset.Iacpath)
// 		if err != nil && err != git.NoErrAlreadyUpToDate {
// 			log.Panic(err)
// 		}
// 		log.Println("doneiacrepo")
// 		doneiacrepo <- 1
// 	}()
// 	go func() {

// 		log.Println("inTemplrepo")
// 		_, err := CloneGetrepo(*constset.Templurl, constset.Templepath)
// 		if err != nil {
// 			if err == git.NoErrAlreadyUpToDate {
// 				templ.ClsGempl()
// 			} else {
// 				log.Panic(err)
// 			}
// 		}

// 		files := fileops.GetAllFiles(constset.Templepath)
// 		for _, file := range files {
// 			fileinfo := strings.Split(file, ".")
// 			if len(fileinfo) > 1 {
// 				switch fileinfo[0] {
// 				case "defaultconfig":
// 					defig.GenDefig(true)
// 				case "Dockerfile":
// 					dockfig.MakeDockemple(fileinfo[1], true)
// 				case "jenkins":
// 					jenfig.MakeJenkimple(fileinfo[1], true)
// 				case "values":
// 					valfig.MakeValuemple(fileinfo[1], true)

// 				}
// 			}
// 		}
// 		log.Println("doneTemplrepo")
// 		doneTemplrepo <- 1
// 	}()
// 	<-doneArchrepo
// 	<-doneiacrepo
// 	<-doneTemplrepo
// 	// close(doneArchrepo)
// 	// close(doneiacrepo)
// 	// close(doneTemplrepo)
// 	log.Println("doneinit")

// }

func CloneGetrepo(repourl, branch, localpath string, c context.Context) (*git.Repository, error) {

	var r *git.Repository
	var err error
	var w *git.Worktree
	log := logagent.Inst(c)
	if _, err = os.Stat(localpath); os.IsNotExist(err) {
		// os.MkdirAll(localpath, 0777)
		// os.MkdirAll(localpath, 0777) //os.ModeDir.Perm())
		err = nil
		publicKey, err := ssh.NewPublicKeys("git", []byte(constset.Sshkey), "")
		if err != nil {
			log.Panic(err)
		}
		r, err = git.PlainClone(localpath, false, &git.CloneOptions{
			URL:               repourl,
			Auth:              publicKey,
			ReferenceName:     plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
			SingleBranch:      true,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		})
		if err != nil {
			log.Panic(err)
		}
	} else {
		r, err = git.PlainOpen(localpath)
		if err != nil {
			log.Panic(err)
		}
		w, err = r.Worktree()
		if err != nil {
			log.Panic(err)
		}
		publicKey, err := ssh.NewPublicKeys("git", []byte(constset.Sshkey), "")
		if err != nil {
			log.Panic(err)
		}
		// _, err = git.PlainClone(localpath, false, &git.CloneOptions{
		// 	URL:               repourl,
		// 	Auth:              publicKey,
		// 	RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		// })
		err = w.Pull(&git.PullOptions{
			RemoteName:        "origin",
			Auth:              publicKey,
			ReferenceName:     plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch)),
			SingleBranch:      true,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth})
		if err != nil {
			log.Print(err)
		}
	}

	return r, err
}

func CommitPushFiles(filesinfo []Writeinfo, repo *git.Repository, perfixstr string, c context.Context) bool {

	log := logagent.Inst(c)
	if len(filesinfo) > 0 {
		// perfixstr := constset.Iacpath
		w, err := repo.Worktree()
		if err != nil {
			log.Panic(err)
		}
		for _, chfile := range filesinfo {
			log.Print(strings.TrimPrefix(chfile.Filepath, perfixstr))
			if chfile.Del {
				os.Remove(chfile.Filepath)
				// _, err = w.Remove(strings.TrimPrefix(chfile.Filepath, perfistr))
				// if err != nil && err != index.ErrEntryNotFound {
				// 	log.Panic(err)
				// }
			} else {
				fileops.Writeover(chfile.Filepath, chfile.Content, c)

			}
			_, err = w.Add(strings.TrimPrefix(chfile.Filepath, perfixstr))
			if err != nil {
				log.Panic(err)
			}

		}

		isupdate := commitPush(w, repo, c)
		return isupdate
	} else {
		return false
	}
}

func commitPush(w *git.Worktree, r *git.Repository, c context.Context) bool {
	log := logagent.Inst(c)
	st, err := w.Status()
	if err != nil {
		log.Print(err)
	}
	fmt.Println(st)
	if st.IsClean() {
		return false
	}
	_, err = w.Commit(*constset.Commitmsg, &git.CommitOptions{ //提交
		Author: &object.Signature{
			Name:  *constset.Gitname,
			Email: *constset.Gitemail,
			When:  time.Now(),
		}})
	if err != nil {
		log.Panic(err)
	}
	publicKey, err := ssh.NewPublicKeys("git", []byte(constset.Sshkey), "")

	if err != nil {
		log.Panic(err)
	}
	// auth, _ := publicKey()
	err = r.Push(&git.PushOptions{Auth: publicKey})
	if err != nil {
		log.Panic(err)
	}

	return true
}
