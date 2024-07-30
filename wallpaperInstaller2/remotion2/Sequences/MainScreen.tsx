import React, { useEffect, useMemo, useState } from "react";
import { interpolate, Sequence, Easing, useCurrentFrame } from "remotion";
import { MainContent } from "./MainContent/Content";
import { Screen1 } from "./MainContent/Screen1";
import { Screen2 } from "./MainContent/Screen2";
import logo from "../../public/applogo.svg";
import { Button, Checkbox, CheckboxProps } from "@fluentui/react-components";
import { PlayerRef } from "@remotion/player";
import {InstallerData} from '../Datas';
// Frame 216 - 600
export function MainScreen(props: {
  playerRef: React.RefObject<PlayerRef> | null;
  playbackRate: number;
  setPlaybackRate: (rate: number) => void;
  data: InstallerData;
  setData: (data: InstallerData) => void;
}) {
  const applogo = useMemo(() => logo, []);
  const [mainConfig, setMainConfig] = useState({
    nextscreen: "screen1",
  });
  const [screen1config, setScreen1Config] = useState({
    data1: "",
    data2: "",
  });
  const [screen2config, setScreen2Config] = useState({
    data1: "",
    data2: "",
  });

  const frame = useCurrentFrame();
  const transitionTargets = useMemo(() => {
    return {
      screen1to2: 275,
    };
  }, []);
  useEffect(() => {
    if (frame == 245) {
      props.playerRef?.current?.pause();
      props.setPlaybackRate(1);
    }
    if (frame == transitionTargets.screen1to2) {
      props.playerRef?.current?.pause();
    }
  }, [frame]);

  const smalleningcurve = useMemo(() => Easing.bezier(0, 1.21, 0.83, 0.92), []);
  const transitioncurve = useMemo(() => Easing.bezier(0, 1,1,1), []);

  return (
    <>
      <Sequence from={200} durationInFrames={80}>
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
                easing: smalleningcurve,
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
                easing: smalleningcurve,
                extrapolateLeft: "clamp",
                extrapolateRight: "clamp",
              })}%`,
              opacity: interpolate(frame, [200, 245], [0, 1], {
                easing: smalleningcurve,
                extrapolateLeft: "clamp",
                extrapolateRight: "clamp",
              }),
            }}
          ></div>

          {/* CONTENT */}
          <div
            className="pr-4 pl-4 mt-2 flex flex-row"
            style={{
              width: "100%",
              height: "100%",
              opacity: interpolate(frame, [200, 245], [0, 1], {
                easing: smalleningcurve,
                extrapolateLeft: "clamp",
                extrapolateRight: "clamp",
              }),
            }}
          >
            <div
              className="flex flex-col items-start pb-4"
              style={{
                position: "relative",
                opacity: interpolate(
                  frame,
                  [245, transitionTargets.screen1to2],
                  [1, 0],
                  {
                    easing: transitioncurve,
                    extrapolateLeft: "extend",
                    extrapolateRight: "extend",
                  }
                ),
                height: "100%",
                minWidth: `100%`,
                right: `${interpolate(
                  frame,
                  [245, transitionTargets.screen1to2],
                  [0, 100],
                  {
                    easing: transitioncurve,
                    extrapolateLeft: "extend",
                    extrapolateRight: "extend",
                  }
                )}%`,
              }}
            >
              <div
                style={{
                  height: `${interpolate(frame, [200, 245], [100, 0], {
                    easing: Easing.bezier(0, 1.21, 0.83, 0.92),
                    extrapolateLeft: "clamp",
                    extrapolateRight: "clamp",
                  })}%`,
                }}
              ></div>
              {useMemo(
                () => (
                  <MainContent
                    mainConfig={mainConfig}
                    setMainConfig={setMainConfig}
                    playerRef={props.playerRef}
                    data={props.data}
                    setData={props.setData}
                    setPlaybackRate={props.setPlaybackRate}
                  />
                ),
                [mainConfig, props.data]
              )}
            </div>

            <div
              className="flex flex-col items-start"
              style={{
                minWidth: `100%`,
                maxWidth: `100%`,
                overflow: "clip",
                opacity: interpolate(
                  frame,
                  [245, transitionTargets.screen1to2],
                  [0, 1],
                  {
                    easing: transitioncurve,
                    extrapolateLeft: "clamp",
                    extrapolateRight: "clamp",
                  }
                ),
                position: "relative",
                right: `${interpolate(
                  frame,
                  [245, transitionTargets.screen1to2],
                  [0, 100],
                  {
                    easing: transitioncurve,
                    extrapolateLeft: "extend",
                    extrapolateRight: "extend",
                  }
                )}%`,
              }}
            >
              {useMemo(() => {
                if (mainConfig.nextscreen == "screen1") {
                  return (
                    <Screen1
                      playerRef={props.playerRef}
                      config={screen1config}
                      setstate={setScreen1Config}
                      playbackRate={props.playbackRate}
                      setPlaybackRate={props.setPlaybackRate}
                      data={props.data}
                      setData={props.setData}
                    />
                  );
                } else {
                  return <></>;
                }
              }, [screen1config, mainConfig.nextscreen])}
              {useMemo(() => {
                if (mainConfig.nextscreen == "screen2") {
                  return (
                    <Screen2
                      playerRef={props.playerRef}
                      config={screen2config}
                      setstate={setScreen2Config}
                      setPlaybackRate={props.setPlaybackRate}
                      data={props.data}
                      setData={props.setData}
                    />
                  );
                } else {
                  return <></>;
                }
              }, [screen2config, mainConfig.nextscreen])}
            </div>
          </div>
        </div>
      </Sequence>
    </>
  );
}
