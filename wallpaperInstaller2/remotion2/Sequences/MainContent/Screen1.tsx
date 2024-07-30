import React from 'react';
import {PlayerRef} from '@remotion/player';
import { Button } from '@fluentui/react-components';
export function Screen1(props: {
  config: { data1: string; data2: string };
  setstate: (data: { data1: string; data2: string }) => void;
  playerRef: React.RefObject<PlayerRef> | null;
  playbackRate: number;
  setPlaybackRate: (rate: number) => void;
}) {
  return (
    <div>
      screen1
      <Button
        onClick={() => {
                  console.log("clicked");
                  props.setPlaybackRate(-1)
          props.playerRef?.current?.play();
        }}
      >
        Back
      </Button>
    </div>
  );
}