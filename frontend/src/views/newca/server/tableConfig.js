export const tableConfig = {
  data() {
    return {
      tableColumns: [
        { prop: 'Dn', label: '证书名称', align: 'center', minWidth: 100, maxWidth: 220},
        { prop: 'UserId', label: '创建者', align: 'center', minWidth: 100, maxWidth: 220},
      ],
      operates: {
        list: [
          {
            label: '编辑',
            show: true,
            type: 'primary',
            method: row => {
              this.edit(row)
            }
          },
          {
            label: '删除',
            show: true,
            type: 'danger',
            method: row => {
              this.deleteData([row.id])
            }
          }
        ],
        width: 225
      },
      dataQuery: {}
    }
  }
}
