@echo off
cd .\WallpaperPopupApp
cmd /c npx vite build && cmd /c npx electron-forge package
cd .. && exit
