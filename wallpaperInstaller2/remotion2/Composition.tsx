import { PlayerRef } from "@remotion/player";
import React from "react";
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
export const InstallerFrame = (props: {
  playerRef: React.RefObject<PlayerRef> | null;
}) => {
  const frame = useCurrentFrame();
  if (frame == 50 && props.playerRef != null) {
  }
  const applogo = logo;

  const openingcurve = Easing.bezier(0.02, 1.66, 0.47, 0.33);
  const smalleningcurve = Easing.bezier(0, 1.21, 0.83, 0.92);
  return (
    <>
      <Opening/>
      <MainScreen/>
    </>
  );
};
