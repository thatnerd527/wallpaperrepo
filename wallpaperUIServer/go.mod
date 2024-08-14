module wallpaperuiserver

go 1.21.3

replace wallpaperuiserver/protocol => ./protocol

require (
	github.com/elliotchance/pie/v2 v2.8.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.2
	github.com/orcaman/concurrent-map/v2 v2.0.1
)

require github.com/golang/protobuf v1.5.4 // indirect

require (
	github.com/TheTitanrain/w32 v0.0.0-20180517000239-4f5cfb03fabf // indirect
	github.com/getlantern/context v0.0.0-20190109183933-c447772a6520 // indirect
	github.com/getlantern/errors v0.0.0-20190325191628-abdb3e3e36f7 // indirect
	github.com/getlantern/golog v0.0.0-20190830074920-4ef2e798c2d7 // indirect
	github.com/getlantern/hex v0.0.0-20190417191902-c6586a6fe0b7 // indirect
	github.com/getlantern/hidden v0.0.0-20190325191715-f02dbb02be55 // indirect
	github.com/getlantern/ops v0.0.0-20190325191751-d70cb0d6f85f // indirect
	github.com/getlantern/systray v1.2.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/leonelquinteros/gorand v1.0.2 // indirect
	github.com/mitchellh/go-ps v1.0.0 // indirect
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c // indirect
	github.com/sqweek/dialog v0.0.0-20240226140203-065105509627 // indirect
	golang.org/x/exp v0.0.0-20220321173239-a90fa8a75705 // indirect
	golang.org/x/net v0.23.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	wallpaperuiserver/protocol v0.0.0-00010101000000-000000000000
)
