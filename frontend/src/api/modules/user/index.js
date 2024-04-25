import axios from '../../axios'

export const UserInfo = () => {
    const url = '/api/system/current'
    return axios({
        url: url,
        method: 'get'
    })
}

export const GetAllUser = () => {
    const url = '/api/casbin/get_all_user'
    return axios({
        url: url,
        method: 'get'
    })
}

export const GetRoleByUser = (data) => {
    const url = '/api/casbin/get_role_by_user'
    // for (var key in data) {
    //     url += key + '=' + data[key] + '&'
    // }
    return axios({
        url: url,
        method: 'post',
        data
    })
}

export const ModifyRole = (data) => {
    const url = '/api/casbin/modify_role'
    return axios({
        url: url,
        method: 'post',
        data
    })
}

export const getRoleList = () => {
    const url = '/api/casbin/role_list'
    return axios({
        url: url,
        method: 'get',
    })
}

export const handleUserAble = (data) => {
    var url = '/api/casbin/is_enable?'
    for (var key in data) {
        url += key + '=' + data[key] + '&'
    }
    return axios({
        url,
        method: 'get'
    })
}