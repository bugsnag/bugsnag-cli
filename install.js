const axios = require('axios');
const fs = require('fs');
const path = require('path');
const os = require("os");
const cTable = require("console.table");


const {name, repository, version} = require("./package.json");

const error = msg => {
    console.error(msg);
    process.exit(1);
};

const supportedPlatforms = [
    {
        TYPE: "windows",
        ARCHITECTURE: "x86_64",
        ARTIFACT_NAME: "x86_64-windows-bugsnag-cli.exe",
        BINARY_NAME: "bugsnag-cli.exe"
    },
    {
        TYPE: "windows",
        ARCHITECTURE: "i386",
        ARTIFACT_NAME: "i386-windows-bugsnag-cli.exe",
        BINARY_NAME: "bugsnag-cli.exe"
    },
    {
        TYPE: "linux",
        ARCHITECTURE: "x86_64",
        ARTIFACT_NAME: "x86_64-linux-bugsnag-cli",
        BINARY_NAME: "bugsnag-cli"

    },
    {
        TYPE: "linux",
        ARCHITECTURE: "i386",
        ARTIFACT_NAME: "i386-linux-bugsnag-cli",
        BINARY_NAME: "bugsnag-cli"
    },
    {
        TYPE: "Darwin",
        ARCHITECTURE: "x86_64",
        ARTIFACT_NAME: "x86_64-macos-bugsnag-cli",
        BINARY_NAME: "bugsnag-cli"
    },
    {
        TYPE: "Darwin",
        ARCHITECTURE: "arm64",
        ARTIFACT_NAME: "arm64-macos-bugsnag-cli",
        BINARY_NAME: "bugsnag-cli"
    }
];

const getPlatformMetadata = () => {
    const type = os.type();
    const architecture = os.arch();

    for (let supportedPlatform of supportedPlatforms) {
        if (
            type === supportedPlatform.TYPE &&
            architecture === supportedPlatform.ARCHITECTURE
        ) {
            return supportedPlatform;
        }
    }

    error(
        `Platform with type "${type}" and architecture "${architecture}" is not supported by ${name}.\nYour system must be one of the following:\n\n${cTable.getTable(
            supportedPlatforms
        )}`
    );
};

async function downloadBinaryFromGitHub(url, outputPath) {
    try {
        const response = await axios.get(url, { responseType: 'arraybuffer' });
        const binaryData = response.data;
        fs.writeFileSync(outputPath, binaryData, 'binary');
        fs.chmodSync(outputPath, '755');
        console.log('Binary downloaded successfully!');
    } catch (error) {
        console.error('Error downloading binary:', error.message);
    }
}

const platformMetadata = getPlatformMetadata();

const binaryUrl = `${repository.url}/releases/download/${version}/${platformMetadata.ARTIFACT_NAME}`;

const binaryOutputPath = path.join(__dirname, 'node_modules', '.bin', platformMetadata.BINARY_NAME);

downloadBinaryFromGitHub(binaryUrl, binaryOutputPath);
