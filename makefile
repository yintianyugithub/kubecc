.PHONY: build clean run

OUTPUT_DIR = ./api-greet

# 生成代码
build-api-greet:
	goctl api go -api ${OUTPUT_DIR}/greet.api -dir ${OUTPUT_DIR} -style go_zero
