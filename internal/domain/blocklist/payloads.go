package blocklist

type PayloadChanged struct {
	Changes []BlocklistChange `json:"changes"`
}
