import {Easing, interpolate, Sequence, useCurrentFrame} from 'remotion';
import React from 'react';
import logo from '../../public/applogo.svg';
// Frame 0 - 216
export function Opening() {
    const frame = useCurrentFrame();
    const applogo = logo;

  const openingcurve = Easing.bezier(0.02, 1.66, 0.47, 0.33);
  const smalleningcurve = Easing.bezier(0, 1.21, 0.83, 0.92);
    return <>
    <Sequence from={0} durationInFrames={96}>
        <div className="w-full h-full flex flex-col justify-center items-center">
          <img
            src={applogo}
            style={{
              width: interpolate(frame, [0, 96], [0, 150], {
                easing: openingcurve,
                extrapolateLeft: "clamp",
                extrapolateRight: "clamp",
              }),
              height: interpolate(frame, [0, 96], [0, 150], {
                easing: openingcurve,
                extrapolateLeft: "clamp",
                extrapolateRight: "clamp",
              }),
              opacity: interpolate(frame, [0, 30], [0, 1], {
                easing: Easing.linear,
              }),
            }}
          />
        </div>
      </Sequence>
      <Sequence from={96} durationInFrames={60}>
        <div className="w-full h-full flex flex-col justify-center items-center">
          <img
            src={applogo}
            style={{
              width: interpolate(frame, [96, 156], [150, 160], {
                extrapolateLeft: "clamp",
                extrapolateRight: "clamp",
              }),
              height: interpolate(frame, [96, 156], [150, 160], {
                extrapolateLeft: "clamp",
                extrapolateRight: "clamp",
              }),
            }}
          />
        </div>
      </Sequence>
      <Sequence from={156} durationInFrames={60}>
        <div className="w-full h-full flex flex-col justify-start items-center">
          <img
            src={applogo}
            style={{
              width: interpolate(frame, [156, 216], [160, 64], {
                easing: smalleningcurve,
              }),
              height: interpolate(frame, [156, 216], [160, 64], {
                easing: smalleningcurve,
              }),
              position: "relative",
              marginTop: `${interpolate(frame, [156, 216], [25, 0], {
                easing: smalleningcurve,
              })}%`,
            }}
          />
        </div>
      </Sequence>
    </>;
}