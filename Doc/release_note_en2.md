English / [Japanese](release_note_ja2.md)

Changes from master to second branch
========================================
* Remove built-in command `sudo`
* Add built-in command `more` (support color and unicode)
* readline: support C-q,C-v (`QUOTED_INSERT`)
* pwd: add options -L(use PWD from environment) and -P(avoid all symlinks)
* Output `nyagos.dump` if panic occurs.
* `__du__` : implemented the prototype of du
* `rmdir` prints the progress as before.
