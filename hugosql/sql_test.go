package hugosql

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"time"

	"github.com/spf13/afero"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// go test -coverprofile=xcover.out .
// go tool cover -html=xcover.out
// coverage: 87.7% of statements

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	// http://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go A note on compiler optimisations
	benchSqlResult               []*stringEntities
	benchErr                     error
	benchDn, benchDsn, benchJoin string
)

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func validTestFile() afero.Fs {
	fs := new(afero.MemMapFs)
	fn := "path/to/mysql_accessdata.txt"
	viper.Set("SqlSource", fn)
	inMemFile, _ := fs.Create(fn)
	inMemFile.Write([]byte(`userName:passw0rd@tcp(:3306)/databaseName`))
	return fs
}

func inValidTestFile() afero.Fs {
	fs := new(afero.MemMapFs)
	fn := "path/to/wrong/mysql-accessdata.txt"
	viper.Set("SqlSource", fn)
	inMemFile, _ := fs.Create(fn)
	inMemFile.Write([]byte("Hello World"))
	return fs
}

func TestColumn(t *testing.T) {
	test := newStringEntities([]string{"k1", "k2"})
	test.row = rowStrings{
		"k2": "",
		"k1": "value1",
	}
	assert.Equal(t, "value1", test.Column("k1"))
	assert.Equal(t, "", test.Column("k2"))
	assert.Equal(t, "", test.Column("k3"))
}

func TestJoin(t *testing.T) {
	test := newStringEntities([]string{"k1", "k2", "k3", "k4", "k5", "k6"})
	test.row = rowStrings{
		"k1": "value1",
		"k2": "value2",
		"k3": "value3",
		"k4": "value4",
		"k5": "value5",
		"k6": "value6",
	}
	assert.Equal(t, "value1", test.JoinValues("", "k1"))
	assert.Equal(t, "", test.JoinValues("", "k0"))
	assert.Equal(t, "", test.JoinValues(""))
	assert.Equal(t, "value1,value2", test.JoinValues(",", "k1", "k2"))
	assert.Equal(t, "value1,value2,value3", test.JoinValues(",", "k1", "k2", "k3"))
	assert.Equal(t, "value1,value2,value3,value4,value5,value6", test.JoinValues(",", "*"))
	// do it twice to check if we correctly range over the map to preserve the order
	assert.Equal(t, "value1,value2,value3,value4,value5,value6", test.JoinValues(",", "*"))
	assert.Equal(t, "k1,k2,k3,k4,k5,k6", test.JoinColumns(","))
}

// $ go test -run=NONE -bench=BenchmarkJoin > BenchmarkJoin_old|new.txt
// $ benchcmp BenchmarkJoin_old.txt BenchmarkJoin_new.txt
// BenchmarkJoinValues	  500000	      2474 ns/op	     800 B/op	       4 allocs/op
func BenchmarkJoinValues(b *testing.B) {
	exp := `PHPPHPPHPPHPPHP;GolangGolangGolangGolangGolang;JavaJavaJavaJavaJava;JavaScriptJavaScriptJavaScriptJavaScriptJavaScript;PerlPerlPerlPerlPerl;RubyRubyRubyRubyRuby;ASPASPASPASPASP;C++C++C++C++C++;DDDDD;CCCCC`
	b.ReportAllocs()
	cols := []string{"k0", "k1", "k2", "k3", "k4", "k5", "k6", "k7", "k8", "k9"}
	test := newStringEntities(cols)
	test.row = rowStrings{
		"k0": strings.Repeat("PHP", 5),
		"k1": strings.Repeat("Golang", 5),
		"k2": strings.Repeat("Java", 5),
		"k3": strings.Repeat("JavaScript", 5),
		"k4": strings.Repeat("Perl", 5),
		"k5": strings.Repeat("Ruby", 5),
		"k6": strings.Repeat("ASP", 5),
		"k7": strings.Repeat("C++", 5),
		"k8": strings.Repeat("D", 5),
		"k9": strings.Repeat("C", 5),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchJoin = test.JoinValues(";", cols...)
		if benchJoin != exp {
			b.Fatalf("\nExpected %s\nGot %s\n", exp, benchJoin)
		}
	}
}

func TestInt(t *testing.T) {
	test := newStringEntities([]string{"k1", "k2", "k3", "k4"})
	test.row = rowStrings{
		"k1": "value1",
		"k2": "0815",
		"k3": "10000000000000",
		"k4": "3.141",
	}
	assert.Equal(t, 0, test.Int("k1"))
	assert.Equal(t, 815, test.Int("k2"))
	assert.Equal(t, 10000000000000, test.Int("k3"))
	assert.Equal(t, 0, test.Int("k4"))
	assert.Equal(t, 0, test.Int("k5"))
}

