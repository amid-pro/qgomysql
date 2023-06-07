package qgomysql

type map_si map[string]interface{}

type SSHconfig struct {
	Use_ssh  bool
	User     string
	Password string
	Port     int
	Address  string
}

type MYSQLconfig struct {
	User     string
	Password string
}

type SERVERconfig struct {
	Port      int
	Delimiter string
}

type Config struct {
	SSH    SSHconfig
	MYSQL  MYSQLconfig
	SERVER SERVERconfig
	query  string
	result string
	raw    []map_si
}
