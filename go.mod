module newApple

go 1.18

require (
	github.com/DataDog/go-python3 v0.0.0-20211102160307-40adc605f1fe
	github.com/Djoulzy/Tools/clog v0.0.0-20220609190146-71af779f6ddc
	github.com/Djoulzy/Tools/confload v0.0.0-20220609190146-71af779f6ddc
	github.com/Djoulzy/emutools/mem v0.0.0-20220728172040-42916fc274d5
	github.com/Djoulzy/emutools/mem2 v0.0.0-20220728172040-42916fc274d5
	github.com/Djoulzy/emutools/mos6510 v0.0.0-20220728172040-42916fc274d5
	github.com/Djoulzy/emutools/render v0.0.0-20220728172040-42916fc274d5
	github.com/Djoulzy/godsk v0.0.0-20220705093616-8da8bd02989a
	github.com/Djoulzy/gowoz v0.0.0-20220707173953-ab69067b349a
	github.com/mattn/go-tty v0.0.4
)

require (
	github.com/Djoulzy/emutools/charset v0.0.0-20220728172040-42916fc274d5 // indirect
	github.com/albenik/bcd v0.0.0-20170831201648-635201416bc7 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-ini/ini v1.66.6 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/veandco/go-sdl2 v0.4.24 // indirect
	golang.org/x/exp v0.0.0-20220722155223-a9213eeb770e // indirect
	golang.org/x/image v0.0.0-20220722155232-062f8c9fd539 // indirect
	golang.org/x/sys v0.0.0-20220728004956-3c1f35247d10 // indirect
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
