package altconfig

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/max-gui/logagent/pkg/logagent"
	"github.com/max-gui/spells/internal/githelp"
	"github.com/max-gui/spells/internal/iac/archfig"
	"github.com/max-gui/spells/internal/iac/defig"
	"github.com/max-gui/spells/internal/iac/dockfig"
	"github.com/max-gui/spells/internal/iac/jenfig"
	"github.com/max-gui/spells/internal/iac/templ"
	"github.com/max-gui/spells/internal/pkg/constset"
)

func Arch_commits(commitinfo CommitCheckHookInfo, install bool, c context.Context) {
	// var archr *git.Repository
	log := logagent.Inst(c)
	var iacr *git.Repository
	if install {
		gitres := githelp.UpdateAll(c)
		iacr = gitres[constset.Iacname].Repo
		// _, err := githelp.CloneGetrepo(*constset.Archurl, constset.Archpath)
		// if err != nil && err != git.NoErrAlreadyUpToDate {
		// 	log.Panic(err)
		// }

		// _, err = githelp.CloneGetrepo(*constset.Templurl, constset.Templepath)
		// if err != nil && err != git.NoErrAlreadyUpToDate {
		// 	log.Panic(err)
		// }
		// iacr, err = githelp.CloneGetrepo(*constset.IacUrl, constset.Iacpath)
		// if err != nil && err != git.NoErrAlreadyUpToDate {
		// 	log.Panic(err)
		// }
	}
	var filesinfo = []githelp.Writeinfo{}
	// var addeds = map[string]struct{}{}
	// var dels = map[string]struct{}{}
	// for _, c := range commitinfo.ChangedFiles {
	// 	if c.Status == "A" {
	// 		spath := strings.Split(c.Path, "/")
	// 		addeds[spath[len(spath)-1]] = struct{}{}
	// 		// addeds = append(addeds, spath[len(spath)-1])
	// 	}
	// 	if c.Status == "D" {
	// 		spath := strings.Split(c.Path, "/")
	// 		dels[spath[len(spath)-1]] = struct{}{}
	// 	}
	// }

	// extfiles := Findarchfile(addeds)

	// for del := range dels {
	// 	// if v = del {
	// 	delete(extfiles, del)
	// 	// }
	// 	// }
	// }

	// if len(extfiles) > 0 {
	// 	log.Panicf("these apps need to be removed:%v", extfiles)
	// }

	// for _, c := range commitinfo.ChangedFiles {
	// 	// for _, a := range c..Added {
	// 	if strings.Contains(c.Path, "yaml") && (c.Status == "D") {
	// 		// projinfo := strings.Split(c.Path, constset.PthSep)

	// 		appconf := archfig.GetArchfig(c)

	// 		// os.RemoveAll(constset.iacpath + "app" + constset.PthSep + "iac" + constset.PthSep + appconf.Application.Appid + constset.PthSep + appconf.Application.Name)
	// 		changes := ArchAltGenWithChanges(appconf, install, true)
	// 		archfig.RMArch(appconf)
	// 		filesinfo = append(filesinfo, changes...)
	// 	}
	// 	// } else if c.Status == "D" {
	// 	// 	strs := strings.Split(c.Path, "/")
	// 	// 	templ.RMtempl(strs[len(strs)-1])
	// 	// }

	// 	// }
	// }

	for _, changed := range commitinfo.ChangedFiles {
		// for _, a := range c..Added {
		if strings.Contains(changed.Status, "R") {
			log.Panic("delete and add can't be in one commit,please split and push")
		} else if strings.Contains(changed.Path, "yaml") {
			if changed.Status == "D" && install {
				// projinfo := strings.Split(c.Path, constset.PthSep)
				appconf := archfig.GetArchfig(changed, c)

				// os.RemoveAll(constset.iacpath + "app" + constset.PthSep + "iac" + constset.PthSep + appconf.Application.Appid + constset.PthSep + appconf.Application.Name)
				changes, _ := ArchAltGenWithChanges(appconf, install, true, c)
				appconf.RMArch(c)
				filesinfo = append(filesinfo, changes...)
			} else if changed.Status == "A" || changed.Status == "M" {

				appconf := archfig.GetArchfigByGitContent(changed, false, c)

				if changed.Status == "A" {
					if flag, v := archfig.AppExist(appconf.Application.Name, c); flag {
						log.Panicf("these apps need to be removed:%v", v)
					}
					// path := orgappconf.Application.Team + constset.PthSep +
					// 	orgappconf.Application.Project + constset.PthSep +
					// 	orgappconf.Application.Name + ".yaml"
					// log.Panicf("these apps need to be removed:%v", path)

				}

				changes, _ := ArchAltGenWithChanges(appconf, install, false, c)
				if install {
					appconf.Install(c)
				}
				filesinfo = append(filesinfo, changes...)
			}
		}
		// } else if c.Status == "D" {
		// 	strs := strings.Split(c.Path, "/")
		// 	templ.RMtempl(strs[len(strs)-1])
		// }

		// }
	}
	githelp.CommitPushFiles(filesinfo, iacr, constset.Iacpath, c)
}

