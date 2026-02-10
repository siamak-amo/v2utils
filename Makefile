GOCC = go

CMD_DIR = ./cmd
SRC_DIRS = ./pkg ./getopt

BUILD_DIR = ./build
BINS = v2utils

main: $(BINS)

v2utils:
	$(GOCC) build -o $(BUILD_DIR)/v2utils -v $(CMD_DIR)/$@

debug:
	$(GOCC) build -tags debug -o v2utils -v $(CMD_DIR)/

test:
	$(GOCC) test -v $(SRC_DIRS)
