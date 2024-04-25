// import * as user from '@/api/modules/user'
// import * as root from '@/api/modules/newca/root'
// import * as personal from "@/api/modules/newca/personal"
// import * as server from "@/api/modules/newca/server"
// import * as sign from "@/api/modules/newca/sign"
// import * as system from "@/api/modules/system"
const modulesFiles = require.context("./modules", true, /\.js$/)
const modules = {}

// recycle and alias the per moudle
for (const key of modulesFiles.keys()) {
    // split the module path
    const keyArr = key.split('/')
    // alias the module
    const keys = keyArr[keyArr.length - 2]
    // set the module
    modules[keys] = modulesFiles(key)
}

export default modules