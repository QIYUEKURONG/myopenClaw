package tools

type Tool interface {
	//名字
	Name() string
	//工具的描述
	Description() string
	//具体的执行
	Execute(args map[string]interface{}) (string, error)
	//参数的转化
	Parameters() map[string]interface{}
}
