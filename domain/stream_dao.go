package domain

import (
	"encoding/base64"
	"fmt"
	"os"
	"streamapi/datasource"
	"streamapi/structs"
	"streamapi/utils/cache"
	"streamapi/utils/errors"
	"streamapi/utils/logger"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	queryInsert = "insert into streamdata(hash,filename,stream_link,download_link,has_thumb) values(?,?,?,?,?)"
)

type Data structs.Data

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

func generateThumbUrlValue(data *structs.Data) {
	if data.HasThumb != 0 {
		thumbBaseDir := os.Getenv("THUMBNAILS_DIR")
		if thumbBaseDir == ""{
			thumbBaseDir = os.Getenv("THUMB_BASE_DIR")
		}
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

func checkIfFileExistsInDB(hash string) (bool, *errors.RestErr) {
	existQuery := fmt.Sprintf("select * from streamfrequent where hash='%s'", hash)
	// Check if already exists
	rows, getErr := datasource.Client.Query(existQuery)
	if getErr != nil {
		logger.Debug.Println(getErr)
		err := rows.Close()
		if err != nil {
			return false, errors.NewBadRequestError(
				"Frequent Database query error and row close error",
			)
		}
		return false, errors.NewBadRequestError("Frequent Database query error")
	}

	if rows.Next() {
		// Already exists, no need to insert again
		err := rows.Close()
		if err != nil {
			return false, errors.NewBadRequestError("Frequent Database row close error")
		}
		return true, nil
	}
	err := rows.Close()
	if err != nil {
		return false, errors.NewBadRequestError("Frequent Database row close error")
	}
	return false, nil
}

func updateSearchedFileInDB(hash string) *errors.RestErr {
	//update the search_frequency counter if needed
	updateQuery := fmt.Sprintf(
		"update streamfrequent set search_frequency = search_frequency + 1 where hash='%s'",
		hash,
	)
	stmt, err := datasource.Client.Prepare(updateQuery)
	if err != nil {
		logger.Debug.Println(err)
		err := stmt.Close()
		if err != nil {
			return errors.NewBadRequestError("Frequent Database update error and row close error")
		}
		return errors.NewBadRequestError("Frequent Database update error")
	}

	_, err = stmt.Exec()
	if err != nil {
		err := stmt.Close()
		if err != nil {
			return errors.NewBadRequestError("Frequent Database exec error and row close error")
		}
		return errors.NewBadRequestError("Frequent Database exec error")
	}

	err = stmt.Close()
	if err != nil {
		return errors.NewBadRequestError("Close error")
	}
	return nil

}

func addSearchedFileToDB(data *structs.Data) *errors.RestErr {
	// Check if already exists
	exists, gErr := checkIfFileExistsInDB(data.Hash)
	if gErr != nil {
		return gErr
	}
	if exists {
		// If exists, update the search frequency counter
		return updateSearchedFileInDB(data.Hash)
	}

	// If not exists, insert new record
	insertQuery := "insert into streamfrequent(hash) values (?)"
	stmt, err := datasource.Client.Prepare(insertQuery)
	if err != nil {
		logger.Debug.Println(err)
		err := stmt.Close()
		if err != nil {
			return errors.NewBadRequestError("Frequent Database insert error and row close error")
		}
		return errors.NewBadRequestError("Frequent Database insert error")
	}

	_, err = stmt.Exec(data.Hash)
	if err != nil {
		err := stmt.Close()
		if err != nil {
			return errors.NewBadRequestError("Frequent Database exec error and row close error")
		}
		return errors.NewBadRequestError("Frequent Database exec error")
	}

	err = stmt.Close()
	if err != nil {
		return errors.NewBadRequestError("Close error")
	}
	return nil
}

func GetFileDataFromDB(
	c *gin.Context,
	queryParams structs.QueryParams,
) (*structs.Data, *errors.RestErr, *structs.CacheHeaders) {

	var cacheHeaders structs.CacheHeaders
	cacheHeaders.CacheHit = false

	hash := queryParams.Hash
	ipAddress := queryParams.IPAddress
	action := queryParams.Action


	// Check for cached data first
	ctx := c.Request.Context()
	if action != "" {
		cacheHeaders.CacheKey = cache.CreateNavCacheKey(hash, action)
	} else {
		cacheHeaders.CacheKey = cache.CreateRecordCacheKey(hash)
	}

	cachedData, hit := cache.GetSingleDataFromCache(ctx, cacheHeaders.CacheKey)
	// Cache HIT: Return Cached Data
	if hit {
		cacheHeaders.CacheHit = true
		return cachedData, nil, &cacheHeaders
	}

	// Cache MISS : Move on to regular flow
	var data structs.Data
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
			return nil, errors.NewBadRequestError("Database query error and row close error"), &cacheHeaders
		}
		return nil, errors.NewBadRequestError("Database query error"), &cacheHeaders
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
				return nil, errors.NewBadRequestError("Fetch error and row close error"), &cacheHeaders
			}
			return nil, errors.NewBadRequestError("Fetch error"), &cacheHeaders
		}
	}
	err := rows.Close()
	if err != nil {
		logger.Debug.Println("Error in closing row")
		return nil, errors.NewBadRequestError("Close Error"), &cacheHeaders
	}
	generateThumbUrlValue(&data)

	// Set data to redis
	cache.SetDataToCache(ctx, cacheHeaders.CacheKey, data)

	return &data, nil, &cacheHeaders

}

func GetRecentFilesFromDB(page int, limit int, sort string) (*[]structs.Data, *errors.RestErr) {

	var data []structs.Data
	offset := (page - 1) * limit
	sortOption := ""

	switch sort {
	case "searched":
		sortOption += "search_frequency desc"
	case "oldest":
		sortOption += "sf.updated_at asc"
	default:
		sortOption += "sf.updated_at desc"
	}

	rows, getErr := datasource.Client.Query(fmt.Sprintf(
		"select s.hash,filename,download_link,stream_link,has_thumb,sf.created_at,sf.updated_at,search_frequency from streamdata s inner join streamfrequent sf where s.hash=sf.hash order by %s limit %d,%d",
		sortOption,
		offset,
		limit,
	),
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
		var rowData structs.Data
		err := rows.Scan(
			&rowData.Hash,
			&rowData.Filename,
			&rowData.DownloadLink,
			&rowData.StreamLink,
			&rowData.HasThumb,
			&rowData.CreatedAt,
			&rowData.UpdatedAt,
			&rowData.SearchFrequency,
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

func SearchDataFromDB(query string) (*[]structs.Data, *errors.RestErr) {

	var data []structs.Data
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
		var rowData structs.Data
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
		addSearchedFileToDB(&rowData)
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
