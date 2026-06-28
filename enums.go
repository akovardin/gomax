package gomax

type EventType string

const (
	EventTypeMessageNew     EventType = "MESSAGE_NEW"
	EventTypeMessageEdit    EventType = "MESSAGE_EDIT"
	EventTypeMessageDelete  EventType = "MESSAGE_DELETE"
	EventTypeMessageRead    EventType = "MESSAGE_READ"
	EventTypeTyping         EventType = "TYPING"
	EventTypePresence       EventType = "PRESENCE"
	EventTypeReactionUpdate EventType = "REACTION_UPDATE"
	EventTypeChatUpdate     EventType = "CHAT_UPDATE"
	EventTypeUserUpdate     EventType = "USER_UPDATE"
	EventTypeVideoReady     EventType = "VIDEO_READY"
	EventTypeFileReady      EventType = "FILE_READY"
	EventTypeRaw            EventType = "RAW"
	EventTypeOnStart        EventType = "ON_START"
)

type ChatType string

const (
	ChatTypeDialog  ChatType = "DIALOG"
	ChatTypeGroup   ChatType = "GROUP"
	ChatTypeChannel ChatType = "CHANNEL"
)

type AccessType string

const (
	AccessTypePublic  AccessType = "PUBLIC"
	AccessTypePrivate AccessType = "PRIVATE"
)

type MessageStatus string

const (
	MessageStatusSent     MessageStatus = "SENT"
	MessageStatusDelivered MessageStatus = "DELIVERED"
	MessageStatusRead     MessageStatus = "READ"
	MessageStatusEdited   MessageStatus = "EDITED"
	MessageStatusRemoved  MessageStatus = "REMOVED"
)

type AttachmentType string

const (
	AttachmentTypePhoto         AttachmentType = "PHOTO"
	AttachmentTypeVideo         AttachmentType = "VIDEO"
	AttachmentTypeFile          AttachmentType = "FILE"
	AttachmentTypeSticker       AttachmentType = "STICKER"
	AttachmentTypeAudio         AttachmentType = "AUDIO"
	AttachmentTypeControl       AttachmentType = "CONTROL"
	AttachmentTypeContact       AttachmentType = "CONTACT"
	AttachmentTypeCall          AttachmentType = "CALL"
	AttachmentTypeShare         AttachmentType = "SHARE"
	AttachmentTypeInlineKeyboard AttachmentType = "INLINE_KEYBOARD"
	AttachmentTypeUnknown       AttachmentType = "UNKNOWN"
)

type AuthType int

const (
	AuthTypeStartAuth           AuthType = 0
	AuthTypeResendCode          AuthType = 1
	AuthTypeCheckCode           AuthType = 2
	AuthTypeRegistration        AuthType = 3
	AuthTypePassword            AuthType = 4
	AuthTypeLogin               AuthType = 5
	AuthTypeRestorePassword     AuthType = 6
	AuthTypeQRLogin             AuthType = 7
	AuthTypeRestoreByEmail      AuthType = 8
	AuthTypeSendCodeCall        AuthType = 9
	AuthTypeVerifyEmail         AuthType = 10
	AuthTypeChangePhone         AuthType = 11
)

type ErrorScope string

const (
	ErrorScopeGlobal ErrorScope = "GLOBAL"
	ErrorScopeLocal  ErrorScope = "LOCAL"
)

type ReadAction int

const (
	ReadActionRead   ReadAction = 0
	ReadActionUnread ReadAction = 1
)

type ItemType string

const (
	ItemTypeMessage ItemType = "message"
	ItemTypeChat    ItemType = "chat"
	ItemTypeContact ItemType = "contact"
	ItemTypeFile    ItemType = "file"
	ItemTypePhoto   ItemType = "photo"
	ItemTypeVideo   ItemType = "video"
	ItemTypeAudio   ItemType = "audio"
	ItemTypeSticker ItemType = "sticker"
)

type ControlEvent string

const (
	ControlEventTypingStart  ControlEvent = "typing_start"
	ControlEventTypingStop   ControlEvent = "typing_stop"
	ControlEventRecordingStart ControlEvent = "recording_start"
	ControlEventRecordingStop  ControlEvent = "recording_stop"
	ControlEventSeen         ControlEvent = "seen"
	ControlEventDelivered    ControlEvent = "delivered"
)

type ChatMemberOperation string

const (
	ChatMemberOperationAdd         ChatMemberOperation = "add"
	ChatMemberOperationRemove      ChatMemberOperation = "remove"
	ChatMemberOperationSetAdmin    ChatMemberOperation = "set_admin"
	ChatMemberOperationRemoveAdmin ChatMemberOperation = "remove_admin"
	ChatMemberOperationLeave       ChatMemberOperation = "leave"
	ChatMemberOperationJoin        ChatMemberOperation = "join"
)

type ChatOption string

const (
	ChatOptionShowPreview    ChatOption = "show_preview"
	ChatOptionSilent         ChatOption = "silent"
	ChatOptionMute           ChatOption = "mute"
	ChatOptionPin            ChatOption = "pin"
	ChatOptionUnpin          ChatOption = "unpin"
	ChatOptionBlock          ChatOption = "block"
	ChatOptionUnblock        ChatOption = "unblock"
	ChatOptionReport         ChatOption = "report"
	ChatOptionClearHistory   ChatOption = "clear_history"
	ChatOptionDelete         ChatOption = "delete"
	ChatOptionArchive        ChatOption = "archive"
	ChatOptionUnarchive      ChatOption = "unarchive"
)

type ContactAction string

const (
	ContactActionAccept  ContactAction = "accept"
	ContactActionReject  ContactAction = "reject"
	ContactActionBlock   ContactAction = "block"
	ContactActionUnblock ContactAction = "unblock"
	ContactActionDelete  ContactAction = "delete"
	ContactActionHide    ContactAction = "hide"
)

type AvatarType string

const (
	AvatarTypePhoto AvatarType = "photo"
	AvatarTypeVideo AvatarType = "video"
)
