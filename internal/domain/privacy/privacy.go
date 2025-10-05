package privacy

type PrivacySetting string

const (
	PrivacySettingUndefined        PrivacySetting = ""
	PrivacySettingAll              PrivacySetting = "all"
	PrivacySettingContacts         PrivacySetting = "contacts"
	PrivacySettingContactBlacklist PrivacySetting = "contact_blacklist"
	PrivacySettingMatchLastSeen    PrivacySetting = "match_last_seen"
	PrivacySettingKnown            PrivacySetting = "known"
	PrivacySettingNone             PrivacySetting = "none"
)
