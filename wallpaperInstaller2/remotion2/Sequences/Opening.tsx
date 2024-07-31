import { Easing, interpolate, Sequence, useCurrentFrame } from "remotion";
import React, { useMemo } from "react";
import logo from "../../public/applogo.svg";
import { InstallerData } from "remotion2/Datas";
import { Lottie } from "@remotion/lottie";
import animationData from "../../public/animation.json";
// Frame 0 - 216
export function Opening(props: {
  data: InstallerData;
  setData: (data: InstallerData) => void;
}) {
  const frame = useCurrentFrame();
  const applogo = useMemo(() => logo, []);

  const openingcurve = useMemo(() => Easing.bezier(0.02, 1.66, 0.47, 0.33), []);
  const smalleningcurve = useMemo(() => Easing.bezier(0, 1.21, 0.83, 0.92), []);
  const anim = useMemo(() => animationData, []);
  return (
    <>
      <Sequence from={0} durationInFrames={156}>
        <div className="w-full h-full flex flex-col justify-center items-center">
          <Lottie
            animationData={anim}
            loop={true}
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
        </div>
      </Sequence>
      <Sequence from={156} durationInFrames={60}>
        <div className="w-full h-full flex flex-col items-center justify-center relative">
          <div
            className="w-full h-full flex flex-col justify-start items-center"
            style={{
              height: `calc(max(${interpolate(frame, [156, 216], [0, 100], {
                easing: smalleningcurve,
              })}%, 145px))`,
            }}
          >
            <img
              src={applogo}
              style={{
                width: interpolate(frame, [156, 216], [145, 64], {
                  easing: smalleningcurve,
                }),
                height: interpolate(frame, [156, 216], [145, 64], {
                  easing: smalleningcurve,
                }),
                position: "relative",
              }}
            />
          </div>
        </div>
      </Sequence>
    </>
  );
}
