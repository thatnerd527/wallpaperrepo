import { Player, PlayerRef } from "@remotion/player";
import { InstallerFrame } from "../remotion2/Composition";
import { useEffect, useRef, useState } from "react";
export function Player2() {
  const playerRef = useRef<PlayerRef>(null);
  const [playbackRate, setPlaybackRate] = useState(1);
  return (
    <Player
      ref={playerRef}
      responsiveSize={true}
      component={InstallerFrame}
      inputProps={{
        playerRef: playerRef,
        setPlaybackRate: setPlaybackRate,
        playbackRate: playbackRate,
      }}
      playbackRate={playbackRate}
      durationInFrames={900}
      fps={60}
      loop
      autoPlay
      compositionHeight={450}
      compositionWidth={300}
      style={{ width: "100%", height: "100%" }}
    />
  );
}
