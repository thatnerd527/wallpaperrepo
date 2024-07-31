import React, { useCallback, useEffect, useState } from "react";
import { PlayerRef } from "@remotion/player";
import { interpolate, Sequence, useCurrentFrame, Easing } from "remotion";
import { useMemo } from "react";
import logo from "../../public/applogo.svg";
import { InstallerData } from "../Datas";
import { Lottie } from "@remotion/lottie";
import animationData from "../../public/animation.json";
import { AnimationItem } from "lottie-web";
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
  const anim = useMemo(() => animationData, []);

  const openingcurve = useMemo(() => Easing.bezier(0.02, 1.66, 0.47, 0.33), []);
  const [looped, setLooped] = useState(false)
  const [onetime, setOnetime] = useState(false)
  const [item, setItem] = useState<AnimationItem | null>(null);
  useEffect(() => {
    if (frame == 156 && props.data.pressedbutton) {
      props.playerRef?.current?.pause();
      props.playerRef?.current?.seekTo(281);
      props.setPlaybackRate(1);
      props.playerRef?.current?.play();
    }
    if (frame == 515 && !props.data.installfinished && !looped) {
      setLooped(true)
      console.log("Looped case1")
      if (item != null) {
        item.firstFrame = 0;
      }
      props.playerRef?.current?.seekTo(281);
    }
    if (frame == 641 && !props.data.installfinished && looped) {
      console.log("Looped case2");
      if (item != null) {
        console.log("inner trig")
        item.firstFrame = 0;
      }
      props.playerRef?.current?.seekTo(281);
    }
  }, [frame,item]);
  const onAnimationLoaded = useCallback((item: AnimationItem) => {
    setItem(item);
    if (!onetime) {
      console.log("onetime")
      setOnetime(true)
      item.firstFrame = 126;
    }
  }, [onetime,setItem]);
  return (
    <Sequence from={281} durationInFrames={360}>
      <div className="w-full h-full flex flex-col justify-center items-center">
        <Lottie
          animationData={anim}
          loop={true}

          onAnimationLoaded={onAnimationLoaded}

          style={{
            width: interpolate(frame, [0, 96], [0, 200], {
              easing: openingcurve,
              extrapolateLeft: "clamp",
              extrapolateRight: "clamp",
            }),
            height: interpolate(frame, [0, 96], [0, 200], {
              easing: openingcurve,
              extrapolateLeft: "clamp",
              extrapolateRight: "clamp",
            }),
            opacity: interpolate(frame, [0, 30], [0, 1], {
              easing: Easing.linear,
            }),
          }}
        ></Lottie>
        <h1>Installing</h1>
      </div>
    </Sequence>
  );
}
