@echo off
cd .\tools\protoc\bin
echo Compiling with protoc.
mkdir ..\..\..\protodefs\proto\builtcs
mkdir ..\..\..\protodefs\proto\builtgo
.\bash.exe -c "./protoc.exe ../../../protodefs/proto/*.proto --csharp_out=../../../protodefs/proto/builtcs --proto_path=../../../protodefs/proto --go_out="../../../protodefs/proto/builtgo" --go_opt=paths=source_relative"
echo Copying protocols
mkdir ..\..\..\wallpaperAPI\src\protocols
rem xcopy ..\..\..\protodefs\proto\builtjs\* ..\..\..\wallpaperAPI\src\protocols /y /e
mkdir ..\..\..\WallpaperUI\WallpaperUI\Cs\Protocol
mkdir ..\..\..\wallpaperUIServer\protocol
xcopy ..\..\..\protodefs\proto\builtcs\* ..\..\..\WallpaperUI\WallpaperUI\Cs\Protocol /y /e
xcopy ..\..\..\protodefs\proto\builtgo\* ..\..\..\wallpaperUIServer\protocol /y /e

echo Compiling with pbjs and pbts.
cd ..\..\..\tools\pbjs
mkdir .\result
cmd /c npx pbjs -t static-module -w es6 -o .\result\protocol.js ..\..\protodefs\proto\application.proto
cmd /c npx pbts -o .\result\protocol.d.ts .\result\protocol.js

echo Copying protocols.
mkdir ..\..\WallpaperPopupApp\src\protocol
xcopy .\result\* ..\..\WallpaperPopupApp\src\protocol /y /e
cd ..\..\WallpaperPopupApp\subbuild
node .\build.js
cd ..\..\tools\pbjs
mkdir ..\..\wallpaperAPI\src\protocols
xcopy .\result\* ..\..\wallpaperAPI\src\protocols /y /e
cd ..\..\wallpaperAPI\subbuild
node .\build.js
echo All done.