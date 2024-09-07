package domain

import (
	"encoding/base64"
	"os"
	"streamapi/datasource"
	"streamapi/utils/errors"
	"streamapi/utils/logger"
	"strings"
)

var (
	queryInsert = "insert into streamdata(hash,filename,stream_link,download_link,has_thumb) values(?,?,?,?,?)"
)

func (data *Data) AddDataToDB() *errors.RestErr {

	stmt, err := datasource.Client.Prepare(queryInsert)
	if err != nil {
		logger.Debug.Println(err)
		err := stmt.Close()
		if err != nil {
			return errors.NewBadRequestError("Database insert error and row close error")
		}
		return errors.NewBadRequestError("Database insert error")
	}

	_, err = stmt.Exec(data.Hash, data.Filename, data.StreamLink, data.DownloadLink, data.HasThumb)
	if err != nil {
		err := stmt.Close()
		if err != nil {
			return errors.NewBadRequestError("Database exec error and row close error")
		}
		return errors.NewBadRequestError("Database exec error")
	}

	err = stmt.Close()
	if err != nil {
		return errors.NewBadRequestError("Close error")
	}
	return nil

}

func generateThumbUrlValue(data *Data) {
	if data.HasThumb != 0 {
		thumbBaseDir := os.Getenv("THUMB_BASE_DIR")
		filePath := thumbBaseDir + "/" + data.Hash + ".jpg"
		if _, err := os.Stat(filePath); err == nil {
			bytes, err := os.ReadFile(filePath)
			if err != nil {
				return
			}
			base64Encoding := "data:image/jpeg;base64,"
			base64Encoding += base64.StdEncoding.EncodeToString(bytes)
			data.ThumbUrl = base64Encoding
		}
	}
}

func GetFileDataFromDB(hash string, ipAddress string, action string) (*Data, *errors.RestErr) {

	var data Data
	var query string
	var actionOperator string
	var orderingMethod string

	if action != "" && ipAddress != "" {
		if action == "next" {
			actionOperator = ">"
			orderingMethod = "asc"
		} else {
			actionOperator = "<"
			orderingMethod = "desc"
		}
		currentMgsId := hash[6:]
		searchIpAddress := "'http://" + ipAddress + "/%'"
		query = "select id,hash,filename,download_link,stream_link,has_thumb,created_at from streamdata where stream_link like " + searchIpAddress + "and SUBSTRING(hash,7)" + actionOperator + " " + currentMgsId + " and '" + hash + "' in (select hash from streamdata where stream_link like " + searchIpAddress + ") order by id " + orderingMethod + " limit 1 "
	} else {
		query = "select id,hash,filename,download_link,stream_link,has_thumb,created_at from streamdata where hash='" + hash + "'"
	}

	rows, getErr := datasource.Client.Query(query)

	if getErr != nil {
		logger.Debug.Println(getErr)
		err := rows.Close()
		if err != nil {
			return nil, errors.NewBadRequestError("Database query error and row close error")
		}
		return nil, errors.NewBadRequestError("Database query error")
	}

	if rows.Next() {
		err := rows.Scan(
			&data.ID,
			&data.Hash,
			&data.Filename,
			&data.DownloadLink,
			&data.StreamLink,
			&data.HasThumb,
			&data.CreatedAt,
		)
		if err != nil {
			logger.Debug.Println(err)
			err := rows.Close()
			if err != nil {
				return nil, errors.NewBadRequestError("Fetch error and row close error")
			}
			return nil, errors.NewBadRequestError("Fetch error")
		}
	}
	err := rows.Close()
	if err != nil {
		logger.Debug.Println("Error in closing row")
		return nil, errors.NewBadRequestError("Close Error")
	}
	generateThumbUrlValue(&data)

	return &data, nil

}

func SearchDataFromDB(query string) (*[]Data, *errors.RestErr) {

	var data []Data
	queryP := strings.ReplaceAll(query, "-", "_")
	queryP = strings.ReplaceAll(queryP, " ", "_")
	rows, getErr := datasource.Client.Query(
		"select id,hash,filename,download_link,stream_link,has_thumb,created_at from streamdata where filename like '%" + query + "%' or filename like '%" + queryP + "%'",
	)

	if getErr != nil {
		logger.Debug.Println(getErr)
		err := rows.Close()
		if err != nil {
			return nil, errors.NewBadRequestError("Database query error and row close error")
		}
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
			&rowData.HasThumb,
			&rowData.CreatedAt,
		)
		if err != nil {
			err := rows.Close()
			if err != nil {
				return nil, errors.NewBadRequestError("Fetch error and row close error")
			}
			return nil, errors.NewBadRequestError("Fetch error")
		}
		generateThumbUrlValue(&rowData)
		data = append(data, rowData)
	}
	err := rows.Close()
	if err != nil {
		logger.Debug.Println("Error in closing row")
		return nil, errors.NewBadRequestError("Close Error")
	}

	return &data, nil

}

func DeleteFileDataFromDB(hash string) *errors.RestErr {

	queryDelete := "delete from streamdata where hash='" + hash + "'"

	stmt, err := datasource.Client.Prepare(queryDelete)
	if err != nil {
		logger.Debug.Println(err)
		return errors.NewBadRequestError("Database delete smt error")
	}
	defer stmt.Close()

	result, err := stmt.Exec()
	if err != nil {
		logger.Debug.Println(err)
		return errors.NewBadRequestError("Database delete exec error")
	}

	rowsEff, err := result.RowsAffected()
	if err != nil {
		logger.Debug.Println(err)
		return errors.NewBadRequestError("Database delete exec error")
	}

	if rowsEff == 0 {
		logger.Debug.Println("Rows not affected for delete query")
		return errors.NewBadRequestError("Delete query didn't affect database")
	}

	return nil
}
