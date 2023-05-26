module newApple

go 1.20

replace (
	github.com/Djoulzy/emutools/mos6510 => ../github.com/Djoulzy/emutools/mos6510
	github.com/Djoulzy/mmu => ../github.com/Djoulzy/mmu
)

require (
	github.com/Djoulzy/Tools/clog v0.0.0-20220609190146-71af779f6ddc
	github.com/Djoulzy/Tools/confload v0.0.0-20220609190146-71af779f6ddc
	github.com/Djoulzy/emutools/mem v0.0.0-20221020174520-027303f02bbe
	github.com/Djoulzy/emutools/mos6510 v0.0.0-00010101000000-000000000000
	github.com/Djoulzy/emutools/render v0.0.0-20221020174520-027303f02bbe
	github.com/Djoulzy/godsk v0.0.0-20221012182138-3f22e902d449
	github.com/Djoulzy/gowoz v0.0.0-20221012182153-b80301c0b697
	github.com/Djoulzy/mmu v0.0.0-20221015154434-3927fedd1199
	github.com/mattn/go-tty v0.0.5
)

require (
	github.com/Djoulzy/emutools/charset v0.0.0-20221020174520-027303f02bbe // indirect
	github.com/albenik/bcd v0.0.0-20170831201648-635201416bc7 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/veandco/go-sdl2 v0.4.35 // indirect
	golang.org/x/exp v0.0.0-20230519143937-03e91628a987 // indirect
	golang.org/x/image v0.7.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
