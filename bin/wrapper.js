#!/usr/bin/env node
const { spawnSync } = require('child_process');
const path = require('path');
const os = require('os');
const fs = require('fs');

const ext = os.platform() === 'win32' ? '.exe' : '';
const localBin = path.join(__dirname, 'hostbind' + ext);

let commandToRun = null;

// If the postinstall script successfully downloaded the binary locally, use it
if (fs.existsSync(localBin)) {
    commandToRun = localBin;
} else {
    // Try finding 'hostbind' on the system PATH (installed via go install)
    const { execFileSync } = require('child_process');
    try {
        const which = os.platform() === 'win32' ? 'where' : 'which';
        execFileSync(which, ['hostbind'], { stdio: 'pipe' });
        commandToRun = 'hostbind';
    } catch (_) {
        // Not on PATH either
    }
}

if (!commandToRun) {
    console.error('❌ HostBind binary not found.');
    console.error('');
    console.error('The binary was not downloaded during install (no GitHub Release exists yet),');
    console.error('and hostbind was not found on your system PATH.');
    console.error('');
    console.error('To fix this, choose one option:');
    console.error('  1. Install Go (https://go.dev/dl/) then run:');
    console.error('       go install github.com/RudraMakwana257/hostbind/cmd/hostbind@latest');
    console.error('     Then ensure your GOPATH/bin is in your PATH.');
    console.error('  2. Download the binary from: https://github.com/RudraMakwana257/hostbind/releases');
    process.exit(1);
}

const args = process.argv.slice(2);
const result = spawnSync(commandToRun, args, { stdio: 'inherit' });

if (result.error) {
    console.error('❌ Failed to run HostBind:', result.error.message);
    process.exit(1);
}

process.exit(result.status || 0);
