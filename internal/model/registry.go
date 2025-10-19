package model

var models []any

func Register(m any) {
	models = append(models, m)
}

func Models() []any {
	return models
}
