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
	"path/filepath"
	"strconv"
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
			Name: "Origin, o",
			Value: "--Origin",
			Usage: "Origin of the repository",
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
			Usage:     "Initialize a subsplit from Origin",
			Action: func(c *cli.Context) {
				checkRequirments()


				fmt.Printf("Initialize Task: Starting")

				cmd := fmt.Sprintf("git clone -q %s ", c.Args().First())
				_, err := exec.Command(SHELL, SHELL_ARG_C, cmd).Output()

				absolute_repo := strings.Split(c.Args().First(), "/")

				dir := strings.Replace(absolute_repo[len(absolute_repo) - 1], ".git", "", -1)

				f, err := os.Create(dir + "/.subsplit")
				f.Close()

				if err != nil {
					fmt.Printf("Initialize Task: Failed : %s", err)

				}else {
					file, err := os.Open(dir + "/.gitignore")

					if err != nil {
						fmt.Println(err)

						f, err := os.Create(dir + "/.gitignore")

						if err == nil {

							n, err := io.WriteString(f, dir + "/.subsplit")

							if err != nil {
								fmt.Println(n, err)
							}
						}else {
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
				checkRequirments()
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

						cmd := fmt.Sprintf("git remote add r.REMOTE_NAME r.REMOTE_URL")
						_, err := exec.Command(SHELL, SHELL_ARG_C, cmd).Output()

						if (err != nil) {
							r.sync("origin", "master", "sample annotation", true);
						}


					}

					println("Publish task: Complete")
				}else {
					println("Error: no .subsplit found, has the repo been initalized with git-subsplit init ?")
				}
			},
		}, {
			Name:      "Update",
			Aliases:     []string{"-u"},
			Usage:     "Update ",
			Action: func(c *cli.Context) {
				println("Updating subsplit from Origin ")
			},
		},

	}

	app.Run(os.Args)
}


func getHeads(data string) []Head {
	println("Querying heads from Origin")
	cmd := fmt.Sprintf("git ls-remote Origin")
	output, _ := exec.Command(SHELL, SHELL_ARG_C, cmd).Output()
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
	println("Querying tags from Origin")
	cmd := fmt.Sprintf("git ls-remote Origin")
	output, _ := exec.Command(SHELL, SHELL_ARG_C, cmd).Output()
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
	byteArray, _ := exec.Command(SHELL, SHELL_ARG_C, cmd).Output()
	result := string(byteArray[:])
	return strings.Replace(result, "\n", "", -1)
}

func (r *Repo)syncTags(DRY_RUN bool, ANNOTATE string) {
	println("Syncing Tags ")


	for _, Tag := range r.TAGS {
		cmd := fmt.Sprintf("git show-ref --quiet --verify -- \"refs/tags/" + Tag.tag + "\"")
		output, _ := exec.Command(SHELL, SHELL_ARG_C, cmd).Output()
		fmt.Println(string(output))
		if (output != nil) {
			LOCAL_TAG := r.REMOTE_NAME + "-tag-" + Tag.tag

			cmd = fmt.Sprintf("git branch")
			output, _ = exec.Command(SHELL, SHELL_ARG_C, cmd).Output()
			regex := regexp.MustCompile(LOCAL_TAG)
			matches := regex.FindStringSubmatch(string(output))

			if (matches != nil) {
				println("- skipping tag " + LOCAL_TAG + " (already synced)")
			}else {
				println("Syncing Tag: " + Tag.tag)
				println("Deleting Tag: " + LOCAL_TAG)
				cmd = fmt.Sprintf("git branch -D \"" + LOCAL_TAG + " \" ")
				output, _ = exec.Command(SHELL, SHELL_ARG_C, cmd).Output()
				println("Syncing Tag " + Tag.tag + ": Complete ")

				cmd = fmt.Sprintf("git subtree split -q --annotate=\"" + ANNOTATE + "\" --prefix=\"" + r.SUBPATH + "\" --branch=\"" + LOCAL_TAG + "\" \"" + Tag.tag + "\"")

				test := ""
			if (DRY_RUN == true) {
				test = "--dry-run"
			}

				output, _ = exec.Command(SHELL, SHELL_ARG_C, cmd).Output()
				cmd = fmt.Sprintf("git push -q " + test + " --force " + r.REMOTE_NAME + " " + LOCAL_TAG + ":" + Tag.tag)

				output, _ = exec.Command(SHELL, SHELL_ARG_C, cmd).Output()
			}
			println(output)
		}else {
			println(" - skipping Tag: " + Tag.tag + " (does not exist) ")
		}
	}

}


