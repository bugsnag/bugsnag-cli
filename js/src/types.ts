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