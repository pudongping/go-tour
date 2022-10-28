package sql2struct

import (
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/pudongping/go-tour/internal/word"
)

// {{ if ne "" }} default_value： {{.DefaultValue}}{{ end }}
const strcutTpl = `type {{.TableName | ToCamelCase}} struct {
{{range .Columns}}	{{ $length := len .Comment}} {{ if gt $length 0 }}// {{.Comment}} {{else}}// {{.Name}} {{ end }} {{.Extra}} {{.ColumnKey}} {{.ColumnType}} is_nullable: {{.IsNullable}} {{ if ne .DefaultValue "" }} default_value: {{.DefaultValue}}{{ end }}
	{{ $typeLen := len .Type }} {{ if gt $typeLen 0 }}{{.Name | ToCamelCase}}	{{.Type}}	{{.Tag}}{{ else }}{{.Name}}{{ end }}
{{end}}}

func (model {{.TableName | ToCamelCase}}) TableName() string {
	return "{{.TableName}}"
}` + "\n"

type StructTemplate struct {
	strcutTpl string
}

type StructColumn struct {
	Name         string // 列名
	Type         string // 列的数据类型，仅包含类型信息
	Tag          string // 标签
	Comment      string // 列的注释信息
	IsNullable   string // 列是否允许为 null
	ColumnType   string // 列的数据类型，包含类型名称和可能的其他信息。eg：精度，长度，是否无符号等
	ColumnKey    string // 列是否被索引
	DefaultValue string // 列的默认值
	Extra        string // 列的额外信息。eg：pri、uni
}

type StructTemplateDB struct {
	TableName string
	Columns   []*StructColumn
}

func NewStructTemplate() *StructTemplate {
	return &StructTemplate{strcutTpl: strcutTpl}
}

// AssemblyColumns 组装所需数据
func (t *StructTemplate) AssemblyColumns(tbColumns []*TableColumn, orm string) []*StructColumn {
	tplColumns := make([]*StructColumn, 0, len(tbColumns))
	for _, column := range tbColumns {
		tplColumns = append(tplColumns, &StructColumn{
			Name:         column.ColumnName,
			Type:         t.BuildGoType(column),
			Tag:          t.BuildTag(column, orm),
			Comment:      column.ColumnComment,
			IsNullable:   column.IsNullable,
			ColumnType:   column.ColumnType,
			ColumnKey:    column.ColumnKey,
			DefaultValue: column.ColumnDefault,
			Extra:        column.Extra,
		})
	}

	return tplColumns
}

// BuildTag 构建结构体字段标签
func (t *StructTemplate) BuildTag(column *TableColumn, orm string) string {
	jsonTag := `json:"` + column.ColumnName + `"`
	var ormTag string
	if orm == "gorm" {
		ormTag = t.gormTag(column)
	} else if orm == "xorm" {
		ormTag = t.xormTag(column)
	}

	return "`" + ormTag + " " + jsonTag + "`"
}

// xormTag
// xorm 标签定义 https://www.kancloud.cn/xormplus/xorm/167137
func (t *StructTemplate) xormTag(column *TableColumn) string {
	var ormTag string
	// 因为 xorm 不需要 `unsigned` 属性作为标签的一部分，因此去掉
	// 这里默认只展示类似于 `int(11)` 因为有可能会出现 `int(11) unsigned`
	columnType := strings.Split(column.ColumnType, " ")[0]

	ormTag += `xorm:"` + columnType + " "
	if "PRI" == column.ColumnKey {
		ormTag += "pk "
	}

	if "auto_increment" == column.Extra {
		ormTag += "autoincr "
	}

	if "NO" == column.IsNullable {
		ormTag += "notnull "
	}

	if "UNI" == column.ColumnKey {
		ormTag += "unique "
	}

	if column.ColumnName == "created_at" || column.ColumnName == "create_time" {
		ormTag += "created "
	}
	if column.ColumnName == "updated_at" || column.ColumnName == "update_time" {
		ormTag += "updated "
	}
	if column.ColumnName == "deleted_at" || column.ColumnName == "delete_time" {
		ormTag += "deleted "
	}

	ormTag += "'" + column.ColumnName + `'"`
	return ormTag
}

// gormTag
// gorm 标签定义 https://gorm.io/zh_CN/docs/models.html
func (t *StructTemplate) gormTag(column *TableColumn) string {
	var ormTag string
	ormTag += `gorm:"column:` + column.ColumnName + ";"
	if "PRI" == column.ColumnKey {
		ormTag += "primaryKey;"
	}

	if "UNI" == column.ColumnKey {
		ormTag += "unique;"
	}

	if "auto_increment" == column.Extra {
		ormTag += "autoIncrement;"
	}

	if "NO" == column.IsNullable {
		ormTag += "not null;"
	}

	ormTag += `"`
	return ormTag
}

// BuildGoType 映射出 go 语言的数据类型
func (t *StructTemplate) BuildGoType(column *TableColumn) string {
	// 先精确匹配
	if columnType, ok := DBTypeToStructType[column.DataType]; ok {
		return columnType
	}

	// 模糊正则匹配
	for _, item := range TypeMysqlMatchList {
		if ok, _ := regexp.MatchString(item.Key, column.DataType); ok {
			return item.Value
		}
	}

	return "string"
}

func (t *StructTemplate) Generate(tableName string, tplColumns []*StructColumn) error {
	tpl := template.Must(template.New("sql2struct").Funcs(template.FuncMap{
		"ToCamelCase": word.UnderscoreToUpperCamelCase,
	}).Parse(t.strcutTpl))

	tplDB := StructTemplateDB{
		TableName: tableName,
		Columns:   tplColumns,
	}
	err := tpl.Execute(os.Stdout, tplDB)
	if err != nil {
		return err
	}

	return nil
}
