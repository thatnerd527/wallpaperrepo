@echo off
cd .\WallpaperUI
dotnet publish --output dist --configuration Debug
cd .. && exit