package main

import (
	"io"
	"os"
	"os/exec"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"github.com/codegangsta/cli"
	"github.com/davecgh/go-spew/spew"
	"path/filepath"
)

type myRegexp struct {
	*regexp.Regexp
}

type Repo struct {
	SUBPATH     string
	REMOTE_URL  string
	REMOTE_NAME string
	HEADS       []Head
	TAGS        []Tag
}
type Tag struct {
	tag string
}
type Head struct {
	head string

}


func exe_cmd(cmd string, wg *sync.WaitGroup) {
	fmt.Println("command is ", cmd)
	// splitting head => g++ parts => rest of the command
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:len(parts)]

	out, err := exec.Command(head, parts...).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s", out)
	wg.Done() // Need to signal to waitgroup that this goroutine is done
}

func main() {
	app := cli.NewApp()
	app.Name = "git-subsplit"
	app.Usage = "Automate and simplify the process of managing one-way read-only subtree splits."
	app.Version = "1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "work-dir, w",
			Value: "--work-dir",
			Usage: "directory that contains the subsplit working directory",
		},
		cli.StringFlag{
			Name: "heads, g",
			Value: "--heads",
			Usage: "Only publish for listed heads instead of all heads",
		},
		cli.StringFlag{
			Name: "no-heads, i",
			Value: "--no-heads",
			Usage: "Do not publish any heads",
		},
		cli.StringFlag{
			Name: "tags, t",
			Value: "--tags",
			Usage: "Only publish for listed tags instead of all tags",
		},
		cli.StringFlag{
			Name: "no-tags, j",
			Value: "--no-tags",
			Usage: "Do not publish any tag",
		},
		cli.StringFlag{
			Name: "rebuild-tags,r",
			Value: "--rebuild-tags",
			Usage: "Rebuild all tags (as opposed to skipping tags that are already synced)",
		},
		cli.StringFlag{
			Name: "update, u",
			Value: "--update",
			Usage: "Fetch updates from repository before publishing",
		},
		cli.StringFlag{
			Name: "dry-run, n",
			Value: "--dry-run",
			Usage: "Do everything except actually send the updates",
		},
		cli.StringFlag{
			Name: "annotate, a",
			Value: "--annotate",
			Usage: "annotate the repository",
		},

		cli.StringFlag{
			Name: "quiet, q",
			Value: "--quiet",
			Usage: "Do not display output",
		},
		cli.StringFlag{
			Name: "debug, d",
			Value: "--debug",
			Usage: "Show Debugging Output",
		},

	}

	app.Commands = []cli.Command{
		{
			Name:      "init",
			Aliases:     []string{"-i"},
			Usage:     "Initialize a subsplit from origin",
			Action: func(c *cli.Context) {
				fmt.Printf("Initialize Task: Starting")

				cmd := fmt.Sprintf("git clone -q %s ", c.Args().First())
				_, err := exec.Command("sh", "-c", cmd).Output()

				abs_repo := strings.Split(c.Args().First(), "/")

				dir := strings.Replace(abs_repo[len(abs_repo)-1], ".git", "", -1)

				f, err := os.Create(dir + "/.subsplit")
				f.Close()

				if err != nil {
					fmt.Printf("Initialize Task: Failed : %s", err)

				}else {
					file, err := os.Open(dir + "/.gitignore")

					if err != nil {
						fmt.Println(err)

						f, err := os.Create(dir +"/.gitignore")

						if err == nil {

							n, err := io.WriteString(f, dir + "/.subsplit")

							if err != nil {
								fmt.Println(n, err)
							}
						}else{
							fmt.Println(err)
						}
						f.Close()
					}else {
						n, err := io.WriteString(file, dir + "/.subsplit")

						if err != nil {
								fmt.Println(n, err)
							}
					}

					fmt.Printf("Initialize Task: Successful")
				}

			},
		},
		{
			Name:      "publish",
			Aliases:     []string{"-p"},
			Usage:     "This command will create subtree splits of the project's repository branches and tags. It will then push each branch and tag to the repository dedicated to the subtree.",
			Action: func(c *cli.Context) {

				dir, _ := filepath.Abs(filepath.Dir(os.Args[0]));
				dir += "/.subsplit"

				if _, err := os.Stat(dir); err == nil {
					println("Publish Task: Starting")
					input_repos := strings.Split(c.Args().First(), ",")
					tags := ""
					heads := ""

					Repos := []Repo{}

					for _, val := range input_repos {

						r := Repo{
							SUBPATH:getSubPath(val),
							REMOTE_URL:getRemoteUrl(val),
							REMOTE_NAME:getRemoteName(val),
							HEADS:getHeads(heads),
							TAGS :getTags(tags),
						}
						Repos = append(Repos, r)
					}

					for _, r := range Repos {
						spew.Dump(r.HEADS)

						//cmd := fmt.Sprintf("git remote add r.REMOTE_NAME r.REMOTE_URL")
						//output, _ := exec.Command("sh", "-c", cmd).Output()
					}

					println("Syncing Task: Started")
					println("Syncing Task: Successful")

					println("Publish task: Successful")
				}else {
					println("Error: no .subsplit found, has the repo been initalized with git-subsplit init ?")
				}
			},
		}, {
			Name:      "Update",
			Aliases:     []string{"-u"},
			Usage:     "Update ",
			Action: func(c *cli.Context) {
				println("Updating subsplit from origin ")
			},
		},

	}

	app.Run(os.Args)
}


