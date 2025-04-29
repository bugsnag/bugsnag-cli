declare module '@bugsnag/cli' {
    interface BugsnagCreateBuildOptions {
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

    interface BugsnagUploadReactNativeOptions {
        bundle?: string
        codeBundleId?: string
        sourceMap?: string
        versionName?: string
        androidAppManifest?: string
        androidVariant?: string
        androidVersionCode?: string
        iosBundleVersion?: string
        iosPlist?: string
        iosScheme?: string
        iosXcodeProject?: string
    }

    interface BugsnagUploadiOSOptions {
        bundle?: string
        codeBundleId?: string
        sourceMap?: string
        versionName?: string
        bundleVersion?: string
        plist?: string
        scheme?: string
        xcodeProject?: string
    }

    interface BugsnagUploadAndroidOptions {
        bundle?: string
        codeBundleId?: string
        sourceMap?: string
        versionName?: string
        appManifest?: string
        variant?: string
        versionCode?: string
    }

    interface BugsnagUploadJsOptions {
        baseUrl?: string
        bundle?: string
        bundleUrl?: string
        projectRoot?: string
        sourceMap?: string
        versionName?: string
        codeBundleId?: string
    }

    /**
     * Convert camelCase to kebab-case
     */
    export function camelToKebab (str: string): string;

    /**
     * Execute a Bugsnag CLI command
     */
    export function run (command: string, options?: object, target?: string): Promise<string>;

    /**
     * Send build information to Bugsnag
     */
    export function CreateBuild (options: BugsnagCreateBuildOptions, target?: string): Promise<string>;

    /**
     * Upload sourcemaps to Bugsnag
     * Provides nested methods for specific upload types.
     */
    export const Upload: {
        ReactNative: ((options: BugsnagUploadReactNativeOptions, target?: string) => Promise<string>) & {
            Android: (options: BugsnagUploadAndroidOptions, target?: string) => Promise<string>,
            iOS: (options: BugsnagUploadiOSOptions, target?: string) => Promise<string>
        },
        Js: (options: BugsnagUploadJsOptions, target?: string) => Promise<string>,
    }
}
