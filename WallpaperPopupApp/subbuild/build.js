
import { readFileSync, writeFileSync } from "fs";
import { join } from "path";

let file = readFileSync(
  join(".", "..", "src", "protocol", "protocol.js"),
  "utf8"
);
file = file.replace(
  `import * as $protobuf from "protobufjs/minimal";`,
    `
  import * as $protobuf1 from "protobufjs/minimal.js";
  const $protobuf = $protobuf1.default;
`
);
console.log("Writing to protocol.js");
writeFileSync(join(".", "..", "src", "protocol", "protocol.js"), file);