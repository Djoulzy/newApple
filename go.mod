module newApple

go 1.19

require (
	github.com/Djoulzy/Tools/clog v0.0.0-20220609190146-71af779f6ddc
	github.com/Djoulzy/Tools/confload v0.0.0-20220609190146-71af779f6ddc
	github.com/Djoulzy/emutools/mem/v2 v2.0.0-20221020174520-027303f02bbe
	github.com/Djoulzy/emutools/mos6510 v0.0.0-20221020174520-027303f02bbe
	github.com/Djoulzy/emutools/render v0.0.0-20221020174520-027303f02bbe
	github.com/Djoulzy/godsk v0.0.0-20221012182138-3f22e902d449
	github.com/Djoulzy/gowoz v0.0.0-20221012182153-b80301c0b697
	github.com/mattn/go-tty v0.0.4
)

require (
	github.com/Djoulzy/emutools/charset v0.0.0-20221020174520-027303f02bbe // indirect
	github.com/Djoulzy/emutools/mem v0.0.0-20221020174520-027303f02bbe // indirect
	github.com/albenik/bcd v0.0.0-20170831201648-635201416bc7 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/veandco/go-sdl2 v0.4.30 // indirect
	golang.org/x/exp v0.0.0-20230127140709-cafedaf64729 // indirect
	golang.org/x/image v0.3.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/Djoulzy/emutools/mem v0.0.0-20221019172136-ccfe72c37753 => ../github.com/Djoulzy/emutools/mem
	github.com/Djoulzy/emutools/mem/v2 v2.0.0-20221019172136-ccfe72c37753 => ../github.com/Djoulzy/emutools/mem/v2
	github.com/Djoulzy/emutools/mos6510 v0.0.0-20221019172136-ccfe72c37753 => ../github.com/Djoulzy/emutools/mos6510
	github.com/Djoulzy/emutools/render v0.0.0-20221019172136-ccfe72c37753 => ../github.com/Djoulzy/emutools/render
)
