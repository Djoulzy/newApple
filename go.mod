module newApple

go 1.21

replace (
	github.com/Djoulzy/emutools/mos6510 => ../github.com/Djoulzy/emutools/mos6510
	github.com/Djoulzy/emutools/render => ../github.com/Djoulzy/emutools/render
	github.com/Djoulzy/godsk => ../github.com/Djoulzy/godsk
	github.com/Djoulzy/gowoz => ../github.com/Djoulzy/gowoz
	github.com/Djoulzy/mmu => ../github.com/Djoulzy/mmu
)

require (
	github.com/Djoulzy/Tools/clog v0.0.0-20220609190146-71af779f6ddc
	github.com/Djoulzy/Tools/confload v0.0.0-20220609190146-71af779f6ddc
	github.com/Djoulzy/emutools/mos6510 v0.0.0-20230911073411-473d153118d9
	github.com/Djoulzy/emutools/render v0.0.0-20230911073411-473d153118d9
	github.com/Djoulzy/godsk v0.0.0-20230912184023-cc720e3e79d5
	github.com/Djoulzy/gowoz v0.0.0-20230912174927-2b179db0043b
	github.com/Djoulzy/mmu v0.0.0-20230605062009-e48b6d54957a
	github.com/mattn/go-tty v0.0.5
)

require (
	github.com/Djoulzy/emutools/charset v0.0.0-20230911073411-473d153118d9 // indirect
	github.com/albenik/bcd v0.0.0-20170831201648-635201416bc7 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/veandco/go-sdl2 v0.4.35 // indirect
	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
	golang.org/x/image v0.12.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