func TestFloat(t *testing.T) {
	test := newStringEntities([]string{"k1", "k2", "k3", "k4"})
	test.row = rowStrings{
		"k1": "value1",
		"k2": "0815",
		"k3": "10000000000000",
		"k4": "3.141",
	}
	assert.Equal(t, 0, test.Float("k1"))
	assert.Equal(t, 815, test.Float("k2"))
	assert.Equal(t, 1e+13, test.Float("k3"))
	assert.Equal(t, 3.141, test.Float("k4"))
	assert.Equal(t, 0, test.Float("k5"))
}

func TestDateTime(t *testing.T) {
	testSSR := newStringEntities([]string{"k1", "k2", "k3", "k4", "k5"})
	testSSR.row = rowStrings{
		// @todo add more database time formats
		"k1": "2012-07-09 11:16:13",
		"k2": "2014-08-26 02:15:47",
		"k3": "10000000000000",
		"k4": "2015-02-27 09:58:59",
		"k5": "1424991728",
	}

	tests := []struct {
		expected, column, layout, format string
	}{
		{
			"Mon Jul 9 11:16:13 +0000 UTC 2012",
			"k1",
			"2006-01-02 15:04:05.999999",
			"Mon Jan 2 15:04:05 -0700 MST 2006",
		},
		{
			"2014-08-26",
			"k2",
			"2006-01-02 15:04:05.999999",
			"2006-01-02",
		},
		{
			"0001-01-01",
			"k3",
			"2006-01-02 15:04:05.999999",
			"2006-01-02",
		},
		{
			"2015/02/27 09:58",
			"k4",
			"",
			"2006/01/02 15:04",
		},
		{
			time.Unix(1424991728, 0).Format("2006/01/02 15:04"),
			"k5",
			"unix",
			"2006/01/02 15:04",
		},
		{
			"0001/01/01 00:00",
			"kNOTFOUND",
			"",
			"2006/01/02 15:04",
		},
	}

	for _, test := range tests {
		assert.Equal(
			t,
			test.expected,
			testSSR.DateTime(test.column, test.layout).Format(test.format),
		)
	}
}

func TestGenerateQuery(t *testing.T) {

	tests := []struct {
		query, eQuery string
		isFile, eErr  bool
	}{
		{
			"SELECT * FROM gopher1", "SELECT * FROM gopher1", false, false,
		},
		{
			"\t SELECT * FROM gopher2\n", "SELECT * FROM gopher2", false, false,
		},
		{
			" UPDATE FROM gopher\n", "", false, true,
		},
		{
			"path/to/benchmark.sql", "SELECT * FROM gopher3", true, false,
		},
		{
			"path/to/benchmark....sql", "", true, true,
		},
	}

	for _, test := range tests {
		var fs afero.Fs
		if test.isFile {
			fs = new(afero.MemMapFs)
			fn := test.query
			if test.eErr {
				fn = "non-existent" + test.query
			}
			inMemFile, err := fs.Create(fn)
			if err != nil {
				t.Error(err, test.query)
			}
			inMemFile.Write([]byte(test.eQuery))
		}
		qry, err := generateQuery(fs, test.query)
		if test.eErr {
			assert.Error(t, err, test.query)
		} else {
			assert.NoError(t, err, test.query)
		}
		assert.Equal(t, test.eQuery, qry, test.query)
	}

}

// $ go test -run=NONE -bench=BenchmarkGetDriver > BenchmarkGetDriver_old|new.txt
// $ benchcmp BenchmarkGetDriver_old.txt BenchmarkGetDriver_new.txt
// BenchmarkGetDriver	  200000	      7976 ns/op	    1098 B/op	      18 allocs/op
func BenchmarkGetDriver(b *testing.B) {
	b.ReportAllocs()
	fs := validTestFile()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchDn, benchDsn, benchErr = getDriver(fs)
	}
}

