import React from "react";
import { interpolate, Sequence, Easing, useCurrentFrame } from "remotion";
import logo from "../../public/applogo.svg";
// Frame 216 - 600
export function MainScreen() {
  const applogo = logo;

  const frame = useCurrentFrame();
  const smalleningcurve = Easing.bezier(0, 1.21, 0.83, 0.92);
  return (
    <>
      <Sequence from={200} durationInFrames={600}>
        <div className="w-full h-full flex flex-col justify-start items-center">
          <img
            src={applogo}
            style={{
              width: interpolate(frame, [156, 216], [160, 64], {
                extrapolateLeft: "clamp",
                extrapolateRight: "clamp",
                easing: smalleningcurve,
              }),
              height: interpolate(frame, [156, 216], [160, 64], {
                extrapolateLeft: "clamp",
                extrapolateRight: "clamp",
                easing: smalleningcurve,
              }),
              display: frame >= 216 ? "none" : "block",
              position: "relative",
              opacity: 0,
              marginBottom: "8px",
            }}
          />
          <img
            src={applogo}
            style={{
              width: "64px",
              height: "64px",
              position: "relative",
              marginTop: "2px",
              marginBottom: "6px",
              display: frame >= 216 ? "block" : "none",
            }}
          />
          <h1
            className="font-semibold text-3xl mt-3 mb-3"
            style={{
              opacity: interpolate(frame, [200, 245], [0, 1], {
                easing: Easing.bezier(0, 1.21, 0.83, 0.92),
                extrapolateLeft: "clamp",
                extrapolateRight: "clamp",
              }),
            }}
          >
            WallpaperUI
          </h1>
          <div
            className="h-[2px] bg-gray-500 mt-2"
            style={{
              width: `${interpolate(frame, [200, 245], [0, 95], {
                easing: Easing.bezier(0, 1.21, 0.83, 0.92),
                extrapolateLeft: "clamp",
                extrapolateRight: "clamp",
              })}%`,
              opacity: interpolate(frame, [200, 245], [0, 1], {
                easing: Easing.bezier(0, 1.21, 0.83, 0.92),
                extrapolateLeft: "clamp",
                extrapolateRight: "clamp",
              }),
            }}
          ></div>
        </div>
      </Sequence>
    </>
  );
}
