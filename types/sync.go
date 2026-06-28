package types

const DefaultConfigHash = ""

type SyncState struct {
	ChatsSync    int    `json:"chatsSync"`
	ContactsSync int    `json:"contactsSync"`
	DraftsSync   int    `json:"draftsSync"`
	PresenceSync int    `json:"presenceSync"`
	ConfigHash   string `json:"configHash"`
}

func DefaultSyncState() SyncState {
	return SyncState{
		ChatsSync:    -1,
		ContactsSync: -1,
		DraftsSync:   -1,
		PresenceSync: -1,
		ConfigHash:   DefaultConfigHash,
	}
}

type SyncOverrides struct {
	ChatsSync    *int    `json:"chatsSync"`
	ContactsSync *int    `json:"contactsSync"`
	DraftsSync   *int    `json:"draftsSync"`
	PresenceSync *int    `json:"presenceSync"`
	ConfigHash   *string `json:"configHash"`
}

func (o *SyncOverrides) Resolve(saved SyncState) SyncState {
	result := saved
	if o.ChatsSync != nil {
		result.ChatsSync = *o.ChatsSync
	}
	if o.ContactsSync != nil {
		result.ContactsSync = *o.ContactsSync
	}
	if o.DraftsSync != nil {
		result.DraftsSync = *o.DraftsSync
	}
	if o.PresenceSync != nil {
		result.PresenceSync = *o.PresenceSync
	}
	if o.ConfigHash != nil {
		result.ConfigHash = *o.ConfigHash
	}
	return result
}
