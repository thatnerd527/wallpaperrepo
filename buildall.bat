@echo off
echo Step 0: Cleaning up output folder
rmdir /s /q .\output


echo Step 1: Building all projects in parallel
(
    start call .\scripts\APIbuild.bat
    start call .\scripts\BootstrapBuild.bat
    start call .\scripts\PopupBuild.bat
    start call .\scripts\UIServerBuild.bat
    start call .\scripts\UIBuild.bat
) | set /P "="
echo All projects built

echo Step 2: Copying files to output folder
mkdir .\output
echo Copy UI Server
echo touch > .\output\wallpaperuiserver.exe
xcopy .\wallpaperUIServer\wallpaperuiserver.exe .\output /y
echo Copy Electron App
mkdir .\output\popupapp
xcopy .\WallpaperPopupApp\out\wallpaperpopupapp-win32-x64\* .\output\popupapp /s /e /y
echo Copy Blazor UI
mkdir .\output\public
xcopy .\WallpaperUI\dist .\output\public /s /e /y
echo Copy API
mkdir .\output\public\wwwroot\javascript\API
xcopy .\wallpaperAPI\dist .\output\public\wwwroot\javascript\API /s /e /y
echo Copy Bootstrap
mkdir .\output\bootstrap
xcopy .\wallpaperbootstrap\dist .\output\bootstrap /s /e /y

echo Step 3: Done
pause