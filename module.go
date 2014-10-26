package conf

type module struct {
	Name string
	Port int
}

type modules map[string]interface{}
