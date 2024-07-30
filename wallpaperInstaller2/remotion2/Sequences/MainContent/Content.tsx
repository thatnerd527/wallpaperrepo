import React from "react";
import { Button, Checkbox } from "@fluentui/react-components";
import { PlayerRef } from "@remotion/player";
import { InstallerData } from "../../Datas";

export function MainContent(props: {
  mainConfig: { nextscreen: string };
  setMainConfig: (data: { nextscreen: string }) => void;
  playerRef: React.RefObject<PlayerRef> | null;
  data: InstallerData;
  setPlaybackRate: (rate: number) => void;
  setData: (data: InstallerData) => void;
}) {
  //console.log("ReRendr");
  return (
    <>
      <div className="w-full h-full flex flex-col">
        <div className="w-full flex flex-row justify-start items-center">
          <Checkbox
            checked={props.data.runatstartup}
            onChange={() =>
              props.data.safeMerge(
                {
                  ...props.data,
                  runatstartup: !props.data.runatstartup,
                },
                props.setData
              )
            }
          ></Checkbox>
          Run at login
        </div>
        <div className="w-full flex flex-row justify-start items-center">
          <Checkbox
            checked={props.data.createdesktopshortcuts}
            onChange={() =>
              props.data.safeMerge(
                {
                  ...props.data,
                  createdesktopshortcuts: !props.data.createdesktopshortcuts,
                },
                props.setData
              )
            }
          ></Checkbox>
          Create desktop shortcuts
        </div>
        <div className="w-full flex flex-row justify-start items-center">
          <Checkbox
            checked={props.data.enableautoupdates}
            onChange={() =>
              props.data.safeMerge(
                {
                  ...props.data,
                  enableautoupdates: !props.data.enableautoupdates,
                },
                props.setData
              )
            }
          ></Checkbox>
          Enable automatic updates
        </div>
        <div className="w-full flex flex-row justify-start items-center">
          <Checkbox
            checked={props.data.sendlogs}
            onChange={() =>
              props.data.safeMerge(
                {
                  ...props.data,
                  sendlogs: !props.data.sendlogs,
                },
                props.setData
              )
            }
          ></Checkbox>
          Enable sending anonymous usage data
        </div>
        <div className="h-full w-full flex flex-row justify-end items-end">
          <Button
            onClick={() => {
              props.setMainConfig({ nextscreen: "screen2" });

              props.playerRef?.current?.play();
            }}
          >
            Advanced configuration
          </Button>
          <div className="w-4"></div>
          <Button
            appearance="primary"
            onClick={() => {
              props.data.safeMerge(
                { ...props.data, pressedbutton: true },
                props.setData
                );
                props.setPlaybackRate(-1);
              props.playerRef?.current?.play();
            }}
          >
            Install
          </Button>
        </div>
      </div>
    </>
  );
}
