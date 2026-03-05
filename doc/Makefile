ifeq ($(OS),Windows_NT)
    DEL=del
else
    DEL=rm
endif

MINIPAGE=go run github.com/hymkor/minipage@latest

MD_EN := $(filter-out index%,$(filter-out _%.md,$(wildcard *_en.md)))
MD_JA := $(filter-out index%,$(filter-out _%.md,$(wildcard *_ja.md)))

HTML_EN := $(patsubst %.md,../docs/%.html,$(MD_EN))
HTML_JA := $(patsubst %.md,../docs/%.html,$(MD_JA))

ALL := $(HTML_EN) $(HTML_JA) ../docs/index.html ../docs/index_ja.html ../docs/CHANGELOG.html ../docs/CHANGELOG_ja.html
all: $(ALL)

../docs/%_en.html: %_en.md
	$(MINIPAGE) -title "NYAOS.ORG/NYAGOS - $*" -sidebar _Sidebar_en.md _Header.md $< > $@

../docs/%_ja.html: %_ja.md
	$(MINIPAGE) -title "NYAOS.ORG/NYAGOS - $*" -sidebar _Sidebar_ja.md _Header.md $< > $@

../docs/index.html:    index.md
	$(MINIPAGE) -title "NYAOS.ORG/NYAGOS" -sidebar _Sidebar_en.md _Header.md $< > $@

../docs/index_ja.html: index_ja.md
	$(MINIPAGE) -title "NYAOS.ORG/NYAGOS" -sidebar _Sidebar_ja.md _Header.md $< > $@

../docs/CHANGELOG.html: CHANGELOG.md
	$(MINIPAGE) -title "NYAOS.ORG/NYAGOS - Changelog" -sidebar _Sidebar_en.md _Header.md $< > $@

../docs/CHANGELOG_ja.html: CHANGELOG_ja.md
	$(MINIPAGE) -title "NYAOS.ORG/NYAGOS - Changelog" -sidebar _Sidebar_ja.md _Header.md $< > $@

clean:
	$(DEL) $(ALL)