func (r *Repo)syncHeads(Origin string, DRY_RUN bool) {
	println("Syncing heads: Started ")

	for _, head := range r.HEADS {
		cmd := fmt.Sprintf("git show-ref --quiet --verify -- \"refs/remotes/" + Origin + "/" + head.head + "\"")
		output, _ := exec.Command(SHELL, SHELL_ARG_C, cmd).Output()
		fmt.Println(string(output))
		if (output != nil) {
			println("Syncing Branch: " + head.head)
			cmd := fmt.Sprintf("git checkout " + head.head + " >/dev/null 2>&1")
			output, _ := exec.Command(SHELL, SHELL_ARG_C, cmd).Output()

			LOCAL_BRANCH := r.REMOTE_NAME + "-branch-" + head.head

			cmd = fmt.Sprintf("git branch -D \"" + LOCAL_BRANCH + "\" >/dev/null 2>&1")
			output, _ = exec.Command(SHELL, SHELL_ARG_C, cmd).Output()

			cmd = fmt.Sprintf("git branch -D \"" + LOCAL_BRANCH + " - checkout \"")
			output, _ = exec.Command(SHELL, SHELL_ARG_C, cmd).Output()

			cmd = fmt.Sprintf("git checkout -b \"" + LOCAL_BRANCH + "-checkout\" \"" + Origin + "/" + head.head + "\" ")
			output, _ = exec.Command(SHELL, SHELL_ARG_C, cmd).Output()

			cmd = fmt.Sprintf("git subtree split -q --prefix=\"" + r.SUBPATH + "\" --branch=\"" + LOCAL_BRANCH + "\" \"" + Origin + "/" + head.head + "\" >/dev/null")
			output, _ = exec.Command(SHELL, SHELL_ARG_C, cmd).Output()//
			test := ""
			if (DRY_RUN == true) {
				test = "--dry-run"
			}
			cmd = fmt.Sprintf("git push -q " + test + " --force " + r.REMOTE_NAME + " " + LOCAL_BRANCH + ":" + head.head)
			output, _ = exec.Command(SHELL, SHELL_ARG_C, cmd).Output()
			println(" Syncing Branch " + head.head + ": Complete")
			println(output)

		}else {
			println(" - skipping head: " + head.head + " (does not exist) ")
		}
	}
	println("Syncing heads: Completed")
}

func (r *Repo)sync(Origin string, Branch string, Annotate string, Dry_Run bool) {
	println("Syncing Task: Started")
	r.syncHeads(Origin, Dry_Run);
	r.syncTags(Origin, Annotate);
	println("Syncing Task: Completed")
}

func (r *Repo)update(Origin string, Branch string) {
	println("Updating subsplit from Origin: Starting")
	cmd := fmt.Sprintf("git fetch -q " + Origin)
	output, _ := exec.Command(SHELL, SHELL_ARG_C, cmd).Output()
	cmd = fmt.Sprintf("git fetch -q -t " + Origin)
	output, _ = exec.Command(SHELL, SHELL_ARG_C, cmd).Output()
	cmd = fmt.Sprintf("git checkout master")
	output, _ = exec.Command(SHELL, SHELL_ARG_C, cmd).Output()
	cmd = fmt.Sprintf("git reset --hard " + Origin + "/" + Branch)
	output, _ = exec.Command(SHELL, SHELL_ARG_C, cmd).Output()
	println("Updating subsplit from Origin: Completed")
	println(output)
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

func checkRequirments() {
	valid := false
	cmd := fmt.Sprintf("git version")
	output, _ := exec.Command(SHELL, SHELL_ARG_C, cmd).Output()

	result := strings.Split(output, " ")
	f, _ := strconv.ParseFloat(result[len(result) - 1], 64)
	if ( f < "1.7.11" ) {
		println("Git subplit needs git subtree; upgrade git to >=1.7.11")
		os.Exit(1)
	}else{
		valid = true
	}

	return valid
}

