import axios from '../../../axios'

export const GetServerById = (id) => {
    const url = '/api/ca/get_certificate_server'
    return axios({
        url: url,
        method: 'get',
        id
    })
}

export const CreateServer = (data) => {
    const url = '/api/ca/create_certificate_server'
    return axios({
        url: url,
        method: 'post',
        data
    })
}

export const GetServerList = (data) => {
    let url = '/api/ca/server_list?'
    for (var key in data) {
        url += key + '=' + data[key] + '&'
    }
    return axios({
        url: url,
        method: 'get',
    })
}

export const DownloadCA = (data) => {
    const url = '/api/ca/download_certificate_server'
    return axios({
        url: url,
        method: 'post',
        data,
        responseType: 'blob',
        headers:{ 'Content-Type': 'application/json; application/octet-stream'},
    })
}