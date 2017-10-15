package backend

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"go.uber.org/zap"

	"core/utils"
)

type MemSQL struct {
	db *sql.DB
}

type Integration struct {
	RepoId         int
	AppId          int
	InstallationId int
}

func (m *MemSQL) Open() {
	mysql, err := sql.Open("mysql", "root@/heupr?interpolateParams=true")
	if err != nil {
		panic(err.Error()) // TODO: Proper error handling.
	}
	m.db = mysql
}

func (m *MemSQL) Close() {
	m.db.Close()
}

func (d *MemSQL) ReadIntegrations() ([]Integration, error) {
	integrations := []Integration{}
	results, err := d.db.Query("select repo_id, app_id, installation_id from integrations")
	if err != nil {
		return nil, err
	}

	defer results.Close()
	for results.Next() {
		integration := Integration{}
		err := results.Scan(&integration.RepoId, &integration.AppId, &integration.InstallationId)
		if err != nil {
			return nil, err
		}
		integrations = append(integrations, integration)
		err = results.Err()
		if err != nil {
			return nil, err
		}
	}
	return integrations, nil
}

func (d *MemSQL) ReadIntegrationByRepoId(repoId int) (*Integration, error) {
	integration := new(Integration)
	err := d.db.QueryRow("select repo_id, app_id, installation_id from integrations where repo_id = ?", repoId).Scan(&integration.RepoId, &integration.AppId, &integration.InstallationId)
	if err != nil {
		utils.AppLog.Error("ReadIntegrationByRepoId Database Read Failure", zap.Error(err))
		return nil, err
	}
	return integration, nil
}
