package connection_strings

var config = map[string]string{
	"DB_CONNECTION": "sqlserver",
	"DB_HOST":       "10.64.5.151",
	"DB_INSTANCE":   "sqlexpress",
	"DB_PORT":       "1433",
	"DB_DATABASE":   "CisSysDB",
	"DB_USERNAME":   "dev",
	"DB_PASSWORD":   "devpwd",
	"CRYPTO_KEY":    "kVEkdsor4Ms=", // base 64 of double encrypted key
	"PRIV_KEY":      "TestPriv",
}
