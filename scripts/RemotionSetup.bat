@echo off
cd .\remotion-fork
cmd /c bun install
cd .\packages\core
cmd /c bun install
cd ..\..
cd .. && exit.