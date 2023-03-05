package db

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	_ "embed"
	"golang/jtunnel-client/data"
	"io"
	"log"
	_ "modernc.org/sqlite"
)

//go:embed tables.sql
var sqlFile []byte

type LocalDb struct {
	FilePath string
	db       *sql.DB
}

func NewLocalDb(file string) LocalDb {
	var d *sql.DB
	if file == "" {
		d, _ = sql.Open("sqlite", ":memory:")
	} else {
		log.Println("Saving all requests to file ", file)
		d, _ = sql.Open("sqlite", file)
	}

	d.Exec("BEGIN TRANSACTION")
	_, sqlExecErr := d.Exec(string(sqlFile))
	if sqlExecErr != nil {
		d.Exec("ROLLBACK")
		log.Fatal("Could not create tables", sqlExecErr.Error())
	}
	d.Exec("COMMIT")
	return LocalDb{
		FilePath: file,
		db:       d,
	}
}

func (d *LocalDb) Get(id int64) (data.TunnelData, error) {
	row := d.db.QueryRow("select ID,TUNNEL_ID,REQUEST_DATA,RESPONSE_DATA,LOCAL_PORT from ALL_REQUESTS where id=?", id)
	var tunnelData data.TunnelData
	err := row.Scan(
		&tunnelData.Id,
		&tunnelData.TunnelId,
		&tunnelData.RequestData,
		&tunnelData.ResponseData,
		&tunnelData.LocalPort,
	)
	tunnelData.RequestData = uncompress(tunnelData.RequestData)
	tunnelData.ResponseData = uncompress(tunnelData.ResponseData)
	return tunnelData, err
}

func (d *LocalDb) ListWithoutData(start int64, limit int) ([]data.TunnelData, error) {
	results, err := d.db.Query("select ID,TUNNEL_ID from ALL_REQUESTS where id>? ORDER BY ID DESC limit ? ", start, limit)
	if err != nil {
		return nil, err
	}
	var rows []data.TunnelData
	var temp data.TunnelData
	var tunnelId string
	for results.Next() {
		err := results.Scan(&temp.Id, &tunnelId)
		temp.TunnelId = tunnelId
		if err != nil {
			return nil, err
		}
		rows = append(rows, temp)
	}

	return rows, nil

}

func uncompress(in []byte) []byte {
	buf := bytes.NewBuffer(in)
	reader, _ := gzip.NewReader(buf)
	bytes, _ := io.ReadAll(reader)
	return bytes
}

func compress(in []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(in); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (d *LocalDb) Save(t data.TunnelData) error {
	reqData, err := compress(t.RequestData)
	if err != nil {
		log.Println("Could not compress request Data", err.Error())
	}
	respData, err := compress(t.ResponseData)
	if err != nil {
		log.Println("Could not compress response Data", err.Error())
	}
	_, err = d.db.Exec("INSERT INTO ALL_REQUESTS VALUES (?,?,?,?,?,?,?) ",
		t.Id, t.TunnelId, t.IsReplay, reqData, respData, t.LocalPort, "")
	if err != nil {
		log.Println("Cloud not save tunnelData into database ", err.Error())
		return err
	}
	return nil
}
