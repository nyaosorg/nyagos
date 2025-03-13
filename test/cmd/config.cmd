@setlocal
@set "PROMPT=$G$S"
"%~dp0..\..\nyagos" --norc    --glob -e "if nyagos.option.glob == false then print('[NG]') ; os.exit(1) else print('[OK]') end" || exit /b 1
"%~dp0..\..\nyagos" --norc --no-glob -e "if nyagos.option.glob == true  then print('[NG]') ; os.exit(1) else print('[OK]') end" || exit /b 1
@endlocal
