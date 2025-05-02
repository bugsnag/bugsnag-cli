import { exec } from 'child_process'

interface BaseOptions {
    apiKey?: string
    dryRun?: boolean
    logLevel?: string
    port?: number
    failOnUploadError?: boolean
    verbose?: boolean
    overwrite?: boolean
    retries?: number
    timeout?: number
}

export interface BugsnagCreateBuildOptions extends BaseOptions {
    autoAssignRelease?: boolean
    buildApiRootUrl?: string
    builderName?: string
    metadata?: object
    provider?: string
    releaseStage?: string
    repository?: string
    revision?: string
    versionName?: string
    androidAab?: string
    appManifest?: string
    versionCode?: string
    bundleVersion?: string
}

interface UploadOptions extends BaseOptions {
    uploadApiRootUrl?: string
    projectRoot?: string
    dev?: boolean
    bundle?: string
    versionName?: string
    sourceMap?: string
    codeBundleId?: string
}

export interface BugsnagUploadReactNativeOptions extends UploadOptions {
    androidAppManifest?: string
    androidVariant?: string
    androidVersionCode?: string
    iosBundleVersion?: string
    iosPlist?: string
    iosScheme?: string
    iosXcodeProject?: string
}

export interface BugsnagUploadiOSOptions extends UploadOptions {
    sourceMap?: string
    bundleVersion?: string
    plist?: string
    scheme?: string
    xcodeProject?: string
}

export interface BugsnagUploadAndroidOptions extends UploadOptions {
    appManifest?: string
    variant?: string
    versionCode?: string
}

export interface BugsnagUploadJsOptions extends UploadOptions {
    baseUrl?: string
    bundleUrl?: string
    projectRoot?: string
}

/**
 * Wrapper for Bugsnag CLI
 */
class BugsnagCLI {
    /**
     * Convert camelCase to kebab-case
     */
    static camelToKebab(str: string) {
        return str.replace(/([a-z0-9])([A-Z])/g, '$1-$2').toLowerCase()
    }

    /**
     * Execute a Bugsnag CLI command
     */
    static run(command: string, options = {}, target = ''): Promise<string> {
        return new Promise((resolve, reject) => {
            // Convert the options keys from camelCase to kebab-case
            const kebabCaseOptions = Object.entries(options)
                .map(([key, value]) => {
                    const kebabKey = BugsnagCLI.camelToKebab(key)
                    if (typeof value === 'boolean' && value === true) {
                        return `--${kebabKey}`
                    } else if (typeof value !== 'boolean') {
                        return `--${kebabKey}=${value}`
                    }
                    return ''
                })
                .filter(Boolean)
                .join(' ')

            const positionalArg = target ? `"${target}"` : ''
            const cliCommand = `npx bugsnag-cli ${command} ${kebabCaseOptions} ${positionalArg}`.trim()

            // Execute the command
            exec(cliCommand, (error, stdout, stderr) => {
                if (error) {
                    const errorMessage = `Command failed: ${cliCommand}\n` +
                        `Error: ${error.message}\n` +
                        `${stdout.trim()}`
                    reject(errorMessage)
                } else {
                    resolve(stdout.trim())
                }
            })
        })
    }

    /**
     * Upload sourcemaps to Bugsnag
     * Provides nested methods for specific upload types.
     */
    static Upload = {

        ReactNative: Object.assign(
            (options: BugsnagUploadReactNativeOptions = {}, target = ''): Promise<string> =>
                BugsnagCLI.run('upload react-native', options, target), // Default ReactNative command
            {
                iOS: (options: BugsnagUploadiOSOptions = {}, target = ''): Promise<string> =>
                    BugsnagCLI.run('upload react-native-ios', options, target),
            
                Android: (options: BugsnagUploadAndroidOptions = {}, target = ''): Promise<string> =>
                    BugsnagCLI.run('upload react-native-android', options, target),
            }
        ),
        Js: (options: BugsnagUploadJsOptions = {}, target = ''): Promise<string> =>
            BugsnagCLI.run('upload js', options, target),
    }

    /**
     * Send build information to Bugsnag
     */
    static CreateBuild(options: BugsnagCreateBuildOptions = {}, target = ''): Promise<string> {
        return new Promise((resolve, reject) => {
            try {
                const output = BugsnagCLI.run('create-build', options, target)
                if (output instanceof Promise) {
                    output.then(resolve).catch(reject)
                } else {
                    resolve(output)
                }
            } catch (error) {
                reject(error)
            }
        })
    }

}

export default BugsnagCLI
