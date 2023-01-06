NAME=warning-suppressor
EXE=$(NAME).exe
VERSION=0.8
VER=08
DIST=$(NAME)-$(VERSION)
ZIP=$(NAME)$(VER).zip

$(EXE): main.go config.go
	go build -ldflags="-s -w" -trimpath

.PHONY: release
release: $(EXE)
	rm -rf $(DIST)
	mkdir $(DIST)
	cp $(EXE) README.md sample-config.yml $(DIST)
	7z a $(ZIP) $(DIST)
	rm -rf $(DIST)

