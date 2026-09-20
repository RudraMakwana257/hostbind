const fs = require('fs');
const os = require('os');
const path = require('path');
const https = require('https');
const child_process = require('child_process');

// Read version dynamically from package.json to stay in sync
const pkg = require('../package.json');
const VERSION = `v${pkg.version}`;
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
    console.error('Supported: linux/darwin/windows on amd64/arm64');
    process.exit(0); // Exit 0 so npm install doesn't crash completely
}

let binaryName = `hostbind-${goOs}-${goArch}`;
const ext = goOs === 'windows' ? '.exe' : '';
binaryName += ext;

const url = `https://github.com/${REPO}/releases/download/${VERSION}/${binaryName}`;
const binDir = path.join(__dirname, '..', 'bin');
const dest = path.join(binDir, `hostbind${ext}`);

if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
}

console.log(`HostBind: Attempting to download ${VERSION} for ${goOs}-${goArch}...`);

function showManualInstructions() {
    console.error('');
    console.error('─────────────────────────────────────────────────────────');
    console.error('  HostBind binary could not be downloaded automatically.');
    console.error('');
    console.error('  To finish setup, choose one option:');
    console.error('');
    console.error('  Option 1 — Install Go (https://go.dev/dl/) then run:');
    console.error('    go install github.com/RudraMakwana257/hostbind/cmd/hostbind@latest');
    console.error('    (then ensure your GOPATH/bin is in your PATH)');
    console.error('');
    console.error('  Option 2 — Download the binary manually:');
    console.error(`    https://github.com/${REPO}/releases`);
    console.error(`    Place the binary as: ${dest}`);
    console.error('─────────────────────────────────────────────────────────');
}

function fallbackToGo() {
    console.log("HostBind: Release not found. Attempting fallback via 'go install'...");
    try {
        child_process.execSync(
            'go install github.com/RudraMakwana257/hostbind/cmd/hostbind@latest',
            { stdio: 'inherit' }
        );
        console.log("✅ Successfully installed via 'go install'. Ensure your GOPATH/bin is in your PATH.");
    } catch (e) {
        console.error("❌ 'go install' fallback also failed (Go may not be installed).");
        showManualInstructions();
    }
}

function download(downloadUrl) {
    https.get(downloadUrl, (res) => {
        if (res.statusCode === 301 || res.statusCode === 302) {
            download(res.headers.location);
        } else if (res.statusCode === 404) {
            console.warn(`HostBind: No release found at ${downloadUrl}`);
            fallbackToGo();
        } else if (res.statusCode !== 200) {
            console.warn(`HostBind: Unexpected HTTP ${res.statusCode} from GitHub.`);
            fallbackToGo();
        } else {
            const file = fs.createWriteStream(dest);
            res.pipe(file);
            file.on('finish', () => {
                file.close();
                if (goOs !== 'windows') {
                    fs.chmodSync(dest, 0o755);
                }
                console.log(`✅ HostBind ${VERSION} installed successfully!`);
            });
            file.on('error', (err) => {
                console.error('HostBind: Failed to write binary:', err.message);
                fallbackToGo();
            });
        }
    }).on('error', (err) => {
        console.error('HostBind: Network error during download:', err.message);
        fallbackToGo();
    });
}

download(url);
