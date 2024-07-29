@echo off
cd .\remotion-fork
cmd /c npm run build
cd .\packages\core
echo Building core..
cmd /c npm run build
echo Building with tsc..
cmd /c bun run tsc
cd ..\player
echo Building player..
cmd /c npm run build
echo Building with tsc..
cmd /c bun run tsc
cd ..\..
cd .. && exit.