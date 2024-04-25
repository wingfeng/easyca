const getters = {
    user: state => state.user,
    device: state => state.app.device,
    size: state => state.app.size,
    treeWidth: state => state.app.treeWidth,
    minMainHeight: state => state.app.minMainHeight,
    pageTableHeaderHeight: state => state.app.pageTableHeaderHeight
}

export default getters;