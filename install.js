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
        TYPE: 'Windows',
        ARCHITECTURE: 'x64',
        ARTIFACT_NAME: 'x86_64-windows-bugsnag-cli.exe',
        BINARY_NAME: 'bugsnag-cli.exe'
    },
    {
        TYPE: 'Windows',
        ARCHITECTURE: 'i386',
        ARTIFACT_NAME: 'i386-windows-bugsnag-cli.exe',
        BINARY_NAME: 'bugsnag-cli.exe'
    },
    {
        TYPE: 'Linux',
        ARCHITECTURE: 'x64',
        ARTIFACT_NAME: 'x86_64-linux-bugsnag-cli',
        BINARY_NAME: 'bugsnag-cli'
    },
    {
        TYPE: 'Linux',
        ARCHITECTURE: 'i386',
        ARTIFACT_NAME: 'i386-linux-bugsnag-cli',
        BINARY_NAME: 'bugsnag-cli'
    },
    {
        TYPE: 'Darwin',
        ARCHITECTURE: 'x64',
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
        const binDir = path.resolve(process.cwd(),'..','..','.bin');
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
    fs.readFile(packageJsonPath, 'utf8', (err, data) => {
        if (err) {
            console.error(`Error reading package.json: ${err}`);
            return;
        }

        try {
            const packageJson = JSON.parse(data);

            packageJson.scripts = {
                ...packageJson.scripts,
                "bugsnag:create-build": "./node_modules/.bin/bugsnag-cli create-build",
                "bugsnag:upload-android": "./node_modules/.bin/bugsnag-cli upload react-native-android"
            };

            const updatedPackageJson = JSON.stringify(packageJson, null, 2);

            fs.writeFile(packageJsonPath, updatedPackageJson, 'utf8', (err) => {
                if (err) {
                    console.error(`Error writing package.json: ${err}`);
                    return;
                }
            });
        } catch (err) {
            console.error(`Error parsing package.json: ${err}`);
        }
    })
}

const platformMetadata = getPlatformMetadata();
const repoUrl = removeGitPrefixAndSuffix(repository.url);
const binaryUrl = `${repoUrl}/releases/download/v${version}/${platformMetadata.ARTIFACT_NAME}`;
const binaryOutputPath = path.join(process.cwd(),'..','..','.bin', platformMetadata.BINARY_NAME);
const projectPackageJsonPath = path.join(process.cwd(),'..','..', '..','package.json');

downloadBinaryFromGitHub(binaryUrl, binaryOutputPath);
writeToPackageJson(projectPackageJsonPath)