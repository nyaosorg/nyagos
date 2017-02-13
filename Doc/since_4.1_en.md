English / [Japanese](./since_4.1_ja.md)

Since 4.1
=========

NYAGOS 4.0 has the problem that a panic sometimes occurs because 
multi-goroutines call only one Lua-instance.

On NYAGOS 4.1, I tryied to prevent from panic and to make stable 
by creating Lua-instances for each lua-call and not allowing 
multi-goroutines use same one.

But, it has the new problem that the functions assigned to nyagos[] ,
can not see variables set on .nyagos including themselves.
Because Lua-instances have their own variable-areas not seen by 
other instances.

To share data between the different Lua-instance, On 4.1, the values 
assigned on the global table nyagos[] and share[] are linked to 
Go's memory. They are able to be accessed by all Lua-instance.

It breaks the compatibility on Lua scripts. To let Lua scripts for 4.0
run on 4.1, these modifing are required.

- Value assigned to global variable should be assigned to `share[]`.
    - Values and functions in `share[]` are able to be access from
      all lua-instances.
    - nyagos.exe can find modifying the member of share[] and nyagos[] only.
        - OK: `share.foo = { '1','2','3' }`
        - NG: `share.foo[1] = 'x'`
        - OK: `local t=share.foo ; t[1] = 'x' ; share.foo = t`
    - Do not assign closure to `nyagos.alias[]` ! 
      The code in the function can not access the bind variables.
