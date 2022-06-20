module newApple

go 1.18

require (
	github.com/DataDog/go-python3 v0.0.0-20211102160307-40adc605f1fe
	github.com/Djoulzy/Tools/clog v0.0.0-20220609190146-71af779f6ddc
	github.com/Djoulzy/Tools/confload v0.0.0-20220609190146-71af779f6ddc
	github.com/Djoulzy/emutools/mem v0.0.0-20220620100952-42aec5b63f5b
	github.com/Djoulzy/emutools/mos6510 v0.0.0-20220620100952-42aec5b63f5b
	github.com/Djoulzy/emutools/render v0.0.0-20220620100952-42aec5b63f5b
	github.com/mattn/go-tty v0.0.4
)

require (
	github.com/Djoulzy/emutools/charset v0.0.0-20220620100952-42aec5b63f5b // indirect
	github.com/albenik/bcd v0.0.0-20170831201648-635201416bc7 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-ini/ini v1.66.6 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/veandco/go-sdl2 v0.4.24 // indirect
	golang.org/x/exp v0.0.0-20220613132600-b0d781184e0d // indirect
	golang.org/x/image v0.0.0-20220617043117-41969df76e82 // indirect
	golang.org/x/sys v0.0.0-20220615213510-4f61da869c0c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/Djoulzy/Tools/clog v0.0.0-20220429054701-4c221b41ecdf => ../github.com/Djoulzy/Tools/clog
	github.com/Djoulzy/emutools/mem v0.0.0-20220618111135-6d466b4f81cf => ../github.com/Djoulzy/emutools/mem
	github.com/Djoulzy/emutools/mos6510 v0.0.0-20220618111135-6d466b4f81cf => ../github.com/Djoulzy/emutools/mos6510
	github.com/Djoulzy/emutools/render v0.0.0-20220618111135-6d466b4f81cf => ../github.com/Djoulzy/emutools/render
)
