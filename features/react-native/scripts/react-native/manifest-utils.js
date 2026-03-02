const fs = require('fs');
const path = require('path');

/**
 * Recursively find and remove the package attribute from all AndroidManifest.xml files
 * under the given directory.
 */
function removePackageAttributeFromManifests(rootDir) {
  const manifests = [];
  function findManifests(dir) {
    const entries = fs.readdirSync(dir, { withFileTypes: true });
    for (const entry of entries) {
      const fullPath = path.join(dir, entry.name);
      if (entry.isDirectory()) {
        findManifests(fullPath);
      } else if (entry.isFile() && entry.name === 'AndroidManifest.xml') {
        manifests.push(fullPath);
      }
    }
  }
  findManifests(rootDir);
  for (const manifestPath of manifests) {
    let contents = fs.readFileSync(manifestPath, 'utf8');
    // Remove package attribute from <manifest ...>
    contents = contents.replace(/<manifest\s+([^>]*?)package="[^"]*"([^>]*)>/, (match, p1, p2) => {
      // Remove just the package attribute
      let tag = `<manifest ${p1}${p2}>`;
      tag = tag.replace(/\s+package="[^"]*"/, '');
      return tag;
    });
    fs.writeFileSync(manifestPath, contents);
  }
}

module.exports = { removePackageAttributeFromManifests };
