package frontend

import (
	"context"
	"fmt"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"

	"go.uber.org/zap"

	"core/utils"
)

const secretKey = "figrin-dan-and-the-modal-nodes"

func (fs *FrontendServer) NewHook(repo *github.Repository, client *github.Client) error {
	boltDB, err := bolt.Open(utils.Config.BoltDBPath, 0644, nil)
	if err != nil {
		utils.AppLog.Error("failed opening bolt", zap.Error(err))
		return err
	}
	database := BoltDB{DB: boltDB}
	if check, err := fs.hookExists(repo, client, database); check {
		return errors.Wrap(err, "hook exists")
	}
	name := *repo.Name
	owner := *repo.Owner.Login
	url := "http://heupr.io:8020/hook"
	hook, _, err := client.Repositories.CreateHook(context.Background(), owner, name, &github.Hook{
		Name:   github.String("web"),
		Events: []string{"issues", "pull_request"},
		Active: github.Bool(true),
		Config: map[string]interface{}{
			"url":          url,
			"secret":       secretKey,
			"content_type": "json",
			"insecure_ssl": false,
		},
	})
	if err != nil {
		return errors.Wrap(err, "adding new hook to repo")
	}
	if err = database.Store("hook", *repo.ID, []byte(strconv.Itoa(*hook.ID))); err != nil {
		return errors.Wrap(err, "error storing hook info")
	}
	boltDB.Close()
	return nil
}

func (fs *FrontendServer) hookExists(repo *github.Repository, client *github.Client, database BoltDB) (bool, error) {
	name, owner := "", ""
	if repo.Name != nil && repo.Owner.Login != nil {
		name = *repo.Name
		owner = *repo.Owner.Login
	}
	hook, err := database.Retrieve("hook", *repo.ID)
	if err != nil {
		return false, errors.Wrap(err, "error retrieving hook info")
	}

	hookID, err := strconv.Atoi(string(hook))
	if err != nil {
		return false, errors.Wrap(err, "failed hook string conversion")
	}
	_, _, err = client.Repositories.GetHook(context.Background(), owner, name, hookID)
	if err != nil {
		return false, errors.Wrap(err, "error getting GitHub hook info")
	}
	return true, fmt.Errorf("hook previously established for repo %v", name)
}
