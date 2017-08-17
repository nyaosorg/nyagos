[English](release_note_en2.md) / Japanese

master ブランチから second ブランチへの変更
===========================================

* 内蔵コマンドの sudo を削除
* 内蔵コマンド more を追加(カラー & utf8 サポート)
* 一行入力で `C-q`,`C-v` をサポート(`QUOTED_INSERT`)
* 内蔵コマンド pwd に -P(全てのリンクをたどる) ,-L(環境からPWDを得る) を追加
* パニックが発生した時、nyagos.dump を出力するようにした
* `diskused`: du ライクな新コマンド
* `rmdir` : 進捗を表示する仕様を復活させた
* `diskfree`: df ライクな新コマンド
