import fs from 'fs'
import { Readable } from 'stream'
import { ReadableStream } from 'stream/web'
import { createWriteStream } from 'fs'
import path from 'path'
import os from 'os'
import YAML from 'yaml'
import packageJson from './package.json'

interface SupportedPlatform {
    TYPE: string
    ARCHITECTURE: string
    ARTIFACT_NAME: string
    BINARY_NAME: string
}

interface PackageJson {
    name: string
    repository: { url: string }
    version: string
}

const supportedPlatformsConfig: string = fs.readFileSync(
    path.join(__dirname, 'supported-platforms.yml'),
    'utf8'
)

const supportedPlatforms: SupportedPlatform[] = YAML.parse(supportedPlatformsConfig)
const { name, repository, version }: PackageJson = packageJson as PackageJson

const handleError = (msg: string) => {
    console.error(msg)
    process.exit(1)
}

const removeGitPrefixAndSuffix = (input: string): string => {
    let result = input.replace(/^git\+/, '')
    result = result.replace(/\.git$/, '')
    return result
}

const getPlatformMetadata = (): SupportedPlatform | undefined => {
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

const downloadBinaryFromGitHub = async (downloadUrl: string, outputPath: string): Promise<void> => {
    try {
        const binDir = path.resolve(process.cwd(), 'bin')
        if (!fs.existsSync(binDir)) {
            fs.mkdirSync(binDir, { recursive: true })
        }

        const resp = await fetch(downloadUrl)

        if (resp.ok && resp.body) {
            const writer = createWriteStream(outputPath)
            Readable.fromWeb(resp.body as ReadableStream).pipe(writer)
        } else {
            throw new Error(`Failed to download. Status: ${resp.status}`)
        }

        fs.chmodSync(outputPath, 0o755)
        console.log('Binary downloaded successfully!')
    } catch (err: any) {
        console.error('Error downloading binary:', err.message)
    }
}

const platformMetadata = getPlatformMetadata()
const repoUrl = removeGitPrefixAndSuffix(repository.url)
const binaryUrl = `${repoUrl}/releases/download/v${version}/${platformMetadata?.ARTIFACT_NAME}`
const binaryOutputPath = path.join(process.cwd(), 'bin', platformMetadata?.BINARY_NAME || '')

downloadBinaryFromGitHub(binaryUrl, binaryOutputPath)
