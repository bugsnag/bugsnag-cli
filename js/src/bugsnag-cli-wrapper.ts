import { execFile } from 'child_process'
import { BugsnagCreateBuildOptions, BugsnagUploadiOSOptions, BugsnagUploadJsOptions, BugsnagUploadAndroidOptions, BugsnagUploadReactNativeOptions } from './types'
import * as path from "path"

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
            const kebabCaseOptions: string[] = Object.entries(options)
                .map(([key, value]) => {
                    const kebabKey = BugsnagCLI.camelToKebab(key)
                    if (typeof value === 'boolean' && value === true) {
                        return [`--${kebabKey}`]
                    } else if (typeof value !== 'boolean') {
                        return [`--${kebabKey}`, String(value)]
                    }
                    return []
                })
                .flat()

            const binPath = path.resolve(__dirname, path.join('..', 'bin', 'bugsnag-cli'))

            // Prepare CLI arguments
            const args = [
                ...command.split(' ').filter(Boolean),
                ...kebabCaseOptions,
                ...(target.trim() ? [target.trim()] : [])
            ]

            // Debug log
            console.log(args)

            // Execute the command
            execFile(binPath, args, (error, stdout, stderr) => {
                if (error) {
                    const errorMessage = `Command failed: ${binPath}\n` +
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

export = BugsnagCLI
