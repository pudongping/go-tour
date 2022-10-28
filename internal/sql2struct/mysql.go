package sql2struct

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type DBModel struct {
	DBEngine *sql.DB
	DBInfo   *DBInfo
}

type DBInfo struct {
	DBType   string
	Host     string
	UserName string
	Password string
	Charset  string
}

type TableColumn struct {
	ColumnName    string // 类的名称
	DataType      string // 列的数据类型，仅包含类型信息
	IsNullable    string // 列是否允许为 null
	ColumnKey     string // 列是否被索引
	ColumnType    string // 列的数据类型，包含类型名称和可能的其他信息。eg：精度，长度，是否无符号等
	ColumnComment string // 列的注释信息
	ColumnDefault string // 列的默认值
	Extra         string // 额外信息
}

// DBTypeToStructType 表字段类型映射
var DBTypeToStructType = map[string]string{
	"integer":            "int64",
	"int":                "int",
	"int unsigned":       "uint",
	"tinyint":            "int8",
	"tinyint unsigned":   "uint8",
	"smallint":           "int16",
	"smallint unsigned":  "uint16",
	"mediumint":          "int64",
	"mediumint unsigned": "uint64",
	"bigint":             "int64",
	"bigint unsigned":    "uint64",
	"bit":                "int",
	"bool":               "bool",
	"enum":               "string",
	"set":                "string",
	"varchar":            "string",
	"char":               "string",
	"tinytext":           "string",
	"mediumtext":         "string",
	"text":               "string",
	"longtext":           "string",
	"blob":               "string",
	"tinyblob":           "string",
	"mediumblob":         "string",
	"longblob":           "string",
	"json":               "string",
	"date":               "time.Time",
	"datetime":           "time.Time",
	"timestamp":          "time.Time",
	"time":               "time.Time",
	"float":              "float64",
	"float unsigned":     "float64",
	"double":             "float64",
	"double unsigned":    "float64",
	// 考虑到整数溢出以及浮点数精度丢失风险，建议将 `decimal` 类型以 `string` 类型接收。（当然也可以使用 float64 接收，操作比较方便）
	"decimal": "float64",
}

var TypeMysqlMatchList = []struct {
	Key   string
	Value string
}{
	{`^(tinyint)[(]\d+[)] unsigned`, "uint8"},
	{`^(smallint)[(]\d+[)] unsigned`, "uint16"},
	{`^(int)[(]\d+[)] unsigned`, "uint32"},
	{`^(bigint)[(]\d+[)] unsigned`, "uint64"},
	{`^(float)[(]\d+,\d+[)] unsigned`, "float64"},
	{`^(double)[(]\d+,\d+[)] unsigned`, "float64"},
	{`^(tinyint)[(]\d+[)]`, "int8"},
	{`^(smallint)[(]\d+[)]`, "int16"},
	{`^(int)[(]\d+[)]`, "int"},
	{`^(bigint)[(]\d+[)]`, "int64"},
	{`^(char)[(]\d+[)]`, "string"},
	{`^(enum)[(](.)+[)]`, "string"},
	{`^(varchar)[(]\d+[)]`, "string"},
	{`^(varbinary)[(]\d+[)]`, "[]byte"},
	{`^(blob)[(]\d+[)]`, "[]byte"},
	{`^(binary)[(]\d+[)]`, "[]byte"},
	{`^(decimal)[(]\d+,\d+[)]`, "float64"},
	{`^(mediumint)[(]\d+[)]`, "string"},
	{`^(double)[(]\d+,\d+[)]`, "float64"},
	{`^(float)[(]\d+,\d+[)]`, "float64"},
	{`^(datetime)[(]\d+[)]`, "time.Time"},
	{`^(bit)[(]\d+[)]`, "[]uint8"},
	{`^(text)[(]\d+[)]`, "string"},
	{`^(integer)[(]\d+[)]`, "int"},
	{`^(timestamp)[(]\d+[)]`, "time.Time"},
	{`^(geometry)[(]\d+[)]`, "[]byte"},
}

func NewDBModel(info *DBInfo) *DBModel {
	return &DBModel{DBInfo: info}
}

func (m *DBModel) Connect() error {
	var err error
	s := "%s:%s@tcp(%s)/information_schema?" +
		"charset=%s&parseTime=True&loc=Local"
	dsn := fmt.Sprintf(
		s,
		m.DBInfo.UserName,
		m.DBInfo.Password,
		m.DBInfo.Host,
		m.DBInfo.Charset,
	)
	m.DBEngine, err = sql.Open(m.DBInfo.DBType, dsn)
	if err != nil {
		return err
	}

	return nil
}

// GetColumns 获取表中列的信息
func (m *DBModel) GetColumns(dbName, tableName string) ([]*TableColumn, error) {
	// use information_schema;
	// SELECT * FROM COLUMNS WHERE TABLE_SCHEMA = '数据库名称' and TABLE_NAME = '数据表名称';
	query := "SELECT COLUMN_NAME, DATA_TYPE, IS_NULLABLE, COLUMN_KEY, " +
		"COLUMN_TYPE, COLUMN_COMMENT, COALESCE(COLUMN_DEFAULT, ''), EXTRA " +
		"FROM COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? "
	rows, err := m.DBEngine.Query(query, dbName, tableName)
	if err != nil {
		return nil, err
	}
	if rows == nil {
		return nil, errors.New("没有数据")
	}
	defer rows.Close()

	var columns []*TableColumn
	for rows.Next() {
		var column TableColumn
		err := rows.Scan(
			&column.ColumnName,
			&column.DataType,
			&column.IsNullable,
			&column.ColumnKey,
			&column.ColumnType,
			&column.ColumnComment,
			&column.ColumnDefault,
			&column.Extra,
		)
		if err != nil {
			return nil, err
		}

		columns = append(columns, &column)
	}

	return columns, nil
}
