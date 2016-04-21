package datastore

import "fmt"

func (db *datastore) GetConfig(keys ...string) (map[string]string, error) {
	query := fmt.Sprintf(configGetQuery)
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
