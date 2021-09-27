const {execSync} = require("child_process")
const fs = require('fs');
const path = require('path');

function compress(distFolder, readPath, fileName, outPath) {
    execSync(`tar -C ./${distFolder} -zcf ${outPath}.tar.gz ${fileName}`)
    fs.unlinkSync(readPath)
}

function cleanDir(dirPath) {
    for (const file of fs.readdirSync(dirPath)) {
        fs.unlinkSync(path.join(dirPath, file))
    }
}

const OS = {
    LINUX: "linux",
    WINDOWS: "windows",
    DARWIN: "darwin"
}

const ARCH = {
    386: "386",
    AMD64: "amd64"
}

const BUILDS = [
    [OS.WINDOWS, ARCH.AMD64],
    [OS.WINDOWS, ARCH["386"]],
    [OS.LINUX, ARCH["386"]],
    [OS.LINUX, ARCH.AMD64],
    [OS.DARWIN, ARCH.AMD64],
]

const DIST_PATH = "dist";

function release(tag) {
    console.log(`releasing ${tag}`)

    cleanDir(DIST_PATH)

    for (let [os, arch] of BUILDS) {
        console.log(`processing ${os} ${arch}`)
        try {
            execSync(`GOOS=${os} GOARCH=${arch} go build -o ./dist ./cmd/tracker`)
            const fileName = os === "windows" ? "tracker.exe" : "tracker"
            const filePath = path.join(DIST_PATH, fileName)
            const distName = `tracker-${tag}-${os}-${arch}`
            const outPath = path.join(DIST_PATH, distName)
            compress(DIST_PATH, filePath, fileName, outPath)
        } catch (e) {
            console.error({os, arch})
            console.error(e)
            throw e
        }
    }
}

const tag = process.argv[2]
if (!tag) {
    throw new Error("tag is missing")
}

release(tag)
