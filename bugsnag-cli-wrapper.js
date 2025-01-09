const { exec } = require('child_process');

/**
 *
 * @typedef {Object} BugsnagOptions
 * @property {string|undefined} apiKey - The API key for authentication. Optional.
 * @property {boolean} dryRun - Whether to perform a dry run without actually uploading. Defaults to `false`.
 * @property {boolean} failOnUploadError - Whether to fail the process on upload error. Defaults to `false`.
 * @property {string|undefined} logLevel - The log level for output. Optional.
 * @property {number} port - The port for the connection. Defaults to `443`.
 * @property {boolean} verbose - Whether to enable verbose logging. Defaults to `false`.
 * @property {boolean} overwrite - Whether to overwrite existing files. Defaults to `false`.
 * @property {number} retries - The number of retries for failed requests. Defaults to `0`.
 * @property {number} timeout - The timeout value in seconds for requests. Defaults to `300`.
 * @property {string} uploadApiRootUrl - The root URL for the upload API. Defaults to `'https://upload.bugsnag.com'`.
 * @property {string|undefined} projectRoot - The root directory for the project. Optional.
 * @property {boolean} dev - Whether this is a development environment. Defaults to `false`.
 *
 */

/**
 * Wrapper for Bugsnag CLI
 */
class BugsnagCLI {
    /**
     * Convert camelCase to kebab-case
     * @param {string} str - The string in camelCase format.
     * @returns {string} - The string in kebab-case format.
     */
    static camelToKebab(str) {
        return str.replace(/([a-z0-9])([A-Z])/g, '$1-$2').toLowerCase();
    }

    /**
     * Execute a Bugsnag CLI command
     * @param {string} command - The main Bugsnag CLI command (e.g., "upload").
     * @param {Object} options - Key-value pairs of options for the CLI (e.g., { 'api-key': 'YOUR_API_KEY' }).
     * @param {string} target - Positional argument for the command (e.g., a file or folder path).
     * @returns {Promise<string>} - Resolves with the command's stdout or rejects with an error including stderr.
     */
    static run(command, options = {}, target = '') {
        return new Promise((resolve, reject) => {
            // Convert the options keys from camelCase to kebab-case
            const kebabCaseOptions = Object.entries(options)
                .map(([key, value]) => {
                    const kebabKey = BugsnagCLI.camelToKebab(key);
                    if (typeof value === 'boolean' && value === true) {
                        return `--${kebabKey}`;
                    } else if (typeof value !== 'boolean') {
                        return `--${kebabKey}=${value}`;
                    }
                    return '';
                })
                .filter(Boolean)
                .join(' ');

            const positionalArg = target ? `"${target}"` : '';
            const cliCommand = `npx bugsnag-cli ${command} ${kebabCaseOptions} ${positionalArg}`.trim();

            // Execute the command
            exec(cliCommand, (error, stdout, stderr) => {
                if (error) {
                    const errorMessage = `Command failed: ${cliCommand}\n` +
                        `Error: ${error.message}\n` +
                        `${stdout.trim()}`;
                    reject(errorMessage);
                } else {
                    resolve(stdout.trim());
                }
            });
        });
    }

    /**
     * Upload sourcemaps to Bugsnag
     * Provides nested methods for specific upload types.
     */
    static Upload = {
        ReactNative: Object.assign(
            /**
             *
             * @param {BugsnagOptions} options - Common Key-value pairs of options for the CLI.
             * @param {string|undefined} options.bundle - The bundle identifier. Optional.
             * @param {string|undefined} options.codeBundleId - The code bundle identifier. Optional.
             * @param {string|undefined} options.sourceMap - The source map file for the project. Optional.
             * @param {string|undefined} options.versionName - The version name for the project. Optional.
             * @param {string|undefined} options.androidAppManifest - The Android app manifest file. Optional.
             * @param {string|undefined} options.androidVariant - The variant for Android (e.g., production, staging). Optional.
             * @param {string|undefined} options.androidVersionCode - The version code for the Android project. Optional.
             * @param {string|undefined} options.iosBundleVersion - The iOS bundle version. Optional.
             * @param {string|undefined} options.iosPlist - The iOS plist file. Optional.
             * @param {string|undefined} options.iosScheme - The iOS build scheme. Optional.
             * @param {string|undefined} options.iosXcodeProject - The Xcode project file for iOS. Optional.
             * @param {string} target - The path to the file or directory to upload (e.g., a React Native bundle or folder).
             * @returns {Promise<string>} - Resolves with the command's output or rejects with an error message.
             */
            (options = {}, target = '') =>
                BugsnagCLI.run('upload react-native', options, target), // Default ReactNative command
            {
                /**
                 *
                 * @param {BugsnagOptions} options - Common Key-value pairs of options for the CLI.
                 * @param {string|undefined} options.bundle - The bundle identifier. Optional.
                 * @param {string|undefined} options.codeBundleId - The code bundle identifier. Optional.
                 * @param {string|undefined} options.sourceMap - The source map file for the project. Optional.
                 * @param {string|undefined} options.versionName - The version name for the project. Optional.
                 * @param {string|undefined} options.bundleVersion - The version of the bundle. Optional.
                 * @param {string|undefined} options.plist - The plist file for iOS projects. Optional.
                 * @param {string|undefined} options.scheme - The scheme for iOS builds. Optional.
                 * @param {string|undefined} options.xcodeProject - The Xcode project file for iOS builds. Optional.
                 * @param {string} target - The path to the file or directory to upload (e.g., a bundle file or folder).
                 * @returns {Promise<string>} - Resolves with the command's output or rejects with an error message.
                 */
                iOS: (options = {}, target = '') =>
                    BugsnagCLI.run('upload react-native-ios', options, target),
                /**
                 *
                 * @param {BugsnagOptions} options - Common Key-value pairs of options for the CLI.
                 * @param {string|undefined} options.bundle - The bundle identifier. Optional.
                 * @param {string|undefined} options.codeBundleId - The code bundle identifier. Optional.
                 * @param {string|undefined} options.sourceMap - The source map file for the project. Optional.
                 * @param {string|undefined} options.versionName - The version name for the project. Optional.
                 * @param {string|undefined} options.appManifest - The manifest file for the app. Optional.
                 * @param {string|undefined} options.variant - The variant of the project (e.g., production, staging). Optional.
                 * @param {string|undefined} options.versionCode - The version code for the project. Optional.
                 * @param {string} target - The path to the file or directory to upload (e.g., a bundle file or folder).
                 * @returns {Promise<string>}
                 */
                Android: (options = {}, target = '') =>
                    BugsnagCLI.run('upload react-native-android', options, target),
            }
        ),
        /**
         *
         * @param {BugsnagOptions} options - Common Key-value pairs of options for the CLI.
         * @param {string|undefined} options.baseUrl - The base URL for the project. Optional.
         * @param {string|undefined} options.bundle - The bundle identifier. Optional.
         * @param {string|undefined} options.bundleUrl - The URL of the bundle. Optional.
         * @param {string|undefined} options.projectRoot - The root directory for the project. Optional.
         * @param {string|undefined} options.sourceMap - The source map file for the project. Optional.
         * @param {string|undefined} options.versionName - The version name for the project. Optional.
         * @param {string|undefined} options.codeBundleId - The code bundle identifier. Optional.
         * @param {string} target - The path to the file or directory to upload (e.g., a JavaScript bundle or folder).
         * @returns {Promise<string>} - Resolves with the command's output or rejects with an error message.
         */
        Js: (options = {}, target = '') =>
            BugsnagCLI.run('upload js', options, target),
    };
}

module.exports = BugsnagCLI;
