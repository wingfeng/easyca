const checkDomain = (rule, value, cb) => {
    const reg = /^[*A-Za-z0-9]+\.[A-Za-z0-9]+[\/=\?%\-&_~`@[\]\':+!]*([^<>\"\"])*$/

    if (reg.test(value)) {
        return cb()
    }

    cb(new Error('请输入正确的域名格式'))
}

export let domain = [
    { required: true, message: '请输入域名', trigger: 'blur' },
    { validator: checkDomain, trigger: 'blur' }
]