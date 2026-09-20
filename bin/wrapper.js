#!/usr/bin/env node
const { spawnSync } = require('child_process');
const path = require('path');
const os = require('os');
const fs = require('fs');

const ext = os.platform() === 'win32' ? '.exe' : '';
const localBin = path.join(__dirname, 'hostbind' + ext);

let commandToRun = 'hostbind'; // Default to global PATH lookup

// If the postinstall script successfully downloaded the binary locally, use it
if (fs.existsSync(localBin)) {
    commandToRun = localBin;
}

const args = process.argv.slice(2);
const result = spawnSync(commandToRun, args, { stdio: 'inherit' });

if (result.error) {
    if (result.error.code === 'ENOENT') {
        console.error("❌ HostBind binary not found.");
        console.error("Please run 'npm install -g hostbind' again or download it manually.");
    } else {
        console.error(result.error);
    }
    process.exit(1);
}

process.exit(result.status || 0);
