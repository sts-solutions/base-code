package postgres

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/sts-solutions/base-code/ccvalidation"
)

type DBConnection struct {
	Host      string
	DBName    string
	SSLMode   string
	Port      int
	UserName  string
	Password  string
	cnnString string
}

func NewDbConnectionFromCnnString(cnnString string) DBConnection {
	c := DBConnection{
		cnnString: cnnString,
	}
	c.parseConnectionString()

	return c
}

func NewDBConnectionFromCnnStringUrl(cssStr string) (c DBConnection) {

	parseUrl, err := url.Parse(cssStr)
	if err != nil {
		return
	}

	port, err := strconv.Atoi(parseUrl.Port())
	if err != nil {
		return
	}

	password, _ := parseUrl.User.Password()

	c.Password = password
	c.Port = port
	c.Host = parseUrl.Hostname()
	c.UserName = parseUrl.User.Username()
	c.Password, _ = parseUrl.User.Password()
	c.DBName = parseUrl.Path
	c.SSLMode = parseUrl.Query().Get("sslmode")
	c.cnnString = cssStr

	return
}

func NewDBConnectionParams(host string,
	dbName string,
	sslMode string,
	port int,
	userName string,
	password string) *DBConnection {

	c := &DBConnection{
		Host:     host,
		DBName:   dbName,
		SSLMode:  sslMode,
		Port:     port,
		UserName: userName,
		Password: password,
	}
	c.buildConnectionString()
	return c
}

func (c *DBConnection) ConnectionString() string {
	return c.cnnString
}

func (c *DBConnection) Validate() error {
	result := ccvalidation.Result{}

	if c.Host == "" {
		result.AddFailureMessage(fmt.Sprintf("host is not valid: %s", c.Host))
	}

	if c.DBName == "" {
		result.AddFailureMessage(fmt.Sprintf("database name is not valid: %s", c.DBName))
	}

	if c.Port <= 0 {
		result.AddFailureMessage(fmt.Sprintf("port must be greater than zero, got %d", c.Port))
	}

	if c.UserName == "" {
		result.AddFailureMessage(fmt.Sprintf("user name is not valid: %s", c.UserName))
	}

	if c.Password == "" {
		result.AddFailureMessage("password is not valid: <empty>")
	}

	if c.SSLMode == "" {
		result.AddFailureMessage(fmt.Sprintf("ssl mode is not valid: %s", c.SSLMode))
	}

	if result.IsFailure() {
		return result
	}

	return nil
}

func (c *DBConnection) buildConnectionString() {
	c.cnnString = "postgresql://" + c.UserName + ":" + c.Password + "@" + c.Host + ":" + strconv.Itoa(c.Port) + "/" + c.DBName + "?sslmode=" + c.SSLMode
}

func (c *DBConnection) parseConnectionString() {
	parseUrl, err := url.Parse(c.cnnString)
	if err != nil {
		return
	}

	port, err := strconv.Atoi(parseUrl.Port())
	if err != nil {
		return
	}

	password, _ := parseUrl.User.Password()

	c.Password = password
	c.Port = port
	c.Host = parseUrl.Hostname()
	c.UserName = parseUrl.User.Username()
	c.Password, _ = parseUrl.User.Password()
	c.DBName = parseUrl.Path[1:] // remove leading '/'
	c.SSLMode = parseUrl.Query().Get("sslmode")
}
