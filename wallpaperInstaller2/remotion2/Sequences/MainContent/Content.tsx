import React from "react";
import { Button, Checkbox } from "@fluentui/react-components";
import { PlayerRef } from "@remotion/player";

export function MainContent(props: {
  config: { option1: boolean; option2: boolean; option3: boolean };
  mainConfig: { nextscreen: string };
  setstate: (data: {
    option1: boolean;
    option2: boolean;
    option3: boolean;
  }) => void;
  setMainConfig: (data: { nextscreen: string }) => void;
  playerRef: React.RefObject<PlayerRef> | null;
}) {
  console.log("ReRendr");
  return (
    <>
      <div className="flex flex-row items-center mb-2">
        <Checkbox
          checked={props.config.option1}
          onChange={(_, data) => {
            data &&
              props.setstate({
                ...props.config,
                option1: data.checked as boolean,
              });
          }}
        ></Checkbox>
        <div className="w-1.5"></div>
        Option 1
      </div>
      <div className="flex flex-row items-center mb-2">
        <Checkbox
          checked={props.config.option2}
          onChange={(_, data) => {
            data &&
              props.setstate({
                ...props.config,
                option2: data.checked as boolean,
              });
          }}
        ></Checkbox>
        <div className="w-1.5"></div>
        Option 2
      </div>
      <div className="flex flex-row items-center mb-2">
        <Checkbox
          checked={props.config.option3}
          onChange={(_, data) => {
            data &&
              props.setstate({
                ...props.config,
                option3: data.checked as boolean,
                option2: data.checked as boolean,
              });
          }}
        ></Checkbox>
        <div className="w-1.5"></div>
        Option 3
      </div>
      <Button
        onClick={() => {
                  props.setMainConfig({ nextscreen: "screen1" });
                  props.playerRef?.current?.play();
        }}
      >
        Set next screen to screen 1
      </Button>
      <Button
        onClick={() => {
                  props.setMainConfig({ nextscreen: "screen2" });

                  props.playerRef?.current?.play();
        }}
      >
        Set next screen to screen 2
      </Button>
    </>
  );
}
