package store

import (
	"fmt"

	"github.com/kusubooru/shimmie"
)

const (
	configKeyTitle          = "title"
	configKeyDescription    = "site_description"
	configKeyKeywords       = "site_keywords"
	configKeyAnalyticsID    = "google_analytics_id"
	configKeyAnalyticsIDOld = "ga_profile_id"
)

func (db *datastore) GetCommon() (*shimmie.Common, error) {
	keys := []string{
		configKeyTitle,
		configKeyDescription,
		configKeyKeywords,
		configKeyAnalyticsID,
		configKeyAnalyticsIDOld,
	}
	c, err := db.GetConfig(keys...)
	if err != nil {
		return nil, err
	}
	conf := shimmie.Common{
		Title:       c[configKeyTitle],
		Description: c[configKeyDescription],
		Keywords:    c[configKeyKeywords],
		AnalyticsID: c[configKeyAnalyticsID],
	}
	if conf.AnalyticsID == "" {
		conf.AnalyticsID = c[configKeyAnalyticsIDOld]
	}
	return &conf, nil
}

func (db *datastore) GetConfig(keys ...string) (map[string]string, error) {
	query := fmt.Sprint(configGetQuery)
	if len(keys) != 0 {
		query = fmt.Sprintf("%vWHERE\n", query)
		for _, k := range keys {
			query = fmt.Sprintf("%v  name = '%v' OR\n", query, k)
		}
	}
	query = query[:len(query)-4] // chomp last ' OR\n'

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	var (
		config = struct {
			Name  string
			Value string
		}{}
		m = map[string]string{}
	)
	for rows.Next() {
		err := rows.Scan(&config.Name, &config.Value)
		if err != nil {
			return nil, err
		}
		m[config.Name] = config.Value
	}
	return m, nil
}

const configGetQuery = `
SELECT *
FROM config
`
