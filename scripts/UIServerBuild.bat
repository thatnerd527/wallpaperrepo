@echo off
cd .\wallpaperUIServer
go build -ldflags "-H windowsgui"
cd ..
.\tools\rcedit-x64.exe .\wallpaperUIServer\wallpaperuiserver.exe --set-icon .\icon.ico
.\tools\rcedit-x64.exe .\wallpaperUIServer\wallpaperuiserver.exe --set-version-string "FileDescription" "A daemon process that is used to support all of the operations of the Wallpaper System"
.\tools\rcedit-x64.exe .\wallpaperUIServer\wallpaperuiserver.exe --set-version-string "ProductName" "Wallpaper UI Server"
exit