func Tmpt(filesinfo []githelp.Writeinfo, repoUrl, branch, repoPath string, install bool, c context.Context) {
	log := logagent.Inst(c)
	if install {
		repor, err := githelp.CloneGetrepo(repoUrl, branch, repoPath, c)
		if err != nil && err != git.NoErrAlreadyUpToDate {
			log.Panic(err)
		}

		githelp.CommitPushFiles(filesinfo, repor, repoPath, c)
	}
}

func Arch_commit(filecontentinfo archfig.FileContentInfo, repourl string, install bool, c context.Context) ([]githelp.Writeinfo, []archfig.Arch_config) {
	var changes []githelp.Writeinfo
	var archalts []archfig.Arch_config
	log := logagent.Inst(c)
	if strings.Contains(filecontentinfo.Path, ".yaml") {

		if strings.Contains(filecontentinfo.Status, "R") {
			log.Panic("delete and add can't be in one commit,please split and push")
		} else {
			pathinfo := strings.Split(filecontentinfo.Path, constset.PthSep)
			appfname := strings.TrimSuffix(pathinfo[len(pathinfo)-1], ".yaml")
			if filecontentinfo.Status == "D" && install {
				// projinfo := strings.Split(c.Path, constset.PthSep)

				appconf := archfig.GetArchfigSin(appfname, c)

				// os.RemoveAll(constset.iacpath + "app" + constset.PthSep + "iac" + constset.PthSep + appconf.Application.Appid + constset.PthSep + appconf.Application.Name)
				changes, _ = ArchAltGenWithChanges(appconf, install, true, c)
				appconf.RMArch(c)
			} else if filecontentinfo.Status == "A" || filecontentinfo.Status == "M" {
				appconf := archfig.GetArchfigByGitContentSin(filecontentinfo.Content, appfname, false, c)
				appconf.Application.Repositry = repourl
				if filecontentinfo.Status == "A" {
					if flag, v := archfig.AppExist(appconf.Application.Name, c); flag {
						log.Panicf("these apps need to be removed:%v", v)
					}
					// path := orgappconf.Application.Team + constset.PthSep +
					// 	orgappconf.Application.Project + constset.PthSep +
					// 	orgappconf.Application.Name + ".yaml"
					// log.Panicf("these apps need to be removed:%v", path)

				}

				changes, archalts = ArchAltGenWithChanges(appconf, install, false, c)
			}
		}
	}

	return changes, archalts
}

func Findarchfile(filenames map[string]struct{}, c context.Context) map[string]struct{} {
	dirPth := constset.Archpath
	log := logagent.Inst(c)

	var fullpath = map[string]struct{}{}
	filepath.WalkDir(dirPth, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			log.Panic(err)
		}
		// log.Print(path)
		if !d.IsDir() {
			for k := range filenames {
				if d.Name() == k {
					fullpath[strings.TrimPrefix(path, dirPth)] = struct{}{}
					// fullpath = append(fullpath, strings.TrimPrefix(path, dirPth))
				}
			}
		}
		return nil
	})
	return fullpath
}

