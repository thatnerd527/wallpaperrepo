import { useState } from 'react'
import { Player } from '@remotion/player'
import UpdateElectron from '@/components/update'
import logoVite from './assets/logo-vite.svg'
import logoElectron from './assets/logo-electron.svg'
import { Player2 } from './Player'
import './App.css'
import { InstallerFrame } from "../remotion2/Composition";
import {
  FluentProvider,
  webLightTheme,
  Button,
  webDarkTheme,
} from "@fluentui/react-components";

function App() {
  const [count, setCount] = useState(0)
  return (
    <div className="w-full h-full flex flex-col" style={{
      backgroundColor: 'rgba(0, 0, 0, 0.5)'
    }}>
      <div
        className="w-full min-h-14 p-4 flex flex-row justify-end items-center align-middle"
        style={
          {
            backgroundColor: "rgba(0, 0, 0, 0.2)",
            WebkitAppRegion: "drag",
          } as any
        }
      >
        <div
          style={
            {
              WebkitAppRegion: "no-drag",
            } as any
          }
        >
          <img src={logoVite} alt="logo-vite" className="w-12 h-12" />
        </div>
      </div>
      <FluentProvider theme={webDarkTheme} className='w-full h-full' style={{
        backgroundColor: 'transparent',
      }}>
        <Player2 />
      </FluentProvider>
    </div>
  );
}

export default App