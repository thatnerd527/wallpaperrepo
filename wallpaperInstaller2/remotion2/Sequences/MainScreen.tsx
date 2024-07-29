import React, { useState } from "react";
import { interpolate, Sequence, Easing, useCurrentFrame } from "remotion";
import logo from "../../public/applogo.svg";
import { Button, Checkbox, CheckboxProps } from "@fluentui/react-components";
// Frame 216 - 600
export function MainScreen() {
  const applogo = logo;
  const [checked, setChecked] = React.useState<CheckboxProps["checked"]>(true);
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
              position: "relative",
              marginTop: `${interpolate(frame, [156, 216], [25, 0], {
                extrapolateLeft: "clamp",
                extrapolateRight: "clamp",
                easing: smalleningcurve,
              })}%`,
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
          <Checkbox
            value={"Test"}
            checked={checked}
            onChange={(ev, data) => {
              setChecked(data?.checked);
            }}
          ></Checkbox>
        </div>
      </Sequence>
    </>
  );
}
