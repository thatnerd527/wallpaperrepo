import { PlayerRef } from "@remotion/player";
import React, { useState } from "react";
import {
  AbsoluteFill,
  Easing,
  interpolate,
  Sequence,
  useCurrentFrame,
} from "remotion";
import logo from "../public/applogo.svg";
import { Opening } from "./Sequences/Opening";
import { MainScreen } from "./Sequences/MainScreen";
import { InstallingSequence } from "./Sequences/Installing";
import {InstallerData} from './Datas';
export const InstallerFrame = (props: {
  playerRef: React.RefObject<PlayerRef> | null;
  playbackRate: number;
  setPlaybackRate: (rate: number) => void;
}) => {
  const [data, setData] = useState(new InstallerData)
  return (
    <>
      <Opening data={data} setData={setData} />
      <MainScreen
        playerRef={props.playerRef}
        playbackRate={props.playbackRate}
        setPlaybackRate={props.setPlaybackRate}
        data={data}
        setData={setData}
      />
      <InstallingSequence data={data} playbackRate={props.playbackRate} playerRef={props.playerRef} setData={setData} setPlaybackRate={props.setPlaybackRate}/>
    </>
  );
};