func TestGetDriver(t *testing.T) {
	defer viper.Set("SqlSource", "")
	tests := []struct {
		fileName, fileContent, eFileContent, eDriverName string
		eErr                                             bool
	}{
		{
			"thisIsRandom.txt",
			"mysql_userName:passw0rd@tcp(:3306)/databaseName\n",
			"userName:passw0rd@tcp(:3306)/databaseName",
			"mysql",
			false,
		},
		{
			"this_IsRandom.txt",
			"\tmysql_userName:passw0rd@tcp(:3306)/databaseName\n",
			"userName:passw0rd@tcp(:3306)/databaseName",
			"mysql",
			false,
		},
		{
			"thisIsAlsoRandom.txt",
			"\n\n\n",
			"",
			"",
			true,
		},
		{
			"thisIsAlsoRandom2.txt",
			"",
			"",
			"",
			true,
		},
		{
			"path" + helpers.FilePathSeparator + "to" + helpers.FilePathSeparator + "mysql-thisIsRandom2.txt",
			"mysql-userName:passw0rd@tcp(:3306)/database_Name",
			"userName:passw0rd@tcp(:3306)/databaseName",
			"",
			true,
		},
		{
			"C:" + helpers.FilePathSeparator + "windows" + helpers.FilePathSeparator + "thisIsRandom3.txt",
			"sqlite3_userKJHDJKHFSDFbaseName",
			"userKJHDJKHFSDFbaseName",
			"sqlite3",
			false,
		},
		{
			"C:" + helpers.FilePathSeparator + "windows" + helpers.FilePathSeparator + "thisIsRandom4.txt",
			"sqlite3_",
			"",
			"sqlite3",
			true,
		},
	}

	fs := new(afero.MemMapFs)
	viper.Set("SqlSource", "non-existent.txt")

	aDn, aDsn, aErr := getDriver(fs)
	assert.Error(t, aErr)
	assert.Equal(t, "", aDn)
	assert.Equal(t, "", aDsn)

	for _, test := range tests {
		viper.Set("SqlSource", test.fileName)
		fs := new(afero.MemMapFs)
		inMemFile, err := fs.Create(test.fileName)
		if err != nil {
			t.Fatal(err)
		}
		inMemFile.Write([]byte(test.fileContent))

		aDn, aDsn, aErr := getDriver(fs)
		if test.eErr {
			assert.Error(t, aErr, test.fileName)
		} else {
			assert.NoError(t, aErr, test.fileName)

			assert.Equal(t, test.eDriverName, aDn, test.fileName, "Driver Name")
			assert.Equal(t, test.eFileContent, aDsn, test.fileName, "Data Source Name")
		}
	}
}

func TestGetDriverEnv(t *testing.T) {
	viper.Set("SqlSource", "mysql_userName:passw0rd@tcp(127.0.0.1:5432)/databaseName")
	fs := new(afero.MemMapFs)
	aDn, aDsn, aErr := getDriver(fs)
	assert.NoError(t, aErr)
	assert.Equal(t, "mysql", aDn, "Driver Name via env")
	assert.Equal(t, "userName:passw0rd@tcp(127.0.0.1:5432)/databaseName", aDsn, "Driver Source Name via env")

	viper.Set("SqlSource", "")
	aDn, aDsn, aErr = getDriver(fs)
	assert.Equal(t, "", aDn)
	assert.Equal(t, "", aDsn)
	assert.Error(t, aErr)
}

func TestGetDb(t *testing.T) {
	tmpDir := helpers.GetTempDir("", hugofs.SourceFs)
	tests := []struct {
		sqlSource string
		dbFile    string
		isNil     bool
	}{
		{
			tmpDir + "tmp_" + randString(10) + ".txt",
			"sqlite_" + tmpDir + "hugo_testing-getSql_sqlite_" + randString(10) + ".db",
			true,
		},
		{
			tmpDir + "tmp_" + randString(10) + ".txt",
			"sqlite3_" + tmpDir + "hugo-testing_getSql_sqlite_" + randString(10) + ".db",
			false,
		},
	}

	for _, test := range tests {
		defer os.Remove(test.sqlSource)
		defer os.Remove(test.dbFile)
		if err := ioutil.WriteFile(test.sqlSource, []byte(test.dbFile), 0600); err != nil {
			t.Fatal(err)
		}
		viper.Set("SqlSource", test.sqlSource)
		if test.isNil {
			assert.Nil(t, getDb())
		} else {
			assert.NotNil(t, getDb())
		}

	}
	viper.Set("SqlSource", "")
}

func TestGetSqlFail(t *testing.T) {

	tdb := GetSql("INSERT INTO", "sales")
	assert.Nil(t, tdb)

	var tmpFs afero.Fs
	tmpFs, hugofs.SourceFs = hugofs.SourceFs, inValidTestFile()
	defer func() { hugofs.SourceFs = tmpFs }()
	tdb = GetSql("SELECT * from ", "sales")
	assert.Nil(t, tdb)
}

