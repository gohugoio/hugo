package hugosql

// sqlite3 is the default driver and always included.
// Other drivers can be included with their build tag or the build tag alldb
// to include all drivers.

import (
	"bytes"
	"database/sql"
	"errors"
	"strings"
	"time"

	"io/ioutil"

	"strconv"

	"fmt"

	"github.com/spf13/afero"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

var (
	db *sql.DB
	// disable ping in Benchmark with sqlite3 ... with ping the bench fails ...
	pingDb = true
)

// getDriver reads the content of a file or from the ENV for the dataSourceName and
// splits the file name or ENV variable via first underscore to get the driver name.
func getDriver(fs afero.Fs) (dn, dsn string, err error) {

	fileDsn := viper.GetString("SqlSource") // can be a file name or ENV setting

	fOk, err := helpers.Exists(fileDsn, fs)
	if err != nil {
		return dn, dsn, err
	}

	if fOk {
		fp, err := fs.Open(fileDsn)
		if err != nil {
			return dn, dsn, err
		}
		dsnBytes, err := ioutil.ReadAll(fp)
		if err != nil {
			return dn, dsn, err
		}

		if len(dsnBytes) == 0 {
			return dn, dsn, errors.New("Content of file: " + fileDsn + " is empty")
		}
		dsn = string(dsnBytes)
	} else {
		dsn = fileDsn // env var
	}

	if false == strings.ContainsRune(dsn, '_') {
		return "", "", errors.New("Cannot find driver config neither in a file nor in env var HUGO_SQL_SOURCE: " + dsn)
	}

	dnBuf := make([]byte, 0, 100)
	for i := 0; i < len(dsn); i++ {
		if dsn[i] == '_' {
			dn = string(dnBuf)
			break
		}
		dnBuf = append(dnBuf, dsn[i])
	}

	dsn = strings.TrimSpace(dsn[len(dn)+1:])
	dn = strings.TrimSpace(dn)

	found := false
	for _, d := range sql.Drivers() {
		if dn == d {
			found = true
			break
		}
	}

	if false == found {
		return "", "", fmt.Errorf("Your driver name %s cannot be found in the list of available drivers: %s", dn, strings.Join(sql.Drivers(), ", "))
	}
	if dsn == "" {
		return "", "", errors.New("SqlSource is empty. No credentials found")
	}

	return dn, dsn, err
}

// getDb returns the DB instance. Currently DB is unclosed.
func getDb() *sql.DB {

	if db != nil {
		return db
	}

	dn, dsn, err := getDriver(hugofs.SourceFs)
	if err != nil {
		jww.FATAL.Printf("Failed to initialize driver: %s", err)
		return nil
	}

	db, err := sql.Open(dn, dsn)
	if err != nil {
		jww.FATAL.Printf("Cannot open driver: %s. Driver name: %s", err, dn)
		return nil
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	if pingDb {
		if err = db.Ping(); err != nil {
			jww.FATAL.Printf("Database connection failed: %s. Driver name: %s", err, dn)
			return nil
		}
	}
	return db
}

// generateQuery assembles the query from all the query parts. If the query contains
// .sql at the end then the query is considered to be file and the file content gets loaded.
func generateQuery(fs afero.Fs, queryParts ...string) (string, error) {

	qry := strings.TrimSpace(strings.Join(queryParts, " "))
	lQry := strings.ToLower(qry)

	if lQry != "" && strings.Index(lQry, ".sql") == len(qry)-4 {
		ok, err := helpers.Exists(qry, fs)
		if err != nil {
			return "", fmt.Errorf("generateQuery: Error: %s with file: %s", err, qry)
		}
		if !ok {
			return "", fmt.Errorf("generateQuery: File %s not found", qry)
		}

		fp, err := fs.Open(qry)
		if err != nil {
			return "", fmt.Errorf("generateQuery: Cannot open file: %s", qry)
		}
		c, err := ioutil.ReadAll(fp)
		if err != nil {
			return "", fmt.Errorf("generateQuery: Cannot read from file descriptor: %s", err)
		}
		qry = string(c)
		lQry = strings.ToLower(qry)
	}
	if strings.Index(lQry, "select") != 0 {
		return "", fmt.Errorf("SELECT key word not found at beginning of query: %s", qry)
	}
	return qry, nil
}

// GetSql executes a SELECT query and returns a slice containing columns names and its string values
func GetSql(queryParts ...string) []*stringEntities {

	qry, err := generateQuery(hugofs.SourceFs, queryParts...)
	if err != nil {
		jww.ERROR.Print(err)
		return nil
	}

	if getDb() == nil {
		return nil
	}

	rows, err := getDb().Query(qry)
	if err != nil {
		jww.ERROR.Println(err, "Query was:", qry)
		return nil
	}
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		jww.ERROR.Println("GetSql: Cannot retrieve column names", err)
		return nil
	}

	ret := make([]*stringEntities, 0, 2000)
	rss := newRowTransformer(columnNames)
	for rows.Next() {

		if err := rows.Scan(rss.cp...); err != nil {
			jww.ERROR.Println("GetSql: Cannot scan a row", err)
			return nil
		}
		err := rss.toString()
		if err != nil {
			jww.ERROR.Println("GetSql: ", err)
			return nil
		}
		rss.append(&ret)
	}
	return ret
}

type rowTransformer struct {
	// cp are the column pointers
	cp []interface{}
	// row contains the final row result
	se       *stringEntities
	colCount int
	colNames []string
}

func newRowTransformer(columnNames []string) *rowTransformer {
	lenCN := len(columnNames)
	s := &rowTransformer{
		cp:       make([]interface{}, lenCN),
		se:       newStringEntities(columnNames),
		colCount: lenCN,
		colNames: columnNames,
	}
	for i := 0; i < lenCN; i++ {
		s.cp[i] = new(sql.RawBytes)
	}
	return s
}

func (s *rowTransformer) toString() error {
	for i := 0; i < s.colCount; i++ {
		if rb, ok := s.cp[i].(*sql.RawBytes); ok {
			s.se.row[s.colNames[i]] = string(*rb)
			*rb = nil // reset pointer to discard current value to avoid a bug
		} else {
			return errors.New("Cannot convert index " + strconv.Itoa(i) + " column " + s.colNames[i] + " to type *sql.RawBytes")
		}
	}
	return nil
}

// append appends the current row to the ret return value and clears the row result
func (s *rowTransformer) append(ret *[]*stringEntities) {
	*ret = append(*ret, s.se)
	s.se = newStringEntities(s.colNames)
}

func newStringEntities(col []string) *stringEntities {
	return &stringEntities{
		col: col,
		row: make(rowStrings, len(col)),
	}
}

type (
	rowStrings map[string]string

	stringEntities struct {
		// col columns to preserve the order of columns when ranging over the map.
		col []string
		row rowStrings
	}
)

func (s *stringEntities) Column(c string) string {
	if v, ok := s.row[c]; ok {
		return v
	}
	return ""
}

func (s *stringEntities) Columns() []string {
	return s.col
}

func (s *stringEntities) JoinColumns(sep string) string {
	return strings.Join(s.col, sep)
}

// Join joins uses separator sep to join n columns or all columns when using * as 2nd argument.
func (s *stringEntities) JoinValues(sep string, columns ...string) string {

	if len(columns) == 0 {
		return ""
	}
	if len(columns) == 1 {
		if v, ok := s.row[columns[0]]; ok {
			return v
		}
		if columns[0] == "*" {
			columns = s.col
		} else {
			return ""
		}
	}

	var buf bytes.Buffer
	lc1 := len(columns) - 1
	for i, c := range columns {
		if v, ok := s.row[c]; ok {
			buf.WriteString(v)
			if i < lc1 {
				buf.WriteString(sep)
			}
		}
	}
	return buf.String()
}

// Int parses the value from column c to an int value
func (s *stringEntities) Int(c string) int64 {
	if v, ok := s.row[c]; ok {
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			jww.WARN.Printf("Failed to convert column %s with value %s to int", c, v)
			return 0
		}
		return i
	}
	return 0
}

// Float parses the value from column c to a float64 value
func (s *stringEntities) Float(c string) float64 {
	if v, ok := s.row[c]; ok {
		i, err := strconv.ParseFloat(v, 64)
		if err != nil {
			jww.WARN.Printf("Failed to convert column %s with value %s to float64", c, v)
			return 0
		}
		return i
	}
	return 0
}

// DateTime takes c as the column name and parses its layout according to parameter l.
// If you provide for layout l the term unix then the value is to be considered a unix timestamp.
func (s *stringEntities) DateTime(c, l string) time.Time {
	if l == "" {
		l = "2006-01-02 15:04:05.999999"
	}
	var err error
	var t time.Time
	if v, ok := s.row[c]; ok {
		if l == "unix" {
			i, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				jww.WARN.Printf("Failed to parse column %s with value %s to int for unixtime", c, v)
				return t
			}
			t = time.Unix(i, 0)
		} else {
			t, err = time.Parse(l, v)
			if err != nil {
				jww.WARN.Printf("Failed to parse column %s with value %s to time", c, v)
				return t
			}
		}
		return t
	}
	return t
}
