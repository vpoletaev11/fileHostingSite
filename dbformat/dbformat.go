package dbformat

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

const (
	filenameLen    = 20
	descriptionLen = 35
)

const (
	getUserTimezone = "SELECT timezone FROM users WHERE username = ?;"
)

// FileInfo contains formatted file info from MySQL database
type FileInfo struct {
	Label        string
	DownloadLink string
	FilesizeMb   string
	Description  string
	Owner        string
	Category     string
	UploadDate   string
	Rating       int

	LabelComment         string
	FilesizeBytesComment string
	DescriptionComment   string
}

// DownloadFileInfo contains formatted file info getted from MySQL database
type DownloadFileInfo struct {
	DownloadLink string
	Label        string
	FilesizeMB   string
	Description  string
	Owner        string
	Category     string
	UploadDate   string
	Rating       int
}

func userLocalTime(db *sql.DB, globalTime time.Time, username string) (time.Time, error) {
	userTimezone := ""
	err := db.QueryRow(getUserTimezone, username).Scan(&userTimezone)
	if err != nil {
		return time.Time{}, err
	}

	location, err := time.LoadLocation(userTimezone)
	if err != nil {
		return time.Time{}, err

	}

	globalTime = globalTime.In(location)
	return globalTime, nil
}

// FormatedDownloadFileInfo returns fromatted download file info
func FormatedDownloadFileInfo(username string, db *sql.DB, query, argument string) (DownloadFileInfo, error) {
	fi := DownloadFileInfo{}
	var uploadDateTime time.Time
	id := 0
	filesizeBytes := 0
	err := db.QueryRow(query, argument).Scan(
		&id,
		&fi.Label,
		&filesizeBytes,
		&fi.Description,
		&fi.Owner,
		&fi.Category,
		&uploadDateTime,
		&fi.Rating,
	)
	if err != nil {
		return DownloadFileInfo{}, err
	}
	fi.DownloadLink = "/files/" + strconv.Itoa(id)
	userTime, err := userLocalTime(db, uploadDateTime, username)
	if err != nil {
		return DownloadFileInfo{}, err
	}
	fi.UploadDate = userTime.Format("2006-01-02 15:04:05")
	fi.FilesizeMB = fmt.Sprintf("%.6f", float64(filesizeBytes)/1024/1024) + " MB"

	return fi, nil
}

// FormatedFilesInfo returns array of formatted file information
func FormatedFilesInfo(username string, db *sql.DB, query string, args ...interface{}) ([]FileInfo, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return []FileInfo{}, err
	}
	defer rows.Close()

	var fiTableCollection []FileInfo
	fiTable := new(FileInfo)

	id := 0
	var uploadDateTime time.Time
	for rows.Next() {
		err := rows.Scan(
			&id,
			&fiTable.LabelComment,
			&fiTable.FilesizeBytesComment,
			&fiTable.DescriptionComment,
			&fiTable.Owner,
			&fiTable.Category,
			&uploadDateTime,
			&fiTable.Rating,
		)
		if err != nil {
			return []FileInfo{}, err
		}
		userTime, err := userLocalTime(db, uploadDateTime, username)
		if err != nil {
			return []FileInfo{}, err
		}
		fiTable.UploadDate = userTime.Format("2006-01-02 15:04:05")

		if len(fiTable.LabelComment) > filenameLen {
			fiTable.Label = fiTable.LabelComment[:filenameLen] + "..."
		} else {
			fiTable.Label = fiTable.LabelComment
		}

		if len(fiTable.DescriptionComment) > descriptionLen {
			fiTable.Description = fiTable.DescriptionComment[:descriptionLen] + "..."
		} else {
			fiTable.Description = fiTable.DescriptionComment
		}

		fsBytes, err := strconv.Atoi(fiTable.FilesizeBytesComment)
		if err != nil {
			return []FileInfo{}, err
		}
		fiTable.FilesizeMb = fmt.Sprintf("%.4f", float64(fsBytes)/1024/1024) + " MB"
		fiTable.DownloadLink = "/download?id=" + strconv.Itoa(id)
		fiTable.FilesizeBytesComment = fiTable.FilesizeBytesComment + " Bytes"

		fiTableCollection = append(fiTableCollection, *fiTable)
	}
	return fiTableCollection, nil
}
