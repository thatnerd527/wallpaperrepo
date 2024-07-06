@echo off
cd .\WallpaperPopupApp
cmd /c npx vite build && cmd /c npx electron-forge package
cd ..
timeout /t 4
.\tools\rcedit-x64.exe .\WallpaperPopupApp\out\wallpaperpopupapp-win32-x64\wallpaperpopupapp.exe --set-icon .\icon.ico
.\tools\rcedit-x64.exe .\WallpaperPopupApp\out\wallpaperpopupapp-win32-x64\wallpaperpopupapp.exe --set-version-string "FileDescription" "A small application that is used to host pop ups and typing requests from backgrounds, panels and the wallpaper system. "
.\tools\rcedit-x64.exe .\WallpaperPopupApp\out\wallpaperpopupapp-win32-x64\wallpaperpopupapp.exe --set-version-string "ProductName" "Wallpaper Popup App"
exit                                                                                            
