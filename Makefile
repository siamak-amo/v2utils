GOCC = go

CMD_DIR = ./cmd
PKG_DIR = ./pkg
MAIN_CMD = v2utils.go

main:
	$(GOCC) build -v $(CMD_DIR)/$(MAIN_CMD)

test:
	$(GOCC) test -v $(PKG_DIR)
