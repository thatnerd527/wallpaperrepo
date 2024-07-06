@echo off
echo Step -1: Removing output folders
rmdir /s /q .\packaged

echo Step 0: Building all subprojects.
call buildall.bat

echo Step 1: Building installer from source.
cd .\wallpaperInstaller
go get
go build -ldflags "-H windowsgui"
rem go build
cd ..
echo Step 2: Packaging application.
.\tools\7z\7za.exe a -mmt=on .\packaged\install.zip .\output\*
echo Step 3: Copying.
xcopy .\wallpaperInstaller\wallpaperInstaller.exe .\packaged
xcopy .\installsplash.bapng .\packaged
mkdir .\packaged\tools
xcopy .\wallpaperInstaller\tools .\packaged\tools /s /e /y