package session

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"

	"github.com/akovardin/gomax/types"
)

type Store struct {
	db     *sql.DB
	dbPath string
}

type SessionInfo struct {
	Token        string
	DeviceID     string
	Phone        string
	MtInstanceID string
	Sync         types.SyncState
}

func NewStore(workDir string, dbName string) (*Store, error) {
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return nil, err
	}
	dbPath := filepath.Join(workDir, dbName)
	store := &Store{dbPath: dbPath}
	if err := store.getConnection(); err != nil {
		return nil, err
	}
	if err := store.initializeDB(store.db); err != nil {
		store.db.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) getConnection() error {
	if s.db != nil {
		return nil
	}
	db, err := sql.Open("sqlite", s.dbPath)
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *Store) initializeDB(db *sql.DB) error {
	schema := `CREATE TABLE IF NOT EXISTS sessions (
		token TEXT NOT NULL PRIMARY KEY,
		device_id TEXT NOT NULL,
		phone TEXT NOT NULL,
		mt_instance_id TEXT NOT NULL DEFAULT '',
		chats_sync INTEGER NOT NULL DEFAULT -1,
		contacts_sync INTEGER NOT NULL DEFAULT -1,
		drafts_sync INTEGER NOT NULL DEFAULT -1,
		presence_sync INTEGER NOT NULL DEFAULT -1,
		config_hash TEXT NOT NULL DEFAULT ''
	)`
	_, err := db.Exec(schema)
	return err
}

func (s *Store) SaveSession(info *SessionInfo) error {
	if err := s.getConnection(); err != nil {
		return err
	}
	_, err := s.db.Exec(
		`INSERT OR REPLACE INTO sessions (token, device_id, phone, mt_instance_id, chats_sync, contacts_sync, drafts_sync, presence_sync, config_hash)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		info.Token,
		info.DeviceID,
		info.Phone,
		info.MtInstanceID,
		info.Sync.ChatsSync,
		info.Sync.ContactsSync,
		info.Sync.DraftsSync,
		info.Sync.PresenceSync,
		info.Sync.ConfigHash,
	)
	return err
}

func (s *Store) LoadSession() (*SessionInfo, error) {
	if err := s.getConnection(); err != nil {
		return nil, err
	}
	row := s.db.QueryRow("SELECT token, device_id, phone, mt_instance_id, chats_sync, contacts_sync, drafts_sync, presence_sync, config_hash FROM sessions LIMIT 1")
	return s.rowToSession(row)
}

func (s *Store) LoadSessionByDeviceID(deviceID string) (*SessionInfo, error) {
	if err := s.getConnection(); err != nil {
		return nil, err
	}
	row := s.db.QueryRow("SELECT token, device_id, phone, mt_instance_id, chats_sync, contacts_sync, drafts_sync, presence_sync, config_hash FROM sessions WHERE device_id = ? LIMIT 1", deviceID)
	return s.rowToSession(row)
}

func (s *Store) LoadSessionByPhone(phone string) (*SessionInfo, error) {
	if err := s.getConnection(); err != nil {
		return nil, err
	}
	row := s.db.QueryRow("SELECT token, device_id, phone, mt_instance_id, chats_sync, contacts_sync, drafts_sync, presence_sync, config_hash FROM sessions WHERE phone = ? LIMIT 1", phone)
	return s.rowToSession(row)
}

func (s *Store) DeleteSession(token string) error {
	if err := s.getConnection(); err != nil {
		return err
	}
	_, err := s.db.Exec("DELETE FROM sessions WHERE token = ?", token)
	return err
}

func (s *Store) DeleteAllSessions() error {
	if err := s.getConnection(); err != nil {
		return err
	}
	_, err := s.db.Exec("DELETE FROM sessions")
	return err
}

func (s *Store) UpdateToken(oldToken string, newToken string) error {
	if err := s.getConnection(); err != nil {
		return err
	}
	_, err := s.db.Exec("UPDATE sessions SET token = ? WHERE token = ?", newToken, oldToken)
	return err
}

func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *Store) rowToSession(row *sql.Row) (*SessionInfo, error) {
	var info SessionInfo
	err := row.Scan(
		&info.Token,
		&info.DeviceID,
		&info.Phone,
		&info.MtInstanceID,
		&info.Sync.ChatsSync,
		&info.Sync.ContactsSync,
		&info.Sync.DraftsSync,
		&info.Sync.PresenceSync,
		&info.Sync.ConfigHash,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &info, nil
}
