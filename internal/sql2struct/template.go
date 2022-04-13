package sql2struct

import (
	"fmt"
	"os"
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

func (t *StructTemplate) AssemblyColumns(tbColumns []*TableColumn) []*StructColumn {
	tplColumns := make([]*StructColumn, 0, len(tbColumns))
	for _, column := range tbColumns {
		tag := fmt.Sprintf("`"+"json:"+"\"%s\""+"`", column.ColumnName) // 标签
		tplColumns = append(tplColumns, &StructColumn{
			Name:         column.ColumnName,
			Type:         DBTypeToStructType[column.DataType],
			Tag:          tag,
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
