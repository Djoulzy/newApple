module newApple

go 1.18

require (
	github.com/DataDog/go-python3 v0.0.0-20211102160307-40adc605f1fe
	github.com/Djoulzy/Tools/clog v0.0.0-20220609190146-71af779f6ddc
	github.com/Djoulzy/Tools/confload v0.0.0-20220609190146-71af779f6ddc
	github.com/Djoulzy/emutools/mem v0.0.0-20220729172030-9232a1f7b306
	github.com/Djoulzy/emutools/mem2 v0.0.0-20220729172030-9232a1f7b306
	github.com/Djoulzy/emutools/mos6510 v0.0.0-20220729172030-9232a1f7b306
	github.com/Djoulzy/emutools/render v0.0.0-20220729172030-9232a1f7b306
	github.com/Djoulzy/godsk v0.0.0-20220705093616-8da8bd02989a
	github.com/Djoulzy/gowoz v0.0.0-20220707173953-ab69067b349a
	github.com/mattn/go-tty v0.0.4
)

require (
	github.com/Djoulzy/emutools/charset v0.0.0-20220729172030-9232a1f7b306 // indirect
	github.com/albenik/bcd v0.0.0-20170831201648-635201416bc7 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/veandco/go-sdl2 v0.4.25 // indirect
	golang.org/x/exp v0.0.0-20220930202632-ec3f01382ef9 // indirect
	golang.org/x/image v0.0.0-20220902085622-e7cb96979f69 // indirect
	golang.org/x/sys v0.0.0-20220928140112-f11e5e49a4ec // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/Djoulzy/Tools/clog v0.0.0-20220429054701-4c221b41ecdf => ../github.com/Djoulzy/Tools/clog
	github.com/Djoulzy/emutools/mem v0.0.0-20220728172040-42916fc274d5 => ../github.com/Djoulzy/emutools/mem
	github.com/Djoulzy/emutools/mem2 v0.0.0-20220728172040-42916fc274d5 => ../github.com/Djoulzy/emutools/mem2
	github.com/Djoulzy/emutools/mos6510 v0.0.0-20220728172040-42916fc274d5 => ../github.com/Djoulzy/emutools/mos6510
	github.com/Djoulzy/emutools/render v0.0.0-20220728172040-42916fc274d5 => ../github.com/Djoulzy/emutools/render
	github.com/Djoulzy/godsk v0.0.0-20220705093616-8da8bd02989a => ../github.com/Djoulzy/godsk
	github.com/Djoulzy/gowoz v0.0.0-20220707173953-ab69067b349a => ../github.com/Djoulzy/gowoz
)