func getHeads(data string) []Head {
	println("Querying heads from origin")
	cmd := fmt.Sprintf("git ls-remote origin")
	output, _ := exec.Command("sh", "-c", cmd).Output()
	regex := regexp.MustCompile(`refs/heads/(.*)`)
	matches := regex.FindStringSubmatch(string(output))

	Heads := []Head{}

	for _, val := range matches {

		h := Head{
			head : val,
		}
		Heads = append(Heads, h)
	}
	return Heads
}

func getTags(data string) []Tag {
	println("Querying tags from origin")
	cmd := fmt.Sprintf("git ls-remote origin")
	output, _ := exec.Command("sh", "-c", cmd).Output()
	regex := regexp.MustCompile(`refs/tags/(.*)`)
	matches := regex.FindStringSubmatch(string(output))

	Tags := []Tag{}

	for _, val := range matches {

		t := Tag{
			tag : val,
		}
		Tags = append(Tags, t)
	}
	return Tags
}

func getSubPath(data string) string {
	println("Parsing the path ")
	result := strings.Split(data, ":")
	return result[0]
}

func getRemoteUrl(data string) string {
	println("Parsing the remote url ")
	result := strings.Split(data, ":")
	return result[1] + ":" + result[2]
}

func getRemoteName(data string) string {
	println("Generating the remote name ")
	cmd := fmt.Sprintf("echo " + data + " | git hash-object --stdin")
	byteArray, _ := exec.Command("sh", "-c", cmd).Output()
	result := string(byteArray[:])
	return strings.Replace(result, "\n", "", -1)
}

func (r *Repo)syncTags(){
	println("Syncing Tags ")
	//  Todo: convert from sh to go
//	for TAG in $TAGS
//		do
//			if [ -n "$VERBOSE" ];
//			then
//				echo "${DEBUG} git show-ref --quiet --verify -- \"refs/tags/${TAG}\""
//			fi
//
//			if ! git show-ref --quiet --verify -- "refs/tags/${TAG}"
//			then
//				say " - skipping tag '${TAG}' (does not exist)"
//				continue
//			fi
//			LOCAL_TAG="${REMOTE_NAME}-tag-${TAG}"
//
//			if [ -n "$VERBOSE" ];
//			then
//				echo "${DEBUG} LOCAL_TAG="${LOCAL_TAG}""
//			fi
//
//			if git branch | grep "${LOCAL_TAG}$" >/dev/null && [ -z "$REBUILD_TAGS" ]
//			then
//				say " - skipping tag '${TAG}' (already synced)"
//				continue
//			fi
//
//			if [ -n "$VERBOSE" ];
//			then
//				echo "${DEBUG} git branch | grep \"${LOCAL_TAG}$\" >/dev/null && [ -z \"${REBUILD_TAGS}\" ]"
//			fi
//
//			say " - syncing tag '${TAG}'"
//			say " - deleting '${LOCAL_TAG}'"
//			git branch -D "$LOCAL_TAG" >/dev/null 2>&1
//
//			if [ -n "$VERBOSE" ];
//			then
//				echo "${DEBUG} git branch -D \"${LOCAL_TAG}\" >/dev/null 2>&1"
//			fi
//
//			say " - subtree split for '${TAG}'"
//			git subtree split -q --annotate="${ANNOTATE}" --prefix="$SUBPATH" --branch="$LOCAL_TAG" "$TAG" >/dev/null
//
//			if [ -n "$VERBOSE" ];
//			then
//				echo "${DEBUG} git subtree split -q --annotate=\"${ANNOTATE}\" --prefix=\"$SUBPATH\" --branch=\"$LOCAL_TAG\" \"$TAG\" >/dev/null"
//			fi
//
//			say " - subtree split for '${TAG}' [DONE]"
//			if [ $? -eq 0 ]
//			then
//				PUSH_CMD="git push -q ${DRY_RUN} --force ${REMOTE_NAME} ${LOCAL_TAG}:refs/tags/${TAG}"
//
//				if [ -n "$VERBOSE" ];
//				then
//					echo "${DEBUG} PUSH_CMD=\"${PUSH_CMD}\""
//				fi
//
//				if [ -n "$DRY_RUN" ]
//				then
//					echo \# $PUSH_CMD
//					$PUSH_CMD
//				else
//					$PUSH_CMD
//				fi
//			fi
//		done
}

