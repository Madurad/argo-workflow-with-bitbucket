package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var Version = "0.1.0"

func main() {
	log.Println("Argo webhook started with version", Version)
	// For prod only
	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
	log.SetFlags(log.Lshortfile)
	r := gin.Default()

	// Whitelist repo prefix that will allow to process
	whRepo := os.Getenv("WH_REPO")
	if whRepo == "" {
		whRepo = "TheYkk"
	}

	r.POST("/", func(c *gin.Context) {
		c.JSON(202, gin.H{})

		var bc Bitbucket

		err := c.BindJSON(&bc)
		if err != nil {
			log.Fatal(err)
		}

		// Get repo url
		// Check repo is in white list
		if !strings.HasPrefix(strings.ToLower(bc.Repository.FullName), strings.ToLower(whRepo)) {
			return
		}

		// Git repo details
		gitRevision := bc.Push.Changes[0].New.Name
		gitRepoName := strings.Split(bc.Repository.FullName, "/")[1]
		fullGitRepo := "git@bitbucket.org:" + bc.Repository.FullName + ".git"

		timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		argoFilename := "argo" + timestamp + ".yml"

		input, _ := ioutil.ReadFile("argo.yml")
		temp := bytes.Replace(input, []byte("<git_repo_name>"), []byte(gitRepoName), -1)
		temp2 := bytes.Replace(temp, []byte("<git_repo_full>"), []byte(fullGitRepo), -1)
		output := bytes.Replace(temp2, []byte("<git_revision>"), []byte(gitRevision), -1)
		_ = ioutil.WriteFile(argoFilename, output, 0666)

		commandOutput, err := exec.Command("sh", "-c", "./argo submit "+argoFilename).CombinedOutput()
		if err != nil {
			fmt.Printf("Accepted webhook request, did NOT start Argo workflow: git_repo=%q,git_revision=%q, because of: %q\n", bc.Repository.FullName, gitRevision, string(err.Error()))
		} else {
			fmt.Printf("Accepted webhook request, started Argo workflow: git_repo=%q,git_revision=%q, with message: %q\n", bc.Repository.FullName, gitRevision, string(commandOutput))
		}

	})
	_ = r.Run(":3000")
}
