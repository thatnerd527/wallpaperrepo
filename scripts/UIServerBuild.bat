@echo off
cd .\wallpaperUIServer
go build -ldflags "-H windowsgui"
cd .. && exit