echo Step 1: NPM Install in paralell
(
    start call .\scripts\APISetup.bat
    start call .\scripts\BootstrapSetup.bat
    start call .\scripts\PopupSetup.bat
    start call .\scripts\UIServerSetup.bat
    start call .\scripts\UISetup.bat
) | set /P "="