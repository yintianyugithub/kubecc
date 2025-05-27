.PHONY: build clean run

OUTPUT_DIR = ./api

# 生成代码
build-api:
	goctl api go -api ${OUTPUT_DIR}/api.api -dir ${OUTPUT_DIR} -style go_zero --home ./core/tpl/1.6.4
