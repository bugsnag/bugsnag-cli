const esbuild = require("esbuild");
const path = require("node:path");

async function run() {
    try {
        await esbuild.build({
            entryPoints: {
                index: path.join(__dirname, "src/index.js"),
                utils: path.join(__dirname, "src/utils/index.js"),
            },
            outdir: path.join(__dirname, "dist"),
            bundle: true,
            sourcemap: true,
            platform: "node",
            format: "cjs",
            target: "node18",
        });

        console.log("Build complete. Sourcemap generated.");
    } catch (err) {
        console.error(err);
        process.exit(1);
    }
}

run();
