package main

//go:generate go run merge_yamls.go
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=./oapi-config.yaml ../api/gen.api.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=./cfg-models.yaml ../api/gen.api.yaml
// Generation below is for making easier generating client for manual tests
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../api/manualtests/cfg-client.yaml ../api/gen.api.yaml
