const state = {
    userInfo: {}
}

const mutations = {
    SET_USER: (state, data) => {
        state.userInfo = data
    }
}

const actions = {
    setUser({ commit }, data) {
        commit("SET_USER", data)
    }
}

export default {
    namespaced: true,
    state,
    mutations,
    actions
}