package golang

// 报错行号+7
const goCodeTemplate = `// Generated by github.com/davyxu/protoplus
// DO NOT EDIT!
package {{.PackageName}}

import (
	"github.com/davyxu/protoplus/proto"
	"fmt"
)
var (
	_ *proto.Buffer
	_ fmt.Stringer
)

{{range $a, $enumobj := .Enums}}
type {{.Name}} int32
const (	{{range .Fields}}
	{{$enumobj.Name}}_{{.Name}} {{.Type}} = {{TagNumber $enumobj .}} {{end}}
)

var (
{{$enumobj.Name}}MapperValueByName = map[string]int32{ {{range .Fields}}
	"{{.Name}}": {{TagNumber $enumobj .}}, {{end}}
}

{{$enumobj.Name}}MapperNameByValue = map[int32]string{ {{range .Fields}}
	{{TagNumber $enumobj .}}: "{{.Name}}" , {{end}}
}
)

func (self {{$enumobj.Name}}) String() string {
	return {{$enumobj.Name}}MapperNameByValue[int32(self)]
}
{{end}}

{{range $a, $obj := .Structs}}
{{ObjectLeadingComment .}}
type {{.Name}} struct{	{{range .Fields}}
	{{GoFieldName .}} {{ProtoTypeName .}} {{GoStructTag .}}{{FieldTrailingComment .}} {{end}}
}

func (self *{{.Name}}) String() string { return fmt.Sprintf("%+v",*self) }

func (self *{{.Name}}) Size() (ret int) {
{{range .Fields}}
	{{if IsStructSlice .}}
	if len(self.{{GoFieldName .}}) > 0 {
		for _, elm := range self.{{GoFieldName .}} {
			ret += proto.SizeStruct({{TagNumber $obj .}}, elm)
		}
	}
	{{else}}
	ret += proto.Size{{CodecName .}}({{TagNumber $obj .}}, self.{{GoFieldName .}})
	{{end}}
{{end}}
	return
}

func (self *{{.Name}}) Marshal(buffer *proto.Buffer) error {
{{range .Fields}}
	{{if IsStructSlice .}}
		for _, elm := range self.{{GoFieldName .}} {
			proto.MarshalStruct(buffer, {{TagNumber $obj .}}, elm)
		}
	{{else}}	
		proto.Marshal{{CodecName .}}(buffer, {{TagNumber $obj .}}, self.{{GoFieldName .}})
	{{end}}
{{end}}
	return nil
}

func (self *{{.Name}}) Unmarshal(buffer *proto.Buffer, fieldIndex uint64, wt proto.WireType) error {
	switch fieldIndex {
	{{range .Fields}} case {{TagNumber $obj .}}: {{if IsStruct .}}
		self.{{GoFieldName .}} = new({{.Type}})
		return proto.UnmarshalStruct(buffer, wt, self.{{GoFieldName .}}) {{else if IsStructSlice .}}
		elm := new({{.Type}})
		if err := proto.UnmarshalStruct(buffer, wt, elm); err != nil {
			return err
		} else {
			self.{{GoFieldName .}} = append(self.{{GoFieldName .}}, elm)
			return nil
		}{{else}}
		return proto.Unmarshal{{CodecName .}}(buffer, wt, &self.{{GoFieldName .}}) {{end}}
	{{end}}
	}

	return proto.ErrUnknownField
}
{{end}}

`