func (r *Repo)syncHeads(){
	println("Syncing heads ")
	//  Todo: convert from sh to go
//	if [ -n "$VERBOSE" ];
//			then
//				echo "${DEBUG} git show-ref --quiet --verify -- \"refs/remotes/origin/${HEAD}\""
//			fi
//
//			if ! git show-ref --quiet --verify -- "refs/remotes/origin/${HEAD}"
//			then
//				say " - skipping head '${HEAD}' (does not exist)"
//				continue
//			fi
//			LOCAL_BRANCH="${REMOTE_NAME}-branch-${HEAD}"
//
//			if [ -n "$VERBOSE" ];
//			then
//				echo "${DEBUG} LOCAL_BRANCH=\"${LOCAL_BRANCH}\""
//			fi
//
//			say " - syncing branch '${HEAD}'"
//			git checkout master >/dev/null 2>&1
//			git branch -D "$LOCAL_BRANCH" >/dev/null 2>&1
//			git branch -D "${LOCAL_BRANCH}-checkout" >/dev/null 2>&1
//			git checkout -b "${LOCAL_BRANCH}-checkout" "origin/${HEAD}" >/dev/null 2>&1
//			git subtree split -q --prefix="$SUBPATH" --branch="$LOCAL_BRANCH" "origin/${HEAD}" >/dev/null
//
//			if [ -n "$VERBOSE" ];
//			then
//				echo "${DEBUG} git checkout master >/dev/null 2>&1"
//				echo "${DEBUG} git branch -D \"$LOCAL_BRANCH\" >/dev/null 2>&1"
//				echo "${DEBUG} git branch -D \"${LOCAL_BRANCH}-checkout\" >/dev/null 2>&1"
//				echo "${DEBUG} git checkout -b \"${LOCAL_BRANCH}-checkout\" \"origin/${HEAD}\" >/dev/null 2>&1"
//				echo "${DEBUG} git subtree split -q --prefix=\"$SUBPATH\" --branch=\"$LOCAL_BRANCH\" \"origin/${HEAD}\" >/dev/null"
//			fi
//
//			if [ $? -eq 0 ]
//			then
//				PUSH_CMD="git push -q ${DRY_RUN} --force $REMOTE_NAME ${LOCAL_BRANCH}:${HEAD}"
//
//				if [ -n "$VERBOSE" ];
//				then
//					echo "${DEBUG} $PUSH_CMD"
//				fi
//
//				if [ -n "$DRY_RUN" ]
//				then
//					echo \# $PUSH_CMD
//					$PUSH_CMD
//				else
//					$PUSH_CMD
//				fi
//			fi
}

func (r *Repo)sync(){
	println("Starting Sync ")
	r.syncHeads();
	r.syncTags();
	println("Sync Completed ")
}

func (r *Repo)update(){
	println("Updating subsplit from origin ")
	//  Todo: convert from sh to go
//	subsplit_require_work_dir
//
//
//	git fetch -q origin
//	git fetch -q -t origin
//	git checkout master
//	git reset --hard origin/master
//
//	if [ -n "$VERBOSE" ];
//	then
//		echo "${DEBUG} git fetch -q origin"
//		echo "${DEBUG} git fetch -q -t origin"
//		echo "${DEBUG} git checkout master"
//		echo "${DEBUG} git reset --hard origin/master"
//	fi
}


func (r *myRegexp) FindStringSubmatchMap(s string) map[string]string {
	captures := make(map[string]string)

	match := r.FindStringSubmatch(s)
	if match == nil {
		return captures
	}

	for i, name := range r.SubexpNames() {
		// Ignore the whole regexp match and unnamed groups
		if i == 0 || name == "" {
			continue
		}

		captures[name] = match[i]

	}
	return captures
}
