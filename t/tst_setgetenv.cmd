@more +1 "%~f0" | "%~dp0..\nyagos" - & exit /b 0
lua_e "nyagos.setenv('AHAHA','成功') ; nyagos.write(nyagos.getenv('AHAHA'),'\n')"
lua_e "nyagos.env['AHAHA']='成功2' ; nyagos.write(nyagos.env['AHAHA'],'\n')"
