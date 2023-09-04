package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"go-gatefuse/src/config"
	"reflect"
	"strings"
)

func InitializeDatabaseTables(db *sql.DB) error {
	createTableSQL := `
	DROP TABLE IF EXISTS gate_records;
	DROP TABLE IF EXISTS app_settings;
    CREATE TABLE gate_records (
        src_port INTEGER,
        src_addr TEXT,
        dst_port INTEGER,
        dst_addr TEXT,
		proto TEXT,
        comment TEXT,
        active BOOLEAN,
        uuid TEXT PRIMARY KEY
    );
	CREATE TABLE app_settings (
		main_domain TEXT NOT NULL,
		nginx_conf_path TEXT NOT NULL,
		unbound_conf_path TEXT NOT NULL,
  		unbound_remote BOOLEAN NOT NULL,
  		unbound_remote_host TEXT NOT NULL,
		username TEXT NOT NULL,
		password TEXT NOT NULL
	);
	INSERT INTO app_settings(main_domain, nginx_conf_path, unbound_conf_path, unbound_remote, unbound_remote_host, username, password) VALUES ("", "", "", false, "", "admin", "admin");
    `
	_, err := db.Exec(createTableSQL)
	return err
}

func AddNewRecord(db *sql.DB, s config.GateRecord) error {
	stmt, err := db.Prepare("INSERT INTO gate_records(src_port, src_addr, dst_port, dst_addr, proto, comment, active, uuid) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(s.SrcPort, s.SrcAddr, s.DstPort, s.DstAddr, s.Protocol, s.Comment, s.Active, s.UUID)
	return err
}

func RetrieveAllGateRecords(db *sql.DB) ([]config.GateRecord, error) {
	rows, err := db.Query("SELECT src_port, src_addr, dst_port, dst_addr, proto, comment, active, uuid FROM gate_records")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []config.GateRecord
	for rows.Next() {
		var record config.GateRecord
		err := rows.Scan(&record.SrcPort, &record.SrcAddr, &record.DstPort, &record.DstAddr, &record.Protocol, &record.Comment, &record.Active, &record.UUID)
		if err != nil {
			return nil, err
		}
		items = append(items, record)
	}

	// Check for errors from iterating over rows.
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func RetrieveOneGateRecord(db *sql.DB, uuid string) (record config.GateRecord, err error) {
	row := db.QueryRow("SELECT src_port, src_addr, dst_port, dst_addr, proto, comment, active, uuid FROM gate_records WHERE uuid=?", uuid)
	err = row.Scan(&record.SrcPort, &record.SrcAddr, &record.DstPort, &record.DstAddr, &record.Protocol, &record.Comment, &record.Active, &record.UUID)
	return record, err
}

func UpdateGateRecord(db *sql.DB, record config.GateRecord) error {
	stmt, err := db.Prepare("UPDATE gate_records SET src_port = ?, src_addr = ?, dst_port = ?, dst_addr = ?, comment = ?, active = ? WHERE uuid = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(record.SrcPort, record.SrcAddr, record.DstPort, record.DstAddr, record.Comment, record.Active, record.UUID)
	return err
}

func DeleteGateRecord(db *sql.DB, uuid string) error {
	query := "DELETE FROM gate_records WHERE uuid = ?"
	_, err := db.Exec(query, uuid)
	return err
}

func SaveAppSettings(db *sql.DB) error {
	t := reflect.TypeOf(config.Settings)
	v := reflect.ValueOf(config.Settings)

	if t.Kind() != reflect.Struct {
		return errors.New("src should be of struct type")
	}

	columns := make([]string, t.NumField())
	values := make([]interface{}, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		jsonTag := t.Field(i).Tag.Get("json")
		if jsonTag == "" {
			return fmt.Errorf("missing json tag for field %s", t.Field(i).Name)
		}

		columns[i] = fmt.Sprintf("%s = ?", jsonTag)
		values[i] = v.Field(i).Interface()
	}
	query := fmt.Sprintf("UPDATE app_settings SET %s;", strings.Join(columns, ","))
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(values...)
	return err

}

// LoadAppSettings Load settings from SQL to the structure described as `dst`
func LoadAppSettings(db *sql.DB, dst interface{}) error {
	t := reflect.TypeOf(dst).Elem()
	v := reflect.ValueOf(dst).Elem()

	if t.Kind() != reflect.Struct {
		return errors.New("dst should be a pointer to struct type")
	}

	columns := make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		columns[i] = t.Field(i).Tag.Get("json")
	}

	row := db.QueryRow(fmt.Sprintf("SELECT %s FROM app_settings LIMIT 1", strings.Join(columns, ",")))

	values := make([]interface{}, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		values[i] = v.Field(i).Addr().Interface()
	}

	return row.Scan(values...)
}