func initTestDb(t testing.TB, randFileName bool) (sqlSource, dbName string) {

	checkFail := func(err error) {
		if err != nil {
			t.Fatal("Maybe remove test files in $TMPDIR", err)
		}
	}
	fileNamePart := ""
	if randFileName {
		fileNamePart = "_" + randString(8)
	}

	dbName = helpers.GetTempDir("", hugofs.SourceFs) + "hugo_testing_getSql_sqlite" + fileNamePart + ".db"
	sqlSource = helpers.GetTempDir("", hugofs.SourceFs) + "testing_getSql" + fileNamePart + ".txt"
	isSource, err := helpers.Exists(sqlSource, hugofs.SourceFs)
	checkFail(err)
	if isSource {
		return sqlSource, dbName
	}

	db, err := sql.Open("sqlite3", dbName)
	checkFail(err)
	defer db.Close()

	sqlStmt := `
	create table foo (id integer not null primary key, name text,city varchar(20) not null,street varchar(10) null);
	delete from foo;
	`
	_, err = db.Exec(sqlStmt)
	checkFail(err)

	tx, err := db.Begin()
	checkFail(err)

	stmt, err := tx.Prepare("insert into foo(id, name, city) values(?, ?, ?)")
	defer stmt.Close()
	checkFail(err)

	for i := 0; i < 50; i++ {
		_, err = stmt.Exec(
			i,
			fmt.Sprintf("こんにちわ世界 %03d %s", i, randString(i)),
			fmt.Sprintf("ゴーファー都市 %03d", i),
		)
		checkFail(err)
	}
	tx.Commit()

	// un-comment this to check if sqlite3 has been successfully written to
	//	rows, err := db.Query("select id, city from foo")
	//	checkFail(err)
	//	defer rows.Close()
	//	for rows.Next() {
	//		var id int
	//		var city string
	//		rows.Scan(&id, &city)
	//		fmt.Println(id, city)
	//	}

	err = ioutil.WriteFile(sqlSource, []byte("sqlite3_"+dbName+"\n"), 0600)
	checkFail(err)

	return sqlSource, dbName
}

func TestGetSql(t *testing.T) {
	sqlSource, dbName := initTestDb(t, true)
	viper.Set("SqlSource", sqlSource)
	defer os.Remove(sqlSource)
	defer os.Remove(dbName)

	fooResult := GetSql("SELECT * from foo")
	assert.NotNil(t, fooResult)
	assert.Equal(t, 50, len(fooResult))

	for i, rows := range fooResult {
		assert.Equal(t, fmt.Sprintf("ゴーファー都市 %03d", i), rows.Column("city"))
		assert.Equal(t, "", rows.Column("typo_in_column_name"))
	}

	fooResult2 := GetSql("SELECT id from ", "foo where id <= 10")
	assert.NotNil(t, fooResult2)
	assert.Equal(t, 11, len(fooResult2))

	for i, rows := range fooResult2 {
		assert.Equal(t, strconv.Itoa(i), rows.Column("id"))
		assert.Equal(t, "", rows.Column("non-existent_column"))
	}

	fooResult3 := GetSql("SELECT id froooom", " foo where id <= 10")
	assert.Nil(t, fooResult3)

	qryFile := helpers.GetTempDir("", hugofs.SourceFs) + "hugo_testing_query_" + randString(8) + ".sql"
	defer os.Remove(qryFile)
	err := ioutil.WriteFile(qryFile, []byte(`SELECT * FROM foo where id <= 12`), 0600)
	if err != nil {
		t.Error(err)
	}

	fooResult4 := GetSql(qryFile)
	assert.NotNil(t, fooResult4)
	assert.Equal(t, 13, len(fooResult4))
}

// $ go test -run=NONE -bench=BenchmarkGetSql > BenchmarkGetSql_old|new.txt
// $ benchcmp BenchmarkGetSql_old.txt BenchmarkGetSql_new.txt
// BenchmarkGetSql	    2000	    921062 ns/op	   61874 B/op	     797 allocs/op
func BenchmarkGetSql(b *testing.B) {
	pingDb = false
	reset := func() { pingDb = true }
	defer reset()
	b.ReportAllocs()
	sqlSource, dbName := initTestDb(b, false)
	viper.Set("SqlSource", sqlSource)
	// cannot remove files otherwise benchmark fails ... and no chance to run an afterBench function :-(
	//	defer os.Remove(sqlSource)
	//	defer os.Remove(dbName)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchSqlResult = GetSql("SELECT * from foo")
		if benchSqlResult == nil {
			b.Log("\nsqlSource File:", sqlSource, "\nsqlite3 file:", dbName, "\n")
			b.FailNow()
		} else {
			if 50 != len(benchSqlResult) {
				b.Fatalf("BenchmarkGetSql: Expected 50 rows but got %d", len(benchSqlResult))
			}
		}
	}
}
