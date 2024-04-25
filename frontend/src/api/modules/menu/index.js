import axios from '../../axios'

export const GetMenuList = () => {
    const url = '/api/system/menu'
    return axios({
        url: url,
        method: 'get'
    })
}
