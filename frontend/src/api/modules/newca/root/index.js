import axios from '../../../axios'

export const RootInfo = () => {
    const url = '/api/ca/get_base_cert'
    return axios({
        url: url,
        method: 'get',
    })
}

export const RootCreate = (data) => {
    const url = '/api/ca/create_ca'
    return axios({
        url: url,
        method: 'post',
        data
    })
}

export const DownloadCA = () => {
    const url = '/api/ca/download_base_cert'
    return axios({
        url: url,
        method: 'get',
        // responseType: 'blob',
        // headers:{ 'Content-Type': 'application/json; application/octet-stream'},
    })
}