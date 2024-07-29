import React from "react";
import { Composition } from "remotion";
import { InstallerFrame } from "./Composition";
import {
  FluentProvider,
  webLightTheme,
  Button,
} from "@fluentui/react-components";

export const RemotionRoot: React.FC = () => {
  return (
    <>
      <FluentProvider theme={webLightTheme}>
        <Composition
          id="Empty"
          component={InstallerFrame}
          defaultProps={{
            playerRef: null,
          }}
          durationInFrames={900}
          fps={60}
          width={300}
          height={450}
        />
      </FluentProvider>
    </>
  );
};
