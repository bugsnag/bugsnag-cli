const { execFileSync } = require('child_process')
const fs = require('fs')

module.exports = {
  configureAndroidProject: function configureAndroidProject (fixtureDir, newArchEnabled) {
    // set android:usesCleartextTraffic="true" in AndroidManifest.xml
    const androidManifestPath = `${fixtureDir}/android/app/src/main/AndroidManifest.xml`
    let androidManifestContents = fs.readFileSync(androidManifestPath, 'utf8')

// Ensure `android:usesCleartextTraffic="true"` is added to the <application> tag
    if (!androidManifestContents.includes('android:usesCleartextTraffic="true"')) {
      androidManifestContents = androidManifestContents.replace(
          '<application',
          '<application android:usesCleartextTraffic="true"'
      )
    }

// Check if the Bugsnag <meta-data> tag already exists
    if (!androidManifestContents.includes('com.bugsnag.android.API_KEY')) {
      androidManifestContents = androidManifestContents.replace(
          /<application[^>]*>/,
          match => `${match}\n    <meta-data android:name="com.bugsnag.android.API_KEY" android:value="1234567890ABCDEF1234567890ABCDEF"/>`
      )
    }

    fs.writeFileSync(androidManifestPath, androidManifestContents)

    // enable/disable the new architecture in gradle.properties
    const gradlePropertiesPath = `${fixtureDir}/android/gradle.properties`
    let gradlePropertiesContents = fs.readFileSync(gradlePropertiesPath, 'utf8')
    gradlePropertiesContents = gradlePropertiesContents.replace(/newArchEnabled\s*=\s*(true|false)/, `newArchEnabled=${newArchEnabled}`)
    fs.writeFileSync(gradlePropertiesPath, gradlePropertiesContents)

    const buildGradlePath = `${fixtureDir}/android/app/build.gradle`
    let buildGradleContents = fs.readFileSync(buildGradlePath, 'utf8')
    buildGradleContents = buildGradleContents.replace(/def\s*enableProguardInReleaseBuilds\s*=\s*false/, 'def enableProguardInReleaseBuilds = true')
    fs.writeFileSync(buildGradlePath, buildGradleContents)
  },
  buildAAB: function buildAAB (fixtureDir) {
    execFileSync('./gradlew', ['bundleRelease'], { cwd: `${fixtureDir}/android`, stdio: 'inherit' })
  }
}
