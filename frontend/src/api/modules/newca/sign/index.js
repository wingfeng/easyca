import axios from '../../../axios'

export const GetSignById = (id) => {
    const url = '/api/ca/get_certificate_sign'
    return axios({
        url: url,
        method: 'get',
        id
    })
}

export const CreateSign = (data) => {
    const url = '/api/ca/create_certificate_sign'
    return axios({
        url: url,
        method: 'post',
        data
    })
}

export const GetSignList = (data) => {
    let url = '/api/ca/sign_list?'
    for (var key in data) {
        url += key + '=' + data[key] + '&'
    }
    return axios({
        url: url,
        method: 'get'
    })
}

export const DownloadCA = (data) => {
    const url = '/api/ca/download_certificate_sign'
    return axios({
        url: url,
        method: 'post',
        data,
        responseType: 'blob',
        headers:{ 'Content-Type': 'application/json; application/octet-stream'},
    })
}