MD_EN := $(filter-out index%,$(filter-out _%.md,$(wildcard *_en.md)))
MD_JA := $(filter-out index%,$(filter-out _%.md,$(wildcard *_ja.md)))

HTML_EN := $(patsubst %.md,../docs/%.html,$(MD_EN))
HTML_JA := $(patsubst %.md,../docs/%.html,$(MD_JA))

ALL := $(HTML_EN) $(HTML_JA) ../docs/index.html ../docs/index_ja.html
all: $(ALL)

../docs/%_en.html: %_en.md
	gawk 'match($$0,/^##* */) { print substr($$0,RLENGTH+1); nextfile }' $< | minipage -title-from-file - -sidebar _Sidebar_en.md _Header.md $< > $@

../docs/%_ja.html: %_ja.md
	gawk 'match($$0,/^##* */) { print substr($$0,RLENGTH+1); nextfile }' $< | minipage -title-from-file - -sidebar _Sidebar_ja.md _Header.md $< > $@

../docs/index.html:    index.md
	minipage -title "NYAOS.ORG/NYAGOS" -sidebar _Sidebar_en.md _Header.md $< > $@

../docs/index_ja.html: index_ja.md
	minipage -title "NYAOS.ORG/NYAGOS" -sidebar _Sidebar_ja.md _Header.md $< > $@

clean:
	rm $(ALL)
