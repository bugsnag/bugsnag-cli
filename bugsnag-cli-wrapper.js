const { exec } = require('child_process');

/**
 * Wrapper for Bugsnag CLI
 */
class BugsnagCli {
    /**
     * Execute a Bugsnag CLI command
     * @param {string} command - The main Bugsnag CLI command (e.g., "upload").
     * @param {Object} options - Key-value pairs of options for the CLI (e.g., { 'api-key': 'YOUR_API_KEY' }).
     * @param {string} target - Positional argument for the command (e.g., a file or folder path).
     * @returns {Promise<string>} - Resolves with the command's stdout or rejects with an error including stderr.
     */
    static run(command, options = {}, target = '') {
        return new Promise((resolve, reject) => {
            // Build the command string
            const flagArgs = Object.entries(options)
                .map(([key, value]) => {
                    if (typeof value === 'boolean' && value === true) {
                        return `--${key}`;
                    } else if (typeof value !== 'boolean') {
                        return `--${key}=${value}`;
                    }
                    return '';
                })
                .filter(Boolean)
                .join(' ');

            const positionalArg = target ? `"${target}"` : '';
            const cliCommand = `npx bugsnag-cli ${command} ${flagArgs} ${positionalArg}`.trim();

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
     * @param {string} type - The type of upload (e.g., "js", "xcode-archive").
     * @param {Object} options - Key-value pairs of options for the upload command.
     * @param {string} target - File or folder path to upload.
     * @returns {Promise<string>} - Resolves with the command's stdout.
     */
    static Upload(type, options = {}, target = '') {
        const command = `upload ${type}`;
        return BugsnagCli.run(command, options, target);
    }
}

module.exports = BugsnagCli;
