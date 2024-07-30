import React, { useEffect } from "react";
import { PlayerRef } from "@remotion/player";
import { interpolate, Sequence, useCurrentFrame, Easing } from "remotion";
import { useMemo } from "react";
import logo from "../../public/applogo.svg";
import { InstallerData } from "../Datas";
export function InstallingSequence(props: {
  playerRef: React.RefObject<PlayerRef> | null;
  playbackRate: number;
  setPlaybackRate: (rate: number) => void;
  data: InstallerData;
  setData: (data: InstallerData) => void;
}) {
  const applogo = useMemo(() => logo, []);
  const frame = useCurrentFrame();
  const smalleningcurve = useMemo(() => Easing.bezier(0, 1.21, 0.83, 0.92), []);
  useEffect(() => {
    if (frame == 156 && props.data.pressedbutton) {
      props.playerRef?.current?.pause();
      props.playerRef?.current?.seekTo(281);
      props.setPlaybackRate(1);
      //props.playerRef?.current?.play();
    }
  }, [frame]);
  return (
    <Sequence from={281} durationInFrames={60}>
      <div className="w-full h-full flex flex-col justify-start items-center">
        <img
          src={applogo}
          style={{
            width: "160px",
            height: "160px",
            position: "relative",
            marginTop: `25%`,
          }}
        />
        <div className="h-4"></div>
        <h1>Installing</h1>
      </div>
    </Sequence>
  );
}
