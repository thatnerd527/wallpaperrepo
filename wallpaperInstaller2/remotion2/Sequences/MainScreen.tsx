import React, { useEffect, useMemo, useState } from "react";
import { interpolate, Sequence, Easing, useCurrentFrame } from "remotion";
import { MainContent } from "./MainContent/Content";
import { Screen1 } from "./MainContent/Screen1";
import { Screen2 } from "./MainContent/Screen2";
import logo from "../../public/applogo.svg";
import { Button, Checkbox, CheckboxProps } from "@fluentui/react-components";
import { PlayerRef } from "@remotion/player";
// Frame 216 - 600
export function MainScreen(props: {
  playerRef: React.RefObject<PlayerRef> | null;
  playbackRate: number;
  setPlaybackRate: (rate: number) => void;
}) {
  const applogo = useMemo(() => logo, []);
  const [mainConfig, setMainConfig] = useState({
    nextscreen: "screen1",
  });
  const [page1config, setPage1Config] = useState({
    option1: false,
    option2: false,
    option3: false,
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
      screen1to2: 290,
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
          <div
            className="w-full h-full pr-4 pl-4 mt-2 flex flex-row"
            style={{
              opacity: interpolate(frame, [200, 245], [0, 1], {
                easing: smalleningcurve,
                extrapolateLeft: "clamp",
                extrapolateRight: "clamp",
              }),
            }}
          >
            <div
              className="flex flex-col items-start"
              style={{
                minWidth: `${interpolate(
                  frame,
                  [245, transitionTargets.screen1to2],
                  [100, 0],
                  {
                    easing: smalleningcurve,
                    extrapolateLeft: "clamp",
                    extrapolateRight: "clamp",
                  }
                )}%`,
                maxWidth: `${interpolate(
                  frame,
                  [245, transitionTargets.screen1to2],
                  [100, 0],
                  {
                    easing: smalleningcurve,
                    extrapolateLeft: "clamp",
                    extrapolateRight: "clamp",
                  }
                )}%`,
                overflow: "clip",
                opacity: interpolate(
                  frame,
                  [245, transitionTargets.screen1to2],
                  [1, 0],
                  {
                    easing: smalleningcurve,
                    extrapolateLeft: "clamp",
                    extrapolateRight: "clamp",
                  }
                ),
              }}
            >
              <div
                style={{
                  height: `${interpolate(frame, [200, 245], [50, 0], {
                    easing: Easing.bezier(0, 1.21, 0.83, 0.92),
                    extrapolateLeft: "clamp",
                    extrapolateRight: "clamp",
                  })}%`,
                }}
              ></div>
              {useMemo(
                () => (
                  <MainContent
                    config={page1config}
                    mainConfig={mainConfig}
                    setMainConfig={setMainConfig}
                    setstate={setPage1Config}
                    playerRef={props.playerRef}
                  />
                ),
                [page1config, mainConfig]
              )}
            </div>

            <div
              className="flex flex-col items-start"
              style={{
                minWidth: `${interpolate(
                  frame,
                  [245, transitionTargets.screen1to2],
                  [0, 100],
                  {
                    easing: smalleningcurve,
                    extrapolateLeft: "clamp",
                    extrapolateRight: "clamp",
                  }
                )}%`,
                maxWidth: `${interpolate(
                  frame,
                  [245, transitionTargets.screen1to2],
                  [0, 100],
                  {
                    easing: smalleningcurve,
                    extrapolateLeft: "clamp",
                    extrapolateRight: "clamp",
                  }
                )}%`,
                overflow: "clip",
                opacity: interpolate(
                  frame,
                  [245, transitionTargets.screen1to2],
                  [0, 1],
                  {
                    easing: smalleningcurve,
                    extrapolateLeft: "clamp",
                    extrapolateRight: "clamp",
                  }
                ),
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
