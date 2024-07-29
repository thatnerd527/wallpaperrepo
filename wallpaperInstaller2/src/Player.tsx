import { Player, PlayerRef } from "@remotion/player";
import { InstallerFrame } from "../remotion2/Composition";
import { useEffect, useRef } from "react";
export function Player2() {
  const playerRef = useRef<PlayerRef>(null);
  useEffect(() => {
      const { current } = playerRef;
      if (!current) {
        return;
      }

      const listener = () => {
        console.log("started playing");
        let element = document.getElementsByClassName("__remotion-player")[0];
        console.log(current.getContainerNode())
        console.log(element);
       // debugger;
      };
      current.addEventListener("play", listener);
      current.play();
      return () => {
        current.removeEventListener("play", listener);
      };
    }, []);
  return (
    <Player
      ref={playerRef}
      responsiveSize={true}
      component={InstallerFrame}
      inputProps={{
        playerRef: playerRef,
      }}
      durationInFrames={900}
      fps={60}
      loop
      controls
      compositionHeight={450}
      compositionWidth={300}
      style={{ width: "100%", height: "100%" }}
    />
  );
}
