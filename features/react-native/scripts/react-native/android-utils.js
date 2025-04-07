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
  },
  configureReactNavigationAndroid: function configureReactNavigationAndroid (fixtureDir, reactNativeVersion) {
    const fileExtension = parseFloat(reactNativeVersion) < 0.73 ? 'java' : 'kt'
    let mainActivityPattern, mainActivityReplacement
    if (fileExtension === 'java') {
      mainActivityPattern = 'public class MainActivity extends ReactActivity {'
      mainActivityReplacement = `
  import android.os.Bundle;
  
  public class MainActivity extends ReactActivity {
  
    /**
     * Required for react-navigation/native implementation
     * https://reactnavigation.org/docs/getting-started/#installing-dependencies-into-a-bare-react-native-project
     */
    @Override
    protected void onCreate(Bundle savedInstanceState) {
      super.onCreate(null);
    }
  `
    } else if (fileExtension === 'kt') {
      mainActivityPattern = 'class MainActivity : ReactActivity() {'
      mainActivityReplacement = `
  import android.os.Bundle
  
  class MainActivity : ReactActivity() {
  
    /**
     * Required for react-navigation/native implementation
     * https://reactnavigation.org/docs/getting-started/#installing-dependencies-into-a-bare-react-native-project
     */
    override fun onCreate(savedInstanceState: Bundle?) {
      super.onCreate(null)
    }
  `
    }

    const mainActivityPath = `${fixtureDir}/android/app/src/main/java/com/reactnative/MainActivity.${fileExtension}`
    let mainActivityContents = fs.readFileSync(mainActivityPath, 'utf8')
    mainActivityContents = mainActivityContents.replace(mainActivityPattern, mainActivityReplacement)
    fs.writeFileSync(mainActivityPath, mainActivityContents)
  },
  buildAPK: function buildAPK (fixtureDir) {
    execFileSync('./gradlew', ['assembleRelease'], { cwd: `${fixtureDir}/android`, stdio: 'inherit' })
    fs.copyFileSync(`${fixtureDir}/android/app/build/outputs/apk/release/app-release.apk`, `${fixtureDir}/reactnative.apk`)
  }
}
