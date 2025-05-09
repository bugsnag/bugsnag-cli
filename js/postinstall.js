const fs = require('fs')
const { Readable } = require('stream')
const { createWriteStream } = require('fs')
const path = require('path')
const os = require('os')
const YAML = require('yaml')
const packageJson = require('./package.json')

const supportedPlatformsConfig = fs.readFileSync(
    path.join(__dirname, 'supported-platforms.yml'),
    'utf8'
)

const supportedPlatforms = YAML.parse(supportedPlatformsConfig)
const { name, repository, version } = packageJson

const handleError = (msg) => {
    console.error(msg)
    process.exit(1)
}

const removeGitPrefixAndSuffix = (input) => {
    let result = input.replace(/^git\+/, '')
    result = result.replace(/\.git$/, '')
    return result
}

const getPlatformMetadata = () => {
    const type = os.type()
    const architecture = os.arch()

    for (const supportedPlatform of supportedPlatforms) {
        if (type === supportedPlatform.TYPE && architecture === supportedPlatform.ARCHITECTURE) {
            return supportedPlatform
        }
    }

    const supportedPlatformsTable = supportedPlatforms.map((platform) => ({
        Type: platform.TYPE,
        Architecture: platform.ARCHITECTURE,
        Artifact: platform.ARTIFACT_NAME
    }))

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

        const resp = await fetch(downloadUrl)

        if (resp.ok && resp.body) {
            const writer = createWriteStream(outputPath)
            Readable.fromWeb(resp.body).pipe(writer)
        } else {
            throw new Error(`Failed to download. Status: ${resp.status}`)
        }

        fs.chmodSync(outputPath, 0o755)
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
