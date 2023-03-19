package domain

import (
	"streamapi/datasource"
	"streamapi/utils/errors"
	"streamapi/utils/logger"
)

var (
	queryInsert = "insert into streamdata(hash,filename,stream_link,download_link) values(?,?,?,?)"
)

func (data *Data) AddDataToDB() *errors.RestErr {

	stmt, err := datasource.Client.Prepare(queryInsert)
	if err != nil {
		logger.Debug.Println(err)
		return errors.NewBadRequestError("Database insert error")
	}

	_, err = stmt.Exec(data.Hash, data.Filename, data.StreamLink, data.DownloadLink)
	if err != nil {
		return errors.NewBadRequestError("Database exec error")
	}

	return nil

}

func GetFileDataFromDB(hash string) (*Data, *errors.RestErr) {

	var data Data

	rows, getErr := datasource.Client.Query("select * from streamdata where hash='" + hash + "'")

	if getErr != nil {
		logger.Debug.Println(getErr)
		return nil, errors.NewBadRequestError("Database query error")
	}

	if rows.Next() {
		err := rows.Scan(
			&data.ID,
			&data.Hash,
			&data.Filename,
			&data.DownloadLink,
			&data.StreamLink,
			&data.CreatedAt,
		)
		if err != nil {
			logger.Debug.Println(err)
			return nil, errors.NewBadRequestError("Fetch error")
		}
	}

	return &data, nil

}

func SearchDataFromDB(query string) (*[]Data, *errors.RestErr) {

	var data []Data

	rows, getErr := datasource.Client.Query(
		"select * from streamdata where filename like '%" + query + "%'",
	)

	if getErr != nil {
		logger.Debug.Println(getErr)
		return nil, errors.NewBadRequestError("Database query error")
	}

	for rows.Next() {
		var rowData Data
		err := rows.Scan(
			&rowData.ID,
			&rowData.Hash,
			&rowData.Filename,
			&rowData.DownloadLink,
			&rowData.StreamLink,
			&rowData.CreatedAt,
		)
		if err != nil {
			return nil, errors.NewBadRequestError("Fetch error")
		}
		data = append(data, rowData)
	}

	return &data, nil

}
