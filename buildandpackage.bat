@echo off
echo Step -1: Removing output folders
rmdir /s /q .\packaged

echo Step 0: Building all subprojects.
call buildall.bat

echo Step 1: Building installer and uninstaller from source.
cd .\wallpaperInstaller
go get
go build -ldflags "-H windowsgui"
go build -ldflags "-H windowsgui" -tags uninstall -o uninstall.exe
rem go build
cd ..
echo Step 2: Packaging application.
xcopy .\wallpaperInstaller\uninstall.exe .\output
.\tools\7z\7za.exe a -mmt=on .\packaged\install.zip .\output\*
echo Step 3: Copying.
xcopy .\wallpaperInstaller\wallpaperInstaller.exe .\packaged
xcopy .\wallpaperInstaller\registrytemplate.reg .\packaged
xcopy .\installsplash.bapng .\packaged
mkdir .\packaged\tools
xcopy .\wallpaperInstaller\tools .\packaged\tools /s /e /y