package main

//go:generate go run merge_yamls.go
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=./oapi-config.yaml ../api/gen.api.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=./cfg-models.yaml ../api/gen.api.yaml

func main() {}
