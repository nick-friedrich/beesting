import concurrently from "concurrently";
import path from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const { result } = concurrently(
  [
    "npm:watch-*",
    { command: "npm run watch-css", name: "Tailwind" },
    {
      command: "cd ../.. && go run ./cmd/beesting dev example-app",
      name: "Beesting",
    },
    // { command: "deploy", name: "deploy", env: { PUBLIC_KEY: "..." } },
    // {
    //   command: "watch",
    //   name: "watch",
    //   cwd: path.resolve(__dirname, "scripts/watchers"),
    // },
  ],
  {
    prefix: "",
    killOthers: ["failure", "success"],
    restartTries: 3,
    cwd: path.resolve(__dirname, ".."),
  },
);

function success() {
  console.log("All processes completed successfully");
}

function failure() {
  console.log("One or more processes failed");
}

result.then(success, failure);
