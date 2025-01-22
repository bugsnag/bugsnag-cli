const fs = require('fs')
const path = require('path')
const os = require('os')
const YAML = require('yaml')
const createWriteStream = require('fs').createWriteStream
const Readable = require('stream').Readable

const supportedPlatformsConfig = fs.readFileSync(path.join(__dirname, 'supported-platforms.yml'), 'utf8')
const { name, repository, version } = require('./package.json')

const handleError = (msg) => {
    console.error(msg)
    process.exit(1)
}

// Remove the git prefix and suffix from the repository URL
const removeGitPrefixAndSuffix = (input) => {
    let result = input.replace(/^git\+/, '')
    result = result.replace(/\.git$/, '')
    return result
}

// Parse supported-platforms.yml into an iterable array
const supportedPlatforms = YAML.parse(supportedPlatformsConfig)

const getPlatformMetadata = () => {
    const type = os.type()
    const architecture = os.arch()

    for (const supportedPlatform of supportedPlatforms) {
        if (type === supportedPlatform.TYPE && architecture === supportedPlatform.ARCHITECTURE) {
            return supportedPlatform
        }
    }

    const supportedPlatformsTable = supportedPlatforms.map((platform) => {
        return {
            Type: platform.TYPE,
            Architecture: platform.ARCHITECTURE,
            Artifact: platform.ARTIFACT_NAME
        }
    })

    handleError(
        `Platform with type "${type}" and architecture "${architecture}" is not supported by ${name}.\nYour system must be one of the following:\n\n${JSON.stringify(
            supportedPlatformsTable,
            null,
            2
        )}`
    )
}

const downloadBinaryFromGitHub = async (downloadUrl, outputPath) => {
    try {
        const binDir = path.resolve(process.cwd(), 'bin')
        if (!fs.existsSync(binDir)) {
            fs.mkdirSync(binDir, { recursive: true })
        }

        const fileName = downloadUrl.split("/").pop()
        const resp = await fetch(downloadUrl)

        if (resp.ok && resp.body) {
            let writer = createWriteStream(outputPath)
            Readable.fromWeb(resp.body).pipe(writer)
        }

        fs.chmodSync(outputPath, '755')
        console.log('Binary downloaded successfully!')
    } catch (err) {
        console.error('Error downloading binary:', err.message)
    }
}

const platformMetadata = getPlatformMetadata()
const repoUrl = removeGitPrefixAndSuffix(repository.url)
const binaryUrl = `${repoUrl}/releases/download/v${version}/${platformMetadata.ARTIFACT_NAME}`
const binaryOutputPath = path.join(process.cwd(), 'bin', platformMetadata.BINARY_NAME)

downloadBinaryFromGitHub(binaryUrl, binaryOutputPath)
