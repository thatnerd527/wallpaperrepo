import React from "react";
import { PlayerRef } from "@remotion/player";
import { Button } from "@fluentui/react-components";
import { InstallerData } from "../../Datas";
export function Screen2(props: {
  config: { data1: string; data2: string };
  setstate: (data: { data1: string; data2: string }) => void;
  playerRef: React.RefObject<PlayerRef> | null;
  setPlaybackRate: (rate: number) => void;
  data: InstallerData;
  setData: (data: InstallerData) => void;
}) {
  return (
    <div>
      screen2
          <Button
        onClick={() => {
          console.log("clicked");
          props.setPlaybackRate(-1);
          props.data.safeMerge(
            {
              ...props.data,
            },
            props.setData
          );
          props.playerRef?.current?.play();
        }}
      >
        Back
      </Button>
    </div>
  );
}
