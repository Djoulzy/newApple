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
	github.com/Djoulzy/emutools/mos6510 v0.0.0-20240123173627-4140dd715cad
	github.com/Djoulzy/emutools/render v0.0.0-20240123173627-4140dd715cad
	github.com/Djoulzy/godsk v0.0.0-20230918100154-614368a3a7c0
	github.com/Djoulzy/gowoz v0.0.0-20230913121504-97232d4d9c93
	github.com/Djoulzy/mmu v0.0.0-20230605062009-e48b6d54957a
	github.com/mattn/go-tty v0.0.5
)

require (
	github.com/Djoulzy/emutools/charset v0.0.0-20240123173627-4140dd715cad // indirect
	github.com/albenik/bcd v0.0.0-20170831201648-635201416bc7 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/veandco/go-sdl2 v0.4.38 // indirect
	golang.org/x/exp v0.0.0-20240119083558-1b970713d09a // indirect
	golang.org/x/image v0.15.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