func ArchAltGenWithChanges(appconf archfig.Arch_config, install bool, del bool, c context.Context) ([]githelp.Writeinfo, []archfig.Arch_config) {

	var filesinfo = []githelp.Writeinfo{}
	log := logagent.Inst(c)
	if appconf.Application.Ungenfig {
		return filesinfo, []archfig.Arch_config{}
	}
	dockstr, jenconfig, jenstr := archAltGen(appconf, c)
	// changes = append(changes, basepath+"Jenkinsfile")
	log.Print(jenconfig)
	log.Print(jenstr)
	log.Print(dockstr)

	if install {
		// appconf.Install()
		basepath := constset.Iacpath + "app" + constset.PthSep + "iac" + constset.PthSep + appconf.Application.Appid + constset.PthSep + appconf.Application.Name + constset.PthSep
		filesinfo = append(filesinfo, githelp.Writeinfo{Filepath: basepath + "Dockerfile", Content: dockstr, Del: del})
		filesinfo = append(filesinfo, githelp.Writeinfo{Filepath: basepath + "Jenkinsfile", Content: jenstr, Del: del})

		return filesinfo, []archfig.Arch_config{appconf}
	}
	return filesinfo, []archfig.Arch_config{}
}

func archAltGen(appconf archfig.Arch_config, c context.Context) (string, jenfig.JenkinsInfo, string) {
	// valconfig := valfig.GenValfig(appconf, envinfo, envdc)
	dockstr := dockfig.GenDocfile(appconf, c)

	jenconfig := jenfig.GenJenfig(appconf)

	jenstr := jenfig.GenJenfile(jenconfig, c)
	return dockstr, jenconfig, jenstr
}

func Archdef_commit(commitcheck CommitCheckHookInfo, install bool, c context.Context) {
	for _, changed := range commitcheck.ChangedFiles {
		if changed.Status == "A" || changed.Status == "M" {

			templAlt(changed.Path, changed.Content, install, c)

		} else if changed.Status == "D" {
			strs := strings.Split(changed.Path, "/")
			templ.RMtempl(strs[len(strs)-1])
		}
	}
}

func templAlt(filepath, content string, install bool, c context.Context) {
	strs := strings.Split(filepath, "/")
	confinfo := strings.Split(strs[len(strs)-1], ".")

	switch confinfo[0] {
	case "defaultconfig":
		// iac.GenDefig(install)
		defig.GenDefigFrom([]byte(content), install, c)
	case "Dockerfile":
		// iac.MakeDockemple(confinfo[1], install)
		templ.Getempl(content, strs[len(strs)-1], install, c)
	case "jenkins":
		// iac.MakeJenkimple(confinfo[1], install)
		templ.Getempl(content, strs[len(strs)-1], install, c)
	case "values":
		// iac.MakeValuemple(confinfo[1], install)
		templ.Getempl(content, strs[len(strs)-1], install, c)
	}
}

type CommitInfo struct {
	CreatedDate int64
	CreatedBy   string
	Message     string   // 提交信息
	NodeId      string   // ": "6204231647c1c813c832f6b9819f0f2758152c7b",   ##提交版本号
	Added       []string // ": ["README.md"],        ##新增文件
	Modified    []string // ":["com/test.java"],  ##修改文件
	Removed     []string //": ["test.txt"]        ##删除文件
}

type RepositoryInfo struct {
	SshUrl      string //": "ssh://",  ##代码库ssh地址
	HttpUrl     string //": "http://",  ##代码库http地址
	Name        string //": "sqms-code",      ##代码库名称
	Namespace   string //": "git"        ##项目组名称
	Description string //": "123",     ##代码库描述
	Id          string //": 866,                ##代码库id
}

type CommitHookInfo struct {
	CurrentBranch       string
	Commit              CommitInfo
	Commits             []CommitInfo
	Total_commits_count int    //": 2           ##本次提交事件包含commit个数
	Before              string //": "eb531858347c34486455c0d47e90ad8229930daf",       ##推送前提交版本号
	After               string //":  "6204231647c1c813c832f6b9819f0f2758152c7b",       ##推送后sha值
	Repository          RepositoryInfo
}

type CommitCheckHookInfo struct {
	Commiter       string //提交人
	Branch         string //分支
	Type           string //git,svn
	ChangedFiles   []archfig.FileContentInfo
	CommitIds      []string
	Messages       []string
	RepositoryName string
	Namespace      string
}
