import axios from '../../../axios'

export const GetPersonal = () => {
    const url = '/api/ca/get_certificate_personal'
    return axios({
        url: url,
        method: 'get',
    })
}

export const CreatePersonal = () => {
    const url = '/api/ca/create_certificate_personal'
    return axios({
        url: url,
        method: 'post',
    })
}

export const DownloadCA = () => {
    const url = '/api/ca/download_certificate_personal'
    return axios({
        url: url,
        method: 'get',
        responseType: 'blob',
        headers:{ 'Content-Type': 'application/json; application/octet-stream'},
    })
}