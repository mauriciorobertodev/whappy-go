package privacy

// To listen all privacy events use "privacy:*"
const (
	// To listen all change events, use the prefix "privacy:changed/*"
	EventChangedLastSeen     = "privacy:changed/last_seen"
	EventChangedStatus       = "privacy:changed/status"
	EventChangedProfile      = "privacy:changed/profile"
	EventChangedGroupAdd     = "privacy:changed/group_add"
	EventChangedReadReceipts = "privacy:changed/read_receipts"
	EventChangedCallAdd      = "privacy:changed/call_add"
	EventChangedOnline       = "privacy:changed/online"
)
