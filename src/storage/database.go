package storage

import (
	"database/sql"
	"go-gatefuse/src/config"
)

func CreateTables(db *sql.DB) error {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS gate_records (
        src_port INTEGER,
        src_addr TEXT,
        dst_port INTEGER,
        dst_addr TEXT,
		proto TEXT,
        comment TEXT,
        active BOOLEAN,
        uuid TEXT PRIMARY KEY
    );
	CREATE TABLE IF NOT EXISTS app_settings (
		main_domain TEXT,
		nginx_conf_path TEXT,
		username TEXT,
		password TEXT
	);
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
	stmt, err := db.Prepare("UPDATE app_settings SET main_domain = ?, nginx_conf_path = ?, username = ?, password = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(config.Settings.MainDomain, config.Settings.NginxConfPath, config.Settings.Username, config.Settings.Password)
	return err
}

func LoadAppSettings(db *sql.DB) (s config.AppSettings) {
	row := db.QueryRow("SELECT main_domain, nginx_conf_path, username, password FROM app_settings")
	row.Scan(&s.MainDomain, &s.NginxConfPath, &s.Username, &s.Password)
	return s
}
