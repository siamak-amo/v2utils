GOCC = go

CMD_DIR = ./cmd
SRC_DIRS = ./pkg ./getopt
BINS = v2utils

main: $(BINS)

$(BINS):
	$(GOCC) build -o v2utils -v $(CMD_DIR)/$@

debug:
	$(GOCC) build -tags debug -o v2utils -v $(CMD_DIR)/

test:
	$(GOCC) test -v $(SRC_DIRS)
