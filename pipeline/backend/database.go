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
	RepoID         int64
	AppID          int
	InstallationID int
}

func (m *MemSQL) Open() {
	mysql, err := sql.Open("mysql", "root@/heupr?interpolateParams=true&parseTime=true")
	if err != nil {
		panic(err.Error()) // TODO: Proper error handling.
	}
	m.db = mysql
}

func (m *MemSQL) Close() {
	m.db.Close()
}

func (m *MemSQL) ReadIntegrations() ([]Integration, error) {
	integrations := []Integration{}
	results, err := m.db.Query("select repo_id, app_id, installation_id from integrations")
	if err != nil {
		return nil, err
	}

	defer results.Close()
	for results.Next() {
		integration := Integration{}
		err := results.Scan(&integration.RepoID, &integration.AppID, &integration.InstallationID)
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

func (m *MemSQL) ReadIntegrationByRepoID(repoID int64) (*Integration, error) {
	integration := new(Integration)
	err := m.db.QueryRow("select repo_id, app_id, installation_id from integrations where repo_id = ?", repoID).Scan(&integration.RepoID, &integration.AppID, &integration.InstallationID)
	if err != nil {
		utils.AppLog.Error("ReadIntegrationByRepoId Database Read Failure", zap.Error(err))
		return nil, err
	}
	return integration, nil
}
