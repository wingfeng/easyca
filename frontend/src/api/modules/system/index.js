import axios from '../../axios'

export const SystemVersion = () => {
    const url = '/api/system/version'
    return axios({
        url: url,
        method: 'get'
    })
}