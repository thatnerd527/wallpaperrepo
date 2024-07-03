@echo off
cd .\WallpaperUI\WallpaperUI
cmd /c npx tailwindcss -i .\wwwroot\css\app.css -o .\wwwroot\css\app.min.css
cd ..
dotnet publish --output dist --configuration Debug
cd .. && exit