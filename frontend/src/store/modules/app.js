const state = {
    device: 'desktop',
    size: 'small',
    treeWidth: '18%',
    minMainHeight: window.screen.henght - 465,
    pageTableHeaderHeight: 0
}

const mutations = {
    PAGE_TABLE_HEADER_HEIGHT: (state, height) => {
        state.pageTableHeaderHeight = height
    },
    TOGGLE_DEVICE: (state, device) => {
        state.device = device
    }
}

const actions = {
    pageTableHeaderHeight({ commit }, { height}) {
        commit("PAGE_TABLE_HEADER_HEIGHT", height)
    },

    toggleDevice({ commit }, device) {
        commit("TOGGLE_DEVICE", device)
    }
}

export default {
    namespaced: true,
    state,
    mutations,
    actions
}