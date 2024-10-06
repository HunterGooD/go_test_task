package main

//go:generate swag init -g cmd/music_service/main.go
//go:generate oapi-codegen --config=config/music_info_cfg.yaml api/music_info.yaml
