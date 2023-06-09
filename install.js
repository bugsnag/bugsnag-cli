const axios = require('axios');
const fs = require('fs');
const path = require('path');
const os = require('os');

const { name, repository, version } = require('./package.json');

const handleError = (msg) => {
    console.error(msg);
    process.exit(1);
};

// Remove the git prefix and suffix from the repository URL
const removeGitPrefixAndSuffix = (input) => {
    let result = input.replace(/^git\+/, '');
    result = result.replace(/\.git$/, '');
    return result;
};

const supportedPlatforms = [
    {
        TYPE: 'windows',
        ARCHITECTURE: 'x86_64',
        ARTIFACT_NAME: 'x86_64-windows-bugsnag-cli.exe',
        BINARY_NAME: 'bugsnag-cli.exe'
    },
    {
        TYPE: 'windows',
        ARCHITECTURE: 'i386',
        ARTIFACT_NAME: 'i386-windows-bugsnag-cli.exe',
        BINARY_NAME: 'bugsnag-cli.exe'
    },
    {
        TYPE: 'linux',
        ARCHITECTURE: 'x86_64',
        ARTIFACT_NAME: 'x86_64-linux-bugsnag-cli',
        BINARY_NAME: 'bugsnag-cli'
    },
    {
        TYPE: 'linux',
        ARCHITECTURE: 'i386',
        ARTIFACT_NAME: 'i386-linux-bugsnag-cli',
        BINARY_NAME: 'bugsnag-cli'
    },
    {
        TYPE: 'Darwin',
        ARCHITECTURE: 'x86_64',
        ARTIFACT_NAME: 'x86_64-macos-bugsnag-cli',
        BINARY_NAME: 'bugsnag-cli'
    },
    {
        TYPE: 'Darwin',
        ARCHITECTURE: 'arm64',
        ARTIFACT_NAME: 'arm64-macos-bugsnag-cli',
        BINARY_NAME: 'bugsnag-cli'
    }
];

const getPlatformMetadata = () => {
    const type = os.type();
    const architecture = os.arch();

    for (const supportedPlatform of supportedPlatforms) {
        if (type === supportedPlatform.TYPE && architecture === supportedPlatform.ARCHITECTURE) {
            return supportedPlatform;
        }
    }

    const supportedPlatformsTable = supportedPlatforms.map((platform) => {
        return {
            Type: platform.TYPE,
            Architecture: platform.ARCHITECTURE,
            Artifact: platform.ARTIFACT_NAME
        };
    });

    handleError(
        `Platform with type "${type}" and architecture "${architecture}" is not supported by ${name}.\nYour system must be one of the following:\n\n${JSON.stringify(
            supportedPlatformsTable,
            null,
            2
        )}`
    );
};

const downloadBinaryFromGitHub = async (downloadUrl, outputPath) => {
    try {
        const binDir = path.resolve(__dirname, '..', '.bin');
        if (!fs.existsSync(binDir)) {
            fs.mkdirSync(binDir, { recursive: true });
        }

        const response = await axios.get(downloadUrl, { responseType: 'arraybuffer' });
        const binaryData = response.data;
        fs.writeFileSync(outputPath, binaryData, 'binary');
        fs.chmodSync(outputPath, '755');
        console.log('Binary downloaded successfully!');
    } catch (err) {
        console.error('Error downloading binary:', err.message);
    }
};

const writeToPackageJson = (packageJsonPath) => {
    const packageJson = require(packageJsonPath);

    packageJson.scripts.bugsnagCreateBuild = './node_modules/.bin/bugsnag-cli create-build';
    packageJson.scripts.bugsnagUpload = './node_modules/.bin/bugsnag-cli upload react-native-android';

    fs.writeFileSync(packageJsonPath, JSON.stringify(packageJson, null, 2));
}

const platformMetadata = getPlatformMetadata();
const repoUrl = removeGitPrefixAndSuffix(repository.url);
const binaryUrl = `${repoUrl}/releases/download/v${version}/${platformMetadata.ARTIFACT_NAME}`;
const binaryOutputPath = path.join(__dirname, '..', '.bin', platformMetadata.BINARY_NAME);
const projectPackageJsonPath = path.join(__dirname, '..', '..', 'package.json');

downloadBinaryFromGitHub(binaryUrl, binaryOutputPath);
writeToPackageJson(projectPackageJsonPath)
