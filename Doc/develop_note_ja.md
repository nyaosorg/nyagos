- **lua53.dll のかわりに Gopher-Lua を採用** #300
    - 旧来の lua53.dll 版 nyagos.exe は `cd mains ; go build` でビルド可能
    - Lua無し版 nyagos.exe を `cd ngs ; go build` でビルド可能
- `nyagos.option.cleanup_buffer` を追加(デフォルトは false)。true の場合、一行入力の前にコンソールバッファをクリアする
- `set -o OPTION_NAME` と `set +o OPTION_NAME` を新設(`nyagos.option.OPTION_NAME=` on Lua と等価)
