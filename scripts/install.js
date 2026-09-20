const fs = require('fs');
const os = require('os');
const path = require('path');
const https = require('https');
const child_process = require('child_process');

const VERSION = 'v1.0.0'; // When you create a GitHub release, use this exact tag.
const REPO = 'RudraMakwana257/hostbind';

const osMap = {
    win32: 'windows',
    darwin: 'darwin',
    linux: 'linux'
};

const archMap = {
    x64: 'amd64',
    arm64: 'arm64'
};

const goOs = osMap[os.platform()];
const goArch = archMap[os.arch()];

if (!goOs || !goArch) {
    console.error(`HostBind: Unsupported platform (${os.platform()} ${os.arch()})`);
    process.exit(0); // Exit 0 so npm install doesn't crash completely
}

let binaryName = `hostbind-${goOs}-${goArch}`;
let ext = goOs === 'windows' ? '.exe' : '';
binaryName += ext;

const url = `https://github.com/${REPO}/releases/download/${VERSION}/${binaryName}`;
const binDir = path.join(__dirname, '..', 'bin');
const dest = path.join(binDir, `hostbind${ext}`);

if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir);
}

console.log(`Downloading HostBind ${VERSION} for ${goOs}-${goArch}...`);

function fallbackToGo() {
    console.log("Release not found (or download failed). Attempting fallback to 'go install'...");
    try {
        child_process.execSync('go install github.com/RudraMakwana257/hostbind/cmd/hostbind@latest', { stdio: 'inherit' });
        console.log("✅ Successfully installed via 'go install'. Ensure your GOPATH/bin is in your PATH.");
    } catch (e) {
        console.error("❌ Fallback failed. Please install Go or download the binary manually from GitHub Releases.");
    }
}

function download(downloadUrl) {
    https.get(downloadUrl, (res) => {
        if (res.statusCode === 301 || res.statusCode === 302) {
            download(res.headers.location);
        } else if (res.statusCode !== 200) {
            fallbackToGo();
        } else {
            const file = fs.createWriteStream(dest);
            res.pipe(file);
            file.on('finish', () => {
                file.close();
                if (goOs !== 'windows') {
                    fs.chmodSync(dest, 0o755);
                }
                console.log('✅ HostBind binary installed successfully!');
            });
        }
    }).on('error', (err) => {
        console.error('Download error:', err.message);
        fallbackToGo();
    });
}

download(url);
